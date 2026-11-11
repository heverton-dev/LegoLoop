package matriz

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"legoloop/internal/assets"
	"legoloop/internal/processo"
	"legoloop/internal/reputation"
)

// maxLeitura limita quanto de cada arquivo de peça é lido (a nota e o catálogo só precisam do começo).
const maxLeitura = 64 << 10

// Options controla a geração.
type Options struct {
	Vault  string // cofre onde ficam registry_matrix.json, matrix/ e o vocabulário
	Forcar bool   // reprocessa todas as peças, mesmo sem mudança
}

// Deps são as bordas injetáveis da geração.
type Deps struct {
	Scan      func(ctx context.Context) ([]assets.Item, error)
	Catalogar Catalogador
	ReadFile  func(path string) ([]byte, error)                      // nil = lê os primeiros 64 KiB do disco
	Score     func(ref, name, content string) reputation.PieceRating // nil = reputation.ScorePiece sem histórico
	FS        processo.FS                                            // nil = disco real
	Home      string                                                 // para abreviar caminhos como "~/..."
	Now       func() time.Time
}

// Summary é o que a geração fez.
type Summary struct {
	Total          int            `json:"total"`
	Reprocessadas  int            `json:"reprocessadas"`  // peças novas, alteradas ou à espera do cérebro
	Reaproveitadas int            `json:"reaproveitadas"` // peças iguais à geração anterior (hash igual): nem nota nem cérebro
	Copias         int            `json:"copias"`         // cópias idênticas de uma peça já catalogada (outro caminho): perfil compartilhado, sem cérebro
	PorTipo        map[string]int `json:"por_tipo"`
	Pendentes      int            `json:"pendentes"`                        // propostas novas de termo, à espera do usuário
	Alertas        []string       `json:"alertas,omitempty"`                // textos de peça que parecem ordem, citados como dado
	Erros          []string       `json:"erros,omitempty"`                  // peças que o cérebro não conseguiu catalogar
	Invalidadas    []string       `json:"confirmacoes_invalidas,omitempty"` // confirmações do usuário que perderam a validade (conteúdo mudou)
	Arquivos       []string       `json:"arquivos"`                         // arquivos gravados, relativos ao cofre
}

var kindProcesso = tipoDeProcesso()

func tipoDeProcesso() assets.Kind {
	conhecidos := map[assets.Kind]bool{
		assets.KindSkill: true, assets.KindMCP: true, assets.KindScript: true, assets.KindHook: true,
		assets.KindGate: true, assets.KindHarness: true, assets.KindBinary: true,
	}
	for _, k := range assets.Kinds() {
		if !conhecidos[k] {
			return k
		}
	}
	return ""
}

var tiposCatalogados = map[assets.Kind]bool{
	assets.KindSkill: true, assets.KindMCP: true, assets.KindScript: true,
	assets.KindHook: true, assets.KindGate: true, kindProcesso: true,
}

// Generate varre as peças, pontua, cataloga e grava a matriz geral e a matriz por tipo numa só operação atômica.
// É incremental: peça com o mesmo hash da geração anterior é reaproveitada, sem nota nem chamada ao cérebro.
func Generate(ctx context.Context, opts Options, d Deps) (*Summary, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	items, err := d.Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("matriz: a varredura falhou: %w", err)
	}
	g := &gerador{opts: opts, d: d, sum: &Summary{PorTipo: map[string]int{}}}
	if g.conf, err = lerConfirmados(opts.Vault); err != nil {
		return nil, err
	}
	g.carregarAnterior()
	entradas, err := g.catalogarTudo(ctx, filtrar(items, d))
	if err != nil {
		return nil, err
	}
	g.sum.Invalidadas = invalidadas(g.conf, entradas)
	files, err := arquivos(opts.Vault, entradas, g.novas, d, g.sum)
	if err != nil {
		return nil, err
	}
	st := processo.Store{Roots: processo.Roots{Vault: opts.Vault}, FS: d.FS}
	if err := st.Write(ctx, files, processo.WriteOpts{Force: true}); err != nil {
		return nil, err
	}
	for _, f := range files {
		rel, _ := filepath.Rel(opts.Vault, f.Path)
		g.sum.Arquivos = append(g.sum.Arquivos, filepath.ToSlash(rel))
	}
	return g.sum, nil
}

// gerador guarda o estado de uma execução: a matriz anterior (por Ref), as propostas novas e o resumo.
type gerador struct {
	opts    Options
	d       Deps
	prev    map[string]Entrada
	porHash map[string]Entrada     // perfil já catalogado por (tipo, hash): cópias idênticas não chamam o cérebro de novo
	conf    map[string]confirmacao // perfis que o usuário confirmou ou definiu, por (tipo, hash)
	novas   []Pendente
	sum     *Summary
}

// carregarAnterior lê a matriz da geração anterior. Sem arquivo ou na v1 (sem perfil nem hash) tudo é novo;
// arquivo ilegível também, com um aviso no resumo.
func (g *gerador) carregarAnterior() {
	g.prev, g.porHash = map[string]Entrada{}, map[string]Entrada{}
	reg, err := LerRegistro(filepath.Join(g.opts.Vault, "registry_matrix.json"))
	if errors.Is(err, os.ErrNotExist) {
		return
	}
	if err != nil {
		g.sum.Erros = append(g.sum.Erros, "a matriz anterior estava ilegível; gerando do zero")
		return
	}
	for _, e := range reg.Entradas {
		g.prev[e.Ref] = e
		if !e.Retentar && !g.opts.Forcar {
			g.porHash[chaveHash(e.Tipo, e.Hash)] = e
		}
	}
}

// catalogarTudo cataloga cada peça. A falha do cérebro numa peça não aborta: a peça fica sem perfil, marcada
// para nova tentativa, e o erro vai ao resumo.
func (g *gerador) catalogarTudo(ctx context.Context, items []assets.Item) ([]Entrada, error) {
	var entradas []Entrada
	for _, it := range items {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		e, err := g.item(ctx, it)
		if err != nil {
			g.sum.Erros = append(g.sum.Erros, it.Name+": "+err.Error())
		}
		for _, a := range e.Alertas {
			if l := it.Name + ": " + a; !slices.Contains(g.sum.Alertas, l) {
				g.sum.Alertas = append(g.sum.Alertas, l)
			}
		}
		entradas = append(entradas, e)
		g.sum.PorTipo[e.Tipo]++
	}
	g.sum.Total = len(entradas)
	return entradas, nil
}

// item devolve a entrada da peça, já com o perfil que o usuário confirmou (se houver) por cima de qualquer outro.
func (g *gerador) item(ctx context.Context, it assets.Item) (Entrada, error) {
	e, err := g.itemBase(ctx, it)
	if c, ok := g.conf[chaveHash(e.Tipo, e.Hash)]; ok {
		p := c.Perfil
		e.Perfil, e.Retentar = &p, false
	}
	return e, err
}

// itemBase devolve a entrada da peça: a anterior, se o hash não mudou, ou uma catalogação nova.
func (g *gerador) itemBase(ctx context.Context, it assets.Item) (Entrada, error) {
	r, hash, conteudo, lido := identidade(it, g.d)
	if prev, ok := g.prev[r]; ok && !g.opts.Forcar && prev.Hash == hash && !prev.Retentar {
		prev.Tipado = tipado(it)
		g.sum.Reaproveitadas++
		return prev, nil
	}
	if c, ok := g.porHash[chaveHash(string(it.Kind), hash)]; ok {
		g.sum.Copias++
		return copiaDe(c, it, r), nil
	}
	if leConteudo(it.Kind) && !lido {
		conteudo = lerArquivo(it.Path, g.d)
	}
	g.sum.Reprocessadas++
	e, err := g.catalogar(ctx, it, r, hash, conteudo)
	if !e.Retentar {
		g.porHash[chaveHash(e.Tipo, hash)] = e
	}
	return e, err
}

func chaveHash(tipo, hash string) string { return tipo + "|" + hash }

// copiaDe monta a entrada de uma cópia idêntica de uma peça já catalogada: mesmo perfil e mesma nota, Ref própria.
func copiaDe(c Entrada, it assets.Item, r string) Entrada {
	c.Ref, c.Tipado = r, tipado(it)
	c.Rating.Ref = r
	return c
}

func (g *gerador) catalogar(ctx context.Context, it assets.Item, r, hash, conteudo string) (Entrada, error) {
	e := Entrada{Ref: r, Nome: it.Name, Tipo: string(it.Kind), Hash: hash, Tipado: tipado(it)}
	score := g.d.Score
	if score == nil {
		score = func(ref, nome, texto string) reputation.PieceRating {
			return reputation.ScorePiece(ref, nome, texto, 0, 0)
		}
	}
	e.Rating = score(r, it.Name, conteudo)
	p := Peca{Ref: r, Nome: it.Name, Descricao: it.Description, Tipo: string(it.Kind)}
	if vaiAoCerebro(it.Kind) {
		p.Conteudo = conteudo
	}
	cat := g.d.Catalogar
	if _, ok := g.conf[chaveHash(e.Tipo, hash)]; ok {
		cat.Brain = nil // o usuário já decidiu o perfil: o cérebro não gasta modelo nesta peça
	}
	res, err := cat.Catalogar(ctx, p)
	e.Perfil, e.Alertas = res.Perfil, res.Alertas
	// sem cérebro, a peça sem casos de uso fica à espera dele; com falha do cérebro, tenta de novo
	e.Retentar = err != nil || (cat.Brain == nil && (e.Perfil == nil || len(e.Perfil.CasosDeUso) == 0))
	completaIO(e.Perfil, it)
	g.novas = append(g.novas, res.Pendentes...)
	return e, err
}

// filtrar fica só com os tipos catalogados e ordena por caminho: a saída não depende da ordem da varredura.
func filtrar(items []assets.Item, d Deps) []assets.Item {
	var out []assets.Item
	for _, it := range items {
		if tiposCatalogados[it.Kind] && it.Path != "" {
			out = append(out, it)
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return refDe(out[i], d.Home) < refDe(out[j], d.Home) })
	return out
}

// ref abrevia o caminho com "~/" quando está sob o HOME, como a matriz v1.
func ref(path, home string) string {
	if home != "" {
		if rel, err := filepath.Rel(home, path); err == nil && !strings.HasPrefix(rel, "..") {
			return "~/" + filepath.ToSlash(rel)
		}
	}
	return filepath.ToSlash(path)
}

// refDe é a Ref da peça: o caminho abreviado, mais "#nome" em MCP e hook, que dividem um mesmo arquivo de configuração.
func refDe(it assets.Item, home string) string {
	r := ref(it.Path, home)
	if it.Kind == assets.KindMCP || it.Kind == assets.KindHook {
		return r + "#" + it.Name
	}
	return r
}

// leConteudo diz se o texto do arquivo pode ser lido: MCP e hook apontam para o arquivo de configuração
// inteiro, que pode guardar segredo; esses nunca são lidos.
func leConteudo(k assets.Kind) bool { return k != assets.KindMCP && k != assets.KindHook }

// vaiAoCerebro diz se o texto vai ao catalogador (e portanto ao modelo): só skills e Processos.
func vaiAoCerebro(k assets.Kind) bool { return k == assets.KindSkill || k == kindProcesso }

// identidade devolve a Ref e o hash da peça. Só lê o arquivo quando o tipo permite e a varredura não trouxe o
// sha256; nesse caso devolve também o texto lido, para a catalogação não ler duas vezes.
func identidade(it assets.Item, d Deps) (r, hash, conteudo string, lido bool) {
	r = refDe(it, d.Home)
	if leConteudo(it.Kind) && it.Meta["sha256"] == "" {
		conteudo, lido = lerArquivo(it.Path, d), true
	}
	return r, hashDe(it, conteudo), conteudo, lido
}

func lerArquivo(path string, d Deps) string {
	if d.ReadFile != nil {
		b, _ := d.ReadFile(path)
		return string(b)
	}
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	b, _ := io.ReadAll(io.LimitReader(f, maxLeitura))
	return string(b)
}

// hashDe: o sha256 que a varredura já calculou; senão o do conteúdo lido; em MCP e hook, o da identidade.
func hashDe(it assets.Item, conteudo string) string {
	if h := it.Meta["sha256"]; h != "" {
		return h
	}
	texto := conteudo
	if !leConteudo(it.Kind) {
		chaves := make([]string, 0, len(it.Meta))
		for k := range it.Meta {
			chaves = append(chaves, k+"="+it.Meta[k])
		}
		slices.Sort(chaves)
		texto = string(it.Kind) + "|" + it.Name + "|" + it.Description + "|" + strings.Join(chaves, ",")
	}
	h := sha256.Sum256([]byte(texto))
	return hex.EncodeToString(h[:])
}

func completaIO(p *PerfilDeUso, it assets.Item) {
	if p == nil {
		return
	}
	if len(p.Consome) == 0 && it.Meta["consumes"] != "" {
		p.Consome = []string{it.Meta["consumes"]}
	}
	if len(p.Produz) == 0 && it.Meta["produces"] != "" {
		p.Produz = []string{it.Meta["produces"]}
	}
}

// camposPorTipo: o que a matriz por tipo guarda de cada tipo de peça (nada que possa ser segredo).
var camposPorTipo = map[assets.Kind][]string{
	assets.KindSkill: {"consumes", "produces", "artifact", "gate", "suggested_model"},
	assets.KindMCP:   {"transport", "command", "host"},
	assets.KindHook:  {"matcher"},
	kindProcesso:     {"trigger", "artifact"},
}

func tipado(it assets.Item) map[string]any {
	out := map[string]any{}
	for _, k := range camposPorTipo[it.Kind] {
		if v := it.Meta[k]; v != "" {
			out[k] = v
		}
	}
	if it.Kind == assets.KindScript {
		out["ext"] = strings.ToLower(filepath.Ext(it.Path))
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// arquivos monta, sem gravar, o registro geral, os arquivos por tipo e a fila de propostas de termo.
func arquivos(vault string, entradas []Entrada, novas []Pendente, d Deps, sum *Summary) ([]processo.File, error) {
	now := time.Now
	if d.Now != nil {
		now = d.Now
	}
	vocab := d.Catalogar.Vocab
	reg := Registro{GeneratedAt: now().UTC().Format("2006-01-02"), Versao: VersaoMatriz, VocabVersao: vocab.Versao, TotalPecas: len(entradas)}
	porTipo := map[string][]Entrada{}
	for _, e := range entradas {
		reg.Ratings = append(reg.Ratings, e.Rating)
		porTipo[e.Tipo] = append(porTipo[e.Tipo], e)
		geral := e
		geral.Tipado = nil
		reg.Entradas = append(reg.Entradas, geral)
	}
	files := []processo.File{json2(filepath.Join(vault, "registry_matrix.json"), "registro", reg)}
	for _, tipo := range sortedKeys(porTipo) {
		a := arquivoPorTipo{Versao: VersaoMatriz, Tipo: tipo, Total: len(porTipo[tipo]), Entradas: porTipo[tipo]}
		files = append(files, json2(filepath.Join(vault, "matrix", tipo+".json"), "tipo", a))
	}
	return comFila(vault, files, novas, sum)
}

// comFila junta as propostas novas à fila existente e acrescenta o arquivo da fila, se houve novidade.
func comFila(vault string, files []processo.File, novas []Pendente, sum *Summary) ([]processo.File, error) {
	if len(novas) == 0 {
		return files, nil
	}
	fila, err := LerPendentes(vault)
	if err != nil {
		return nil, err
	}
	antes := len(fila)
	for _, n := range novas {
		fila = mesclar(fila, n)
	}
	sum.Pendentes = len(fila) - antes
	return append(files, json2(filepath.Join(vault, pendenteFile), "pendente", fila)), nil
}

func json2(path, kind string, v any) processo.File {
	b, _ := json.MarshalIndent(v, "", "  ")
	return processo.File{Kind: kind, Path: path, Data: append(b, '\n')}
}

func sortedKeys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

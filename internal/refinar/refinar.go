// Package refinar muda uma etapa de um processo do cofre. Mostra o que muda em INTENT, SPEC, PLAN, ROUTE, DAG,
// loop e skill, pede a confirmação, guarda cada original como .legoloop-bak, sobe a versão e grava tudo de uma
// vez. É uma fatia independente: depende só do núcleo (processo, ste) e da matriz, nunca de outra função.
package refinar

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"legoloop/internal/caminhos"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"legoloop/internal/matriz"
	"legoloop/internal/processo"
	"legoloop/internal/ste"
	"legoloop/internal/studiodoc"
)

var (
	// ErrNome: o nome não é de processo (use legoloop-<nome>).
	ErrNome = errors.New("refinar: nome inválido")
	// ErrNaoExiste: o cofre não tem o processo da v2 com esse nome.
	ErrNaoExiste = errors.New("refinar: processo não encontrado")
	// ErrMudanca: a mudança pedida não vale para o processo.
	ErrMudanca = errors.New("refinar: mudança inválida")
	// ErrCancelado: o usuário não confirmou. Nenhum arquivo mudou.
	ErrCancelado = errors.New("refinar: cancelado")
)

var reNome = regexp.MustCompile(`^legoloop-[a-z0-9]+(-[a-z0-9]+)*$`)

const (
	// sufixoBak é o sufixo que processo.Store e canon usam para guardar o original. O Store cria o backup;
	// aqui o sufixo só serve para informar o caminho dele.
	sufixoBak  = ".legoloop-bak"
	maxLinhas  = 4   // linhas de cada lado que o diff mostra por arquivo
	maxLargura = 100 // largura de cada linha do diff
	formato    = "%-8s %-6s %-9s %s\n"
)

func init() {
	for id, msg := range map[string]string{
		"refinar.nome_invalido":     "O nome %s não é de processo. Use legoloop-<nome>, em letras minúsculas.",
		"refinar.nao_existe":        "O processo %s não existe no cofre %s. Veja os processos com legoloop listar.",
		"refinar.so_v1":             "O processo %s é da v1 e não tem o arquivo .processo.json. Crie o processo de novo com legoloop fabricar.",
		"refinar.sem_etapa":         "O processo %s não tem a etapa %s. As etapas são: %s.",
		"refinar.sem_mudanca":       "Diga o que mudar na etapa %s: --peca, --caso ou --titulo.",
		"refinar.sem_caso":          "A etapa %s precisa de ao menos um caso de uso. Informe com --caso.",
		"refinar.etapa_composta":    "A etapa %s compõe outro processo. Para mudar a composição, use legoloop compor.",
		"refinar.peca_desconhecida": "A peça %s não está na matriz do cofre. Veja as peças com legoloop matriz ver.",
		"refinar.cancelado":         "Operação cancelada. Nenhum arquivo mudou.",
		"refinar.concluido":         "Processo %s refinado: versão %s para %s.",
		"refinar.entrada_invalida":  "Entrada inválida: responda 1, 2 ou 0.",
		"refinar.diff_titulo":       "REFINAR %s: versão %s para %s.",
		"refinar.diff_etapa":        "Etapa %s (%s).",
		"refinar.diff_peca":         "Peça: %s para %s.",
		"refinar.diff_titulo_etapa": "Título: %s para %s.",
		"refinar.diff_casos":        "Casos: %s para %s.",
		"refinar.diff_colunas":      "ARQUIVO  ESTADO MUDANÇA   CAMINHO",
		"refinar.diff_mais":         "... mais %d linhas.",
		"refinar.diff_aviso":        "AVISO: %s",
		"refinar.cartao_acao":       "Refinar o processo %s",
		"refinar.cartao_aviso":      "O programa troca %d arquivos. Cada original fica guardado como .legoloop-bak.",
		"refinar.cartao_etapa":      "O programa muda a etapa %s: %s.",
		"refinar.cartao_versao":     "O programa sobe a versão de %s para %s e guarda %s no histórico.",
		"refinar.cartao_gera":       "O programa gera de novo INTENT, SPEC, PLAN, ROUTE, DAG, loop e skill.",
		"refinar.cartao_se_errado":  "Se um passo falhar, o programa desfaz tudo. O processo volta à versão %s. Para voltar depois, restaure cada arquivo .legoloop-bak.",
		"refinar.efeito_troca":      "troca %s",
		"refinar.efeito_novo":       "grava %s",
		"refinar.param_peca":        "Nome da peça que faz a etapa. Enter mantém o valor.",
		"refinar.param_casos":       "Casos de uso da etapa, separados por vírgula. Enter mantém o valor.",
	} {
		ste.Default.Register(id, msg)
	}
}

// Mudanca é a alteração de uma etapa. Peca, Titulo e Casos vazios mantêm o valor atual; ao menos um deles
// precisa trazer valor novo. Depois da mudança a etapa precisa ter ao menos um caso de uso.
type Mudanca struct {
	Etapa  string   `json:"etapa"`
	Peca   string   `json:"peca,omitempty"`
	Titulo string   `json:"titulo,omitempty"`
	Casos  []string `json:"casos,omitempty"`
}

// Request é o pedido: o processo e a mudança.
type Request struct {
	Nome    string
	Mudanca Mudanca
}

// ConfirmFunc mostra o plano e devolve true para gravar. replan recalcula o plano com outra mudança (o
// "personalizar"); o último plano pedido com replan é o que vale na gravação.
type ConfirmFunc func(p Plano, replan func(Mudanca) (Plano, error)) (bool, error)

// Deps são as bordas do refinar. Store.Roots traz o cofre e a raiz do projeto; Store.FS nil usa o disco real.
// Confirm nil segue sem perguntar. Clock nil usa time.Now.
type Deps struct {
	Store   processo.Store
	Confirm ConfirmFunc
	Clock   func() time.Time
}

// Arquivo é um dos arquivos que o refinar regera. Tipo é INTENT, SPEC, PLAN, ROUTE, DAG, loop, skill ou
// processo. Estado é troca, novo ou igual; só troca e novo são gravados.
type Arquivo struct {
	Tipo    string   `json:"tipo"`
	Caminho string   `json:"caminho"`
	Estado  string   `json:"estado"`
	Mais    int      `json:"mais"`   // linhas que entram
	Menos   int      `json:"menos"`  // linhas que saem
	Linhas  []string `json:"linhas"` // amostra do diff: "+ linha" e "- linha"
	Backup  string   `json:"backup,omitempty"`
}

// Plano é tudo que o refinar vai fazer, mostrado ao usuário antes de qualquer efeito.
type Plano struct {
	Nome     string         `json:"nome"`
	De       string         `json:"de"`
	Para     string         `json:"para"`
	Mudanca  Mudanca        `json:"mudanca"` // a mudança com os valores finais: serve de base para o replan
	Antes    processo.Etapa `json:"antes"`
	Depois   processo.Etapa `json:"depois"`
	Arquivos []Arquivo      `json:"arquivos"`
	Avisos   []string       `json:"avisos,omitempty"`
	files    []processo.File
}

// Resultado é o que o refinar gravou.
type Resultado struct {
	Nome     string    `json:"nome"`
	De       string    `json:"de"`
	Para     string    `json:"para"`
	Mudanca  Mudanca   `json:"mudanca"`
	Arquivos []Arquivo `json:"arquivos"`
	Avisos   []string  `json:"avisos,omitempty"`
}

// Run planeja a mudança, pede a confirmação e grava. Cancelar ou falhar deixa o cofre e o projeto como estavam.
func Run(ctx context.Context, req Request, d Deps) (*Resultado, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	atual, err := planejar(req, d)
	if err != nil {
		return nil, err
	}
	if d.Confirm != nil {
		replan := func(m Mudanca) (Plano, error) {
			p, err := planejar(Request{Nome: req.Nome, Mudanca: m}, d)
			if err == nil {
				atual = p
			}
			return p, err
		}
		ok, err := d.Confirm(atual, replan)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, ErrCancelado
		}
	}
	return gravar(ctx, atual, d)
}

// gravar troca os arquivos de uma vez só: o Store guarda os originais e desfaz tudo se um passo falhar.
func gravar(ctx context.Context, pl Plano, d Deps) (*Resultado, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := d.Store.Write(ctx, pl.files, processo.WriteOpts{Force: true}); err != nil {
		return nil, err
	}
	return &Resultado{Nome: pl.Nome, De: pl.De, Para: pl.Para, Mudanca: pl.Mudanca, Arquivos: pl.Arquivos, Avisos: pl.Avisos}, nil
}

func (d Deps) agora() time.Time {
	if d.Clock == nil {
		return time.Now().UTC()
	}
	return d.Clock().UTC()
}

// Planejar calcula a mudança na etapa e o diff dos arquivos sem tocar no disco.
func Planejar(req Request, d Deps) (Plano, error) {
	return planejar(req, d)
}

// planejar lê o processo, aplica a mudança, sobe a versão e regera os arquivos. Não grava nada.
func planejar(req Request, d Deps) (Plano, error) {
	roots := d.Store.Roots
	if roots.Vault == "" || roots.Docs == "" {
		return Plano{}, fmt.Errorf("%w: informe o cofre e a raiz do projeto", processo.ErrPlano)
	}
	if !reNome.MatchString(req.Nome) {
		return Plano{}, fmt.Errorf("%w: %s", ErrNome, ste.Msg("refinar.nome_invalido", req.Nome))
	}
	atual, err := lerProcesso(roots.Vault, req.Nome)
	if err != nil {
		return Plano{}, err
	}
	novo, antes, err := aplicar(atual, req.Mudanca)
	if err != nil {
		return Plano{}, err
	}
	if novo.Versao, err = subirMinor(atual.Versao); err != nil {
		return Plano{}, err
	}
	novo.Historico = append(slices.Clone(atual.Historico), atual.Versao)
	arqs, files, err := montar(novo, req.Mudanca, d.agora(), roots)
	if err != nil {
		return Plano{}, err
	}
	depois := novo.Etapas[indiceDe(novo.Etapas, req.Mudanca.Etapa)]
	return Plano{Nome: novo.Nome, De: atual.Versao, Para: novo.Versao, Antes: antes, Depois: depois, Arquivos: arqs,
		Mudanca: Mudanca{Etapa: depois.ID, Peca: depois.Peca, Titulo: depois.Titulo, Casos: depois.Casos},
		Avisos:  avisosDaPeca(roots.Vault, depois.Peca, antes.Peca), files: files}, nil
}

// lerProcesso lê loops/<nome>.processo.json. Processo da v1 (só o .loop.json) não se refina.
func lerProcesso(vault, nome string) (processo.Processo, error) {
	caminho := processo.CaminhoNoCofre(vault, nome)
	if _, err := os.Stat(caminho); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return processo.Processo{}, err
		}
		msg := ste.Msg("refinar.nao_existe", nome, vault)
		if _, err := os.Stat(filepath.Join(vault, caminhos.DirLoops, nome+".loop.json")); err == nil {
			msg = ste.Msg("refinar.so_v1", nome)
		}
		return processo.Processo{}, fmt.Errorf("%w: %s", ErrNaoExiste, msg)
	}
	p, err := processo.Ler(caminho)
	if err == nil && p.Nome != nome {
		err = fmt.Errorf("%w: o arquivo de %s guarda o processo %s", processo.ErrInvalido, nome, p.Nome)
	}
	return p, err
}

func indiceDe(etapas []processo.Etapa, id string) int {
	return slices.IndexFunc(etapas, func(e processo.Etapa) bool { return e.ID == id })
}

// aplicar devolve o processo com a etapa mudada (sem tocar no original) e a etapa como era.
func aplicar(p processo.Processo, m Mudanca) (processo.Processo, processo.Etapa, error) {
	i := indiceDe(p.Etapas, m.Etapa)
	if i < 0 {
		ids := make([]string, len(p.Etapas))
		for j, e := range p.Etapas {
			ids[j] = e.ID
		}
		return p, processo.Etapa{}, erroMudanca("refinar.sem_etapa", p.Nome, m.Etapa, strings.Join(ids, ", "))
	}
	antes := p.Etapas[i]
	depois, err := mudar(antes, m)
	if err != nil {
		return p, antes, err
	}
	p.Etapas = slices.Clone(p.Etapas)
	p.Etapas[i] = depois
	return p, antes, nil
}

func erroMudanca(id string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrMudanca, ste.Msg(id, args...))
}

// mudar aplica Peca, Titulo e Casos na etapa. Fica com os casos pedidos; sem eles, com os que a etapa já tinha.
func mudar(e processo.Etapa, m Mudanca) (processo.Etapa, error) {
	peca, titulo, casos := strings.TrimSpace(m.Peca), strings.TrimSpace(m.Titulo), limpar(m.Casos)
	switch {
	case e.Compoe != "":
		return e, erroMudanca("refinar.etapa_composta", e.ID)
	case peca == "" && titulo == "" && len(casos) == 0:
		return e, erroMudanca("refinar.sem_mudanca", e.ID)
	}
	if peca != "" {
		e.Peca = peca
	}
	if titulo != "" {
		e.Titulo = titulo
	}
	if len(casos) > 0 {
		e.Casos = casos
	}
	if len(e.Casos) == 0 {
		return e, erroMudanca("refinar.sem_caso", e.ID)
	}
	return e, nil
}

// limpar tira os espaços, os vazios e os repetidos, mantendo a ordem.
func limpar(termos []string) []string {
	var out []string
	for _, t := range termos {
		if t = strings.TrimSpace(t); t != "" && !slices.Contains(out, t) {
			out = append(out, t)
		}
	}
	return out
}

// subirMinor sobe a versão do meio: 1.0.0 vira 1.1.0. Uma mudança de etapa não quebra quem já usa o processo.
func subirMinor(v string) (string, error) {
	var maior, menor, correcao int
	if _, err := fmt.Sscanf(v, "%d.%d.%d", &maior, &menor, &correcao); err != nil {
		return "", fmt.Errorf("%w: versão %q: %v", processo.ErrInvalido, v, err)
	}
	return fmt.Sprintf("%d.%d.0", maior, menor+1), nil
}

// avisosDaPeca avisa, sem barrar, quando a peça nova não consta na matriz. Sem matriz legível, não avisa.
func avisosDaPeca(vault, nova, antiga string) []string {
	if nova == antiga {
		return nil
	}
	reg, err := matriz.LerRegistro(filepath.Join(vault, "registry_matrix.json"))
	if err != nil || len(reg.Entradas) == 0 {
		return nil
	}
	if slices.ContainsFunc(reg.Entradas, func(e matriz.Entrada) bool { return e.Nome == nova }) {
		return nil
	}
	return []string{ste.Msg("refinar.peca_desconhecida", nova)}
}

// montar regera as projeções do processo novo e compara cada uma com o que está no disco. Devolve o diff de
// todos os arquivos e a lista do que será gravado (o que não muda fica de fora). O processo.json vai por
// último: ele é quem diz a versão.
func montar(p processo.Processo, m Mudanca, now time.Time, roots processo.Roots) ([]Arquivo, []processo.File, error) {
	slug := strings.ToUpper(strings.ReplaceAll(strings.TrimPrefix(p.Nome, "legoloop-"), "-", "_"))
	docs := filepath.Join(roots.Docs, "docs")
	dataDocs := dataDaSpec(filepath.Join(docs, "specs"), slug, now)
	pr, err := projetar(p, now, dataDocs)
	if err != nil {
		return nil, nil, err
	}
	pj, err := p.Marshal()
	if err != nil {
		return nil, nil, err
	}
	intentPath, intent, err := intentDe(filepath.Join(docs, "intents"), slug, p, m, now)
	if err != nil {
		return nil, nil, err
	}
	vault := func(parts ...string) string { return filepath.Join(append([]string{roots.Vault}, parts...)...) }
	todos := []processo.File{
		{Kind: "INTENT", Path: intentPath, Data: intent},
		{Kind: "SPEC", Path: filepath.Join(docs, "specs", pr.SpecNome), Data: []byte(pr.Spec)},
		{Kind: "PLAN", Path: filepath.Join(docs, "plans", pr.PlanNome), Data: []byte(pr.Plan)},
		{Kind: "ROUTE", Path: filepath.Join(docs, "routes", pr.RouteNome), Data: pr.Route},
		{Kind: "DAG", Path: vault(caminhos.DirLoops, p.Nome+".dag.json"), Data: pr.DAG},
		{Kind: "loop", Path: vault(caminhos.DirLoops, p.Nome+".loop.json"), Data: pr.Loop},
		{Kind: "skill", Path: vault("skills", p.Nome, "SKILL.md"), Data: []byte(pr.Skill)},
		{Kind: "processo", Path: processo.CaminhoNoCofre(roots.Vault, p.Nome), Data: pj},
	}
	return comparar(todos)
}

// projetar usa a data do refinar para o loop e a data da SPEC que já existe para SPEC, PLAN e ROUTE: assim os
// três documentos mantêm o nome, e o ROUTE continua citando a SPEC que está no projeto.
func projetar(p processo.Processo, now, dataDocs time.Time) (*processo.Projecoes, error) {
	pr, err := processo.Projetar(p, now)
	if err != nil {
		return nil, err
	}
	if pr.Loop, err = manterCriacao(pr.Loop, p.CriadoEm); err != nil {
		return nil, err
	}
	if dataDocs.Format(time.DateOnly) == now.Format(time.DateOnly) {
		return pr, nil
	}
	pd, err := processo.Projetar(p, dataDocs)
	if err != nil {
		return nil, err
	}
	pr.Spec, pr.Plan, pr.Route = pd.Spec, pd.Plan, pd.Route
	pr.SpecNome, pr.PlanNome, pr.RouteNome = pd.SpecNome, pd.PlanNome, pd.RouteNome
	return pr, nil
}

// manterCriacao devolve o loop com created_at igual à criação do processo: só updated_at leva a data do refinar.
func manterCriacao(loopJSON []byte, criado time.Time) ([]byte, error) {
	var lp studiodoc.Processo
	if err := json.Unmarshal(loopJSON, &lp); err != nil {
		return nil, err
	}
	lp.CreatedAt = criado.UTC()
	b, err := json.MarshalIndent(lp, "", "  ")
	return append(b, '\n'), err
}

// maisRecente devolve o arquivo de dir com a data mais nova no nome e o sufixo dado, e essa data.
func maisRecente(dir, sufixo string) (string, time.Time) {
	var nome string
	var data time.Time
	entradas, _ := os.ReadDir(dir)
	for _, e := range entradas {
		n := e.Name()
		if e.IsDir() || !strings.HasSuffix(n, sufixo) || len(n) < len(time.DateOnly) {
			continue
		}
		if d, err := time.Parse(time.DateOnly, n[:len(time.DateOnly)]); err == nil && n > nome {
			nome, data = n, d
		}
	}
	return nome, data
}

// dataDaSpec acha a SPEC do processo em docs/specs. A data dela vale para os documentos; sem SPEC, vale agora.
func dataDaSpec(dir, slug string, now time.Time) time.Time {
	nome, d := maisRecente(dir, "_SPEC_"+slug+".md")
	if nome == "" {
		return now
	}
	return time.Date(d.Year(), d.Month(), d.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.UTC)
}

// intentDe devolve o INTENT com a seção do refinamento no fim. O texto da entrevista fica como estava. Sem
// INTENT no projeto, nasce um novo, com o objetivo do processo.
func intentDe(dir, slug string, p processo.Processo, m Mudanca, now time.Time) (string, []byte, error) {
	nome, _ := maisRecente(dir, "_INTENT_"+slug+".md")
	base := fmt.Sprintf("# INTENT: %s\n\n## Objetivo\n\n%s\n", p.Nome, strings.TrimSpace(p.Objetivo))
	caminho := filepath.Join(dir, now.Format(time.DateOnly)+"_INTENT_"+slug+".md")
	if nome != "" {
		caminho = filepath.Join(dir, nome)
		b, err := os.ReadFile(caminho)
		if err != nil {
			return "", nil, err
		}
		base = strings.TrimRight(string(b), "\n") + "\n"
	}
	return caminho, []byte(base + secaoIntent(p, m, now)), nil
}

func secaoIntent(p processo.Processo, m Mudanca, now time.Time) string {
	e := p.Etapas[indiceDe(p.Etapas, m.Etapa)]
	var b strings.Builder
	fmt.Fprintf(&b, "\n## Refinamento %s\n\n* Data: %s\n* Versão: %s para %s\n", p.Versao, now.Format(time.DateOnly),
		p.Historico[len(p.Historico)-1], p.Versao)
	fmt.Fprintf(&b, "* Etapa: %s (%s)\n* Peça: %s\n* Título: %s\n* Casos: %s\n", e.ID, e.Tipo, e.Peca, e.Titulo, strings.Join(e.Casos, ", "))
	return b.String()
}

// comparar lê o que está no disco e monta o diff de cada arquivo. Só entram na gravação os arquivos que mudam.
func comparar(todos []processo.File) ([]Arquivo, []processo.File, error) {
	var arqs []Arquivo
	var gravar []processo.File
	for _, f := range todos {
		a := Arquivo{Tipo: f.Kind, Caminho: f.Path, Estado: "novo", Linhas: []string{}}
		velho, err := os.ReadFile(f.Path)
		switch {
		case errors.Is(err, os.ErrNotExist):
		case err != nil:
			return nil, nil, err
		case string(velho) == string(f.Data):
			a.Estado = "igual"
		default:
			a.Estado, a.Backup = "troca", f.Path+sufixoBak
		}
		if a.Estado != "igual" {
			a.Mais, a.Menos, a.Linhas = difLinhas(string(velho), string(f.Data))
			gravar = append(gravar, f)
		}
		arqs = append(arqs, a)
	}
	return arqs, gravar, nil
}

// linhasUteis separa o texto em linhas e descarta as vazias.
func linhasUteis(s string) []string {
	var out []string
	for _, l := range strings.Split(s, "\n") {
		if l = strings.TrimSpace(l); l != "" {
			out = append(out, l)
		}
	}
	return out
}

// difLinhas compara os dois textos como conjuntos de linhas (com repetição). Devolve quantas linhas entram,
// quantas saem e uma amostra: as linhas que saem, com "- ", e as que entram, com "+ ".
func difLinhas(velho, novo string) (int, int, []string) {
	v, n := linhasUteis(velho), linhasUteis(novo)
	sobraV, sobraN := contar(v), contar(n)
	var entram, saem []string
	for _, l := range n {
		if sobraV[l] > 0 {
			sobraV[l]--
		} else {
			entram = append(entram, l)
		}
	}
	for _, l := range v {
		if sobraN[l] > 0 {
			sobraN[l]--
		} else {
			saem = append(saem, l)
		}
	}
	return len(entram), len(saem), amostra(saem, entram)
}

func contar(linhas []string) map[string]int {
	out := map[string]int{}
	for _, l := range linhas {
		out[l]++
	}
	return out
}

func amostra(saem, entram []string) []string {
	out := []string{}
	for _, g := range []struct {
		marca  string
		linhas []string
	}{{"- ", saem}, {"+ ", entram}} {
		for _, l := range g.linhas[:min(len(g.linhas), maxLinhas)] {
			out = append(out, g.marca+cortar(l))
		}
		if mais := len(g.linhas) - maxLinhas; mais > 0 {
			out = append(out, ste.Msg("refinar.diff_mais", mais))
		}
	}
	return out
}

func cortar(s string) string {
	if r := []rune(s); len(r) > maxLargura {
		return string(r[:maxLargura-3]) + "..."
	}
	return s
}

func (a Arquivo) resumo() string {
	if a.Estado == "igual" {
		return "-"
	}
	return fmt.Sprintf("+%d -%d", a.Mais, a.Menos)
}

func traco(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

// Diff é o texto que mostra a mudança: a versão, a etapa e, para cada arquivo, o estado e uma amostra das linhas.
func (p Plano) Diff() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n", ste.Msg("refinar.diff_titulo", p.Nome, p.De, p.Para))
	fmt.Fprintf(&b, "  %s\n", ste.Msg("refinar.diff_etapa", p.Depois.ID, p.Depois.Tipo))
	if p.Antes.Peca != p.Depois.Peca {
		fmt.Fprintf(&b, "  %s\n", ste.Msg("refinar.diff_peca", traco(p.Antes.Peca), p.Depois.Peca))
	}
	if p.Antes.Titulo != p.Depois.Titulo {
		fmt.Fprintf(&b, "  %s\n", ste.Msg("refinar.diff_titulo_etapa", p.Antes.Titulo, p.Depois.Titulo))
	}
	fmt.Fprintf(&b, "  %s\n", ste.Msg("refinar.diff_casos", traco(strings.Join(p.Antes.Casos, ", ")), strings.Join(p.Depois.Casos, ", ")))
	fmt.Fprintf(&b, "\n%s\n", ste.Msg("refinar.diff_colunas"))
	for _, a := range p.Arquivos {
		fmt.Fprintf(&b, formato, a.Tipo, a.Estado, a.resumo(), a.Caminho)
		for _, l := range a.Linhas {
			fmt.Fprintf(&b, "  %s\n", l)
		}
	}
	for _, av := range p.Avisos {
		fmt.Fprintf(&b, "\n%s\n", ste.Msg("refinar.diff_aviso", av))
	}
	return b.String()
}

// Cartao traduz o plano para o cartão STE-PT: o aviso de risco vem antes dos passos.
func (p Plano) Cartao() ste.Cartao {
	trocas := 0
	c := ste.Cartao{
		Acao: ste.Msg("refinar.cartao_acao", p.Nome),
		Passos: []string{
			ste.Msg("refinar.cartao_etapa", p.Depois.ID, p.Depois.Titulo),
			ste.Msg("refinar.cartao_versao", p.De, p.Para, p.De),
			ste.Msg("refinar.cartao_gera"),
		},
		SeDerErrado: ste.Msg("refinar.cartao_se_errado", p.De),
		Parametros: []ste.Parametro{
			{Nome: "peca", Valor: p.Depois.Peca, Ajuda: ste.Msg("refinar.param_peca")},
			{Nome: "casos", Valor: strings.Join(p.Depois.Casos, ", "), Ajuda: ste.Msg("refinar.param_casos")},
		},
	}
	for _, a := range p.Arquivos {
		switch a.Estado {
		case "troca":
			trocas++
			c.Efeitos = append(c.Efeitos, ste.Msg("refinar.efeito_troca", a.Caminho))
		case "novo":
			c.Efeitos = append(c.Efeitos, ste.Msg("refinar.efeito_novo", a.Caminho))
		}
	}
	c.Aviso = strings.Join(append([]string{ste.Msg("refinar.cartao_aviso", trocas)}, p.Avisos...), " ")
	return c
}

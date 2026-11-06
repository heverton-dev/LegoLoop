// Package compor liga um processo do cofre como etapa de outro. A etapa nova executa o processo composto
// (nó subloop no DAG). Antes de gravar, recusa o ciclo entre processos (ErrCiclo), o processo inexistente e a
// versão que não confere. Mostra o plano, pede a confirmação, guarda cada original como .legoloop-bak, sobe a
// versão e grava tudo de uma vez. É uma fatia independente: depende só do núcleo (processo, ste) e nunca de
// outra função.
package compor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"legoloop/internal/loop"
	"legoloop/internal/processo"
	"legoloop/internal/ste"
)

var (
	// ErrNome: o nome do processo que recebe a etapa não é de processo.
	ErrNome = errors.New("compor: nome inválido")
	// ErrNaoExiste: o cofre não tem o processo da v2 com esse nome.
	ErrNaoExiste = errors.New("compor: processo não encontrado")
	// ErrVersao: o processo composto existe, mas em outra versão.
	ErrVersao = errors.New("compor: versão não confere")
	// ErrPedido: a ligação pedida não vale (formato, etapa inexistente ou processo já ligado).
	ErrPedido = errors.New("compor: pedido inválido")
	// ErrCancelado: o usuário não confirmou. Nenhum arquivo mudou.
	ErrCancelado = errors.New("compor: cancelado")
)

var (
	reNome = regexp.MustCompile(`^legoloop-[a-z0-9]+(-[a-z0-9]+)*$`)
	reID   = regexp.MustCompile(`^e(\d+)$`)
)

const (
	sufixoBak = ".legoloop-bak" // sufixo do original que processo.Store guarda
	formato   = "%-8s %-6s %-9s %s\n"
)

func init() {
	for id, msg := range map[string]string{
		"compor.nome_invalido":      "O nome %s não é de processo. Use legoloop-<nome>, em letras minúsculas.",
		"compor.composto_invalido":  "Informe o processo que entra como legoloop-<nome>@<versão>, por exemplo legoloop-b@1.0.0. Recebido: %q.",
		"compor.nao_existe":         "O processo %s não existe no cofre %s. Veja os processos com legoloop listar.",
		"compor.so_v1":              "O processo %s é da v1 e não tem o arquivo .processo.json. Crie o processo de novo com legoloop fabricar.",
		"compor.versao":             "O processo %s está na versão %s, não na %s. Use %s@%s.",
		"compor.ciclo":              "Ligar %s em %s cria um ciclo: %s.",
		"compor.ja_compoe":          "O processo %s já usa %s na etapa %s.",
		"compor.sem_etapa":          "O processo %s não tem a etapa %s. As etapas são: %s.",
		"compor.artefato_repetido":  "A etapa %s de %s já grava %s. A etapa nova grava o artefato de %s no mesmo arquivo. O DAG não aceita dois nós no mesmo arquivo.",
		"compor.titulo_padrao":      "Executar o processo %s",
		"compor.cancelado":          "Operação cancelada. Nenhum arquivo mudou.",
		"compor.concluido":          "Processo %s ligado a %s na etapa %s: versão %s para %s.",
		"compor.entrada_invalida":   "Entrada inválida: responda 1, 2 ou 0.",
		"compor.diff_titulo":        "COMPOR %s: versão %s para %s.",
		"compor.diff_etapa":         "Etapa nova %s: %s.",
		"compor.diff_composto":      "Processo que entra: %s.",
		"compor.diff_depende":       "Depende de: %s.",
		"compor.diff_produz":        "Artefato: %s. Origem: %s, no processo %s.",
		"compor.diff_colunas":       "ARQUIVO  ESTADO MUDANÇA   CAMINHO",
		"compor.cartao_acao":        "Ligar o processo %s em %s",
		"compor.cartao_aviso":       "O programa troca %d arquivos. Cada original fica guardado como .legoloop-bak.",
		"compor.cartao_etapa":       "O programa cria a etapa %s, que executa o processo %s.",
		"compor.cartao_versao":      "O programa sobe a versão de %s para %s e guarda %s no histórico.",
		"compor.cartao_gera":        "O programa gera de novo INTENT, SPEC, PLAN, ROUTE, DAG, loop e skill.",
		"compor.cartao_se_errado":   "Se um passo falhar, o programa desfaz tudo. O processo volta à versão %s. Para voltar depois, restaure cada arquivo .legoloop-bak.",
		"compor.efeito_troca":       "troca %s",
		"compor.efeito_novo":        "grava %s",
		"compor.param_titulo":       "Título da etapa nova. Enter mantém o valor.",
		"compor.param_depende":      "Etapas que vêm antes da etapa nova, separadas por vírgula. Enter mantém o valor.",
		"compor.intent_secao":       "Composição",
		"compor.intent_sem_entrada": "Sem entrevista: a etapa nova veio do legoloop compor.",
	} {
		ste.Default.Register(id, msg)
	}
}

// Ligacao é o pedido de composição. Composto é "legoloop-<nome>@<versão>". Titulo vazio usa "Executar o processo
// <nome>". Depende vazio liga a etapa nova depois das etapas finais do processo (as que nenhuma outra usa).
type Ligacao struct {
	Composto string   `json:"composto"`
	Titulo   string   `json:"titulo,omitempty"`
	Depende  []string `json:"depende,omitempty"`
}

// Request é o pedido: o processo que recebe a etapa e a ligação.
type Request struct {
	Nome    string
	Ligacao Ligacao
}

// ConfirmFunc mostra o plano e devolve true para gravar. replan recalcula o plano com outra ligação (o
// "personalizar"); o último plano pedido com replan é o que vale na gravação.
type ConfirmFunc func(p Plano, replan func(Ligacao) (Plano, error)) (bool, error)

// Deps são as bordas do compor. Store.Roots traz o cofre e a raiz do projeto; Store.FS nil usa o disco real.
// Confirm nil segue sem perguntar. Clock nil usa time.Now.
type Deps struct {
	Store   processo.Store
	Confirm ConfirmFunc
	Clock   func() time.Time
}

// Arquivo é um dos arquivos que o compor regera. Estado é troca, novo ou igual; só troca e novo são gravados.
type Arquivo struct {
	Tipo    string `json:"tipo"`
	Caminho string `json:"caminho"`
	Estado  string `json:"estado"`
	Mais    int    `json:"mais"`  // linhas que entram
	Menos   int    `json:"menos"` // linhas que saem
	Backup  string `json:"backup,omitempty"`
}

// Plano é tudo que o compor vai fazer, mostrado ao usuário antes de qualquer efeito.
type Plano struct {
	Nome     string         `json:"nome"`
	De       string         `json:"de"`
	Para     string         `json:"para"`
	Ligacao  Ligacao        `json:"ligacao"` // a ligação com os valores finais: serve de base para o replan
	Etapa    processo.Etapa `json:"etapa"`
	Arquivos []Arquivo      `json:"arquivos"`
	files    []processo.File
}

// Resultado é o que o compor gravou.
type Resultado struct {
	Nome     string    `json:"nome"`
	De       string    `json:"de"`
	Para     string    `json:"para"`
	Composto string    `json:"composto"`
	Etapa    string    `json:"etapa"`
	Arquivos []Arquivo `json:"arquivos"`
}

// Run planeja a ligação, pede a confirmação e grava. Cancelar ou falhar deixa o cofre e o projeto como estavam.
func Run(ctx context.Context, req Request, d Deps) (*Resultado, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	atual, err := planejar(req, d)
	if err != nil {
		return nil, err
	}
	if d.Confirm != nil {
		replan := func(l Ligacao) (Plano, error) {
			p, err := planejar(Request{Nome: req.Nome, Ligacao: l}, d)
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
	return &Resultado{Nome: pl.Nome, De: pl.De, Para: pl.Para, Composto: pl.Etapa.Compoe, Etapa: pl.Etapa.ID,
		Arquivos: pl.Arquivos}, nil
}

func (d Deps) agora() time.Time {
	if d.Clock == nil {
		return time.Now().UTC()
	}
	return d.Clock().UTC()
}

// Planejar valida o composto, anti-ciclo e gera a prévia dos arquivos sem tocar no disco.
func Planejar(req Request, d Deps) (Plano, error) {
	return planejar(req, d)
}

// planejar lê os dois processos, recusa o ciclo e a versão errada, cria a etapa e regera os arquivos. Não grava.
func planejar(req Request, d Deps) (Plano, error) {
	roots := d.Store.Roots
	if roots.Vault == "" || roots.Docs == "" {
		return Plano{}, fmt.Errorf("%w: informe o cofre e a raiz do projeto", processo.ErrPlano)
	}
	if !reNome.MatchString(req.Nome) {
		return Plano{}, fmt.Errorf("%w: %s", ErrNome, ste.Msg("compor.nome_invalido", req.Nome))
	}
	alvo, versao, err := lerRef(req.Ligacao.Composto)
	if err != nil {
		return Plano{}, err
	}
	a, b, err := lerPar(roots.Vault, req.Nome, alvo, versao)
	if err != nil {
		return Plano{}, err
	}
	e, err := novaEtapa(a, b, req.Ligacao)
	if err != nil {
		return Plano{}, err
	}
	novo, err := comEtapa(a, e)
	if err != nil {
		return Plano{}, err
	}
	arqs, files, err := montar(novo, e, d.agora(), roots)
	if err != nil {
		return Plano{}, err
	}
	l := Ligacao{Composto: e.Compoe, Titulo: e.Titulo, Depende: e.Depende}
	return Plano{Nome: a.Nome, De: a.Versao, Para: novo.Versao, Ligacao: l, Etapa: e, Arquivos: arqs, files: files}, nil
}

// lerRef separa "legoloop-b@1.0.0" em nome e versão. A versão só é conferida contra o processo do cofre.
func lerRef(ref string) (string, string, error) {
	nome, versao, ok := strings.Cut(strings.TrimSpace(ref), "@")
	if !ok || !reNome.MatchString(nome) || versao == "" {
		return "", "", fmt.Errorf("%w: %s", ErrPedido, ste.Msg("compor.composto_invalido", ref))
	}
	return nome, versao, nil
}

// lerPar lê o processo que recebe a etapa e o composto. O ciclo é conferido antes da versão: ligar a em b,
// quando b já usa a, é ciclo em qualquer versão.
func lerPar(vault, nome, alvo, versao string) (processo.Processo, processo.Processo, error) {
	a, err := lerProcesso(vault, nome)
	if err != nil {
		return a, processo.Processo{}, err
	}
	b, err := lerProcesso(vault, alvo)
	if err != nil {
		return a, b, err
	}
	if err := semCiclo(vault, nome, alvo); err != nil {
		return a, b, err
	}
	if b.Versao != versao {
		return a, b, fmt.Errorf("%w: %s", ErrVersao, ste.Msg("compor.versao", alvo, b.Versao, versao, alvo, b.Versao))
	}
	return a, b, nil
}

// lerProcesso lê loops/<nome>.processo.json. Processo da v1 (só o .loop.json) não entra na composição.
func lerProcesso(vault, nome string) (processo.Processo, error) {
	caminho := processo.CaminhoNoCofre(vault, nome)
	if _, err := os.Stat(caminho); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return processo.Processo{}, err
		}
		msg := ste.Msg("compor.nao_existe", nome, vault)
		if _, err := os.Stat(filepath.Join(vault, "loops", nome+".loop.json")); err == nil {
			msg = ste.Msg("compor.so_v1", nome)
		}
		return processo.Processo{}, fmt.Errorf("%w: %s", ErrNaoExiste, msg)
	}
	p, err := processo.Ler(caminho)
	if err == nil && p.Nome != nome {
		err = fmt.Errorf("%w: o arquivo de %s guarda o processo %s", processo.ErrInvalido, nome, p.Nome)
	}
	return p, err
}

// semCiclo recusa a ligação quando o composto já chega, direta ou transitivamente, ao processo que o recebe.
func semCiclo(vault, nome, alvo string) error {
	cadeia, err := caminhoAte(vault, alvo, nome, map[string]bool{})
	if err != nil || cadeia == nil {
		return err
	}
	texto := strings.Join(append([]string{nome}, cadeia...), " -> ")
	return fmt.Errorf("%w: %s", processo.ErrCiclo, ste.Msg("compor.ciclo", alvo, nome, texto))
}

// caminhoAte procura, a partir de de, uma cadeia de composições que chega em fim. Devolve a cadeia (de ... fim)
// ou nil. Um processo composto que não está no cofre não continua a cadeia.
func caminhoAte(vault, de, fim string, visto map[string]bool) ([]string, error) {
	if de == fim {
		return []string{de}, nil
	}
	if visto[de] {
		return nil, nil
	}
	visto[de] = true
	p, err := processo.Ler(processo.CaminhoNoCofre(vault, de))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	for _, e := range p.Etapas {
		filho, _, _ := strings.Cut(e.Compoe, "@")
		if filho == "" {
			continue
		}
		if cadeia, err := caminhoAte(vault, filho, fim, visto); err != nil || cadeia != nil {
			return append([]string{de}, cadeia...), err
		}
	}
	return nil, nil
}

func indiceDe(etapas []processo.Etapa, id string) int {
	return slices.IndexFunc(etapas, func(e processo.Etapa) bool { return e.ID == id })
}

func erroPedido(id string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrPedido, ste.Msg(id, args...))
}

// novaEtapa cria a etapa que executa b. Ela entra no fim de a, com o próximo id eNN. Produz o artefato e usa o
// gate da etapa final de b: o subloop roda b num diretório isolado e só devolve esse artefato.
func novaEtapa(a, b processo.Processo, l Ligacao) (processo.Etapa, error) {
	for _, e := range a.Etapas {
		if usado, _, _ := strings.Cut(e.Compoe, "@"); usado == b.Nome {
			return processo.Etapa{}, erroPedido("compor.ja_compoe", a.Nome, b.Nome, e.ID)
		}
	}
	deps, err := dependencias(a, l.Depende)
	if err != nil {
		return processo.Etapa{}, err
	}
	ord, err := b.Ordem()
	if err != nil {
		return processo.Etapa{}, err
	}
	fim := ord[len(ord)-1]
	produz := fim.Produz
	produz.Caminho = noPai(b.Nome, fim.Produz.Caminho)
	if err := artefatoLivre(a, b, produz.Caminho); err != nil {
		return processo.Etapa{}, err
	}
	titulo := strings.TrimSpace(l.Titulo)
	if titulo == "" {
		titulo = ste.Msg("compor.titulo_padrao", b.Nome)
	}
	e := processo.Etapa{ID: proximoID(a.Etapas), Titulo: titulo, Tipo: fim.Tipo, Compoe: b.Nome + "@" + b.Versao,
		Origem: fim.Produz.Caminho, Produz: produz, Gate: gateNoPai(fim.Gate, fim.Produz.Caminho, produz.Caminho), Depende: deps}
	for _, d := range deps {
		e.Consome = append(e.Consome, a.Etapas[indiceDe(a.Etapas, d)].Produz.Caminho)
	}
	return e, nil
}

// dependencias devolve as etapas pedidas (sem espaço e sem repetição) ou, sem pedido, as etapas finais de a.
func dependencias(a processo.Processo, pedidas []string) ([]string, error) {
	var deps []string
	for _, d := range pedidas {
		if d = strings.TrimSpace(d); d != "" && !slices.Contains(deps, d) {
			deps = append(deps, d)
		}
	}
	if len(deps) == 0 {
		return finais(a.Etapas), nil
	}
	for _, d := range deps {
		if indiceDe(a.Etapas, d) < 0 {
			ids := make([]string, len(a.Etapas))
			for i, e := range a.Etapas {
				ids[i] = e.ID
			}
			return nil, erroPedido("compor.sem_etapa", a.Nome, d, strings.Join(ids, ", "))
		}
	}
	return deps, nil
}

// finais são as etapas de que nenhuma outra depende, na ordem em que foram declaradas.
func finais(es []processo.Etapa) []string {
	usada := map[string]bool{}
	for _, e := range es {
		for _, d := range e.Depende {
			usada[d] = true
		}
	}
	var out []string
	for _, e := range es {
		if !usada[e.ID] {
			out = append(out, e.ID)
		}
	}
	return out
}

// proximoID devolve eNN, com NN um a mais que o maior número entre os ids no formato eNN.
func proximoID(es []processo.Etapa) string {
	maior := 0
	for _, e := range es {
		if m := reID.FindStringSubmatch(e.ID); m != nil {
			n, _ := strconv.Atoi(m[1])
			maior = max(maior, n)
		}
	}
	return fmt.Sprintf("e%02d", maior+1)
}

// comEtapa devolve a cópia de a com a etapa nova no fim, a versão do meio somada (1.0.0 vira 1.1.0) e a
// versão anterior no histórico.
func comEtapa(a processo.Processo, e processo.Etapa) (processo.Processo, error) {
	var maior, menor, correcao int
	if _, err := fmt.Sscanf(a.Versao, "%d.%d.%d", &maior, &menor, &correcao); err != nil {
		return a, fmt.Errorf("%w: versão %q: %v", processo.ErrInvalido, a.Versao, err)
	}
	novo := a
	novo.Etapas = append(slices.Clone(a.Etapas), e)
	novo.Versao = fmt.Sprintf("%d.%d.0", maior, menor+1)
	novo.Historico = append(slices.Clone(a.Historico), a.Versao)
	return novo, nil
}

// noPai devolve o caminho do artefato de b no processo que recebe a etapa: out/<nome-de-b>/<arquivo>. O
// fabricar grava sempre out/NN_<tipo>.md, então dois processos reais repetem os nomes; no pai, o artefato de b
// fica numa pasta só dele. O subloop lê o artefato no caminho de dentro de b (a origem) e grava neste.
func noPai(nomeB, origem string) string {
	return path.Join("out", nomeB, path.Base(origem))
}

// gateNoPai troca, no gate da etapa final de b, o literal com o caminho de dentro de b pelo caminho no pai: o
// gate do nó subloop é conferido no diretório do pai.
func gateNoPai(gate, origem, destino string) string {
	return strings.ReplaceAll(gate, "'"+origem+"'", "'"+destino+"'")
}

// artefatoLivre recusa a ligação quando uma etapa de a já grava o caminho que a etapa nova vai gravar no pai:
// o DAG não aceita dois nós que gravam o mesmo arquivo.
func artefatoLivre(a, b processo.Processo, caminho string) error {
	for _, x := range a.Etapas {
		if x.Produz.Caminho == caminho {
			return erroPedido("compor.artefato_repetido", x.ID, a.Nome, x.Produz.Caminho, b.Nome)
		}
	}
	return nil
}

// montar regera as projeções do processo novo e compara cada uma com o que está no disco. O processo.json vai
// por último: ele é quem diz a versão.
func montar(p processo.Processo, e processo.Etapa, now time.Time, roots processo.Roots) ([]Arquivo, []processo.File, error) {
	slug := strings.ToUpper(strings.ReplaceAll(strings.TrimPrefix(p.Nome, "legoloop-"), "-", "_"))
	docs := filepath.Join(roots.Docs, "docs")
	pr, err := projetar(p, now, dataDaSpec(filepath.Join(docs, "specs"), slug, now))
	if err != nil {
		return nil, nil, err
	}
	pj, err := p.Marshal()
	if err != nil {
		return nil, nil, err
	}
	intentPath, intent, err := intentDe(filepath.Join(docs, "intents"), slug, p, e, now)
	if err != nil {
		return nil, nil, err
	}
	vault := func(parts ...string) string { return filepath.Join(append([]string{roots.Vault}, parts...)...) }
	return comparar([]processo.File{
		{Kind: "INTENT", Path: intentPath, Data: intent},
		{Kind: "SPEC", Path: filepath.Join(docs, "specs", pr.SpecNome), Data: []byte(pr.Spec)},
		{Kind: "PLAN", Path: filepath.Join(docs, "plans", pr.PlanNome), Data: []byte(pr.Plan)},
		{Kind: "ROUTE", Path: filepath.Join(docs, "routes", pr.RouteNome), Data: pr.Route},
		{Kind: "DAG", Path: vault("loops", p.Nome+".dag.json"), Data: pr.DAG},
		{Kind: "loop", Path: vault("loops", p.Nome+".loop.json"), Data: pr.Loop},
		{Kind: "skill", Path: vault("skills", p.Nome, "SKILL.md"), Data: []byte(pr.Skill)},
		{Kind: "processo", Path: processo.CaminhoNoCofre(roots.Vault, p.Nome), Data: pj},
	})
}

// projetar usa a data do compor para o loop e a data da SPEC que já existe para SPEC, PLAN e ROUTE: assim os
// três documentos mantêm o nome no projeto.
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

// manterCriacao devolve o loop com created_at igual à criação do processo: só updated_at leva a data do compor.
func manterCriacao(loopJSON []byte, criado time.Time) ([]byte, error) {
	var lp loop.Loop
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

// intentDe devolve o INTENT com a seção da composição no fim. O texto anterior fica como estava. Sem INTENT no
// projeto, nasce um novo, com o objetivo do processo.
func intentDe(dir, slug string, p processo.Processo, e processo.Etapa, now time.Time) (string, []byte, error) {
	nome, _ := maisRecente(dir, "_INTENT_"+slug+".md")
	base := fmt.Sprintf("# INTENT: %s\n\n## Objetivo\n\n%s\n\n%s\n", p.Nome, strings.TrimSpace(p.Objetivo),
		ste.Msg("compor.intent_sem_entrada"))
	caminho := filepath.Join(dir, now.Format(time.DateOnly)+"_INTENT_"+slug+".md")
	if nome != "" {
		caminho = filepath.Join(dir, nome)
		b, err := os.ReadFile(caminho)
		if err != nil {
			return "", nil, err
		}
		base = strings.TrimRight(string(b), "\n") + "\n"
	}
	var s strings.Builder
	fmt.Fprintf(&s, "\n## %s %s\n\n* Data: %s\n* Versão: %s para %s\n", ste.Msg("compor.intent_secao"), p.Versao,
		now.Format(time.DateOnly), p.Historico[len(p.Historico)-1], p.Versao)
	fmt.Fprintf(&s, "* Etapa: %s (%s)\n* Processo: %s\n* Depende de: %s\n", e.ID, e.Titulo, e.Compoe, strings.Join(e.Depende, ", "))
	return caminho, []byte(base + s.String()), nil
}

// comparar lê o que está no disco e conta as linhas que entram e saem. Só entram na gravação os que mudam.
func comparar(todos []processo.File) ([]Arquivo, []processo.File, error) {
	var arqs []Arquivo
	var gravar []processo.File
	for _, f := range todos {
		a := Arquivo{Tipo: f.Kind, Caminho: f.Path, Estado: "novo"}
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
			a.Mais, a.Menos = contarDif(string(velho), string(f.Data))
			gravar = append(gravar, f)
		}
		arqs = append(arqs, a)
	}
	return arqs, gravar, nil
}

// contarDif compara os dois textos como conjuntos de linhas (com repetição): quantas entram e quantas saem.
func contarDif(velho, novo string) (int, int) {
	sobra := map[string]int{}
	for _, l := range strings.Split(velho, "\n") {
		sobra[l]++
	}
	entram := 0
	for _, l := range strings.Split(novo, "\n") {
		if sobra[l] > 0 {
			sobra[l]--
		} else {
			entram++
		}
	}
	saem := 0
	for _, n := range sobra {
		saem += n
	}
	return entram, saem
}

func (a Arquivo) resumo() string {
	if a.Estado == "igual" {
		return "-"
	}
	return fmt.Sprintf("+%d -%d", a.Mais, a.Menos)
}

// Diff é o texto que mostra a ligação: a versão, a etapa nova e, para cada arquivo, o estado e a contagem.
func (p Plano) Diff() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n", ste.Msg("compor.diff_titulo", p.Nome, p.De, p.Para))
	fmt.Fprintf(&b, "  %s\n", ste.Msg("compor.diff_etapa", p.Etapa.ID, p.Etapa.Titulo))
	fmt.Fprintf(&b, "  %s\n", ste.Msg("compor.diff_composto", p.Etapa.Compoe))
	fmt.Fprintf(&b, "  %s\n", ste.Msg("compor.diff_depende", strings.Join(p.Etapa.Depende, ", ")))
	fmt.Fprintf(&b, "  %s\n", ste.Msg("compor.diff_produz", p.Etapa.Produz.Caminho, p.Etapa.Origem, p.Etapa.Compoe))
	fmt.Fprintf(&b, "\n%s\n", ste.Msg("compor.diff_colunas"))
	for _, a := range p.Arquivos {
		fmt.Fprintf(&b, formato, a.Tipo, a.Estado, a.resumo(), a.Caminho)
	}
	return b.String()
}

// Cartao traduz o plano para o cartão STE-PT: o aviso de risco vem antes dos passos.
func (p Plano) Cartao() ste.Cartao {
	trocas := 0
	c := ste.Cartao{
		Acao: ste.Msg("compor.cartao_acao", p.Etapa.Compoe, p.Nome),
		Passos: []string{
			ste.Msg("compor.cartao_etapa", p.Etapa.ID, p.Etapa.Compoe),
			ste.Msg("compor.cartao_versao", p.De, p.Para, p.De),
			ste.Msg("compor.cartao_gera"),
		},
		SeDerErrado: ste.Msg("compor.cartao_se_errado", p.De),
		Parametros: []ste.Parametro{
			{Nome: "titulo", Valor: p.Etapa.Titulo, Ajuda: ste.Msg("compor.param_titulo")},
			{Nome: "depende", Valor: strings.Join(p.Etapa.Depende, ", "), Ajuda: ste.Msg("compor.param_depende")},
		},
	}
	for _, a := range p.Arquivos {
		switch a.Estado {
		case "troca":
			trocas++
			c.Efeitos = append(c.Efeitos, ste.Msg("compor.efeito_troca", a.Caminho))
		case "novo":
			c.Efeitos = append(c.Efeitos, ste.Msg("compor.efeito_novo", a.Caminho))
		}
	}
	c.Aviso = ste.Msg("compor.cartao_aviso", trocas)
	return c
}

// MensagemConcluido é o texto de quando a ligação termina.
func MensagemConcluido(r *Resultado) string {
	return ste.Msg("compor.concluido", r.Composto, r.Nome, r.Etapa, r.De, r.Para)
}

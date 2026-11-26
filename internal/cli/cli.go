// Package cli implementa os comandos do binário legoloop. Toda E/S passa por Deps,
// permitindo testar a CLI real sem estado global.
package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/heverton-dev/LegoLoop/internal/caminhos"
	"io"
	"os"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/heverton-dev/LegoLoop/internal/dag"
	"github.com/heverton-dev/LegoLoop/internal/provider"
	"github.com/heverton-dev/LegoLoop/internal/videorea"
)

// Version é a versão do binário.
const Version = "0.1.0"

// Deps injeta saídas, descoberta de provedores e sinal de prontidão do Studio.
type Deps struct {
	Stdin    io.Reader
	Stdout   io.Writer
	Stderr   io.Writer
	Discover func(context.Context) provider.Inventory
	Ready    chan<- string
	// Home substitui o diretório do usuário no inventário de ativos (vazio = HOME real; usado em testes).
	Home string
	// Fabricar injeta LLM, scanner e relógio do comando fabricar (nil = ambiente real; usado em testes).
	Fabricar *FabricarDeps
	// Matriz injeta varredura, cérebro e relógio do comando matriz (nil = ambiente real; usado em testes).
	Matriz *MatrizDeps
	// VideoRunner injeta o executor do yt-dlp do video-rea ingest (nil = yt-dlp real; usado em testes).
	VideoRunner videorea.Runner
	// Rota injeta inventário e perfis do comando rota (nil = ambiente real; usado em testes).
	Rota *RotaDeps
	// Homologar injeta o harness da bateria e a memória livre do comando homologar (nil = ambiente real; usado em testes).
	Homologar *HomologarDeps
	// Provider injeta cliente HTTP, arquivos e host canônico da sonda de provider (nil = ambiente real; usado em testes).
	Provider *ProviderDeps
	// Fechar injeta seams do comando fechar (nil = ambiente real; usado em testes).
	Fechar *FecharDeps
	// LegoPress injeta compilador e auditor do comando legopress (nil = ambiente real; usado em testes).
	LegoPress *LegoPressDeps
	// WebQARunner injeta o executor do script Node do comando web-qa (nil = node real; usado em testes).
	WebQARunner func(ctx context.Context, args []string, stdout, stderr io.Writer) int
}

func (d Deps) discover(ctx context.Context) provider.Inventory {
	if d.Discover != nil {
		return d.Discover(ctx)
	}
	return provider.Discover(ctx, provider.DefaultConfig())
}

const usage = `Uso: legoloop <comando> [opções]

Comandos:
  doctor     [--json]                        inventário de harnesses, modelos e chaves (Cofre de Segredos)
  canonize   --root DIR... (--dry-run|--apply) [--vault DIR] [--workers N] [--json]
  validate   <` + caminhos.ArquivoProcessoDAG + `> [--regenerar-studio]  valida contrato, gates, ciclos e a paridade do documento do Studio derivado
  run        <` + caminhos.ArquivoProcessoDAG + `> [--workdir DIR] [--parallel N] [--timeout DUR]
             [--skills DIR] [--var chave=valor]... [--json]
  exec       <slug> [--processos-dir DIR] [--workdir DIR] [--json]  executa o processo pelo slug e grava os artefatos
  route      --from Artifact:X --to Artifact:Y  rota ótima no grafo semântico
  studio     [--addr 127.0.0.1:7878] [--project DIR]  abre o Looping Studio
  forge      [--addr 127.0.0.1:7878] [--project DIR]  abre o Studio e dispara primeira varredura
  import     <github-url> [--name N] [--json]         importa e canoniza skill do GitHub
  upgrade    <arquivo> [--diagnose|--apply] [--json]  diagnóstica ou valida elevação para Padrão Ouro
  fabricar   "<nome>" [--yes] [--goal T] [--model M] [--dry-run] [--json]  cria um processo novo (skill, DAG e documentos)
  rota       --plan P --out ROUTE [--perfil F] [--json]  gera a ROUTE com harness, via, modelo e esforço resolvidos
  homologar  [--pares H|V,...] [--yes] [--non-interactive] [--teto DUR] [--provider] [--corrigir-provider] [--json]  mede a conformidade de cada harness e via
  matriz     gerar|status|vocab|perfil|ver  gera e consulta a matriz de peças (qualidade e perfil de uso)
  listar     [--vault DIR] [--json]  mostra os processos do cofre (versão, natureza, harnesses, usado por)
  dissecar   <alvo> [opções]                 dissecação técnica e extração estrutural de alvo
  refinar    <nome> --etapa ID [--peca P] [--caso C]... [--yes]  muda uma etapa, sobe a versão e guarda a anterior
  remover    <nome> [--vault DIR] [--dry-run] [--force] [--yes] [--json]  apaga um processo do cofre e os vínculos dele
  renomear   <antigo> <novo> [--dry-run] [--vault DIR] [--docs DIR] [--yes] [--json]  renomeia um processo com backup e registro
  compor     <nome> <legoloop-x@versão> [--depende ID]... [--yes]  liga um processo como etapa de outro (recusa ciclo)
  forjar          [--dry-run|--apply|--desfazer DATA|--verificar --orcamento ARQ] [--perfil P] [--sem H] [--home DIR] [--json]  executa a Forja ou verifica o orçamento de tokens por harness
  instalar        [--dry-run|--apply] [--perfil P] [--home DIR] [--json]  executa a Forja, grava os relatórios e configura o atalho do grafo
  atualizar       [--canal estavel|beta] [--base-url URL] [--binario PATH] [--home DIR] [--json]  atualiza o binário com verificação SHA-256, cópia .prev e revalidação da Forja
  instalar-skills [--vault DIR] [--dry-run] [--json]  instala e vincula as skills nativas universais
  skills-custo    [--json]  mede os tokens de frontmatter que cada harness carrega por turno
  skills-migrar (--dry-run|--apply) [--vault DIR]  tira as skills eager de todos os harnesses (backup) e deixa só o índice lazy
  bootstrap       [--vault DIR] [--dry-run] [--json]  orquestra skills, mcp-router e ambiente sem config manual
  stept      [status|set|on|off|scan] [--harness H] [--modo M] [--pausa P] [--json]  alternador de governança STE-PT
  sync-harnesses (--dry-run|--apply) [--json]         conecta harnesses ao legoloop mcp-router
  mcp-router [--gateway]                     inicia servidor MCP STDIO com Lazy discovery (--gateway: só call_mcp_tool)
  mcp-registry (list|refresh)                lista ou indexa os MCPs externos que o router sobe sob demanda
  legopress  [doctor|auditar|criar|editar|atualizar|adequar|preview]  processo industrial de livros técnicos
  video-motion <subcomando> [opções]          processo determinístico de edição e motion vídeo
  video-rea  <subcomando> [opções]             processo universal de engenharia reversa audiovisual
  canal      redc|gerar [opções]              processo REDC de canais e clonagem estrutural
  cartoon-shorts --theme "<tema>" [opções]    processo soberano de geração de Cartoon Shorts (1080x1920)
  swarm      <add|post|read|wait|status|kill|monitor|ui> [-dir DIR] [-addr HOST:PORT] [-root DIR]  painel de coordenação de agentes em enxame
  fechar     <slug> [--dry-run] [--skip-build] [--json]  cadeia determinística de 9 etapas de canonização e fechamento
  web-qa     <preflight|run|corrigir> [opções]  QA web determinístico (Playwright, axe-core, lychee, Obscura, renderlint, unlighthouse, lost-pixel)
  sessoes    <slug> [--json] [--processos-dir DIR]  lista as sessões ativas do processo (memória de execução)
  anticorpos <sync|status> [slug] [--raiz DIR]   sincroniza ou inspeciona o espelho executivo ANTICORPOS.md
  espinha    verificar [--raiz DIR] [--json]     verifica a integridade deterministica da espinha dorsal
  progresso  <gerar|iniciar|validar|status> [opções]  acompanhamento vivo de Tickets validado por Gate
  diorama-roteiro <prevoo|validar-config|coletar|validar-roteiro> [opções]  pré-voo, descoberta e validação epistêmica de roteiro diorama
  diorama-ativos  <prevoo|executar|validar-storyboard> [opções]             elenco, sprites, cenários SVG/Flow e storyboard cenas.json
  video-inventario <validar|verificar> --estilo ID [opções]                 validação humana de estilo e conferência por hash do inventário
  version
`

func mapaComandos() map[string]func(context.Context, []string, Deps) int {
	return map[string]func(context.Context, []string, Deps) int{
		"doctor": cmdDoctor, "canonize": cmdCanonizeJSONSafe, "validate": cmdValidate,
		"run": cmdRun, "exec": cmdExec, "route": cmdRoute, "studio": cmdStudio, "forge": cmdForge,
		"import": cmdImport, "upgrade": cmdUpgrade, "fabricar": cmdFabricar, "matriz": cmdMatriz, "listar": cmdListar,
		"dissecar": cmdDissecar,
		"refinar":  cmdRefinar, "remover": cmdRemover, "renomear": cmdRenomear, "compor": cmdCompor, "forjar": cmdForjar, "instalar": cmdInstalar, "atualizar": cmdAtualizar, "instalar-skills": cmdInstalarSkills, "skills-custo": cmdSkillsCusto, "skills-migrar": cmdSkillsMigrar,
		"stept": cmdStept, "sync-harnesses": cmdSyncHarnesses, "mcp-router": cmdMcpRouter, "mcp-registry": cmdMcpRegistry, "bootstrap": cmdBootstrap,
		"legopress": cmdLegoPress, "video-motion": cmdVideoMotion, "video-rea": cmdVideoRea, "fechar": cmdFechar, "commit": cmdCommit,
		"rota": cmdRota, "homologar": cmdHomologar, "canal": cmdCanal, "cartoon-shorts": cmdCartoonShorts, "legoloop-cartoon-shorts": cmdCartoonShorts,
		"progresso": cmdProgresso, "swarm": cmdSwarm, "web-qa": cmdWebQA, "sessoes": cmdSessoes, "anticorpos": cmdAnticorpos, "espinha": cmdEspinha,
		"diorama-roteiro": cmdDioramaRoteiro, "diorama-ativos": cmdDioramaAtivos, "video-inventario": cmdVideoInventario,
	}
}

// Verbos devolve, em ordem alfabetica deterministica, todos os nomes de verbos
// registrados no mapa real de comandos da CLI.
func Verbos() []string {
	cmds := mapaComandos()
	nomes := make([]string, 0, len(cmds))
	for nome := range cmds {
		nomes = append(nomes, nome)
	}
	slices.Sort(nomes)
	return nomes
}

// Run despacha o subcomando e devolve o exit code (0 ok, 1 falha, 2 uso incorreto).
func Run(ctx context.Context, args []string, d Deps) int {
	if len(args) == 0 {
		fmt.Fprint(d.Stderr, usage)
		return 2
	}
	cmds := mapaComandos()
	switch name := args[0]; name {
	case "help", "-h", "--help":
		fmt.Fprint(d.Stdout, usage)
		return 0
	case "version", "--version":
		fmt.Fprintf(d.Stdout, "legoloop %s\n", Version)
		return 0
	default:
		cmd, ok := cmds[name]
		if !ok {
			fmt.Fprintf(d.Stderr, "comando desconhecido %q\n\n%s", name, usage)
			return 2
		}
		return executarSubcomando(ctx, cmd, args, d)
	}
}

// executarSubcomando aplica o guia de migração do caminho legado, o help do
// subcomando e só então despacha o comando.
func executarSubcomando(ctx context.Context, cmd func(context.Context, []string, Deps) int, args []string, d Deps) int {
	if alvo := caminhoLegado(args[1:]); alvo != "" {
		return guiaMigracaoEscrever(d.Stderr, alvo)
	}
	if hasHelpArg(args[1:]) {
		return runSubcommandHelp(ctx, cmd, args[1:], d)
	}
	return cmd(ctx, args[1:], d)
}

func hasHelpArg(args []string) bool {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			return true
		}
	}
	return false
}

func runSubcommandHelp(ctx context.Context, cmd func(context.Context, []string, Deps) int, args []string, d Deps) int {
	dHelp := d
	dHelp.Stderr = d.Stdout
	_ = cmd(ctx, args, dHelp)
	return 0
}

func cmdUpgrade(_ context.Context, args []string, d Deps) int {
	return RunUpgrade(d.Stdout, d.Stderr, args)
}

func cmdDissecar(_ context.Context, args []string, d Deps) int {
	cmd := NewDissecarCommand()
	cmd.output = d.Stdout
	if err := cmd.Run(args); err != nil {
		fmt.Fprintln(d.Stderr, err)
		return 1
	}
	return 0
}

func newFlags(name string, d Deps) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(d.Stderr)
	return fs
}

// splitPositional aceita `cmd ARQ --flags` e `cmd --flags ARQ`.
func splitPositional(args []string) (string, []string) {
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		return args[0], args[1:]
	}
	return "", args
}

func fail(d Deps, err error) int {
	fmt.Fprintln(d.Stderr, "erro:", err)
	return 1
}

func printJSON(d Deps, v any) int {
	enc := json.NewEncoder(d.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return fail(d, err)
	}
	return 0
}

func table(w io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
}

// listFlag acumula flags repetíveis (--root A --root B).
type listFlag []string

func (l *listFlag) String() string     { return strings.Join(*l, ",") }
func (l *listFlag) Set(v string) error { *l = append(*l, v); return nil }

func humanBytes(n int64) string {
	units := []string{"KiB", "MiB", "GiB", "TiB"}
	if n < 1024 {
		return fmt.Sprintf("%d B", n)
	}
	v, i := float64(n)/1024, 0
	for v >= 1024 && i < len(units)-1 {
		v, i = v/1024, i+1
	}
	return fmt.Sprintf("%.1f %s", v, units[i])
}

// readFileArg lê no máximo dag.MaxDAGBytes (+1 para detectar excesso).
func readFileArg(path string) ([]byte, error) {
	if path == "" {
		return nil, errors.New("informe o arquivo " + caminhos.ArquivoDAG)
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, dag.MaxDAGBytes+1))
	if err == nil && len(data) > dag.MaxDAGBytes {
		err = fmt.Errorf("%s excede %d bytes", path, dag.MaxDAGBytes)
	}
	return data, err
}

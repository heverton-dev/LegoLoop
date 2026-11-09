// Package cli implementa os comandos do binário legoloop. Toda E/S passa por Deps,
// permitindo testar a CLI real sem estado global.
package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"legoloop/internal/dag"
	"legoloop/internal/provider"
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
}

func (d Deps) discover(ctx context.Context) provider.Inventory {
	if d.Discover != nil {
		return d.Discover(ctx)
	}
	return provider.Discover(ctx, provider.DefaultConfig())
}

const usage = `Uso: legoloop <comando> [opções]

Comandos:
  doctor     [--json]                        inventário de harnesses, modelos e chaves
  canonize   --root DIR... (--dry-run|--apply) [--vault DIR] [--workers N] [--json]
  validate   <pipeline.dag.json>             valida contrato, gates e ciclos
  run        <pipeline.dag.json> [--workdir DIR] [--parallel N] [--timeout DUR]
             [--skills DIR] [--var chave=valor]... [--json]
  route      --from Artifact:X --to Artifact:Y  rota ótima no grafo semântico
  studio     [--addr 127.0.0.1:7878] [--project DIR]  abre o Looping Studio
  forge      [--addr 127.0.0.1:7878] [--project DIR]  abre o Studio e dispara primeira varredura
  import     <github-url> [--name N] [--json]         importa e canoniza skill do GitHub
  upgrade    <arquivo> [--diagnose|--apply] [--json]  diagnóstica ou valida elevação para Padrão Ouro
  matrix     [--kind K] [--tier T] [--min N] [--json] DESCONTINUADO: use "matriz ver"
  fabricar   "<nome>" [--yes] [--goal T] [--model M] [--dry-run] [--json]  cria um processo novo (skill, loop, dag e docs)
  matriz     gerar|status|vocab|perfil|ver  gera e consulta a matriz de peças (qualidade e perfil de uso)
  listar     [--vault DIR] [--json]  mostra os processos do cofre (versão, natureza, harnesses, usado por)
  refinar    <nome> --etapa ID [--peca P] [--caso C]... [--yes]  muda uma etapa, sobe a versão e guarda a anterior
  remover    <nome> [--vault DIR] [--yes]  apaga um processo do cofre e os vínculos dele (falha se outro processo o usa)
  compor     <nome> <legoloop-x@versão> [--depende ID]... [--yes]  liga um processo como etapa de outro (recusa ciclo)
  instalar-skills [--vault DIR] [--dry-run] [--json]  instala e vincula as skills nativas universais
  stept      [status|set|on|off|scan] [--harness H] [--modo M] [--pausa P] [--json]  alternador de governança STE-PT
  sync-harnesses (--dry-run|--apply)         conecta harnesses ao legoloop mcp-router
  mcp-router                                 inicia servidor MCP STDIO com Lazy discovery
  legopress  [doctor|auditar|criar|editar|atualizar|adequar|preview]  esteira industrial de livros técnicos
  video-motion <subcomando> [opções]          pipeline determinístico de edição e motion vídeo
  video-capcut <subcomando> [opções]          pipeline soberano de edição talking-head com CapCut e FCPXML
  video-rea  <subcomando> [opções]             pipeline universal de engenharia reversa audiovisual
  canal      redc|gerar [opções]              pipeline REDC de canais e clonagem estrutural
  cartoon-shorts --theme "<tema>" [opções]    pipeline soberano de geração de Cartoon Shorts (1080x1920)
  fechar     <slug> [--dry-run] [--skip-build] [--json]  cadeia determinística de 8 etapas de canonização e fechamento
  version
`

// Run despacha o subcomando e devolve o exit code (0 ok, 1 falha, 2 uso incorreto).
func Run(ctx context.Context, args []string, d Deps) int {
	if len(args) == 0 {
		fmt.Fprint(d.Stderr, usage)
		return 2
	}
	cmds := map[string]func(context.Context, []string, Deps) int{
		"doctor": cmdDoctor, "canonize": cmdCanonize, "validate": cmdValidate,
		"run": cmdRun, "route": cmdRoute, "studio": cmdStudio, "forge": cmdForge,
		"import": cmdImport, "upgrade": cmdUpgrade, "matrix": cmdMatrix, "fabricar": cmdFabricar, "matriz": cmdMatriz, "listar": cmdListar,
		"refinar": cmdRefinar, "remover": cmdRemover, "compor": cmdCompor, "instalar-skills": cmdInstalarSkills,
		"stept": cmdStept, "sync-harnesses": cmdSyncHarnesses, "mcp-router": cmdMcpRouter,
		"legopress": cmdLegoPress, "video-motion": cmdVideoMotion, "video-rea": cmdVideoRea, "fechar": cmdFechar,
		"canal": cmdCanal, "cartoon-shorts": cmdCartoonShorts, "legoloop-cartoon-shorts": cmdCartoonShorts,
	}
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
		return cmd(ctx, args[1:], d)
	}
}

func cmdUpgrade(_ context.Context, args []string, d Deps) int {
	return RunUpgrade(d.Stdout, d.Stderr, args)
}

func cmdMatrix(_ context.Context, args []string, d Deps) int {
	return RunMatrix(d.Stdout, d.Stderr, args)
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

func cmdDoctor(ctx context.Context, args []string, d Deps) int {
	fs := newFlags("doctor", d)
	asJSON := fs.Bool("json", false, "saída JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	inv := d.discover(ctx)
	if *asJSON {
		return printJSON(d, map[string]any{"inventory": inv, "models": inv.Models()})
	}
	w := table(d.Stdout)
	fmt.Fprintln(w, "MODELO\tEXECUTOR")
	for _, m := range inv.Models() {
		fmt.Fprintf(w, "%s\t%s\n", m.ID, m.Source)
	}
	if err := w.Flush(); err != nil {
		return fail(d, err)
	}
	router := "offline"
	if inv.RouterURL != "" {
		router = "online em " + inv.RouterURL
	}
	fmt.Fprintf(d.Stdout, "\n9Router: %s\n", router)
	fmt.Fprintf(d.Stdout, "Harnesses no PATH: %s\n", joinOrNone(cliNames(inv.CLIs)))
	fmt.Fprintf(d.Stdout, "Chaves de API detectadas: %s\n", joinOrNone(inv.APIKeys))
	for _, e := range inv.Errors {
		fmt.Fprintf(d.Stdout, "aviso: %s\n", e)
	}
	return 0
}

func cliNames(cs []provider.CLI) []string {
	var out []string
	for _, c := range cs {
		out = append(out, c.Name)
	}
	return out
}

func joinOrNone(xs []string) string {
	if len(xs) == 0 {
		return "nenhum"
	}
	return strings.Join(xs, ", ")
}

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
		return nil, errors.New("informe o arquivo pipeline.dag.json")
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

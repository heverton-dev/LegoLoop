// Comando ticketgate é o Gate padrão de cada Ticket: uma bateria única de verificação.
//
// Uso: go run ./tools/ticketgate [-pkgs a,b] [-tax a,b] [-tax-max N] [-gold] [-build] [-livros a,b] [-trace a,b]
//
// O DAG aceita um único command por nó. Esta peça reúne a bateria em um único argv e
// falha com exit 1 se qualquer passo falhar. O exit 2 indica erro de uso.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// init faz execCmd achar o gofmt do GOROOT quando ele não está no PATH.
// O nome lógico "gofmt" continua o mesmo para quem substitui execCmd.
func init() {
	base := execCmd
	execCmd = func(dir, name string, args ...string) (string, error) {
		return base(dir, resolveTool(name), args...)
	}
}

// resolveTool devolve o caminho do gofmt em $GOROOT/bin se o PATH não o tem.
func resolveTool(name string) string {
	if name != "gofmt" {
		return name
	}
	if _, err := exec.LookPath(name); err == nil {
		return name
	}
	goroot := strings.TrimSpace(os.Getenv("GOROOT"))
	if goroot == "" {
		if out, err := exec.Command("go", "env", "GOROOT").Output(); err == nil {
			goroot = strings.TrimSpace(string(out))
		}
	}
	if goroot == "" {
		return name
	}
	exe := "gofmt"
	if runtime.GOOS == "windows" {
		exe += ".exe"
	}
	return filepath.Join(goroot, "bin", exe)
}

type config struct {
	Root      string
	Pkgs      []string
	Tax       []string
	TaxMax    int
	Gold      bool
	Build     bool
	Livros    []string
	Trace     []string
	Interview []string
	Bin       string
}

// step é um passo da bateria. Run devolve a saída e um erro quando o passo falha.
type step struct {
	Name string
	Run  func() (string, error)
}

func splitList(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func parseArgs(args []string, errOut io.Writer) (config, error) {
	fs := flag.NewFlagSet("ticketgate", flag.ContinueOnError)
	fs.SetOutput(errOut)
	pkgs := fs.String("pkgs", "", "pacotes Go separados por vírgula (ex.: ./internal/dag/...)")
	tax := fs.String("tax", "", "pastas para o taxlinter, separadas por vírgula")
	taxMax := fs.Int("tax-max", 0, "tolera até N violações no taxlinter (padrão 0)")
	gold := fs.Bool("gold", false, "roda o teste dourado (./internal/golden/...)")
	build := fs.Bool("build", false, "roda go build ./...")
	livros := fs.String("livros", "", "pastas com PDFs para legopress auditar")
	trace := fs.String("trace", "", "pastas ou arquivos da cadeia de rastreabilidade, separados por vírgula")
	interview := fs.String("interview", "", "entrevistas para a regra de prova, separadas por vírgula")
	root := fs.String("root", ".", "raiz do repositório Go")
	bin := fs.String("bin", defaultBin(), "binário legoloop")
	if err := fs.Parse(args); err != nil {
		return config{}, err
	}
	if *taxMax < 0 {
		return config{}, fmt.Errorf("-tax-max deve ser maior ou igual a 0")
	}
	return config{Root: *root, Pkgs: splitList(*pkgs), Tax: splitList(*tax), TaxMax: *taxMax,
		Gold: *gold, Build: *build, Livros: splitList(*livros), Trace: splitList(*trace), Interview: splitList(*interview), Bin: *bin}, nil
}

func buildSteps(c config) []step {
	var steps []step
	if len(c.Pkgs) > 0 {
		steps = append(steps, testStep(c))
	}
	steps = append(steps, vetStep(c), funclenStep(c), gofmtStep(c), modStep(c))
	if len(c.Tax) > 0 {
		steps = append(steps, taxStep(c))
	}
	if c.Gold {
		steps = append(steps, goldStep(c))
	}
	if c.Build {
		steps = append(steps, buildStep(c))
	}
	if len(c.Livros) > 0 {
		steps = append(steps, livrosStep(c))
	}
	if len(c.Trace) > 0 {
		steps = append(steps, traceStep(c))
	}
	if len(c.Interview) > 0 {
		steps = append(steps, interviewStep(c.Root, c.Interview))
	}
	return steps
}

// runSteps executa todos os passos e devolve o número de falhas.
func runSteps(steps []step, out io.Writer) int {
	failed := 0
	for _, s := range steps {
		output, err := s.Run()
		if err != nil {
			failed++
			fmt.Fprintf(out, "[FALHA] %s: %v\n%s\n", s.Name, err, strings.TrimSpace(output))
			continue
		}
		fmt.Fprintf(out, "[ok]    %s\n", s.Name)
	}
	return failed
}

// run executa a bateria e devolve o exit code: 0 passou, 1 falhou, 2 erro de uso.
func run(args []string, out, errOut io.Writer) int {
	cfg, err := parseArgs(args, errOut)
	if err != nil {
		return 2
	}
	if n := runSteps(buildSteps(cfg), out); n > 0 {
		fmt.Fprintf(out, "ticketgate: %d passo(s) falharam\n", n)
		return 1
	}
	fmt.Fprintln(out, "ticketgate: todos os passos passaram")
	return 0
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// Comando ticketgate é o Gate padrão de cada Ticket: uma bateria única de verificação.
//
// Uso: go run ./tools/ticketgate [-pkgs a,b] [-tax a,b] [-gold] [-build] [-livros a,b]
//
// O DAG aceita um único command por nó. Esta peça reúne a bateria em um único argv e
// falha com exit 1 se qualquer passo falhar. O exit 2 indica erro de uso.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type config struct {
	Root   string
	Pkgs   []string
	Tax    []string
	Gold   bool
	Build  bool
	Livros []string
	Bin    string
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
	gold := fs.Bool("gold", false, "roda o teste dourado (./internal/golden/...)")
	build := fs.Bool("build", false, "roda go build ./...")
	livros := fs.String("livros", "", "pastas com PDFs para legopress auditar")
	root := fs.String("root", ".", "raiz do repositório Go")
	bin := fs.String("bin", defaultBin(), "binário legoloop")
	if err := fs.Parse(args); err != nil {
		return config{}, err
	}
	return config{Root: *root, Pkgs: splitList(*pkgs), Tax: splitList(*tax),
		Gold: *gold, Build: *build, Livros: splitList(*livros), Bin: *bin}, nil
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

func main() {
	cfg, err := parseArgs(os.Args[1:], os.Stderr)
	if err != nil {
		os.Exit(2)
	}
	if n := runSteps(buildSteps(cfg), os.Stdout); n > 0 {
		fmt.Fprintf(os.Stdout, "ticketgate: %d passo(s) falharam\n", n)
		os.Exit(1)
	}
	fmt.Println("ticketgate: todos os passos passaram")
}

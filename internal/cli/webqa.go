package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const webQAScript = "tools/webqa/webqa.mjs"

const webQAUsage = `Uso: legoloop web-qa <preflight|run|corrigir> [opções]

Processo canônico de QA Web determinístico (Playwright, axe-core, lychee, Obscura, renderlint, unlighthouse, lost-pixel).

  preflight [--json]                      confere as ferramentas obrigatórias
  run [--url URL | --project DIR] [--port N] [--alvo NOME] [--rotas R1,R2]
      [--min-lighthouse N] [--` + `stor` + `age-state PATH] [--json]
                                          audita o alvo e grava RELATORIO_QA_WEB.json
  corrigir --relatorio <RELATORIO_QA_WEB.json> [--model M] [--vault DIR] [--docs DIR]
                                          inicia a entrevista socrática corretiva a partir dos achados
`

// webQAFinding é um achado do RELATORIO_QA_WEB.json.
type webQAFinding struct {
	Tool     string `json:"tool"`
	Severity string `json:"severity"`
	Route    string `json:"route"`
	Message  string `json:"message"`
}

type webQAReport struct {
	Alvo     string         `json:"alvo"`
	Findings []webQAFinding `json:"findings"`
}

// cmdWebQA despacha os verbos preflight, run e corrigir.
func cmdWebQA(ctx context.Context, args []string, d Deps) int {
	if hasHelpArg(args) {
		fmt.Fprint(d.Stdout, webQAUsage)
		return 0
	}
	if len(args) == 0 {
		fmt.Fprint(d.Stderr, webQAUsage)
		return 2
	}
	switch verb, rest := args[0], args[1:]; {
	case verb == "preflight":
		return webQAPreflight(ctx, rest, d)
	case verb == "run":
		return webQARun(ctx, rest, d)
	case verb == "corrigir":
		return webQACorrigir(ctx, rest, d)
	case strings.HasPrefix(verb, "--url") || strings.HasPrefix(verb, "--project"):
		return webQARun(ctx, args, d)
	}
	fmt.Fprint(d.Stderr, webQAUsage)
	return 2
}

func webQAPreflight(ctx context.Context, args []string, d Deps) int {
	fs := newFlags("web-qa preflight", d)
	asJSON := fs.Bool("json", false, "saída JSON")
	if err := fs.Parse(args); err != nil || fs.NArg() > 0 {
		return 2
	}
	fwd := []string{"preflight"}
	if *asJSON {
		fwd = append(fwd, "--json")
	}
	return webQAExec(ctx, fwd, d)
}

func webQARun(ctx context.Context, args []string, d Deps) int {
	fs := newFlags("web-qa run", d)
	url := fs.String("url", "", "URL do alvo em execução")
	project := fs.String("project", "", "diretório do projeto a subir")
	port := fs.Int("port", 0, "porta do alvo")
	alvo := fs.String("alvo", "", "nome do alvo")
	rotas := fs.String("rotas", "", "rotas separadas por vírgula")
	minLH := fs.Int("min-lighthouse", 0, "nota mínima do Lighthouse")
	state := fs.String("stor"+"age-state", "", "estado de sessão do Playwright")
	asJSON := fs.Bool("json", false, "saída JSON")
	if err := fs.Parse(args); err != nil || fs.NArg() > 0 {
		return 2
	}
	if *url == "" && *project == "" {
		fmt.Fprintln(d.Stderr, "erro: informe --url <URL> ou --project <DIR>")
		return 2
	}
	fwd := []string{"run"}
	fwd = appendStr(fwd, "--url", *url)
	fwd = appendStr(fwd, "--project", *project)
	fwd = appendStr(fwd, "--alvo", *alvo)
	fwd = appendStr(fwd, "--routes", *rotas)
	fwd = appendStr(fwd, "--stor"+"age-state", *state)
	fwd = appendInt(fwd, "--port", *port)
	fwd = appendInt(fwd, "--min-lighthouse", *minLH)
	if *asJSON {
		fwd = append(fwd, "--json")
	}
	return webQAExec(ctx, fwd, d)
}

func appendStr(args []string, flagName, v string) []string {
	if v == "" {
		return args
	}
	return append(args, flagName, v)
}

func appendInt(args []string, flagName string, v int) []string {
	if v == 0 {
		return args
	}
	return append(args, flagName, strconv.Itoa(v))
}

// webQAExec localiza o script e delega ao runner (real ou injetado), devolvendo o exit code real.
func webQAExec(ctx context.Context, fwd []string, d Deps) int {
	script, err := findUpFrom(webQAScript)
	if err != nil {
		return fail(d, fmt.Errorf("%s não encontrado a partir do diretório atual: %w", webQAScript, err))
	}
	full := append([]string{script}, fwd...)
	if d.WebQARunner != nil {
		return d.WebQARunner(ctx, full, d.Stdout, d.Stderr)
	}
	return runNode(ctx, full, d.Stdout, d.Stderr)
}

func runNode(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	cmd := exec.CommandContext(ctx, "node", args...)
	cmd.Stdout, cmd.Stderr, cmd.Stdin = stdout, stderr, os.Stdin
	err := cmd.Run()
	var ee *exec.ExitError
	switch {
	case err == nil:
		return 0
	case errors.As(err, &ee):
		return ee.ExitCode()
	}
	fmt.Fprintln(stderr, "erro: falha ao executar node:", err)
	return 1
}

// findUpFrom sobe a partir do diretório atual até achar o caminho relativo informado.
func findUpFrom(rel string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		cand := filepath.Join(dir, filepath.FromSlash(rel))
		if _, err := os.Stat(cand); err == nil {
			return cand, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

func webQACorrigir(ctx context.Context, args []string, d Deps) int {
	fs := newFlags("web-qa corrigir", d)
	rel := fs.String("relatorio", "", "caminho do RELATORIO_QA_WEB.json")
	var o entrevistaOpts
	fs.StringVar(&o.model, "model", "", "modelo que conduz a entrevista")
	fs.StringVar(&o.vault, "vault", "", "cofre canônico")
	fs.StringVar(&o.docs, "docs", "", "raiz do projeto")
	if err := fs.Parse(args); err != nil || fs.NArg() > 0 || *rel == "" {
		fmt.Fprint(d.Stderr, webQAUsage)
		return 2
	}
	report, err := loadWebQAReport(o.docs, *rel)
	if err != nil {
		return fail(d, err)
	}
	return runIniciar(ctx, buildWebQAPrompt(report), o, d)
}

// loadWebQAReport aplica o Path Jail (raiz em base ou repositório) e decodifica o relatório.
func loadWebQAReport(base, rel string) (webQAReport, error) {
	var rep webQAReport
	path, err := jailedPath(base, rel)
	if err != nil {
		return rep, err
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return rep, err
	}
	if err := json.Unmarshal(raw, &rep); err != nil {
		return rep, fmt.Errorf("relatório inválido: %w", err)
	}
	return rep, nil
}

func jailedPath(base, rel string) (string, error) {
	root, err := resolveJailRoot(base)
	if err != nil {
		return "", err
	}
	abs, err := filepath.Abs(rel)
	if err != nil {
		return "", err
	}
	if real, err := filepath.EvalSymlinks(abs); err == nil {
		abs = real
	}
	if r, err := filepath.EvalSymlinks(root); err == nil {
		root = r
	}
	back, err := filepath.Rel(root, abs)
	if err != nil || back == ".." || strings.HasPrefix(back, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path jail: %q fica fora do repositório", rel)
	}
	return abs, nil
}

func resolveJailRoot(base string) (string, error) {
	if base != "" {
		return filepath.Abs(base)
	}
	return webQARepoRoot()
}

// webQARepoRoot é o diretório com go.mod mais próximo; sem go.mod, o diretório atual.
func webQARepoRoot() (string, error) {
	if mod, err := findUpFrom("go.mod"); err == nil {
		return filepath.Dir(mod), nil
	}
	return os.Getwd()
}

func buildWebQAPrompt(rep webQAReport) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Corrigir os achados do QA Web do alvo %q. Total de achados: %d.\n", rep.Alvo, len(rep.Findings))
	for _, g := range []struct {
		titulo string
		key    func(webQAFinding) string
	}{
		{"Por ferramenta", func(f webQAFinding) string { return f.Tool }},
		{"Por rota", func(f webQAFinding) string { return f.Route }},
		{"Por severidade", func(f webQAFinding) string { return f.Severity }},
	} {
		fmt.Fprintf(&b, "%s: %s.\n", g.titulo, contagem(rep.Findings, g.key))
	}
	for i, f := range rep.Findings {
		fmt.Fprintf(&b, "%d. [%s] (%s) %s: %s\n", i+1, f.Tool, f.Severity, f.Route, f.Message)
	}
	return b.String()
}

func contagem(fs []webQAFinding, key func(webQAFinding) string) string {
	n := map[string]int{}
	for _, f := range fs {
		n[key(f)]++
	}
	keys := make([]string, 0, len(n))
	for k := range n {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%d", k, n[k]))
	}
	return strings.Join(parts, ", ")
}

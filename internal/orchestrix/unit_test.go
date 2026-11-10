package orchestrix

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/dag"
)

func nodeUnit() dag.Node {
	return dag.Node{
		ID: "TICK-X.1", Type: dag.TypeUnitTask,
		ContextIn:  []string{"a/in.go", "a/ainda_nao.go"},
		ContextOut: []string{"a/out.go", "scratch/x.txt"},
		TestSeam:   map[string]string{"file": "a/out_test.go", "signature": "TestOut"},
		SpecRef:    "AC-X.1",
		Command:    []string{"go", "test", "./a/..."},
		Output:     dag.HandoffContract{Artifact: "scratch/gates/TICK-X.1.json", Gate: "exit_code == 0"},
	}
}

func TestCheckBoundaryAceitaPermitidos(t *testing.T) {
	err := checkBoundary([]string{"a/out.go", `a\out_test.go`, "scratch/x.txt"}, allowedWrites(nodeUnit()))
	if err != nil {
		t.Fatalf("deveria aceitar: %v", err)
	}
}

// Mutação: cada desvio precisa reprovar o nó.
func TestCheckBoundaryRecusaDesvios(t *testing.T) {
	allowed := allowedWrites(nodeUnit())
	casos := map[string]string{
		"fora da lista":  "a/outro.go",
		"sobe diretório": "../fora.go",
		"sobe escondido": "a/../../fora.go",
		"absoluto":       "/etc/passwd",
		"unidade":        `C:\Windows\x`,
		"git":            ".git/hooks/pre-commit",
	}
	for nome, c := range casos {
		if err := checkBoundary([]string{c}, allowed); err == nil {
			t.Errorf("%s: %q deveria reprovar", nome, c)
		}
	}
}

func TestCheckBoundaryJuntaTodosOsDesvios(t *testing.T) {
	err := checkBoundary([]string{"x.go", "y.go"}, []string{"a.go"})
	if err == nil || !strings.Contains(err.Error(), "x.go") || !strings.Contains(err.Error(), "y.go") {
		t.Fatalf("erro deveria citar os dois desvios: %v", err)
	}
}

func TestBuildUnitPrompt(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "a"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "a", "in.go"), []byte("package a // CONTEUDO"), 0o644); err != nil {
		t.Fatal(err)
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	p, err := buildUnitPrompt(nodeUnit(), root)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"TICK-X.1", "AC-X.1", "TestOut", "exit_code == 0", "go test ./a/...",
		"- a/out.go", "- a/out_test.go", "package a // CONTEUDO", "(arquivo ainda não existe)", "Não use sub-agentes"} {
		if !strings.Contains(p, want) {
			t.Errorf("prompt sem %q", want)
		}
	}
}

func TestBuildUnitPromptRecusaContextoForaDoProjeto(t *testing.T) {
	root, err := os.OpenRoot(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	n := nodeUnit()
	n.ContextIn = []string{"../segredo.txt"}
	if _, err := buildUnitPrompt(n, root); err == nil {
		t.Fatal("context_in com .. deveria falhar")
	}
}

func TestUnitModelAliasesRouterCombos(t *testing.T) {
	o := unitOrchestrix("n1.txt=x")
	cases := []struct{ model, fallback, want string }{
		{"code-fast", "flash", "sonnet"},
		{"code-pro", "flash", "opus"},
		{"sonnet", "sonnet", "sonnet"},
	}
	for _, c := range cases {
		n := unitNode("n1", "out/n1.json", "exit_code == 0", nil)
		n.Routing = dag.RoutingConfig{Mode: dag.ModeSuggested, Model: c.model, Fallback: c.fallback}
		got, err := o.unitModel(n)
		if err != nil || got != c.want {
			t.Errorf("%s/%s: modelo = %q, err = %v; esperado %q", c.model, c.fallback, got, err, c.want)
		}
	}
	n := unitNode("n1", "out/n1.json", "exit_code == 0", nil)
	n.Routing = dag.RoutingConfig{Mode: dag.ModeManual, Model: "inexistente"}
	if _, err := o.unitModel(n); err == nil {
		t.Fatal("modelo sem alias e fora do inventário deveria falhar")
	}
}

package orchestrix

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/provider"
)

func TestNoDeCodigoEmHarnessHomologado(t *testing.T) {
	dir := t.TempDir()
	o := &Orchestrix{
		Inventory: provider.Inventory{CLIs: []provider.CLI{{Name: "omp", Path: os.Args[0]}}},
		ExtraEnv:  []string{"LEGOLOOP_HELPER=omp"},
	}
	text, usage, err := o.callHarnessInDir(context.Background(), provider.SourceOmp, "homologado", "omp:x", "PROMPT-DO-NO", dir, "low")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(text, "PROMPT-DO-NO") || usage.Input != 40 || usage.Output != 7 {
		t.Fatalf("text = %q, usage = %+v", text, usage)
	}
	if _, err := os.Lstat(filepath.Join(dir, promptFile)); err == nil {
		t.Fatal("prompt.md deveria sair da pasta de trabalho")
	}
}

func TestNoDeCodigoEmHarnessReprovadoRecusado(t *testing.T) {
	o := &Orchestrix{
		Inventory: provider.Inventory{CLIs: []provider.CLI{{Name: "omp", Path: filepath.Join(t.TempDir(), "inexistente")}}},
	}
	for _, estado := range []string{"reprovado", "pendente", ""} {
		_, _, err := o.callHarnessInDir(context.Background(), provider.SourceOmp, estado, "omp:x", "p", t.TempDir(), "")
		if err == nil || !strings.Contains(err.Error(), "recusado") {
			t.Fatalf("estado %q: err = %v", estado, err)
		}
	}
}

func nodeOutroHarness(estado string) (*Orchestrix, func(string) dag.Node) {
	o := unitOrchestrix("")
	o.Inventory.CLIs = append(o.Inventory.CLIs, provider.CLI{Name: "omp", Path: os.Args[0]})
	o.ExtraEnv = []string{"LEGOLOOP_HELPER=omp"}
	o.EstadoHarness = func(harness, via string) string {
		if harness == "omp" && via == "oficial" {
			return estado
		}
		return ""
	}
	return o, func(id string) dag.Node {
		n := unitNode(id, "out/"+id+".json", "exit_code == 0", nil)
		n.ContextOut = nil
		n.Routing = dag.RoutingConfig{Mode: dag.ModeManual, Harness: "omp", Via: "oficial", Model: "omp:x"}
		return n
	}
}

func TestNoOutroHarnessReprovadoRecusado(t *testing.T) {
	repo := newRepoT(t)
	o, no := nodeOutroHarness("reprovado")
	_, err := o.Execute(context.Background(), no("n1"), repo)
	if err == nil || !strings.Contains(err.Error(), "recusado") {
		t.Fatalf("err = %v", err)
	}
	if c := worktreeCount(t, repo); c != 1 {
		t.Fatalf("worktree vazou: %d", c)
	}
}

func TestNoOutroHarnessHomologadoExecuta(t *testing.T) {
	repo := newRepoT(t)
	o, no := nodeOutroHarness("homologado")
	if _, err := o.Execute(context.Background(), no("n1"), repo); err != nil {
		t.Fatal(err)
	}
	if c := worktreeCount(t, repo); c != 1 {
		t.Fatalf("worktree vazou: %d", c)
	}
}

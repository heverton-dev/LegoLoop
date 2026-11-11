package orchestrix_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/orchestrix"
	"legoloop/internal/runner"
	"legoloop/internal/studiodoc"
)

// subloopFinal grava um sub-loop de um nó que escreve FINAL.txt no diretório dele.
func subloopFinal(t *testing.T) *orchestrix.Orchestrix {
	t.Helper()
	store, err := studiodoc.NewStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	sub := &studiodoc.Processo{
		ID: "sub-final", Name: "Sub Final", Version: "1.0.0", Trigger: &dag.Trigger{Type: dag.TriggerManual},
		Nodes: []dag.Node{{
			ID: "n1", Type: dag.TypeScript, Command: []string{os.Args[0], "FINAL.txt"},
			Routing: dag.RoutingConfig{Mode: dag.ModeManual, Model: "local"},
			Output:  dag.HandoffContract{Artifact: "FINAL.txt", Gate: "file_exists('FINAL.txt')"},
		}},
	}
	if err := store.Save(sub); err != nil {
		t.Fatal(err)
	}
	return &orchestrix.Orchestrix{LoopStore: store, AllowScripts: true, ExtraEnv: []string{"LEGOLOOP_HELPER=substep"}}
}

func noSub(artefato, sub string) dag.Node {
	return dag.Node{
		ID: "sub_node", Type: dag.TypeSubloop, Ref: "sub-final",
		Routing: dag.RoutingConfig{Mode: dag.ModeManual, Model: "subloop"},
		Output:  dag.HandoffContract{Artifact: artefato, SubArtifact: sub, Gate: "file_exists('" + artefato + "')"},
	}
}

// TestSubloopGravaNoCaminhoDoPai: com sub_artifact, o artefato é lido no sub-loop por sub_artifact e gravado no
// pai por artifact. O arquivo do pai com o mesmo nome do sub-loop fica intacto.
func TestSubloopGravaNoCaminhoDoPai(t *testing.T) {
	orch := subloopFinal(t)
	work := t.TempDir()
	if err := os.WriteFile(filepath.Join(work, "FINAL.txt"), []byte("do pai"), 0o600); err != nil {
		t.Fatal(err)
	}
	n := noSub("out/sub-final/FINAL.txt", "FINAL.txt")
	res, err := orch.Execute(context.Background(), n, work)
	if err != nil || res.ExitCode != 0 || res.Model != "subloop/sub-final" {
		t.Fatalf("Execute: %v %+v", err, res)
	}
	if b, err := os.ReadFile(filepath.Join(work, "out", "sub-final", "FINAL.txt")); err != nil || string(b) != "ok FINAL.txt" {
		t.Errorf("artefato no caminho do pai: %q (%v)", b, err)
	}
	if b, _ := os.ReadFile(filepath.Join(work, "FINAL.txt")); string(b) != "do pai" {
		t.Errorf("o arquivo do pai com o mesmo nome mudou: %q", b)
	}
	d, err := dag.Parse([]byte(`{"name":"pai","nodes":[{"id":"sub_node","type":"subloop","ref":"sub-final",
		"routing":{"mode":"manual","model":"subloop"},
		"output":{"artifact":"out/sub-final/FINAL.txt","sub_artifact":"FINAL.txt","gate":"file_exists('out/sub-final/FINAL.txt')"}}]}`))
	if err != nil {
		t.Fatal(err)
	}
	rep, err := runner.Run(context.Background(), d, runner.Options{Workdir: t.TempDir(), Executor: orch})
	if err != nil || !rep.OK {
		t.Errorf("runner.Run com sub_artifact: %v %+v", err, rep)
	}
}

// TestSubloopSemSubArtifactUsaOMesmoCaminho: o DAG antigo (sem sub_artifact) segue gravando no mesmo caminho.
func TestSubloopSemSubArtifactUsaOMesmoCaminho(t *testing.T) {
	orch := subloopFinal(t)
	work := t.TempDir()
	if _, err := orch.Execute(context.Background(), noSub("FINAL.txt", ""), work); err != nil {
		t.Fatal(err)
	}
	if b, err := os.ReadFile(filepath.Join(work, "FINAL.txt")); err != nil || string(b) != "ok FINAL.txt" {
		t.Errorf("sem sub_artifact o artefato vai para o mesmo caminho: %q (%v)", b, err)
	}
}

// TestSubloopSubArtifactAusenteOuForaDoSubloop: o caminho que o sub-loop não grava falha citando esse caminho; o
// caminho que sai do diretório do sub-loop não é lido. Nada é gravado no pai.
func TestSubloopSubArtifactAusenteOuForaDoSubloop(t *testing.T) {
	orch := subloopFinal(t)
	fora := filepath.Join(os.TempDir(), "legoloop-fora-subloop.txt")
	if err := os.WriteFile(fora, []byte("segredo"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(fora) })
	for _, sub := range []string{"NAO_EXISTE.txt", "../legoloop-fora-subloop.txt"} {
		work := t.TempDir()
		_, err := orch.Execute(context.Background(), noSub("out/x.txt", sub), work)
		if err == nil || !strings.Contains(err.Error(), sub) {
			t.Errorf("%s: esperava erro citando o caminho, veio %v", sub, err)
		}
		if _, err := os.Stat(filepath.Join(work, "out", "x.txt")); err == nil {
			t.Errorf("%s: o pai não pode receber artefato", sub)
		}
	}
}

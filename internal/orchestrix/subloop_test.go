package orchestrix_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/orchestrix"
	"legoloop/internal/provider"
	"legoloop/internal/runner"
	"legoloop/internal/studiodoc"
)

func TestRunSubloopSuccess(t *testing.T) {
	loopsDir := t.TempDir()
	store, err := studiodoc.NewStore(loopsDir)
	if err != nil {
		t.Fatal(err)
	}

	// Subloop com 2 nós: n1 (intermediário) e n2 (final)
	sub := &studiodoc.Processo{
		ID:      "pesquisa-mercado",
		Name:    "Pesquisa de Mercado",
		Version: "2.0.0",
		Trigger: &dag.Trigger{Type: dag.TriggerManual},
		Nodes: []dag.Node{
			{
				ID:      "sub_n1",
				Type:    dag.TypeScript,
				Command: []string{os.Args[0], "INTERMEDIARIO.txt"},
				Routing: dag.RoutingConfig{Mode: dag.ModeManual, Model: "local"},
				Output: dag.HandoffContract{
					Artifact: "INTERMEDIARIO.txt",
					Gate:     "file_exists('INTERMEDIARIO.txt')",
				},
			},
			{
				ID:        "sub_n2",
				Type:      dag.TypeScript,
				DependsOn: []string{"sub_n1"},
				Command:   []string{os.Args[0], "FINAL.txt"},
				Routing:   dag.RoutingConfig{Mode: dag.ModeManual, Model: "local"},
				Output: dag.HandoffContract{
					Artifact: "FINAL.txt",
					Gate:     "file_exists('FINAL.txt')",
				},
			},
		},
	}
	if err := store.Save(sub); err != nil {
		t.Fatal(err)
	}

	workdir := t.TempDir()

	orch := &orchestrix.Orchestrix{
		Inventory:    provider.Inventory{},
		LoopStore:    store,
		AllowScripts: true,
		ExtraEnv:     []string{"LEGOLOOP_HELPER=substep"},
	}

	subNode := dag.Node{
		ID:      "sub_node",
		Type:    dag.TypeSubloop,
		Ref:     "pesquisa-mercado",
		Routing: dag.RoutingConfig{Mode: dag.ModeManual, Model: "subloop"},
		Output: dag.HandoffContract{
			Artifact: "FINAL.txt",
			Gate:     "file_exists('FINAL.txt')",
		},
	}

	res, err := orch.Execute(context.Background(), subNode, workdir)
	if err != nil {
		t.Fatalf("orch.Execute falhou no subloop: %v", err)
	}

	if res.ExitCode != 0 {
		t.Errorf("ExitCode esperado 0, obteve %d", res.ExitCode)
	}

	// 1. O artefato final DEVE estar no workdir pai
	finalPath := filepath.Join(workdir, "FINAL.txt")
	if _, err := os.Stat(finalPath); err != nil {
		t.Errorf("FINAL.txt deveria existir no workdir pai: %v", err)
	}

	// 2. O artefato intermediário NÃO deve ter vazado para o workdir pai
	interPath := filepath.Join(workdir, "INTERMEDIARIO.txt")
	if _, err := os.Stat(interPath); err == nil {
		t.Errorf("INTERMEDIARIO.txt NÃO deveria existir no workdir pai (contexto estéril)")
	}

	// 3. Execução completa via runner.Run no DAG pai
	parentDAGBytes := []byte(`{
		"name": "parent",
		"version": "2.0.0",
		"trigger": {"type": "manual"},
		"nodes": [{
			"id": "sub_node",
			"type": "subloop",
			"ref": "pesquisa-mercado",
			"routing": {"mode": "manual", "model": "subloop"},
			"output": {"artifact": "FINAL.txt", "gate": "file_exists('FINAL.txt')"}
		}]
	}`)
	parentDAG, err := dag.Parse(parentDAGBytes)
	if err != nil {
		t.Fatal(err)
	}

	rep, err := runner.Run(context.Background(), parentDAG, runner.Options{
		Workdir:  workdir,
		Executor: orch,
	})
	if err != nil {
		t.Fatalf("runner.Run falhou: %v", err)
	}
	if !rep.OK {
		t.Errorf("report deveria ser OK: %+v", rep)
	}
}

func TestRunSubloopNotFound(t *testing.T) {
	loopsDir := t.TempDir()
	store, err := studiodoc.NewStore(loopsDir)
	if err != nil {
		t.Fatal(err)
	}

	orch := &orchestrix.Orchestrix{
		LoopStore:    store,
		AllowScripts: true,
	}

	subNode := dag.Node{
		ID:      "sub_node",
		Type:    dag.TypeSubloop,
		Ref:     "loop-fantasma",
		Routing: dag.RoutingConfig{Mode: dag.ModeManual, Model: "subloop"},
		Output:  dag.HandoffContract{Artifact: "out.txt", Gate: "true"},
	}

	_, err = orch.Execute(context.Background(), subNode, t.TempDir())
	if err == nil {
		t.Fatal("esperava erro para subloop inexistente")
	}
	if !strings.Contains(err.Error(), "não encontrado") {
		t.Fatalf("mensagem esperada com 'não encontrado', obteve: %v", err)
	}
}

func TestRunSubloopInputTransfer(t *testing.T) {
	loopsDir := t.TempDir()
	store, err := studiodoc.NewStore(loopsDir)
	if err != nil {
		t.Fatal(err)
	}

	sub := &studiodoc.Processo{
		ID:      "echo-input",
		Name:    "Echo Input",
		Version: "2.0.0",
		Trigger: &dag.Trigger{Type: dag.TriggerManual},
		Nodes: []dag.Node{
			{
				ID:      "n1",
				Type:    dag.TypeScript,
				Command: []string{os.Args[0], "OUT.txt"},
				Routing: dag.RoutingConfig{Mode: dag.ModeManual, Model: "local"},
				Output:  dag.HandoffContract{Artifact: "OUT.txt", Gate: "file_exists('INPUT.txt')"},
			},
		},
	}
	if err := store.Save(sub); err != nil {
		t.Fatal(err)
	}

	workdir := t.TempDir()
	// Escreve INPUT.txt no pai
	if err := os.WriteFile(filepath.Join(workdir, "INPUT.txt"), []byte("dado"), 0600); err != nil {
		t.Fatal(err)
	}

	orch := &orchestrix.Orchestrix{
		LoopStore:    store,
		AllowScripts: true,
		ExtraEnv:     []string{"LEGOLOOP_HELPER=substep"},
	}

	subNode := dag.Node{
		ID:      "sub_node",
		Type:    dag.TypeSubloop,
		Ref:     "echo-input",
		Inputs:  map[string]string{"in": "INPUT.txt"},
		Routing: dag.RoutingConfig{Mode: dag.ModeManual, Model: "subloop"},
		Output:  dag.HandoffContract{Artifact: "OUT.txt", Gate: "file_exists('OUT.txt')"},
	}

	res, err := orch.Execute(context.Background(), subNode, workdir)
	if err != nil {
		t.Fatalf("Execute falhou: %v", err)
	}
	if res.ExitCode != 0 {
		t.Errorf("ExitCode != 0: %d", res.ExitCode)
	}
}

package runner

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/triagem"
)

func parseBudgetDAG(t *testing.T, node dag.Node) *dag.DAG {
	t.Helper()
	if node.Type == "" {
		node.Type = "script"
		node.Command = []string{"x"}
	}
	if node.Routing.Mode == "" {
		node.Routing = dag.RoutingConfig{Mode: "manual", Model: "none"}
	}
	if node.Output.Artifact == "" {
		node.Output = dag.HandoffContract{Artifact: node.ID + ".md", Gate: "exit_code == 0"}
	}
	raw, err := json.Marshal(map[string]any{
		"name":  "budget-test",
		"nodes": []dag.Node{node},
	})
	if err != nil {
		t.Fatal(err)
	}
	d, err := dag.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

type spyExecutor struct {
	called     bool
	lastPrompt string
}

func (s *spyExecutor) Execute(_ context.Context, n dag.Node, _ string) (Result, error) {
	s.called = true
	s.lastPrompt = n.Inputs["prompt"]
	return Result{ExitCode: 0}, nil
}

func TestTokenBudgetGate_AbortAndonCord(t *testing.T) {
	dir := t.TempDir()
	exec := &spyExecutor{}
	d := parseBudgetDAG(t, dag.Node{
		ID:          "n1",
		Name:        "no-estourado",
		Inputs:      map[string]string{"prompt": strings.Repeat("palavra ", 200)}, // ~1600 chars => ~400 tokens
		TokenBudget: map[string]any{"max_prompt_tokens": 50},
	})

	rep, err := Run(context.Background(), d, Options{
		Workdir:  dir,
		Executor: exec,
	})
	if err != nil {
		t.Fatalf("Run retornou erro inesperado: %v", err)
	}
	if rep.OK {
		t.Fatalf("esperava rep.OK == false quando token_budget é excedido")
	}
	if exec.called {
		t.Fatalf("executor não deveria ser chamado quando Andon Cord aborta por token_budget")
	}
	if len(rep.Nodes) != 1 || rep.Nodes[0].Status != StatusError {
		t.Fatalf("esperava StatusError no nó, obtido: %+v", rep.Nodes)
	}
	if !strings.Contains(rep.Nodes[0].Error, "max_prompt_tokens") {
		t.Fatalf("erro deveria mencionar max_prompt_tokens, obtido: %q", rep.Nodes[0].Error)
	}

	mem, memErr := triagem.LerMemoria(dir)
	if memErr != nil {
		t.Fatalf("esperava memoria_execucao.json gravada, erro: %v", memErr)
	}
	if mem.TokenBudget == nil || mem.TokenBudget.Action != "abort" {
		t.Fatalf("esperava telemetria de abort em memoria_execucao.json, obtido: %+v", mem.TokenBudget)
	}
}

func TestTokenBudgetGate_TruncatePolicy(t *testing.T) {
	dir := t.TempDir()
	exec := &spyExecutor{}
	d := parseBudgetDAG(t, dag.Node{
		ID:     "n1",
		Name:   "no-compactado",
		Inputs: map[string]string{"prompt": strings.Repeat("a", 400)}, // ~100 tokens
		TokenBudget: map[string]any{
			"max_prompt_tokens": 25,
			"policy":            "truncate",
		},
	})

	rep, err := Run(context.Background(), d, Options{
		Workdir:  dir,
		Executor: exec,
	})
	if err != nil {
		t.Fatalf("Run retornou erro inesperado: %v", err)
	}
	if !rep.OK {
		t.Fatalf("esperava rep.OK == true com política truncate, erro: %+v", rep.Nodes)
	}
	if !exec.called {
		t.Fatalf("executor deveria ter sido chamado após truncamento")
	}
	if len(exec.lastPrompt) > 100 {
		t.Fatalf("prompt deveria ter sido truncado para <= 100 chars (25 tokens), obtido %d", len(exec.lastPrompt))
	}

	mem, memErr := triagem.LerMemoria(dir)
	if memErr != nil {
		t.Fatalf("esperava memoria_execucao.json gravada, erro: %v", memErr)
	}
	if mem.TokenBudget == nil || mem.TokenBudget.Action != "truncate" {
		t.Fatalf("esperava telemetria de truncate em memoria_execucao.json, obtido: %+v", mem.TokenBudget)
	}
}

func TestTokenBudgetGate_WithinLimit(t *testing.T) {
	dir := t.TempDir()
	exec := &spyExecutor{}
	d := parseBudgetDAG(t, dag.Node{
		ID:          "n1",
		Name:        "no-ok",
		Inputs:      map[string]string{"prompt": "prompt curto"},
		TokenBudget: map[string]any{"max_prompt_tokens": 1000},
	})

	rep, err := Run(context.Background(), d, Options{
		Workdir:  dir,
		Executor: exec,
	})
	if err != nil {
		t.Fatalf("Run retornou erro inesperado: %v", err)
	}
	if !rep.OK || !exec.called {
		t.Fatalf("esperava execução normal dentro do limite")
	}
}

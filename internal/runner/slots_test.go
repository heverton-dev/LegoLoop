package runner

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/dag"
)

type slotTrackerExec struct {
	steps []string
}

func (s *slotTrackerExec) Execute(ctx context.Context, n dag.Node, workdir string) (Result, error) {
	s.steps = append(s.steps, "main:"+n.ID)
	_ = os.WriteFile(filepath.Join(workdir, n.Output.Artifact), []byte("ok"), 0o644)
	return Result{ExitCode: 0, Model: "test"}, nil
}

func TestRunExecutesSlotGatesAndRequiresAllToPass(t *testing.T) {
	workdir := t.TempDir()

	src := `{"name":"test-slot-gates","version":"2.0.0","trigger":{"type":"manual"},"nodes":[{
		"id":"n-gate",
		"type":"script",
		"command":["echo","1"],
		"routing":{"mode":"manual","model":"none"},
		"slots":{
			"gates":["file_exists('extra.txt')"]
		},
		"output":{"artifact":"res.txt","gate":"file_exists('res.txt') && exit_code == 0"}
	}]}`

	d, err := dag.Parse([]byte(src))
	if err != nil {
		t.Fatalf("falha ao parsear dag: %v", err)
	}

	exec := &slotTrackerExec{}
	// Sem extra.txt, deve falhar no gate do slot
	rep, err := Run(context.Background(), d, Options{
		Workdir:      workdir,
		Executor:     exec,
		AllowScripts: true,
	})
	if err != nil {
		t.Fatalf("Run falhou inesperadamente com erro: %v", err)
	}
	if rep.OK || rep.Nodes[0].Status != StatusGateFailed {
		t.Fatalf("esperava StatusGateFailed pelo gate de slot, veio: %+v", rep.Nodes[0])
	}

	// Agora criando extra.txt, deve passar
	_ = os.WriteFile(filepath.Join(workdir, "extra.txt"), []byte("extra"), 0o644)
	rep2, err := Run(context.Background(), d, Options{
		Workdir:      workdir,
		Executor:     exec,
		AllowScripts: true,
	})
	if err != nil {
		t.Fatalf("Run falhou na 2a vez: %v", err)
	}
	if !rep2.OK || rep2.Nodes[0].Status != StatusPassed {
		t.Fatalf("esperava StatusPassed, veio: %+v", rep2.Nodes[0])
	}
}

func TestRunSlotHooksRequireConsent(t *testing.T) {
	workdir := t.TempDir()

	src := `{"name":"test-slot-consent","version":"2.0.0","trigger":{"type":"manual"},"nodes":[{
		"id":"n-hook",
		"type":"script",
		"command":["echo","1"],
		"routing":{"mode":"manual","model":"none"},
		"slots":{
			"hooks":{"pre":[["echo","pre"]]}
		},
		"output":{"artifact":"res.txt","gate":"true"}
	}]}`

	d, err := dag.Parse([]byte(src))
	if err != nil {
		t.Fatalf("falha ao parsear dag: %v", err)
	}

	exec := &slotTrackerExec{}
	// Sem AllowScripts deve recusar
	rep, err := Run(context.Background(), d, Options{
		Workdir:      workdir,
		Executor:     exec,
		AllowScripts: false,
	})
	if err == nil && rep.OK {
		t.Fatalf("esperava falha por falta de consentimento de script/hook")
	}
}

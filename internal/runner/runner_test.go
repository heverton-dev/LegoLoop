package runner

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"legoloop/internal/dag"
)

// fakeExec grava o artefato do nó e mede a concorrência real observada.
type fakeExec struct {
	delay    time.Duration
	exitCode map[string]int
	fail     map[string]bool
	skipFile map[string]bool
	active   atomic.Int32
	maxSeen  atomic.Int32
	mu       sync.Mutex
	order    []string
}

func (f *fakeExec) Execute(ctx context.Context, n dag.Node, workdir string) (Result, error) {
	cur := f.active.Add(1)
	defer f.active.Add(-1)
	for {
		prev := f.maxSeen.Load()
		if cur <= prev || f.maxSeen.CompareAndSwap(prev, cur) {
			break
		}
	}
	f.mu.Lock()
	f.order = append(f.order, n.ID)
	f.mu.Unlock()
	select {
	case <-time.After(f.delay):
	case <-ctx.Done():
		return Result{}, ctx.Err()
	}
	if f.fail[n.ID] {
		return Result{}, errors.New("boom")
	}
	if !f.skipFile[n.ID] {
		p := filepath.Join(workdir, n.Output.Artifact)
		if err := os.WriteFile(p, []byte("ok "+n.ID), 0o644); err != nil {
			return Result{}, err
		}
	}
	return Result{ExitCode: f.exitCode[n.ID], Tokens: TokenUsage{Input: 10, Output: 5},
		Model: "fake/" + n.ID, Notes: []string{"nota " + n.ID}}, nil
}

func diamond(t *testing.T, timeoutSec int) *dag.DAG {
	t.Helper()
	mk := func(id, deps string) string {
		return `{"id":"` + id + `","type":"script","command":["x"],"depends_on":[` + deps + `],
		"timeout_sec":` + strconv.Itoa(timeoutSec) + `,"routing":{"mode":"manual","model":"none"},
		"output":{"artifact":"` + id + `.md","gate":"file_exists('` + id + `.md') && exit_code == 0"}}`
	}
	src := `{"name":"d","nodes":[` + strings.Join([]string{mk("a", ""), mk("b", `"a"`), mk("c", `"a"`), mk("d", `"b","c"`)}, ",") + `]}`
	d, err := dag.Parse([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func statuses(r Report) map[string]Status {
	m := map[string]Status{}
	for _, n := range r.Nodes {
		m[n.ID] = n.Status
	}
	return m
}

func TestRunAllPassParallel(t *testing.T) {
	ex := &fakeExec{delay: 50 * time.Millisecond}
	rep, err := Run(context.Background(), diamond(t, 0), Options{Workdir: t.TempDir(), Parallelism: 4, Executor: ex})
	if err != nil {
		t.Fatal(err)
	}
	if !rep.OK || len(rep.Nodes) != 4 {
		t.Fatalf("report = %+v", rep)
	}
	for id, s := range statuses(rep) {
		if s != StatusPassed {
			t.Errorf("%s = %s", id, s)
		}
	}
	if rep.Tokens.Input != 40 || rep.Tokens.Output != 20 {
		t.Fatalf("telemetria de tokens = %+v", rep.Tokens)
	}
	if ex.maxSeen.Load() != 2 {
		t.Fatalf("b e c deveriam rodar em paralelo, max=%d", ex.maxSeen.Load())
	}
	if ex.order[0] != "a" || ex.order[3] != "d" {
		t.Fatalf("ordem topológica violada: %v", ex.order)
	}
}

func TestRunSequential(t *testing.T) {
	ex := &fakeExec{delay: 10 * time.Millisecond}
	rep, err := Run(context.Background(), diamond(t, 0), Options{Workdir: t.TempDir(), Parallelism: 1, Executor: ex})
	if err != nil || !rep.OK {
		t.Fatalf("rep=%+v err=%v", rep, err)
	}
	if ex.maxSeen.Load() != 1 {
		t.Fatalf("modo sequencial executou %d em paralelo", ex.maxSeen.Load())
	}
}

func TestGateFailureSkipsDownstream(t *testing.T) {
	ex := &fakeExec{exitCode: map[string]int{"b": 1}}
	rep, err := Run(context.Background(), diamond(t, 0), Options{Workdir: t.TempDir(), Executor: ex})
	if err != nil {
		t.Fatal(err)
	}
	got := statuses(rep)
	if rep.OK || got["a"] != StatusPassed || got["b"] != StatusGateFailed || got["c"] != StatusPassed || got["d"] != StatusSkipped {
		t.Fatalf("statuses = %v", got)
	}
}

func TestMissingArtifactFailsGate(t *testing.T) {
	ex := &fakeExec{skipFile: map[string]bool{"a": true}}
	rep, _ := Run(context.Background(), diamond(t, 0), Options{Workdir: t.TempDir(), Executor: ex})
	if statuses(rep)["a"] != StatusGateFailed {
		t.Fatalf("statuses = %v", statuses(rep))
	}
}

func TestExecutorErrorAndTimeout(t *testing.T) {
	ex := &fakeExec{fail: map[string]bool{"a": true}}
	rep, _ := Run(context.Background(), diamond(t, 0), Options{Workdir: t.TempDir(), Executor: ex})
	if s := statuses(rep); s["a"] != StatusError || s["b"] != StatusSkipped {
		t.Fatalf("statuses = %v", s)
	}
	if rep.Nodes[0].Error == "" {
		t.Fatal("erro do executor não registrado")
	}
	slow := &fakeExec{delay: 3 * time.Second}
	start := time.Now()
	rep, _ = Run(context.Background(), diamond(t, 1), Options{Workdir: t.TempDir(), Executor: slow})
	if s := statuses(rep); s["a"] != StatusTimeout {
		t.Fatalf("statuses = %v", s)
	}
	if time.Since(start) > 2*time.Second {
		t.Fatal("timeout individual não respeitado")
	}
}

func TestDefaultTimeoutApplies(t *testing.T) {
	slow := &fakeExec{delay: 2 * time.Second}
	rep, _ := Run(context.Background(), diamond(t, 0), Options{Workdir: t.TempDir(), Executor: slow, DefaultTimeout: 50 * time.Millisecond})
	if s := statuses(rep); s["a"] != StatusTimeout {
		t.Fatalf("statuses = %v", s)
	}
}

func TestInvalidOptions(t *testing.T) {
	d := diamond(t, 0)
	if _, err := Run(context.Background(), d, Options{Workdir: t.TempDir()}); err == nil {
		t.Fatal("executor nil deveria falhar")
	}
	if _, err := Run(context.Background(), d, Options{Workdir: filepath.Join(t.TempDir(), "nao-existe"), Executor: &fakeExec{}}); err == nil {
		t.Fatal("workdir inexistente deveria falhar")
	}
	if _, err := Run(context.Background(), nil, Options{Workdir: t.TempDir(), Executor: &fakeExec{}}); err == nil {
		t.Fatal("dag nil deveria falhar")
	}
}

func TestParentCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	rep, _ := Run(ctx, diamond(t, 0), Options{Workdir: t.TempDir(), Executor: &fakeExec{delay: time.Second}})
	if s := statuses(rep); s["a"] != StatusError || rep.OK {
		t.Fatalf("statuses = %v", s)
	}
}

// A rota efetiva (modelo) e as notas de failover chegam ao relatório.
func TestReportCarriesModelAndNotes(t *testing.T) {
	rep, _ := Run(context.Background(), diamond(t, 0), Options{Workdir: t.TempDir(), Executor: &fakeExec{}})
	for _, n := range rep.Nodes {
		if n.Model != "fake/"+n.ID || len(n.Notes) != 1 || n.Notes[0] != "nota "+n.ID {
			t.Fatalf("telemetria incompleta: %+v", n)
		}
	}
}

type schemaExec struct {
	payload string
}

func (s *schemaExec) Execute(ctx context.Context, n dag.Node, workdir string) (Result, error) {
	p := filepath.Join(workdir, n.Output.Artifact)
	if err := os.WriteFile(p, []byte(s.payload), 0o644); err != nil {
		return Result{}, err
	}
	return Result{ExitCode: 0, Model: "test"}, nil
}

func TestRunnerSchemaEnforcement(t *testing.T) {
	schemaJSON := `{"type":"object","required":["title","version"],"properties":{"title":{"type":"string"},"version":{"type":"number"}}}`
	dagSrc := `{
		"name": "schema-runner-test",
		"nodes": [{
			"id": "node_json",
			"type": "script",
			"command": ["go", "version"],
			"routing": {"mode": "manual", "model": "test"},
			"output": {
				"artifact": "result.json",
				"schema": "schema.json",
				"gate": "exit_code == 0"
			}
		}]
	}`
	d, err := dag.Parse([]byte(dagSrc))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	// 1. Caso Válido
	dirValid := t.TempDir()
	if err := os.WriteFile(filepath.Join(dirValid, "schema.json"), []byte(schemaJSON), 0o644); err != nil {
		t.Fatal(err)
	}
	execValid := &schemaExec{payload: `{"title": "LegoLoop", "version": 2.0}`}
	repValid, err := Run(context.Background(), d, Options{Workdir: dirValid, Executor: execValid})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !repValid.OK || repValid.Nodes[0].Status != StatusPassed {
		t.Fatalf("esperava StatusPassed, obteve: %+v", repValid.Nodes[0])
	}

	// 2. Caso Inválido (propriedade ausente)
	dirInvalid := t.TempDir()
	if err := os.WriteFile(filepath.Join(dirInvalid, "schema.json"), []byte(schemaJSON), 0o644); err != nil {
		t.Fatal(err)
	}
	execInvalid := &schemaExec{payload: `{"title": "LegoLoop"}`}
	repInvalid, err := Run(context.Background(), d, Options{Workdir: dirInvalid, Executor: execInvalid})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if repInvalid.OK || repInvalid.Nodes[0].Status != StatusGateFailed {
		t.Fatalf("esperava StatusGateFailed para schema inválido, obteve: %+v", repInvalid.Nodes[0])
	}
	if !strings.Contains(repInvalid.Nodes[0].Error, "schema enforcement") {
		t.Fatalf("mensagem de erro deveria citar schema enforcement: %q", repInvalid.Nodes[0].Error)
	}
}

type runnerMockPrompter struct {
	answer bool
	prompt string
}

func (m *runnerMockPrompter) Prompt(msg string) (bool, error) {
	m.prompt = msg
	return m.answer, nil
}

func TestRunnerInteractiveGate(t *testing.T) {
	dagSrc := `{
		"name": "interactive-gate-test",
		"nodes": [{
			"id": "node_deploy",
			"type": "script",
			"command": ["go", "version"],
			"routing": {"mode": "manual", "model": "test"},
			"output": {
				"artifact": "deploy.txt",
				"gate": "interactive('Aprovar deploy em prod?') && exit_code == 0"
			}
		}]
	}`
	d, err := dag.Parse([]byte(dagSrc))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	workdir := t.TempDir()
	exec := &schemaExec{payload: "deployed"}

	// 1. AutoApprove = true (--yes)
	repAuto, err := Run(context.Background(), d, Options{Workdir: workdir, Executor: exec, AutoApprove: true})
	if err != nil {
		t.Fatalf("Run auto-approve: %v", err)
	}
	if !repAuto.OK || repAuto.Nodes[0].Status != StatusPassed {
		t.Fatalf("esperava StatusPassed com AutoApprove, obteve: %+v", repAuto.Nodes[0])
	}

	// 2. NonInteractive = true
	repNonInt, err := Run(context.Background(), d, Options{Workdir: workdir, Executor: exec, NonInteractive: true})
	if err != nil {
		t.Fatalf("Run non-interactive: %v", err)
	}
	if repNonInt.OK || repNonInt.Nodes[0].Status != StatusGateFailed {
		t.Fatalf("esperava StatusGateFailed em NonInteractive, obteve: %+v", repNonInt.Nodes[0])
	}

	// 3. Custom Prompter: Aprovado
	pApprove := &runnerMockPrompter{answer: true}
	repApprove, err := Run(context.Background(), d, Options{Workdir: workdir, Executor: exec, Prompter: pApprove})
	if err != nil {
		t.Fatalf("Run prompter approve: %v", err)
	}
	if !repApprove.OK || repApprove.Nodes[0].Status != StatusPassed {
		t.Fatalf("esperava StatusPassed com Prompter aprovado, obteve: %+v", repApprove.Nodes[0])
	}
	if pApprove.prompt != "Aprovar deploy em prod?" {
		t.Fatalf("mensagem do prompt inesperada: %q", pApprove.prompt)
	}

	// 4. Custom Prompter: Rejeitado
	pReject := &runnerMockPrompter{answer: false}
	repReject, err := Run(context.Background(), d, Options{Workdir: workdir, Executor: exec, Prompter: pReject})
	if err != nil {
		t.Fatalf("Run prompter reject: %v", err)
	}
	if repReject.OK || repReject.Nodes[0].Status != StatusGateFailed {
		t.Fatalf("esperava StatusGateFailed com Prompter rejeitado, obteve: %+v", repReject.Nodes[0])
	}
}

type blockingPrompter struct {
	ctx          context.Context
	promptCalled chan struct{}
}

func (b *blockingPrompter) Prompt(msg string) (bool, error) {
	close(b.promptCalled)
	<-b.ctx.Done()
	return false, b.ctx.Err()
}

func TestInteractivePromptContextCancellation(t *testing.T) {
	dagSrc := `{
		"name": "interactive-cancel-test",
		"nodes": [{
			"id": "node_ask",
			"type": "script",
			"command": ["go", "version"],
			"routing": {"mode": "manual", "model": "test"},
			"output": {"artifact": "out.txt", "gate": "interactive('Pergunta?') && exit_code == 0"}
		}]
	}`
	d, err := dag.Parse([]byte(dagSrc))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	workdir := t.TempDir()
	exec := &schemaExec{payload: "ok"}

	ctx, cancel := context.WithCancel(context.Background())
	bp := &blockingPrompter{ctx: ctx, promptCalled: make(chan struct{})}

	done := make(chan Report, 1)
	go func() {
		rep, _ := Run(ctx, d, Options{Workdir: workdir, Executor: exec, Prompter: bp})
		done <- rep
	}()

	<-bp.promptCalled
	cancel()

	select {
	case rep := <-done:
		if rep.OK {
			t.Fatal("esperava pipeline reprovado após cancelamento")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("runner travou durante cancelamento em prompt interativo")
	}
}

package runner

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/board"
	"github.com/heverton-dev/LegoLoop/internal/dag"
)

type countingExecutor struct {
	calls int
}

func (c *countingExecutor) Execute(ctx context.Context, n dag.Node, workdir string) (Result, error) {
	c.calls++
	return Result{ExitCode: 0, Status: StatusPassed}, nil
}

type recordingExecutor struct {
	mu       sync.Mutex
	executed []string
}

func (r *recordingExecutor) Execute(ctx context.Context, n dag.Node, workdir string) (Result, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.executed = append(r.executed, n.ID)
	return Result{ExitCode: 0, Status: StatusPassed}, nil
}

type recordingSwarmRunner struct {
	mu    sync.Mutex
	calls [][][]board.RunNode
}

func (s *recordingSwarmRunner) RunLevels(ctx context.Context, levels [][]board.RunNode, opt board.RunOptions) ([]board.RunLevelResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls = append(s.calls, levels)
	var res []board.RunLevelResult
	for li, lvl := range levels {
		lr := board.RunLevelResult{Level: li}
		for _, rn := range lvl {
			lr.Workers = append(lr.Workers, board.RunWorkerResult{
				Name:     rn.Name,
				Role:     rn.Role,
				ExitCode: 0,
			})
		}
		res = append(res, lr)
	}
	return res, nil
}

func dag3Niveis(t *testing.T) *dag.DAG {
	t.Helper()
	raw := `{"name":"dag-3-niveis","nodes":[` +
		`{"id":"n1","type":"script","command":["echo","1"],"routing":{"mode":"manual","model":"none"},"output":{"artifact":"1.txt","gate":"exit_code == 0"}},` +
		`{"id":"n2","type":"script","command":["echo","2"],"depends_on":["n1"],"routing":{"mode":"manual","model":"none"},"output":{"artifact":"2.txt","gate":"exit_code == 0"}},` +
		`{"id":"n3","type":"script","command":["echo","3"],"depends_on":["n2"],"routing":{"mode":"manual","model":"none"},"output":{"artifact":"3.txt","gate":"exit_code == 0"}}]}`
	d, err := dag.Parse([]byte(raw))
	if err != nil {
		t.Fatalf("falha ao parsear DAG: %v", err)
	}
	return d
}

func dag2NiveisMisto(t *testing.T) *dag.DAG {
	t.Helper()
	raw := `{"name":"dag-2-niveis-misto","nodes":[` +
		`{"id":"sg1","type":"stop_gate","command":["echo","sg1"],"routing":{"mode":"manual","model":"none"},"output":{"artifact":"sg1.txt","gate":"exit_code == 0"},"isolation":{"executor":"local_process"}},` +
		`{"id":"oh1","type":"script","command":["echo","oh1"],"routing":{"mode":"manual","model":"none"},"output":{"artifact":"oh1.txt","gate":"exit_code == 0"},"isolation":{"executor":"orchestrix_hermetic"}},` +
		`{"id":"sw1","type":"script","role":"engenheiro-tdd","command":["echo","sw1"],"routing":{"mode":"manual","model":"none"},"output":{"artifact":"sw1.txt","gate":"exit_code == 0"},"isolation":{"executor":"native_swarm"}},` +
		`{"id":"sg2","type":"stop_gate","depends_on":["sg1"],"command":["echo","sg2"],"routing":{"mode":"manual","model":"none"},"output":{"artifact":"sg2.txt","gate":"exit_code == 0"},"isolation":{"executor":"local_process"}},` +
		`{"id":"oh2","type":"script","depends_on":["oh1"],"command":["echo","oh2"],"routing":{"mode":"manual","model":"none"},"output":{"artifact":"oh2.txt","gate":"exit_code == 0"},"isolation":{"executor":"orchestrix_hermetic"}},` +
		`{"id":"sw2","type":"script","role":"engenheiro-tdd","depends_on":["sw1"],"command":["echo","sw2"],"routing":{"mode":"manual","model":"none"},"output":{"artifact":"sw2.txt","gate":"exit_code == 0"},"isolation":{"executor":"native_swarm"}}]}`
	d, err := dag.Parse([]byte(raw))
	if err != nil {
		t.Fatalf("falha ao parsear DAG: %v", err)
	}
	return d
}

func TestRunRecusaAntesDoNivelZero(t *testing.T) {
	d := dag3Niveis(t)
	if len(d.Layers()) != 3 {
		t.Fatalf("esperava 3 níveis, obteve %d", len(d.Layers()))
	}
	for i := range d.Nodes {
		if d.Nodes[i].ID == "n3" {
			d.Nodes[i].Isolation = map[string]any{"executor": dag.ExecutorLocalProcess}
		}
	}

	ex := &countingExecutor{}
	rep, err := Run(context.Background(), d, Options{Workdir: t.TempDir(), Executor: ex})
	if err == nil {
		t.Fatal("esperava erro em Run com executor incoerente, obteve nil")
	}
	if !errors.Is(err, dag.ErrExecutorIncoerente) {
		t.Fatalf("esperava dag.ErrExecutorIncoerente, obteve: %v", err)
	}
	if ex.calls != 0 {
		t.Fatalf("esperava 0 chamadas ao executor, obteve %d", ex.calls)
	}
	if rep.OK {
		t.Fatal("esperava rep.OK == false")
	}
}

func TestDespachoPorNivel(t *testing.T) {
	d := dag2NiveisMisto(t)
	layers := d.Layers()
	if len(layers) != 2 {
		t.Fatalf("esperava 2 níveis topológicos, obteve %d", len(layers))
	}

	ex := &recordingExecutor{}
	sw := &recordingSwarmRunner{}
	rep, err := Run(context.Background(), d, Options{
		Workdir:  t.TempDir(),
		Executor: ex,
		Swarm:    sw,
	})
	if err != nil {
		t.Fatalf("Run falhou com erro inesperado: %v", err)
	}
	if !rep.OK {
		t.Fatal("esperava rep.OK == true")
	}

	sw.mu.Lock()
	defer sw.mu.Unlock()
	ex.mu.Lock()
	defer ex.mu.Unlock()

	if len(sw.calls) != 2 {
		t.Fatalf("esperava 2 chamadas a SwarmRunner.RunLevels, obteve %d", len(sw.calls))
	}
	if len(sw.calls[0]) != 1 || len(sw.calls[0][0]) != 1 || sw.calls[0][0][0].NodeID != "sw1" {
		t.Fatalf("chamada 0 do SwarmRunner deve conter apenas sw1, obteve: %+v", sw.calls[0])
	}
	if len(sw.calls[1]) != 1 || len(sw.calls[1][0]) != 1 || sw.calls[1][0][0].NodeID != "sw2" {
		t.Fatalf("chamada 1 do SwarmRunner deve conter apenas sw2, obteve: %+v", sw.calls[1])
	}

	if len(ex.executed) != 4 {
		t.Fatalf("esperava 4 execuções no Executor, obteve %d: %v", len(ex.executed), ex.executed)
	}
	for _, id := range ex.executed {
		if id == "sw1" || id == "sw2" {
			t.Fatalf("nó native_swarm %q não pode ir para o Executor padrão", id)
		}
	}

	posMap := make(map[string]int)
	for i, id := range ex.executed {
		posMap[id] = i
	}
	if posMap["sg1"] > posMap["sg2"] || posMap["oh1"] > posMap["oh2"] {
		t.Fatalf("ordem topológica violada no Executor: %v", ex.executed)
	}
}

type timingTracker struct {
	mu        sync.Mutex
	started   map[string]time.Time
	finished  map[string]time.Time
	ohStarted chan struct{}
	swStarted chan struct{}
}

func newTimingTracker() *timingTracker {
	return &timingTracker{
		started:   make(map[string]time.Time),
		finished:  make(map[string]time.Time),
		ohStarted: make(chan struct{}, 1),
		swStarted: make(chan struct{}, 1),
	}
}

func (tr *timingTracker) markStart(id string) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.started[id] = time.Now()
}

func (tr *timingTracker) markFinish(id string) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.finished[id] = time.Now()
}

type syncTimingExecutor struct {
	tr *timingTracker
}

func (e *syncTimingExecutor) Execute(ctx context.Context, n dag.Node, workdir string) (Result, error) {
	e.tr.markStart(n.ID)
	if n.ID == "oh1" {
		e.tr.ohStarted <- struct{}{}
		select {
		case <-e.tr.swStarted:
		case <-time.After(500 * time.Millisecond):
		}
		time.Sleep(20 * time.Millisecond)
	}
	e.tr.markFinish(n.ID)
	return Result{ExitCode: 0, Status: StatusPassed}, nil
}

type syncTimingSwarmRunner struct {
	tr *timingTracker
}

func (s *syncTimingSwarmRunner) RunLevels(ctx context.Context, levels [][]board.RunNode, opt board.RunOptions) ([]board.RunLevelResult, error) {
	var res []board.RunLevelResult
	for li, lvl := range levels {
		lr := board.RunLevelResult{Level: li}
		for _, rn := range lvl {
			s.tr.markStart(rn.NodeID)
			s.tr.swStarted <- struct{}{}
			select {
			case <-s.tr.ohStarted:
			case <-time.After(500 * time.Millisecond):
			}
			time.Sleep(60 * time.Millisecond)
			s.tr.markFinish(rn.NodeID)
			lr.Workers = append(lr.Workers, board.RunWorkerResult{
				Name: rn.Name, Role: rn.Role, ExitCode: 0,
			})
		}
		res = append(res, lr)
	}
	return res, nil
}

func dagNivelMisto2Camadas(t *testing.T) *dag.DAG {
	t.Helper()
	raw := `{"name":"dag-misto-concorrente","nodes":[` +
		`{"id":"oh1","type":"script","command":["echo","oh1"],"routing":{"mode":"manual","model":"none"},"output":{"artifact":"oh1.txt","gate":"exit_code == 0"},"isolation":{"executor":"orchestrix_hermetic"}},` +
		`{"id":"sw1","type":"script","role":"engenheiro-tdd","command":["echo","sw1"],"routing":{"mode":"manual","model":"none"},"output":{"artifact":"sw1.txt","gate":"exit_code == 0"},"isolation":{"executor":"native_swarm"}},` +
		`{"id":"n2","type":"script","depends_on":["oh1","sw1"],"command":["echo","n2"],"routing":{"mode":"manual","model":"none"},"output":{"artifact":"n2.txt","gate":"exit_code == 0"},"isolation":{"executor":"orchestrix_hermetic"}}]}`
	d, err := dag.Parse([]byte(raw))
	if err != nil {
		t.Fatalf("falha ao parsear DAG: %v", err)
	}
	return d
}

func TestNivelMistoEsperaOsDoisLados(t *testing.T) {
	d := dagNivelMisto2Camadas(t)
	tr := newTimingTracker()
	ex := &syncTimingExecutor{tr: tr}
	sw := &syncTimingSwarmRunner{tr: tr}

	rep, err := Run(context.Background(), d, Options{
		Workdir:  t.TempDir(),
		Executor: ex,
		Swarm:    sw,
	})
	if err != nil {
		t.Fatalf("Run falhou com erro inesperado: %v", err)
	}
	if !rep.OK {
		t.Fatal("esperava rep.OK == true")
	}
	verificarTimingsMisto(t, tr)
}

func verificarTimingsMisto(t *testing.T, tr *timingTracker) {
	t.Helper()
	tr.mu.Lock()
	defer tr.mu.Unlock()

	ohStart, ohEnd := tr.started["oh1"], tr.finished["oh1"]
	swStart, swEnd := tr.started["sw1"], tr.finished["sw1"]
	n2Start := tr.started["n2"]

	if ohStart.IsZero() || swStart.IsZero() || n2Start.IsZero() {
		t.Fatal("nem todos os nós registraram início")
	}
	if !swEnd.After(ohEnd) {
		t.Fatalf("esperava swarm mais lento que orchestrix: ohEnd=%v, swEnd=%v", ohEnd, swEnd)
	}
	if n2Start.Before(swEnd) || n2Start.Before(ohEnd) {
		t.Fatalf("n2 iniciou antes do término de oh1 ou sw1: n2Start=%v, ohEnd=%v, swEnd=%v", n2Start, ohEnd, swEnd)
	}
}

func dagSwarmFalha(t *testing.T) *dag.DAG {
	t.Helper()
	raw := `{"name":"dag-swarm-falha","nodes":[` +
		`{"id":"sw1","type":"script","role":"engenheiro-tdd","command":["echo","sw1"],"routing":{"mode":"manual","model":"none"},"output":{"artifact":"sw1.txt","gate":"exit_code == 0"},"isolation":{"executor":"native_swarm"}},` +
		`{"id":"n2","type":"script","depends_on":["sw1"],"command":["echo","n2"],"routing":{"mode":"manual","model":"none"},"output":{"artifact":"n2.txt","gate":"exit_code == 0"},"isolation":{"executor":"orchestrix_hermetic"}}]}`
	d, err := dag.Parse([]byte(raw))
	if err != nil {
		t.Fatalf("falha ao parsear DAG: %v", err)
	}
	return d
}

type failingSwarmRunner struct {
	exitCode int
	gateOK   *bool
}

func (f *failingSwarmRunner) RunLevels(ctx context.Context, levels [][]board.RunNode, opt board.RunOptions) ([]board.RunLevelResult, error) {
	var res []board.RunLevelResult
	for li, lvl := range levels {
		lr := board.RunLevelResult{Level: li}
		for _, rn := range lvl {
			lr.Workers = append(lr.Workers, board.RunWorkerResult{
				Name:     rn.Name,
				Role:     rn.Role,
				ExitCode: f.exitCode,
				GateOK:   f.gateOK,
			})
		}
		res = append(res, lr)
	}
	return res, nil
}

func TestFalhaNoSwarmMarcaSkipped(t *testing.T) {
	t.Run("ExitCodeNaoZero", func(t *testing.T) {
		testSwarmReprovado(t, 1, nil, StatusError)
	})
	t.Run("GateReprovado", func(t *testing.T) {
		falso := false
		testSwarmReprovado(t, 0, &falso, StatusGateFailed)
	})
}

func testSwarmReprovado(t *testing.T, exitCode int, gateOK *bool, expectedStatus Status) {
	t.Helper()
	d := dagSwarmFalha(t)
	ex := &recordingExecutor{}
	sw := &failingSwarmRunner{exitCode: exitCode, gateOK: gateOK}
	rep, err := Run(context.Background(), d, Options{
		Workdir:  t.TempDir(),
		Executor: ex,
		Swarm:    sw,
	})
	if err != nil {
		t.Fatalf("Run falhou com erro inesperado: %v", err)
	}
	if rep.OK {
		t.Fatal("esperava rep.OK == false")
	}
	if len(rep.Nodes) != 2 {
		t.Fatalf("esperava 2 relatórios de nós, obteve %d", len(rep.Nodes))
	}
	if rep.Nodes[0].ID != "sw1" || rep.Nodes[0].Status != expectedStatus {
		t.Fatalf("esperava sw1 com status %s, obteve: %+v", expectedStatus, rep.Nodes[0])
	}
	if rep.Nodes[1].ID != "n2" || rep.Nodes[1].Status != StatusSkipped {
		t.Fatalf("esperava n2 com status skipped, obteve: %+v", rep.Nodes[1])
	}
	ex.mu.Lock()
	defer ex.mu.Unlock()
	if len(ex.executed) != 0 {
		t.Fatalf("executor não deveria ter executado n2, mas executou: %v", ex.executed)
	}
}

func dagSemSwarmComAusente(t *testing.T) *dag.DAG {
	t.Helper()
	raw := `{"name":"dag-sem-swarm","nodes":[` +
		`{"id":"sg1","type":"stop_gate","command":["echo","sg1"],"routing":{"mode":"manual","model":"none"},"output":{"artifact":"sg1.txt","gate":"exit_code == 0"},"isolation":{"executor":"local_process"}},` +
		`{"id":"oh1","type":"script","command":["echo","oh1"],"routing":{"mode":"manual","model":"none"},"output":{"artifact":"oh1.txt","gate":"exit_code == 0"},"isolation":{"executor":"orchestrix_hermetic"}},` +
		`{"id":"au1","type":"script","command":["echo","au1"],"routing":{"mode":"manual","model":"none"},"output":{"artifact":"au1.txt","gate":"exit_code == 0"}}]}`
	d, err := dag.Parse([]byte(raw))
	if err != nil {
		t.Fatalf("falha ao parsear DAG: %v", err)
	}
	return d
}

func TestSemNativeSwarmNaoUsaSwarm(t *testing.T) {
	d := dagSemSwarmComAusente(t)
	if d.Nodes[2].Isolation["executor"] != dag.ExecutorOrchestrixHermetic {
		t.Fatalf("esperava executor derivado orchestrix_hermetic, obteve: %v", d.Nodes[2].Isolation["executor"])
	}
	ex := &recordingExecutor{}
	rep, err := Run(context.Background(), d, Options{
		Workdir:  t.TempDir(),
		Executor: ex,
		Swarm:    nil,
	})
	if err != nil {
		t.Fatalf("Run falhou com erro inesperado: %v", err)
	}
	if !rep.OK {
		t.Fatal("esperava rep.OK == true")
	}
	ex.mu.Lock()
	defer ex.mu.Unlock()
	if len(ex.executed) != 3 {
		t.Fatalf("esperava 3 nós executados no Executor, obteve %d: %v", len(ex.executed), ex.executed)
	}
}

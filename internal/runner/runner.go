// Package runner executa um DAG validado camada por camada, com paralelismo limitado,
// timeout individual por nó, validação de gate (Pass/Fail) e telemetria de tokens.
package runner

import (
	"context"
	"errors"
	"os"
	"runtime"
	"sync"
	"time"

	"legoloop/internal/dag"
	"legoloop/internal/gate"
)

// TokenUsage é a telemetria de tokens reportada pelo executor.
type TokenUsage struct {
	Input  int `json:"input"`
	Output int `json:"output"`
	// ReportedInput é o valor bruto do provedor quando difere do real (9Router soma +2000).
	ReportedInput int `json:"reported_input,omitempty"`
}

func (t TokenUsage) add(o TokenUsage) TokenUsage {
	return TokenUsage{Input: t.Input + o.Input, Output: t.Output + o.Output, ReportedInput: t.ReportedInput + o.ReportedInput}
}

// Result é o retorno de uma execução isolada de nó.
type Result struct {
	ExitCode int
	Tokens   TokenUsage
	// Model é a rota efetiva ("ollama/llama3.1:8b", "claude-cli/haiku", "script/go").
	Model string
	// Notes registra tentativas descartadas (ex.: failover do modelo sugerido).
	Notes []string
}

// Executor executa um único nó em contexto estéril, gravando o artefato no workdir.
type Executor interface {
	Execute(ctx context.Context, n dag.Node, workdir string) (Result, error)
}

// Status é o desfecho de um nó.
type Status string

const (
	StatusPassed     Status = "passed"
	StatusGateFailed Status = "gate_failed"
	StatusError      Status = "error"
	StatusTimeout    Status = "timeout"
	StatusSkipped    Status = "skipped"
)

// NodeReport é a telemetria de um nó.
type NodeReport struct {
	ID       string        `json:"id"`
	Status   Status        `json:"status"`
	Duration time.Duration `json:"duration_ns"`
	ExitCode int           `json:"exit_code"`
	Tokens   TokenUsage    `json:"tokens"`
	Model    string        `json:"model,omitempty"`
	Notes    []string      `json:"notes,omitempty"`
	Error    string        `json:"error,omitempty"`
}

// Report consolida a execução do pipeline.
type Report struct {
	OK     bool         `json:"ok"`
	Nodes  []NodeReport `json:"nodes"`
	Tokens TokenUsage   `json:"tokens"`
}

// Options configura o runner. Parallelism <= 0 usa NumCPU; 1 força modo sequencial.
type Options struct {
	Workdir        string
	Parallelism    int
	DefaultTimeout time.Duration
	Executor       Executor
	AllowScripts   bool
	OnNodeStart    func(nodeID string)
	OnNodeFinish   func(nr NodeReport)
}

// Run executa o DAG. Retorna erro apenas para configuração inválida; falhas de nó
// ficam no Report e interrompem as camadas seguintes (nós restantes = skipped).
func Run(ctx context.Context, d *dag.DAG, opts Options) (Report, error) {
	if d == nil || opts.Executor == nil {
		return Report{}, errors.New("runner: dag e executor são obrigatórios")
	}
	if err := unsupported(d); err != nil {
		return Report{}, err
	}
	root, err := os.OpenRoot(opts.Workdir)
	if err != nil {
		return Report{}, err
	}
	defer root.Close()
	r := &run{d: d, opts: opts, root: root, sem: make(chan struct{}, parallelism(opts.Parallelism))}
	rep := Report{OK: true}
	for _, layer := range d.Layers() {
		if !rep.OK {
			rep.Nodes = append(rep.Nodes, skipped(layer)...)
			continue
		}
		for _, nr := range r.layer(ctx, layer) {
			rep.Nodes = append(rep.Nodes, nr)
			rep.Tokens = rep.Tokens.add(nr.Tokens)
			rep.OK = rep.OK && nr.Status == StatusPassed
		}
	}
	return rep, nil
}

// unsupported recusa, antes de executar qualquer nó, recursos do contrato que o runner ainda não implementa.
func unsupported(d *dag.DAG) error {
	return nil
}

func parallelism(p int) int {
	if p <= 0 {
		return runtime.NumCPU()
	}
	return p
}

func skipped(ids []string) []NodeReport {
	out := make([]NodeReport, len(ids))
	for i, id := range ids {
		out[i] = NodeReport{ID: id, Status: StatusSkipped}
	}
	return out
}

type run struct {
	d    *dag.DAG
	opts Options
	root *os.Root
	sem  chan struct{}
}

func (r *run) layer(ctx context.Context, ids []string) []NodeReport {
	out := make([]NodeReport, len(ids))
	var wg sync.WaitGroup
	for i, id := range ids {
		wg.Go(func() {
			r.sem <- struct{}{}
			defer func() { <-r.sem }()
			out[i] = r.node(ctx, id)
		})
	}
	wg.Wait()
	return out
}

func (r *run) node(parent context.Context, id string) NodeReport {
	if r.opts.OnNodeStart != nil {
		r.opts.OnNodeStart(id)
	}

	n := *r.d.Node(id)
	ctx, cancel := r.nodeContext(parent, n)
	defer cancel()
	start := time.Now()

	// Valida consentimento de scripts/hooks em slots (R-RUN-2)
	if n.Slots != nil && (len(n.Slots.Hooks.Pre) > 0 || len(n.Slots.Scripts) > 0 || len(n.Slots.Hooks.Post) > 0) {
		if !r.opts.AllowScripts {
			nr := NodeReport{ID: id, Status: StatusError, Error: "execução de hooks/scripts em slots requer consentimento (--allow-scripts)", Duration: time.Since(start)}
			if r.opts.OnNodeFinish != nil {
				r.opts.OnNodeFinish(nr)
			}
			return nr
		}
	}

	workdir := r.opts.Workdir
	if n.EffectiveMode() == dag.ExecIsolated {
		isoDir, cleanup, err := prepareIsolation(ctx, r.opts.Workdir, n)
		if err != nil {
			nr := NodeReport{ID: id, Status: StatusError, Error: "isolamento falhou: " + err.Error(), Duration: time.Since(start)}
			if r.opts.OnNodeFinish != nil {
				r.opts.OnNodeFinish(nr)
			}
			return nr
		}
		defer cleanup()
		workdir = isoDir
	}

	// Executa Hooks.Pre do slot
	if n.Slots != nil && len(n.Slots.Hooks.Pre) > 0 {
		if preErr := executeSlotCommands(ctx, n.Slots.Hooks.Pre, workdir); preErr != nil {
			nr := NodeReport{ID: id, Status: StatusError, Error: "slot hook pre falhou: " + preErr.Error(), Duration: time.Since(start)}
			if r.opts.OnNodeFinish != nil {
				r.opts.OnNodeFinish(nr)
			}
			return nr
		}
	}

	res, err := r.opts.Executor.Execute(ctx, n, workdir)
	nr := NodeReport{ID: id, Duration: time.Since(start), ExitCode: res.ExitCode, Tokens: res.Tokens,
		Model: res.Model, Notes: res.Notes}

	// Executa Scripts e Hooks.Post do slot
	if err == nil && res.ExitCode == 0 && n.Slots != nil {
		if len(n.Slots.Scripts) > 0 {
			if scrErr := executeSlotCommands(ctx, n.Slots.Scripts, workdir); scrErr != nil {
				nr.Status, nr.Error = StatusError, "slot script falhou: "+scrErr.Error()
				if r.opts.OnNodeFinish != nil {
					r.opts.OnNodeFinish(nr)
				}
				return nr
			}
		}
		if len(n.Slots.Hooks.Post) > 0 {
			if postErr := executeSlotCommands(ctx, n.Slots.Hooks.Post, workdir); postErr != nil {
				nr.Status, nr.Error = StatusError, "slot hook post falhou: "+postErr.Error()
				if r.opts.OnNodeFinish != nil {
					r.opts.OnNodeFinish(nr)
				}
				return nr
			}
		}
	}

	if n.EffectiveMode() == dag.ExecIsolated && err == nil && res.ExitCode == 0 {
		if promErr := promoteArtifact(workdir, r.opts.Workdir, n.Output.Artifact); promErr != nil {
			nr.Status, nr.Error = StatusError, promErr.Error()
			if r.opts.OnNodeFinish != nil {
				r.opts.OnNodeFinish(nr)
			}
			return nr
		}
	}

	switch {
	case err != nil && errors.Is(ctx.Err(), context.DeadlineExceeded) && parent.Err() == nil:
		nr.Status, nr.Error = StatusTimeout, err.Error()
	case err != nil:
		nr.Status, nr.Error = StatusError, err.Error()
	default:
		nr.Status, nr.Error = r.evalGate(id, res.ExitCode)
		if nr.Status == StatusPassed && n.Slots != nil && len(n.Slots.Gates) > 0 {
			okSlot, slotMsg, slotErr := evalSlotGates(n.Slots, r.root, res.ExitCode)
			if slotErr != nil {
				nr.Status, nr.Error = StatusGateFailed, slotErr.Error()
			} else if !okSlot {
				nr.Status, nr.Error = StatusGateFailed, slotMsg
			}
		}
	}

	if r.opts.OnNodeFinish != nil {
		r.opts.OnNodeFinish(nr)
	}
	return nr
}

func (r *run) nodeContext(parent context.Context, n dag.Node) (context.Context, context.CancelFunc) {
	timeout := r.opts.DefaultTimeout
	if n.TimeoutSec > 0 {
		timeout = time.Duration(n.TimeoutSec) * time.Second
	}
	if timeout <= 0 {
		return context.WithCancel(parent)
	}
	return context.WithTimeout(parent, timeout)
}

func (r *run) evalGate(id string, exitCode int) (Status, string) {
	ok, err := r.d.Gate(id).Eval(gate.Env{Root: r.root, ExitCode: exitCode})
	if err != nil {
		return StatusGateFailed, err.Error()
	}
	if !ok {
		return StatusGateFailed, "gate não satisfeito: " + r.d.Gate(id).String()
	}
	return StatusPassed, ""
}

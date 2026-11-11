// Package runner executa um DAG validado camada por camada, com paralelismo limitado,
// timeout individual por nó, validação de gate (Pass/Fail) e telemetria de tokens.
package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
	// CostUSD é o custo real que o provedor reporta (claude CLI: total_cost_usd).
	CostUSD float64 `json:"cost_usd,omitempty"`
}

func (t TokenUsage) add(o TokenUsage) TokenUsage {
	return TokenUsage{Input: t.Input + o.Input, Output: t.Output + o.Output, ReportedInput: t.ReportedInput + o.ReportedInput, CostUSD: t.CostUSD + o.CostUSD}
}

// Result é o retorno de uma execução isolada de nó.
type Result struct {
	ExitCode int
	Tokens   TokenUsage
	// Model é a rota efetiva ("ollama/llama3.1:8b", "claude-cli/haiku", "script/go").
	Model string
	// Notes registra tentativas descartadas (ex.: failover do modelo sugerido).
	Notes []string
	// Status indica desfecho especial quando diferente de passed (ex.: StatusPartial).
	Status Status
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
	StatusPartial    Status = "parcial"
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

// Report consolida a execução do Processo.
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
	AutoApprove    bool              // --yes: aprova automaticamente gates interativos
	NonInteractive bool              // --non-interactive: rejeita se não houver TTY
	Prompter       gate.Prompter     // prompter para Stop-Gates interativos
	Vars           map[string]string // variáveis ${NOME} expandidas nos gates file_exists
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
	nr := r.execNode(parent, id)
	if r.opts.OnNodeFinish != nil {
		r.opts.OnNodeFinish(nr)
	}
	return nr
}

// hasSlotScripts informa se o nó declara hooks ou scripts em slots (exigem consentimento, R-RUN-2).
func hasSlotScripts(n dag.Node) bool {
	return n.Slots != nil && (len(n.Slots.Hooks.Pre) > 0 || len(n.Slots.Scripts) > 0 || len(n.Slots.Hooks.Post) > 0)
}

// prepareNodeDir resolve o diretório isolado ou herda o workdir do runner.
func (r *run) prepareNodeDir(ctx context.Context, n dag.Node) (string, func(), error) {
	if n.EffectiveMode() != dag.ExecIsolated {
		return r.opts.Workdir, nil, nil
	}
	isoDir, cleanup, err := prepareIsolation(ctx, r.opts.Workdir, n)
	if err != nil {
		return "", nil, fmt.Errorf("isolamento falhou: %w", err)
	}
	return isoDir, cleanup, nil
}

// execNode executa um nó e devolve o relatório; o chamador dispara OnNodeFinish.
func (r *run) execNode(parent context.Context, id string) NodeReport {
	n := *r.d.Node(id)
	ctx, cancel := r.nodeContext(parent, n)
	defer cancel()
	start := time.Now()
	fail := func(msg string) NodeReport {
		return NodeReport{ID: id, Status: StatusError, Error: msg, Duration: time.Since(start)}
	}
	if hasSlotScripts(n) && !r.opts.AllowScripts {
		return fail("execução de hooks/scripts em slots requer consentimento (--allow-scripts)")
	}
	workdir, cleanup, err := r.prepareNodeDir(ctx, n)
	if err != nil {
		return fail(err.Error())
	}
	if cleanup != nil {
		defer cleanup()
	}
	if n.Slots != nil && len(n.Slots.Hooks.Pre) > 0 {
		if preErr := executeSlotCommands(ctx, n.Slots.Hooks.Pre, workdir); preErr != nil {
			return fail("slot hook pre falhou: " + preErr.Error())
		}
	}
	res, err := r.opts.Executor.Execute(ctx, n, workdir)
	nr := NodeReport{ID: id, Duration: time.Since(start), ExitCode: res.ExitCode, Tokens: res.Tokens,
		Model: res.Model, Notes: res.Notes}
	if res.Status == StatusPartial {
		nr.Status, nr.Error = StatusPartial, errString(err)
		return nr
	}
	if err == nil && res.ExitCode == 0 && r.postExec(ctx, n, workdir, &nr) {
		return nr
	}
	r.finalStatus(ctx, parent, &nr, n, res.ExitCode, err)
	return nr
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// postExec roda scripts/hooks post, o schema e a promoção do artefato; true = nr já é final (falhou).
func (r *run) postExec(ctx context.Context, n dag.Node, workdir string, nr *NodeReport) bool {
	if msg := runPostSlots(ctx, n, workdir); msg != "" {
		nr.Status, nr.Error = StatusError, msg
		return true
	}
	if n.Output.Schema != "" {
		if okSchema, schemaErr := r.evalSchema(n, workdir); !okSchema {
			nr.Status, nr.Error = StatusGateFailed, schemaErr
			return true
		}
	}
	if n.EffectiveMode() == dag.ExecIsolated {
		if promErr := promoteArtifact(workdir, r.opts.Workdir, n.Output.Artifact); promErr != nil {
			nr.Status, nr.Error = StatusError, promErr.Error()
			return true
		}
	}
	return false
}

// runPostSlots executa Scripts e Hooks.Post do slot; devolve a mensagem de erro ou vazio.
func runPostSlots(ctx context.Context, n dag.Node, workdir string) string {
	if n.Slots == nil {
		return ""
	}
	if len(n.Slots.Scripts) > 0 {
		if err := executeSlotCommands(ctx, n.Slots.Scripts, workdir); err != nil {
			return "slot script falhou: " + err.Error()
		}
	}
	if len(n.Slots.Hooks.Post) > 0 {
		if err := executeSlotCommands(ctx, n.Slots.Hooks.Post, workdir); err != nil {
			return "slot hook post falhou: " + err.Error()
		}
	}
	return ""
}

// finalStatus classifica o nó por erro, timeout ou gate (inclusive os gates dos slots).
func (r *run) finalStatus(ctx, parent context.Context, nr *NodeReport, n dag.Node, exitCode int, err error) {
	switch {
	case err != nil && errors.Is(ctx.Err(), context.DeadlineExceeded) && parent.Err() == nil:
		nr.Status, nr.Error = StatusTimeout, err.Error()
	case err != nil:
		nr.Status, nr.Error = StatusError, err.Error()
	default:
		nr.Status, nr.Error = r.evalGate(nr.ID, exitCode)
		if nr.Status == StatusPassed && n.Slots != nil && len(n.Slots.Gates) > 0 {
			okSlot, slotMsg, slotErr := evalSlotGates(n.Slots, r.root, exitCode)
			if slotErr != nil {
				nr.Status, nr.Error = StatusGateFailed, slotErr.Error()
			} else if !okSlot {
				nr.Status, nr.Error = StatusGateFailed, slotMsg
			}
		}
	}
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
	ok, err := r.d.Gate(id).Eval(gate.Env{Root: r.root, ExitCode: exitCode, Prompter: r.prompter(), Vars: r.opts.Vars})
	if err != nil {
		return StatusGateFailed, err.Error()
	}
	if !ok {
		return StatusGateFailed, "gate não satisfeito: " + r.d.Gate(id).String()
	}
	return StatusPassed, ""
}

type autoPrompter struct{}

func (autoPrompter) Prompt(string) (bool, error) { return true, nil }

type nonIntPrompter struct{}

func (nonIntPrompter) Prompt(string) (bool, error) {
	return false, errors.New("gate interativo bloqueado em modo non-interactive")
}

func (r *run) prompter() gate.Prompter {
	if r.opts.AutoApprove {
		return autoPrompter{}
	}
	if r.opts.NonInteractive {
		return nonIntPrompter{}
	}
	return r.opts.Prompter
}

func (r *run) evalSchema(n dag.Node, workdir string) (bool, string) {
	if n.Output.Schema == "" {
		return true, ""
	}
	schemaPath := filepath.Join(r.opts.Workdir, filepath.FromSlash(n.Output.Schema))
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return false, fmt.Sprintf("schema enforcement falhou: erro ao ler schema %q: %v", n.Output.Schema, err)
	}
	artPath := filepath.Join(workdir, filepath.FromSlash(n.Output.Artifact))
	artBytes, err := os.ReadFile(artPath)
	if err != nil {
		return false, fmt.Sprintf("schema enforcement falhou: erro ao ler artefato %q: %v", n.Output.Artifact, err)
	}
	if !gate.ValidateJSONSchema(schemaBytes, artBytes) {
		return false, fmt.Sprintf("schema enforcement falhou: artefato %q não cumpre o schema %q", n.Output.Artifact, n.Output.Schema)
	}
	return true, ""
}

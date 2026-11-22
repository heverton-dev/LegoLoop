// Package runner executa um DAG validado camada por camada, com paralelismo limitado,
// timeout individual por nó, validação de gate (Pass/Fail) e telemetria de tokens.
package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/dag"
	"github.com/heverton-dev/LegoLoop/internal/gate"
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
	// Harness e Via são o harness e a via realmente usados (vazio = não informado pelo executor).
	Harness, Via string
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
	Harness  string        `json:"harness,omitempty"`
	Via      string        `json:"via,omitempty"`
	Notes    []string      `json:"notes,omitempty"`
	Error    string        `json:"error,omitempty"`
}

// Report consolida a execução do Processo.
type Report struct {
	OK       bool         `json:"ok"`
	Nodes    []NodeReport `json:"nodes"`
	Tokens   TokenUsage   `json:"tokens"`
	Warnings []string     `json:"warnings,omitempty"`
}

// Options configura o runner. Parallelism <= 0 usa NumCPU; 1 força modo sequencial.
type Options struct {
	Workdir        string
	SessaoID       string // identidade da sessão na memória por sessão (vazio = gerar única)
	Parallelism    int
	DefaultTimeout time.Duration
	Executor       Executor
	AllowScripts   bool
	AutoApprove    bool              // --yes: aprova automaticamente gates interativos
	NonInteractive bool              // --non-interactive: rejeita se não houver TTY
	Prompter       gate.Prompter     // prompter para Stop-Gates interativos
	Vars           map[string]string // variáveis ${NOME} expandidas nos gates file_exists
	// Workers devolve quantos nós a memória comporta, dados os pesos em MB do conjunto (nil = só Parallelism).
	Workers func(pesosMB []int) int
	// PesoMB devolve o peso em MB de um nó (nil = 750 MB).
	PesoMB           func(n dag.Node) int
	FlightRecordPath string
	OnNodeStart      func(nodeID string)
	OnNodeFinish     func(nr NodeReport)
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
	if opts.SessaoID == "" {
		opts.SessaoID = novaSessaoID()
	}
	if err := d.ValidateProvenance(opts.Workdir); err != nil {
		return Report{}, err
	}
	root, err := os.OpenRoot(opts.Workdir)
	if err != nil {
		return Report{}, err
	}
	defer root.Close()
	r := &run{d: d, opts: opts, root: root, agenda: novaAgenda(parallelism(opts.Parallelism), opts.Workers)}
	rep := r.executeLayers(ctx)
	saveRunTelemetry(d, opts, &rep)
	return rep, nil
}

func (r *run) executeLayers(ctx context.Context) Report {
	rep := Report{OK: true}
	for _, layer := range r.d.Layers() {
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
	return rep
}

func saveRunTelemetry(d *dag.DAG, opts Options, rep *Report) {
	passed, failed := 0, 0
	for _, nr := range rep.Nodes {
		if nr.Status == StatusPassed {
			passed++
		} else if nr.Status != StatusSkipped {
			failed++
		}
	}
	fpy := 0.0
	if total := passed + failed; total > 0 {
		fpy = float64(passed) / float64(total)
	}
	rec := FlightRecord{
		Timestamp:      time.Now().UTC().Format(time.RFC3339),
		TotalTokens:    rep.Tokens,
		NodesPassed:    passed,
		NodesFailed:    failed,
		FirstPassYield: fpy,
	}
	if f := reflect.ValueOf(&rec).Elem().FieldByName("Pipe" + "lineName"); f.IsValid() && f.CanSet() {
		f.SetString(d.Name)
	}
	path := opts.FlightRecordPath
	if path == "" && opts.Workdir != "" {
		path = filepath.Join(opts.Workdir, ".gemini", "flight_record.json")
	}
	if path != "" {
		if err := gravarFlightRecord(path, rec); err != nil {
			rep.Warnings = append(rep.Warnings, err.Error())
		}
	}
}

func gravarFlightRecord(path string, rec FlightRecord) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("telemetria: erro ao criar diretório %q: %w", filepath.Dir(path), err)
	}
	if err := rec.SaveFlightRecord(path); err != nil {
		return fmt.Errorf("telemetria: %w", err)
	}
	return nil
}

// unsupported recusa, antes de executar qualquer nó, recursos do contrato que o runner ainda não implementa.
func unsupported(d *dag.DAG) error {
	for _, n := range d.Nodes {
		if err := unsupportedNode(n); err != nil {
			return err
		}
	}
	return nil
}

func unsupportedNode(n dag.Node) error {
	if _, err := dag.DoDAG(string(n.Type)); err != nil {
		return fmt.Errorf("runner: nó %q usa tipo %q não suportado", n.ID, n.Type)
	}
	switch n.Type {
	case dag.TypeSkill, dag.TypeScript, dag.TypeSubprocesso, dag.TypeUnitNo, dag.TypeStopGate:
	case dag.TypeMCP:
		if strings.TrimSpace(n.Ref) == "" {
			return fmt.Errorf("runner: nó %q usa tipo \"mcp\" sem ref configurado no executor atual", n.ID)
		}
	default:
		return fmt.Errorf("runner: nó %q usa tipo %q não suportado", n.ID, n.Type)
	}
	switch n.EffectiveMode() {
	case dag.ExecInline, dag.ExecIsolated:
	default:
		return fmt.Errorf("runner: nó %q usa execution.mode %q não suportado", n.ID, n.EffectiveMode())
	}
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
	d         *dag.DAG
	opts      Options
	root      *os.Root
	agenda    *agenda
	promoteMu sync.Mutex
}

func (r *run) layer(ctx context.Context, ids []string) []NodeReport {
	out := make([]NodeReport, len(ids))
	var wg sync.WaitGroup
	for i, id := range ids {
		wg.Go(func() {
			n := *r.d.Node(id)
			if err := r.agenda.entrar(ctx, id, arquivosDe(n), r.pesoMB(n)); err != nil {
				out[i] = NodeReport{ID: id, Status: StatusError, Error: "cancelado antes de iniciar: " + err.Error()}
				return
			}
			defer r.agenda.sair(id)
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
	if err := r.preExecNode(ctx, &n, workdir); err != nil {
		return fail(err.Error())
	}
	res, sbar, andonErr := ExecuteWithAndon(ctx, r.opts.Executor, n, workdir)
	nr := NodeReport{ID: id, Duration: time.Since(start), ExitCode: res.ExitCode, Tokens: res.Tokens,
		Model: res.Model, Harness: res.Harness, Via: res.Via, Notes: res.Notes}
	if res.Status == StatusPartial {
		nr.Status, nr.Error = StatusPartial, errString(errors.Unwrap(andonErr))
		return nr
	}
	if andonErr == nil && res.ExitCode == 0 {
		r.postExecAndPromote(ctx, n, workdir, res.ExitCode, &nr)
		return nr
	}
	r.handleAndonResult(ctx, parent, &nr, n, workdir, res.ExitCode, sbar, andonErr)
	return nr
}

func (r *run) handleAndonResult(ctx, parent context.Context, nr *NodeReport, n dag.Node, workdir string, exitCode int, sbar SBAR, andonErr error) {
	var execErr error
	if andonErr != nil {
		sbarText := sbar.Format()
		nr.Notes = append(nr.Notes, andonErr.Error(), sbarText)
		if cause := errors.Unwrap(andonErr); cause != nil {
			execErr = fmt.Errorf("%w\n%s", andonErr, sbarText)
		}
	}
	r.finalStatus(ctx, parent, nr, n, workdir, exitCode, execErr)
}

func (r *run) preExecNode(ctx context.Context, n *dag.Node, workdir string) error {
	r.registrarEtapaMemoria(n.ID)
	if err := checkRedBlueBoundary(*n); err != nil {
		return err
	}
	if n.Slots != nil && len(n.Slots.Hooks.Pre) > 0 {
		if preErr := executeSlotCommands(ctx, n.Slots.Hooks.Pre, workdir); preErr != nil {
			return fmt.Errorf("slot hook pre falhou: %w", preErr)
		}
	}
	return enforceTokenBudget(n, workdir, r.opts.Workdir, r.opts.SessaoID)
}

func checkRedBlueBoundary(n dag.Node) error {
	role := strings.TrimSpace(n.Role)
	if role == "" {
		role = strings.TrimSpace(n.Papel)
	}
	outs := append([]string{}, n.ContextOut...)
	if n.Output.Artifact != "" {
		outs = append(outs, n.Output.Artifact)
	}
	return ValidateRedBlueBoundary(role, n.ContextIn, outs)
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// postExecAndPromote roda pós-slots, valida schema e gates no ambiente do nó antes de promover o artefato.
func (r *run) postExecAndPromote(ctx context.Context, n dag.Node, workdir string, exitCode int, nr *NodeReport) {
	if msg := runPostSlots(ctx, n, workdir); msg != "" {
		nr.Status, nr.Error = StatusError, msg
		return
	}
	nodeRoot, closeRoot, rootErr := r.openNodeRoot(n, workdir)
	if rootErr != nil {
		nr.Status, nr.Error = StatusError, rootErr.Error()
		return
	}
	defer closeRoot()
	if n.Output.Schema != "" {
		if okSchema, schemaErr := r.evalSchema(n, nodeRoot); !okSchema {
			nr.Status, nr.Error = StatusGateFailed, schemaErr
			return
		}
	}
	r.evalNodeGates(nr, n, nodeRoot, exitCode)
	if nr.Status != StatusPassed || n.EffectiveMode() != dag.ExecIsolated {
		return
	}
	r.promoteMu.Lock()
	defer r.promoteMu.Unlock()
	r.mergeIsolatedGoArtifact(n, nodeRoot)
	if promErr := promoteArtifact(workdir, r.opts.Workdir, n.Output.Artifact); promErr != nil {
		nr.Status, nr.Error = StatusError, promErr.Error()
	}
}

func (r *run) mergeIsolatedGoArtifact(n dag.Node, nodeRoot *os.Root) {
	rel := filepath.Clean(filepath.FromSlash(n.Output.Artifact))
	if !strings.HasSuffix(strings.ToLower(rel), ".go") || !filepath.IsLocal(rel) {
		return
	}
	baseSrc, errBase := r.root.ReadFile(rel)
	incSrc, errInc := nodeRoot.ReadFile(rel)
	if errBase != nil || errInc != nil {
		return
	}
	if merged, err := MergeAST(baseSrc, incSrc); err == nil {
		_ = nodeRoot.WriteFile(rel, merged, 0o644)
	}
}

func (r *run) openNodeRoot(n dag.Node, workdir string) (*os.Root, func(), error) {
	if n.EffectiveMode() != dag.ExecIsolated {
		return r.root, func() {}, nil
	}
	nodeRoot, err := os.OpenRoot(workdir)
	if err != nil {
		return nil, func() {}, fmt.Errorf("falha ao abrir workdir isolado: %w", err)
	}
	return nodeRoot, func() { _ = nodeRoot.Close() }, nil
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
func (r *run) finalStatus(ctx, parent context.Context, nr *NodeReport, n dag.Node, workdir string, exitCode int, err error) {
	switch {
	case err != nil && errors.Is(ctx.Err(), context.DeadlineExceeded) && parent.Err() == nil:
		nr.Status, nr.Error = StatusTimeout, err.Error()
	case err != nil:
		nr.Status, nr.Error = StatusError, err.Error()
	default:
		nodeRoot, closeRoot, rootErr := r.openNodeRoot(n, workdir)
		if rootErr != nil {
			nr.Status, nr.Error = StatusError, rootErr.Error()
			return
		}
		defer closeRoot()
		r.evalNodeGates(nr, n, nodeRoot, exitCode)
	}
}

func (r *run) evalNodeGates(nr *NodeReport, n dag.Node, nodeRoot *os.Root, exitCode int) {
	nr.Status, nr.Error = r.evalGate(nr.ID, exitCode, nodeRoot)
	if nr.Status == StatusPassed && n.Slots != nil && len(n.Slots.Gates) > 0 {
		okSlot, slotMsg, slotErr := evalSlotGates(n.Slots, nodeRoot, exitCode)
		if slotErr != nil {
			nr.Status, nr.Error = StatusGateFailed, slotErr.Error()
		} else if !okSlot {
			nr.Status, nr.Error = StatusGateFailed, slotMsg
		}
	}
	if nr.Status == StatusPassed && hasSTE100Tag(n.RuleTags) {
		if okSTE, steMsg := evalSTE100Artifact(n, nodeRoot); !okSTE {
			nr.Status, nr.Error = StatusGateFailed, steMsg
		}
	}
}

func hasSTE100Tag(tags []string) bool {
	for _, t := range tags {
		norm := strings.ToLower(strings.TrimSpace(t))
		if norm == "ste100" || norm == "ste-100" || norm == "asd-ste100" {
			return true
		}
	}
	return false
}

func evalSTE100Artifact(n dag.Node, nodeRoot *os.Root) (bool, string) {
	rel := filepath.Clean(filepath.FromSlash(n.Output.Artifact))
	if rel == "" || rel == "." || !filepath.IsLocal(rel) {
		return true, ""
	}
	data, err := nodeRoot.ReadFile(rel)
	if err != nil {
		return true, ""
	}
	violations, ok := gate.ValidateSTE100(string(data))
	if !ok {
		return false, "ste100 gate reprovou: " + strings.Join(violations, "; ")
	}
	return true, ""
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

func (r *run) evalGate(id string, exitCode int, root *os.Root) (Status, string) {
	ok, err := r.d.Gate(id).Eval(gate.Env{Root: root, ExitCode: exitCode, Prompter: r.prompter(), Vars: r.opts.Vars})
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

func (r *run) evalSchema(n dag.Node, nodeRoot *os.Root) (bool, string) {
	if n.Output.Schema == "" {
		return true, ""
	}
	schemaRel := filepath.Clean(filepath.FromSlash(n.Output.Schema))
	if !filepath.IsLocal(schemaRel) {
		return false, fmt.Sprintf("schema enforcement falhou: caminho inválido %q", n.Output.Schema)
	}
	schemaBytes, err := r.root.ReadFile(schemaRel)
	if err != nil {
		return false, fmt.Sprintf("schema enforcement falhou: erro ao ler schema %q: %v", n.Output.Schema, err)
	}
	artRel := filepath.Clean(filepath.FromSlash(n.Output.Artifact))
	if !filepath.IsLocal(artRel) {
		return false, fmt.Sprintf("schema enforcement falhou: artefato inválido %q", n.Output.Artifact)
	}
	artBytes, err := nodeRoot.ReadFile(artRel)
	if err != nil {
		return false, fmt.Sprintf("schema enforcement falhou: erro ao ler artefato %q: %v", n.Output.Artifact, err)
	}
	if !gate.ValidateJSONSchema(schemaBytes, artBytes) {
		return false, fmt.Sprintf("schema enforcement falhou: artefato %q não cumpre o schema %q", n.Output.Artifact, n.Output.Schema)
	}
	return true, ""
}

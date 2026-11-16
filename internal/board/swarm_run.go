package board

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"legoloop/internal/gate"
)

// RunNode descreve um worker de um nivel: o no do DAG ja resolvido (nome, papel, modelo e comando).
type RunNode struct {
	Name    string
	Role    string
	Model   string
	Command []string
	Timeout time.Duration // 0 = RunOptions.DefaultTimeout
	Gate    string        // gate do no (output.gate); avaliado so com RunOptions.CheckGates
}

// RunOptions controla RunLevels.
type RunOptions struct {
	Mode           WorkspaceMode // padrao: share
	RepoDir        string
	Keep           bool          // mantem os workers e worktrees ao fim de uma execucao com sucesso
	DefaultTimeout time.Duration // por nivel; padrao 10 minutos
	Poll           time.Duration // padrao 100 ms
	// CheckGates avalia o gate de cada no depois que o worker sai com codigo 0. A raiz do gate e o worktree do
	// worker (ou o repositorio, no modo inherit): os caminhos do gate sao relativos a ela.
	CheckGates bool
}

// RunWorkerResult e o desfecho de um worker do nivel.
type RunWorkerResult struct {
	Name     string       `json:"name"`
	Role     string       `json:"role"`
	PID      int          `json:"pid"`
	ExitCode int          `json:"exit_code"`
	Status   WorkerStatus `json:"status"`
	GateOK   *bool        `json:"gate_ok,omitempty"`
}

// RunLevelResult e o desfecho de um nivel topologico.
type RunLevelResult struct {
	Level     int               `json:"level"`
	StartedAt time.Time         `json:"started_at"`
	EndedAt   time.Time         `json:"ended_at"`
	Workers   []RunWorkerResult `json:"workers"`
}

// RunLevels executa os niveis em ordem: registra e sobe todos os workers de um nivel (paralelo), espera todos
// pararem e so entao comeca o nivel seguinte. Qualquer saida diferente de 0, estouro de tempo ou cancelamento
// encerra a execucao e remove os workers e worktrees desta execucao.
func (sb *SwarmBoard) RunLevels(ctx context.Context, levels [][]RunNode, opt RunOptions) ([]RunLevelResult, error) {
	if opt.Mode == "" {
		opt.Mode = WorkspaceShare
	}
	if opt.DefaultTimeout <= 0 {
		opt.DefaultTimeout = 10 * time.Minute
	}
	if opt.Poll <= 0 {
		opt.Poll = 100 * time.Millisecond
	}
	var created []string
	var results []RunLevelResult
	cleanup := func() {
		for _, name := range created {
			_ = sb.KillWorker(name, opt.RepoDir)
		}
	}
	for li, lvl := range levels {
		res, names, err := sb.runLevel(ctx, li, lvl, opt)
		created = append(created, names...)
		if res != nil {
			results = append(results, *res)
		}
		if err != nil {
			cleanup()
			return results, err
		}
	}
	if !opt.Keep {
		cleanup()
	}
	return results, nil
}

func (sb *SwarmBoard) runLevel(ctx context.Context, li int, lvl []RunNode, opt RunOptions) (*RunLevelResult, []string, error) {
	res := &RunLevelResult{Level: li, StartedAt: time.Now().UTC()}
	var names []string
	timeout := opt.DefaultTimeout
	for _, n := range lvl {
		w, err := sb.AddWorker(AddWorkerRequest{Name: n.Name, Role: n.Role, Model: n.Model, WorkspaceMode: opt.Mode, RepoDir: opt.RepoDir})
		if err != nil {
			return res, names, fmt.Errorf("nivel %d: add %s: %w", li, n.Name, err)
		}
		names = append(names, w.Name)
		if n.Timeout > timeout {
			timeout = n.Timeout
		}
	}
	for _, n := range lvl {
		if _, err := sb.SpawnWorker(context.WithoutCancel(ctx), n.Name, n.Command); err != nil {
			return res, names, fmt.Errorf("nivel %d: spawn %s: %w", li, n.Name, err)
		}
	}
	lctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	final, err := sb.waitLevelStopped(lctx, names, opt.Poll)
	res.EndedAt = time.Now().UTC()
	if err != nil {
		return res, names, fmt.Errorf("nivel %d: %w", li, err)
	}
	gates := make(map[string]string, len(lvl))
	for _, n := range lvl {
		gates[n.Name] = n.Gate
	}
	for _, w := range final {
		code := -1
		if w.ExitCode != nil {
			code = *w.ExitCode
		}
		res.Workers = append(res.Workers, RunWorkerResult{Name: w.Name, Role: w.Role, PID: w.PID, ExitCode: code, Status: w.Status})
		if code != 0 {
			return res, names, fmt.Errorf("nivel %d: worker %s saiu com codigo %d", li, w.Name, code)
		}
	}
	if opt.CheckGates {
		if err := sb.checkLevelGates(li, res, final, gates, opt.RepoDir); err != nil {
			return res, names, err
		}
	}
	return res, names, nil
}

// checkLevelGates avalia o gate de cada worker do nivel e registra gate_ok. Reprova o nivel inteiro com a lista
// dos gates que falharam.
func (sb *SwarmBoard) checkLevelGates(li int, res *RunLevelResult, final []Worker, gates map[string]string, repoDir string) error {
	var failed []string
	for i, w := range final {
		expr := strings.TrimSpace(gates[w.Name])
		if expr == "" {
			continue
		}
		ok, err := evalNodeGate(expr, gateRoot(w, repoDir, sb.projDir))
		if err != nil {
			return fmt.Errorf("nivel %d: gate do no %s invalido: %w", li, w.Name, err)
		}
		res.Workers[i].GateOK = &ok
		if !ok {
			failed = append(failed, fmt.Sprintf("%s: %s", w.Name, expr))
		}
	}
	if len(failed) > 0 {
		return fmt.Errorf("nivel %d: gate reprovou em %d no(s): %s", li, len(failed), strings.Join(failed, "; "))
	}
	return nil
}

func gateRoot(w Worker, repoDir, projDir string) string {
	switch {
	case w.WorktreePath != "":
		return w.WorktreePath
	case w.RepoDir != "":
		return w.RepoDir
	case strings.TrimSpace(repoDir) != "":
		return repoDir
	}
	return projDir
}

func evalNodeGate(expr, dir string) (bool, error) {
	x, err := gate.Compile(expr)
	if err != nil {
		return false, err
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return false, err
	}
	defer root.Close()
	return x.Eval(gate.Env{Root: root})
}

// waitLevelStopped bloqueia ate todos os workers pararem com codigo de saida registrado.
func (sb *SwarmBoard) waitLevelStopped(ctx context.Context, names []string, poll time.Duration) ([]Worker, error) {
	want := make(map[string]bool, len(names))
	for _, n := range names {
		want[n] = true
	}
	ticker := time.NewTicker(poll)
	defer ticker.Stop()
	for {
		st, err := sb.Status()
		if err != nil {
			return nil, err
		}
		var done []Worker
		for _, w := range st.Workers {
			if want[w.Name] && w.Status == WorkerStopped && w.ExitCode != nil {
				done = append(done, w)
			}
		}
		if len(done) == len(names) {
			return done, nil
		}
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("espera dos workers: %w", ctx.Err())
		case <-ticker.C:
		}
	}
}

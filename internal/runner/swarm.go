package runner

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/board"
	"github.com/heverton-dev/LegoLoop/internal/dag"
)

// SwarmRunner interface para execução de níveis do DAG via Swarm.
type SwarmRunner interface {
	RunLevels(ctx context.Context, levels [][]board.RunNode, opt board.RunOptions) ([]board.RunLevelResult, error)
}

func isNativeSwarm(n dag.Node) bool {
	if n.Isolation == nil {
		return false
	}
	exec, ok := n.Isolation["executor"].(string)
	return ok && strings.TrimSpace(exec) == dag.ExecutorNativeSwarm
}

func (r *run) executarNivelSwarm(ctx context.Context, ids []string, mu *sync.Mutex, reports map[string]NodeReport) {
	if r.opts.Swarm == nil {
		mu.Lock()
		for _, id := range ids {
			reports[id] = NodeReport{ID: id, Status: StatusError, Error: "runner: swarm runner não configurado"}
		}
		mu.Unlock()
		return
	}
	res := r.runSwarmNivel(ctx, ids)
	mu.Lock()
	for _, nr := range res {
		reports[nr.ID] = nr
	}
	mu.Unlock()
}

func (r *run) runSwarmNivel(ctx context.Context, ids []string) []NodeReport {
	nodes, workerToID := r.montarRunNodes(ids)
	r.notificarInicioSwarm(ids)
	opt := board.RunOptions{RepoDir: r.opts.Workdir, DefaultTimeout: r.opts.DefaultTimeout}
	results, err := r.opts.Swarm.RunLevels(ctx, [][]board.RunNode{nodes}, opt)
	reports := r.mapearResultadosSwarm(ids, workerToID, results, err)
	r.notificarFimSwarm(reports)
	return reports
}

func (r *run) montarRunNodes(ids []string) ([]board.RunNode, map[string]string) {
	usados := make(map[string]bool, len(ids))
	nodes := make([]board.RunNode, 0, len(ids))
	workerToID := make(map[string]string, len(ids))
	for _, id := range ids {
		n := r.d.Node(id)
		if n == nil {
			continue
		}
		rn := nodeToRunNode(*n, usados)
		nodes = append(nodes, rn)
		workerToID[rn.Name] = id
	}
	return nodes, workerToID
}

func nodeToRunNode(n dag.Node, usados map[string]bool) board.RunNode {
	role := n.Role
	if role == "" {
		role = n.Papel
	}
	return board.RunNode{
		NodeID:  n.ID,
		Name:    board.WorkerNameFromNode(n.ID, usados),
		Role:    role,
		Model:   n.Routing.Model,
		Command: n.Command,
		Timeout: time.Duration(n.TimeoutSec) * time.Second,
		Gate:    n.Output.Gate,
	}
}

func (r *run) notificarInicioSwarm(ids []string) {
	if r.opts.OnNodeStart == nil {
		return
	}
	for _, id := range ids {
		r.opts.OnNodeStart(id)
	}
}

func (r *run) notificarFimSwarm(reports []NodeReport) {
	if r.opts.OnNodeFinish == nil {
		return
	}
	for _, nr := range reports {
		r.opts.OnNodeFinish(nr)
	}
}

func (r *run) mapearResultadosSwarm(ids []string, workerToID map[string]string, results []board.RunLevelResult, runErr error) []NodeReport {
	repMap := make(map[string]NodeReport, len(ids))
	for _, res := range results {
		for _, w := range res.Workers {
			id := workerToID[w.Name]
			if id == "" {
				id = w.Name
			}
			repMap[id] = workerToReport(id, w, runErr)
		}
	}
	out := make([]NodeReport, len(ids))
	for i, id := range ids {
		if nr, ok := repMap[id]; ok {
			out[i] = nr
		} else {
			out[i] = fallbackSwarmReport(id, runErr)
		}
	}
	return out
}

func workerToReport(id string, w board.RunWorkerResult, runErr error) NodeReport {
	nr := NodeReport{ID: id, ExitCode: w.ExitCode, Status: StatusPassed}
	if w.ExitCode != 0 {
		nr.Status = StatusError
		if runErr != nil {
			nr.Error = runErr.Error()
		} else {
			nr.Error = fmt.Sprintf("worker %s saiu com código %d", w.Name, w.ExitCode)
		}
	} else if w.GateOK != nil && !*w.GateOK {
		nr.Status = StatusGateFailed
		nr.Error = "gate reprovado"
	}
	return nr
}

func fallbackSwarmReport(id string, runErr error) NodeReport {
	if runErr != nil {
		return NodeReport{ID: id, Status: StatusError, Error: runErr.Error()}
	}
	return NodeReport{ID: id, Status: StatusPassed, ExitCode: 0}
}

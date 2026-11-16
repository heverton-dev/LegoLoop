package cli

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"legoloop/internal/board"
	"legoloop/internal/dag"
)

var reNomeInvalido = regexp.MustCompile(`[^a-zA-Z0-9_-]`)

// workerNameFromNode transforma o id do no em nome valido de worker e garante unicidade.
func workerNameFromNode(id string, usados map[string]bool) string {
	n := reNomeInvalido.ReplaceAllString(id, "-")
	if n == "" || !(n[0] >= 'a' && n[0] <= 'z' || n[0] >= 'A' && n[0] <= 'Z' || n[0] >= '0' && n[0] <= '9') {
		n = "n" + n
	}
	if len(n) > 60 {
		n = n[:60]
	}
	base, i := n, 2
	for usados[n] {
		n = fmt.Sprintf("%s-%d", base, i)
		i++
	}
	usados[n] = true
	return n
}

// swarmLevelsFromDAG monta os niveis a partir das camadas topologicas. No sem papel e stop_gate nao viram
// worker (o stop-gate e humano) e voltam em skipped.
func swarmLevelsFromDAG(p *dag.DAG, global []string) ([][]board.RunNode, []string, error) {
	usados := map[string]bool{}
	var levels [][]board.RunNode
	var skipped []string
	for _, layer := range p.Layers() {
		var lvl []board.RunNode
		for _, id := range layer {
			n := p.Node(id)
			role := n.Role
			if role == "" {
				role = n.Papel
			}
			if n.Type == dag.TypeStopGate || role == "" {
				skipped = append(skipped, id)
				continue
			}
			cmd := n.Command
			if len(cmd) == 0 {
				cmd = global
			}
			if len(cmd) == 0 {
				return nil, nil, fmt.Errorf("no %s sem comando: informe -- COMANDO ARGS... ou command no no", id)
			}
			lvl = append(lvl, board.RunNode{
				Name: workerNameFromNode(id, usados), Role: role, Model: n.Routing.Model,
				Command: cmd, Timeout: time.Duration(n.TimeoutSec) * time.Second,
			})
		}
		if len(lvl) > 0 {
			levels = append(levels, lvl)
		}
	}
	return levels, skipped, nil
}

// cmdSwarmRun executa o DAG por nivel topologico: todos os nos de um nivel sobem em paralelo, cada um em
// worktree proprio, e o nivel seguinte so comeca depois que todos pararam com codigo 0.
func cmdSwarmRun(ctx context.Context, args []string, d Deps) int {
	own, global := splitSpawnArgs(args)
	path, rest := splitPositional(own)
	fs := newFlags("swarm run", d)
	dir := fs.String("dir", defaultBoardDir(), "diretorio do quadro")
	repo := fs.String("repo", ".", "diretorio do repositorio git")
	ws := fs.String("workspace", "share", "share|branch|inherit")
	timeout := fs.Duration("timeout", 10*time.Minute, "tempo maximo por nivel")
	keep := fs.Bool("keep", false, "mantem workers e worktrees apos o sucesso")
	if err := fs.Parse(rest); err != nil {
		return 2
	}
	if path == "" {
		path = fs.Arg(0)
	}
	if path == "" {
		fmt.Fprintln(d.Stderr, "uso: legoloop swarm run <dag.json> [-dir D] [-repo R] [-workspace M] [-timeout T] [-keep] [-- comando args...]")
		return 2
	}
	p, err := loadDAG(path)
	if err != nil {
		return fail(d, err)
	}
	levels, skipped, err := swarmLevelsFromDAG(p, global)
	if err != nil {
		return fail(d, err)
	}
	sb, err := board.OpenSwarmBoard(*dir)
	if err != nil {
		return fail(d, err)
	}
	results, runErr := sb.RunLevels(ctx, levels, board.RunOptions{
		Mode: board.WorkspaceMode(*ws), RepoDir: *repo, Keep: *keep, DefaultTimeout: *timeout,
	})
	out := map[string]any{"dag": p.Name, "levels": results, "skipped": skipped}
	if code := printJSON(d, out); code != 0 {
		return code
	}
	if runErr != nil {
		return fail(d, runErr)
	}
	return 0
}

package cli

import (
	"context"
	"fmt"
	"time"

	"legoloop/internal/board"
)

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
	gates := fs.Bool("gates", false, "avalia o gate de cada no no worktree do worker; gate falso reprova o nivel")
	if err := fs.Parse(rest); err != nil {
		return 2
	}
	if path == "" {
		path = fs.Arg(0)
	}
	if path == "" {
		fmt.Fprintln(d.Stderr, "uso: legoloop swarm run <dag.json> [-dir D] [-repo R] [-workspace M] [-timeout T] [-keep] [-gates] [-- comando args...]")
		return 2
	}
	p, err := loadDAG(path)
	if err != nil {
		return fail(d, err)
	}
	levels, skipped, err := board.LevelsFromDAG(p, global)
	if err != nil {
		return fail(d, err)
	}
	sb, err := board.OpenSwarmBoard(*dir)
	if err != nil {
		return fail(d, err)
	}
	results, runErr := sb.RunLevels(ctx, levels, board.RunOptions{
		Mode: board.WorkspaceMode(*ws), RepoDir: *repo, Keep: *keep, DefaultTimeout: *timeout, CheckGates: *gates,
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

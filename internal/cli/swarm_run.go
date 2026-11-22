package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/board"
)

// swarmRunOpts consolida as opcoes da linha de comando do swarm run.
type swarmRunOpts struct {
	dir, repo, ws string
	timeout       time.Duration
	keep, gates   bool
	progresso     string
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
	gates := fs.Bool("gates", false, "avalia o gate de cada no no worktree do worker; gate falso reprova o nivel")
	prog := fs.String("progresso", "", "progresso.json do PLAN, atualizado a cada nó (padrão: campo progresso do DAG)")
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
	return execSwarmRun(ctx, d, path, global, swarmRunOpts{
		dir: *dir, repo: *repo, ws: *ws, timeout: *timeout, keep: *keep, gates: *gates, progresso: *prog,
	})
}

// execSwarmRun carrega o DAG, sobe os niveis e emite o JSON do resultado.
func execSwarmRun(ctx context.Context, d Deps, path string, global []string, o swarmRunOpts) int {
	p, err := loadDAG(path)
	if err != nil {
		return fail(d, err)
	}
	levels, skipped, err := board.LevelsFromDAG(p, global)
	if err != nil {
		return fail(d, err)
	}
	sb, err := board.OpenSwarmBoard(o.dir)
	if err != nil {
		return fail(d, err)
	}
	ro := board.RunOptions{
		Mode: board.WorkspaceMode(o.ws), RepoDir: o.repo, Keep: o.keep, DefaultTimeout: o.timeout, CheckGates: o.gates,
	}
	arq, err := caminhoProgressoVivo(o.repo, o.progresso, p.Progresso)
	if err != nil {
		return fail(d, err)
	}
	obs, err := ligarProgressoSwarm(ctx, &ro, arq, o.repo)
	if err != nil {
		return fail(d, err)
	}
	results, runErr := sb.RunLevels(ctx, levels, ro)
	out := map[string]any{"dag": p.Name, "levels": results, "skipped": skipped}
	if code := printJSON(d, out); code != 0 {
		return code
	}
	if runErr != nil {
		return codigoComProgresso(d, obs, fail(d, runErr))
	}
	return codigoComProgresso(d, obs, 0)
}

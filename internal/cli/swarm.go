package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/board"
)

// defaultBoardDir prioriza o quadro herdado do worker (LEGOLOOP_SWARM_BOARD) e cai no diretorio atual.
func defaultBoardDir() string {
	if env := strings.TrimSpace(os.Getenv("LEGOLOOP_SWARM_BOARD")); env != "" {
		return env
	}
	wd, err := os.Getwd()
	if err != nil {
		return ".legoloop-swarm"
	}
	return wd + "/.legoloop-swarm"
}

func printSwarmUsage(w io.Writer) {
	fmt.Fprintln(w, "uso: legoloop swarm <add|spawn|run|post|read|wait|status|kill|monitor|ui> [opcoes]")
	fmt.Fprintln(w, "subcomandos:")
	fmt.Fprintln(w, "  add      registra worker no enxame (-dir, -role, -name, -workspace, -repo)")
	fmt.Fprintln(w, "  spawn    inicia o processo real do worker (-dir, -name, -- comando args...)")
	fmt.Fprintln(w, "  run      executa o DAG por nivel topologico, workers paralelos (-dir, -repo, -workspace, -timeout, -keep, -- comando)")
	fmt.Fprintln(w, "  post     publica mensagem no quadro (-dir, -from, -to, -body)")
	fmt.Fprintln(w, "  read     le mensagens do quadro (-dir, -since)")
	fmt.Fprintln(w, "  wait     aguarda mensagem para o worker (-dir, -target, -since, -timeout)")
	fmt.Fprintln(w, "  status   exibe status dos workers (-dir)")
	fmt.Fprintln(w, "  kill     encerra worker do enxame (-dir, -name, -repo)")
	fmt.Fprintln(w, "  monitor  observa quadro do enxame (-dir)")
	fmt.Fprintln(w, "  ui       inicia painel web (-addr, -dir, -root)")
}

func cmdSwarm(ctx context.Context, args []string, d Deps) int {
	if len(args) == 0 {
		printSwarmUsage(d.Stderr)
		return 2
	}
	if args[0] == "--help" || args[0] == "-h" || args[0] == "help" {
		printSwarmUsage(d.Stdout)
		return 0
	}
	sub := map[string]func(context.Context, []string, Deps) int{
		"add": cmdSwarmAdd, "spawn": cmdSwarmSpawn, "run": cmdSwarmRun, "post": cmdSwarmPost, "read": cmdSwarmRead,
		"wait": cmdSwarmWait, "status": cmdSwarmStatus, "kill": cmdSwarmKill, "monitor": cmdSwarmMonitor,
		"ui": cmdSwarmUI,
	}
	fn, ok := sub[args[0]]
	if !ok {
		fmt.Fprintf(d.Stderr, "subcomando swarm desconhecido: %s\n", args[0])
		return 2
	}
	return fn(ctx, args[1:], d)
}

func cmdSwarmAdd(_ context.Context, args []string, d Deps) int {
	fs := newFlags("swarm add", d)
	role := fs.String("role", "", "papel canonico")
	name := fs.String("name", "", "nome do worker")
	parent := fs.String("parent", "", "processo pai (vazio = principal)")
	ws := fs.String("workspace", "inherit", "share|branch|inherit")
	model := fs.String("model", "", "modelo")
	dir := fs.String("dir", defaultBoardDir(), "diretorio do quadro")
	repo := fs.String("repo", ".", "diretorio do repositorio git")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	sb, err := board.OpenSwarmBoard(*dir)
	if err != nil {
		return fail(d, err)
	}
	w, err := sb.AddWorker(board.AddWorkerRequest{
		Name: *name, Role: *role, Parent: *parent, Model: *model,
		WorkspaceMode: board.WorkspaceMode(*ws), RepoDir: *repo,
	})
	if err != nil {
		return fail(d, err)
	}
	return printJSON(d, w)
}

func cmdSwarmPost(_ context.Context, args []string, d Deps) int {
	fs := newFlags("swarm post", d)
	from := fs.String("from", envOr("LEGOLOOP_SWARM_WORKER", "coordenador"), "remetente")
	to := fs.String("to", "@all", "destinatario")
	body := fs.String("body", "", "mensagem")
	dir := fs.String("dir", defaultBoardDir(), "diretorio do quadro")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	sb, err := board.OpenSwarmBoard(*dir)
	if err != nil {
		return fail(d, err)
	}
	msg, err := sb.PostMessage(*from, *to, *body)
	if err != nil {
		return fail(d, err)
	}
	return printJSON(d, msg)
}

func parseSeqs(raw string) []int64 {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var out []int64
	for _, p := range strings.Split(raw, ",") {
		if n, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64); err == nil {
			out = append(out, n)
		}
	}
	return out
}

func cmdSwarmRead(_ context.Context, args []string, d Deps) int {
	fs := newFlags("swarm read", d)
	since := fs.Int64("since", 0, "sequencia inicial")
	seqs := fs.String("seqs", "", "sequencias separadas por virgula")
	dir := fs.String("dir", defaultBoardDir(), "diretorio do quadro")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	sb, err := board.OpenSwarmBoard(*dir)
	if err != nil {
		return fail(d, err)
	}
	msgs, err := sb.ReadMessages(*since, parseSeqs(*seqs))
	if err != nil {
		return fail(d, err)
	}
	return printJSON(d, msgs)
}

func cmdSwarmWait(ctx context.Context, args []string, d Deps) int {
	fs := newFlags("swarm wait", d)
	target := fs.String("target", workerTargetFromEnv(), "alvo (@nome ou @all)")
	since := fs.Int64("since", 0, "sequencia inicial")
	timeout := fs.Duration("timeout", 30*time.Second, "tempo maximo de espera")
	dir := fs.String("dir", defaultBoardDir(), "diretorio do quadro")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	sb, err := board.OpenSwarmBoard(*dir)
	if err != nil {
		return fail(d, err)
	}
	wctx, cancel := context.WithTimeout(ctx, *timeout)
	defer cancel()
	msgs, err := sb.WaitMessages(wctx, *target, *since, 100*time.Millisecond)
	if err != nil {
		return fail(d, err)
	}
	return printJSON(d, msgs)
}

func cmdSwarmStatus(_ context.Context, args []string, d Deps) int {
	fs := newFlags("swarm status", d)
	dir := fs.String("dir", defaultBoardDir(), "diretorio do quadro")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	sb, err := board.OpenSwarmBoard(*dir)
	if err != nil {
		return fail(d, err)
	}
	st, err := sb.Status()
	if err != nil {
		return fail(d, err)
	}
	return printJSON(d, st)
}

func cmdSwarmKill(_ context.Context, args []string, d Deps) int {
	fs := newFlags("swarm kill", d)
	name := fs.String("name", "", "nome do worker")
	dir := fs.String("dir", defaultBoardDir(), "diretorio do quadro")
	repo := fs.String("repo", ".", "diretorio do repositorio git")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	sb, err := board.OpenSwarmBoard(*dir)
	if err != nil {
		return fail(d, err)
	}
	if err := sb.KillWorker(*name, *repo); err != nil {
		return fail(d, err)
	}
	return printJSON(d, map[string]string{"status": "killed", "name": *name})
}

func cmdSwarmMonitor(_ context.Context, args []string, d Deps) int {
	fs := newFlags("swarm monitor", d)
	timeout := fs.Duration("timeout", 60*time.Second, "limite de inatividade")
	dir := fs.String("dir", defaultBoardDir(), "diretorio do quadro")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	sb, err := board.OpenSwarmBoard(*dir)
	if err != nil {
		return fail(d, err)
	}
	stalled, err := sb.CheckStall(*timeout)
	if err != nil {
		return fail(d, err)
	}
	return printJSON(d, map[string]any{"stalled": stalled})
}

func envOr(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func workerTargetFromEnv() string {
	if w := strings.TrimSpace(os.Getenv("LEGOLOOP_SWARM_WORKER")); w != "" {
		return "@" + w
	}
	return ""
}

// cmdSwarmSpawn inicia o processo real do worker no seu worktree. O comando vem apos "--".
// O processo sobrevive ao CLI; o encerramento posterior usa o PID gravado no quadro (swarm kill).
func cmdSwarmSpawn(ctx context.Context, args []string, d Deps) int {
	own, cmdArgs := splitSpawnArgs(args)
	fs := newFlags("swarm spawn", d)
	name := fs.String("name", "", "nome do worker ja registrado")
	dir := fs.String("dir", defaultBoardDir(), "diretorio do quadro")
	if err := fs.Parse(own); err != nil {
		return 2
	}
	if strings.TrimSpace(*name) == "" || len(cmdArgs) == 0 {
		fmt.Fprintln(d.Stderr, "uso: legoloop swarm spawn -name W [-dir D] -- comando [args...]")
		return 2
	}
	sb, err := board.OpenSwarmBoard(*dir)
	if err != nil {
		return fail(d, err)
	}
	pid, err := sb.SpawnWorker(context.WithoutCancel(ctx), *name, cmdArgs)
	if err != nil {
		return fail(d, err)
	}
	return printJSON(d, map[string]any{"name": *name, "pid": pid, "status": board.WorkerWorking})
}

func splitSpawnArgs(args []string) ([]string, []string) {
	for i, a := range args {
		if a == "--" {
			return args[:i], args[i+1:]
		}
	}
	return args, nil
}

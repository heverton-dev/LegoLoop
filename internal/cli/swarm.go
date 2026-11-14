package cli

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"legoloop/internal/board"
)

func defaultBoardDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return ".legoloop-swarm"
	}
	return wd + "/.legoloop-swarm"
}

func cmdSwarm(ctx context.Context, args []string, d Deps) int {
	if len(args) == 0 {
		fmt.Fprintln(d.Stderr, "uso: legoloop swarm <add|post|read|wait|status|kill|monitor> [opcoes]")
		return 2
	}
	sub := map[string]func(context.Context, []string, Deps) int{
		"add": cmdSwarmAdd, "post": cmdSwarmPost, "read": cmdSwarmRead,
		"wait": cmdSwarmWait, "status": cmdSwarmStatus, "kill": cmdSwarmKill, "monitor": cmdSwarmMonitor,
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
		Name: *name, Role: *role, Model: *model,
		WorkspaceMode: board.WorkspaceMode(*ws), RepoDir: *repo,
	})
	if err != nil {
		return fail(d, err)
	}
	return printJSON(d, w)
}

func cmdSwarmPost(_ context.Context, args []string, d Deps) int {
	fs := newFlags("swarm post", d)
	from := fs.String("from", "coordenador", "remetente")
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
	target := fs.String("target", "", "alvo (@nome ou @all)")
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

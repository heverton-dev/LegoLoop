package cli

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"

	"legoloop/internal/graph"
	"legoloop/internal/studio"
)

func resolveSwarmProjectDir(root, dir string) string {
	target := root
	if target == "" {
		target = dir
	}
	abs, err := filepath.Abs(target)
	if err != nil {
		return target
	}
	base := filepath.Base(abs)
	if base == ".legoloop-swarm" || base == ".swarm" {
		return filepath.Dir(abs)
	}
	return abs
}

func cmdSwarmUI(ctx context.Context, args []string, d Deps) int {
	fs := newFlags("swarm ui", d)
	addr := fs.String("addr", "127.0.0.1:7879", "endereco HTTP do painel")
	dir := fs.String("dir", defaultBoardDir(), "diretorio do quadro")
	root := fs.String("root", "", "diretorio raiz com multiplos quadros")
	_ = fs.String("process", "swarm", "nome do processo em execucao")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	idx, err := graph.NewIndex(graph.DefaultCatalog())
	if err != nil {
		return fail(d, err)
	}
	srv := &studio.Server{
		Index:      idx,
		ProjectDir: resolveSwarmProjectDir(*root, *dir),
		Discover:   d.discover,
		Home:       d.Home,
	}
	ready := make(chan string, 1)
	go func() {
		if a, ok := <-ready; ok {
			fmt.Fprintf(d.Stdout, "Painel Swarm em http://%s/swarm\n", a)
			if d.Ready != nil {
				d.Ready <- a
			}
		}
	}()
	err = srv.ListenAndServe(ctx, *addr, ready)
	close(ready)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fail(d, err)
	}
	return 0
}

package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"legoloop/internal/mcp"
)

func cmdMcpRouter(_ context.Context, args []string, d Deps) int {
	if d.Stdin == nil || d.Stdout == nil {
		fmt.Fprintln(d.Stderr, "mcp-router requer stdin e stdout conectados")
		return 1
	}
	fs := newFlags("mcp-router", d)
	vaultFlag := fs.String("vault", "", "cofre canônico (padrão ~/.aidd/canonical)")
	projectFlag := fs.String("project", ".", "diretório do projeto")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	home, vault := resolveMcpRouterPaths(d.Home, *vaultFlag)
	router := mcp.NewConfiguredRouter(home, vault, *projectFlag)
	defer router.Close()
	if err := router.CarregarToolsDeProcessos(*projectFlag); err != nil {
		fmt.Fprintf(d.Stderr, "mcp-router: %v\n", err)
		return 1
	}
	if err := router.Serve(d.Stdin, d.Stdout); err != nil {
		fmt.Fprintf(d.Stderr, "mcp-router encerrou com erro: %v\n", err)
		return 1
	}
	return 0
}

func resolveMcpRouterPaths(home, vault string) (string, string) {
	if home == "" {
		home, _ = os.UserHomeDir()
	}
	if vault == "" && home != "" {
		vault = filepath.Join(home, ".aidd", "canonical")
	}
	return home, vault
}

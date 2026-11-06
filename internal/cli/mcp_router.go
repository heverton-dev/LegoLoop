package cli

import (
	"context"
	"fmt"
	"legoloop/internal/mcp"
)

func cmdMcpRouter(_ context.Context, _ []string, d Deps) int {
	if d.Stdin == nil || d.Stdout == nil {
		fmt.Fprintln(d.Stderr, "mcp-router requer stdin e stdout conectados")
		return 1
	}
	router := mcp.NewRouter(nil, nil)
	if err := router.Serve(d.Stdin, d.Stdout); err != nil {
		fmt.Fprintf(d.Stderr, "mcp-router encerrou com erro: %v\n", err)
		return 1
	}
	return 0
}

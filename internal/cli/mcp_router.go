package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/heverton-dev/LegoLoop/internal/a2ui"
	"github.com/heverton-dev/LegoLoop/internal/mcp"
	"github.com/heverton-dev/LegoLoop/internal/studio"
)

func cmdMcpRouter(_ context.Context, args []string, d Deps) int {
	if d.Stdin == nil || d.Stdout == nil {
		fmt.Fprintln(d.Stderr, "mcp-router requer stdin e stdout conectados")
		return 1
	}
	fs := newFlags("mcp-router", d)
	vaultFlag := fs.String("vault", "", "cofre canônico (padrão ~/.aidd/canonical)")
	projectFlag := fs.String("project", ".", "diretório do projeto")
	gatewayFlag := fs.Bool("gateway", false, "expõe só as 4 ferramentas de gateway (piso de contexto); ws e swarm_* seguem por call_mcp_tool")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	home, vault := resolveMcpRouterPaths(d.Home, *vaultFlag)
	router := mcp.NewConfiguredRouter(home, vault, *projectFlag)
	defer router.Close()
	defer attachEmbeddedA2UI(router, *projectFlag, d.Stderr)()
	if *gatewayFlag {
		router.EnableGateway()
	}
	if err := router.CarregarToolsDeProcessos(*projectFlag); err != nil {
		fmt.Fprintf(d.Stderr, "mcp-router: %v\n", err)
		return mcp.CodigoDeSaida(err)
	}
	if err := router.Serve(d.Stdin, d.Stdout); err != nil {
		fmt.Fprintf(d.Stderr, "mcp-router encerrou com erro: %v\n", err)
		return 1
	}
	return 0
}

func attachEmbeddedA2UI(router *mcp.Router, projectDir string, stderr io.Writer) func() {
	hub := a2ui.NewHub(a2ui.DefaultCatalog())
	baseURL, stopUI, err := studio.StartEmbedded(projectDir, hub)
	if err != nil {
		return func() {}
	}
	router.SetSessaoHub(hub, baseURL)
	router.SetSessaoAbridor(abrirURLSistema)
	return func() {
		if cErr := stopUI(); cErr != nil && stderr != nil {
			fmt.Fprintf(stderr, "mcp-router ui close: %v\n", cErr)
		}
	}
}

func abrirURLSistema(url string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return exec.Command("xdg-open", url).Start()
	}
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

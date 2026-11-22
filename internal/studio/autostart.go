package studio

import (
	"context"
	"net"
	"net/http"

	"github.com/heverton-dev/LegoLoop/internal/a2ui"
	"github.com/heverton-dev/LegoLoop/internal/provider"
)

// StartEmbedded sobe um servidor local do Studio em loopback (porta livre) conectado ao hub A2UI,
// dispensando iniciar `legoloop studio` manualmente nas vias CLI e MCP.
func StartEmbedded(projectDir string, hub *a2ui.Hub) (string, func() error, error) {
	if hub == nil {
		hub = a2ui.NewHub(a2ui.DefaultCatalog())
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", nil, err
	}
	host := ln.Addr().String()
	srv := &Server{
		ProjectDir: projectDir,
		Host:       host,
		A2UI:       hub,
		Discover: func(ctx context.Context) provider.Inventory {
			return provider.Discover(ctx, provider.DefaultConfig())
		},
	}
	httpSrv := &http.Server{Handler: srv.Handler()}
	go func() {
		if sErr := httpSrv.Serve(ln); sErr != nil && sErr != http.ErrServerClosed {
			return
		}
	}()
	return "http://" + host, httpSrv.Close, nil
}

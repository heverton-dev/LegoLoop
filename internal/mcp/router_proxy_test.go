package mcp

import (
	"strings"
	"testing"
	"time"
)

func routerComMcpExterno(t *testing.T) *Router {
	t.Helper()
	home := t.TempDir()
	specs := map[string]ServerSpec{"fake": fakeSpec()}
	if err := SaveRegistry(RegistryPath(home), specs); err != nil {
		t.Fatal(err)
	}
	cat := map[string][]CatalogEntry{"fake": {{Name: "echo", Description: "ecoa o texto recebido"}}}
	if err := SaveCatalog(CatalogPath(home), cat); err != nil {
		t.Fatal(err)
	}
	r := NewConfiguredRouter(home, t.TempDir(), t.TempDir())
	t.Cleanup(r.Close)
	return r
}

func TestRouterBuscaMcpExternoNoCatalogoSemSubirProcesso(t *testing.T) {
	r := routerComMcpExterno(t)
	resp := r.callMcpSearch(1, map[string]any{"query": "ecoa"})
	txt := respText(t, resp)
	if !strings.Contains(txt, "[fake]") || !strings.Contains(txt, "echo") {
		t.Fatalf("busca nao achou o MCP externo: %s", txt)
	}
	if r.proxy.ativos() != 0 {
		t.Fatal("buscar nao pode subir servidor")
	}
}

func TestRouterBuscaPeloNomeDoServidorNaoIndexado(t *testing.T) {
	home := t.TempDir()
	_ = SaveRegistry(RegistryPath(home), map[string]ServerSpec{"novo": fakeSpec()})
	r := NewConfiguredRouter(home, t.TempDir(), t.TempDir())
	t.Cleanup(r.Close)
	txt := respText(t, r.callMcpSearch(1, map[string]any{"query": "novo"}))
	if !strings.Contains(txt, "[novo]") || !strings.Contains(txt, "mcp-registry refresh") {
		t.Fatalf("deve orientar a indexar: %s", txt)
	}
}

func TestRouterExecutaMcpExternoSobDemanda(t *testing.T) {
	r := routerComMcpExterno(t)
	args := map[string]any{"server": "fake", "tool": "echo", "arguments": map[string]any{"text": "oi"}}
	txt := respText(t, r.callMcpTool(1, args))
	if txt != "eco:oi" {
		t.Fatalf("esperava eco:oi, veio %q", txt)
	}
	r.Close()
	deadline := time.Now().Add(3 * time.Second)
	for r.proxy.ativos() != 0 && time.Now().Before(deadline) {
		time.Sleep(50 * time.Millisecond)
	}
	if r.proxy.ativos() != 0 {
		t.Fatal("Close deve encerrar os servidores")
	}
}

func respText(t *testing.T, resp JSONRPCResponse) string {
	t.Helper()
	m, ok := resp.Result.(map[string]any)
	if !ok {
		t.Fatalf("resultado inesperado: %v", resp)
	}
	items, _ := m["content"].([]map[string]any)
	if len(items) == 0 {
		t.Fatalf("sem conteudo: %v", m)
	}
	s, _ := items[0]["text"].(string)
	return s
}

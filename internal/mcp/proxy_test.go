package mcp

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestMain transforma o binario de teste em um servidor MCP falso quando FAKE_MCP=1.
func TestMain(m *testing.M) {
	if os.Getenv("FAKE_MCP") == "1" {
		runFakeMCP()
		return
	}
	os.Exit(m.Run())
}

func runFakeMCP() {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		var req struct {
			ID     any             `json:"id"`
			Method string          `json:"method"`
			Params json.RawMessage `json:"params"`
		}
		if json.Unmarshal(sc.Bytes(), &req) != nil || req.ID == nil {
			continue
		}
		out, _ := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": req.ID, "result": fakeResult(req.Method, req.Params)})
		os.Stdout.Write(append(out, '\n'))
	}
}

func fakeResult(method string, params json.RawMessage) any {
	switch method {
	case "tools/list":
		return map[string]any{"tools": []map[string]any{{"name": "echo", "description": "ecoa o texto recebido"}}}
	case "tools/call":
		var p struct {
			Arguments map[string]any `json:"arguments"`
		}
		_ = json.Unmarshal(params, &p)
		txt, _ := p.Arguments["text"].(string)
		return map[string]any{"content": []map[string]any{{"type": "text", "text": "eco:" + txt}}}
	}
	return map[string]any{"protocolVersion": "2024-11-05", "capabilities": map[string]any{}}
}

func fakeSpec() ServerSpec {
	return ServerSpec{Command: os.Args[0], Env: map[string]string{"FAKE_MCP": "1"}}
}

func TestRegistryRoundTripEMerge(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "mcp-registry.json")
	got, err := LoadRegistry(path)
	if err != nil || len(got) != 0 {
		t.Fatalf("registro ausente deve ser vazio: %v %v", got, err)
	}
	if err := SaveRegistry(path, map[string]ServerSpec{"a": {Command: "x"}}); err != nil {
		t.Fatal(err)
	}
	merged := MergeRegistry(map[string]ServerSpec{"a": {Command: "x"}}, map[string]ServerSpec{"a": {Command: "outro"}, "b": {Command: "y"}})
	if merged["a"].Command != "x" || merged["b"].Command != "y" {
		t.Fatalf("merge nao pode sobrescrever o existente: %v", merged)
	}
	back, _ := LoadRegistry(path)
	if back["a"].Command != "x" {
		t.Fatalf("round trip falhou: %v", back)
	}
}

func TestProxySobeSoNaPrimeiraChamadaEEcoa(t *testing.T) {
	p := NewProxy(map[string]ServerSpec{"fake": fakeSpec()}, time.Minute)
	defer p.Close()
	if p.ativos() != 0 {
		t.Fatal("nenhum processo deve existir antes da primeira chamada")
	}
	out, err := p.CallTool("fake", "echo", map[string]any{"text": "oi"})
	if err != nil || out != "eco:oi" {
		t.Fatalf("chamada falhou: %q %v", out, err)
	}
	if p.ativos() != 1 {
		t.Fatalf("esperava 1 servidor ativo, veio %d", p.ativos())
	}
}

func TestProxyEncerraPorOciosidade(t *testing.T) {
	p := NewProxy(map[string]ServerSpec{"fake": fakeSpec()}, 150*time.Millisecond)
	defer p.Close()
	if _, err := p.CallTool("fake", "echo", map[string]any{"text": "a"}); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(3 * time.Second)
	for p.ativos() != 0 && time.Now().Before(deadline) {
		time.Sleep(50 * time.Millisecond)
	}
	if p.ativos() != 0 {
		t.Fatal("servidor ocioso deveria ter sido encerrado")
	}
	if out, err := p.CallTool("fake", "echo", map[string]any{"text": "b"}); err != nil || out != "eco:b" {
		t.Fatalf("deve religar sob demanda: %q %v", out, err)
	}
}

func TestProxyServidorDesconhecidoEErro(t *testing.T) {
	p := NewProxy(map[string]ServerSpec{}, time.Minute)
	if _, err := p.CallTool("nao-existe", "x", nil); err == nil || !strings.Contains(err.Error(), "nao-existe") {
		t.Fatalf("erro deve citar o servidor: %v", err)
	}
}

func TestCatalogoRefreshEBusca(t *testing.T) {
	home := t.TempDir()
	p := NewProxy(map[string]ServerSpec{"fake": fakeSpec()}, time.Minute)
	defer p.Close()
	cat, errs := p.Refresh()
	if len(errs) != 0 || len(cat["fake"]) != 1 || cat["fake"][0].Name != "echo" {
		t.Fatalf("refresh falhou: %v %v", cat, errs)
	}
	if p.ativos() != 0 {
		t.Fatal("refresh deve encerrar os servidores que subiu")
	}
	if err := SaveCatalog(CatalogPath(home), cat); err != nil {
		t.Fatal(err)
	}
	loaded, _ := LoadCatalog(CatalogPath(home))
	if loaded["fake"][0].Description != "ecoa o texto recebido" {
		t.Fatalf("catalogo nao persistiu: %v", loaded)
	}
}

func TestStartStdioUsaCwdDoSpec(t *testing.T) {
	spec := fakeSpec()
	spec.Cwd = filepath.Join(t.TempDir(), "nao-existe")
	if _, err := startStdio(spec); err == nil {
		t.Fatal("cwd inexistente deve falhar, provando que o Dir foi aplicado")
	}
}

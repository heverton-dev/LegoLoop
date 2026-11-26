package mcp

import (
	"os"
	"path/filepath"
	"testing"
)

// O router fica vivo por horas; editar o registro tem de valer na próxima subida do servidor.
func TestProxySpecReleRegistroAntesDeSubir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mcp-registry.json")
	escrever := func(conteudo string) {
		if err := os.WriteFile(path, []byte(conteudo), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	escrever(`{"servers":{"flow":{"command":"node"}}}`)
	specs, _ := LoadRegistry(path)
	p := NewProxy(specs, 0)
	p.recarregar = func() (map[string]ServerSpec, error) { return LoadRegistry(path) }
	escrever(`{"servers":{"flow":{"command":"node","env":{"FLOW_MCP_CDP_PORT":"9222"}}}}`)
	spec, ok := p.spec("flow")
	if !ok || spec.Env["FLOW_MCP_CDP_PORT"] != "9222" {
		t.Fatalf("spec deveria refletir o registro atual, veio %+v", spec)
	}
}

func TestProxySpecMantemCacheSeRegistroQuebrar(t *testing.T) {
	p := NewProxy(map[string]ServerSpec{"flow": {Command: "node"}}, 0)
	p.recarregar = func() (map[string]ServerSpec, error) { return nil, os.ErrNotExist }
	if spec, ok := p.spec("flow"); !ok || spec.Command != "node" {
		t.Fatalf("falha de leitura deveria manter o cache, veio %+v %v", spec, ok)
	}
	if _, ok := p.spec("outro"); ok {
		t.Fatal("servidor ausente nao pode existir")
	}
}

func TestProxySpecSemRecarregarUsaCache(t *testing.T) {
	p := NewProxy(map[string]ServerSpec{"flow": {Command: "node"}}, 0)
	if spec, ok := p.spec("flow"); !ok || spec.Command != "node" {
		t.Fatalf("sem recarregar deveria usar o cache, veio %+v", spec)
	}
}

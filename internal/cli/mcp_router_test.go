package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/caminhos"
)

// TestMcpRouterExpoeToolsDeProcessos prova a AC-E09.1 pelo CLI real: o comando
// legoloop mcp-router expõe a tool do Processo em tools/list do STDIO único.
func TestMcpRouterExpoeToolsDeProcessos(t *testing.T) {
	base := t.TempDir()
	manifesto := map[string]any{"tools": []any{map[string]any{
		"name":        "status_alpha",
		"description": "Devolve o status do Processo alpha.",
		"comando":     []string{"legoloop", "listar", "--json"},
	}}}
	data, err := json.Marshal(manifesto)
	if err != nil {
		t.Fatalf("marshal manifesto: %v", err)
	}
	dir := filepath.Join(caminhos.DirProcessosEm(base), "alpha")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("montar Processo: %v", err)
	}
	if err := os.WriteFile(caminhos.CaminhoToolsDoProcesso(dir), data, 0o600); err != nil {
		t.Fatalf("gravar manifesto: %v", err)
	}

	reqs := strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05"}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`,
	}, "\n") + "\n"
	var stdout, stderr bytes.Buffer
	deps := Deps{Stdin: strings.NewReader(reqs), Stdout: &stdout, Stderr: &stderr, Home: t.TempDir()}
	if code := Run(context.Background(), []string{"mcp-router", "--project", base}, deps); code != 0 {
		t.Fatalf("mcp-router exit %d: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "status_alpha") {
		t.Errorf("tools/list do mcp-router deve expor status_alpha, obteve: %s", stdout.String())
	}
	if !strings.Contains(stdout.String(), "mcp_search") {
		t.Errorf("bootstrap do router ausente em tools/list: %s", stdout.String())
	}
}

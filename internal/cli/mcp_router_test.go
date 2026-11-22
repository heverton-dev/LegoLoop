package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/caminhos"
)

// TestMcpRouterConflitoSaidaNaoZero prova a AC-E09.2 no CLI real: conflito de
// nome de tool entre Processos recusa o registro e o mcp-router sai com exit ≠ 0.
func TestMcpRouterConflitoSaidaNaoZero(t *testing.T) {
	base := t.TempDir()
	for _, nome := range []string{"alpha", "beta"} {
		manifesto := map[string]any{"tools": []any{map[string]any{
			"name":        "eco_dup",
			"description": "Tool em conflito entre Processos.",
			"comando":     []string{"legoloop", "listar", "--json"},
		}}}
		data, err := json.Marshal(manifesto)
		if err != nil {
			t.Fatalf("marshal manifesto: %v", err)
		}
		dir := filepath.Join(caminhos.DirProcessosEm(base), nome)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("montar Processo: %v", err)
		}
		if err := os.WriteFile(caminhos.CaminhoToolsDoProcesso(dir), data, 0o600); err != nil {
			t.Fatalf("gravar manifesto: %v", err)
		}
	}
	var stdout, stderr bytes.Buffer
	deps := Deps{Stdin: strings.NewReader(""), Stdout: &stdout, Stderr: &stderr, Home: t.TempDir()}
	code := Run(context.Background(), []string{"mcp-router", "--project", base}, deps)
	if code == 0 {
		t.Fatalf("esperava exit ≠ 0 no conflito, obteve 0: %s", stderr.String())
	}
	for _, trecho := range []string{"eco_dup", "alpha", "beta"} {
		if !strings.Contains(stderr.String(), trecho) {
			t.Errorf("conflito devia nomear %q no stderr, obteve: %s", trecho, stderr.String())
		}
	}
}

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

// TestMcpRouterGatewayExpoeSoQuatroFerramentas prova a flag --gateway no CLI real:
// tools/list devolve apenas call_mcp_tool; sem a flag, o swarm_* aparece.
func TestMcpRouterGatewayExpoeUmaFerramenta(t *testing.T) {
	contar := func(args ...string) int {
		var stdout, stderr bytes.Buffer
		entrada := `{"jsonrpc":"2.0","id":1,"method":"tools/list"}` + "\n"
		deps := Deps{Stdin: strings.NewReader(entrada), Stdout: &stdout, Stderr: &stderr, Home: t.TempDir()}
		if code := Run(context.Background(), append([]string{"mcp-router", "--project", t.TempDir()}, args...), deps); code != 0 {
			t.Fatalf("mcp-router %v saiu com %d: %s", args, code, stderr.String())
		}
		var resp struct {
			Result struct {
				Tools []struct{ Name string } `json:"tools"`
			} `json:"result"`
		}
		if err := json.Unmarshal([]byte(strings.TrimSpace(stdout.String())), &resp); err != nil {
			t.Fatalf("resposta inválida %q: %v", stdout.String(), err)
		}
		return len(resp.Result.Tools)
	}
	if n := contar("--gateway"); n != 1 {
		t.Errorf("--gateway expõe %d ferramentas, quer 1", n)
	}
	if n := contar(); n <= 1 {
		t.Errorf("sem --gateway o router deve expor mais de 1 ferramenta, expôs %d", n)
	}
}

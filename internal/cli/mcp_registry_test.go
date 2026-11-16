package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"legoloop/internal/mcp"
)

func rodarRegistry(home string, args ...string) (int, string, string) {
	out, errb := &bytes.Buffer{}, &bytes.Buffer{}
	deps := Deps{Stdin: strings.NewReader(""), Stdout: out, Stderr: errb, Home: home}
	code := Run(context.Background(), append([]string{"mcp-registry"}, args...), deps)
	return code, out.String(), errb.String()
}

func TestMcpRegistryListNuncaImprimeValoresDeEnv(t *testing.T) {
	home := t.TempDir()
	specs := map[string]mcp.ServerSpec{"supa": {Command: "npx", Args: []string{"-y", "x"}, Env: map[string]string{"TOKEN": "segredo-123"}}}
	if err := mcp.SaveRegistry(mcp.RegistryPath(home), specs); err != nil {
		t.Fatal(err)
	}
	code, out, _ := rodarRegistry(home, "list")
	if code != 0 || !strings.Contains(out, "supa") || !strings.Contains(out, "npx") {
		t.Fatalf("list falhou (%d): %s", code, out)
	}
	if strings.Contains(out, "segredo-123") {
		t.Fatalf("list vazou valor de env: %s", out)
	}
}

func TestMcpRegistryRefreshReportaFalhaPorServidor(t *testing.T) {
	home := t.TempDir()
	specs := map[string]mcp.ServerSpec{"quebrado": {Command: "comando-que-nao-existe-xyz"}}
	if err := mcp.SaveRegistry(mcp.RegistryPath(home), specs); err != nil {
		t.Fatal(err)
	}
	code, out, _ := rodarRegistry(home, "refresh")
	if code != 1 || !strings.Contains(out, "quebrado") || !strings.Contains(out, "falha") {
		t.Fatalf("esperava exit 1 citando o servidor (%d): %s", code, out)
	}
}

func TestMcpRegistrySemSubcomandoDaUso(t *testing.T) {
	if code, _, errb := rodarRegistry(t.TempDir()); code != 2 || !strings.Contains(errb, "Uso:") {
		t.Fatalf("esperava exit 2 com uso: %d %s", code, errb)
	}
}

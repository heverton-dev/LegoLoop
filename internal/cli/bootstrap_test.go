package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBootstrapDryRunJSON(t *testing.T) {
	home := t.TempDir()
	vault := filepath.Join(t.TempDir(), "canonical")
	var out, errOut bytes.Buffer
	d := Deps{Stdout: &out, Stderr: &errOut, Home: home}

	code := Run(context.Background(), []string{"bootstrap", "--vault", vault, "--dry-run", "--json"}, d)
	if code != 0 {
		t.Fatalf("esperado exit 0, veio %d; stderr=%s", code, errOut.String())
	}

	var res bootstrapResult
	if err := json.Unmarshal(out.Bytes(), &res); err != nil {
		t.Fatalf("JSON inválido: %v\nsaída=%s", err, out.String())
	}
	if !res.DryRun || res.Status != "ok" {
		t.Fatalf("resultado inesperado: %+v", res)
	}
	if len(res.Skills) == 0 {
		t.Fatalf("esperava skills listadas no dry-run")
	}
	if _, err := os.Stat(vault); !os.IsNotExist(err) {
		t.Fatalf("dry-run não deveria criar vault %s", vault)
	}
}

func TestBootstrapApplyExecutaOrquestracao(t *testing.T) {
	home := t.TempDir()
	vault := filepath.Join(t.TempDir(), "canonical")
	claudeCfg := filepath.Join(home, ".claude.json")
	if err := os.WriteFile(claudeCfg, []byte(`{"mcpServers":{}}`), 0644); err != nil {
		t.Fatal(err)
	}

	var out, errOut bytes.Buffer
	d := Deps{Stdout: &out, Stderr: &errOut, Home: home}

	code := Run(context.Background(), []string{"bootstrap", "--vault", vault, "--json"}, d)
	if code != 0 {
		t.Fatalf("esperado exit 0, veio %d; stderr=%s", code, errOut.String())
	}

	var res bootstrapResult
	if err := json.Unmarshal(out.Bytes(), &res); err != nil {
		t.Fatalf("JSON inválido: %v", err)
	}
	if res.Status != "ok" || !res.HarnessesSynced || !res.EnvReady {
		t.Fatalf("esperava bootstrap completo, veio: %+v", res)
	}

	data, err := os.ReadFile(claudeCfg)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "mcp-router") {
		t.Fatalf("esperava mcp-router injetado em %s, veio: %s", claudeCfg, string(data))
	}
}

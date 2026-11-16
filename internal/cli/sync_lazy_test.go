package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/mcp"
)

func rodarSync(t *testing.T, home string, args ...string) string {
	t.Helper()
	out, errb := &bytes.Buffer{}, &bytes.Buffer{}
	deps := Deps{Stdin: strings.NewReader(""), Stdout: out, Stderr: errb, Home: home}
	if code := Run(context.Background(), append([]string{"sync-harnesses"}, args...), deps); code != 0 {
		t.Fatalf("sync saiu %d: %s", code, errb.String())
	}
	return out.String()
}

func escrever(t *testing.T, path, conteudo string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
}

func lerJSON(t *testing.T, path string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var obj map[string]any
	if err := json.Unmarshal(stripJSONC(raw), &obj); err != nil {
		t.Fatalf("%s invalido: %v", path, err)
	}
	return obj
}

func TestSyncMoveMcpStdioParaRegistroEMantemRemoto(t *testing.T) {
	home := t.TempDir()
	cfg := filepath.Join(home, ".gemini", "config", "mcp_config.json")
	escrever(t, cfg, `{"mcpServers":{
	  "context7":{"command":"npx","args":["-y","@upstash/context7-mcp"]},
	  "github":{"command":"docker","args":["run","-i"],"env":{"TOKEN":"abc"}},
	  "cloudflare":{"serverUrl":"https://mcp.cloudflare.com/mcp"}}}`)
	rodarSync(t, home, "--apply")

	servers := lerJSON(t, cfg)["mcpServers"].(map[string]any)
	if servers["legoloop"] == nil || servers["cloudflare"] == nil {
		t.Fatalf("deve manter legoloop e o remoto: %v", servers)
	}
	if servers["context7"] != nil || servers["github"] != nil {
		t.Fatalf("stdio eager deve sair do harness: %v", servers)
	}
	reg, _ := mcp.LoadRegistry(mcp.RegistryPath(home))
	if reg["context7"].Command != "npx" || reg["github"].Env["TOKEN"] != "abc" {
		t.Fatalf("registro nao recebeu os servidores: %v", reg)
	}
	if reg["cloudflare"].Command != "" {
		t.Fatalf("remoto nao deve ir ao registro stdio: %v", reg["cloudflare"])
	}
}

func TestSyncOpencodeLocalEntradaSoEnabledIntocada(t *testing.T) {
	home := t.TempDir()
	cfg := filepath.Join(home, ".config", "opencode", "opencode.json")
	escrever(t, cfg, `{"mcp":{"rea":{"type":"local","command":["cmd.exe","/c","rea","mcp"],"environment":{"K":"v"},"enabled":true},"so":{"enabled":true}}}`)
	rodarSync(t, home, "--apply")

	mcpMap := lerJSON(t, cfg)["mcp"].(map[string]any)
	if mcpMap["rea"] != nil || mcpMap["legoloop"] == nil || mcpMap["so"] == nil {
		t.Fatalf("estado inesperado do opencode: %v", mcpMap)
	}
	reg, _ := mcp.LoadRegistry(mcp.RegistryPath(home))
	got := reg["rea"]
	if got.Command != "cmd.exe" || len(got.Args) != 3 || got.Env["K"] != "v" {
		t.Fatalf("conversao do formato opencode falhou: %+v", got)
	}
}

func TestSyncNaoSobrescreveRegistroNemBackupOriginal(t *testing.T) {
	home := t.TempDir()
	if err := mcp.SaveRegistry(mcp.RegistryPath(home), map[string]mcp.ServerSpec{"a": {Command: "original"}}); err != nil {
		t.Fatal(err)
	}
	cfg := filepath.Join(home, ".cursor", "mcp.json")
	escrever(t, cfg, `{"mcpServers":{"a":{"command":"outro"},"b":{"command":"novo"}}}`)
	rodarSync(t, home, "--apply")
	rodarSync(t, home, "--apply")

	reg, _ := mcp.LoadRegistry(mcp.RegistryPath(home))
	if reg["a"].Command != "original" || reg["b"].Command != "novo" {
		t.Fatalf("registro deve ser add-only: %v", reg)
	}
	bak, _ := os.ReadFile(cfg + ".bak_legoloop")
	if !strings.Contains(string(bak), `"outro"`) {
		t.Fatalf("segundo apply nao pode sobrescrever o backup original: %s", bak)
	}
}

func TestSyncIgnoraServidorDesabilitadoECobreChaveServers(t *testing.T) {
	home := t.TempDir()
	ide := filepath.Join(home, ".gemini", "antigravity-ide", "mcp_config.json")
	escrever(t, ide, `{"mcpServers":{"off":{"command":"x","disabled":true}}}`)
	vs := filepath.Join(home, "AppData", "Roaming", "Code", "User", "mcp.json")
	escrever(t, vs, `{"servers":{"rea":{"type":"stdio","command":"rea","args":["mcp"]}}}`)
	rodarSync(t, home, "--apply")

	if lerJSON(t, ide)["mcpServers"].(map[string]any)["off"] == nil {
		t.Fatal("servidor desabilitado deve ficar como esta")
	}
	reg, _ := mcp.LoadRegistry(mcp.RegistryPath(home))
	if _, ok := reg["off"]; ok {
		t.Fatal("desabilitado nao entra no registro")
	}
	if reg["rea"].Command != "rea" || lerJSON(t, vs)["servers"].(map[string]any)["rea"] != nil {
		t.Fatalf("chave 'servers' do VS Code nao foi tratada: %v", reg)
	}
}

func TestSyncVSCodeUsaChaveServersComTypeStdio(t *testing.T) {
	home := t.TempDir()
	vs := filepath.Join(home, "AppData", "Roaming", "Code", "User", "mcp.json")
	escrever(t, vs, `{"servers":{"rea":{"type":"stdio","command":"rea","args":["mcp"]}}}`)
	rodarSync(t, home, "--apply")

	obj := lerJSON(t, vs)
	if obj["mcpServers"] != nil {
		t.Fatalf("VS Code nao usa 'mcpServers': %v", obj)
	}
	lego, _ := obj["servers"].(map[string]any)["legoloop"].(map[string]any)
	if lego == nil || lego["type"] != "stdio" || lego["command"] != "legoloop" {
		t.Fatalf("legoloop deve entrar em servers com type stdio: %v", obj["servers"])
	}
}

func TestSyncConfigExtraViaFlag(t *testing.T) {
	home := t.TempDir()
	extra := filepath.Join(t.TempDir(), "perfil", ".claude.json")
	escrever(t, extra, `{"mcpServers":{"playwright":{"command":"npx","args":["-y","@playwright/mcp@latest"]}}}`)
	rodarSync(t, home, "--apply", "--config", extra)

	if lerJSON(t, extra)["mcpServers"].(map[string]any)["playwright"] != nil {
		t.Fatal("config extra deve virar lazy")
	}
	reg, _ := mcp.LoadRegistry(mcp.RegistryPath(home))
	if reg["playwright"].Command != "npx" {
		t.Fatalf("registro deve receber o servidor do config extra: %v", reg)
	}
}

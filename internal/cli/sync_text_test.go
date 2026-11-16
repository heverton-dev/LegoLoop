package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/mcp"
)

const codexFixture = `model = "gpt-5"

[mcp_servers.cbm]
command = "C:\tools\cbm.exe"
args = []
cwd = "C:\proj"
type = "stdio"

[mcp_servers.rea]
command = "cmd.exe"
args = [ "/d", "/c", "npx",
  "-y", "rea-agents@6.7.0", "mcp" ]
startup_timeout_sec = 30

[mcp_servers.rea.env]
JAVA_HOME = "C:\jdk"

[mcp_servers.remoto]
url = "https://exemplo.com/mcp"

[profiles.x]
model = "o3"
`

func TestSyncCodexTomlMoveStdioEPreservaResto(t *testing.T) {
	home := t.TempDir()
	cfg := filepath.Join(home, ".codex", "config.toml")
	escrever(t, cfg, codexFixture)
	rodarSync(t, home, "--apply")

	raw, _ := os.ReadFile(cfg)
	out := string(raw)
	for _, quer := range []string{`model = "gpt-5"`, "[mcp_servers.remoto]", "[profiles.x]", "[mcp_servers.legoloop]", `command = "legoloop"`} {
		if !strings.Contains(out, quer) {
			t.Fatalf("faltou %q em:\n%s", quer, out)
		}
	}
	if strings.Contains(out, "[mcp_servers.cbm]") || strings.Contains(out, "[mcp_servers.rea]") || strings.Contains(out, "rea.env") {
		t.Fatalf("stdio eager deve sair do toml:\n%s", out)
	}
	reg, _ := mcp.LoadRegistry(mcp.RegistryPath(home))
	rea := reg["rea"]
	if rea.Command != "cmd.exe" || len(rea.Args) != 6 || rea.Env["JAVA_HOME"] != `C:\jdk` {
		t.Fatalf("rea mal convertido: %+v", rea)
	}
	if reg["cbm"].Cwd != `C:\proj` || reg["cbm"].Command != `C:\tools\cbm.exe` {
		t.Fatalf("cbm mal convertido: %+v", reg["cbm"])
	}
}

func TestSyncCodexTomlIdempotente(t *testing.T) {
	home := t.TempDir()
	cfg := filepath.Join(home, ".codex", "config.toml")
	escrever(t, cfg, codexFixture)
	rodarSync(t, home, "--apply")
	rodarSync(t, home, "--apply")
	raw, _ := os.ReadFile(cfg)
	if n := strings.Count(string(raw), "[mcp_servers.legoloop]"); n != 1 {
		t.Fatalf("legoloop duplicado %d vezes:\n%s", n, raw)
	}
}

const hermesFixture = `model: x
mcp_servers:
  rea:
    command: cmd.exe
    args:
      - /d
      - /c
      - mcp
    env:
      JAVA_HOME: 'C:\jdk'
    enabled: true
  off:
    command: parado
    enabled: false
  remoto:
    url: https://exemplo.com/mcp

# =====
dashboard:
  x: 1
`

func TestSyncHermesYamlMoveStdioEPreservaResto(t *testing.T) {
	home := t.TempDir()
	cfg := filepath.Join(home, "AppData", "Local", "hermes", "config.yaml")
	escrever(t, cfg, hermesFixture)
	rodarSync(t, home, "--apply")

	raw, _ := os.ReadFile(cfg)
	out := string(raw)
	for _, quer := range []string{"model: x", "  legoloop:", "    command: legoloop", "      - mcp-router", "  off:", "  remoto:", "# =====", "dashboard:"} {
		if !strings.Contains(out, quer) {
			t.Fatalf("faltou %q em:\n%s", quer, out)
		}
	}
	if strings.Contains(out, "  rea:") {
		t.Fatalf("rea deve sair do yaml:\n%s", out)
	}
	reg, _ := mcp.LoadRegistry(mcp.RegistryPath(home))
	rea := reg["rea"]
	if rea.Command != "cmd.exe" || len(rea.Args) != 3 || rea.Env["JAVA_HOME"] != `C:\jdk` {
		t.Fatalf("rea mal convertido: %+v", rea)
	}
	if _, ok := reg["off"]; ok {
		t.Fatal("desabilitado nao entra no registro")
	}
}

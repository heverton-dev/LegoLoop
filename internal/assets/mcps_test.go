package assets

import (
	"path/filepath"
	"testing"
)

// Só nomes e metadados seguros: env, args, headers e query de URL NUNCA entram no inventário.
func TestScanMCPsNeverLeaksSecrets(t *testing.T) {
	cfg, home, project := testCfg(t)
	write(t, filepath.Join(home, ".claude.json"), `{
	  "mcpServers": {
	    "rea": {"command": "npx", "args": ["-y", "pacote-secreto-xyz"], "env": {"API_KEY": "SEGREDO123"}},
	    "remoto": {"type": "http", "url": "https://api.exemplo.com/mcp?token=ABC999", "headers": {"Authorization": "Bearer ZZZ"}}
	  },
	  "projects": {"C:/p": {"mcpServers": {"playwright": {"command": "node"}}}}
	}`)
	write(t, filepath.Join(project, ".mcp.json"), `{"mcpServers": {"local": {"command": "C:\\Program Files\\x\\srv.exe"}}}`)
	items := scanKind(t, cfg, KindMCP)

	rea := mustFind(t, items, "rea")
	if rea.Scope != "global" || rea.Origin != "claude" || rea.Meta["transport"] != "stdio" || rea.Meta["command"] != "npx" {
		t.Errorf("rea = %+v", rea)
	}
	rem := mustFind(t, items, "remoto")
	if rem.Meta["transport"] != "http" || rem.Meta["host"] != "api.exemplo.com" {
		t.Errorf("remoto = %+v", rem)
	}
	if p := mustFind(t, items, "playwright"); p.Scope != "project" {
		t.Errorf("playwright = %+v", p)
	}
	if l := mustFind(t, items, "local"); l.Scope != "project" || l.Path != filepath.Join(project, ".mcp.json") {
		t.Errorf("local = %+v", l)
	}
	mustNotLeak(t, items, "pacote-secreto-xyz", "SEGREDO123", "ABC999", "ZZZ")
}

func TestScanMCPsOtherHarnessFormats(t *testing.T) {
	cfg, home, _ := testCfg(t)
	write(t, filepath.Join(home, ".cursor", "mcp.json"), `{"mcpServers": {"cur": {"command": "node"}}}`)
	write(t, filepath.Join(home, ".gemini", "settings.json"), `{"mcpServers": {"gem": {"url": "https://g.exemplo.com/sse"}}}`)
	write(t, filepath.Join(home, "AppData", "Roaming", "Claude", "claude_desktop_config.json"), `{"mcpServers": {"desk": {"command": "uvx"}}}`)
	write(t, filepath.Join(home, ".config", "opencode", "opencode.json"), `{"mcp": {"oc": {"type": "local", "command": ["npx", "-y", "algo"]}}}`)
	items := scanKind(t, cfg, KindMCP)
	for name, origin := range map[string]string{"cur": "cursor", "gem": "gemini", "desk": "claude-desktop", "oc": "opencode"} {
		if it := mustFind(t, items, name); it.Origin != origin {
			t.Errorf("%s origin = %q, quer %q", name, it.Origin, origin)
		}
	}
	if oc := mustFind(t, items, "oc"); oc.Meta["command"] != "npx" {
		t.Errorf("opencode (comando em array) = %+v", oc)
	}
	if g := mustFind(t, items, "gem"); g.Meta["transport"] != "sse" && g.Meta["transport"] != "http" {
		t.Errorf("gem sem transporte remoto: %+v", g)
	}
}

func TestScanMCPsToleratesBrokenFiles(t *testing.T) {
	cfg, home, _ := testCfg(t)
	write(t, filepath.Join(home, ".claude.json"), `{ não é json`)
	write(t, filepath.Join(home, ".cursor", "mcp.json"), `{"mcpServers": ["errado"]}`)
	write(t, filepath.Join(home, ".gemini", "settings.json"), `{"mcpServers": {"ok": {"command": "x"}, "ruim": 42}}`)
	items := scanKind(t, cfg, KindMCP)
	mustFind(t, items, "ok")
	if _, found := find(items, "ruim"); found {
		t.Error("entrada de servidor inválida deve ser ignorada")
	}
}

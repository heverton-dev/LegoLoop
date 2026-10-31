package assets_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/assets"
)

func TestCascadeHooksAndMCPsAcrossDiscoveredProjects(t *testing.T) {
	home := t.TempDir()
	desktop := filepath.Join(home, "Desktop")

	// Projeto A com .claude/settings.json (project hook) e .claude/settings.local.json (local hook)
	projA := filepath.Join(desktop, "projeto-a")
	claudeDirA := filepath.Join(projA, ".claude")
	if err := os.MkdirAll(claudeDirA, 0755); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(projA, "package.json"), []byte("{}"), 0644)
	_ = os.WriteFile(filepath.Join(claudeDirA, "settings.json"), []byte(`{
		"hooks": {"Stop": [{"hooks": [{"type": "command", "command": "npm test"}]}]}
	}`), 0644)
	_ = os.WriteFile(filepath.Join(claudeDirA, "settings.local.json"), []byte(`{
		"hooks": {"Stop": [{"hooks": [{"type": "command", "command": "echo secret-local"}]}]}
	}`), 0644)

	// Projeto B com .mcp.json (project mcp)
	projB := filepath.Join(desktop, "projeto-b")
	if err := os.MkdirAll(projB, 0755); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(projB, "go.mod"), []byte("module b"), 0644)
	_ = os.WriteFile(filepath.Join(projB, ".mcp.json"), []byte(`{
		"mcpServers": {
			"db-tool": {"command": "docker", "args": ["run", "db"]}
		}
	}`), 0644)

	cfg := assets.Config{
		Home:           home,
		AutoDiscover:   true,
		DiscoveryRoots: []string{desktop},
	}
	scanner := assets.NewScanner(cfg)
	snap := scanner.Scan(context.Background())

	// 1. Verifica hooks nos escopos project e local
	var foundProjHook, foundLocalHook bool
	for _, it := range snap.Of(assets.KindHook) {
		if it.Command == "npm test" && it.Scope == "project" {
			foundProjHook = true
		}
		if it.Command == "echo secret-local" && it.Scope == "local" {
			foundLocalHook = true
		}
	}
	if !foundProjHook {
		t.Fatal("hook de escopo 'project' do projeto-a não encontrado")
	}
	if !foundLocalHook {
		t.Fatal("hook de escopo 'local' do projeto-a não encontrado")
	}

	// 2. Verifica MCP no escopo project
	var foundProjMCP bool
	for _, it := range snap.Of(assets.KindMCP) {
		if it.Name == "db-tool" && it.Scope == "project" {
			foundProjMCP = true
		}
	}
	if !foundProjMCP {
		t.Fatal("MCP 'db-tool' de escopo 'project' do projeto-b não encontrado")
	}
}

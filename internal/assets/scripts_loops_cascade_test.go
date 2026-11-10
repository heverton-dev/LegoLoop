package assets_test

import (
	"context"
	"legoloop/internal/caminhos"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/assets"
)

func TestCascadeScriptsAndLoopsAcrossDiscoveredProjects(t *testing.T) {
	home := t.TempDir()
	desktop := filepath.Join(home, "Desktop")

	// Projeto C com script em tools/ e loop em loops/
	projC := filepath.Join(desktop, "projeto-c")
	toolsDir := filepath.Join(projC, "tools")
	loopsDir := caminhos.DirLoopsEm(projC)
	if err := os.MkdirAll(toolsDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(loopsDir, 0755); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(projC, "go.mod"), []byte("module c"), 0644)
	_ = os.WriteFile(filepath.Join(toolsDir, "lint.ps1"), []byte("Write-Host lint"), 0644)
	_ = os.WriteFile(filepath.Join(loopsDir, "build.loop.json"), []byte(`{
		"id": "build-loop",
		"name": "Build Loop",
		"version": "2.0.0",
		"trigger": {"type": "manual"},
		"nodes": [{"id": "n1", "output": {"artifact": "BUILD.log"}}]
	}`), 0644)

	cfg := assets.Config{
		Home:           home,
		AutoDiscover:   true,
		DiscoveryRoots: []string{desktop},
	}
	scanner := assets.NewScanner(cfg)
	snap := scanner.Scan(context.Background())

	// 1. Verifica script em tools/ com scope "project"
	var foundScript bool
	for _, it := range snap.Of(assets.KindScript) {
		if it.Name == "lint.ps1" && it.Scope == "project" {
			foundScript = true
		}
	}
	if !foundScript {
		t.Fatal("script 'lint.ps1' de tools/ do projeto-c não foi inventariado com scope 'project'")
	}

	// 2. Verifica loop em loops/ com scope "project"
	var foundLoop bool
	for _, it := range snap.Of(assets.KindLoop) {
		if it.Name == "Build Loop" && it.Scope == "project" {
			foundLoop = true
		}
	}
	if !foundLoop {
		t.Fatal("loop 'Build Loop' do projeto-c não foi inventariado com scope 'project'")
	}
}

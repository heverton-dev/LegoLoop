package assets_test

import (
	"context"
	"legoloop/internal/caminhos"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/assets"
)

func TestCascadeScriptsAndProcessosAcrossDiscoveredProjects(t *testing.T) {
	home := t.TempDir()
	desktop := filepath.Join(home, "Desktop")

	projC := filepath.Join(desktop, "projeto-c")
	toolsDir := filepath.Join(projC, "tools")
	processosDir := filepath.Join(projC, caminhos.DirProcessos)
	if err := os.MkdirAll(toolsDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(processosDir, 0755); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(projC, "go.mod"), []byte("module c"), 0644)
	_ = os.WriteFile(filepath.Join(toolsDir, "lint.ps1"), []byte("Write-Host lint"), 0644)
	_ = os.WriteFile(filepath.Join(processosDir, "build"+caminhos.SufixoStudio), []byte(`{
		"id": "build-processo",
		"name": "Build Processo",
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

	var foundScript bool
	for _, it := range snap.Of(assets.KindScript) {
		if it.Name == "lint.ps1" && it.Scope == "project" {
			foundScript = true
		}
	}
	if !foundScript {
		t.Fatal("script 'lint.ps1' de tools/ do projeto-c não foi inventariado com scope 'project'")
	}

	var foundProcesso bool
	for _, it := range snap.Of(kindProcesso) {
		if it.Name == "Build Processo" && it.Scope == "project" {
			foundProcesso = true
		}
	}
	if !foundProcesso {
		t.Fatal("Processo 'Build Processo' do projeto-c não foi inventariado com scope 'project'")
	}
}

package assets_test

import (
	"context"
	"legoloop/internal/caminhos"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/assets"
)

func TestScaffoldsLayer(t *testing.T) {
	home := t.TempDir()
	scaffoldsDir := filepath.Join(home, ".aidd", "scaffolds")

	// 1. Skill de scaffold
	scaffoldSkill := filepath.Join(scaffoldsDir, "starter-kit", "skills", "setup")
	if err := os.MkdirAll(scaffoldSkill, 0755); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(scaffoldSkill, "SKILL.md"), []byte("---\nname: Starter Setup\ndescription: skill de scaffold\n---\n"), 0644)

	// 2. Loop de scaffold
	scaffoldLoopDir := filepath.Join(scaffoldsDir, "starter-kit", caminhos.DirLoops)
	if err := os.MkdirAll(scaffoldLoopDir, 0755); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(scaffoldLoopDir, "starter.loop.json"), []byte(`{
		"id": "starter-loop",
		"name": "Starter Loop",
		"version": "2.0.0",
		"trigger": {"type": "manual"},
		"nodes": [{"id": "n1", "output": {"artifact": "A.md"}}]
	}`), 0644)

	// 3. Hook de scaffold
	scaffoldHookDir := filepath.Join(scaffoldsDir, "starter-kit", ".claude")
	if err := os.MkdirAll(scaffoldHookDir, 0755); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(scaffoldHookDir, "settings.json"), []byte(`{
		"hooks": {
			"Stop": [{"hooks": [{"type": "command", "command": "scaffold-check"}]}]
		}
	}`), 0644)

	scanner := assets.NewScanner(assets.Config{Home: home})
	snap := scanner.Scan(context.Background())

	// Verifica se a skill de scaffold foi inventariada com scope "scaffold"
	var foundSkill bool
	for _, it := range snap.Of(assets.KindSkill) {
		if it.Name == "Starter Setup" {
			foundSkill = true
			if it.Scope != "scaffold" {
				t.Errorf("escopo esperado 'scaffold', obteve %s", it.Scope)
			}
		}
	}
	if !foundSkill {
		t.Fatal("skill de scaffold não foi inventariada")
	}

	// Verifica se o loop de scaffold foi inventariado com scope "scaffold"
	var foundLoop bool
	for _, it := range snap.Of(assets.KindLoop) {
		if it.Name == "Starter Loop" {
			foundLoop = true
			if it.Scope != "scaffold" {
				t.Errorf("escopo de loop esperado 'scaffold', obteve %s", it.Scope)
			}
		}
	}
	if !foundLoop {
		t.Fatal("loop de scaffold não foi inventariado")
	}

	// Verifica se o hook de scaffold foi inventariado com scope "scaffold"
	var foundHook bool
	for _, it := range snap.Of(assets.KindHook) {
		if it.Name == "Stop" && it.Command == "scaffold-check" {
			foundHook = true
			if it.Scope != "scaffold" {
				t.Errorf("escopo de hook esperado 'scaffold', obteve %s", it.Scope)
			}
		}
	}
	if !foundHook {
		t.Fatal("hook de scaffold não foi inventariado")
	}
}

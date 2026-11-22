package processo_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/processo"
)

// TestDetectaTodosOsHarnessesDoUsuario prova que todo harness com pasta de skills própria é
// coberto pela política lazy (índice no lugar das skills), inclusive os de layout fora do padrão
// ~/.<nome>/skills: opencode e windsurf, crush/devin/goose em ~/.config e pi em ~/.pi/agent.
func TestDetectaTodosOsHarnessesDoUsuario(t *testing.T) {
	home := t.TempDir()
	esperado := map[string]string{
		"cline":    filepath.Join(".cline", "skills"),
		"hermes":   filepath.Join(".hermes", "skills"),
		"mimocode": filepath.Join(".mimocode", "skills"),
		"augment":  filepath.Join(".augment", "skills"),
		"junie":    filepath.Join(".junie", "skills"),
		"grok":     filepath.Join(".grok", "skills"),
		"copilot":  filepath.Join(".copilot", "skills"),
		"trae":     filepath.Join(".trae", "skills"),
		"roo":      filepath.Join(".roo", "skills"),
		"continue": filepath.Join(".continue", "skills"),
		"qwen":     filepath.Join(".qwen", "skills"),
		"windsurf": filepath.Join(".codeium", "windsurf", "skills"),
		"opencode": filepath.Join(".config", "opencode", "skills"),
		"crush":    filepath.Join(".config", "crush", "skills"),
		"devin":    filepath.Join(".config", "devin", "skills"),
		"goose":    filepath.Join(".config", "goose", "skills"),
		"pi":       filepath.Join(".pi", "agent", "skills"),
	}
	for nome, skills := range esperado {
		if err := os.MkdirAll(filepath.Join(home, filepath.Dir(skills)), 0o755); err != nil {
			t.Fatalf("montar %s: %v", nome, err)
		}
	}
	got := map[string]string{}
	for _, h := range processo.DetectarHarnesses(processo.OSFS(), home) {
		got[h.Nome] = h.Skills
	}
	for nome, skills := range esperado {
		want := filepath.Join(home, skills)
		if got[nome] != want {
			t.Errorf("harness %q: pasta de skills = %q, quer %q", nome, got[nome], want)
		}
	}
}

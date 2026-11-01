package studio

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/loop"
)

func TestMaterializeCreatesDAGAndOrchestratorSkillAndRegistersInLoopStore(t *testing.T) {
	s, _ := newServer(t)
	vaultDir := t.TempDir()
	s.Vault = vaultDir
	loopStore, err := loop.NewStore(filepath.Join(vaultDir, "loops"))
	if err != nil {
		t.Fatal(err)
	}
	s.Loops = loopStore

	h := s.Handler()

	body := `{
		"name": "minha-esteira",
		"nodes": [
			{"id":"n1","block":"interview","routing":{"mode":"suggested","model":"sonnet","fallback":"qwen2.5-coder:7b"}},
			{"id":"n2","block":"aidd-spec","routing":{"mode":"auto","combo":"code-pro"}}
		],
		"edges": [{"from":"n1","to":"n2"}]
	}`

	rec := do(t, h, "POST", "/api/materialize", body)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body = %s", rec.Code, rec.Body.String())
	}

	// 1. pipeline.dag.json deve existir no ProjectDir
	dagFile := filepath.Join(s.ProjectDir, PipelineFile)
	if _, err := os.Stat(dagFile); err != nil {
		t.Fatalf("pipeline.dag.json não foi gravado: %v", err)
	}

	// 2. Skill orquestradora deve existir no Vault (< 20 linhas)
	skillFile := filepath.Join(vaultDir, "skills", "loop-minha-esteira", "SKILL.md")
	content, err := os.ReadFile(skillFile)
	if err != nil {
		t.Fatalf("SKILL.md orquestradora não foi gerada: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) >= 20 {
		t.Fatalf("skill orquestradora deve ter < 20 linhas, tem %d", len(lines))
	}
	if !strings.Contains(string(content), "legoloop run") {
		t.Fatalf("skill orquestradora deve invocar legoloop run: %s", string(content))
	}

	// 3. Loop deve estar registrado no Store do acervo
	saved, err := s.Loops.Get("minha-esteira")
	if err != nil {
		t.Fatalf("loop não foi salvo no acervo: %v", err)
	}
	if saved.Name != "minha-esteira" {
		t.Errorf("nome do loop salvo = %q", saved.Name)
	}
}

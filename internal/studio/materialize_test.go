package studio

import (
	"legoloop/internal/caminhos"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/studiodoc"
)

func TestMaterializeCreatesDAGAndOrchestratorSkillAndRegistersInProcessoStore(t *testing.T) {
	s, _ := newServer(t)
	vaultDir := t.TempDir()
	s.Vault = vaultDir
	acervo, err := studiodoc.NewStore(caminhos.DirLoopsEm(vaultDir))
	if err != nil {
		t.Fatal(err)
	}
	s.Processos = acervo

	h := s.Handler()

	body := `{
		"name": "meu-processo",
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

	// 1. o arquivo DAG deve existir no ProjectDir
	dagFile := filepath.Join(s.ProjectDir, ArquivoDAGGravado)
	if _, err := os.Stat(dagFile); err != nil {
		t.Fatalf("%s não foi gravado: %v", caminhos.ArquivoDAG, err)
	}

	// 2. Skill orquestradora deve existir no Vault (< 20 linhas)
	skillFile := filepath.Join(vaultDir, "skills", "loop-meu-processo", "SKILL.md")
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

	saved, err := s.Processos.Get("meu-processo")
	if err != nil {
		t.Fatalf("Processo não foi salvo no acervo: %v", err)
	}
	if saved.Name != "meu-processo" {
		t.Errorf("nome do Processo salvo = %q", saved.Name)
	}
}

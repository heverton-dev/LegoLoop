package studio

import (
	"legoloop/internal/caminhos"
	"legoloop/internal/dag"
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
	skillFile := filepath.Join(vaultDir, "skills", "processo-meu-processo", "SKILL.md")
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

func TestWriteOrchestratorSkillSlugJailAndErrors(t *testing.T) {
	vault := t.TempDir()
	s, _ := newServer(t)
	s.Vault = vault

	// 1. Slugs inválidos devem retornar erro e nunca gravar no disco.
	for _, badName := range []string{"", "../escape", `a\b`, "a/b", "foo..bar", "NomeComMaiuscula", "nome com espaco"} {
		err := s.writeOrchestratorSkill(&dag.DAG{Name: badName})
		if err == nil {
			t.Errorf("esperava erro para slug inválido %q, obteve nil", badName)
		}
	}

	// 2. Endpoint /api/materialize também bloqueia nome não-slug sem gravar fora do cofre.
	badBody := `{
		"name": "../escape",
		"nodes": [
			{"id":"n1","block":"interview","routing":{"mode":"auto","combo":"t-chat-128k"}}
		],
		"edges": []
	}`
	recBad := do(t, s.Handler(), "POST", "/api/materialize", badBody)
	if recBad.Code == http.StatusOK {
		t.Fatalf("POST /api/materialize deveria falhar para name inválido, obteve 200")
	}

	// 3. Symlink/junction em vault/skills apontando para fora deve ser bloqueado.
	vaultLink := t.TempDir()
	outside := t.TempDir()
	skillsLink := filepath.Join(vaultLink, "skills")
	createHardeningLink(t, outside, skillsLink)
	sLink := &Server{Vault: vaultLink}
	if err := sLink.writeOrchestratorSkill(&dag.DAG{Name: "via-link"}); err == nil {
		t.Fatalf("writeOrchestratorSkill deveria bloquear symlink em skills")
	}
	if _, err := os.Stat(filepath.Join(outside, "processo-via-link", "SKILL.md")); err == nil {
		t.Fatalf("arquivo foi criado através do symlink fora do cofre")
	}

	// 4. Erro de I/O quando skills é um arquivo regular não deve ser engolido silenciosamente.
	vaultFileErr := t.TempDir()
	if err := os.WriteFile(filepath.Join(vaultFileErr, "skills"), []byte("bloqueio"), 0o644); err != nil {
		t.Fatal(err)
	}
	sErr := &Server{Vault: vaultFileErr}
	if err := sErr.writeOrchestratorSkill(&dag.DAG{Name: "teste-erro"}); err == nil {
		t.Fatalf("esperava erro de escrita quando skills não é diretório")
	}
}

package studio

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/interview"
)

func TestPieceActions_FavoriteAndUpgrade(t *testing.T) {
	s, _ := newServer(t)
	s.Vault = t.TempDir()

	mockLLM := &mockStudioLLM{
		responses: []string{
			"HIPÓTESE: Elevar a peça para padrão ouro.\nCONFIANÇA: 60%\nPERGUNTA: Quais guardrails faltam?\nPALPITE: Zero Stubs",
		},
	}
	s.Interview = interview.NewEngine(mockLLM)
	h := s.Handler()

	// 1. Testar Favoritar
	rec := do(t, h, http.MethodPost, "/api/assets/skill/test-skill/favorite", "{}")
	if rec.Code != http.StatusOK {
		t.Fatalf("esperava 200 em favorite, obteve %d: %s", rec.Code, rec.Body.String())
	}

	var favResp map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &favResp)
	if favResp["favorited"] != true {
		t.Errorf("esperava favorited = true, obteve %v", favResp["favorited"])
	}

	// 2. Testar Iniciar Upgrade
	rec = do(t, h, http.MethodPost, "/api/assets/skill/test-skill/upgrade", "{}")
	if rec.Code != http.StatusOK {
		t.Fatalf("esperava 200 em upgrade, obteve %d: %s", rec.Code, rec.Body.String())
	}

	var upgResp map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &upgResp)
	if upgResp["status"] != "UPGRADE_INITIATED" {
		t.Errorf("esperava UPGRADE_INITIATED, obteve %v", upgResp["status"])
	}
}

func TestPieceActions_Delete(t *testing.T) {
	s, _ := newServer(t)
	h := s.Handler()

	tempFile := filepath.Join(s.ProjectDir, "deleteme.md")
	os.WriteFile(tempFile, []byte("lixo"), 0644)

	// Se não existe no inventário de assets, retorna 404
	rec := do(t, h, http.MethodDelete, "/api/assets/skill/non-existent?mode=local", "")
	if rec.Code != http.StatusNotFound {
		t.Errorf("esperava 404 para peça inexistente, obteve %d", rec.Code)
	}
}

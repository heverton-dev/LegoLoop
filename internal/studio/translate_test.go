package studio

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTranslateEndpointCacheAndImmutability(t *testing.T) {
	home := t.TempDir()
	skillDir := filepath.Join(home, ".claude", "skills", "test-skill")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}
	skillFile := filepath.Join(skillDir, "SKILL.md")
	originalContent := "---\nname: english-skill\ndescription: A tool to do magic\n---\n# English Content\nDo this and that."
	if err := os.WriteFile(skillFile, []byte(originalContent), 0600); err != nil {
		t.Fatal(err)
	}

	srv := &Server{Home: home, Host: host}
	reqBody, _ := json.Marshal(map[string]string{"path": skillFile})
	req := httptest.NewRequest(http.MethodPost, "/api/assets/translate", bytes.NewReader(reqBody))
	req.Host = host
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.handleAssetTranslate(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("esperava 200 OK, obteve %d: %s", w.Code, w.Body.String())
	}

	var resp TranslationRecord
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("erro decodificando json: %v", err)
	}
	if resp.SHA256 == "" || resp.WhatItIs == "" {
		t.Errorf("resposta de tradução incompleta: %+v", resp)
	}
	if strings.Contains(resp.SourcePath, home) || filepath.IsAbs(resp.SourcePath) {
		t.Errorf("source_path vazou caminho absoluto: %q", resp.SourcePath)
	}

	// Verifica se arquivo original foi mantido 100% intacto
	current, err := os.ReadFile(skillFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(current) != originalContent {
		t.Errorf("ARQUIVO ORIGINAL FOI MODIFICADO! Esperava %q, obteve %q", originalContent, string(current))
	}
}

func TestTranslateEndpointRejectsPathTraversal(t *testing.T) {
	srv := &Server{Home: t.TempDir(), ProjectDir: t.TempDir(), Vault: t.TempDir(), Host: host}
	reqBody, _ := json.Marshal(map[string]string{"path": "../../etc/passwd"})
	req := httptest.NewRequest(http.MethodPost, "/api/assets/translate", bytes.NewReader(reqBody))
	req.Host = host
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.handleAssetTranslate(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("esperava 400 Bad Request para traversal, obteve %d", w.Code)
	}
}

func TestTranslateEndpointBlocksOutsideRevealRootsAndSymlinks(t *testing.T) {
	proj := t.TempDir()
	home := t.TempDir()
	vault := t.TempDir()
	outside := t.TempDir()

	secretFile := filepath.Join(outside, "segredo.txt")
	if err := os.WriteFile(secretFile, []byte("SEGREDO_CONFIDENCIAL"), 0o600); err != nil {
		t.Fatal(err)
	}

	srv := &Server{ProjectDir: proj, Home: home, Vault: vault, Host: host}
	h := srv.Handler()

	// 1. Caminho absoluto fora de revealRoots deve retornar 400 ou 403 e nunca ler o arquivo.
	payload, _ := json.Marshal(map[string]string{"path": secretFile})
	rec := do(t, h, http.MethodPost, "/api/assets/translate", string(payload))
	if rec.Code != http.StatusBadRequest && rec.Code != http.StatusForbidden {
		t.Fatalf("caminho absoluto fora das raízes deveria retornar 400/403, obteve %d: %s", rec.Code, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "SEGREDO_CONFIDENCIAL") || strings.Contains(rec.Body.String(), outside) {
		t.Fatalf("resposta vazou conteúdo ou caminho fora das raízes: %s", rec.Body.String())
	}

	// 2. Symlink/junction dentro do projeto apontando para fora deve ser bloqueado.
	linkPath := filepath.Join(proj, "escape-link")
	createHardeningLink(t, outside, linkPath)
	linkPayload, _ := json.Marshal(map[string]string{"path": filepath.Join(linkPath, "segredo.txt")})
	recLink := do(t, h, http.MethodPost, "/api/assets/translate", string(linkPayload))
	if recLink.Code != http.StatusBadRequest && recLink.Code != http.StatusForbidden {
		t.Fatalf("symlink para fora deveria retornar 400/403, obteve %d: %s", recLink.Code, recLink.Body.String())
	}

	// 3. Caminho válido em ProjectDir via <projeto>/ funciona e mascara SourcePath.
	validFile := filepath.Join(proj, "skills", "local", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(validFile), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(validFile, []byte("# Skill Local"), 0o600); err != nil {
		t.Fatal(err)
	}
	okPayload, _ := json.Marshal(map[string]string{"path": "<projeto>/skills/local/SKILL.md"})
	recOK := do(t, h, http.MethodPost, "/api/assets/translate", string(okPayload))
	if recOK.Code != http.StatusOK {
		t.Fatalf("esperava 200 para caminho permitido, obteve %d: %s", recOK.Code, recOK.Body.String())
	}
	if strings.Contains(recOK.Body.String(), proj) || strings.Contains(recOK.Body.String(), home) {
		t.Fatalf("resposta de tradução vazou caminho absoluto: %s", recOK.Body.String())
	}

	// 4. Mensagem de erro não deve vazar caminhos de Home/ProjectDir/Vault.
	missingPayload, _ := json.Marshal(map[string]string{"path": filepath.Join(home, "nao-existe.md")})
	recMiss := do(t, h, http.MethodPost, "/api/assets/translate", string(missingPayload))
	if recMiss.Code != http.StatusBadRequest {
		t.Fatalf("esperava 400 para ausente, obteve %d", recMiss.Code)
	}
	if strings.Contains(recMiss.Body.String(), home) || strings.Contains(recMiss.Body.String(), proj) {
		t.Fatalf("erro vazou caminho absoluto: %s", recMiss.Body.String())
	}
}

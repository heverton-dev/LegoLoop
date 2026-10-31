package studio

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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

	srv := &Server{Home: home}
	reqBody, _ := json.Marshal(map[string]string{"path": skillFile})
	req := httptest.NewRequest(http.MethodPost, "/api/assets/translate", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.handleAssetTranslate(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("esperava 200 OK, obteve %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		SHA256          string `json:"sha256"`
		WhatItIs        string `json:"what_it_is"`
		WhatItDoes      string `json:"what_it_does"`
		WhatItDelivers  string `json:"what_it_delivers"`
		FullContentPTBR string `json:"full_content_ptbr"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("erro decodificando json: %v", err)
	}
	if resp.SHA256 == "" || resp.WhatItIs == "" {
		t.Errorf("resposta de tradução incompleta: %+v", resp)
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
	srv := &Server{Home: t.TempDir()}
	reqBody, _ := json.Marshal(map[string]string{"path": "../../etc/passwd"})
	req := httptest.NewRequest(http.MethodPost, "/api/assets/translate", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.handleAssetTranslate(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("esperava 400 Bad Request para traversal, obteve %d", w.Code)
	}
}

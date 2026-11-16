package studio

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAssetInspectRealMetadataAndSource(t *testing.T) {
	s, _, proj := newAssetsServer(t)
	h := s.Handler()
	skillDir := filepath.Join(proj, "skills", "minha-skill-real")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	skillMD := "---\nname: minha-skill-real\ndescription: Valida contratos reais do projeto sem mocks.\n---\n\n# Guia Prático\nExecuta verificação em disco.\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skillMD), 0o644); err != nil {
		t.Fatal(err)
	}
	s.Assets.Refresh(t.Context())

	// 1. Inspecionar skill real por nome/caminho mascarado
	rec := do(t, h, http.MethodGet, "/api/assets/inspect?name=minha-skill-real&kind=skill", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("inspect skill falhou: %d %s", rec.Code, rec.Body.String())
	}
	var out struct {
		OK            bool     `json:"ok"`
		Name          string   `json:"name"`
		Description   string   `json:"description"`
		WhatItIs      string   `json:"what_it_is"`
		SourceContent string   `json:"source_content"`
		Path          string   `json:"path"`
		Tier          string   `json:"tier"`
		Score         string   `json:"score"`
		Sections      []string `json:"sections"`
	}
	decode(t, rec, &out)
	if !out.OK || !strings.Contains(out.SourceContent, "Guia Prático") {
		t.Fatalf("inspect não retornou conteúdo real: %+v", out)
	}
	if !strings.Contains(out.Description, "Valida contratos reais") {
		t.Fatalf("inspect não extraiu descrição real do frontmatter: %+v", out)
	}
	if strings.Contains(out.Path, s.ProjectDir) {
		t.Fatalf("inspect vazou caminho absoluto: %s", out.Path)
	}

	// 2. Inspecionar trigger nativo do catálogo retorna metadados reais do grafo
	recTrig := do(t, h, http.MethodGet, "/api/assets/inspect?id=trigger-manual&kind=trigger_manual", "")
	if recTrig.Code != http.StatusOK {
		t.Fatalf("inspect trigger falhou: %d %s", recTrig.Code, recTrig.Body.String())
	}
	var trigOut map[string]any
	decode(t, recTrig, &trigOut)
	if trigOut["kind"] != "trigger_manual" {
		t.Fatalf("inspect trigger kind incorreto: %+v", trigOut)
	}
}

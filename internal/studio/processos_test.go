package studio

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/processo"
)

func TestStudioProcessos(t *testing.T) {
	tmp := t.TempDir()
	vault := filepath.Join(tmp, "vault")
	docs := filepath.Join(tmp, "docs")
	_ = os.MkdirAll(filepath.Join(vault, "loops"), 0o755)
	_ = os.MkdirAll(filepath.Join(vault, "skills"), 0o755)
	_ = os.MkdirAll(docs, 0o755)

	procA := processo.Processo{
		Nome:     "legoloop-alfa",
		Versao:   "1.0.0",
		Natureza: "software",
		Objetivo: "Processo de teste alfa",
		Etapas: []processo.Etapa{
			{ID: "e01", Tipo: "coletar", Peca: "skill-a", Titulo: "Etapa 1", Casos: []string{"requisitos"}, Produz: processo.Artefato{Caminho: "out/01.md"}, Gate: "exit_code == 0"},
			{ID: "e02", Tipo: "construir", Peca: "skill-b", Titulo: "Etapa 2", Casos: []string{"tdd"}, Depende: []string{"e01"}, Produz: processo.Artefato{Caminho: "out/02.md"}, Gate: "exit_code == 0"},
		},
	}
	procB := processo.Processo{
		Nome:     "legoloop-beta",
		Versao:   "1.0.0",
		Natureza: "software",
		Objetivo: "Processo de teste beta",
		Etapas: []processo.Etapa{
			{ID: "e01", Tipo: "verificar", Peca: "skill-c", Titulo: "Etapa B1", Casos: []string{"qualidade"}, Produz: processo.Artefato{Caminho: "out/beta.md"}, Gate: "exit_code == 0"},
		},
	}
	store := processo.Store{Roots: processo.Roots{Vault: vault, Docs: docs}}
	dataA, _ := json.Marshal(procA)
	_ = os.WriteFile(processo.CaminhoNoCofre(vault, procA.Nome), dataA, 0o644)
	_ = os.WriteFile(filepath.Join(vault, "loops", procA.Nome+".loop.json"), []byte(`{"version":"1.0.0"}`), 0o644)

	dataB, _ := json.Marshal(procB)
	_ = os.WriteFile(processo.CaminhoNoCofre(vault, procB.Nome), dataB, 0o644)
	_ = os.WriteFile(filepath.Join(vault, "loops", procB.Nome+".loop.json"), []byte(`{"version":"1.0.0"}`), 0o644)

	srv, _ := newServer(t)
	srv.Vault = vault
	srv.ProjectDir = docs
	srv.StoreOverride = &store

	t.Run("GET /api/processos lista os processos", func(t *testing.T) {
		rec := do(t, srv.Handler(), http.MethodGet, "/api/processos", "")
		if rec.Code != http.StatusOK {
			t.Fatalf("esperava 200, obteve %d: %s", rec.Code, rec.Body.String())
		}
		var res struct {
			Vault     string `json:"vault"`
			Processos []struct {
				Nome     string `json:"nome"`
				Versao   string `json:"versao"`
				Natureza string `json:"natureza"`
			} `json:"processos"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
			t.Fatalf("resposta inválida: %v", err)
		}
		if len(res.Processos) != 2 {
			t.Fatalf("esperava 2 processos, obteve %d", len(res.Processos))
		}
	})

	t.Run("POST /api/processos/refinar preview e apply", func(t *testing.T) {
		bodyPrev := map[string]any{
			"name":   "legoloop-alfa",
			"etapa":  "e02",
			"peca":   "skill-nova",
			"caso":   []string{"tdd"},
			"titulo": "Nova Construção",
			"apply":  false,
		}
		buf, _ := json.Marshal(bodyPrev)
		rec := do(t, srv.Handler(), http.MethodPost, "/api/processos/refinar", string(buf))
		if rec.Code != http.StatusOK {
			t.Fatalf("esperava 200 no preview, obteve %d: %s", rec.Code, rec.Body.String())
		}
		var prev struct {
			De     string `json:"de"`
			Para   string `json:"para"`
			Diff   string `json:"diff"`
			Cartao struct {
				Aviso  string   `json:"aviso"`
				Passos []string `json:"passos"`
			} `json:"cartao"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &prev); err != nil {
			t.Fatalf("json preview inválido: %v", err)
		}
		if prev.De != "1.0.0" || prev.Para != "1.1.0" {
			t.Errorf("versões esperadas 1.0.0 -> 1.1.0, obteve %s -> %s", prev.De, prev.Para)
		}
		if prev.Cartao.Aviso == "" {
			t.Errorf("cartão STE-PT deve ter aviso de risco antes dos passos")
		}

		bodyApply := map[string]any{
			"name":   "legoloop-alfa",
			"etapa":  "e02",
			"peca":   "skill-nova",
			"caso":   []string{"tdd"},
			"titulo": "Nova Construção",
			"apply":  true,
		}
		bufApp, _ := json.Marshal(bodyApply)
		recApp := do(t, srv.Handler(), http.MethodPost, "/api/processos/refinar", string(bufApp))
		if recApp.Code != http.StatusOK {
			t.Fatalf("esperava 200 no apply, obteve %d: %s", recApp.Code, recApp.Body.String())
		}
		pAtual, _ := processo.Ler(processo.CaminhoNoCofre(vault, "legoloop-alfa"))
		if pAtual.Versao != "1.1.0" {
			t.Errorf("esperava versão 1.1.0 após refinar, obteve %s", pAtual.Versao)
		}
	})

	t.Run("POST /api/processos/compor preview e apply", func(t *testing.T) {
		bodyPrev := map[string]any{
			"name":     "legoloop-alfa",
			"composto": "legoloop-beta@1.0.0",
			"titulo":   "Subloop Beta",
			"depende":  []string{"e02"},
			"apply":    false,
		}
		buf, _ := json.Marshal(bodyPrev)
		rec := do(t, srv.Handler(), http.MethodPost, "/api/processos/compor", string(buf))
		if rec.Code != http.StatusOK {
			t.Fatalf("esperava 200 no preview compor, obteve %d: %s", rec.Code, rec.Body.String())
		}
		var prev struct {
			Cartao struct {
				Aviso string `json:"aviso"`
			} `json:"cartao"`
			Diff string `json:"diff"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &prev); err != nil {
			t.Fatalf("json preview compor inválido: %v", err)
		}

		bodyApply := map[string]any{
			"name":     "legoloop-alfa",
			"composto": "legoloop-beta@1.0.0",
			"titulo":   "Subloop Beta",
			"depende":  []string{"e02"},
			"apply":    true,
		}
		bufApp, _ := json.Marshal(bodyApply)
		recApp := do(t, srv.Handler(), http.MethodPost, "/api/processos/compor", string(bufApp))
		if recApp.Code != http.StatusOK {
			t.Fatalf("esperava 200 no apply compor, obteve %d: %s", recApp.Code, recApp.Body.String())
		}
	})

	t.Run("POST /api/processos/remover preview e apply", func(t *testing.T) {
		bodyPrev := map[string]any{
			"name":  "legoloop-alfa",
			"apply": false,
		}
		buf, _ := json.Marshal(bodyPrev)
		rec := do(t, srv.Handler(), http.MethodPost, "/api/processos/remover", string(buf))
		if rec.Code != http.StatusOK {
			t.Fatalf("esperava 200 no preview remover, obteve %d: %s", rec.Code, rec.Body.String())
		}
		var prev struct {
			Cartao struct {
				Aviso string `json:"aviso"`
			} `json:"cartao"`
			Arquivos []string `json:"arquivos"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &prev); err != nil {
			t.Fatalf("json preview remover inválido: %v", err)
		}
		if len(prev.Arquivos) == 0 {
			t.Errorf("preview remover deve listar arquivos a apagar")
		}

		bodyApply := map[string]any{
			"name":  "legoloop-alfa",
			"apply": true,
		}
		bufApp, _ := json.Marshal(bodyApply)
		recApp := do(t, srv.Handler(), http.MethodPost, "/api/processos/remover", string(bufApp))
		if recApp.Code != http.StatusOK {
			t.Fatalf("esperava 200 no apply remover, obteve %d: %s", recApp.Code, recApp.Body.String())
		}
		if _, err := os.Stat(processo.CaminhoNoCofre(vault, "legoloop-alfa")); !os.IsNotExist(err) {
			t.Errorf("arquivo do processo deveria ter sido removido")
		}
	})
}

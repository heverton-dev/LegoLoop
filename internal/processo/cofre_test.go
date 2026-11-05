package processo_test

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"legoloop/internal/processo"
)

func TestCofreProcessos(t *testing.T) {
	t.Run("cofre sem loops devolve vazio sem erro", func(t *testing.T) {
		procs, err := processo.ProcessosDoCofre(t.TempDir())
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if len(procs) != 0 {
			t.Fatalf("esperava vazio, veio %v", procs)
		}
	})

	t.Run("lê v2, v1 e anota erro em json corrompido", func(t *testing.T) {
		vault := t.TempDir()
		loopsDir := filepath.Join(vault, "loops")
		if err := os.MkdirAll(loopsDir, 0o755); err != nil {
			t.Fatal(err)
		}

		// v2 válido
		p := processo.Processo{
			Nome:     "legoloop-alfa",
			Versao:   "1.2.0",
			Natureza: "dados",
			Objetivo: "Objetivo do alfa.",
			CriadoEm: time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC),
			Etapas: []processo.Etapa{
				{ID: "E01", Titulo: "Etapa 1", Tipo: "construir", Peca: "tdd-workflow",
					Produz: processo.Artefato{Caminho: "out/1.md", Tipo: "markdown"}, Gate: "file_exists('out/1.md')"},
			},
		}
		data, err := p.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(processo.CaminhoNoCofre(vault, p.Nome), data, 0o644); err != nil {
			t.Fatal(err)
		}

		// v1 (só loop.json)
		if err := os.WriteFile(filepath.Join(loopsDir, "legoloop-beta.loop.json"), []byte(`{"version":"0.5.0"}`), 0o644); err != nil {
			t.Fatal(err)
		}

		// v2 corrompido
		if err := os.WriteFile(processo.CaminhoNoCofre(vault, "legoloop-gama"), []byte(`{quebrado`), 0o644); err != nil {
			t.Fatal(err)
		}

		// esteira não legoloop (deve ser ignorada)
		if err := os.WriteFile(filepath.Join(loopsDir, "outra-esteira.loop.json"), []byte(`{}`), 0o644); err != nil {
			t.Fatal(err)
		}

		// subdiretório (deve ser ignorado)
		if err := os.MkdirAll(filepath.Join(loopsDir, "legoloop-dir.loop.json"), 0o755); err != nil {
			t.Fatal(err)
		}

		procs, err := processo.ProcessosDoCofre(vault)
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if len(procs) != 3 {
			t.Fatalf("esperava 3 processos, veio %d: %+v", len(procs), procs)
		}

		// Ordenação alfabética
		if procs[0].Nome != "legoloop-alfa" || procs[0].Versao != "1.2.0" || procs[0].Natureza != "dados" || procs[0].Origem != "v2" {
			t.Errorf("alfa incorreto: %+v", procs[0])
		}
		if procs[1].Nome != "legoloop-beta" || procs[1].Versao != "0.5.0" || procs[1].Natureza != processo.NaturezaDesconhecida || procs[1].Origem != "v1" {
			t.Errorf("beta incorreto: %+v", procs[1])
		}
		if procs[2].Nome != "legoloop-gama" || procs[2].Erro == "" {
			t.Errorf("gama deveria ter Erro anotado: %+v", procs[2])
		}
	})
}

func TestMapearUsadoPor(t *testing.T) {
	procs := []processo.ProcessoCofre{
		{
			Nome: "legoloop-pai",
			Etapas: []processo.Etapa{
				{ID: "E1", Compoe: "legoloop-filho@1.0.0"},
				{ID: "E2", Compoe: "legoloop-outro@2.0.0"},
				{ID: "E3", Compoe: "legoloop-filho@1.0.0"}, // duplicata ignorada
				{ID: "E4", Compoe: "sem-versao"},           // sem @ ignorada
				{ID: "E5"},                                 // vazia
			},
		},
		{
			Nome: "legoloop-irmao",
			Etapas: []processo.Etapa{
				{ID: "E1", Compoe: "legoloop-filho@1.0.0"},
			},
		},
	}

	usadoPor := processo.MapearUsadoPor(procs)
	filho := usadoPor["legoloop-filho"]
	if !slices.Equal(filho, []string{"legoloop-irmao", "legoloop-pai"}) {
		t.Errorf("esperava ordenado [legoloop-irmao, legoloop-pai], veio %v", filho)
	}

	outro := usadoPor["legoloop-outro"]
	if !slices.Equal(outro, []string{"legoloop-pai"}) {
		t.Errorf("esperava [legoloop-pai], veio %v", outro)
	}

	if len(usadoPor["sem-versao"]) != 0 {
		t.Errorf("etapa sem @ não deveria compor: %v", usadoPor["sem-versao"])
	}
}

func TestHarnessesVinculados(t *testing.T) {
	home := t.TempDir()
	for _, h := range []string{"claude", "codex", "gemini"} {
		if err := os.MkdirAll(filepath.Join(home, "."+h, "skills"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Join(home, ".claude", "skills", "legoloop-x"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(home, ".gemini", "skills", "legoloop-x"), 0o755); err != nil {
		t.Fatal(err)
	}

	hs := processo.DetectarHarnesses(processo.OSFS(), home)
	vinculados := processo.HarnessesVinculados(nil, hs, "legoloop-x")
	if !slices.Equal(vinculados, []string{"claude", "gemini"}) {
		t.Errorf("esperava [claude, gemini], veio %v", vinculados)
	}

	nada := processo.HarnessesVinculados(nil, hs, "legoloop-inexistente")
	if len(nada) != 0 || nada == nil {
		t.Errorf("esperava slice vazio não-nulo, veio %v", nada)
	}
}

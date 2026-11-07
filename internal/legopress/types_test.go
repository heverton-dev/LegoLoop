package legopress_test

import (
	"testing"

	"legoloop/internal/legopress"
)

func TestManifestValidation_Valid(t *testing.T) {
	manifest := legopress.Manifest{
		Title:    "LegoLoop Book",
		Subtitle: "Guia Modular",
		Author:   "Heverton",
		Format:   "170x240mm",
		Parts: []legopress.Part{
			{ID: "P01", Title: "Macro", File: "partes/01-macro.md", Level: 1},
			{ID: "P02", Title: "Meso", File: "partes/02-meso.md", Level: 1},
		},
		OutputPDF: "livro.pdf",
	}

	if err := manifest.Validate(); err != nil {
		t.Fatalf("esperava manifesto valido, obteve erro: %v", err)
	}
}

func TestManifestValidation_EmptyTitle(t *testing.T) {
	manifest := legopress.Manifest{
		Title: "",
		Parts: []legopress.Part{
			{ID: "P01", Title: "Macro", File: "partes/01-macro.md"},
		},
	}

	err := manifest.Validate()
	if err == nil {
		t.Fatal("esperava erro para titulo vazio, obteve nil")
	}
	if err != legopress.ErrManifestoInvalido {
		t.Fatalf("esperava ErrManifestoInvalido, obteve: %v", err)
	}
}

func TestManifestValidation_NoParts(t *testing.T) {
	manifest := legopress.Manifest{
		Title: "Livro Sem Partes",
		Parts: []legopress.Part{},
	}

	err := manifest.Validate()
	if err != legopress.ErrManifestoInvalido {
		t.Fatalf("esperava ErrManifestoInvalido para partes vazias, obteve: %v", err)
	}
}

func TestManifestValidation_EmptyPartFile(t *testing.T) {
	manifest := legopress.Manifest{
		Title: "Livro Incompleto",
		Parts: []legopress.Part{
			{ID: "P01", Title: "Macro", File: "   "},
		},
	}

	err := manifest.Validate()
	if err != legopress.ErrManifestoInvalido {
		t.Fatalf("esperava ErrManifestoInvalido para caminho de arquivo vazio, obteve: %v", err)
	}
}

func TestModeValidation(t *testing.T) {
	validModes := []legopress.Mode{
		legopress.ModeCreate,
		legopress.ModeEdit,
		legopress.ModeSync,
		legopress.ModeAdapt,
		legopress.ModeAudit,
	}

	for _, m := range validModes {
		if err := m.Validate(); err != nil {
			t.Errorf("modo valido %q reportou erro: %v", m, err)
		}
	}

	invalidMode := legopress.Mode("invalido")
	if err := invalidMode.Validate(); err != legopress.ErrModoDesconhecido {
		t.Fatalf("esperava ErrModoDesconhecido para modo invalido, obteve: %v", err)
	}
}

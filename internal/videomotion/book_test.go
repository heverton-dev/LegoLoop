package videomotion_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/videomotion"
)

func TestGenerateDidacticBook(t *testing.T) {
	tempDir := t.TempDir()

	cfg := videomotion.DidacticBookConfig{
		OutputDir:    tempDir,
		ProcessoName: "video-motion",
		Version:      "1.0.0",
	}

	manifestPath, err := videomotion.GenerateDidacticBook(cfg)
	if err != nil {
		t.Fatalf("GenerateDidacticBook falhou: %v", err)
	}

	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		t.Fatalf("manifesto livro.json não existe no caminho: %s", manifestPath)
	}

	// Validar integridade do manifesto LegoPress
	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("falha ao ler manifesto: %v", err)
	}
	var manifest videomotion.LegoPressBookManifest
	if err := json.Unmarshal(raw, &manifest); err != nil {
		t.Fatalf("manifesto inválido no formato LegoPress: %v", err)
	}
	if manifest.VersaoManifesto != "1.0" || manifest.PastaPartes != "partes" {
		t.Errorf("manifesto com campos divergentes do padrão LegoPress: %+v", manifest)
	}

	// Validar que os capítulos foram criados dentro da subpasta partes/
	bookDir := filepath.Dir(manifestPath)
	partesDir := filepath.Join(bookDir, "partes")
	for _, capFile := range manifest.Partes {
		capPath := filepath.Join(partesDir, capFile)
		content, err := os.ReadFile(capPath)
		if err != nil {
			t.Errorf("capítulo %s não foi criado em partes/: %v", capFile, err)
			continue
		}
		if len(content) < 50 {
			t.Errorf("capítulo %s muito curto (< 50 bytes)", capFile)
		}
	}
}

func TestGenerateDidacticBook_Validation(t *testing.T) {
	_, err := videomotion.GenerateDidacticBook(videomotion.DidacticBookConfig{})
	if err != videomotion.ErrInvalidBookConfig {
		t.Errorf("esperava ErrInvalidBookConfig, obteve: %v", err)
	}
}

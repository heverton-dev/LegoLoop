package videomotion_test

import (
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/videomotion"
)

func TestGenerateDidacticBook(t *testing.T) {
	tempDir := t.TempDir()

	cfg := videomotion.DidacticBookConfig{
		OutputDir:    tempDir,
		PipelineName: "video-motion",
		Version:      "1.0.0",
	}

	manifestPath, err := videomotion.GenerateDidacticBook(cfg)
	if err != nil {
		t.Fatalf("GenerateDidacticBook falhou: %v", err)
	}

	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		t.Fatalf("manifesto livro.json não existe no caminho: %s", manifestPath)
	}

	// Validar que os 5 capítulos do livro foram criados
	expectedChapters := []string{
		"capitulo-01-fundacao.md",
		"capitulo-02-arquitetura-4-pilares.md",
		"capitulo-03-passo-a-passo.md",
		"capitulo-04-receitas-multi-aspecto.md",
		"capitulo-05-rotina-diff-updater.md",
	}

	bookDir := filepath.Dir(manifestPath)
	for _, capFile := range expectedChapters {
		capPath := filepath.Join(bookDir, capFile)
		content, err := os.ReadFile(capPath)
		if err != nil {
			t.Errorf("capítulo %s não foi criado: %v", capFile, err)
			continue
		}
		if len(content) < 100 {
			t.Errorf("capítulo %s muito curto (< 100 bytes)", capFile)
		}
	}
}

func TestGenerateDidacticBook_Validation(t *testing.T) {
	_, err := videomotion.GenerateDidacticBook(videomotion.DidacticBookConfig{})
	if err == nil {
		t.Fatal("esperava erro para configuração vazia de livro")
	}
}

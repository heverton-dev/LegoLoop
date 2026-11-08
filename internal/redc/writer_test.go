package redc_test

import (
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/redc"
)

func TestWriteREPC_Success(t *testing.T) {
	tempDir := t.TempDir()
	dna := redc.ChannelDNA{
		CanalSlug:    "canal-teste",
		TomDeVoz:     "Autoritário e dinâmico",
		Arquetipo:    "Mestre",
		RitmoCortes:  "Cortes a cada 2s",
		DesignTokens: map[string]string{"cor_primaria": "#123456"},
		MotionTokens: map[string]string{"transicao": "snap"},
		TopGanchos:   []string{"Gancho de ouro"},
	}

	created, err := redc.WriteREPC(tempDir, dna)
	if err != nil {
		t.Fatalf("esperava sucesso ao gravar repc, obteve: %v", err)
	}

	if len(created) != 4 {
		t.Errorf("esperava 4 arquivos criados, obteve %d", len(created))
	}

	// Verificar se os 4 arquivos existem
	expectedFiles := []string{
		filepath.Join(tempDir, "docs", "channels", "canal-teste", "repc", "01_dossie_dna.json"),
		filepath.Join(tempDir, "docs", "channels", "canal-teste", "repc", "02_design_tokens.json"),
		filepath.Join(tempDir, "docs", "channels", "canal-teste", "repc", "03_motion_tokens.json"),
		filepath.Join(tempDir, "docs", "channels", "canal-teste", "repc", "04_matriz_ganchos.md"),
	}

	for _, f := range expectedFiles {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("arquivo esperado nao foi criado: %s", f)
		}
	}
}

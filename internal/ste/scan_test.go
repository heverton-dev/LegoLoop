package ste_test

import (
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/ste"
)

func TestScanHarnesses(t *testing.T) {
	t.Parallel()

	tmpHome := t.TempDir()

	// Simula apenas Gemini e Claude instalados
	geminiDir := filepath.Join(tmpHome, ".gemini")
	claudeDir := filepath.Join(tmpHome, ".claude")
	if err := os.MkdirAll(geminiDir, 0755); err != nil {
		t.Fatalf("falha ao criar gemini dir: %v", err)
	}
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatalf("falha ao criar claude dir: %v", err)
	}

	harnesses := ste.ScanHarnesses(tmpHome)
	if len(harnesses) == 0 {
		t.Fatalf("esperava lista de harnesses nao vazia")
	}

	mapa := make(map[ste.HarnessID]ste.HarnessInfo)
	for _, h := range harnesses {
		mapa[h.ID] = h
	}

	gemini, okGemini := mapa[ste.HarnessGemini]
	if !okGemini || !gemini.Instalado {
		t.Errorf("esperava gemini instalado=true, obteve %+v", gemini)
	}

	claude, okClaude := mapa[ste.HarnessClaude]
	if !okClaude || !claude.Instalado {
		t.Errorf("esperava claude instalado=true, obteve %+v", claude)
	}

	omp, okOmp := mapa[ste.HarnessOMP]
	if !okOmp || omp.Instalado {
		t.Errorf("esperava omp instalado=false, obteve %+v", omp)
	}

	// Verifica se NÃO criou diretórios vazios para harnesses ausentes
	ompDir := filepath.Join(tmpHome, ".omp")
	if _, err := os.Stat(ompDir); !os.IsNotExist(err) {
		t.Errorf("scanner nao pode criar diretorios fantasmas no disco: %s", ompDir)
	}
}

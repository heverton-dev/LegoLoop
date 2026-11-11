package caminhos

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLeituraDoisDialetos(t *testing.T) {
	base := t.TempDir()
	if got := DirsLeitura(base); len(got) != 0 {
		t.Fatalf("base vazia: %v", got)
	}
	antigo := filepath.Join(base, DirLoops)
	novo := filepath.Join(base, DirProcessos)
	for _, d := range []string{antigo, novo} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	got := DirsLeitura(base)
	if len(got) != 2 || got[0] != novo || got[1] != antigo {
		t.Fatalf("ordem esperada novo,antigo: %v", got)
	}
	if got := CaminhoDAGLeitura(antigo); got != filepath.Join(antigo, ArquivoDAG) {
		t.Errorf("sem arquivo novo: %q", got)
	}
	if err := os.WriteFile(filepath.Join(antigo, ArquivoProcessoDAG), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := CaminhoDAGLeitura(antigo); got != filepath.Join(antigo, ArquivoProcessoDAG) {
		t.Errorf("com arquivo novo: %q", got)
	}
}

package caminhos

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolucao(t *testing.T) {
	if got, want := DirLoopsEm("v"), filepath.Join("v", "loops"); got != want {
		t.Errorf("DirLoopsEm = %q, esperado %q", got, want)
	}
	if got, want := CaminhoDAG("d"), filepath.Join("d", "pipeline.dag.json"); got != want {
		t.Errorf("CaminhoDAG = %q, esperado %q", got, want)
	}
	if ChaveAPILoops != "loops" {
		t.Errorf("ChaveAPILoops mudou: %q", ChaveAPILoops)
	}
}

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

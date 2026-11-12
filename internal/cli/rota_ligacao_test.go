package cli

import (
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/orchestrix"
)

func TestLigarRoteamentoLeEstadoDaHomologacao(t *testing.T) {
	home := t.TempDir()
	if err := os.MkdirAll(filepath.Join(home, ".aidd"), 0o755); err != nil {
		t.Fatal(err)
	}
	dados := `{"omp|oficial":{"estado":"homologado"},"agy|oficial":{"estado":"reprovado"}}`
	if err := os.WriteFile(filepath.Join(home, ".aidd", "homologacao.json"), []byte(dados), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(t.TempDir())
	ox := &orchestrix.Orchestrix{}
	if err := ligarRoteamento(ox, Deps{Home: home}); err != nil {
		t.Fatal(err)
	}
	if ox.EstadoHarness("omp", "oficial") != "homologado" || ox.EstadoHarness("agy", "oficial") != "reprovado" || ox.EstadoHarness("omp", "9router") != "" {
		t.Fatal("estado de homologação não chegou ao motor")
	}
	if ox.Elegiveis == nil {
		t.Fatal("candidatos elegíveis não chegaram ao motor")
	}
}

func TestRunnerOptionsLigaWorkersEPeso(t *testing.T) {
	ox := &orchestrix.Orchestrix{MemoriaLivre: func() (uint64, error) { return 4 << 30, nil }}
	opts := runFlags{parallel: 8}.runnerOptions(ox, nil)
	if opts.Workers == nil || opts.PesoMB == nil {
		t.Fatal("o runner não recebeu o limite de workers e o peso")
	}
	if got := opts.Workers([]int{750, 750}); got != 2 {
		t.Fatalf("workers = %d", got)
	}
}

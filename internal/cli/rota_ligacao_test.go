package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/caminhos"
	"legoloop/internal/orchestrix"
	"legoloop/internal/roteamento/homologacao"
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

func TestRunGravaTabelaFactualComReport(t *testing.T) {
	work := t.TempDir()
	p := caminhos.CaminhoDAG(work)
	if err := os.WriteFile(p, []byte(processoScript), 0o644); err != nil {
		t.Fatal(err)
	}
	saida := filepath.Join(work, "REPORT_FACTUAL.md")
	d, _, errb := deps(t)
	if code := run(t, d, "run", p, "--workdir", work, "--allow-scripts", "--report", saida); code != 0 {
		t.Fatalf("code=%d err=%s", code, errb.String())
	}
	b, err := os.ReadFile(saida)
	if err != nil || !strings.Contains(string(b), "| Nó | Commit | Custo (USD) | Gate | Harness | Via | Modelo | Substituição |") || !strings.Contains(string(b), "| a | - |") {
		t.Fatalf("tabela factual = %q (%v)", b, err)
	}
}

func TestAlarmeDeBanimentoPersisteOParNaHomologacao(t *testing.T) {
	home := t.TempDir()
	t.Chdir(t.TempDir())
	ox := &orchestrix.Orchestrix{}
	if err := ligarRoteamento(ox, Deps{Home: home}); err != nil {
		t.Fatal(err)
	}
	if ox.Reprovados == nil || ox.AoReprovar == nil {
		t.Fatal("o alarme não chegou ao motor")
	}
	ox.AoReprovar("agy", "oficial", "banimento: terms of service")
	store, err := homologacao.Ler(homologacao.Caminho(home))
	if err != nil || store["agy|oficial"].Estado != homologacao.EstadoReprovado {
		t.Fatalf("store = %+v (%v)", store, err)
	}
}

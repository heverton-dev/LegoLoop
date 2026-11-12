package orchestrix

import (
	"errors"
	"testing"

	"legoloop/internal/dag"
)

const gb = uint64(1) << 30

func TestWorkersPorMemoria(t *testing.T) {
	casos := []struct {
		nome     string
		parallel int
		maximo   int
		livre    uint64
		pesos    []int
		quer     int
	}{
		{"memória larga: vale o --parallel", 4, 0, 64 * gb, []int{750, 750, 750, 750}, 4},
		{"memória limita: 4 GB x 0,8 / 1500 MB = 2", 4, 0, 4 * gb, []int{750, 750}, 2},
		{"arredonda para baixo: 4 GB x 0,8 / 1000 MB = 3", 8, 0, 4 * gb, []int{500, 500}, 3},
		{"nunca menos que 1", 4, 0, 100 << 20, []int{750}, 1},
		{"teto próprio do perfil", 8, 2, 64 * gb, []int{750}, 2},
		{"nunca passa do --parallel", 3, 10, 64 * gb, []int{750}, 3},
		{"sem pesos, só o teto", 4, 0, 1 * gb, nil, 4},
	}
	for _, c := range casos {
		if got := CalcularWorkers(c.parallel, c.maximo, c.livre, c.pesos); got != c.quer {
			t.Errorf("%s: workers = %d, esperado %d", c.nome, got, c.quer)
		}
	}
}

func TestPesoDoHarnessSemMedicaoVale750(t *testing.T) {
	o := &Orchestrix{PesoHarnessMB: func(h, via string) int {
		if h == "omp" {
			return 420
		}
		return 0
	}}
	omp := dag.Node{Routing: dag.RoutingConfig{Harness: "omp", Via: "oficial"}}
	agy := dag.Node{Routing: dag.RoutingConfig{Harness: "agy", Via: "oficial"}}
	if o.PesoMB(omp) != 420 || o.PesoMB(agy) != 750 || (&Orchestrix{}).PesoMB(omp) != 750 {
		t.Fatalf("pesos: omp=%d agy=%d", o.PesoMB(omp), o.PesoMB(agy))
	}
}

func TestLimiteWorkersUsaAMemoriaInjetada(t *testing.T) {
	o := &Orchestrix{MemoriaLivre: func() (uint64, error) { return 4 * gb, nil }}
	if got := o.LimiteWorkers(8, 0)([]int{750, 750}); got != 2 {
		t.Fatalf("workers = %d", got)
	}
	o.MemoriaLivre = func() (uint64, error) { return 0, errors.New("sem leitura") }
	if got := o.LimiteWorkers(8, 0)([]int{750, 750}); got != 8 {
		t.Fatalf("sem leitura de memória deveria valer o teto: %d", got)
	}
}

func TestMemoriaLivreDoSOLeUmValorReal(t *testing.T) {
	livre, err := MemoriaLivreSO()
	if err != nil || livre == 0 {
		t.Fatalf("memória livre = %d, err = %v", livre, err)
	}
}

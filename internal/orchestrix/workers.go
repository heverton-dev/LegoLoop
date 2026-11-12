package orchestrix

import (
	"legoloop/internal/dag"
)

// pesoPadraoMB é o peso assumido de um harness sem peso_mb homologado (SPEC 4.11).
const pesoPadraoMB = 750

// fracaoMemoria é a parte da memória livre que os workers podem ocupar (0,8).
const (
	fracaoNum = 8
	fracaoDen = 10
)

// CalcularWorkers aplica a fórmula da SPEC 4.11:
// min(parallel, maxWorkers, floor(memLivre x 0,8 / soma dos pesos)), com mínimo de 1.
// maxWorkers <= 0 significa sem teto próprio. memLivre vem em bytes e os pesos em MB.
func CalcularWorkers(parallel, maxWorkers int, memLivre uint64, pesosMB []int) int {
	teto := max(parallel, 1)
	if maxWorkers > 0 {
		teto = min(teto, maxWorkers)
	}
	var soma uint64
	for _, p := range pesosMB {
		soma += uint64(max(p, 0)) << 20
	}
	if soma == 0 {
		return teto
	}
	porMemoria := memLivre / fracaoDen * fracaoNum / soma
	return max(1, int(min(uint64(teto), porMemoria)))
}

// PesoMB devolve o peso do nó: o peso homologado do harness ou 750 MB sem medição.
func (o *Orchestrix) PesoMB(n dag.Node) int {
	if o.PesoHarnessMB != nil {
		if p := o.PesoHarnessMB(harnessDe(n), n.Routing.Via); p > 0 {
			return p
		}
	}
	return pesoPadraoMB
}

// LimiteWorkers devolve a função que o runner consulta antes de iniciar um nó.
// Sem leitura de memória (nil ou erro), só vale o teto de parallel e de maxWorkers.
func (o *Orchestrix) LimiteWorkers(parallel, maxWorkers int) func([]int) int {
	return func(pesosMB []int) int {
		if o.MemoriaLivre == nil {
			return CalcularWorkers(parallel, maxWorkers, ^uint64(0), nil)
		}
		livre, err := o.MemoriaLivre()
		if err != nil {
			return CalcularWorkers(parallel, maxWorkers, ^uint64(0), nil)
		}
		return CalcularWorkers(parallel, maxWorkers, livre, pesosMB)
	}
}

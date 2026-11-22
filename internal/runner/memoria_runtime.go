package runner

import (
	"log"
	"os"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/dag"
	"github.com/heverton-dev/LegoLoop/internal/triagem"
)

// proximosDe devolve os IDs dos nós que dependem de id, pela ordem do DAG.
func proximosDe(d *dag.DAG, id string) []string {
	var out []string
	for _, n := range d.Nodes {
		for _, dep := range n.DependsOn {
			if dep == id {
				out = append(out, n.ID)
				break
			}
		}
	}
	return out
}

// recordEtapaMemoria grava a etapa atual e o próximo passo na memória de
// execução do processo (controle de progresso da Lei 11).
func recordEtapaMemoria(memDir, sessaoID, etapa, proximo string) {
	if memDir == "" || sessaoID == "" {
		return
	}
	mem := triagem.MemoriaExecucao{
		SessaoID:     sessaoID,
		PID:          os.Getpid(),
		EtapaAtual:   etapa,
		ProximoPasso: proximo,
	}
	if anterior, err := triagem.LerMemoria(memDir, sessaoID); err == nil {
		mem.TokenBudget = anterior.TokenBudget
	}
	if err := triagem.SalvarMemoria(memDir, mem); err != nil {
		log.Printf("memoria: falha ao gravar etapa %s: %v", etapa, err)
	}
}

// registrarEtapaMemoria marca o nó atual na memória e aponta o próximo passo.
func (r *run) registrarEtapaMemoria(id string) {
	proximo := strings.Join(proximosDe(r.d, id), ",")
	if proximo == "" {
		proximo = "fim"
	}
	recordEtapaMemoria(r.opts.Workdir, r.opts.SessaoID, id, proximo)
}

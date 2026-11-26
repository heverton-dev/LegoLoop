package cli

import (
	"os"

	"github.com/heverton-dev/LegoLoop/internal/progresso"
)

// preservarEstadoAnterior copia estado, prova e responsavel dos tickets ja executados do progresso
// existente para o recem-gerado. Regerar o plano nao pode apagar o historico de aprovacao.
func preservarEstadoAnterior(progJSON string, novo *progresso.Progresso) {
	if _, err := os.Stat(progJSON); err != nil {
		return
	}
	antigo, err := progresso.Carregar(progJSON)
	if err != nil {
		return
	}
	por := make(map[string]progresso.ItemTicket, len(antigo.Tickets))
	for _, it := range antigo.Tickets {
		por[it.ID] = it
	}
	for i := range novo.Tickets {
		a, ok := por[novo.Tickets[i].ID]
		if !ok || a.Estado == progresso.EstadoPendente {
			continue
		}
		novo.Tickets[i].Estado, novo.Tickets[i].Responsavel = a.Estado, a.Responsavel
		novo.Tickets[i].ExitCode, novo.Tickets[i].Prova = a.ExitCode, a.Prova
		novo.Tickets[i].AtualizadoEm = a.AtualizadoEm
	}
	novo.RecalcularResumo()
}

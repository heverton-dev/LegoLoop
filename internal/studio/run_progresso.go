package studio

import (
	"context"
	"net/http"
	"os"

	"github.com/heverton-dev/LegoLoop/internal/dag"
	"github.com/heverton-dev/LegoLoop/internal/progresso"
	"github.com/heverton-dev/LegoLoop/internal/runner"
)

// arquivoProgresso resolve e confere o campo `progresso` do DAG; vazio desliga o progresso vivo.
func (s *Server) arquivoProgresso(d *dag.DAG) (string, error) {
	arq, err := progresso.CaminhoDoDAG(s.ProjectDir, d.Progresso)
	if err != nil || arq == `` {
		return ``, err
	}
	if _, err := os.Stat(arq); err != nil {
		return ``, err
	}
	return arq, nil
}

// validarProgresso responde 400 quando o campo `progresso` do DAG é inválido ou aponta para arquivo ausente.
func (s *Server) validarProgresso(w http.ResponseWriter, d *dag.DAG) bool {
	if _, err := s.arquivoProgresso(d); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{`error`: `progresso vivo: ` + s.maskPaths(err.Error())})
		return false
	}
	return true
}

// ligarProgresso encadeia o Observador aos ganchos do runner; devolve nil quando o DAG não declara progresso.
func (s *Server) ligarProgresso(ctx context.Context, d *dag.DAG, opts *runner.Options) *progresso.Observador {
	arq, err := s.arquivoProgresso(d)
	if err != nil || arq == `` {
		return nil
	}
	obs := progresso.NovoObservador(ctx, arq, s.ProjectDir, `studio`)
	inicio, fim := opts.OnNodeStart, opts.OnNodeFinish
	opts.OnNodeStart = func(id string) {
		obs.NoInicio(id)
		inicio(id)
	}
	opts.OnNodeFinish = func(nr runner.NodeReport) {
		obs.NoFim(nr.ID, nr.Status == runner.StatusPassed, nr.ExitCode, nr.Error)
		fim(nr)
	}
	return obs
}

// avisarErrosProgresso transmite as falhas de gravação do progresso no canal SSE.
func (s *Server) avisarErrosProgresso(obs *progresso.Observador, send sseSender) {
	if obs == nil {
		return
	}
	for _, e := range obs.Erros() {
		send(`progresso_error`, map[string]any{`error`: s.maskPaths(e.Error())})
	}
}

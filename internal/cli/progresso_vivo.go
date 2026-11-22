package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"legoloop/internal/board"
	"legoloop/internal/dag"
	"legoloop/internal/gate"
	"legoloop/internal/orchestrix"
	"legoloop/internal/progresso"
	"legoloop/internal/runner"
)

// caminhoProgressoVivo escolhe o arquivo de progresso: a flag vence o campo `progresso` do DAG.
// A flag é do usuário e vale como veio; o campo do DAG passa pelo Path Jail (progresso.CaminhoDoDAG).
// Vazio desliga o progresso vivo.
func caminhoProgressoVivo(base, flag, doDAG string) (string, error) {
	if flag == `` {
		return progresso.CaminhoDoDAG(base, doDAG)
	}
	if filepath.IsAbs(flag) {
		return flag, nil
	}
	return filepath.Join(base, flag), nil
}

func novoObservador(ctx context.Context, arquivo, dir, responsavel string) (*progresso.Observador, error) {
	if arquivo == `` {
		return nil, nil
	}
	if _, err := os.Stat(arquivo); err != nil {
		return nil, fmt.Errorf(`progresso vivo: %w`, err)
	}
	return progresso.NovoObservador(ctx, arquivo, dir, responsavel), nil
}

// runVivo roda o DAG no runner e mantém o progresso vivo a cada nó.
func runVivo(ctx context.Context, p *dag.DAG, f runFlags, ox *orchestrix.Orchestrix, pr gate.Prompter) (runner.Report, *progresso.Observador, error) {
	opts := f.runnerOptions(ox, pr)
	arq, err := caminhoProgressoVivo(f.workdir, f.progresso, p.Progresso)
	if err != nil {
		return runner.Report{}, nil, err
	}
	obs, err := novoObservador(ctx, arq, f.workdir, `runner`)
	if err != nil {
		return runner.Report{}, nil, err
	}
	if obs != nil {
		opts.OnNodeStart = obs.NoInicio
		opts.OnNodeFinish = func(nr runner.NodeReport) {
			obs.NoFim(nr.ID, nr.Status == runner.StatusPassed, nr.ExitCode, nr.Error)
		}
	}
	rep, err := runner.Run(ctx, p, opts)
	return rep, obs, err
}

// ligarProgressoSwarm conecta o swarm ao progresso vivo; devolve nil quando não há arquivo.
func ligarProgressoSwarm(ctx context.Context, ro *board.RunOptions, arquivo, dir string) (*progresso.Observador, error) {
	obs, err := novoObservador(ctx, arquivo, dir, `swarm`)
	if obs != nil {
		ro.OnNodeStart = obs.NoInicio
		ro.OnNodeEnd = obs.NoFim
	}
	return obs, err
}

// codigoComProgresso imprime as falhas de gravação do progresso e reprova a execução se houve alguma.
func codigoComProgresso(d Deps, obs *progresso.Observador, codigo int) int {
	if obs == nil {
		return codigo
	}
	erros := obs.Erros()
	for _, e := range erros {
		fmt.Fprintln(d.Stderr, `erro: progresso não atualizado:`, e)
	}
	if codigo == 0 && len(erros) > 0 {
		return 1
	}
	return codigo
}

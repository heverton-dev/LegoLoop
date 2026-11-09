package runner

import (
	"context"
	"fmt"

	"legoloop/internal/dag"
)

// MaxNodeRetries define o teto do Poka-Yoke de runtime: no máximo 2 tentativas.
const MaxNodeRetries = 2

// ExecuteWithAndon executa um nó com disjuntor Jidoka (Andon Cord).
// Se a 2ª tentativa falhar, puxa o Andon Cord e gera um diagnóstico SBAR.
func ExecuteWithAndon(ctx context.Context, exec Executor, n dag.Node, workdir string) (Result, SBAR, error) {
	var lastRes Result
	var lastErr error

	for attempt := 1; attempt <= MaxNodeRetries; attempt++ {
		res, err := exec.Execute(ctx, n, workdir)
		lastRes = res
		lastErr = err

		if err == nil && res.ExitCode == 0 {
			return res, SBAR{}, nil
		}
	}

	sbar := generateAndonSBAR(n.ID, lastRes.ExitCode, lastErr)
	andonErr := fmt.Errorf("andon_cord: limite de %d tentativas esgotado para nó %q: %w", MaxNodeRetries, n.ID, lastErr)
	return lastRes, sbar, andonErr
}

func generateAndonSBAR(nodeID string, exitCode int, err error) SBAR {
	errStr := "desconhecido"
	if err != nil {
		errStr = err.Error()
	}
	return FormatNodeError(nodeID, exitCode, errStr)
}

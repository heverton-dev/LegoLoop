package runner_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/runner"
)

type failTwiceExecutor struct {
	calls int
}

func (e *failTwiceExecutor) Execute(_ context.Context, _ dag.Node, _ string) (runner.Result, error) {
	e.calls++
	return runner.Result{
		ExitCode: 1,
		Notes:    []string{"falha sintática simulada"},
	}, errors.New("simulated test failure")
}

func TestAndonCordCircuitBreaker(t *testing.T) {
	exec := &failTwiceExecutor{}
	node := dag.Node{
		ID:     "step_andon",
		Type:   dag.TypeUnitTask,
		Output: dag.HandoffContract{Artifact: "out.txt", Gate: "exit_code == 0"},
	}

	res, sbar, err := runner.ExecuteWithAndon(context.Background(), exec, node, t.TempDir())
	if err == nil {
		t.Fatal("esperava erro ao falhar 2 vezes, mas obteve nil")
	}

	if exec.calls != 2 {
		t.Fatalf("esperava exatamente 2 tentativas (Andon Cord), obteve %d", exec.calls)
	}

	if res.ExitCode != 1 {
		t.Errorf("esperava exit_code 1, obteve %d", res.ExitCode)
	}

	if !strings.Contains(err.Error(), "andon_cord") {
		t.Errorf("erro deve mencionar andon_cord, obteve: %v", err)
	}

	if sbar.Situation == "" || sbar.Recommendation == "" {
		t.Errorf("esperava SBAR preenchido, obteve: %+v", sbar)
	}
}

func TestSBAREmissionOnFailure(t *testing.T) {
	sbar := runner.NewSBAR(
		"Falha no nó step_crawler",
		"Tentativa 2/2 falhou ao executar go test",
		"Exit code 1 com falha de asserção na Seam",
		"Verificar assinatura de função ou corrigir teste antes de nova tentativa",
	)

	formatted := sbar.Format()
	requiredMarkers := []string{
		"[S] SITUATION:",
		"[B] BACKGROUND:",
		"[A] ASSESSMENT:",
		"[R] RECOMMENDATION:",
	}

	for _, marker := range requiredMarkers {
		if !strings.Contains(formatted, marker) {
			t.Errorf("formato SBAR deve conter %q, obteve:\n%s", marker, formatted)
		}
	}
}

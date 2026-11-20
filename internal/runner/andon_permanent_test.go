package runner_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/runner"
)

type permanentExecutor struct{ calls int }

func (e *permanentExecutor) Execute(_ context.Context, _ dag.Node, _ string) (runner.Result, error) {
	e.calls++
	return runner.Result{ExitCode: 1}, fmt.Errorf("nó: %w", runner.Permanent(errors.New("gate reprovou")))
}

func TestAndonDoesNotRetryPermanentFailure(t *testing.T) {
	exec := &permanentExecutor{}
	node := dag.Node{ID: "n", Type: dag.NodeType("unit_no"), Output: dag.HandoffContract{Artifact: "o.txt", Gate: "exit_code == 0"}}
	_, _, err := runner.ExecuteWithAndon(context.Background(), exec, node, t.TempDir())
	if err == nil {
		t.Fatal("esperava erro")
	}
	if exec.calls != 1 {
		t.Fatalf("falha definitiva deve ter 1 tentativa, obteve %d", exec.calls)
	}
	if !runner.IsPermanent(err) {
		t.Errorf("o erro do Andon deve continuar definitivo: %v", err)
	}
	if strings.Contains(err.Error(), "limite de") {
		t.Errorf("mensagem não deve falar de tentativas esgotadas: %v", err)
	}
}

func TestAndonStillRetriesOrdinaryFailure(t *testing.T) {
	exec := &failTwiceExecutor{}
	node := dag.Node{ID: "n", Type: dag.TypeScript, Output: dag.HandoffContract{Artifact: "o.txt", Gate: "exit_code == 0"}}
	_, _, err := runner.ExecuteWithAndon(context.Background(), exec, node, t.TempDir())
	if err == nil || exec.calls != runner.MaxNodeRetries {
		t.Fatalf("falha comum deve usar %d tentativas, obteve %d (err=%v)", runner.MaxNodeRetries, exec.calls, err)
	}
	if runner.IsPermanent(err) {
		t.Errorf("falha comum não é definitiva: %v", err)
	}
}

func TestPermanentKeepsCause(t *testing.T) {
	cause := errors.New("causa")
	err := runner.Permanent(cause)
	if !errors.Is(err, cause) || !runner.IsPermanent(err) {
		t.Fatalf("Permanent deve preservar a causa e marcar: %v", err)
	}
	if runner.Permanent(nil) != nil {
		t.Error("Permanent(nil) deve devolver nil")
	}
}

package orchestrix

import (
	"context"
	"strings"
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/runner"
)

type countingExecutor struct {
	inner *Orchestrix
	calls int
}

func (c *countingExecutor) Execute(ctx context.Context, n dag.Node, workdir string) (runner.Result, error) {
	c.calls++
	return c.inner.Execute(ctx, n, workdir)
}

// A reprovação do gate depois das retentativas é definitiva: o runner não repete o nó.
func TestUnitNoGateEsgotadoNaoEhRepetidoPeloRunner(t *testing.T) {
	repo := newRepoT(t)
	n := unitNode("n1", "out/n1.json", "file_exists('nunca.txt')", nil)
	exec := &countingExecutor{inner: unitOrchestrix("n1.txt=x")}
	_, _, err := runner.ExecuteWithAndon(context.Background(), exec, n, repo)
	if err == nil || !strings.Contains(err.Error(), "reprovou") {
		t.Fatalf("err = %v", err)
	}
	if exec.calls != 1 {
		t.Fatalf("o runner deve chamar o nó 1 vez, chamou %d", exec.calls)
	}
	if !runner.IsPermanent(err) {
		t.Errorf("a reprovação do gate deve ser definitiva: %v", err)
	}
}

func TestUnitNoRespeitaRetryPolicyMaxRetries(t *testing.T) {
	cases := []struct{ max, want int }{{0, 1}, {1, 2}, {2, 3}}
	for _, c := range cases {
		repo := newRepoT(t)
		n := unitNode("n1", "out/n1.json", "file_exists('nunca.txt')", nil)
		n.RetryPolicy = &dag.RetryPolicy{MaxRetries: c.max}
		res, err := unitOrchestrix("n1.txt=x").Execute(context.Background(), n, repo)
		if err == nil {
			t.Fatalf("max_retries=%d: esperava reprovação", c.max)
		}
		if want := 5 * c.want; res.Tokens.Input != want {
			t.Errorf("max_retries=%d: tokens = %d, esperado %d (%d chamadas)", c.max, res.Tokens.Input, want, c.want)
		}
	}
}

func TestUnitRetriesCapsAtMaxUnitRetries(t *testing.T) {
	n := dag.Node{RetryPolicy: &dag.RetryPolicy{MaxRetries: 9}}
	if got := unitRetries(n); got != maxUnitRetries {
		t.Errorf("unitRetries = %d, esperado o teto %d", got, maxUnitRetries)
	}
	if got := unitRetries(dag.Node{}); got != maxUnitRetries {
		t.Errorf("sem retry_policy deve usar %d, obteve %d", maxUnitRetries, got)
	}
}

package cli_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"legoloop/internal/cli"
)

func TestCLICommandsConsolidated(t *testing.T) {
	ctx := context.Background()

	t.Run("UsoExibeMatrizUnificada", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		code := cli.Run(ctx, []string{"--help"}, cli.Deps{Stdout: &stdout, Stderr: &stderr})
		if code != 0 {
			t.Fatalf("esperava exit code 0 para --help, obteve %d", code)
		}
		out := stdout.String()
		if !strings.Contains(out, "matriz") {
			t.Errorf("ajuda deve conter 'matriz', obteve:\n%s", out)
		}
		if !strings.Contains(out, "DESCONTINUADO: use \"matriz ver\"") {
			t.Errorf("ajuda deve sinalizar que matrix está descontinuado em prol de matriz ver")
		}
	})

	t.Run("ComandoMatrizStatusRetornaZeroOuAjuda", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		code := cli.Run(ctx, []string{"matriz", "--help"}, cli.Deps{Stdout: &stdout, Stderr: &stderr})
		if code != 0 && code != 2 {
			t.Fatalf("esperava 0 ou 2 para matriz --help, obteve %d", code)
		}
	})
}

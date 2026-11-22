package cli_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/cli"
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
		if strings.Contains(out, "matrix") {
			t.Errorf("ajuda não deve conter o verbo removido 'matrix', obteve:\n%s", out)
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

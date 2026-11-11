package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"legoloop/internal/cli"
)

func TestICAOReadbackPrompter(t *testing.T) {
	conf := cli.ICAOConfirmacao{
		TargetResource:  "internal/videorea/crawler.go",
		MutationCommand: "rm internal/videorea/crawler.go",
		BlastRadius:     "1 arquivo de produção",
		RollbackCommand: "git checkout -- internal/videorea/crawler.go",
	}

	t.Run("AprovaComHearbackApproved", func(t *testing.T) {
		in := strings.NewReader("HEARBACK APPROVED\n")
		var out bytes.Buffer

		approved, err := cli.PromptICAO(&out, in, conf)
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if !approved {
			t.Fatal("esperava approved == true")
		}

		rendered := out.String()
		if !strings.Contains(rendered, "TARGET RESOURCE:") || !strings.Contains(rendered, "ROLLBACK COMMAND:") {
			t.Errorf("saída deve conter campos ICAO, obteve:\n%s", rendered)
		}
	})

	t.Run("RejeitaComEntradaDivergente", func(t *testing.T) {
		in := strings.NewReader("qualquer outra coisa\n")
		var out bytes.Buffer

		approved, err := cli.PromptICAO(&out, in, conf)
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if approved {
			t.Fatal("esperava approved == false para entrada não autorizada")
		}
	})
}

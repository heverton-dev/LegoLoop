package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"legoloop/internal/cli"
)

func TestAutopilotUIMode(t *testing.T) {
	decision := cli.AutopilotDecisionCard{
		Title:         "Aprovação da Curadoria de Peças",
		Summary:       "O pipeline identificou 3 ferramentas gratuitas que substituem APIs pagas.",
		EstimatedTime: "1m 30s",
		EstimatedCost: "$0.00 (100% Soberano)",
	}

	t.Run("RenderizaCartaoExecutivoParaLeigos", func(t *testing.T) {
		rendered := decision.Render()
		expectedKeywords := []string{
			"CARTÃO DE DECISÃO EXECUTIVO",
			"Aprovação da Curadoria de Peças",
			"Tempo Previsto: 1m 30s",
			"Custo Estimado: $0.00",
		}

		for _, kw := range expectedKeywords {
			if !strings.Contains(rendered, kw) {
				t.Errorf("cartão deve conter %q, obteve:\n%s", kw, rendered)
			}
		}
	})

	t.Run("PromptAutopilotAprovaComSim", func(t *testing.T) {
		in := strings.NewReader("sim\n")
		var out bytes.Buffer

		approved, err := cli.PromptAutopilot(&out, in, decision)
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if !approved {
			t.Fatal("esperava approved == true para 'sim'")
		}
	})

	t.Run("PromptAutopilotRejeitaComNao", func(t *testing.T) {
		in := strings.NewReader("nao\n")
		var out bytes.Buffer

		approved, err := cli.PromptAutopilot(&out, in, decision)
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if approved {
			t.Fatal("esperava approved == false para 'nao'")
		}
	})
}

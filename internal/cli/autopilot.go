package cli

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// AutopilotDecisionForm formata a decisão humana em linguagem natural para operadores leigos (Dupla Superfície).
type AutopilotDecisionForm struct {
	Title         string `json:"title"`
	Summary       string `json:"summary"`
	EstimatedTime string `json:"estimated_time"`
	EstimatedCost string `json:"estimated_cost"`
}

// Render formata o cartão de decisão executivo em visual limpo e sem jargões de código.
func (c AutopilotDecisionForm) Render() string {
	return fmt.Sprintf(
		"📋 [CARTÃO DE DECISÃO EXECUTIVO // MODO AUTOPILOT]\n"+
			"· Ação:           %s\n"+
			"· Resumo Simples: %s\n"+
			"· Tempo Previsto: %s\n"+
			"· Custo Estimado: %s\n",
		c.Title,
		c.Summary,
		c.EstimatedTime,
		c.EstimatedCost,
	)
}

// PromptAutopilot solicita confirmação binária amigável [Sim/Não] para leigos.
func PromptAutopilot(w io.Writer, r io.Reader, form AutopilotDecisionForm) (bool, error) {
	fmt.Fprintf(w, "\n%s\nDeseja autorizar esta etapa? [Sim/Não]: ", form.Render())
	scanner := bufio.NewScanner(r)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, err
		}
		return false, nil
	}
	ans := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if ans == "sim" || ans == "s" || ans == "yes" || ans == "y" {
		return true, nil
	}
	return false, nil
}

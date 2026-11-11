package cli

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// ICAOConfirmacao descreve a trava de ação destrutiva no padrão da aviação civil (ICAO Readback/Hearback).
type ICAOConfirmacao struct {
	TargetResource  string `json:"target_resource"`
	MutationCommand string `json:"mutation_command"`
	BlastRadius     string `json:"blast_radius"`
	RollbackCommand string `json:"rollback_command"`
}

// Render formata o cartão ICAO em texto legível para confirmação no terminal.
func (c ICAOConfirmacao) Render() string {
	return fmt.Sprintf(
		"✈️ [ICAO READBACK // CONFIRMAÇÃO DE AÇÃO DESTRUTIVA]\n"+
			"· TARGET RESOURCE:  %s\n"+
			"· MUTATION COMMAND: %s\n"+
			"· BLAST RADIUS:     %s\n"+
			"· ROLLBACK COMMAND: %s\n",
		c.TargetResource,
		c.MutationCommand,
		c.BlastRadius,
		c.RollbackCommand,
	)
}

// PromptICAO exibe o cartão ICAO e aguarda confirmação estrita (HEARBACK APPROVED, s ou sim).
func PromptICAO(w io.Writer, r io.Reader, conf ICAOConfirmacao) (bool, error) {
	fmt.Fprintf(w, "\n%s\nPara autorizar a mutação, confirme com HEARBACK APPROVED: ", conf.Render())
	scanner := bufio.NewScanner(r)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, err
		}
		return false, nil
	}
	ans := strings.TrimSpace(scanner.Text())
	if strings.EqualFold(ans, "HEARBACK APPROVED") || strings.EqualFold(ans, "s") || strings.EqualFold(ans, "sim") {
		return true, nil
	}
	return false, nil
}

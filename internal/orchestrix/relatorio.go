package orchestrix

import (
	"fmt"
	"strings"

	"legoloop/internal/runner"
)

// notaSubstituicao é o prefixo das notas que o despacho grava quando troca de candidato.
const notaSubstituicao = "substituicao: "

type RelatorioNo struct {
	ID       string
	Commit   string
	CustoUSD float64
	Gate     string
	Harness  string
	Via      string
	Modelo   string
	// Substituicao descreve a troca de candidato em execução (vazio = sem troca).
	Substituicao string
}

// RelatorioDoRunner monta as linhas factuais a partir do que o motor realmente executou:
// harness, via e modelo usados, custo, status do gate e a Substituicao registrada nas notas.
// O commit vem de commits (id do nó para hash); nó sem commit conhecido fica vazio.
func RelatorioDoRunner(rep runner.Report, commits map[string]string) []RelatorioNo {
	nos := make([]RelatorioNo, 0, len(rep.Nodes))
	for _, n := range rep.Nodes {
		no := RelatorioNo{ID: n.ID, Commit: commits[n.ID], CustoUSD: n.Tokens.CostUSD, Gate: string(n.Status),
			Harness: n.Harness, Via: n.Via, Modelo: strings.TrimPrefix(n.Model, prefixoModeloNo)}
		for _, nota := range n.Notes {
			if t, ok := strings.CutPrefix(nota, notaSubstituicao); ok {
				no.Substituicao = t
			}
		}
		nos = append(nos, no)
	}
	return nos
}

type RelatorioDivergencia struct {
	No        string
	Descricao string
}

func GerarRelatorioFactual(nos []RelatorioNo, divergencias []RelatorioDivergencia) string {
	var b strings.Builder
	b.WriteString("## Parte factual\n\n")
	b.WriteString("| Nó | Commit | Custo (USD) | Gate | Harness | Via | Modelo | Substituição |\n")
	b.WriteString("| :--- | :--- | :---: | :---: | :--- | :--- | :--- | :--- |\n")
	var total float64
	for _, n := range nos {
		total += n.CustoUSD
		fmt.Fprintf(&b, "| %s | %s | %.4f | %s | %s | %s | %s | %s |\n",
			celula(n.ID), celula(n.Commit), n.CustoUSD, celula(n.Gate), celula(n.Harness), celula(n.Via), celula(n.Modelo), celula(n.Substituicao))
	}
	fmt.Fprintf(&b, "\nCusto total (USD): %.4f\n\n", total)
	b.WriteString("### Divergências do PLAN\n\n")
	if len(divergencias) == 0 {
		b.WriteString("Nenhuma.\n")
		return b.String()
	}
	for _, d := range divergencias {
		fmt.Fprintf(&b, "- %s: %s\n", celula(d.No), celula(d.Descricao))
	}
	return b.String()
}

func celula(s string) string {
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}

package orchestrix

import (
	"fmt"
	"strings"
)

type RelatorioNo struct {
	ID       string
	Commit   string
	CustoUSD float64
	Gate     string
	Harness  string
	Modelo   string
}

type RelatorioDivergencia struct {
	No        string
	Descricao string
}

func GerarRelatorioFactual(nos []RelatorioNo, divergencias []RelatorioDivergencia) string {
	var b strings.Builder
	b.WriteString("## Parte factual\n\n")
	b.WriteString("| Nó | Commit | Custo (USD) | Gate | Harness | Modelo |\n")
	b.WriteString("| :--- | :--- | :---: | :---: | :--- | :--- |\n")
	var total float64
	for _, n := range nos {
		total += n.CustoUSD
		fmt.Fprintf(&b, "| %s | %s | %.4f | %s | %s | %s |\n",
			celula(n.ID), celula(n.Commit), n.CustoUSD, celula(n.Gate), celula(n.Harness), celula(n.Modelo))
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

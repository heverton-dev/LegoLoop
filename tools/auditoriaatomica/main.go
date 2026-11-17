// Comando auditoriaatomica materializa o Processo auditoria-atomica e roda a auditoria universal
// dos seis verbos do LegoLoop, gravando o relatório com o veredito de cada um.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"legoloop/internal/fabricar"
)

// nomeRelatorioVerbos é o relatório de auditoria universal dos seis verbos.
const nomeRelatorioVerbos = "2026-11-17_REPORT_AUDITORIA_6_VERBOS.md"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "uso: auditoriaatomica gerar|auditar")
		os.Exit(2)
	}
	switch os.Args[1] {
	case "gerar":
		gerar()
	case "auditar":
		auditar()
	default:
		fmt.Fprintf(os.Stderr, "comando desconhecido: %s\n", os.Args[1])
		os.Exit(2)
	}
}

// gerar materializa o Processo auditoria-atomica no nascimento 5/5 do scaffold.
func gerar() {
	gravados, err := fabricar.MaterializarAuditoriaAtomica(fabricar.Roots{Vault: ".", Docs: "."}, time.Now().UTC())
	if err != nil {
		fmt.Fprintf(os.Stderr, "materializar: %v\n", err)
		os.Exit(1)
	}
	for _, c := range gravados {
		fmt.Println(c)
	}
}

// auditar roda a auditoria universal dos seis verbos e grava o relatório atômico
// de cada verbo mais o relatório consolidado com o veredito de cada um.
func auditar() {
	dirAtomicos := filepath.Join("docs", "reports", "auditoria-6-verbos")
	rels, codigos, err := fabricar.AuditarVerbosEGravar(".", fabricar.TestePadraoEm("."), dirAtomicos)
	if err != nil {
		fmt.Fprintf(os.Stderr, "auditar: %v\n", err)
		os.Exit(2)
	}
	saida := filepath.Join("docs", "reports", nomeRelatorioVerbos)
	if err := fabricar.EscreverRelatorio(saida, fabricar.RelatorioUniversal(rels, time.Now().UTC())); err != nil {
		fmt.Fprintf(os.Stderr, "gravar relatório: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("relatórios atômicos:", dirAtomicos)
	fmt.Println("relatório:", saida)
	exit := 0
	for i, rel := range rels {
		fmt.Printf("%s: %s (exit %d)\n", rel.Alvo.Nome, rel.Veredito, codigos[i])
		if codigos[i] != 0 {
			exit = 1
		}
	}
	os.Exit(exit)
}

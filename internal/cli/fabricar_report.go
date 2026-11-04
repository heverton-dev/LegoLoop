package cli

import (
	"fmt"
	"io"
	"strings"

	"legoloop/internal/fabricar"
)

// PrintReportTexto imprime o resumo por camada em formato texto human-readable.
// Resumo por camada: espelho e confiança, peças vistas e cortes, decisões com motivo,
// ACs e gates, arquivos e harnesses.
func PrintReportTexto(w io.Writer, r *fabricar.Report) {
	if r == nil {
		return
	}
	fmt.Fprintln(w, "\n--- Resumo por Camada ---")
	printCamada1(w, r)
	printCamada2(w, r)
	printCamada3(w, r)
	printCamada4(w, r)
	printCamada5(w, r)
}

func printCamada1(w io.Writer, r *fabricar.Report) {
	fmt.Fprintf(w, "[Camada 1 - Clarificação] Confiança: %d%% | Resultado: %s\n",
		r.Clarificacao.Confianca, r.Clarificacao.Espelho.Outcome)
}

func printCamada2(w io.Writer, r *fabricar.Report) {
	var totalPecas int
	for _, c := range r.Varredura.PecasVistas {
		totalPecas += c
	}
	cortesStr := "nenhum"
	if len(r.Varredura.Cortes) > 0 {
		cortesStr = strings.Join(r.Varredura.Cortes, ", ")
	}
	fmt.Fprintf(w, "[Camada 2 - Varredura] Peças vistas: %d | Cortes: %s\n", totalPecas, cortesStr)
}

func printCamada3(w io.Writer, r *fabricar.Report) {
	fmt.Fprintf(w, "[Camada 3 - Fusão] Decisões (%d):\n", len(r.Fusao.Decisoes))
	for _, d := range r.Fusao.Decisoes {
		entrou := "não"
		if d.Entrou {
			entrou = "sim"
		}
		peca := d.Peca
		if peca == "" {
			peca = "(nenhuma)"
		}
		fmt.Fprintf(w, "  - Peça: %s | Entrou: %s | Motivo: %s\n", peca, entrou, d.Motivo)
	}
}

func printCamada4(w io.Writer, r *fabricar.Report) {
	fmt.Fprintf(w, "[Camada 4 - Tríade] Critérios (ACs): %d | Tickets: %d | Gates: %s\n",
		r.Triade.ACs, r.Triade.Tickets, strings.Join(r.Triade.Gates, ", "))
}

func printCamada5(w io.Writer, r *fabricar.Report) {
	fmt.Fprintf(w, "[Camada 5 - Gravação] Arquivos: %d (%d bytes) | Harnesses: %s\n",
		len(r.Gravacao.Arquivos), r.Gravacao.Bytes, strings.Join(r.Gravacao.Harnesses, ", "))
}

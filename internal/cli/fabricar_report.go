package cli

import (
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"

	"legoloop/internal/fabricar"
)

// reportLacuna é a saída do dry-run com lacuna: os avisos, depois o cartão e o resumo por camada (ou o JSON).
// O erro e o código de saída 1 ficam com o chamador.
func reportLacuna(d Deps, res *fabricar.Result, o fabricarOpts) {
	for _, w := range res.Warnings {
		fmt.Fprintln(d.Stderr, "aviso:", w)
	}
	if o.asJSON {
		printJSON(d, res)
		return
	}
	fmt.Fprintf(d.Stdout, "\n(dry-run) o processo %s tem lacuna; nada seria gravado.\n\n%s", res.CanonicalName, res.Report.Cartao.Texto())
	PrintReportTexto(d.Stdout, res.Report)
}

// PrintReportTexto imprime o resumo por camada em formato texto human-readable.
// Resumo por camada: espelho e confiança, matriz e cortes, peça ou lacuna por etapa,
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

// printCamada2 resume a matriz: quantas peças, a idade dela e quantos cortes (a lista vai no --json).
func printCamada2(w io.Writer, r *fabricar.Report) {
	var total int
	for _, c := range r.Varredura.PecasVistas {
		total += c
	}
	fmt.Fprintf(w, "[Camada 2 - Matriz] Peças na matriz: %d | Idade: %d h | Cortes: %d\n",
		total, int(r.Varredura.MatrizIdade.Hours()), len(r.Varredura.Cortes))
}

// printCamada3 mostra as necessidades, uma linha por etapa (a peça escolhida ou a lacuna) e o que ficou de
// fora agrupado por motivo. Cada decisão, com o motivo inteiro, vai no --json.
func printCamada3(w io.Writer, r *fabricar.Report) {
	n := r.Fusao.Necessidades
	fmt.Fprintf(w, "[Camada 3 - Fusão] Necessidades: %s | etapas: %s\n", n.Natureza, strings.Join(n.Etapas, ", "))
	var etapas []fabricar.Decisao
	fora := map[string]int{}
	for _, d := range r.Fusao.Decisoes {
		if d.Entrou || d.Etapa != "" {
			etapas = append(etapas, d)
			continue
		}
		fora[motivoCurto(d.Motivo)]++
	}
	printEtapas(w, etapas)
	printFora(w, fora)
}

// printEtapas alinha as etapas numa coluna de largura fixa.
func printEtapas(w io.Writer, ds []fabricar.Decisao) {
	largura := 0
	for _, d := range ds {
		largura = max(largura, len([]rune(rotuloEtapa(d))))
	}
	for _, d := range ds {
		pad := strings.Repeat(" ", largura-len([]rune(rotuloEtapa(d))))
		if d.Entrou {
			fmt.Fprintf(w, "  %s%s  %s (aptidão %.2f) | %s\n", rotuloEtapa(d), pad, d.Peca, d.Aptidao, d.Motivo)
		} else {
			fmt.Fprintf(w, "  %s%s  LACUNA | %s\n", rotuloEtapa(d), pad, d.Motivo)
		}
	}
}

func rotuloEtapa(d fabricar.Decisao) string {
	if d.Etapa == "" {
		return "geral"
	}
	return d.Etapa
}

// printFora conta as peças de fora por motivo, em ordem alfabética do motivo.
func printFora(w io.Writer, fora map[string]int) {
	if len(fora) == 0 {
		return
	}
	total, partes := 0, make([]string, 0, len(fora))
	for _, m := range slices.Sorted(maps.Keys(fora)) {
		total += fora[m]
		partes = append(partes, fmt.Sprintf("%s: %d", m, fora[m]))
	}
	fmt.Fprintf(w, "  Fora: %d peças (%s)\n", total, strings.Join(partes, ", "))
}

// motivoCurto é a parte do motivo antes dos dois-pontos (ex.: "natureza", "aptidão desconhecida").
func motivoCurto(m string) string {
	if i := strings.Index(m, ":"); i > 0 {
		return m[:i]
	}
	return m
}

func printCamada4(w io.Writer, r *fabricar.Report) {
	fmt.Fprintf(w, "[Camada 4 - Tríade] Critérios (ACs): %d | Tickets: %d | Gates: %s\n",
		r.Triade.ACs, r.Triade.Tickets, strings.Join(r.Triade.Gates, ", "))
}

func printCamada5(w io.Writer, r *fabricar.Report) {
	fmt.Fprintf(w, "[Camada 5 - Gravação] Arquivos: %d (%d bytes) | Harnesses: %s\n",
		len(r.Gravacao.Arquivos), r.Gravacao.Bytes, strings.Join(r.Gravacao.Harnesses, ", "))
}

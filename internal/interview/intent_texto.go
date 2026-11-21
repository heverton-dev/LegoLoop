package interview

import (
	"fmt"
	"strings"
)

// tituloMax é o tamanho máximo, em caracteres, do título de uma INTENT.
const tituloMax = 100

// tituloDoTema devolve um título curto a partir da ideia inicial: a primeira frase, cortada na fronteira
// de palavra quando passa de tituloMax. A ideia completa fica só em "O Problema Real".
func tituloDoTema(theme string) string {
	t := strings.TrimSpace(theme)
	if i := strings.IndexAny(t, "\r\n"); i >= 0 {
		t = strings.TrimSpace(t[:i])
	}
	if i := strings.Index(t, ". "); i >= 0 {
		t = t[:i+1]
	}
	r := []rune(t)
	if len(r) <= tituloMax {
		return t
	}
	corte := string(r[:tituloMax])
	if i := strings.LastIndex(corte, " "); i > 0 {
		corte = corte[:i]
	}
	return strings.TrimRight(corte, " ,;:.") + "…"
}

// criteriosDe separa o critério de sucesso do espelho, juntado com "; ", em um item por critério.
func criteriosDe(s string) []string {
	var out []string
	for _, c := range strings.Split(s, "; ") {
		if c = strings.TrimSpace(c); c != "" {
			out = append(out, c)
		}
	}
	return out
}

// criteriosNumerados escreve os critérios em uma linha só, numerados, para a célula da tabela.
func criteriosNumerados(s string) string {
	itens := criteriosDe(s)
	if len(itens) <= 1 {
		return strings.TrimSpace(s)
	}
	partes := make([]string, len(itens))
	for i, c := range itens {
		partes[i] = fmt.Sprintf("(%d) %s", i+1, c)
	}
	return strings.Join(partes, " ")
}

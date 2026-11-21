package mcp

import (
	"sort"
	"strings"
	"unicode"
)

// tamanhoMinimoTermo descarta palavras curtas ("de", "no", "e"), que casam com qualquer texto.
const tamanhoMinimoTermo = 3

// termosDaBusca separa a consulta em palavras minúsculas, sem duplicatas e sem as muito curtas.
// Underscore e hífen separam palavras: "trace_path" vira "trace" e "path".
func termosDaBusca(q string) []string {
	campos := strings.FieldsFunc(strings.ToLower(q), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	vistos := map[string]bool{}
	var out []string
	for _, c := range campos {
		if len([]rune(c)) >= tamanhoMinimoTermo && !vistos[c] {
			vistos[c] = true
			out = append(out, c)
		}
	}
	return out
}

// pontuarBusca conta quantas palavras da consulta aparecem no texto. Consulta sem palavras úteis
// pontua 0, e a consulta vazia pontua 1 para listar tudo.
func pontuarBusca(q, texto string) int {
	if strings.TrimSpace(q) == "" {
		return 1
	}
	hay := strings.ToLower(texto)
	pontos := 0
	for _, t := range termosDaBusca(q) {
		if strings.Contains(hay, t) {
			pontos++
		}
	}
	return pontos
}

// minimoDePontos é o piso para uma ferramenta entrar no resultado: metade das palavras, arredondada para cima.
func minimoDePontos(q string) int {
	n := len(termosDaBusca(q))
	if n == 0 {
		return 1
	}
	return (n + 1) / 2
}

// ordenarMatches põe a ferramenta com mais palavras casadas primeiro; o empate cai em servidor e nome.
func ordenarMatches(q string, ms []ToolMatch) {
	sort.SliceStable(ms, func(i, j int) bool {
		pi := pontuarBusca(q, ms[i].Server+" "+ms[i].Name+" "+ms[i].Description)
		pj := pontuarBusca(q, ms[j].Server+" "+ms[j].Name+" "+ms[j].Description)
		if pi != pj {
			return pi > pj
		}
		if ms[i].Server != ms[j].Server {
			return ms[i].Server < ms[j].Server
		}
		return ms[i].Name < ms[j].Name
	})
}

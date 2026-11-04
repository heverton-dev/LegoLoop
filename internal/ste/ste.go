// Package ste aplica o perfil STE-PT às mensagens do produto: português técnico simples, inspirado no
// ASD-STE100 (frases curtas, um termo por conceito, aviso de risco antes da ação). É um perfil nosso,
// não uma certificação do padrão. O glossário é um arquivo versionado (glossario.json).
package ste

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

//go:embed glossario.json
var glossarioJSON []byte

// Violation é uma regra do perfil que a mensagem não cumpre.
type Violation struct {
	Msg   string `json:"msg"`
	Regra string `json:"regra"` // frase-longa | ambigua | termo | sem-aviso | sem-desfazer
}

type glossario struct {
	Versao      int                 `json:"versao"`
	MaxPalavras int                 `json:"max_palavras"`
	Ambiguas    []string            `json:"ambiguas"`
	Termos      map[string][]string `json:"termos"` // termo preferido -> sinônimos proibidos
}

// proibido liga um sinônimo proibido ao termo que o substitui.
type proibido struct {
	re        *regexp.Regexp
	palavra   string
	preferido string
}

var (
	gl        = mustGlossario()
	ambiguas  = compilar(gl.Ambiguas)
	proibidos = compilarProibidos(gl.Termos)
	fimFrase  = regexp.MustCompile(`[.!?]+(?:\s+|$)|\n+`)
)

func mustGlossario() glossario {
	var g glossario
	if err := json.Unmarshal(glossarioJSON, &g); err != nil || g.MaxPalavras <= 0 {
		panic(fmt.Sprintf("ste: glossario.json inválido: %v", err))
	}
	return g
}

// palavraInteira casa a expressão só quando ela não é parte de outra palavra (letras Unicode).
func palavraInteira(termo string) *regexp.Regexp {
	return regexp.MustCompile(`(?i)(?:^|[^\p{L}\p{N}])` + regexp.QuoteMeta(termo) + `(?:[^\p{L}\p{N}]|$)`)
}

func compilar(termos []string) map[string]*regexp.Regexp {
	out := make(map[string]*regexp.Regexp, len(termos))
	for _, t := range termos {
		out[t] = palavraInteira(t)
	}
	return out
}

func compilarProibidos(termos map[string][]string) []proibido {
	var out []proibido
	for preferido, sins := range termos {
		for _, s := range sins {
			out = append(out, proibido{re: palavraInteira(s), palavra: s, preferido: preferido})
		}
	}
	return out
}

// frases divide o texto em frases. Ponto colado em letra ou número (AC-E01.1) não encerra a frase.
func frases(msg string) []string {
	var out []string
	for _, f := range fimFrase.Split(msg, -1) {
		if f = strings.TrimSpace(f); f != "" {
			out = append(out, f)
		}
	}
	return out
}

// palavras conta os tokens que têm letra ou número; hífen e símbolos soltos não contam.
func palavras(frase string) int {
	n := 0
	for _, tok := range strings.Fields(frase) {
		if strings.IndexFunc(tok, func(r rune) bool { return unicode.IsLetter(r) || unicode.IsNumber(r) }) >= 0 {
			n++
		}
	}
	return n
}

// Lint aplica o perfil STE-PT a uma mensagem e devolve as violações, na ordem em que aparecem.
func Lint(msg string) []Violation {
	var out []Violation
	for _, f := range frases(msg) {
		if n := palavras(f); n > gl.MaxPalavras {
			out = append(out, Violation{Msg: fmt.Sprintf("frase com %d palavras (máximo %d): %q", n, gl.MaxPalavras, f), Regra: "frase-longa"})
		}
	}
	return append(out, lintTermos(msg)...)
}

func lintTermos(msg string) []Violation {
	var out []Violation
	for termo, re := range ambiguas {
		if re.MatchString(msg) {
			out = append(out, Violation{Msg: fmt.Sprintf("palavra ambígua %q: diga o que acontece de forma exata", termo), Regra: "ambigua"})
		}
	}
	for _, p := range proibidos {
		if p.re.MatchString(msg) {
			out = append(out, Violation{Msg: fmt.Sprintf("use %q no lugar de %q", p.preferido, p.palavra), Regra: "termo"})
		}
	}
	return out
}

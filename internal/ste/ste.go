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
	Versao                  int                 `json:"versao"`
	MaxPalavras             int                 `json:"max_palavras"`
	MaxPalavrasProcedimento int                 `json:"max_palavras_procedimento"`
	MaxPalavrasDescricao    int                 `json:"max_palavras_descricao"`
	Ambiguas                []string            `json:"ambiguas"`
	Termos                  map[string][]string `json:"termos"` // termo preferido -> sinônimos proibidos
}

// proibido liga um sinônimo proibido ao termo que o substitui.
type proibido struct {
	re        *regexp.Regexp
	palavra   string
	preferido string
}

var (
	gl           = mustGlossario()
	ambiguas     = compilar(gl.Ambiguas)
	proibidos    = compilarProibidos(gl.Termos)
	fimFrase     = regexp.MustCompile(`[.!?]+(?:\s+|$)|\n+`)
	reVozPassiva = regexp.MustCompile(`(?i)\b(?:é|são|foi|foram|será|serão|sendo|sido)\s+(?:[a-zà-ú]+(?:ado|ada|ados|adas|ido|ida|idos|idas)|feito|feita|feitos|feitas|escrito|escrita|escritos|escritas|aberto|aberta|abertos|abertas)\b`)
	imperativos  = map[string]bool{
		"grave": true, "crie": true, "remova": true, "execute": true,
		"rode": true, "leia": true, "abra": true, "feche": true,
		"escreva": true, "edite": true, "compile": true, "instale": true,
		"valide": true, "verifique": true, "teste": true, "envie": true,
		"configure": true, "defina": true, "gere": true, "atualize": true,
		"adicione": true, "apague": true, "limpe": true, "pare": true,
		"inicie": true, "dê": true, "use": true, "mude": true,
		"veja": true, "liste": true, "substitua": true, "insira": true,
		"marque": true, "faça": true, "mostre": true, "copie": true,
	}
)

func mustGlossario() glossario {
	var g glossario
	if err := json.Unmarshal(glossarioJSON, &g); err != nil || g.MaxPalavras <= 0 {
		panic(fmt.Sprintf("ste: glossario.json inválido: %v", err))
	}
	if g.MaxPalavrasProcedimento <= 0 {
		g.MaxPalavrasProcedimento = g.MaxPalavras
	}
	if g.MaxPalavrasDescricao <= 0 {
		g.MaxPalavrasDescricao = 25
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

func isInstrucao(frase string) bool {
	for _, f := range strings.Fields(frase) {
		clean := strings.ToLower(strings.Trim(f, `*·-_.:,;"'`))
		if clean == "" {
			continue
		}
		return imperativos[clean]
	}
	return false
}

func lintFrases(msg string) []Violation {
	var out []Violation
	for _, f := range frases(msg) {
		n := palavras(f)
		max := gl.MaxPalavrasDescricao
		tipo := "descrição"
		if isInstrucao(f) {
			max = gl.MaxPalavrasProcedimento
			tipo = "instrução"
		}
		if n > max {
			out = append(out, Violation{
				Msg:   fmt.Sprintf("%s com %d palavras (máximo %d): %q", tipo, n, max, f),
				Regra: "frase-longa",
			})
		}
	}
	return out
}

func lintVozPassiva(msg string) []Violation {
	var out []Violation
	for _, m := range reVozPassiva.FindAllString(msg, -1) {
		out = append(out, Violation{
			Msg:   fmt.Sprintf("voz passiva %q detectada: use voz ativa com o ator explícito", m),
			Regra: "voz-passiva",
		})
	}
	return out
}

// Lint aplica o perfil STE-PT a uma mensagem e devolve as violações, na ordem em que aparecem.
func Lint(msg string) []Violation {
	return append(lintFrases(msg), lintTermos(msg)...)
}

// LintStrict aplica o perfil completo ASD-STE100 (80%), incluindo checagem de voz ativa.
func LintStrict(msg string) []Violation {
	return append(Lint(msg), lintVozPassiva(msg)...)
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

package main

import (
	"regexp"
	"strings"
	"unicode"

	"legoloop/internal/taxonomy"
)

// brandRe casa a marca LegoLoop em qualquer caixa. A marca é exceção ao dicionário.
var brandRe = regexp.MustCompile(`(?i)legoloop`)

var wordRe = regexp.MustCompile(`[\p{L}]+`)

// maskBrand troca a marca por espaços do mesmo tamanho.
func maskBrand(s string) string {
	return brandRe.ReplaceAllStringFunc(s, func(m string) string {
		return strings.Repeat(" ", len(m))
	})
}

// splitIdentifiers separa camelCase, snake_case, kebab-case, pontos e barras em palavras.
// A quebra de linha permanece, para o número de linha continuar exato.
func splitIdentifiers(s string) string {
	runes := []rune(s)
	var b strings.Builder
	for i, r := range runes {
		if r == '\n' {
			b.WriteRune(r)
			continue
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			b.WriteRune(' ')
			continue
		}
		if i > 0 && camelBoundary(runes, i) {
			b.WriteRune(' ')
		}
		b.WriteRune(r)
	}
	return b.String()
}

func camelBoundary(runes []rune, i int) bool {
	prev, cur := runes[i-1], runes[i]
	if !unicode.IsUpper(cur) {
		return false
	}
	if unicode.IsLower(prev) || unicode.IsDigit(prev) {
		return true
	}
	return unicode.IsUpper(prev) && i+1 < len(runes) && unicode.IsLower(runes[i+1])
}

// splitAffixes separa o prefixo "sub" e o plural "s" quando o resto é termo proibido.
func splitAffixes(s string) string {
	return wordRe.ReplaceAllStringFunc(s, func(w string) string {
		lw := strings.ToLower(w)
		if len(lw) > 3 && strings.HasPrefix(lw, "sub") {
			if _, ok := taxonomy.CheckForbidden(lw[3:]); ok {
				return w[:3] + " " + w[3:]
			}
		}
		if len(lw) > 2 && strings.HasSuffix(lw, "s") {
			if _, ok := taxonomy.CheckForbidden(lw[:len(lw)-1]); ok {
				return w[:len(w)-1] + " s"
			}
		}
		return w
	})
}

// normalizeText prepara qualquer texto para o dicionário: marca fora, identificadores separados.
func normalizeText(s string) string {
	return splitAffixes(splitIdentifiers(maskBrand(s)))
}

// validateNormalized valida o texto normalizado e devolve as violações com linha relativa.
func validateNormalized(text string) []taxonomy.TaxonomyViolation {
	vs, err := taxonomy.ValidateText(normalizeText(text))
	if err != nil {
		return nil
	}
	return vs
}

package fabricar

import (
	"errors"
	"strings"
	"unicode"
)

// ErrInvalidName: o nome não produz um slug válido após a normalização.
var ErrInvalidName = errors.New("fabricar: nome inválido")

const (
	namePrefix = "legoloop-"
	maxNameLen = 64 // limite de studiodoc.ValidateID
)

// translit mapeia letras latinas acentuadas (Latin-1 e Latin Extended-A) para ASCII.
// Só a biblioteca padrão: sem golang.org/x/text.
var translit = buildTranslit([][2]string{
	{"a", "àáâãäåāăąǎ"}, {"c", "çćĉċč"}, {"d", "ďđ"}, {"e", "èéêëēĕėęě"}, {"g", "ĝğġģ"},
	{"h", "ĥħ"}, {"i", "ìíîïĩīĭįı"}, {"j", "ĵ"}, {"k", "ķ"}, {"l", "ĺļľŀł"},
	{"n", "ñńņňŉŋ"}, {"o", "òóôõöøōŏő"}, {"r", "ŕŗř"}, {"s", "śŝşš"}, {"t", "ţťŧ"},
	{"u", "ùúûüũūŭůűų"}, {"w", "ŵ"}, {"y", "ýÿŷ"}, {"z", "źżž"},
	{"ss", "ß"}, {"ae", "æ"}, {"oe", "œ"}, {"th", "þ"},
})

func buildTranslit(groups [][2]string) map[rune]string {
	m := make(map[rune]string)
	for _, g := range groups {
		for _, r := range g[1] {
			m[r] = g[0]
		}
	}
	return m
}

// ascii devolve a forma ASCII minúscula de r; "" quando não há equivalente.
func ascii(r rune) string {
	r = unicode.ToLower(r)
	if r < unicode.MaxASCII {
		return string(r)
	}
	return translit[r]
}

// words quebra o texto em palavras [a-z0-9], tratando qualquer outro caractere como separador.
func words(raw string) []string {
	var out []string
	var cur strings.Builder
	flush := func() {
		if cur.Len() > 0 {
			out = append(out, cur.String())
			cur.Reset()
		}
	}
	for _, r := range raw {
		for _, c := range ascii(r) {
			if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
				cur.WriteRune(c)
			} else {
				flush()
			}
		}
		if ascii(r) == "" {
			flush()
		}
	}
	flush()
	return out
}

// Slug devolve só a parte variável do nome canônico (sem o prefixo legoloop-), ou ErrInvalidName.
func Slug(raw string) (string, error) {
	w := words(raw)
	if len(w) > 0 && w[0] == "legoloop" {
		w = w[1:]
	}
	slug := strings.Join(w, "-")
	if max := maxNameLen - len(namePrefix); len(slug) > max {
		slug = strings.TrimRight(slug[:max], "-")
	}
	if slug == "" {
		return "", ErrInvalidName
	}
	return slug, nil
}

// NormalizeProcessName devolve o nome canônico legoloop-<kebab-case>, com no máximo 64 caracteres.
// O prefixo é imutável e a função é idempotente.
func NormalizeProcessName(raw string) (string, error) {
	slug, err := Slug(raw)
	if err != nil {
		return "", err
	}
	return namePrefix + slug, nil
}

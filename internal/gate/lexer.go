// Package gate compila e avalia expressões de Gate (Pass/Fail) do arquivo DAG.
//
// Gramática:
//
//	or    := and ( "||" and )*
//	and   := unary ( "&&" unary )*
//	unary := "!" unary | cmp
//	cmp   := term ( ("=="|"!="|">"|">="|"<"|"<=") term )?
//	term  := call | number | string | ident | "(" or ")"
package gate

import (
	"fmt"
	"strings"
	"unicode"
)

type tokKind int

const (
	tkEOF tokKind = iota
	tkIdent
	tkNumber
	tkString
	tkOp
	tkLParen
	tkRParen
	tkComma
)

type token struct {
	kind tokKind
	text string
	pos  int
}

var twoCharOps = []string{"&&", "||", "==", "!=", ">=", "<="}

func lex(src string) ([]token, error) {
	var out []token
	for i := 0; i < len(src); {
		c := rune(src[i])
		switch {
		case unicode.IsSpace(c):
			i++
		case c == '\'':
			tok, next, err := lexString(src, i)
			if err != nil {
				return nil, err
			}
			out, i = append(out, tok), next
		case unicode.IsDigit(c):
			j := scanWhile(src, i, func(r rune) bool { return unicode.IsDigit(r) || r == '.' })
			out, i = append(out, token{tkNumber, src[i:j], i}), j
		case unicode.IsLetter(c) || c == '_':
			j := scanWhile(src, i, func(r rune) bool { return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' })
			out, i = append(out, token{tkIdent, src[i:j], i}), j
		default:
			tok, err := lexPunct(src, i)
			if err != nil {
				return nil, err
			}
			out, i = append(out, tok), i+len(tok.text)
		}
	}
	return append(out, token{tkEOF, "", len(src)}), nil
}

func scanWhile(src string, i int, ok func(rune) bool) int {
	for i < len(src) && ok(rune(src[i])) {
		i++
	}
	return i
}

// lexString lê 'literal' com escapes \' e \\ ; demais barras são preservadas (regex).
func lexString(src string, start int) (token, int, error) {
	var b strings.Builder
	for i := start + 1; i < len(src); i++ {
		c := src[i]
		if c == '\\' && i+1 < len(src) && (src[i+1] == '\'' || src[i+1] == '\\') {
			b.WriteByte(src[i+1])
			i++
			continue
		}
		if c == '\'' {
			return token{tkString, b.String(), start}, i + 1, nil
		}
		b.WriteByte(c)
	}
	return token{}, 0, fmt.Errorf("gate: string sem fechamento na posição %d", start)
}

func lexPunct(src string, i int) (token, error) {
	for _, op := range twoCharOps {
		if strings.HasPrefix(src[i:], op) {
			return token{tkOp, op, i}, nil
		}
	}
	switch src[i] {
	case '(':
		return token{tkLParen, "(", i}, nil
	case ')':
		return token{tkRParen, ")", i}, nil
	case ',':
		return token{tkComma, ",", i}, nil
	case '!', '>', '<':
		return token{tkOp, src[i : i+1], i}, nil
	}
	return token{}, fmt.Errorf("gate: caractere inesperado %q na posição %d", src[i], i)
}

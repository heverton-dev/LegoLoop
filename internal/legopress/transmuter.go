package legopress

import (
	"regexp"
	"strings"
)

var slopPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?m)(?i)^.*com certeza.*$\r?\n?`),
	regexp.MustCompile(`(?m)(?i)^.*vou ficar muito feliz.*$\r?\n?`),
	regexp.MustCompile(`(?m)(?i)^.*aqui está a explicação passo a passo.*$\r?\n?`),
	regexp.MustCompile(`(?m)(?i)^.*espero que isso tenha ajudado.*$\r?\n?`),
	regexp.MustCompile(`(?m)(?i)^.*se precisar de mais alguma coisa.*$\r?\n?`),
	regexp.MustCompile(`(?m)(?i)^.*espero ter ajudado.*$\r?\n?`),
}

type Transmuter struct{}

func NewTransmuter() *Transmuter {
	return &Transmuter{}
}

func (t *Transmuter) Sanitize(input string) string {
	res := input
	for _, p := range slopPatterns {
		res = p.ReplaceAllString(res, "")
	}
	return strings.TrimSpace(res)
}

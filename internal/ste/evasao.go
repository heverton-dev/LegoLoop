package ste

import "regexp"

// reEvasao reconhece frase que declara pendencia, falta de verificacao ou adiamento.
// A lei da prova proibe essas frases: fato sem prova se verifica antes de afirmar e achado se corrige
// na mesma execucao. O \b do RE2 e ASCII; todos os limites caem em letra ASCII.
var reEvasao = regexp.MustCompile(`(?i)\bn[ãa]o verificad[oa]s?\b|\bn[ãa]o consegui\b|\bainda falta(?:m)?\b|\bpend[êe]ncias?\b|` +
	`\bpr[óo]xima etapa\b|\bfica(?:m)? para depois\b|\bfalta(?:m)? (?:verificar|validar|testar|rodar)\b|` +
	`\bn[ãa]o foi poss[íi]vel verificar\b|\bn[ãa]o testei\b|\bn[ãa]o rodei\b`)

// reCitacao reconhece trecho entre aspas, crases ou chevrons: citar a frase proibida nao e usa-la.
var reCitacao = regexp.MustCompile("(?s)\"[^\"]*\"|`[^`]*`|“[^”]*”|«[^»]*»")

// LintEvasao devolve uma violacao por frase evasiva encontrada na mensagem, sem contar as citacoes.
func LintEvasao(msg string) []Violation {
	msg = reCitacao.ReplaceAllString(msg, " ")
	var out []Violation
	for _, f := range frases(msg) {
		if m := reEvasao.FindString(f); m != "" {
			out = append(out, Violation{
				Msg:   "evasão de prova " + quote(m) + ": verifique agora ou corrija o achado na mesma resposta",
				Regra: "evasao-de-prova",
			})
		}
	}
	return out
}

func quote(s string) string { return "\"" + s + "\"" }

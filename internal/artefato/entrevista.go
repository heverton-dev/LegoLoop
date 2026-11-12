package artefato

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	entrevistaID      = regexp.MustCompile(`(?m)^id:\s*INTERVIEW-`)
	linhaPergunta     = regexp.MustCompile(`^\|\s*Q\d+\s*\|`)
	linhaFato         = regexp.MustCompile(`^\|\s*F\d+\s*\|`)
	fatoRef           = regexp.MustCompile(`\bF\d+\b`)
	resultadosValidos = []string{"VERIFICADO", "REFUTADO", "PARCIAL", "NÃO REPRODUZIDO"}
)

type fatoEntrevista struct{ prova, resultado string }

func celulas(linha string) []string {
	partes := strings.Split(strings.Trim(strings.TrimSpace(linha), "|"), "|")
	for i := range partes {
		partes[i] = strings.TrimSpace(partes[i])
	}
	return partes
}

func novoFato(c []string) fatoEntrevista {
	if len(c) < 4 {
		return fatoEntrevista{}
	}
	return fatoEntrevista{
		prova:     strings.TrimSpace(strings.Join(c[2:len(c)-1], "|")),
		resultado: c[len(c)-1],
	}
}

func classificarResultado(resultado string) string {
	for _, r := range resultadosValidos {
		if strings.HasPrefix(resultado, r) {
			return r
		}
	}
	return ""
}

func coletarFatos(linhas []string) map[string]fatoEntrevista {
	fatos := map[string]fatoEntrevista{}
	for _, l := range linhas {
		if linhaFato.MatchString(l) {
			c := celulas(l)
			fatos[c[0]] = novoFato(c)
		}
	}
	return fatos
}

func falhasFato(id string, f fatoEntrevista) []string {
	var out []string
	if f.prova == "" {
		out = append(out, fmt.Sprintf("%s: sem prova", id))
	}
	if classificarResultado(f.resultado) == "" {
		out = append(out, fmt.Sprintf("%s: resultado inválido %q", id, f.resultado))
	}
	return out
}

func falhasPergunta(c []string, fatos map[string]fatoEntrevista) []string {
	var refs []string
	if len(c) > 2 {
		refs = fatoRef.FindAllString(c[2], -1)
	}
	if len(refs) == 0 {
		return []string{c[0] + ": sem fato citado"}
	}
	var out []string
	refutados := 0
	for _, ref := range refs {
		f, ok := fatos[ref]
		if !ok {
			out = append(out, fmt.Sprintf("%s: fato %s não existe", c[0], ref))
		} else if classificarResultado(f.resultado) == "REFUTADO" {
			refutados++
		}
	}
	if refutados == len(refs) {
		out = append(out, c[0]+": cita só fatos REFUTADO")
	}
	return out
}

func ValidarEntrevista(texto string) []string {
	if !entrevistaID.MatchString(texto) {
		return nil
	}
	linhas := strings.Split(strings.ReplaceAll(texto, "\r", ""), "\n")
	fatos := coletarFatos(linhas)
	var falhas []string
	for _, l := range linhas {
		switch {
		case linhaFato.MatchString(l):
			c := celulas(l)
			falhas = append(falhas, falhasFato(c[0], fatos[c[0]])...)
		case linhaPergunta.MatchString(l):
			falhas = append(falhas, falhasPergunta(celulas(l), fatos)...)
		}
	}
	return falhas
}

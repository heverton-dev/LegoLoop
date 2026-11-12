package artefato

import (
	"strings"
	"testing"
)

const cabecalhoEntrevista = "---\nid: INTERVIEW-X\n---\n\n| ID | Pergunta | Fatos | Decisão |\n| :--- | :--- | :--- | :--- |\n"

func entrevista(perguntas, fatos string) string {
	return cabecalhoEntrevista + perguntas + "\n| ID | Afirmação | Prova | Resultado |\n| :--- | :--- | :--- | :--- |\n" + fatos
}

func TestEntrevistaSemFatoFalha(t *testing.T) {
	texto := entrevista(
		"| Q1 | a? | | x |\n| Q2 | b? | F9 | x |\n| Q3 | c? | F1 | x |\n",
		"| F1 | afirma | | VERIFICADO |\n",
	)
	got := strings.Join(ValidarEntrevista(texto), "\n")
	for _, esperado := range []string{"Q1: sem fato citado", "Q2: fato F9 não existe", "F1: sem prova"} {
		if !strings.Contains(got, esperado) {
			t.Fatalf("falhas %q não contêm %q", got, esperado)
		}
	}
}

func TestEntrevistaIntegra(t *testing.T) {
	texto := entrevista(
		"| Q1 | a? | F1, F2 | x |\n",
		"| F1 | afirma | arq.go:1 | VERIFICADO |\n| F2 | outra | grep \"a\\|b\" | REFUTADO para um. |\n",
	)
	if got := ValidarEntrevista(texto); len(got) != 0 {
		t.Fatalf("entrevista íntegra devolveu falhas: %v", got)
	}
}

func TestEntrevistaSoRefutadoFalha(t *testing.T) {
	texto := entrevista("| Q1 | a? | F1 | x |\n", "| F1 | afirma | arq.go:1 | REFUTADO |\n")
	got := strings.Join(ValidarEntrevista(texto), "\n")
	if !strings.Contains(got, "Q1: cita só fatos REFUTADO") {
		t.Fatalf("falhas = %q", got)
	}
}

func TestEntrevistaResultadoInvalido(t *testing.T) {
	texto := entrevista("| Q1 | a? | F1 | x |\n", "| F1 | afirma | arq.go:1 | TALVEZ |\n")
	got := strings.Join(ValidarEntrevista(texto), "\n")
	if !strings.Contains(got, "F1: resultado inválido") {
		t.Fatalf("falhas = %q", got)
	}
}

func TestEntrevistaIgnoraOutrosArtefatos(t *testing.T) {
	texto := "---\nid: SPEC-X\n---\n| Q1 | a? | | x |\n"
	if got := ValidarEntrevista(texto); got != nil {
		t.Fatalf("artefato que não é entrevista devolveu falhas: %v", got)
	}
}

package main

import (
	"strings"
	"testing"
)

const entrevistaBase = "---\nid: INTERVIEW-X\n---\n| Q1 | a? | F1 | x |\n| F1 | afirma | arq.go:1 | VERIFICADO |\n"

func TestInterviewStepIntegra(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "i.md", entrevistaBase)
	if out, err := interviewStep(dir, []string{"."}).Run(); err != nil {
		t.Fatalf("entrevista íntegra falhou: %v %s", err, out)
	}
}

func TestInterviewStepListaFalhas(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "i.md", strings.Replace(entrevistaBase, "| F1 | afirma | arq.go:1 |", "| F1 | afirma | |", 1))
	out, err := interviewStep(dir, []string{"."}).Run()
	if err == nil {
		t.Fatal("esperava falha para fato sem prova")
	}
	if !strings.Contains(out, "i.md: F1: sem prova") {
		t.Fatalf("saída %q não cita a falha", out)
	}
}

func TestInterviewStepCaminhoInexistente(t *testing.T) {
	if _, err := interviewStep(t.TempDir(), []string{"nada"}).Run(); err == nil {
		t.Fatal("esperava erro para caminho inexistente")
	}
}

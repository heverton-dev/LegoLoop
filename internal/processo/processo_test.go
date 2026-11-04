package processo_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"legoloop/internal/processo"
)

func etapa(id string, deps ...string) processo.Etapa {
	return processo.Etapa{
		ID: id, Titulo: "Etapa " + id, Tipo: "construir", Peca: "aidd-tdd", Depende: deps,
		Produz: processo.Artefato{Caminho: "out/" + id + ".md", Tipo: "markdown"},
		Gate:   "file_exists('out/" + id + ".md')",
	}
}

func valido() processo.Processo {
	return processo.Processo{
		Nome: "legoloop-build-engine", Versao: "1.0.0", Natureza: "software", Objetivo: "Construir com TDD.",
		Etapas:   []processo.Etapa{etapa("E01"), etapa("E02", "E01"), etapa("E03", "E02")},
		CriadoEm: time.Date(2026, 11, 4, 0, 0, 0, 0, time.UTC),
	}
}

func TestProcessoValidate_Valido(t *testing.T) {
	if err := valido().Validate(); err != nil {
		t.Fatalf("processo válido foi recusado: %v", err)
	}
}

func TestProcessoValidate_Recusa(t *testing.T) {
	cases := []struct {
		nome string
		msg  string
		muda func(p *processo.Processo)
	}{
		{"nome sem prefixo", "nome", func(p *processo.Processo) { p.Nome = "build-engine" }},
		{"nome com maiúscula", "nome", func(p *processo.Processo) { p.Nome = "legoloop-Build" }},
		{"nome com barra", "nome", func(p *processo.Processo) { p.Nome = "legoloop-a/b" }},
		{"nome só com prefixo", "nome", func(p *processo.Processo) { p.Nome = "legoloop-" }},
		{"versão inválida", "versão", func(p *processo.Processo) { p.Versao = "1.0" }},
		{"versão com letra", "versão", func(p *processo.Processo) { p.Versao = "v1.0.0" }},
		{"sem natureza", "natureza", func(p *processo.Processo) { p.Natureza = " " }},
		{"sem objetivo", "objetivo", func(p *processo.Processo) { p.Objetivo = "" }},
		{"sem etapas", "ao menos uma etapa", func(p *processo.Processo) { p.Etapas = nil }},
		{"id duplicado", "duplicado", func(p *processo.Processo) { p.Etapas[2].ID = "E01" }},
		{"dependência inexistente", "não existe", func(p *processo.Processo) { p.Etapas[1].Depende = []string{"E99"} }},
		{"etapa sem conexão", "sem conexão", func(p *processo.Processo) { p.Etapas = append(p.Etapas, etapa("E09")) }},
		{"ciclo", "volta a si mesma", func(p *processo.Processo) { p.Etapas[0].Depende = []string{"E03"} }},
		{"gate não compila", "não compila", func(p *processo.Processo) { p.Etapas[0].Gate = "file_exists(" }},
		{"gate vazio", "sem gate", func(p *processo.Processo) { p.Etapas[0].Gate = "" }},
		{"artefato fora da raiz", "caminho relativo", func(p *processo.Processo) { p.Etapas[0].Produz.Caminho = "../fora.md" }},
		{"artefato absoluto", "caminho relativo", func(p *processo.Processo) { p.Etapas[0].Produz.Caminho = "/etc/passwd" }},
	}
	for _, c := range cases {
		t.Run(c.nome, func(t *testing.T) {
			p := valido()
			p.Etapas = append([]processo.Etapa(nil), p.Etapas...)
			c.muda(&p)
			err := p.Validate()
			if !errors.Is(err, processo.ErrInvalido) {
				t.Fatalf("esperava ErrInvalido, veio %v", err)
			}
			if !strings.Contains(err.Error(), c.msg) {
				t.Fatalf("o erro deveria falar de %q (a regra certa), veio: %v", c.msg, err)
			}
		})
	}
}

func TestProcessoValidate_UmaEtapaSozinhaEValida(t *testing.T) {
	p := valido()
	p.Etapas = []processo.Etapa{etapa("E01")}
	if err := p.Validate(); err != nil {
		t.Fatalf("processo de uma etapa é válido: %v", err)
	}
}

func TestProcessoValidate_MensagemDizOQueFalhou(t *testing.T) {
	p := valido()
	p.Etapas[1].Depende = []string{"E99"}
	if err := p.Validate(); err == nil || !strings.Contains(err.Error(), "E99") {
		t.Fatalf("o erro deve citar a etapa problemática: %v", err)
	}
}

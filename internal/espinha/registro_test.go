package espinha

import (
	"path/filepath"
	"strings"
	"testing"
)

func registroBaseValido() Registro {
	return Registro{
		Schema: "espinha.v1",
		Verbos: []Verbo{
			{Nome: "validate", Camada: "espinha", Motor: "internal/dag", Contrato: "valida grafo"},
			{Nome: "legopress", Camada: "dominio", Motor: "internal/legopress", Contrato: "compila livros"},
		},
		Motores: []Motor{
			{Pacote: "internal/dag", Interface: "Validate", Deterministico: true},
		},
		Processos: []Processo{
			{Slug: "legoloop-fechar"},
		},
		Requisitos: []Requisito{
			{Lei: "Lei 1", Verificador: "internal/fechar", Mecanico: true},
			{Lei: "Lei 3", Mecanico: false, Justificativa: "Diretriz de concisao sem verificador estatico."},
		},
	}
}

func TestRegistroValida(t *testing.T) {
	t.Run("registro real do repositorio", func(t *testing.T) {
		caminho := filepath.Join("..", "..", "espinha", "registro.json")
		reg, err := CarregarRegistro(caminho)
		if err != nil {
			t.Fatalf("CarregarRegistro(%s) falhou: %v", caminho, err)
		}
		if err := reg.Validar(); err != nil {
			t.Fatalf("registro real invalido: %v", err)
		}
		if len(reg.Processos) != 33 {
			t.Fatalf("esperado 33 processos, obtido %d", len(reg.Processos))
		}
		for _, p := range reg.Processos {
			if !strings.HasPrefix(p.Slug, "legoloop-") {
				t.Errorf("processo sem prefixo legoloop-: %q", p.Slug)
			}
		}
	})

	t.Run("recusa schema invalido", func(t *testing.T) {
		reg := registroBaseValido()
		reg.Schema = "espinha.v0"
		if err := reg.Validar(); err == nil {
			t.Fatal("esperado erro para schema invalido")
		}
	})

	t.Run("recusa listas vazias", func(t *testing.T) {
		casos := []struct {
			nome  string
			mutar func(*Registro)
		}{
			{"verbos vazio", func(r *Registro) { r.Verbos = nil }},
			{"motores vazio", func(r *Registro) { r.Motores = nil }},
			{"processos vazio", func(r *Registro) { r.Processos = nil }},
			{"requisitos vazio", func(r *Registro) { r.Requisitos = nil }},
		}
		for _, c := range casos {
			reg := registroBaseValido()
			c.mutar(&reg)
			if err := reg.Validar(); err == nil {
				t.Errorf("%s: esperado erro de validacao", c.nome)
			}
		}
	})

	t.Run("recusa nomes repetidos", func(t *testing.T) {
		casos := []struct {
			nome  string
			mutar func(*Registro)
		}{
			{"verbo repetido", func(r *Registro) {
				r.Verbos = append(r.Verbos, Verbo{Nome: "validate", Camada: "espinha", Motor: "internal/dag", Contrato: "dup"})
			}},
			{"motor repetido", func(r *Registro) {
				r.Motores = append(r.Motores, Motor{Pacote: "internal/dag", Interface: "Dup", Deterministico: true})
			}},
			{"processo repetido", func(r *Registro) {
				r.Processos = append(r.Processos, Processo{Slug: "legoloop-fechar"})
			}},
			{"requisito repetido", func(r *Registro) {
				r.Requisitos = append(r.Requisitos, Requisito{Lei: "Lei 1", Verificador: "x", Mecanico: true})
			}},
		}
		for _, c := range casos {
			reg := registroBaseValido()
			c.mutar(&reg)
			if err := reg.Validar(); err == nil {
				t.Errorf("%s: esperado erro de duplicata", c.nome)
			}
		}
	})

	t.Run("recusa camada invalida", func(t *testing.T) {
		reg := registroBaseValido()
		reg.Verbos[0].Camada = "nucleo"
		if err := reg.Validar(); err == nil {
			t.Fatal("esperado erro para camada invalida")
		}
	})

	t.Run("recusa requisito sem verificador ou justificativa", func(t *testing.T) {
		regMec := registroBaseValido()
		regMec.Requisitos[0].Verificador = ""
		if err := regMec.Validar(); err == nil {
			t.Error("esperado erro para requisito mecanico sem verificador")
		}
		regNaoMec := registroBaseValido()
		regNaoMec.Requisitos[1].Justificativa = ""
		if err := regNaoMec.Validar(); err == nil {
			t.Error("esperado erro para requisito nao mecanico sem justificativa")
		}
	})
}

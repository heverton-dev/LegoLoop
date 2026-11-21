package espinha

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestTodaLeiTemVerificador(t *testing.T) {
	agentsPath := filepath.Join("..", "..", "AGENTS.md")
	dadosAgents, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Fatalf("ler AGENTS.md: %v", err)
	}
	regPath := filepath.Join("..", "..", "espinha", "registro.json")
	reg, err := CarregarRegistro(regPath)
	if err != nil {
		t.Fatalf("carregar registro: %v", err)
	}

	leis := ExtrairLeisAgents(string(dadosAgents))
	esperadas := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 16, 17, 18}
	if !reflect.DeepEqual(leis, esperadas) {
		t.Fatalf("leis extraidas = %v, esperado %v", leis, esperadas)
	}

	t.Run("registro real cobre todas as leis do AGENTS.md", func(t *testing.T) {
		if err := ConferirLeis(leis, reg); err != nil {
			t.Fatalf("ConferirLeis falhou no repositorio real: %v", err)
		}
	})

	t.Run("extracao ignora sublistas indentadas e outras secoes", func(t *testing.T) {
		md := "# Titulo\n## 1. Leis Invioláveis\n\n1. **A:** texto\n   1. subitem\n2. **B:** texto\n\n## 2. Outra Secao\n3. **C:** fora da secao 1\n"
		got := ExtrairLeisAgents(md)
		want := []int{1, 2}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("ExtrairLeisAgents = %v, esperado %v", got, want)
		}
	})

	t.Run("mutacao omite lei do registro falha", func(t *testing.T) {
		mutado := reg
		mutado.Requisitos = append([]Requisito(nil), reg.Requisitos[:len(reg.Requisitos)-1]...)
		if err := ConferirLeis(leis, mutado); err == nil {
			t.Fatal("esperado erro ao omitir lei do registro")
		}
	})

	t.Run("mutacao remove verificador de lei mecanica falha", func(t *testing.T) {
		mutado := clonarRequisitos(reg)
		alterou := false
		for i := range mutado.Requisitos {
			if mutado.Requisitos[i].Mecanico {
				mutado.Requisitos[i].Verificador = ""
				alterou = true
				break
			}
		}
		if !alterou {
			t.Fatal("nenhuma lei mecanica encontrada para mutacao")
		}
		if err := ConferirLeis(leis, mutado); err == nil {
			t.Fatal("esperado erro ao remover verificador de lei mecanica")
		}
	})

	t.Run("mutacao remove justificativa de lei nao mecanica falha", func(t *testing.T) {
		mutado := clonarRequisitos(reg)
		alterou := false
		for i := range mutado.Requisitos {
			if !mutado.Requisitos[i].Mecanico {
				mutado.Requisitos[i].Justificativa = ""
				alterou = true
				break
			}
		}
		if !alterou {
			t.Fatal("nenhuma lei nao mecanica encontrada para mutacao")
		}
		if err := ConferirLeis(leis, mutado); err == nil {
			t.Fatal("esperado erro ao remover justificativa de lei nao mecanica")
		}
	})
}

func clonarRequisitos(reg Registro) Registro {
	copia := reg
	copia.Requisitos = append([]Requisito(nil), reg.Requisitos...)
	return copia
}

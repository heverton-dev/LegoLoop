package roteamento

import (
	"errors"
	"testing"
)

func perfilDeTeste() Perfil {
	return Perfil{
		"micro":   {{"claude", "sonnet", "low"}, {"claude", "haiku", "low"}},
		"extreme": {{"claude", "opus", "high"}, {"claude", "opus2", "high"}},
	}
}

func invDeTeste(modelos ...string) Inventario {
	return Inventario{SO: "linux", Harnesses: []HarnessInv{{Nome: "claude", Binario: "claude", Modelos: modelos}}}
}

func TestDegradacaoRegistrada(t *testing.T) {
	res, err := Resolver("N1", "micro", perfilDeTeste(), invDeTeste("haiku"))
	if err != nil {
		t.Fatal(err)
	}
	if res.Usado.Modelo != "haiku" || res.StopGate {
		t.Fatalf("usado inesperado: %+v", res)
	}
	if len(res.Substituicoes) != 1 {
		t.Fatalf("esperava 1 substituicao, veio %d", len(res.Substituicoes))
	}
	s := res.Substituicoes[0]
	if s.No != "N1" || s.Pedido.Modelo != "sonnet" || s.Usado.Modelo != "haiku" {
		t.Fatalf("substituicao inesperada: %+v", s)
	}
}

func TestPrimeiroCandidatoSemSubstituicao(t *testing.T) {
	res, err := Resolver("N1", "micro", perfilDeTeste(), invDeTeste("sonnet", "haiku"))
	if err != nil {
		t.Fatal(err)
	}
	if res.Usado.Modelo != "sonnet" || len(res.Substituicoes) != 0 {
		t.Fatalf("resultado inesperado: %+v", res)
	}
}

func TestExtremeSemModeloViraGate(t *testing.T) {
	res, err := Resolver("N2", "extreme", perfilDeTeste(), invDeTeste("haiku"))
	if err != nil {
		t.Fatal(err)
	}
	if !res.StopGate || len(res.Substituicoes) != 0 {
		t.Fatalf("esperava stop_gate sem degradar: %+v", res)
	}
}

func TestInventarioVazioRetornaErro(t *testing.T) {
	_, err := Resolver("N3", "micro", perfilDeTeste(), Inventario{})
	if !errors.Is(err, ErrInventarioVazio) {
		t.Fatalf("erro inesperado: %v", err)
	}
}

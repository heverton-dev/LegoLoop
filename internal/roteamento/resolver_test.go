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

func TestPoliticaDeFalhaPorClasse(t *testing.T) {
	casos := map[string]string{
		ClasseNano:    PoliticaDegrada,
		ClasseMicro:   PoliticaDegrada,
		ClasseGrande:  PoliticaFalha,
		ClasseExtreme: PoliticaFalha,
	}
	for classe, esperado := range casos {
		if got := PoliticaDeFalha(classe); got != esperado {
			t.Fatalf("classe %s: esperava %s, veio %s", classe, esperado, got)
		}
	}
}

func invViaDeTeste(harness string, via Via, modelo string) []ItemInv {
	return []ItemInv{{
		Par:     Par{Harness: harness, Via: via},
		Modelos: []string{modelo},
		Estado:  EstadoHomologado,
	}}
}

func TestResolverViaSegundoCandidato(t *testing.T) {
	casos := []struct {
		nome      string
		elegiveis []CandidatoVia
		inv       []ItemInv
		usado     CandidatoVia
	}{
		{
			nome: "mesma via",
			elegiveis: []CandidatoVia{
				{Harness: "claude", Via: ViaOficial, Modelo: "sonnet"},
				{Harness: "claude", Via: ViaOficial, Modelo: "haiku"},
			},
			inv:   invViaDeTeste("claude", ViaOficial, "haiku"),
			usado: CandidatoVia{Harness: "claude", Via: ViaOficial, Modelo: "haiku"},
		},
		{
			nome: "via diferente",
			elegiveis: []CandidatoVia{
				{Harness: "claude", Via: ViaOficial, Modelo: "sonnet"},
				{Harness: "claude", Via: Via9Router, Modelo: "sonnet"},
			},
			inv:   invViaDeTeste("claude", Via9Router, "sonnet"),
			usado: CandidatoVia{Harness: "claude", Via: Via9Router, Modelo: "sonnet"},
		},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			res, err := ResolverVia("N1", ClasseMicro, map[string][]CandidatoVia{ClasseMicro: c.elegiveis}, c.inv)
			if err != nil {
				t.Fatal(err)
			}
			if res.Usado != c.usado || res.Politica != "" || res.StopGate {
				t.Fatalf("resultado inesperado: %+v", res)
			}
			if len(res.Substituicoes) != 1 {
				t.Fatalf("esperava 1 substituicao, veio %d", len(res.Substituicoes))
			}
			s := res.Substituicoes[0]
			if s.No != "N1" || s.Pedido != c.elegiveis[0] || s.Usado != c.usado {
				t.Fatalf("substituicao inesperada: %+v", s)
			}
		})
	}
}

func TestResolverViaSemCandidato(t *testing.T) {
	inv := invViaDeTeste("claude", ViaOficial, "outro")
	cand := []CandidatoVia{{Harness: "claude", Via: ViaOficial, Modelo: "sonnet"}}
	elegiveis := map[string][]CandidatoVia{ClasseMicro: cand, ClasseGrande: cand, ClasseExtreme: cand}
	casos := []struct {
		classe   string
		politica string
		gate     bool
	}{
		{ClasseMicro, PoliticaDegrada, false},
		{ClasseGrande, PoliticaFalha, false},
		{ClasseExtreme, PoliticaFalha, true},
	}
	for _, c := range casos {
		res, err := ResolverVia("N2", c.classe, elegiveis, inv)
		if err != nil {
			t.Fatal(err)
		}
		if res.Politica != c.politica || res.StopGate != c.gate || len(res.Substituicoes) != 0 {
			t.Fatalf("classe %s: resultado inesperado: %+v", c.classe, res)
		}
	}
}

func TestResolverViaInventarioVazio(t *testing.T) {
	_, err := ResolverVia("N3", ClasseMicro, nil, nil)
	if !errors.Is(err, ErrInventarioVazio) {
		t.Fatalf("erro inesperado: %v", err)
	}
}

func TestInventarioVazioRetornaErro(t *testing.T) {
	_, err := Resolver("N3", "micro", perfilDeTeste(), Inventario{})
	if !errors.Is(err, ErrInventarioVazio) {
		t.Fatalf("erro inesperado: %v", err)
	}
}

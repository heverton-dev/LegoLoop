package refinar_test

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/refinar"
)

func adicionar(m refinar.Mudanca) refinar.Request {
	m.Adicionar = true
	return refinar.Request{Nome: nome, Mudanca: m}
}

func TestRefinar_AdicionarEtapaNova(t *testing.T) {
	c := semear(t, processoDeTeste(), inicio, true)
	m := refinar.Mudanca{Etapa: "e04", Peca: "legoloop-transcrever", Titulo: "Transcrever", Casos: []string{"transcricao"}, Depende: []string{"e02"}}
	res, err := refinar.Run(context.Background(), adicionar(m), depsDe(c, nil, seguir, inicio))
	if err != nil {
		t.Fatalf("adicionar: %v", err)
	}
	if res.Para != "1.1.0" {
		t.Errorf("versão = %s; quer 1.1.0", res.Para)
	}
	p := lerProcesso(t, c)
	if len(p.Etapas) != 4 || p.Etapas[3].ID != "e04" || p.Etapas[3].Peca != "legoloop-transcrever" {
		t.Fatalf("etapa nova ausente: %+v", p.Etapas)
	}
	if !slices.Equal(p.Etapas[3].Depende, []string{"e02"}) || p.Etapas[3].Gate == "" || p.Etapas[3].Produz.Caminho == "" {
		t.Errorf("etapa nova sem ligação, gate ou artefato: %+v", p.Etapas[3])
	}
	if !slices.Equal(p.Etapas[3].Consome, []string{p.Etapas[1].Produz.Caminho}) {
		t.Errorf("a etapa nova deve consumir o artefato de e02: %v", p.Etapas[3].Consome)
	}
}

func TestRefinar_AdicionarSemDependeLigaNaUltimaEtapa(t *testing.T) {
	c := semear(t, processoDeTeste(), inicio, true)
	m := refinar.Mudanca{Etapa: "e04", Peca: "legoloop-x2", Casos: []string{"a"}}
	if _, err := refinar.Run(context.Background(), adicionar(m), depsDe(c, nil, seguir, inicio)); err != nil {
		t.Fatal(err)
	}
	if got := lerProcesso(t, c).Etapas[3].Depende; !slices.Equal(got, []string{"e03"}) {
		t.Errorf("depende = %v; quer [e03]", got)
	}
}

func TestRefinar_AdicionarRecusaIdExistenteEDependeDesconhecido(t *testing.T) {
	for rotulo, m := range map[string]refinar.Mudanca{
		"id existente":         {Etapa: "e02", Peca: "p", Casos: []string{"a"}},
		"depende desconhecido": {Etapa: "e04", Peca: "p", Casos: []string{"a"}, Depende: []string{"e99"}},
		"sem peça":             {Etapa: "e04", Casos: []string{"a"}},
		"sem caso":             {Etapa: "e04", Peca: "p"},
	} {
		c := semear(t, processoDeTeste(), inicio, true)
		antes := arvore(t, c)
		_, err := refinar.Run(context.Background(), adicionar(m), depsDe(c, nil, seguir, inicio))
		if !errors.Is(err, refinar.ErrMudanca) {
			t.Errorf("%s: erro = %v; quer ErrMudanca", rotulo, err)
		}
		igual(t, antes, arvore(t, c))
	}
}

func TestRefinar_SemAdicionarEtapaNovaContinuaRecusada(t *testing.T) {
	c := semear(t, processoDeTeste(), inicio, true)
	m := refinar.Mudanca{Etapa: "e99", Peca: "p", Casos: []string{"a"}}
	_, err := refinar.Run(context.Background(), refinar.Request{Nome: nome, Mudanca: m}, depsDe(c, nil, seguir, inicio))
	if !errors.Is(err, refinar.ErrMudanca) {
		t.Fatalf("erro = %v; quer ErrMudanca", err)
	}
}

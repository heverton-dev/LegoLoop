package refinar_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/caminhos"
	"github.com/heverton-dev/LegoLoop/internal/processo"
	"github.com/heverton-dev/LegoLoop/internal/refinar"
	"github.com/heverton-dev/LegoLoop/internal/triagem"
)

func umaEtapa() processo.Processo {
	p := processoDeTeste()
	p.Etapas = p.Etapas[:1]
	return p
}

func lerTriagem(t *testing.T, c cofre) (triagem.TriagemSpec, error) {
	t.Helper()
	var s triagem.TriagemSpec
	b, err := os.ReadFile(filepath.Join(c.vault, caminhos.DirProcessos, nome, "triagem.json"))
	if err != nil {
		return s, err
	}
	return s, json.Unmarshal(b, &s)
}

func novaEtapa() refinar.Mudanca {
	return refinar.Mudanca{Etapa: "e02", Peca: "legoloop-nova", Casos: []string{"a"}, Adicionar: true}
}

func idsDaTriagem(s triagem.TriagemSpec) []string {
	var ids []string
	for _, e := range s.Etapas {
		ids = append(ids, e.ID)
	}
	return ids
}

func TestRefinar_AdicionarQueFazDoisNosCriaATriagem(t *testing.T) {
	c := semear(t, umaEtapa(), inicio, true)
	canonizar(t, c)
	if _, err := lerTriagem(t, c); err == nil {
		t.Fatal("peça de um nó não tem triagem")
	}
	if _, err := refinar.Run(context.Background(), refinar.Request{Nome: nome, Mudanca: novaEtapa()}, depsDe(c, nil, seguir, inicio)); err != nil {
		t.Fatal(err)
	}
	s, err := lerTriagem(t, c)
	if err != nil {
		t.Fatalf("com dois nós a triagem deve existir: %v", err)
	}
	if err := s.Validar(); err != nil || !slices.Equal(idsDaTriagem(s), []string{"e01", "e02"}) {
		t.Errorf("triagem inválida ou incompleta: %v %v", err, idsDaTriagem(s))
	}
}

func TestRefinar_TriagemExistenteGanhaSoAEtapaNovaEMantemAEdicao(t *testing.T) {
	c := semear(t, umaEtapa(), inicio, true)
	canonizar(t, c)
	ant := triagem.GerarTriagemPadrao(nome, []string{"e01"})
	ant.Etapas[0].PalavrasChave = append(ant.Etapas[0].PalavrasChave, "palavra-do-operador")
	b, _ := json.Marshal(ant)
	if err := os.WriteFile(filepath.Join(c.vault, caminhos.DirProcessos, nome, "triagem.json"), b, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := refinar.Run(context.Background(), refinar.Request{Nome: nome, Mudanca: novaEtapa()}, depsDe(c, nil, seguir, inicio)); err != nil {
		t.Fatal(err)
	}
	s, _ := lerTriagem(t, c)
	if !slices.Equal(idsDaTriagem(s), []string{"e01", "e02"}) || !slices.Contains(s.Etapas[0].PalavrasChave, "palavra-do-operador") {
		t.Errorf("deve completar sem perder a edição: %+v", s.Etapas)
	}
	if !slices.Equal(s.CriterioAmbiguidade.Etapas, []string{"e01", "e02"}) {
		t.Errorf("critério de ambiguidade sem a etapa nova: %v", s.CriterioAmbiguidade.Etapas)
	}
}

func TestRefinar_MudancaDeUmNoNaoCriaTriagem(t *testing.T) {
	c := semear(t, umaEtapa(), inicio, true)
	canonizar(t, c)
	m := refinar.Mudanca{Etapa: "e01", Peca: "aidd-x", Casos: []string{"a"}}
	if _, err := refinar.Run(context.Background(), refinar.Request{Nome: nome, Mudanca: m}, depsDe(c, nil, seguir, inicio)); err != nil {
		t.Fatal(err)
	}
	if _, err := lerTriagem(t, c); err == nil {
		t.Error("um nó não precisa de triagem")
	}
}

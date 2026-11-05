package matriz

import (
	"strings"
	"testing"
)

func TestCalibragemRank(t *testing.T) {
	n := Needs{Natureza: "dados", Etapas: []string{"coletar"}, Casos: []string{"coleta-de-dados"}}
	generica := entrada("super-qa", 0.95, perfil([]string{"coletar"}, []string{"qualquer"}))
	coletor := entrada("coletor", 0.30, perfil([]string{"coletar"}, []string{"dados"}, "coleta-de-dados"))
	r := Rank(n, []Entrada{generica, coletor}, 0)
	if g := achar(t, r, "super-qa"); !g.Excluida || !strings.HasPrefix(g.Motivo, "caso de uso:") || !strings.Contains(g.Motivo, "coleta-de-dados") {
		t.Errorf("com casos pedidos, aptidão 0 é exclusão com motivo: %+v", g)
	}
	if g := achar(t, r, "super-qa"); g.ExcluidaPorUso() {
		t.Errorf("não cobrir o caso não é proibição de uso (não vira aviso): %+v", g)
	}
	if c := achar(t, r, "coletor"); c.Excluida || c.Aptidao != 1 {
		t.Errorf("coletor = %+v", c)
	}
	if r[0].Entrada.Ref != "coletor" {
		t.Errorf("a escolhida vem primeiro: %s", r[0].Entrada.Ref)
	}
	semCasos := Rank(Needs{Natureza: "dados", Etapas: []string{"coletar"}}, []Entrada{generica}, 0)
	if g := achar(t, semCasos, "super-qa"); g.Excluida || g.Aptidao != 1 {
		t.Errorf("sem casos pedidos a etapa basta: %+v", g)
	}
}

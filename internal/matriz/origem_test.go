package matriz

import (
	"slices"
	"strings"
	"testing"
)

func perfilComOrigem(origem string) *PerfilDeUso {
	p := perfil([]string{"construir"}, []string{"software"}, "tdd")
	p.Origem = origem
	return p
}

var precisaConstruirTDD = Needs{Natureza: "software", Etapas: []string{"construir"}, Casos: []string{"tdd"}}

// TestOrigemDoPerfil prova o AC-E03.8: no empate de aptidão, declarado ou confirmado vence inferido antes do
// Score; o motivo cita a origem; a aptidão não muda; fora do empate nada muda.
func TestOrigemDoPerfil(t *testing.T) {
	t.Run("declarada vence inferida com Score maior", func(t *testing.T) {
		r := Rank(precisaConstruirTDD, []Entrada{
			entrada("inferida", 0.95, perfilComOrigem(OrigemInferido)),
			entrada("declarada", 0.60, perfilComOrigem(OrigemDeclarado)),
		}, 0)
		if got := escolhidas(r); !slices.Equal(got, []string{"declarada", "inferida"}) {
			t.Fatalf("ordem = %v; a declarada vence no empate de aptidão", got)
		}
		if m := r[0].Motivo; !strings.Contains(m, "perfil: declarado") {
			t.Errorf("o motivo deve citar a origem: %q", m)
		}
		if m := r[1].Motivo; !strings.Contains(m, "perfil: inferido") {
			t.Errorf("o motivo deve citar a origem: %q", m)
		}
		if r[0].Aptidao != r[1].Aptidao || r[0].Aptidao != 1 {
			t.Errorf("a aptidão não muda com a origem: %v e %v", r[0].Aptidao, r[1].Aptidao)
		}
	})
	t.Run("confirmada vence inferida", func(t *testing.T) {
		r := Rank(precisaConstruirTDD, []Entrada{
			entrada("inferida", 0.95, perfilComOrigem(OrigemInferido)),
			entrada("confirmada", 0.60, perfilComOrigem(OrigemConfirmado)),
		}, 0)
		if got := escolhidas(r); got[0] != "confirmada" || !strings.Contains(r[0].Motivo, "perfil: confirmado") {
			t.Errorf("a confirmada vence: %v %q", got, r[0].Motivo)
		}
	})
	t.Run("só a inferida: é escolhida e o motivo cita inferido", func(t *testing.T) {
		r := Rank(precisaConstruirTDD, []Entrada{entrada("inferida", 0.95, perfilComOrigem(OrigemInferido))}, 0)
		if len(r) != 1 || r[0].Excluida || !strings.Contains(r[0].Motivo, "perfil: inferido") {
			t.Errorf("a inferida sozinha é escolhida: %+v", r)
		}
	})
	t.Run("mesma origem: o Score decide", func(t *testing.T) {
		r := Rank(precisaConstruirTDD, []Entrada{
			entrada("baixa", 0.60, perfilComOrigem(OrigemInferido)),
			entrada("alta", 0.95, perfilComOrigem(OrigemInferido)),
		}, 0)
		if got := escolhidas(r); got[0] != "alta" {
			t.Errorf("ordem = %v", got)
		}
	})
	t.Run("sem empate de aptidão: inferida mais apta vence a declarada", func(t *testing.T) {
		inf := perfilComOrigem(OrigemInferido)
		dec := perfilComOrigem(OrigemDeclarado)
		dec.CasosDeUso = []string{"verificacao-de-qualidade"}
		n := precisaConstruirTDD
		n.Casos = []string{"tdd", "verificacao-de-qualidade"}
		inf.CasosDeUso = []string{"tdd", "verificacao-de-qualidade"}
		r := Rank(n, []Entrada{entrada("declarada", 0.99, dec), entrada("inferida", 0.10, inf)}, 0)
		if got := escolhidas(r); got[0] != "inferida" {
			t.Errorf("a aptidão manda: %v", got)
		}
	})
	t.Run("peça sem perfil nunca é escolhida", func(t *testing.T) {
		r := Rank(precisaConstruirTDD, []Entrada{entrada("sem", 0.99, nil)}, 0)
		if len(r) != 0 {
			t.Errorf("sem perfil fica fora: %+v", r)
		}
	})
}

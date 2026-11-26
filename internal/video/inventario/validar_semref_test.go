package inventario

import (
	"errors"
	"testing"
	"time"
)

// TestValidarEstiloSemReferenciasNoDisco cobre o estilo editado a mao: sem referencias nao valida.
func TestValidarEstiloSemReferenciasNoDisco(t *testing.T) {
	s := novoStoreComEstilo(t, "sem-ref")
	e, err := s.LerEstilo("sem-ref")
	if err != nil {
		t.Fatal(err)
	}
	e.Referencias = nil
	if err := gravarJSON(s.caminhoEstilo("sem-ref"), e); err != nil {
		t.Fatal(err)
	}
	arq := gravarAprovacao(t, t.TempDir(), Aprovacao{Estilo: "sem-ref", AprovadoPor: "usuario", PersonagensTeste: []string{"p1"}})
	if _, err := s.ValidarEstilo("sem-ref", arq, time.Now()); !errors.Is(err, ErrEstiloInvalido) {
		t.Fatalf("err = %v, quer ErrEstiloInvalido", err)
	}
	if got, _ := s.LerEstilo("sem-ref"); got.Status != "rascunho" {
		t.Fatalf("status = %q, quer rascunho", got.Status)
	}
}

package inventario

import (
	"errors"
	"testing"
)

// TestFaixaCincoEtapas cobre o vocabulario de 5 faixas etarias: bebe, crianca, adolescente, adulto, idoso.
func TestFaixaCincoEtapas(t *testing.T) {
	for _, f := range []string{"bebe", "crianca", "adolescente", "adulto", "idoso"} {
		if err := ValidarFaixa(f); err != nil {
			t.Fatalf("ValidarFaixa(%q) = %v, quer nil", f, err)
		}
	}
	for _, f := range []string{"", "jovem", "Adulto", "x"} {
		if err := ValidarFaixa(f); !errors.Is(err, ErrFaixaInvalida) {
			t.Fatalf("ValidarFaixa(%q) = %v, quer ErrFaixaInvalida", f, err)
		}
	}
}

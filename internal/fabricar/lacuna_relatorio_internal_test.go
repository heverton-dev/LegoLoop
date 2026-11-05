package fabricar

import (
	"slices"
	"testing"

	"legoloop/internal/processo"
)

// TestLacunaRelatorioEtapaComposta: etapa que compõe outro processo não é lacuna, mesmo sem peça; só a etapa sem
// peça e sem composição entra no relatório e no cartão da lacuna.
func TestLacunaRelatorioEtapaComposta(t *testing.T) {
	p := processo.Processo{Etapas: []processo.Etapa{
		{Tipo: "coletar", Compoe: "legoloop-coletor"},
		{Tipo: "verificar"},
		{Tipo: "construir", Peca: "alfa"},
	}}
	if got := lacunasDe(p); !slices.Equal(got, []string{"verificar"}) {
		t.Errorf("só verificar é lacuna: %v", got)
	}
}

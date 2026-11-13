package combo_test

import (
	"errors"
	"testing"

	"legoloop/internal/combo"
)

func TestValidarRequisicaoCombo(t *testing.T) {
	// AC-I7-01.2: t-micro-8k com 9000 tokens retorna ErrTetoTokensExcedido; com 4000 tokens retorna nil.
	t.Run("teto_tokens_excedido", func(t *testing.T) {
		err := combo.ValidarRequisicao("t-micro-8k", 9000, false)
		if !errors.Is(err, combo.ErrTetoTokensExcedido) {
			t.Fatalf("esperava ErrTetoTokensExcedido, obteve %v", err)
		}
	})

	t.Run("teto_tokens_aceito", func(t *testing.T) {
		err := combo.ValidarRequisicao("t-micro-8k", 4000, false)
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
	})

	// AC-I7-01.3: t-codigo-128k com exigeFerramentas=true retorna ErrFerramentaNaoSuportada.
	t.Run("ferramentas_nao_suportadas", func(t *testing.T) {
		err := combo.ValidarRequisicao("t-codigo-128k", 50000, true)
		if !errors.Is(err, combo.ErrFerramentaNaoSuportada) {
			t.Fatalf("esperava ErrFerramentaNaoSuportada, obteve %v", err)
		}
	})

	// t-agente-128k com 50000 tokens e exigeFerramentas=true retorna nil.
	t.Run("agente_suporta_ferramentas", func(t *testing.T) {
		err := combo.ValidarRequisicao("t-agente-128k", 50000, true)
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
	})

	// combo desconhecido retorna ErrComboInvalido.
	t.Run("combo_desconhecido", func(t *testing.T) {
		err := combo.ValidarRequisicao("combo-inexistente", 1000, false)
		if !errors.Is(err, combo.ErrComboInvalido) {
			t.Fatalf("esperava ErrComboInvalido, obteve %v", err)
		}
	})
}

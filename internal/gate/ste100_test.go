package gate_test

import (
	"testing"

	"legoloop/internal/gate"
)

func TestSTE100Linter(t *testing.T) {
	t.Run("AceitaFrasesCurtasSemHedging", func(t *testing.T) {
		text := "Execute o teste com o comando go test. Verifique a saída do terminal. Atualize o arquivo de rota."
		violations, ok := gate.ValidateSTE100(text)
		if !ok {
			t.Fatalf("esperava texto válido, mas obteve violações: %v", violations)
		}
	})

	t.Run("RejeitaFrasesComMaisDe20Palavras", func(t *testing.T) {
		// 22 palavras
		longText := "Este é um procedimento muito longo que tem mais de vinte palavras com o objetivo de testar o linter sintático agora mesmo aqui."
		violations, ok := gate.ValidateSTE100(longText)
		if ok {
			t.Fatal("esperava violação de limite de palavras, mas passou")
		}
		if len(violations) == 0 {
			t.Fatal("esperava ao menos 1 violação listada")
		}
	})

	t.Run("RejeitaTermosDeHedging", func(t *testing.T) {
		hedgingText := "Talvez o teste passe se você rodar o comando."
		violations, ok := gate.ValidateSTE100(hedgingText)
		if ok {
			t.Fatal("esperava violação de hedging, mas passou")
		}
		if len(violations) == 0 {
			t.Fatal("esperava violações listadas para hedging")
		}
	})
}

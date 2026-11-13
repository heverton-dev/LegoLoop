package combo

// ValidarRequisicao valida se os parametros da chamada cabem e sao suportados pelo combo.
func ValidarRequisicao(nomeCombo string, tokensEstimados int, exigeFerramentas bool) error {
	c, err := ObterCombo(nomeCombo)
	if err != nil {
		return ErrComboInvalido
	}
	if tokensEstimados > c.MaxTokens {
		return ErrTetoTokensExcedido
	}
	if exigeFerramentas && !c.ExigeFerramentas {
		if !temSuporteFuncao(c.Membros) {
			return ErrFerramentaNaoSuportada
		}
	}
	return nil
}

func temSuporteFuncao(membros []MembroCombo) bool {
	for _, m := range membros {
		if m.SuportaFuncao {
			return true
		}
	}
	return false
}

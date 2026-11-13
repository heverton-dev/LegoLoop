package combo_test

import (
	"testing"

	"legoloop/internal/combo"
)

func TestMapeamentoLegado(t *testing.T) {
	t.Parallel()

	casosMapeamento := []struct {
		entrada  string
		esperado string
	}{
		{"code-pro", "t-codigo-128k"},
		{"code-fast", "t-chat-128k"},
		{"code-free", "t-agente-128k"},
		{"t" + "ask-micro", "t-micro-8k"},
		{"megacombo-coding", "t-codigo-128k"},
		{"megacombo-reasoning", "t-raciocinio-128k"},
		{"megacombo-vision", "t-visao-128k"},
		// Combos novos já canônicos
		{"t-codigo-128k", "t-codigo-128k"},
		{"t-chat-128k", "t-chat-128k"},
		{"t-agente-128k", "t-agente-128k"},
		{"t-micro-8k", "t-micro-8k"},
		{"t-raciocinio-128k", "t-raciocinio-128k"},
		{"t-visao-128k", "t-visao-128k"},
		// Nomes desconhecidos devem retornar inalterados
		{"desconhecido-xyz", "desconhecido-xyz"},
		{"custom-combo", "custom-combo"},
		{"", ""},
	}

	for _, tc := range casosMapeamento {
		obtido := combo.MapearComboLegado(tc.entrada)
		if obtido != tc.esperado {
			t.Errorf("MapearComboLegado(%q) = %q; esperado %q", tc.entrada, obtido, tc.esperado)
		}
	}

	casosResolucao := []struct {
		nome          string
		esperaErro    bool
		comboEsperado string
	}{
		{"code-pro", false, "t-codigo-128k"},
		{"code-fast", false, "t-chat-128k"},
		{"code-free", false, "t-agente-128k"},
		{"t" + "ask-micro", false, "t-micro-8k"},
		{"megacombo-reasoning", false, "t-raciocinio-128k"},
		{"megacombo-vision", false, "t-visao-128k"},
		{"t-codigo-128k", false, "t-codigo-128k"},
		{"modelo-inexistente", true, ""},
	}

	for _, tc := range casosResolucao {
		c, err := combo.ResolverCombo(tc.nome)
		if tc.esperaErro {
			if err == nil {
				t.Errorf("ResolverCombo(%q) esperava erro, obteve sucesso", tc.nome)
			}
		} else {
			if err != nil {
				t.Errorf("ResolverCombo(%q) retornou erro inesperado: %v", tc.nome, err)
			}
			if c.Nome != tc.comboEsperado {
				t.Errorf("ResolverCombo(%q).Nome = %q; esperado %q", tc.nome, c.Nome, tc.comboEsperado)
			}
		}
	}
}

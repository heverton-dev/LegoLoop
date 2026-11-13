package combo_test

import (
	"testing"

	"legoloop/internal/combo"
)

func TestCatalogoHonesto_SeisCombos(t *testing.T) {
	combos := combo.CatalogoHonesto()
	if len(combos) != 6 {
		t.Fatalf("esperava 6 combos, obteve %d", len(combos))
	}

	nomesEsperados := []string{
		"t-agente-128k",
		"t-chat-128k",
		"t-raciocinio-128k",
		"t-visao-128k",
		"t-codigo-128k",
		"t-micro-8k",
	}

	for _, nome := range nomesEsperados {
		c, ok := combos[nome]
		if !ok {
			t.Errorf("combo %q ausente no catalogo", nome)
			continue
		}
		if err := c.Validate(); err != nil {
			t.Errorf("combo %q invalido: %v", nome, err)
		}
	}
}

func TestCatalogoHonesto_LimitesEPropriedades(t *testing.T) {
	combos := combo.CatalogoHonesto()

	agente := combos["t-agente-128k"]
	if agente.MaxTokens != 131072 || !agente.ExigeFerramentas || len(agente.Membros) != 6 {
		t.Errorf("t-agente-128k propriedades invalidas: %+v", agente)
	}

	micro := combos["t-micro-8k"]
	if micro.MaxTokens != 8192 || micro.ExigeFerramentas || len(micro.Membros) != 2 {
		t.Errorf("t-micro-8k propriedades invalidas: %+v", micro)
	}

	codigo := combos["t-codigo-128k"]
	if codigo.MaxTokens != 131072 || len(codigo.Membros) != 4 {
		t.Errorf("t-codigo-128k propriedades invalidas: %+v", codigo)
	}

	visao := combos["t-visao-128k"]
	if !visao.ExigeImagem || len(visao.Membros) != 6 {
		t.Errorf("t-visao-128k propriedades invalidas: %+v", visao)
	}
}

func TestObterComboEListar(t *testing.T) {
	lista := combo.ListarCombos()
	if len(lista) != 6 {
		t.Fatalf("esperava 6 nomes listados, obteve %d", len(lista))
	}

	c, err := combo.ObterCombo("t-agente-128k")
	if err != nil || c.Nome != "t-agente-128k" {
		t.Errorf("falha ao obter t-agente-128k: %v", err)
	}

	_, errInvalido := combo.ObterCombo("combo-inexistente")
	if errInvalido != combo.ErrComboInvalido {
		t.Errorf("esperava ErrComboInvalido, obteve %v", errInvalido)
	}
}

func TestComboValidate_Invalido(t *testing.T) {
	invalidoSemNome := combo.ComboHonesto{
		Nome:      "",
		MaxTokens: 1000,
		Membros:   []combo.MembroCombo{{ID: "m1"}},
	}
	if err := invalidoSemNome.Validate(); err != combo.ErrComboInvalido {
		t.Errorf("esperava ErrComboInvalido para nome vazio, obteve %v", err)
	}

	invalidoSemMembros := combo.ComboHonesto{
		Nome:      "t-teste",
		MaxTokens: 1000,
		Membros:   nil,
	}
	if err := invalidoSemMembros.Validate(); err != combo.ErrComboInvalido {
		t.Errorf("esperava ErrComboInvalido para membros vazios, obteve %v", err)
	}

	invalidoSemTokens := combo.ComboHonesto{
		Nome:      "t-teste",
		MaxTokens: 0,
		Membros:   []combo.MembroCombo{{ID: "m1"}},
	}
	if err := invalidoSemTokens.Validate(); err != combo.ErrComboInvalido {
		t.Errorf("esperava ErrComboInvalido para MaxTokens <= 0, obteve %v", err)
	}
}

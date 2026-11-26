package ativos

import (
	"errors"
	"strings"
	"testing"
)

func TestSchemasValidarFichaPersonagem(t *testing.T) {
	valida := FichaPersonagem{
		ID:                "p1",
		Nome:              "Narrador",
		DescricaoImutavel: "Palito azul com oculos redondos e jaqueta cinza",
		Fichas: map[string]string{
			"turnaround":     "flow/p1_turnaround.png",
			"expressoes_3x2": "flow/p1_expressoes.png",
			"poses_6x":       "flow/p1_poses.png",
		},
	}
	if err := valida.Validar(); err != nil {
		t.Fatalf("ficha valida falhou: %v", err)
	}

	casos := []struct {
		nome  string
		ficha FichaPersonagem
	}{
		{
			nome:  "id vazio",
			ficha: FichaPersonagem{Nome: "A", DescricaoImutavel: "Desc", Fichas: valida.Fichas},
		},
		{
			nome:  "nome vazio",
			ficha: FichaPersonagem{ID: "p1", DescricaoImutavel: "Desc", Fichas: valida.Fichas},
		},
		{
			nome:  "descricao imutavel vazia",
			ficha: FichaPersonagem{ID: "p1", Nome: "A", Fichas: valida.Fichas},
		},
		{
			nome: "falta poses_6x",
			ficha: FichaPersonagem{
				ID:                "p1",
				Nome:              "A",
				DescricaoImutavel: "Desc",
				Fichas: map[string]string{
					"turnaround":     "a.png",
					"expressoes_3x2": "b.png",
				},
			},
		},
		{
			nome: "ficha obrigatoria vazia",
			ficha: FichaPersonagem{
				ID:                "p1",
				Nome:              "A",
				DescricaoImutavel: "Desc",
				Fichas: map[string]string{
					"turnaround":     "a.png",
					"expressoes_3x2": "   ",
					"poses_6x":       "c.png",
				},
			},
		},
	}
	for _, tc := range casos {
		t.Run(tc.nome, func(t *testing.T) {
			err := tc.ficha.Validar()
			if !errors.Is(err, ErrFichaInvalida) {
				t.Fatalf("esperado ErrFichaInvalida, obtido: %v", err)
			}
		})
	}
}

func TestSchemasValidarDescritorSprite(t *testing.T) {
	valido := DescritorSprite{
		ID:           "s1",
		PersonagemID: "p1",
		Tipo:         "pose",
		ArquivoPNG:   "sprites/p1_pose_1.png",
		Largura:      256,
		Altura:       256,
		AlfaValido:   true,
		DespillVerde: true,
		BoundingBox:  BoundingBox{X: 0, Y: 0, W: 128, H: 128},
	}
	if err := valido.Validar(); err != nil {
		t.Fatalf("sprite valido falhou: %v", err)
	}

	invalido := valido
	invalido.AlfaValido = false
	if !errors.Is(invalido.Validar(), ErrSpriteInvalido) {
		t.Fatalf("esperado ErrSpriteInvalido para alfa_valido=false")
	}

	pequeno := valido
	pequeno.Largura = 32
	if !errors.Is(pequeno.Validar(), ErrSpriteInvalido) {
		t.Fatalf("esperado ErrSpriteInvalido para largura abaixo do minimo")
	}

	tipoErrado := valido
	tipoErrado.Tipo = "desconhecido"
	if !errors.Is(tipoErrado.Validar(), ErrSpriteInvalido) {
		t.Fatalf("esperado ErrSpriteInvalido para tipo invalido")
	}

	bboxInvalida := valido
	bboxInvalida.BoundingBox = BoundingBox{X: -1, Y: 0, W: 64, H: 64}
	if !errors.Is(bboxInvalida.Validar(), ErrSpriteInvalido) {
		t.Fatalf("esperado ErrSpriteInvalido para bbox negativa")
	}

	semArquivo := valido
	semArquivo.ArquivoPNG = ""
	if !errors.Is(semArquivo.Validar(), ErrSpriteInvalido) {
		t.Fatalf("esperado ErrSpriteInvalido para arquivo vazio")
	}

	semID := valido
	semID.ID = ""
	if !errors.Is(semID.Validar(), ErrSpriteInvalido) {
		t.Fatalf("esperado ErrSpriteInvalido para id vazio")
	}

	semDespill := valido
	semDespill.DespillVerde = false
	if !errors.Is(semDespill.Validar(), ErrSpriteInvalido) {
		t.Fatalf("esperado ErrSpriteInvalido para despill_verde=false")
	}
}

func TestSchemasSobreposicaoBoundingBox(t *testing.T) {
	a := BoundingBox{X: 0, Y: 0, W: 100, H: 100}
	b := BoundingBox{X: 100, Y: 0, W: 100, H: 100}
	if a.Sobrepoe(b) {
		t.Fatalf("caixas adjacentes nao devem sobrepor")
	}
	c := BoundingBox{X: 50, Y: 50, W: 100, H: 100}
	if !a.Sobrepoe(c) {
		t.Fatalf("caixas cruzadas devem sobrepor")
	}
}

func TestSchemasModoVisualValido(t *testing.T) {
	for _, m := range []string{"codigo", "flow", "hibrido"} {
		if !ModoVisualValido(m) {
			t.Fatalf("modo %q deveria ser valido", m)
		}
	}
	if ModoVisualValido("invalido") || ModoVisualValido("") {
		t.Fatalf("modo invalido nao deveria ser aceito")
	}
}

func inventarioFixture() InventarioEstilo {
	return InventarioEstilo{
		Estilo:           "molde-pilula",
		PartesSuportadas: []string{"coque", "careca", "chapeu-coco", "oculos-redondo"},
		Itens: map[string][]string{
			"pele":     {"pele-clara"},
			"traje":    {"terno-azul", "camisa-branca"},
			"penteado": {"coque", "careca", "moicano"},
			"adorno":   {"chapeu-coco", "cartola", "gravata-vermelha", "oculos-redondo"},
			"pose":     {"apontar", "bracos-cruzados"},
			"humor":    {"feliz", "bravo"},
		},
	}
}

func elencoFixture() Elenco {
	return Elenco{SchemaVersion: 2, Estilo: "molde-pilula", Personagens: []PersonagemElenco{{
		ID: "ana", Nome: "Ana",
		Cast: map[string]string{"skin": "pele-clara", "suit": "terno-azul", "shirt": "camisa-branca", "tie": "gravata-vermelha", "hair": "coque", "hair_col": "#3a2a1a", "hat": "chapeu-coco", "glasses": "oculos-redondo"},
	}}}
}

func TestElencoCastValido(t *testing.T) {
	if err := ValidarElenco(elencoFixture(), inventarioFixture()); err != nil {
		t.Fatalf("elenco valido reprovado: %v", err)
	}
}

func TestElencoCastBloqueios(t *testing.T) {
	casos := []struct {
		nome, campo, valor string
		alvo               error
	}{
		{"hair fora das partes suportadas", "hair", "moicano", ErrParteNaoSuportada},
		{"hat fora das partes suportadas", "hat", "cartola", ErrParteNaoSuportada},
		{"item inexistente no estilo", "suit", "terno-fantasma", ErrEstiloInexistente},
		{"campo desconhecido no cast", "asa", "asa-anjo", ErrParteNaoSuportada},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			el := elencoFixture()
			el.Personagens[0].Cast[c.campo] = c.valor
			err := ValidarElenco(el, inventarioFixture())
			if !errors.Is(err, c.alvo) || !strings.Contains(err.Error(), c.campo) {
				t.Fatalf("esperado %v listando %q, obtido %v", c.alvo, c.campo, err)
			}
		})
	}
}

func TestElencoCastEstiloEListagem(t *testing.T) {
	el := elencoFixture()
	el.Estilo = "aquarela"
	if err := ValidarElenco(el, inventarioFixture()); !errors.Is(err, ErrEstiloInexistente) {
		t.Fatalf("estilo divergente: esperado ErrEstiloInexistente, obtido %v", err)
	}
	el = elencoFixture()
	el.Personagens[0].Cast["hair"] = "moicano"
	el.Personagens[0].Cast["hat"] = "cartola"
	err := ValidarElenco(el, inventarioFixture())
	if err == nil || !strings.Contains(err.Error(), `"hair"`) || !strings.Contains(err.Error(), `"hat"`) {
		t.Fatalf("esperado erro listando hair e hat, obtido %v", err)
	}
	el = elencoFixture()
	el.Personagens[0].ID = ""
	if err := ValidarElenco(el, inventarioFixture()); !errors.Is(err, ErrElencoInvalido) {
		t.Fatalf("personagem sem id: esperado ErrElencoInvalido, obtido %v", err)
	}
}

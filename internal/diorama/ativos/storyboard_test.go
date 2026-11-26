package ativos

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestStoryboardCompilarEValidarSucesso(t *testing.T) {
	dir := criarFixtureCiclo1(t, "codigo", "")
	contratos, err := CarregarContratosEntrada(dir)
	if err != nil {
		t.Fatal(err)
	}
	sprites := []DescritorSprite{
		{ID: "narrador_pose_1", PersonagemID: "narrador", Tipo: "pose", ArquivoPNG: "sprites/n1.png", Largura: 256, Altura: 256, AlfaValido: true, DespillVerde: true, BoundingBox: BoundingBox{X: 0, Y: 0, W: 128, H: 128}},
		{ID: "narrador_expr_1", PersonagemID: "narrador", Tipo: "expressao", ArquivoPNG: "sprites/n2.png", Largura: 256, Altura: 256, AlfaValido: true, DespillVerde: true, BoundingBox: BoundingBox{X: 140, Y: 0, W: 128, H: 128}},
	}
	cenarios := []Cenario{
		{ID: "cen_01", Nome: "Cenario 1", Modo: "codigo", Arquivo: "flow/cen_01.svg", HashEntrada: "h1"},
		{ID: "cen_02", Nome: "Cenario 2", Modo: "codigo", Arquivo: "flow/cen_02.svg", HashEntrada: "h2"},
	}

	sb, err := ExecutarEtapa5Storyboard(dir, contratos.Projeto, contratos.Script, sprites, cenarios)
	if err != nil {
		t.Fatalf("Etapa 5 falhou: %v", err)
	}
	if len(sb.Cenas) != len(contratos.Script.Lines) {
		t.Fatalf("esperado %d cenas, obtido %d", len(contratos.Script.Lines), len(sb.Cenas))
	}
	if _, err := os.Stat(filepath.Join(dir, "cenas.json")); err != nil {
		t.Fatalf("cenas.json nao gravado: %v", err)
	}
}

func TestStoryboardBloqueios(t *testing.T) {
	dir := criarFixtureCiclo1(t, "codigo", "")
	contratos, _ := CarregarContratosEntrada(dir)
	cenarios := []Cenario{{ID: "cen_01"}, {ID: "cen_02"}}
	elems4 := []ElementoMovel{
		{ID: "e1", Tipo: "personagem", Acao: "entra"},
		{ID: "e2", Tipo: "objeto", Acao: "destaca"},
		{ID: "e3", Tipo: "camera", Acao: "zoom_in"},
		{ID: "e4", Tipo: "efeito", Acao: "particulas"},
	}

	t.Run("cobertura incompleta de linhas", func(t *testing.T) {
		sb := Storyboard{SchemaVersion: 1, TotalLinhas: 4, Cenas: []Cena{{IndiceLinha: 0, CenarioID: "cen_01", MetaforaVisual: "m", ElementosMoveis: elems4}}}
		if !errors.Is(ValidarStoryboard(sb, contratos.Script, cenarios), ErrStoryboardInvalido) {
			t.Fatalf("esperado ErrStoryboardInvalido para cobertura incompleta")
		}
	})

	t.Run("metafora vazia", func(t *testing.T) {
		cenas := make([]Cena, len(contratos.Script.Lines))
		for i := range cenas {
			cenas[i] = Cena{IndiceLinha: i, CenarioID: "cen_01", MetaforaVisual: "ok", ElementosMoveis: elems4}
		}
		cenas[2].MetaforaVisual = ""
		sb := Storyboard{SchemaVersion: 1, TotalLinhas: len(cenas), Cenas: cenas}
		if !errors.Is(ValidarStoryboard(sb, contratos.Script, cenarios), ErrStoryboardInvalido) {
			t.Fatalf("esperado ErrStoryboardInvalido para metafora vazia")
		}
	})

	t.Run("menos de 4 elementos moveis", func(t *testing.T) {
		cenas := make([]Cena, len(contratos.Script.Lines))
		for i := range cenas {
			cenas[i] = Cena{IndiceLinha: i, CenarioID: "cen_01", MetaforaVisual: "ok", ElementosMoveis: elems4}
		}
		cenas[1].ElementosMoveis = elems4[:3]
		sb := Storyboard{SchemaVersion: 1, TotalLinhas: len(cenas), Cenas: cenas}
		if !errors.Is(ValidarStoryboard(sb, contratos.Script, cenarios), ErrStoryboardInvalido) {
			t.Fatalf("esperado ErrStoryboardInvalido para < 4 elementos moveis")
		}
	})

	t.Run("cenario inexistente", func(t *testing.T) {
		cenas := make([]Cena, len(contratos.Script.Lines))
		for i := range cenas {
			cenas[i] = Cena{IndiceLinha: i, CenarioID: "cen_01", MetaforaVisual: "ok", ElementosMoveis: elems4}
		}
		cenas[0].CenarioID = "cen_999"
		sb := Storyboard{SchemaVersion: 1, TotalLinhas: len(cenas), Cenas: cenas}
		if !errors.Is(ValidarStoryboard(sb, contratos.Script, cenarios), ErrStoryboardInvalido) {
			t.Fatalf("esperado ErrStoryboardInvalido para cenario inexistente")
		}
	})

	t.Run("sem cenarios falha compilacao", func(t *testing.T) {
		if _, err := ExecutarEtapa5Storyboard(dir, contratos.Projeto, contratos.Script, nil, nil); !errors.Is(err, ErrStoryboardInvalido) {
			t.Fatalf("esperado ErrStoryboardInvalido sem cenarios, obtido: %v", err)
		}
	})
}

func TestExtrairDispositivosPlanoMD(t *testing.T) {
	plano := `# Plano
| Batida | Tempo | Natureza | Narração | Dispositivo Visual (Diorama) | Fonte |
| :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | 00-05s | transicao | Texto 1 | Miniatura de cidade com moedas | — |
| 2 | 05-15s | fato | Texto 2 | Predio do Banco Central com torre | bcb |
`
	got := ExtrairDispositivosPlanoMD(plano)
	if len(got) != 2 || got[0] != "Miniatura de cidade com moedas" || got[1] != "Predio do Banco Central com torre" {
		t.Fatalf("dispositivos extraidos inesperados: %+v", got)
	}
}

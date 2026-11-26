package ativos

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type recortadorFake struct {
	sobrescreverAlfa    *bool
	sobrescreverDespill *bool
	forcarSobreposicao  bool
	falhar              error
	chamadas            int
}

func (r *recortadorFake) RecortarEFatiar(_ context.Context, _ string, ficha FichaPersonagem) ([]DescritorSprite, error) {
	r.chamadas++
	if r.falhar != nil {
		return nil, r.falhar
	}
	alfa := true
	if r.sobrescreverAlfa != nil {
		alfa = *r.sobrescreverAlfa
	}
	despill := true
	if r.sobrescreverDespill != nil {
		despill = *r.sobrescreverDespill
	}
	bbox2 := BoundingBox{X: 140, Y: 0, W: 128, H: 128}
	if r.forcarSobreposicao {
		bbox2 = BoundingBox{X: 20, Y: 20, W: 128, H: 128}
	}
	return []DescritorSprite{
		{
			ID:           ficha.ID + "_pose_1",
			PersonagemID: ficha.ID,
			Tipo:         "pose",
			ArquivoPNG:   "sprites/" + ficha.ID + "_pose_1.png",
			Largura:      256,
			Altura:       256,
			AlfaValido:   alfa,
			DespillVerde: despill,
			BoundingBox:  BoundingBox{X: 0, Y: 0, W: 128, H: 128},
		},
		{
			ID:           ficha.ID + "_expr_1",
			PersonagemID: ficha.ID,
			Tipo:         "expressao",
			ArquivoPNG:   "sprites/" + ficha.ID + "_expr_1.png",
			Largura:      256,
			Altura:       256,
			AlfaValido:   alfa,
			DespillVerde: despill,
			BoundingBox:  bbox2,
		},
	}, nil
}

func TestElencoExecutarEtapa3Sucesso(t *testing.T) {
	dir := criarFixtureCiclo1(t, "codigo", "")
	elenco := []FichaPersonagem{
		{
			ID:                "narrador",
			Nome:              "Narrador",
			DescricaoImutavel: "Palito azul com oculos redondos",
			Fichas: map[string]string{
				"turnaround":     "flow/narrador_turnaround.png",
				"expressoes_3x2": "flow/narrador_expressoes.png",
				"poses_6x":       "flow/narrador_poses.png",
			},
		},
	}
	rec := &recortadorFake{}
	sprites, err := ExecutarEtapa3Elenco(context.Background(), dir, elenco, rec)
	if err != nil {
		t.Fatalf("Etapa 3 falhou: %v", err)
	}
	if len(sprites) != 2 || rec.chamadas != 1 {
		t.Fatalf("esperado 2 sprites e 1 chamada, obtido sprites=%d chamadas=%d", len(sprites), rec.chamadas)
	}
	if _, err := os.Stat(filepath.Join(dir, "flow", "elenco.json")); err != nil {
		t.Fatalf("flow/elenco.json nao gravado: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "sprites", "narrador_pose_1.json")); err != nil {
		t.Fatalf("sprites/narrador_pose_1.json nao gravado: %v", err)
	}
}

func TestElencoEtapa3Bloqueios(t *testing.T) {
	dir := criarFixtureCiclo1(t, "codigo", "")
	fichaValida := FichaPersonagem{
		ID:                "p1",
		Nome:              "P1",
		DescricaoImutavel: "Palito verde",
		Fichas: map[string]string{
			"turnaround":     "a.png",
			"expressoes_3x2": "b.png",
			"poses_6x":       "c.png",
		},
	}

	t.Run("elenco vazio", func(t *testing.T) {
		if _, err := ExecutarEtapa3Elenco(context.Background(), dir, nil, &recortadorFake{}); !errors.Is(err, ErrFichaInvalida) {
			t.Fatalf("esperado ErrFichaInvalida, obtido: %v", err)
		}
	})

	t.Run("ficha ausente", func(t *testing.T) {
		invalida := fichaValida
		invalida.Fichas = map[string]string{"turnaround": "a.png"}
		if _, err := ExecutarEtapa3Elenco(context.Background(), dir, []FichaPersonagem{invalida}, &recortadorFake{}); !errors.Is(err, ErrFichaInvalida) {
			t.Fatalf("esperado ErrFichaInvalida, obtido: %v", err)
		}
	})

	t.Run("alfa invalido", func(t *testing.T) {
		falso := false
		_, err := ExecutarEtapa3Elenco(context.Background(), dir, []FichaPersonagem{fichaValida}, &recortadorFake{sobrescreverAlfa: &falso})
		if !errors.Is(err, ErrSpriteInvalido) {
			t.Fatalf("esperado ErrSpriteInvalido, obtido: %v", err)
		}
	})

	t.Run("sobreposicao de bounding_box", func(t *testing.T) {
		_, err := ExecutarEtapa3Elenco(context.Background(), dir, []FichaPersonagem{fichaValida}, &recortadorFake{forcarSobreposicao: true})
		if !errors.Is(err, ErrSobreposicaoSprite) {
			t.Fatalf("esperado ErrSobreposicaoSprite, obtido: %v", err)
		}
	})

	t.Run("falha no recortador", func(t *testing.T) {
		errRec := errors.New("rembg falhou")
		_, err := ExecutarEtapa3Elenco(context.Background(), dir, []FichaPersonagem{fichaValida}, &recortadorFake{falhar: errRec})
		if !errors.Is(err, errRec) {
			t.Fatalf("esperado erro propagado, obtido: %v", err)
		}
	})

	t.Run("recortador padrao local gera PNGs reais RGBA no disco", func(t *testing.T) {
		sprites, err := ExecutarEtapa3Elenco(context.Background(), dir, []FichaPersonagem{fichaValida}, nil)
		if err != nil || len(sprites) == 0 {
			t.Fatalf("recortador padrao falhou: sprites=%d err=%v", len(sprites), err)
		}
		for _, sp := range sprites {
			if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(sp.ArquivoPNG))); err != nil {
				t.Fatalf("arquivo PNG real ausente %s: %v", sp.ArquivoPNG, err)
			}
		}
	})
}

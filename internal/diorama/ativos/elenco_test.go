package ativos

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/video/inventario"
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

// inventarioVetorialFixture cria o estilo molde-pilula com pele, traje e penteado aprovados.
func inventarioVetorialFixture(t *testing.T, status string) inventario.Store {
	t.Helper()
	s := inventario.Store{Raiz: t.TempDir()}
	base := filepath.Join(s.Raiz, "inventario", "video", "estilos", "molde-pilula")
	itens := map[string]string{"pele": "pele-clara", "traje": "terno-azul", "penteado": "coque"}
	for tipo := range itens {
		if err := os.MkdirAll(filepath.Join(base, "itens", tipo), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	gravarJSONTeste(t, filepath.Join(base, "estilo.json"), inventario.Estilo{ID: "molde-pilula", Status: status, Versao: "1.0.0", PartesSuportadas: []string{"coque"}})
	for tipo, id := range itens {
		gravarJSONTeste(t, filepath.Join(base, "itens", tipo, id+".json"), inventario.Item{ID: id, Tipo: tipo, Estilo: "molde-pilula", Aprovado: true})
	}
	return s
}

func pedidoAnaFixture() []PedidoPersonagem {
	return []PedidoPersonagem{{
		ID: "ana", Nome: "Ana", Tags: []string{"pre-historia"},
		Cast: map[string]string{"skin": "pele-clara", "suit": "terno-azul", "hair": "coque", "hair_col": "#333333", "tie": "colar-de-osso"},
	}}
}

func TestElencoVetorialMontaCastSemFlow(t *testing.T) {
	s := inventarioVetorialFixture(t, "validado")
	proj := t.TempDir()
	rel, err := ExecutarEtapa3Vetorial(proj, s, "molde-pilula", pedidoAnaFixture(), "demo-2026-11-26")
	if err != nil {
		t.Fatalf("etapa 3 vetorial: %v", err)
	}
	if rel.Creditos != 0 || rel.FlowChamado {
		t.Fatalf("esperado 0 credito e Flow nao chamado, obtido %+v", rel)
	}
	var el Elenco
	b, err := os.ReadFile(filepath.Join(proj, "elenco.json"))
	if err != nil {
		t.Fatalf("elenco.json nao gravado: %v", err)
	}
	if err := json.Unmarshal(b, &el); err != nil {
		t.Fatalf("elenco.json invalido: %v", err)
	}
	if el.Estilo != "molde-pilula" || len(el.Personagens) != 1 || el.Personagens[0].Cast["hair"] != "coque" || len(el.Personagens[0].Cast) != 5 {
		t.Fatalf("elenco inesperado: %+v", el)
	}
	if _, err := os.Stat(filepath.Join(proj, "flow")); !os.IsNotExist(err) {
		t.Fatalf("caminho vetorial nao pode escrever em flow/: %v", err)
	}
	if _, err := os.Stat(filepath.Join(proj, "elenco_relatorio.json")); err != nil {
		t.Fatalf("elenco_relatorio.json nao gravado: %v", err)
	}
}

func TestElencoVetorialItemNovoNaoAprovado(t *testing.T) {
	s := inventarioVetorialFixture(t, "validado")
	rel, err := ExecutarEtapa3Vetorial(t.TempDir(), s, "molde-pilula", pedidoAnaFixture(), "demo-2026-11-26")
	if err != nil {
		t.Fatalf("etapa 3 vetorial: %v", err)
	}
	if len(rel.Lacunas) != 1 || rel.Lacunas[0].Item != "colar-de-osso" || rel.Lacunas[0].Tipo != "adorno" {
		t.Fatalf("esperada 1 lacuna adorno colar-de-osso, obtido %+v", rel.Lacunas)
	}
	var it inventario.Item
	caminho := filepath.Join(s.Raiz, "inventario", "video", "estilos", "molde-pilula", "itens", "adorno", "colar-de-osso.json")
	b, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatalf("item novo nao gravado: %v", err)
	}
	if err := json.Unmarshal(b, &it); err != nil {
		t.Fatalf("item novo invalido: %v", err)
	}
	if it.Aprovado || it.Origem != "demo-2026-11-26" || it.EstiloVersao != "1.0.0" || len(it.Tags) != 1 {
		t.Fatalf("item novo deve entrar nao aprovado com origem, versao e tags: %+v", it)
	}
	achados, err := s.Buscar("molde-pilula", "adorno", nil)
	if err != nil || len(achados) != 0 {
		t.Fatalf("item nao aprovado nao pode aparecer na busca: %+v err=%v", achados, err)
	}
}

func TestElencoVetorialBloqueios(t *testing.T) {
	t.Run("estilo rascunho", func(t *testing.T) {
		s := inventarioVetorialFixture(t, "rascunho")
		if _, err := ExecutarEtapa3Vetorial(t.TempDir(), s, "molde-pilula", pedidoAnaFixture(), "o"); !errors.Is(err, ErrEstiloNaoValidado) {
			t.Fatalf("esperado ErrEstiloNaoValidado, obtido %v", err)
		}
	})
	t.Run("estilo inexistente", func(t *testing.T) {
		s := inventarioVetorialFixture(t, "validado")
		if _, err := ExecutarEtapa3Vetorial(t.TempDir(), s, "aquarela", pedidoAnaFixture(), "o"); !errors.Is(err, ErrEstiloInexistente) {
			t.Fatalf("esperado ErrEstiloInexistente, obtido %v", err)
		}
	})
	t.Run("campo desconhecido e personagem sem id", func(t *testing.T) {
		s := inventarioVetorialFixture(t, "validado")
		ruim := []PedidoPersonagem{{ID: "ana", Cast: map[string]string{"asa": "x"}}}
		if _, err := ExecutarEtapa3Vetorial(t.TempDir(), s, "molde-pilula", ruim, "o"); !errors.Is(err, ErrParteNaoSuportada) {
			t.Fatalf("esperado ErrParteNaoSuportada, obtido %v", err)
		}
		if _, err := ExecutarEtapa3Vetorial(t.TempDir(), s, "molde-pilula", []PedidoPersonagem{{}}, "o"); !errors.Is(err, ErrElencoInvalido) {
			t.Fatalf("esperado ErrElencoInvalido, obtido %v", err)
		}
	})
}

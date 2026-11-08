package redc_test

import (
	"testing"

	"legoloop/internal/redc"
)

func TestCatalogVideos_SuccessAndRanking(t *testing.T) {
	rawJSON := []byte(`[
		{
			"id": "v1",
			"titulo": "Metodo A - Como viralizar",
			"views": 100000,
			"likes": 5000,
			"comentarios": 500,
			"duracao_seg": 60,
			"gancho_3s": "Pare tudo que esta fazendo",
			"metodo": "Engenharia Reversa Viral",
			"cluster_tema": "crescimento",
			"transcricao": "Passo 1 analise seus concorrentes..."
		},
		{
			"id": "v2",
			"titulo": "Metodo B - Ganhar inscritos",
			"views": 500000,
			"likes": 25000,
			"comentarios": 2500,
			"duracao_seg": 45,
			"gancho_3s": "O segredo que ninguem te conta",
			"metodo": "Funil de Retencao",
			"cluster_tema": "crescimento",
			"transcricao": "Neste video mostro como reter 80%..."
		}
	]`)

	videos, err := redc.CatalogVideos(rawJSON)
	if err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}

	if len(videos) != 2 {
		t.Fatalf("esperava 2 videos, obteve %d", len(videos))
	}

	// v2 tem 500k views e mais likes, deve vir em 1º lugar no ranking de viralidade
	if videos[0].ID != "v2" {
		t.Errorf("esperava v2 no topo do ranking, obteve %s", videos[0].ID)
	}
	if videos[1].ID != "v1" {
		t.Errorf("esperava v1 em 2º lugar, obteve %s", videos[1].ID)
	}
}

func TestCatalogVideos_Errors(t *testing.T) {
	t.Run("dados corrompidos", func(t *testing.T) {
		_, err := redc.CatalogVideos([]byte(`{ invalido `))
		if err == nil || err != redc.ErrDadosCorrompidos {
			t.Errorf("esperava ErrDadosCorrompidos, obteve %v", err)
		}
	})

	t.Run("lote vazio", func(t *testing.T) {
		_, err := redc.CatalogVideos([]byte(`[]`))
		if err == nil || err != redc.ErrVideosInsuf {
			t.Errorf("esperava ErrVideosInsuf, obteve %v", err)
		}
	})
}

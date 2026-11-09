package redc_test

import (
	"testing"

	"legoloop/internal/redc"
)

func TestClusterMethods_Deduplication(t *testing.T) {
	videos := []redc.VideoEntry{
		{
			ID:          "v1",
			Titulo:      "3 passos para viralizar",
			Metodo:      "Engenharia Reversa Viral",
			ClusterTema: "crescimento",
			Score:       1000,
		},
		{
			ID:          "v2",
			Titulo:      "Como viralizar rapido",
			Metodo:      "Engenharia Reversa Viral",
			ClusterTema: "crescimento",
			Score:       800,
		},
		{
			ID:          "v3",
			Titulo:      "Edicao de video com Remotion",
			Metodo:      "Motion Procedural",
			ClusterTema: "edicao",
			Score:       950,
		},
	}

	clusters := redc.ClusterMethods(videos)

	if len(clusters) != 2 {
		t.Fatalf("esperava 2 clusters unicos, obteve %d", len(clusters))
	}

	// Primeiro cluster deve ser o de maior score acumulado ou pico
	var foundViral, foundMotion bool
	for _, c := range clusters {
		if c.Metodo == "Engenharia Reversa Viral" {
			foundViral = true
			if len(c.VideoIDs) != 2 {
				t.Errorf("esperava 2 videos no cluster viral, obteve %d", len(c.VideoIDs))
			}
		}
		if c.Metodo == "Motion Procedural" {
			foundMotion = true
			if len(c.VideoIDs) != 1 {
				t.Errorf("esperava 1 video no cluster motion, obteve %d", len(c.VideoIDs))
			}
		}
	}

	if !foundViral || !foundMotion {
		t.Errorf("clusters esperados nao foram encontrados: %+v", clusters)
	}
}

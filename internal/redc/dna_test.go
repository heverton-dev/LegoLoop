package redc_test

import (
	"testing"

	"legoloop/internal/redc"
)

func TestExtractDNA_Success(t *testing.T) {
	topVideos := []redc.VideoEntry{
		{
			ID:         "v1",
			Titulo:     "Hack viral 1",
			DuracaoSeg: 45,
			Gancho3s:   "Pare de cometer esse erro agora",
			Metodo:     "Hook Reverso",
			Score:      5000,
		},
		{
			ID:         "v2",
			Titulo:     "Hack viral 2",
			DuracaoSeg: 60,
			Gancho3s:   "O segredo dos canais de 1 milhao",
			Metodo:     "Pattern Interrupt",
			Score:      4200,
		},
	}

	dna, err := redc.ExtractDNA("meu-canal", topVideos)
	if err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}

	if dna.CanalSlug != "meu-canal" {
		t.Errorf("esperava canal meu-canal, obteve %s", dna.CanalSlug)
	}
	if dna.TomDeVoz == "" {
		t.Errorf("tom de voz nao deve ser vazio")
	}
	if dna.RitmoCortes == "" {
		t.Errorf("ritmo de cortes nao deve ser vazio")
	}
	if len(dna.DesignTokens) == 0 {
		t.Errorf("design tokens nao devem ser vazios")
	}
	if len(dna.MotionTokens) == 0 {
		t.Errorf("motion tokens nao devem ser vazios")
	}
	if len(dna.TopGanchos) != 2 {
		t.Errorf("esperava 2 top ganchos, obteve %d", len(dna.TopGanchos))
	}
}

func TestExtractDNA_EmptyList(t *testing.T) {
	_, err := redc.ExtractDNA("canal", nil)
	if err != redc.ErrVideosInsuf {
		t.Errorf("esperava ErrVideosInsuf, obteve %v", err)
	}
}

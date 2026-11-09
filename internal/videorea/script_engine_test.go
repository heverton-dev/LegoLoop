package videorea

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"legoloop/internal/brain"
	"legoloop/internal/interview"
)

type mockBrainScript struct {
	resposta string
	err      error
}

func (m *mockBrainScript) Infer(ctx context.Context, prompt string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.resposta, nil
}
func (m *mockBrainScript) Chat(ctx context.Context, msgs []interview.Message) (string, error) {
	return m.Infer(ctx, "")
}
func (m *mockBrainScript) Name() string           { return "mock-brain" }
func (m *mockBrainScript) Model() string          { return "mock-model" }
func (m *mockBrainScript) Timeout() time.Duration { return time.Second }
func (m *mockBrainScript) Effort() brain.Effort   { return brain.EffortMedium }

func TestBuildCartoonScript_TemaVazio(t *testing.T) {
	ctx := context.Background()
	mock := &mockBrainScript{}

	_, err := BuildCartoonScript(ctx, "   ", mock)
	if !errors.Is(err, ErrInvalidCartoonTheme) {
		t.Fatalf("esperava ErrInvalidCartoonTheme, obteve: %v", err)
	}
}

func TestBuildCartoonScript_SucessoDeterminístico(t *testing.T) {
	ctx := context.Background()
	jsonBrain := `{
		"title": "O Enigma do ET de Varginha",
		"theme": "ET de Varginha",
		"hook": "Uma cidade mineira trocou a produção de café por um visitante alienígena.",
		"fact_or_legend": "Fato Histórico com Relatos Folclóricos",
		"scenes": [
			{
				"index": 1,
				"timestamp": "00:00",
				"narration": "Janeiro de 1996. Após um temporal violento, três garotas avistam uma criatura agachada com olhos vermelhos.",
				"emotion_tag": "[reflexivo]",
				"visual_cue": "Cena de temporal com silhueta da criatura em 3D",
				"keywords": ["varginha", "criatura", "temporal"],
				"accent_color": "#ff3b30"
			},
			{
				"index": 2,
				"timestamp": "00:08",
				"narration": "A notícia corre rápido. Repórteres e caminhões do exército cercam a área, mas o inquérito oficial é arquivado.",
				"emotion_tag": "[animado]",
				"visual_cue": "Caminhões militares e holofotes se movendo",
				"keywords": ["exercito", "inquerito", "boato"],
				"accent_color": "#ff9500"
			},
			{
				"index": 3,
				"timestamp": "00:18",
				"narration": "Nenhum disco pousou de novo, mas a cidade construiu uma caixa d'água de disco voador para eternizar a visita.",
				"emotion_tag": "[dá uma risadinha]",
				"visual_cue": "Caixa d'água em formato de disco voador",
				"keywords": ["disco voador", "caixa dagua", "patrimonio"],
				"accent_color": "#34c759"
			}
		],
		"total_words": 78
	}`

	mock := &mockBrainScript{resposta: jsonBrain}
	script, err := BuildCartoonScript(ctx, "ET de Varginha", mock)
	if err != nil {
		t.Fatalf("erro inesperado ao gerar roteiro: %v", err)
	}

	if script.Title == "" {
		t.Errorf("esperava título não-vazio")
	}
	if !strings.Contains(script.Hook, "visitante") {
		t.Errorf("hook incorreto: %s", script.Hook)
	}
	if len(script.Scenes) < 3 {
		t.Errorf("esperava ao menos 3 cenas, obteve: %d", len(script.Scenes))
	}
	if script.FactOrLegend == "" {
		t.Errorf("fact_or_legend deve estar preenchido")
	}
	if script.TotalWords < 50 {
		t.Errorf("total_words muito baixo para short de 20-30s: %d", script.TotalWords)
	}
}

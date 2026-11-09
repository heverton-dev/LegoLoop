package brain_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"legoloop/internal/brain"
	"legoloop/internal/interview"
)

type mockBrain struct {
	name      string
	model     string
	shouldErr bool
	response  string
	delay     time.Duration
}

func (m *mockBrain) Chat(_ context.Context, _ []interview.Message) (string, error) {
	if m.delay > 0 {
		time.Sleep(m.delay)
	}
	if m.shouldErr {
		return "", errors.New("mock brain error")
	}
	return m.response, nil
}

func (m *mockBrain) Infer(ctx context.Context, prompt string) (string, error) {
	return m.Chat(ctx, []interview.Message{{Role: "user", Content: prompt}})
}

func (m *mockBrain) Name() string           { return m.name }
func (m *mockBrain) Model() string          { return m.model }
func (m *mockBrain) Timeout() time.Duration { return 5 * time.Second }
func (m *mockBrain) Effort() brain.Effort   { return brain.EffortMedium }

func TestAntibodiesMemory(t *testing.T) {
	db := brain.NewAntibodiesDB()

	err := db.Record(
		"internal/videorea",
		"type mismatch in crawler",
		"Types struct omite o campo VideoMetadata",
		"Adicionar VideoMetadata em types.go",
	)
	if err != nil {
		t.Fatalf("falha ao gravar anticorpo: %v", err)
	}

	ab, found := db.Find("internal/videorea", "type mismatch in crawler")
	if !found || ab == nil {
		t.Fatal("esperava encontrar anticorpo gravado")
	}

	if ab.RootCause == "" || ab.DiffSolution == "" {
		t.Fatalf("campos do anticorpo devem estar preenchidos, obteve: %+v", ab)
	}

	contextAdvice := db.InjectPreventiveContext("internal/videorea")
	if !strings.Contains(contextAdvice, "ANTICORPO PREVENTIVO") || !strings.Contains(contextAdvice, "VideoMetadata") {
		t.Errorf("contexto preventivo deve conter os anticorpos do pacote, obteve:\n%s", contextAdvice)
	}
}

func TestSpeculativeCascade(t *testing.T) {
	t.Run("FastBrainRespondeComSucesso", func(t *testing.T) {
		fast := &mockBrain{name: "fast-mock", model: "fast", response: "resposta rápida"}
		heavy := &mockBrain{name: "heavy-mock", model: "heavy", response: "resposta pesada"}

		res, modelUsed, err := brain.ExecuteSpeculative(context.Background(), fast, heavy, "prompt teste", 2*time.Second)
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if modelUsed != "fast-mock" || res != "resposta rápida" {
			t.Errorf("esperava uso do modelo rápido, obteve %s: %s", modelUsed, res)
		}
	})

	t.Run("FastBrainFalhaEHotSwapParaHeavy", func(t *testing.T) {
		fast := &mockBrain{name: "fast-mock", model: "fast", shouldErr: true}
		heavy := &mockBrain{name: "heavy-mock", model: "heavy", response: "resposta pesada após hot-swap"}

		res, modelUsed, err := brain.ExecuteSpeculative(context.Background(), fast, heavy, "prompt teste", 2*time.Second)
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if modelUsed != "heavy-mock" || res != "resposta pesada após hot-swap" {
			t.Errorf("esperava hot-swap para o modelo pesado, obteve %s: %s", modelUsed, res)
		}
	})
}

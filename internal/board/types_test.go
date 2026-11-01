package board_test

import (
	"encoding/json"
	"testing"

	"legoloop/internal/board"
)

func TestCardSerialization(t *testing.T) {
	card := board.Card{
		ID:        "card-123",
		Number:    42,
		Title:     "Implement GitHub Board Ingestor",
		Body:      "Closes #42. Touches internal/board/types.go",
		Status:    board.StatusReady,
		Labels:    []string{"feature", "core"},
		Assignees: []string{"developer-ai"},
		Branch:    "issue-42-github-ingestor",
		PRNumber:  101,
	}

	data, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("falha ao serializar Card: %v", err)
	}

	var parsed board.Card
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("falha ao deserializar Card: %v", err)
	}

	if parsed.ID != "card-123" || parsed.Number != 42 || parsed.Status != board.StatusReady {
		t.Fatalf("card deserializado incorreto: %+v", parsed)
	}
}

func TestWaveSerialization(t *testing.T) {
	wave := board.Wave{
		Index: 0,
		Cards: []board.Card{
			{Number: 1, Title: "Ticket A"},
			{Number: 2, Title: "Ticket B"},
		},
	}

	if len(wave.Cards) != 2 || wave.Index != 0 {
		t.Fatalf("wave inconsistente: %+v", wave)
	}
}

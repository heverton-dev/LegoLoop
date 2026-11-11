package board_test

import (
	"encoding/json"
	"strings"
	"testing"

	"legoloop/internal/board"
)

func TestTicketSerialization(t *testing.T) {
	ticket := board.Ticket{
		ID:        "ticket-123",
		Number:    42,
		Title:     "Implement GitHub Board Ingestor",
		Body:      "Closes #42. Touches internal/board/types.go",
		Status:    board.StatusReady,
		Labels:    []string{"feature", "core"},
		Assignees: []string{"developer-ai"},
		Branch:    "ticket-42-github-ingestor",
		PRNumber:  101,
	}

	data, err := json.Marshal(ticket)
	if err != nil {
		t.Fatalf("falha ao serializar Ticket: %v", err)
	}

	var parsed board.Ticket
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("falha ao deserializar Ticket: %v", err)
	}

	if parsed.ID != "ticket-123" || parsed.Number != 42 || parsed.Status != board.StatusReady {
		t.Fatalf("ticket deserializado incorreto: %+v", parsed)
	}
}

func TestWaveSerialization(t *testing.T) {
	var wave board.Wave
	entrada := `{"index":0,"tickets":[{"number":1,"title":"Ticket A"},{"number":2,"title":"Ticket B"}]}`
	if err := json.Unmarshal([]byte(entrada), &wave); err != nil {
		t.Fatal(err)
	}
	saida, err := json.Marshal(wave)
	if err != nil {
		t.Fatal(err)
	}
	if wave.Index != 0 || strings.Count(string(saida), `"title"`) != 2 {
		t.Fatalf("wave inconsistente: %+v", wave)
	}
}

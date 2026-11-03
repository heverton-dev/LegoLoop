package board_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"legoloop/internal/board"
)

func TestClient_FetchReadyCards_Success(t *testing.T) {
	// Cria servidor HTTP de teste simulando a API do GitHub com headers de rate-limit
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Remaining", "4500")
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(1*time.Hour).Unix(), 10))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[
			{
				"id": "node_card_1",
				"number": 10,
				"title": "Fix memory leak in scanner",
				"body": "Touches internal/assets/scanner.go",
				"status": "Ready",
				"labels": ["bug"]
			},
			{
				"id": "node_card_2",
				"number": 11,
				"title": "Add docs for routing",
				"body": "Touches docs/routes/README.md",
				"status": "Building",
				"labels": ["docs"]
			}
		]`))
	}))
	defer server.Close()

	client := board.NewClient(board.ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
		Guard:   board.NewTokenBucketGuard(100),
		HTTPCli: server.Client(),
	})

	cards, err := client.FetchReadyCards(context.Background())
	if err != nil {
		t.Fatalf("falha ao buscar cards em Ready: %v", err)
	}

	if len(cards) != 1 {
		t.Fatalf("esperava 1 card em Ready, obteve %d", len(cards))
	}

	if cards[0].Number != 10 || cards[0].Status != board.StatusReady {
		t.Fatalf("card incorreto: %+v", cards[0])
	}
}

func TestClient_FetchReadyCards_RateLimitHalt(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simula cota esgotada (apenas 20 requisições restantes)
		w.Header().Set("X-RateLimit-Remaining", "20")
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(15*time.Minute).Unix(), 10))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))
	defer server.Close()

	guard := board.NewTokenBucketGuard(100)
	client := board.NewClient(board.ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
		Guard:   guard,
		HTTPCli: server.Client(),
	})

	_, err := client.FetchReadyCards(context.Background())
	if err == nil {
		t.Fatalf("esperava erro de rate limit devido ao cabeçalho X-RateLimit-Remaining")
	}

	if !guard.IsHalted() {
		t.Fatalf("guard deveria ter sido travado em halted")
	}
}

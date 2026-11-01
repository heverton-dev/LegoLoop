package board_test

import (
	"errors"
	"testing"
	"time"

	"legoloop/internal/board"
)

func TestTokenBucketGuard_AllowedWhenQuotaSufficient(t *testing.T) {
	guard := board.NewTokenBucketGuard(100) // threshold mínimo = 100

	// Simula quota de 1500 chamadas restantes
	err := guard.Check(1500, time.Now().Add(1*time.Hour).Unix())
	if err != nil {
		t.Fatalf("esperava permissão com quota 1500, obteve erro: %v", err)
	}

	if guard.IsHalted() {
		t.Fatalf("guard não deveria estar em estado halted")
	}
}

func TestTokenBucketGuard_HaltedWhenQuotaExhausted(t *testing.T) {
	guard := board.NewTokenBucketGuard(100)

	// Simula quota de apenas 50 chamadas (menor que o piso de 100)
	err := guard.Check(50, time.Now().Add(30*time.Minute).Unix())
	if err == nil {
		t.Fatalf("esperava erro de cota esgotada com 50 chamadas")
	}

	if !errors.Is(err, board.ErrRateLimitExhausted) {
		t.Fatalf("esperava ErrRateLimitExhausted, obteve: %v", err)
	}

	if !guard.IsHalted() {
		t.Fatalf("guard deveria estar em estado halted")
	}

	// Tentativa subsequente deve ser barrada imediatamente
	err = guard.Check(50, time.Now().Add(30*time.Minute).Unix())
	if !errors.Is(err, board.ErrRateLimitExhausted) {
		t.Fatalf("esperava persistência do erro ErrRateLimitExhausted")
	}
}

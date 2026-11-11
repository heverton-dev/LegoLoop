package board

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrRateLimitExhausted indica que o piso seguro de requisições da GitHub API foi atingido.
// Equivalente à saída graciosa Exit 79 do super-board.
var ErrRateLimitExhausted = errors.New("github rate limit exhausted: safe threshold breached (exit 79)")

// TokenBucketGuard controla o consumo e previne esgotamento de quota de APIs externas (GitHub GraphQL/REST).
type TokenBucketGuard struct {
	mu        sync.RWMutex
	threshold int
	halted    bool
	resetAt   time.Time
}

// NewTokenBucketGuard cria um novo guard com o limiar mínimo de segurança (ex: 100 chamadas restantes).
func NewTokenBucketGuard(threshold int) *TokenBucketGuard {
	if threshold <= 0 {
		threshold = 100
	}
	return &TokenBucketGuard{
		threshold: threshold,
	}
}

// Check avalia a quota restante e o timestamp de reset.
// Se a cota for inferior ao limiar, congela o guard em estado halted e retorna ErrRateLimitExhausted.
func (g *TokenBucketGuard) Check(remaining int, resetUnix int64) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.halted {
		return fmt.Errorf("%w: reset scheduled for %s", ErrRateLimitExhausted, g.resetAt.Format(time.RFC3339))
	}

	if remaining < g.threshold {
		g.halted = true
		g.resetAt = time.Unix(resetUnix, 0)
		return fmt.Errorf("%w: remaining calls %d < threshold %d, resets at %s",
			ErrRateLimitExhausted, remaining, g.threshold, g.resetAt.Format(time.RFC3339))
	}

	return nil
}

// IsHalted retorna true se o guard tiver bloqueado operações para evitar banimento de API.
func (g *TokenBucketGuard) IsHalted() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.halted
}

// Reset restaura o estado do guard caso a cota tenha sido renovada externamente.
func (g *TokenBucketGuard) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.halted = false
	g.resetAt = time.Time{}
}

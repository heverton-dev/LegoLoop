package board

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ClientConfig define os parâmetros de conexão do cliente de GitHub Projects.
type ClientConfig struct {
	BaseURL string
	Token   string
	Guard   *TokenBucketGuard
	HTTPCli *http.Client
}

// Client implementa a comunicação com a API do GitHub/Projects respeitando rate limits.
type Client struct {
	cfg ClientConfig
}

// NewClient instancia um novo client de board.
func NewClient(cfg ClientConfig) *Client {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.github.com"
	}
	if cfg.HTTPCli == nil {
		cfg.HTTPCli = &http.Client{Timeout: 10 * time.Second}
	}
	if cfg.Guard == nil {
		cfg.Guard = NewTokenBucketGuard(100)
	}
	return &Client{cfg: cfg}
}

// FetchReadyCards busca cartões na coluna Ready, auditando a cota remanescente do GitHub.
func (c *Client) FetchReadyCards(ctx context.Context) ([]Card, error) {
	if c.cfg.Guard.IsHalted() {
		return nil, ErrRateLimitExhausted
	}

	reqURL := fmt.Sprintf("%s/cards", strings.TrimRight(c.cfg.BaseURL, "/"))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar requisição de cards: %w", err)
	}

	if c.cfg.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.cfg.Token)
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := c.cfg.HTTPCli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na chamada HTTP de cards: %w", err)
	}
	defer resp.Body.Close()

	// Audita rate limit headers caso existam
	if remainingStr := resp.Header.Get("X-RateLimit-Remaining"); remainingStr != "" {
		remaining, _ := strconv.Atoi(remainingStr)
		resetUnix, _ := strconv.ParseInt(resp.Header.Get("X-RateLimit-Reset"), 10, 64)
		if err := c.cfg.Guard.Check(remaining, resetUnix); err != nil {
			return nil, err
		}
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("resposta inesperada da API (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var allCards []Card
	if err := json.NewDecoder(resp.Body).Decode(&allCards); err != nil {
		return nil, fmt.Errorf("falha ao deserializar resposta de cards: %w", err)
	}

	var readyCards []Card
	for _, card := range allCards {
		if card.Status == StatusReady {
			readyCards = append(readyCards, card)
		}
	}

	return readyCards, nil
}

// MoveCard atualiza o status de um card no quadro.
func (c *Client) MoveCard(ctx context.Context, cardID string, toStatus Status) error {
	if c.cfg.Guard.IsHalted() {
		return ErrRateLimitExhausted
	}

	reqURL := fmt.Sprintf("%s/cards/%s/status", strings.TrimRight(c.cfg.BaseURL, "/"), cardID)
	payload := fmt.Sprintf(`{"status":"%s"}`, toStatus)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, strings.NewReader(payload))
	if err != nil {
		return fmt.Errorf("falha ao criar requisição de mover card: %w", err)
	}

	if c.cfg.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.cfg.Token)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.cfg.HTTPCli.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao mover card: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("falha ao mover card (status %d)", resp.StatusCode)
	}

	return nil
}

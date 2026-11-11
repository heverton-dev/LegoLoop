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

// FetchReadyTickets busca tickets na coluna Ready, auditando a cota remanescente do GitHub.
func (c *Client) FetchReadyTickets(ctx context.Context) ([]Ticket, error) {
	if c.cfg.Guard.IsHalted() {
		return nil, ErrRateLimitExhausted
	}
	resp, err := c.doTicketsRequest(ctx)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := c.auditRateLimit(resp); err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("resposta inesperada da API (%d): %s", resp.StatusCode, string(bodyBytes))
	}
	var all []Ticket
	if err := json.NewDecoder(resp.Body).Decode(&all); err != nil {
		return nil, fmt.Errorf("falha ao deserializar resposta de tickets: %w", err)
	}
	return filterReady(all), nil
}

func (c *Client) doTicketsRequest(ctx context.Context) (*http.Response, error) {
	reqURL := fmt.Sprintf("%s/tickets", strings.TrimRight(c.cfg.BaseURL, "/"))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar requisição de tickets: %w", err)
	}
	if c.cfg.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.cfg.Token)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := c.cfg.HTTPCli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na chamada HTTP de tickets: %w", err)
	}
	return resp, nil
}

// auditRateLimit audita os headers de rate limit caso existam.
func (c *Client) auditRateLimit(resp *http.Response) error {
	remainingStr := resp.Header.Get("X-RateLimit-Remaining")
	if remainingStr == "" {
		return nil
	}
	remaining, _ := strconv.Atoi(remainingStr)
	resetUnix, _ := strconv.ParseInt(resp.Header.Get("X-RateLimit-Reset"), 10, 64)
	return c.cfg.Guard.Check(remaining, resetUnix)
}

func filterReady(tickets []Ticket) []Ticket {
	var ready []Ticket
	for _, t := range tickets {
		if t.Status == StatusReady {
			ready = append(ready, t)
		}
	}
	return ready
}

// MoveTicket atualiza o status de um ticket no quadro.
func (c *Client) MoveTicket(ctx context.Context, ticketID string, toStatus Status) error {
	if c.cfg.Guard.IsHalted() {
		return ErrRateLimitExhausted
	}

	reqURL := fmt.Sprintf("%s/tickets/%s/status", strings.TrimRight(c.cfg.BaseURL, "/"), ticketID)
	payload := fmt.Sprintf(`{"status":"%s"}`, toStatus)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, strings.NewReader(payload))
	if err != nil {
		return fmt.Errorf("falha ao criar requisição de mover ticket: %w", err)
	}

	if c.cfg.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.cfg.Token)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.cfg.HTTPCli.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao mover ticket: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("falha ao mover ticket (status %d)", resp.StatusCode)
	}

	return nil
}

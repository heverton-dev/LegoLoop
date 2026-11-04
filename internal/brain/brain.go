// Package brain é o cérebro plugável da fabricação e da catalogação: escolhe, pelo inventário da máquina, o
// executor que conduz a conversa com o modelo (9Router, Ollama, CLI de harness) e explica a escolha. Nenhum
// fornecedor é requisito: o 9Router deixa de ser obrigatório.
package brain

import (
	"context"
	"net/http"
	"strings"
	"time"

	"legoloop/internal/interview"
	"legoloop/internal/orchestrix"
	"legoloop/internal/provider"
)

// DefaultTimeout cobre a geração de documentos longos (a SPEC); é configurável com WithTimeout.
const DefaultTimeout = 5 * time.Minute

// Brain conduz a conversa com o modelo. Satisfaz interview.LLMClient.
type Brain interface {
	interview.LLMClient
	Name() string           // "9router/code-fast", "ollama/qwen2.5", "claude/sonnet", "agy:default"
	Model() string          // o id do modelo, sem o executor: "code-fast", "qwen2.5", "sonnet"
	Timeout() time.Duration // tempo limite de cada chamada ao modelo
}

type config struct {
	http      *http.Client
	ollamaURL string
	timeout   time.Duration
}

// Option ajusta o cérebro.
type Option func(*config)

// WithOllamaURL aponta o Ollama (padrão: o do provider.DefaultConfig).
func WithOllamaURL(u string) Option { return func(c *config) { c.ollamaURL = u } }

// WithHTTP troca o cliente HTTP dos executores que falam por HTTP.
func WithHTTP(h *http.Client) Option { return func(c *config) { c.http = h } }

// WithTimeout muda o tempo limite de cada chamada ao modelo.
func WithTimeout(d time.Duration) Option { return func(c *config) { c.timeout = d } }

func newConfig(opts []Option) config {
	c := config{ollamaURL: provider.DefaultConfig().OllamaURL, timeout: DefaultTimeout}
	for _, o := range opts {
		o(&c)
	}
	if c.http == nil {
		c.http = &http.Client{Timeout: c.timeout}
	}
	return c
}

// routerBrain fala com o 9Router pelo cliente socrático, com os papéis das mensagens e a chave só para a URL confiável.
type routerBrain struct {
	*interview.RouterLLMClient
	name string
}

func (b routerBrain) Name() string { return b.name }

func (b routerBrain) Model() string { return b.RouterLLMClient.Model }

func (b routerBrain) Timeout() time.Duration { return b.RouterLLMClient.HTTPClient.Timeout }

// orchBrain fala com os demais executores pelo despacho do Orchestrix, que recebe um prompt único.
type orchBrain struct {
	o     *orchestrix.Orchestrix
	src   provider.Source
	model string
	name  string
	limit time.Duration // prazo de cada chamada (vem da configuração, não do cliente HTTP)
}

func (b orchBrain) Name() string { return b.name }

func (b orchBrain) Model() string { return b.model }

func (b orchBrain) Timeout() time.Duration { return b.limit }

func (b orchBrain) Chat(ctx context.Context, msgs []interview.Message) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, b.limit)
	defer cancel()
	return b.o.Ask(ctx, b.src, b.model, flatten(msgs))
}

// flatten junta as mensagens num prompt único, na ordem original, com o papel de cada uma.
func flatten(msgs []interview.Message) string {
	parts := make([]string, 0, len(msgs))
	for _, m := range msgs {
		parts = append(parts, "["+m.Role+"]\n"+m.Content)
	}
	return strings.Join(parts, "\n\n")
}

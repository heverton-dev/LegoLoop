package interview

import (
	"context"
	"time"
)

// LLMClient é a interface que o motor socrático usa para inferência.
type LLMClient interface {
	Chat(ctx context.Context, messages []Message) (string, error)
}

// Message representa um turno de conversa no chat.
type Message struct {
	Role    string `json:"role"`    // "system", "user" ou "assistant"
	Content string `json:"content"` // Texto da mensagem
}

// StatusSession representa o ciclo de vida da entrevista socrática.
type StatusSession string

const (
	StatusInProgress StatusSession = "IN_PROGRESS"
	StatusCompleted  StatusSession = "COMPLETED"
)

// Mirror representa o espelho canônico consolidado ao atingir a certeza.
type Mirror struct {
	Outcome     string   `json:"outcome"`      // Resultado Esperado
	User        string   `json:"user"`         // Quem se Beneficia
	WhyNow      string   `json:"why_now"`      // Por que Agora
	Success     string   `json:"success"`      // Critério de Sucesso
	Constraints []string `json:"constraints"`  // Restrições Inegociáveis
	OutOfScope  []string `json:"out_of_scope"` // Fora de Escopo
}

// TurnResult é a resposta estruturada de cada rodada socrática.
type TurnResult struct {
	SessionID   string        `json:"session_id"`
	Status      StatusSession `json:"status"`
	Hypothesis  string        `json:"hypothesis"`
	Confidence  int           `json:"confidence"` // 0 a 100
	Question    string        `json:"question,omitempty"`
	Guess       string        `json:"guess,omitempty"`
	Mirror      *Mirror       `json:"mirror,omitempty"`
	RawResponse string        `json:"raw_response"`
}

// Session mantém o estado em memória da sessão socrática.
type Session struct {
	ID           string        `json:"id"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Status       StatusSession `json:"status"`
	Theme        string        `json:"theme"`
	Confidence   int           `json:"confidence"`
	Hypothesis   string        `json:"hypothesis"`
	Turns        int           `json:"turns"`
	Messages     []Message     `json:"messages"`
	Mirror       *Mirror       `json:"mirror,omitempty"`
	ArtifactPath string        `json:"artifact_path,omitempty"`
}

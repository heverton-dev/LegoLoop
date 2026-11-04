package fabricar

import (
	"context"
	"errors"
	"strings"

	"legoloop/internal/interview"
)

// AskFunc entrega a pergunta (e o palpite) ao usuário e devolve a resposta. A CLI e o Studio a fornecem.
type AskFunc func(question, guess string) (string, error)

// Clarification é o resultado da Camada 1.
type Clarification struct {
	Theme      string           `json:"theme"`
	Mirror     interview.Mirror `json:"mirror"`
	Confidence int              `json:"confidence"` // 0 no modo autônomo: nenhuma entrevista aconteceu
	Turns      int              `json:"turns"`
	Reached    bool             `json:"reached"` // true se a entrevista chegou a >= 95%
	// Interviewed é verdadeiro quando houve entrevista (modo interativo); no autônomo nada foi medido.
	Interviewed bool `json:"interviewed"`
	// Natureza do processo (vocabulário), quando já conhecida; vazia = natureza ainda não identificada.
	Natureza string `json:"natureza,omitempty"`
	// Gaps lista o que o espelho ainda não respondeu quando a entrevista parou abaixo da meta.
	Gaps []string `json:"gaps,omitempty"`
}

// defaultMaxTurns é o limite de turnos da entrevista interativa.
const defaultMaxTurns = 12

// Clarifier é a Camada 1: modela a intenção. Usa interview.Engine com o LLM injetado; não reimplementa a entrevista.
type Clarifier struct {
	llm      interview.LLMClient
	MaxTurns int
	// OnTurn, se definido, recebe a confiança medida antes de cada pergunta ao usuário.
	OnTurn func(confidence int)
	// Context é o que a entrevista sabe do projeto: objetivo, AGENTS.md e peças aptas.
	Context interview.Context
}

// NewClarifier cria o Clarifier. llm pode ser nil se só o modo autônomo for usado.
func NewClarifier(llm interview.LLMClient) *Clarifier {
	return &Clarifier{llm: llm, MaxTurns: defaultMaxTurns}
}

// Clarify devolve a intenção. Em modo autônomo (ou sem ask) não pergunta nada e usa o Goal ou, sem ele,
// o próprio nome. No interativo, funila com o usuário até >= 95% de confiança ou até MaxTurns.
func (c *Clarifier) Clarify(ctx context.Context, req Request, ask AskFunc) (Clarification, error) {
	if !req.Interactive || ask == nil {
		return autonomous(req), nil
	}
	if c.llm == nil {
		return Clarification{}, errors.New("fabricar: o modo interativo exige um LLM")
	}
	return c.interview(ctx, req, ask)
}

func autonomous(req Request) Clarification {
	name := strings.TrimSpace(req.RawName)
	goal := strings.TrimSpace(req.Goal)
	if goal == "" {
		goal = name
	}
	return Clarification{Theme: name, Mirror: interview.Mirror{Outcome: goal}}
}

func (c *Clarifier) interview(ctx context.Context, req Request, ask AskFunc) (Clarification, error) {
	idea := strings.TrimSpace(req.Goal)
	if idea == "" {
		idea = strings.TrimSpace(req.RawName)
	}
	eng := interview.NewEngine(c.llm).WithContext(c.Context)
	sess, turn, err := eng.Start(ctx, idea)
	if err != nil {
		return Clarification{}, err
	}
	turns := 0
	for turn.Status != interview.StatusCompleted && turns < c.maxTurns() {
		if c.OnTurn != nil {
			c.OnTurn(turn.Confidence)
		}
		answer, err := ask(turn.Question, turn.Guess)
		if err != nil {
			return Clarification{}, err
		}
		turns++
		if turn, err = eng.Step(ctx, sess.ID, answer); err != nil {
			return Clarification{}, err
		}
	}
	return outcome(idea, turn, turns), nil
}

func (c *Clarifier) maxTurns() int {
	if c.MaxTurns <= 0 {
		return defaultMaxTurns
	}
	return c.MaxTurns
}

// ClarificationFromSession transforma uma entrevista salva (feita no Chat, passo a passo) na Camada 1 pronta.
// A confiança é a medida na sessão; o espelho só existe se a entrevista o consolidou.
func ClarificationFromSession(s *interview.Session) Clarification {
	turn := &interview.TurnResult{Status: s.Status, Hypothesis: s.Hypothesis, Confidence: s.Confidence, Mirror: s.Mirror}
	return outcome(s.Theme, turn, s.Turns)
}

func outcome(idea string, turn *interview.TurnResult, turns int) Clarification {
	m := interview.Mirror{Outcome: turn.Hypothesis}
	if turn.Mirror != nil {
		m = *turn.Mirror
	}
	if strings.TrimSpace(m.Outcome) == "" {
		m.Outcome = idea
	}
	reached := turn.Status == interview.StatusCompleted && turn.Confidence >= ConfidenceTarget
	clar := Clarification{Theme: idea, Mirror: m, Confidence: turn.Confidence, Turns: turns, Reached: reached, Interviewed: true}
	if !reached {
		clar.Gaps = mirrorGaps(m)
	}
	return clar
}

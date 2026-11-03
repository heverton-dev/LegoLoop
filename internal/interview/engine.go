package interview

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Engine orquestra as sessões da entrevista socrática.
type Engine struct {
	llm      LLMClient
	mu       sync.RWMutex
	sessions map[string]*Session
}

// NewEngine inicializa o motor com o cliente LLM.
func NewEngine(llm LLMClient) *Engine {
	return &Engine{
		llm:      llm,
		sessions: make(map[string]*Session),
	}
}

// Start inicia uma nova sessão a partir de uma ideia inicial.
func (e *Engine) Start(ctx context.Context, initialIdea string) (*Session, *TurnResult, error) {
	if strings.TrimSpace(initialIdea) == "" {
		return nil, nil, errors.New("a ideia inicial não pode ser vazia")
	}

	sessionID := fmt.Sprintf("intv_%d", time.Now().UnixNano())
	now := time.Now()

	sess := &Session{
		ID:         sessionID,
		CreatedAt:  now,
		UpdatedAt:  now,
		Status:     StatusInProgress,
		Theme:      initialIdea,
		Confidence: 10,
		Turns:      1,
		Messages: []Message{
			{Role: "system", Content: SystemPromptCanonical},
			{Role: "user", Content: fmt.Sprintf("Minha ideia inicial é: %s", initialIdea)},
		},
	}

	rawResp, err := e.llm.Chat(ctx, sess.Messages)
	if err != nil {
		return nil, nil, fmt.Errorf("falha ao invocar LLM: %w", err)
	}

	sess.Messages = append(sess.Messages, Message{Role: "assistant", Content: rawResp})
	turn := ParseTurnResult(rawResp)
	turn.SessionID = sessionID

	sess.Confidence = turn.Confidence
	sess.Hypothesis = turn.Hypothesis
	if turn.Status == StatusCompleted {
		sess.Status = StatusCompleted
		sess.Mirror = turn.Mirror
	}

	e.mu.Lock()
	e.sessions[sessionID] = sess
	e.mu.Unlock()

	return sess, turn, nil
}

// Step avança um turno na sessão recebendo a resposta do usuário.
func (e *Engine) Step(ctx context.Context, sessionID, answer string) (*TurnResult, error) {
	if strings.TrimSpace(answer) == "" {
		return nil, errors.New("a resposta não pode ser vazia")
	}

	e.mu.Lock()
	sess, exists := e.sessions[sessionID]
	e.mu.Unlock()

	if !exists {
		return nil, fmt.Errorf("sessão %s não encontrada", sessionID)
	}
	if sess.Status == StatusCompleted {
		return nil, errors.New("esta entrevista já foi finalizada")
	}

	sess.Messages = append(sess.Messages, Message{Role: "user", Content: answer})
	sess.Turns++
	sess.UpdatedAt = time.Now()

	rawResp, err := e.llm.Chat(ctx, sess.Messages)
	if err != nil {
		return nil, fmt.Errorf("falha ao invocar LLM: %w", err)
	}

	sess.Messages = append(sess.Messages, Message{Role: "assistant", Content: rawResp})
	turn := ParseTurnResult(rawResp)
	turn.SessionID = sessionID

	sess.Confidence = turn.Confidence
	sess.Hypothesis = turn.Hypothesis
	if turn.Status == StatusCompleted {
		sess.Status = StatusCompleted
		sess.Mirror = turn.Mirror
	}

	return turn, nil
}

// GetSession retorna os dados de uma sessão existente.
func (e *Engine) GetSession(sessionID string) (*Session, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	sess, exists := e.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("sessão %s não encontrada", sessionID)
	}
	return sess, nil
}

// MaterializeIntent grava o arquivo docs/intents/YYYY-MM-DD_INTENT_*.md seguindo o TEMPLATE_INTENT.md.
func (e *Engine) MaterializeIntent(sessionID, baseDir, fileName string) (string, error) {
	sess, err := e.GetSession(sessionID)
	if err != nil {
		return "", err
	}
	if sess.Mirror == nil {
		return "", errors.New("impossível materializar: espelho socrático ainda não consolidado")
	}

	date := time.Now().Format("2006-01-02")
	filePath, fileName := intentTarget(baseDir, fileName, sess.Theme, date)
	content := RenderIntent(IntentView{
		Theme:      sess.Theme,
		Date:       date,
		SessionID:  sess.ID,
		Confidence: sess.Confidence,
		Mirror:     sess.Mirror,
		SpecRel:    SpecRelFor(fileName, sess.Theme, date),
	})
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("falha ao gravar artefato de intenção: %w", err)
	}

	sess.ArtifactPath = filePath
	return filePath, nil
}

// intentTarget resolve o nome do arquivo (gerando um a partir do tema e da data) e a pasta docs/intents.
func intentTarget(baseDir, fileName, theme, date string) (path, name string) {
	if fileName == "" {
		fileName = fmt.Sprintf("%s_INTENT_%s.md", date, sanitizeFileName(theme))
	}
	dir := filepath.Join(baseDir, "docs", "intents")
	if err := os.MkdirAll(dir, 0755); err != nil {
		dir = baseDir // baseDir plano (diretório temporário de teste)
	}
	return filepath.Join(dir, fileName), fileName
}

var (
	reHypothesis = regexp.MustCompile(`(?i)HIPÓTESE:\s*([^\n]+)`)
	reConfidence = regexp.MustCompile(`(?i)CONFIANÇA:\s*(\d+)%?`)
	reQuestion   = regexp.MustCompile(`(?i)PERGUNTA:\s*([^\n]+)`)
	reGuess      = regexp.MustCompile(`(?i)PALPITE:\s*([^\n]+)`)

	reOutcome     = regexp.MustCompile(`(?i)- Resultado Esperado:\s*([^\n]+)`)
	reUser        = regexp.MustCompile(`(?i)- Quem se Beneficia:\s*([^\n]+)`)
	reWhyNow      = regexp.MustCompile(`(?i)- Por que Agora:\s*([^\n]+)`)
	reSuccess     = regexp.MustCompile(`(?i)- Critério de Sucesso:\s*([^\n]+)`)
	reConstraints = regexp.MustCompile(`(?i)- Restrições Inegociáveis:\s*([^\n]+)`)
	reOutOfScope  = regexp.MustCompile(`(?i)- Fora de Escopo:\s*([^\n]+)`)
)

// ParseTurnResult extrai os campos estruturados da resposta em linguagem natural da LLM.
func ParseTurnResult(raw string) *TurnResult {
	turn := &TurnResult{
		Status:      StatusInProgress,
		RawResponse: raw,
		Confidence:  50,
	}

	if m := reHypothesis.FindStringSubmatch(raw); len(m) > 1 {
		turn.Hypothesis = strings.TrimSpace(m[1])
	}
	if m := reConfidence.FindStringSubmatch(raw); len(m) > 1 {
		if val, err := strconv.Atoi(m[1]); err == nil {
			turn.Confidence = val
		}
	}
	if m := reQuestion.FindStringSubmatch(raw); len(m) > 1 {
		turn.Question = strings.TrimSpace(m[1])
	}
	if m := reGuess.FindStringSubmatch(raw); len(m) > 1 {
		turn.Guess = strings.TrimSpace(m[1])
	}

	// Se houver espelho canônico ou confiança >= 95%
	if strings.Contains(strings.ToUpper(raw), "ESPELHO CANÔNICO") || turn.Confidence >= 95 {
		mirror := &Mirror{}
		if m := reOutcome.FindStringSubmatch(raw); len(m) > 1 {
			mirror.Outcome = strings.TrimSpace(m[1])
		}
		if m := reUser.FindStringSubmatch(raw); len(m) > 1 {
			mirror.User = strings.TrimSpace(m[1])
		}
		if m := reWhyNow.FindStringSubmatch(raw); len(m) > 1 {
			mirror.WhyNow = strings.TrimSpace(m[1])
		}
		if m := reSuccess.FindStringSubmatch(raw); len(m) > 1 {
			mirror.Success = strings.TrimSpace(m[1])
		}
		if m := reConstraints.FindStringSubmatch(raw); len(m) > 1 {
			for _, item := range strings.Split(m[1], ",") {
				mirror.Constraints = append(mirror.Constraints, strings.TrimSpace(item))
			}
		}
		if m := reOutOfScope.FindStringSubmatch(raw); len(m) > 1 {
			for _, item := range strings.Split(m[1], ",") {
				mirror.OutOfScope = append(mirror.OutOfScope, strings.TrimSpace(item))
			}
		}

		if mirror.Outcome != "" || turn.Confidence >= 95 {
			turn.Mirror = mirror
			turn.Status = StatusCompleted
			if turn.Confidence < 95 {
				turn.Confidence = 95
			}
		}
	}

	return turn
}

func sanitizeFileName(s string) string {
	s = strings.ToUpper(s)
	reg := regexp.MustCompile(`[^A-Z0-9]+`)
	s = reg.ReplaceAllString(s, "_")
	s = strings.Trim(s, "_")
	if len(s) > 30 {
		s = s[:30]
	}
	if s == "" {
		s = "NOVA_INTENCAO"
	}
	return s
}

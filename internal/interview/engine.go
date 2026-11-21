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
	ctx      Context
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
	now := time.Now()
	sess := &Session{
		ID:         fmt.Sprintf("intv_%d", now.UnixNano()),
		CreatedAt:  now,
		UpdatedAt:  now,
		Status:     StatusInProgress,
		Theme:      initialIdea,
		Confidence: 10,
		Turns:      1,
		Messages: []Message{
			{Role: "system", Content: e.systemPrompt()},
			{Role: "user", Content: fmt.Sprintf("Minha ideia inicial é: %s", initialIdea)},
		},
	}
	turn, err := e.exchange(ctx, sess)
	if err != nil {
		return nil, nil, err
	}
	e.mu.Lock()
	e.sessions[sess.ID] = sess
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
	return e.exchange(ctx, sess)
}

// exchange envia as mensagens ao cérebro, registra a resposta e atualiza a sessão com o turno lido.
func (e *Engine) exchange(ctx context.Context, sess *Session) (*TurnResult, error) {
	rawResp, err := e.llm.Chat(ctx, sess.Messages)
	if err != nil {
		return nil, fmt.Errorf("falha ao invocar LLM: %w", err)
	}
	sess.Messages = append(sess.Messages, Message{Role: "assistant", Content: rawResp})
	turn := ParseTurnResult(rawResp)
	if emImpasse(turn) {
		if turn, err = e.pedirPremonicao(ctx, sess, turn); err != nil {
			return nil, err
		}
	}
	turn.SessionID = sess.ID
	if !turn.ConfidenceReported {
		turn.Confidence = sess.Confidence // sem medida nova, vale a última; nunca um número inventado
	}
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
	// \*{0,2} tolera o realce em Markdown (**CONFIANÇA**:) que os cérebros emitem.
	reHypothesis = regexp.MustCompile(`(?i)\*{0,2}HIPÓTESE\*{0,2}:\*{0,2}\s*([^\n]+)`)
	reConfidence = regexp.MustCompile(`(?i)\*{0,2}CONFIANÇA\*{0,2}:\*{0,2}\s*(\d+)%?`)
	reQuestion   = regexp.MustCompile(`(?i)\*{0,2}PERGUNTA\*{0,2}:\*{0,2}\s*([^\n]+)`)
	reGuess      = regexp.MustCompile(`(?i)\*{0,2}PALPITE\*{0,2}:\*{0,2}\s*([^\n]+)`)
)

// capture devolve o primeiro grupo de re em raw, sem espaços nas pontas, e se houve casamento.
func capture(re *regexp.Regexp, raw string) (string, bool) {
	if m := re.FindStringSubmatch(raw); len(m) > 1 {
		return strings.TrimSpace(m[1]), true
	}
	return "", false
}

// captureLast devolve a última ocorrência: o cérebro costuma repetir a confiança ao final do espelho.
func captureLast(re *regexp.Regexp, raw string) (string, bool) {
	all := re.FindAllStringSubmatch(raw, -1)
	if len(all) == 0 {
		return "", false
	}
	m := all[len(all)-1]
	if len(m) < 2 {
		return "", false
	}
	return strings.TrimSpace(m[1]), true
}

// splitItems separa uma lista escrita em uma linha, com itens separados por vírgula.
func splitItems(s string) []string {
	var out []string
	for _, item := range strings.Split(s, ",") {
		out = append(out, strings.TrimSpace(item))
	}
	return out
}

// ParseTurnResult extrai os campos estruturados da resposta em linguagem natural da LLM.
func ParseTurnResult(raw string) *TurnResult {
	turn := &TurnResult{
		Status:      StatusInProgress,
		RawResponse: raw,
		Confidence:  50,
	}
	turn.Hypothesis, _ = capture(reHypothesis, raw)
	if s, ok := captureLast(reConfidence, raw); ok {
		if val, err := strconv.Atoi(s); err == nil {
			turn.Confidence, turn.ConfidenceReported = val, true
		}
	}
	turn.Question, _ = capture(reQuestion, raw)
	turn.Guess, _ = capture(reGuess, raw)
	// Se houver espelho canônico ou confiança >= 95%
	if strings.Contains(strings.ToUpper(raw), "ESPELHO CANÔNICO") || turn.Confidence >= 95 {
		completeTurn(turn, parseMirror(raw), hasThreePremonitions(raw))
	}
	return turn
}

// parseMirror lê os seis campos do espelho canônico.
func parseMirror(raw string) *Mirror {
	c := camposEspelho(raw)
	return &Mirror{
		Outcome:     strings.Join(c["resultado esperado"], " "),
		User:        strings.Join(c["quem se beneficia"], " "),
		WhyNow:      strings.Join(c["por que agora"], " "),
		Success:     strings.Join(c["critério de sucesso"], "; "),
		Constraints: itensDoCampo(c["restrições inegociáveis"]),
		OutOfScope:  itensDoCampo(c["fora de escopo"]),
	}
}

var rePremonicaoItem = regexp.MustCompile(`(?im)^\s*[1-3]\.\s*Pergunta:\s*[^\n|]+\|\s*Resposta prevista:\s*[^\n]+`)

// hasThreePremonitions verifica se há seção PREMONIÇÃO com exatamente 3 perguntas e respostas previstas.
func hasThreePremonitions(raw string) bool {
	if !strings.Contains(strings.ToUpper(raw), "PREMONIÇÃO") {
		return false
	}
	return len(rePremonicaoItem.FindAllString(raw, -1)) >= 3
}

// completeTurn dá o turno por concluído quando a confiança chegou a 95% ou quando há espelho + premonição das 3 perguntas.
// Sem Espelho a sessão nunca conclui: a INTENT só materializa a partir dele.
func completeTurn(turn *TurnResult, mirror *Mirror, premonicaoOK bool) {
	if mirror.Outcome == "" {
		return
	}
	turn.Mirror = mirror
	if turn.Confidence < 95 && !premonicaoOK {
		return
	}
	turn.Status = StatusCompleted
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

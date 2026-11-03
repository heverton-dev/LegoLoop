package interview_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"legoloop/internal/interview"
)

type mockLLM struct {
	responses []string
	callCount int
}

func (m *mockLLM) Chat(ctx context.Context, msgs []interview.Message) (string, error) {
	if m.callCount >= len(m.responses) {
		return "HIPÓTESE: Nenhuma\nCONFIANÇA: 50%\nPERGUNTA: Mais detalhes?\nPALPITE: Sim", nil
	}
	resp := m.responses[m.callCount]
	m.callCount++
	return resp, nil
}

func TestInterviewEngine_FullFlow(t *testing.T) {
	mockResp1 := `HIPÓTESE: O usuário deseja criar uma esteira de automação de testes com notificações.
CONFIANÇA: 60%
PERGUNTA: Quais ferramentas de teste específicas você pretende executar no pipeline?
PALPITE: go test ./... e golangci-lint`

	mockResp2 := `HIPÓTESE: O usuário quer uma esteira Go completa com testes e notificações via Slack.
CONFIANÇA: 85%
PERGUNTA: As notificações devem ser disparadas apenas em caso de falha ou também no sucesso?
PALPITE: Apenas em caso de falha para evitar ruído`

	mockResp3 := `HIPÓTESE: Esteira Go automatizada rodando testes/linter com alerta no Slack apenas em caso de falha.
CONFIANÇA: 98%
ESPELHO CANÔNICO:
- Resultado Esperado: Pipeline automatizado de validação Go com webhook do Slack
- Quem se Beneficia: Desenvolvedores da equipe de backend
- Por que Agora: Reduzir tempo de revisão e evitar regressões em produção
- Critério de Sucesso: Testes rodam em < 2min e falhas geram alerta no canal #alerts
- Restrições Inegociáveis: Sem subagentes, contexto estéril, segredos apenas via variáveis de ambiente
- Fora de Escopo: Deploy em homologação ou produção`

	mock := &mockLLM{
		responses: []string{mockResp1, mockResp2, mockResp3},
	}

	engine := interview.NewEngine(mock)

	// Turno 1: Iniciar
	sess, turn1, err := engine.Start(context.Background(), "Quero automatizar os testes do meu projeto Go")
	if err != nil {
		t.Fatalf("Start falhou: %v", err)
	}
	if sess.Status != interview.StatusInProgress {
		t.Errorf("esperava StatusInProgress, obteve %s", sess.Status)
	}
	if turn1.Confidence != 60 {
		t.Errorf("esperava confiança 60%%, obteve %d%%", turn1.Confidence)
	}
	if !strings.Contains(turn1.Question, "ferramentas de teste") {
		t.Errorf("pergunta inesperada: %s", turn1.Question)
	}

	// Turno 2: Responder
	turn2, err := engine.Step(context.Background(), sess.ID, "go test e golangci-lint, avisando no slack")
	if err != nil {
		t.Fatalf("Step 1 falhou: %v", err)
	}
	if turn2.Confidence != 85 {
		t.Errorf("esperava confiança 85%%, obteve %d%%", turn2.Confidence)
	}

	// Turno 3: Responder e atingir confiança >= 95%
	turn3, err := engine.Step(context.Background(), sess.ID, "Apenas falhas no canal #alerts")
	if err != nil {
		t.Fatalf("Step 2 falhou: %v", err)
	}
	if turn3.Status != interview.StatusCompleted {
		t.Errorf("esperava StatusCompleted, obteve %s", turn3.Status)
	}
	if turn3.Confidence < 95 {
		t.Errorf("esperava confiança >= 95%%, obteve %d%%", turn3.Confidence)
	}
	if turn3.Mirror == nil {
		t.Fatal("esperava Mirror preenchido na conclusão")
	}
	if turn3.Mirror.Outcome != "Pipeline automatizado de validação Go com webhook do Slack" {
		t.Errorf("Outcome inesperado: %s", turn3.Mirror.Outcome)
	}

	// Testar Materialização do INTENT.md
	tempDir := t.TempDir()
	path, err := engine.MaterializeIntent(sess.ID, tempDir, "2026-11-02_INTENT_ESTEIRA_TESTES.md")
	if err != nil {
		t.Fatalf("MaterializeIntent falhou: %v", err)
	}
	if !strings.HasSuffix(path, "2026-11-02_INTENT_ESTEIRA_TESTES.md") {
		t.Errorf("caminho de arquivo inesperado: %s", path)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("erro ao ler arquivo gerado: %v", err)
	}
	cStr := string(content)
	if !strings.Contains(cStr, "# 🎯 INTENT:") {
		t.Errorf("arquivo não contém título INTENT: %s", cStr)
	}
	if !strings.Contains(cStr, "Pipeline automatizado de validação Go com webhook do Slack") {
		t.Errorf("arquivo não contém o Outcome esperado: %s", cStr)
	}
	if !strings.Contains(cStr, "STOP-GATE 01") {
		t.Errorf("arquivo não contém Stop-Gate 01: %s", cStr)
	}
}

func TestInterviewEngine_ParseTurnResult(t *testing.T) {
	raw := `HIPÓTESE:   Refatoração do canvas visual.
CONFIANÇA: 75%
PERGUNTA:   O canvas deve usar SVG puro ou Canvas HTML5?
PALPITE:    SVG puro pela facilidade de eventos DOM.`

	turn := interview.ParseTurnResult(raw)
	if turn.Confidence != 75 {
		t.Errorf("esperava confiança 75, obteve %d", turn.Confidence)
	}
	if turn.Hypothesis != "Refatoração do canvas visual." {
		t.Errorf("hipótese inesperada: %s", turn.Hypothesis)
	}
	if turn.Question != "O canvas deve usar SVG puro ou Canvas HTML5?" {
		t.Errorf("pergunta inesperada: %s", turn.Question)
	}
	if turn.Guess != "SVG puro pela facilidade de eventos DOM." {
		t.Errorf("palpite inesperado: %s", turn.Guess)
	}
}

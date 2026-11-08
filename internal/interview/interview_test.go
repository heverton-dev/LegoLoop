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

func TestParseTurnResult_MarkdownEmDestaque(t *testing.T) {
	raw := "**CONFIANÇA**: 60%\n\n**PALPITE**: pipeline 100% local\n\n" +
		"**PERGUNTA**: Você pretende usar algum framework específico para a geração?"

	turn := interview.ParseTurnResult(raw)
	if !turn.ConfidenceReported || turn.Confidence != 60 {
		t.Errorf("esperava confiança 60%% reportada, obteve %d (reportado=%v)", turn.Confidence, turn.ConfidenceReported)
	}
	if strings.TrimSpace(turn.Question) == "" {
		t.Errorf("pergunta não extraída de resposta em markdown: %q", turn.RawResponse)
	}
	if strings.TrimSpace(turn.Guess) == "" {
		t.Errorf("palpite não extraído de resposta em markdown")
	}
}

func TestParseTurnResult_EspelhoEmMarkdown(t *testing.T) {
	raw := "**CONFIANÇA**: 98%\n\n**ESPELHO CANÔNICO**:\n" +
		"- **Resultado Esperado**: gerar vídeos 100% locais\n" +
		"- **Quem se Beneficia**: criadores de PC fraco\n" +
		"- **Por que Agora**: IAs pagas dispararam preço\n" +
		"- **Critério de Sucesso**: rodar offline sem API\n" +
		"- **Restrições Inegociáveis**: sem nuvem\n" +
		"- **Fora de Escopo**: treinar modelos do zero"

	turn := interview.ParseTurnResult(raw)
	if turn.Status != interview.StatusCompleted {
		t.Fatalf("esperava status COMPLETED, obteve %q", turn.Status)
	}
	if turn.Mirror == nil {
		t.Fatal("espelho canônico não extraído")
	}
	if turn.Mirror.Outcome != "gerar vídeos 100% locais" {
		t.Errorf("resultado esperado extraído: %q", turn.Mirror.Outcome)
	}
	if turn.Mirror.Constraints[0] != "sem nuvem" {
		t.Errorf("restrições extraídas: %v", turn.Mirror.Constraints)
	}
}

func TestParseTurnResult_EspelhoComBulletEstrela(t *testing.T) {
	raw := "**HIPÓTESE**: pipeline local\n\n**CONFIANÇA**: 90%\n\n**PERGUNTA**: falta algo?\n\n" +
		"**ESPELHO CANÔNICO**:\n\n" +
		"* **Resultado Esperado**: pipeline offline 1080p\n" +
		"* **Quem se Beneficia**: criadores sem GPU\n" +
		"* **Por que Agora**: IA pagas encareceram\n" +
		"* **Critério de Sucesso**: MP4/H.264 leve e offline\n" +
		"* **Restrições Inegociáveis**: sem nuvem, só CPU\n" +
		"* **Fora de Escopo**: treinar modelos\n\n" +
		"**CONFIANÇA**: 95% (ou mais)"

	turn := interview.ParseTurnResult(raw)
	if turn.Status != interview.StatusCompleted {
		t.Fatalf("esperava COMPLETED com espelho completo, obteve %q (conf=%d)", turn.Status, turn.Confidence)
	}
	if turn.Confidence < 95 {
		t.Errorf("esperava última confiança 95%%, obteve %d", turn.Confidence)
	}
	if turn.Mirror == nil || turn.Mirror.Outcome != "pipeline offline 1080p" {
		t.Errorf("espelho não extraído: %+v", turn.Mirror)
	}
	if turn.Mirror != nil && len(turn.Mirror.Constraints) == 0 {
		t.Errorf("restrições não extraídas")
	}
}

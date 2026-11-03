package studio

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"

	"legoloop/internal/interview"
)

type mockStudioLLM struct {
	responses []string
	callCount int
}

func (m *mockStudioLLM) Chat(ctx context.Context, msgs []interview.Message) (string, error) {
	if m.callCount >= len(m.responses) {
		return "HIPÓTESE: Nenhuma\nCONFIANÇA: 50%\nPERGUNTA: Detalhes?\nPALPITE: Sim", nil
	}
	resp := m.responses[m.callCount]
	m.callCount++
	return resp, nil
}

func TestInterviewAPI_Endpoints(t *testing.T) {
	s, _ := newServer(t)

	mock := &mockStudioLLM{
		responses: []string{
			"HIPÓTESE: Criar módulo de pagamento.\nCONFIANÇA: 60%\nPERGUNTA: Pix ou Cartão?\nPALPITE: Pix",
			"HIPÓTESE: Módulo de pagamento exclusivo Pix com webhook.\nCONFIANÇA: 98%\nESPELHO CANÔNICO:\n- Resultado Esperado: Integração Pix\n- Quem se Beneficia: Clientes\n- Por que Agora: Aumentar conversão\n- Critério de Sucesso: Cobranças geradas e pagas com QR Code\n- Restrições Inegociáveis: Sem log de chaves privadas\n- Fora de Escopo: Boleto e Cartão",
		},
	}
	s.Interview = interview.NewEngine(mock)
	h := s.Handler()

	// 1. Iniciar entrevista
	startBody := `{"idea":"Quero criar integração de pagamento"}`
	rec := do(t, h, http.MethodPost, "/api/interview/start", startBody)
	if rec.Code != http.StatusOK {
		t.Fatalf("esperava 200 em start, obteve %d: %s", rec.Code, rec.Body.String())
	}

	var startResp struct {
		Session interview.Session    `json:"session"`
		Turn    interview.TurnResult `json:"turn"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &startResp); err != nil {
		t.Fatalf("erro ao unmarshal: %v", err)
	}
	if startResp.Session.ID == "" {
		t.Fatal("session ID vazio")
	}
	if startResp.Turn.Confidence != 60 {
		t.Errorf("esperava 60%% confiança, obteve %d", startResp.Turn.Confidence)
	}

	// 2. Responder pergunta
	msgBody := `{"session_id":"` + startResp.Session.ID + `","answer":"Exclusivamente Pix"}`
	rec = do(t, h, http.MethodPost, "/api/interview/message", msgBody)
	if rec.Code != http.StatusOK {
		t.Fatalf("esperava 200 em message, obteve %d: %s", rec.Code, rec.Body.String())
	}

	var msgResp struct {
		Turn interview.TurnResult `json:"turn"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &msgResp); err != nil {
		t.Fatalf("erro ao unmarshal message: %v", err)
	}
	if msgResp.Turn.Confidence != 98 {
		t.Errorf("esperava 98%% de confiança, obteve %d", msgResp.Turn.Confidence)
	}
	if msgResp.Turn.Status != interview.StatusCompleted {
		t.Errorf("esperava status COMPLETED, obteve %s", msgResp.Turn.Status)
	}

	// 3. Finalizar e materializar
	finBody := `{"session_id":"` + startResp.Session.ID + `","file_name":"2026-11-02_INTENT_PAGAMENTO_PIX.md"}`
	rec = do(t, h, http.MethodPost, "/api/interview/finalize", finBody)
	if rec.Code != http.StatusOK {
		t.Fatalf("esperava 200 em finalize, obteve %d: %s", rec.Code, rec.Body.String())
	}

	var finResp struct {
		ArtifactPath string `json:"artifact_path"`
		Status       string `json:"status"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &finResp); err != nil {
		t.Fatalf("erro ao unmarshal finalize: %v", err)
	}
	if !strings.HasSuffix(finResp.ArtifactPath, "2026-11-02_INTENT_PAGAMENTO_PIX.md") {
		t.Errorf("caminho retornado inesperado: %s", finResp.ArtifactPath)
	}

	content, err := os.ReadFile(finResp.ArtifactPath)
	if err != nil {
		t.Fatalf("arquivo não encontrado no disco: %v", err)
	}
	if !strings.Contains(string(content), "Integração Pix") {
		t.Errorf("conteúdo não contém o espelho canônico esperado: %s", string(content))
	}
}

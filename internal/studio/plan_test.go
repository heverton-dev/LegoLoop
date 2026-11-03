package studio

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/interview"
	"legoloop/internal/plangen"
)

type mockStudioPlanLLM struct {
	output string
}

func (m *mockStudioPlanLLM) Chat(ctx context.Context, msgs []interview.Message) (string, error) {
	return m.output, nil
}

func TestStudioPlanAPI_Endpoint(t *testing.T) {
	s, _ := newServer(t)

	mockPlanOutput := `
### [TICKET-E01.1] Conectar websocket do studio
- **Épico:** E01
- **Natureza:** micro
- **Modo:** tdd_red_green
- **Dependências:** nenhuma
- **Arquivos Afetados:** internal/studio/ws.go
- **Critérios Atendidos:** AC-E01.1
- **Comando de Validação:** go test ./internal/studio/...
`

	mockLLM := &mockStudioPlanLLM{output: mockPlanOutput}
	s.PlanGen = plangen.NewGenerator(mockLLM)
	h := s.Handler()

	specFile := filepath.Join(s.ProjectDir, "2026-11-02_SPEC_WS.md")
	os.WriteFile(specFile, []byte("# SPEC WS"), 0644)

	reqBody := `{"spec_path":"` + filepath.ToSlash(specFile) + `","supervision_level":"micro"}`
	rec := do(t, h, http.MethodPost, "/api/plan/generate", reqBody)
	if rec.Code != http.StatusOK {
		t.Fatalf("esperava 200 em plan/generate, obteve %d: %s", rec.Code, rec.Body.String())
	}

	var res plangen.PlanResult
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("erro ao decodificar JSON: %v", err)
	}

	if len(res.Tickets) != 1 {
		t.Fatalf("esperava 1 ticket, obteve %d", len(res.Tickets))
	}
	if res.SupervisionLevel != plangen.SupervisionMicro {
		t.Errorf("supervision_level inesperado: %s", res.SupervisionLevel)
	}
}

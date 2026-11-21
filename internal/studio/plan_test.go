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
	os.WriteFile(specFile, []byte("# SPEC WS\n\n> **Referência da Intenção:** [docs/intents/2026-11-02_INTENT_WS.md](../intents/2026-11-02_INTENT_WS.md)\n"), 0644)
	intentWS := filepath.Join(s.ProjectDir, "docs", "intents", "2026-11-02_INTENT_WS.md")
	os.MkdirAll(filepath.Dir(intentWS), 0755)
	os.WriteFile(intentWS, []byte("---\nid: INTENT-WS\nstatus: APPROVED\n---\n# INTENT WS\n"), 0644)

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

// TestStudioPlanAPI_SpecSemIntencaoDevolve422: sem INTENT aprovada, o PLAN não nasce e o erro é do processo (422).
func TestStudioPlanAPI_SpecSemIntencaoDevolve422(t *testing.T) {
	s, _ := newServer(t)
	s.PlanGen = plangen.NewGenerator(&mockStudioPlanLLM{output: "nunca chamado"})
	specFile := filepath.Join(s.ProjectDir, "2026-11-02_SPEC_SEM_INTENT.md")
	if err := os.WriteFile(specFile, []byte("# SPEC sem referência à INTENT\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	rec := do(t, s.Handler(), http.MethodPost, "/api/plan/generate", `{"spec_path":"`+filepath.ToSlash(specFile)+`","supervision_level":"micro"}`)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("SPEC sem INTENT deveria devolver 422, devolveu %d: %s", rec.Code, rec.Body.String())
	}
}

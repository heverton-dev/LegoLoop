package plangen_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/interview"
	"legoloop/internal/plangen"
)

type mockPlanLLM struct {
	output string
}

func (m *mockPlanLLM) Chat(ctx context.Context, msgs []interview.Message) (string, error) {
	return m.output, nil
}

func TestPlanGenerator_Success(t *testing.T) {
	mockPlanOutput := `
### [TICKET-E01.1] Criar tipos e contratos de serviço
- **Épico:** E01
- **Natureza:** micro
- **Modo:** tdd_red_green
- **Dependências:** nenhuma
- **Arquivos Afetados:** internal/service/types.go, internal/service/types_test.go
- **Critérios Atendidos:** AC-E01.1
- **Comando de Validação:** go test -v ./internal/service/...

### [TICKET-E01.2] Ajustar formatação e lint do serviço
- **Épico:** E01
- **Natureza:** nano
- **Modo:** sequential
- **Dependências:** TICK-E01.1
- **Arquivos Afetados:** internal/service/types.go
- **Critérios Atendidos:** AC-E01.2
- **Comando de Validação:** golangci-lint run ./internal/service/...
`

	tempDir := t.TempDir()
	specFile := filepath.Join(tempDir, "2026-11-02_SPEC_SERVICO.md")
	os.WriteFile(specFile, []byte("# SPEC Servico"), 0644)

	llm := &mockPlanLLM{output: mockPlanOutput}
	gen := plangen.NewGenerator(llm)

	// Teste com nível de supervisão CIRÚRGICO (Micro Gate)
	res, err := gen.Generate(context.Background(), tempDir, specFile, plangen.SupervisionMicro)
	if err != nil {
		t.Fatalf("Generate falhou: %v", err)
	}

	if len(res.Tickets) != 2 {
		t.Fatalf("esperava 2 tickets, obteve %d", len(res.Tickets))
	}

	t1 := res.Tickets[0]
	if t1.ID != "TICK-01.1" || t1.Nature != "micro" || !t1.HumanGateAfter {
		t.Errorf("ticket 1 incorreto: %+v", t1)
	}

	t2 := res.Tickets[1]
	if t2.ID != "TICK-01.2" || t2.Nature != "nano" || t2.HumanGateAfter {
		t.Errorf("ticket 2 incorreto (nano não deve ter gate em nível micro): %+v", t2)
	}

	// Verificar PLAN.md gerado
	planBytes, err := os.ReadFile(res.PlanMarkdownPath)
	if err != nil {
		t.Fatalf("erro ao ler PLAN.md: %v", err)
	}
	if !strings.Contains(string(planBytes), "STOP-GATE 03") {
		t.Error("PLAN.md não contém Stop-Gate 03")
	}

	// Verificar ROUTE.json gerado
	routeBytes, err := os.ReadFile(res.RouteJSONPath)
	if err != nil {
		t.Fatalf("erro ao ler ROUTE.json: %v", err)
	}

	var routeMap map[string]any
	if err := json.Unmarshal(routeBytes, &routeMap); err != nil {
		t.Fatalf("erro no JSON do route: %v", err)
	}

	nodes, ok := routeMap["nodes"].([]any)
	if !ok || len(nodes) != 4 { // TICK-01.1 + Stop-Gate do TICK-01.1 + TICK-01.2 + SG-FINAL
		t.Fatalf("esperava 4 nós na rota (tickets, stop-gate injetado e SG-FINAL), obteve %d", len(nodes))
	}
}

func TestPlanGenerator_AutonomousMode(t *testing.T) {
	mockPlanOutput := `
### [TICKET-E01.1] Implementação direta
- **Épico:** E01
- **Natureza:** micro
- **Modo:** parallel
- **Dependências:** nenhuma
- **Arquivos Afetados:** main.go
- **Critérios Atendidos:** AC-E01.1
- **Comando de Validação:** go test ./...
`
	tempDir := t.TempDir()
	specFile := filepath.Join(tempDir, "SPEC.md")
	os.WriteFile(specFile, []byte("teste"), 0644)

	llm := &mockPlanLLM{output: mockPlanOutput}
	gen := plangen.NewGenerator(llm)

	res, err := gen.Generate(context.Background(), tempDir, specFile, plangen.SupervisionAutonomous)
	if err != nil {
		t.Fatalf("erro no modo autônomo: %v", err)
	}

	routeBytes, _ := os.ReadFile(res.RouteJSONPath)
	var routeMap map[string]any
	_ = json.Unmarshal(routeBytes, &routeMap)
	nodes := routeMap["nodes"].([]any)

	// No modo autônomo, só o SG-FINAL (homologação mecânica) entra: nenhum Stop-Gate humano intermediário
	if len(nodes) != 2 {
		t.Fatalf("modo autônomo deve ter só o ticket e o SG-FINAL, obteve %d nós", len(nodes))
	}
}

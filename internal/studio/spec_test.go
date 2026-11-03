package studio

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/interview"
	"legoloop/internal/specgen"
)

type mockStudioSpecLLM struct {
	output string
}

func (m *mockStudioSpecLLM) Chat(ctx context.Context, msgs []interview.Message) (string, error) {
	return m.output, nil
}

func TestStudioSpecAPI_Endpoint(t *testing.T) {
	s, _ := newServer(t)

	validSpecMD := `# 📋 SPEC: API de Pagamento

## 1. Premissas Técnicas Assumidas
Premissa 1

## 2. Fronteiras de Ação & Regras de Governança
Sempre TDD

## 3. Visão Geral, Arquitetura & Costuras de Teste
Costura em payment.go

## 4. Contratos de Dados & Schemas Formais
Schemas JSON

## 5. Critérios Binários de Aceite por Nó
### E01 — Pagamento
* [ ] AC-E01.1: Validação do payload
  * **Dado que:** Payload correto
  * **Quando:** Disparar request
  * **Então:** exit_code == 0

## 6. Comandos Operacionais & Critério de Conclusão Global
go test ./...

---
🛑 STOP-GATE 02 DE AUDITORIA HUMANA: Aprovado.
`

	mockLLM := &mockStudioSpecLLM{output: validSpecMD}
	s.SpecGen = specgen.NewGenerator(mockLLM, s.Index)
	h := s.Handler()

	intentFile := filepath.Join(s.ProjectDir, "2026-11-02_INTENT_TESTE.md")
	os.WriteFile(intentFile, []byte("# INTENT Teste"), 0644)

	reqBody := `{"intent_path":"` + filepath.ToSlash(intentFile) + `"}`
	rec := do(t, h, http.MethodPost, "/api/spec/generate", reqBody)
	if rec.Code != http.StatusOK {
		t.Fatalf("esperava 200 em spec/generate, obteve %d: %s", rec.Code, rec.Body.String())
	}

	var res specgen.SpecResult
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}

	if len(res.Criteria) != 1 || res.Criteria[0] != "AC-E01.1" {
		t.Errorf("critérios inesperados: %v", res.Criteria)
	}
}

package specgen_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/graph"
	"legoloop/internal/interview"
	"legoloop/internal/specgen"
)

type mockSpecLLM struct {
	output string
}

func (m *mockSpecLLM) Chat(ctx context.Context, msgs []interview.Message) (string, error) {
	return m.output, nil
}

func TestSpecGenerator_Success(t *testing.T) {
	validSpecMD := `# 📋 SPEC: Automação de Testes Go

> **Status:** Proposta para Revisão
> **Data:** 2026-11-02
> **Referência da Intenção:** docs/intents/2026-11-02_INTENT_TESTES.md

## 1. Premissas Técnicas Assumidas (Addy Osmani)
- 1. Go 1.23+ com testes nativos

## 2. Fronteiras de Ação & Regras de Governança (Boundaries — Addy Osmani)
- Sempre: TDD estrito
- Nunca: Zero stubs

## 3. Visão Geral, Arquitetura & Costuras de Teste (Test Seams — Pocock)
Costura primária em internal/service/engine.go

## 4. Contratos de Dados & Schemas Formais
Structs Go.

## 5. Critérios Binários de Aceite por Nó (BDD)
### E01 — Testes
* [ ] AC-E01.1: Validação de suíte Go
  * **Dado que:** Arquivos em internal/
  * **Quando:** Executar go test
  * **Então:** exit_code == 0
* [ ] AC-E01.2: Linter de conformidade
  * **Dado que:** Código escrito
  * **Quando:** Executar golangci-lint
  * **Então:** exit_code == 0

## 6. Comandos Operacionais & Critério de Conclusão Global (Addy Osmani)
- Test: go test ./...

---
🛑 STOP-GATE 02 DE AUDITORIA HUMANA: Aguardando aprovação explícita para gerar PLAN e ROUTE.
`

	tempDir := t.TempDir()
	intentsDir := filepath.Join(tempDir, "docs", "intents")
	if err := os.MkdirAll(intentsDir, 0755); err != nil {
		t.Fatal(err)
	}

	intentFile := filepath.Join(intentsDir, "2026-11-02_INTENT_TESTES.md")
	if err := os.WriteFile(intentFile, []byte("# INTENT de Testes"), 0644); err != nil {
		t.Fatal(err)
	}

	idx, err := graph.NewIndex(graph.DefaultCatalog())
	if err != nil {
		t.Fatal(err)
	}

	llm := &mockSpecLLM{output: validSpecMD}
	gen := specgen.NewGenerator(llm, idx)

	res, err := gen.Generate(context.Background(), tempDir, intentFile)
	if err != nil {
		t.Fatalf("Generate falhou: %v", err)
	}

	if len(res.Criteria) != 2 {
		t.Errorf("esperava 2 critérios, obteve %d", len(res.Criteria))
	}
	if res.Criteria[0] != "AC-E01.1" || res.Criteria[1] != "AC-E01.2" {
		t.Errorf("critérios inesperados: %v", res.Criteria)
	}

	if !strings.HasSuffix(res.FilePath, "2026-11-02_SPEC_TESTES.md") {
		t.Errorf("caminho de spec inesperado: %s", res.FilePath)
	}

	content, err := os.ReadFile(res.FilePath)
	if err != nil {
		t.Fatalf("erro ao ler spec gerada: %v", err)
	}
	if !strings.Contains(string(content), "STOP-GATE 02") {
		t.Error("spec gerada não contém Stop-Gate 02")
	}
}

func TestSpecGenerator_ValidationFailsOnMissingSection(t *testing.T) {
	invalidSpecMD := `# 📋 SPEC: Incompleta
Sem as seções obrigatórias.
`
	tempDir := t.TempDir()
	intentFile := filepath.Join(tempDir, "INTENT.md")
	os.WriteFile(intentFile, []byte("teste"), 0644)

	idx, _ := graph.NewIndex(graph.DefaultCatalog())
	llm := &mockSpecLLM{output: invalidSpecMD}
	gen := specgen.NewGenerator(llm, idx)

	_, err := gen.Generate(context.Background(), tempDir, intentFile)
	if err == nil {
		t.Fatal("esperava erro de validação por seções ausentes, obteve nil")
	}
	if !strings.Contains(err.Error(), "seção obrigatória ausente") {
		t.Errorf("mensagem de erro inesperada: %v", err)
	}
}

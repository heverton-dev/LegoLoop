package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/interview"
	"github.com/heverton-dev/LegoLoop/internal/specgen"
)

type mockSpecLLM struct {
	output     string
	err        error
	calls      int
	failOnCall bool
	t          *testing.T
}

func (m *mockSpecLLM) Chat(_ context.Context, _ []interview.Message) (string, error) {
	m.calls++
	if m.failOnCall && m.t != nil {
		m.t.Errorf("LLM chamado indevidamente para INTENT em status UNDER_REVIEW")
	}
	if m.err != nil {
		return "", m.err
	}
	return m.output, nil
}

const specValidaMock = `# 📋 SPEC: Sistema de Teste

> **Status:** Proposta para Revisão
> **Data:** 2026-11-21
> **Referência da Intenção:** docs/intents/2026-11-21_INTENT_TESTE.md

## 1. Premissas Técnicas Assumidas (Addy Osmani)
- 1. Go 1.23+ com testes nativos

## 2. Fronteiras de Ação & Regras de Governança (Boundaries — Addy Osmani)
- Sempre: TDD estrito
- Nunca: Zero stubs

## 3. Visão Geral, Arquitetura & Costuras de Teste (Test Seams — Pocock)
Costura primária em internal/teste.go

## 4. Contratos de Dados & Schemas Formais
Structs Go.

## 5. Critérios Binários de Aceite por Nó (BDD)
### E01 — Testes
* [ ] AC-E01.1: Validação de suíte Go
  * **Dado que:** Arquivos em internal/
  * **Quando:** Executar go test
  * **Então:** exit_code == 0

## 6. Comandos Operacionais & Critério de Conclusão Global (Addy Osmani)
- Test: go test ./...

---
🛑 STOP-GATE 02 DE AUDITORIA HUMANA: Aguardando aprovação explícita para gerar PLAN e ROUTE.
`

func prepararAmbienteSpec(t *testing.T, status string) (raiz, intentPath string) {
	t.Helper()
	raiz = t.TempDir()
	intentsDir := filepath.Join(raiz, "docs", "intents")
	if err := os.MkdirAll(intentsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	intentPath = filepath.Join(intentsDir, "2026-11-21_INTENT_TESTE.md")
	content := "---\nid: INTENT-TESTE\nstatus: " + status + "\n---\n# INTENT de Teste\n"
	if err := os.WriteFile(intentPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return raiz, intentPath
}

func TestFabricarSpec(t *testing.T) {
	raiz, intentPath := prepararAmbienteSpec(t, "APPROVED")
	mockLLM := &mockSpecLLM{output: specValidaMock}

	var out, errb bytes.Buffer
	d := Deps{
		Stdout: &out,
		Stderr: &errb,
		Fabricar: &FabricarDeps{
			LLM: mockLLM,
		},
	}

	code := cmdFabricar(context.Background(), []string{"spec", "--raiz", raiz, "--intent", intentPath}, d)
	if code != 0 {
		t.Fatalf("esperava exit code 0, obteve %d. Stderr: %s", code, errb.String())
	}

	saida := out.String()
	if !strings.Contains(saida, "SPEC gravada em") {
		t.Errorf("saida esperada nao encontrada: %s", saida)
	}

	matches, err := filepath.Glob(filepath.Join(raiz, "docs", "specs", "*_SPEC_TESTE.md"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 {
		t.Fatalf("esperava 1 arquivo SPEC gerado, encontrou %d", len(matches))
	}

	specConteudo, err := os.ReadFile(matches[0])
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(specConteudo), "STOP-GATE 02") {
		t.Errorf("arquivo de SPEC gerado nao contem STOP-GATE 02")
	}
}

func TestFabricarSpec_JSON(t *testing.T) {
	raiz, intentPath := prepararAmbienteSpec(t, "APPROVED")
	mockLLM := &mockSpecLLM{output: specValidaMock}

	var out, errb bytes.Buffer
	d := Deps{
		Stdout: &out,
		Stderr: &errb,
		Fabricar: &FabricarDeps{
			LLM: mockLLM,
		},
	}

	code := cmdFabricar(context.Background(), []string{"spec", "--raiz", raiz, "--intent", intentPath, "--json"}, d)
	if code != 0 {
		t.Fatalf("esperava exit code 0 com --json, obteve %d. Stderr: %s", code, errb.String())
	}

	var res specgen.SpecResult
	if err := json.Unmarshal(out.Bytes(), &res); err != nil {
		t.Fatalf("saida nao e JSON valido: %v", err)
	}
	if !strings.Contains(res.Theme, "SPEC_TESTE.md") {
		t.Errorf("tema inesperado no JSON: %s", res.Theme)
	}
	if !strings.HasSuffix(res.FilePath, res.Theme) {
		t.Errorf("filepath inesperado no JSON: %s", res.FilePath)
	}
}

func TestFabricarSpec_IntentObrigatoria(t *testing.T) {
	var out, errb bytes.Buffer
	d := Deps{Stdout: &out, Stderr: &errb}

	code := cmdFabricar(context.Background(), []string{"spec"}, d)
	if code != 2 {
		t.Fatalf("esperava exit code 2 sem flag --intent, obteve %d", code)
	}
	if !strings.Contains(errb.String(), "uso:") {
		t.Errorf("esperava mensagem de uso no stderr: %s", errb.String())
	}
}

func TestFabricarSpecRecusa(t *testing.T) {
	raiz, intentPath := prepararAmbienteSpec(t, "UNDER_REVIEW")
	mockLLM := &mockSpecLLM{
		output:     specValidaMock,
		failOnCall: true,
		t:          t,
	}

	var out, errb bytes.Buffer
	d := Deps{
		Stdout: &out,
		Stderr: &errb,
		Fabricar: &FabricarDeps{
			LLM: mockLLM,
		},
	}

	code := cmdFabricar(context.Background(), []string{"spec", "--raiz", raiz, "--intent", intentPath}, d)
	if code != 2 {
		t.Fatalf("esperava exit code 2 para INTENT nao aprovada, obteve %d", code)
	}
	if mockLLM.calls != 0 {
		t.Errorf("LLM nao deveria ser chamado para INTENT nao aprovada, chamadas: %d", mockLLM.calls)
	}
}

func TestFabricarSpec_ErroLLM(t *testing.T) {
	raiz, intentPath := prepararAmbienteSpec(t, "APPROVED")
	mockLLM := &mockSpecLLM{err: errors.New("falha de conexao com LLM")}

	var out, errb bytes.Buffer
	d := Deps{
		Stdout: &out,
		Stderr: &errb,
		Fabricar: &FabricarDeps{
			LLM: mockLLM,
		},
	}

	code := cmdFabricar(context.Background(), []string{"spec", "--raiz", raiz, "--intent", intentPath}, d)
	if code != 1 {
		t.Fatalf("esperava exit code 1 em caso de erro do LLM, obteve %d", code)
	}
}

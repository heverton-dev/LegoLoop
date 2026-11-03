package specgen_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/docs/templates"
	"legoloop/internal/graph"
	"legoloop/internal/interview"
	"legoloop/internal/specgen"
)

type capturingLLM struct {
	output string
	got    []interview.Message
}

func (c *capturingLLM) Chat(_ context.Context, msgs []interview.Message) (string, error) {
	c.got = msgs
	return c.output, nil
}

const bodyOK = `# 📋 SPEC: Automação de Testes Go

> **Status:** texto que o motor deve substituir
> **Referência da Intenção:** C:\\qualquer\\caminho\\absoluto.md

## 1. Premissas Técnicas Assumidas (Addy Osmani)
1. Go conforme o go.mod.

## 2. Fronteiras de Ação & Regras de Governança (Boundaries — Addy Osmani)
* Sempre: TDD.

## 3. Visão Geral, Arquitetura & Costuras de Teste (Test Seams — Pocock)
Costura primária: Run.

## 4. Contratos de Dados & Schemas Formais
Structs.

## 5. Critérios Binários de Aceite por Nó / Ticket (BDD)
* [ ] AC-E01.1: Caminho feliz
  * **Dado que:** x
  * **Quando:** y
  * **Então:** exit_code == 0

## 6. Comandos Operacionais & Critério de Conclusão Global (Addy Osmani)
* Test: go test ./...

🛑 **STOP-GATE 02 DE AUDITORIA HUMANA:**
`

func setup(t *testing.T, intentBody string) (string, string) {
	t.Helper()
	base := t.TempDir()
	dir := filepath.Join(base, "docs", "intents")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	intent := filepath.Join(dir, "2026-11-02_INTENT_TESTES.md")
	if err := os.WriteFile(intent, []byte(intentBody), 0o644); err != nil {
		t.Fatal(err)
	}
	return base, intent
}

func generate(t *testing.T, llm interview.LLMClient, intentBody string) (*specgen.SpecResult, error) {
	t.Helper()
	base, intent := setup(t, intentBody)
	idx, err := graph.NewIndex(graph.DefaultCatalog())
	if err != nil {
		t.Fatal(err)
	}
	return specgen.NewGenerator(llm, idx).Generate(context.Background(), base, intent)
}

func TestSpecGenerator_SystemPromptEmbedsCanonicalTemplate(t *testing.T) {
	llm := &capturingLLM{output: bodyOK}
	if _, err := generate(t, llm, "# INTENT"); err != nil {
		t.Fatal(err)
	}
	sys := llm.got[0].Content
	for _, h := range templates.Headings(templates.Spec()) {
		if !strings.Contains(sys, "## "+h) {
			t.Errorf("o prompt do sistema não embute a seção do template: %q", h)
		}
	}
	if !strings.Contains(sys, "TEMPLATE_SPEC.md v"+templates.Version) {
		t.Error("o prompt deve identificar a versão do template")
	}
	if strings.Contains(sys, "GUIA DE PREENCHIMENTO") {
		t.Error("o comentário de orientação ao autor não deve ir para o prompt")
	}
}

func TestSpecGenerator_StampsTemplateHeaderDeterministically(t *testing.T) {
	res, err := generate(t, &capturingLLM{output: bodyOK}, "# INTENT\n\n> **Alvo:** internal/service\n")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(res.FilePath)
	doc := string(b)
	for _, want := range []string{
		"# 📋 SPEC: Automação de Testes Go",
		"> **Status:** Proposta para Revisão",
		"> **Versão:** 1.0.0",
		"[docs/intents/2026-11-02_INTENT_TESTES.md](../intents/2026-11-02_INTENT_TESTES.md)",
		"> **Alvo Primário:** internal/service",
		"> **Consome:** INTENT aprovada · **Produz:** PLAN e ROUTE",
		"> **Template:** TEMPLATE_SPEC.md v" + templates.Version,
	} {
		if !strings.Contains(doc, want) {
			t.Errorf("cabeçalho sem %q\n%s", want, doc)
		}
	}
	if strings.Contains(doc, "texto que o motor deve substituir") || strings.Contains(doc, `C:\qualquer`) {
		t.Error("o cabeçalho escrito pelo modelo deve ser substituído pelo do motor")
	}
	if !strings.Contains(doc, "## 1. Premissas Técnicas Assumidas") {
		t.Error("o corpo do modelo deve ser preservado")
	}
}

func TestSpecGenerator_UnknownTargetIsDeclaredNotInvented(t *testing.T) {
	res, err := generate(t, &capturingLLM{output: bodyOK}, "# INTENT sem alvo")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(res.FilePath)
	if !strings.Contains(string(b), "> **Alvo Primário:** A definir") {
		t.Errorf("sem Alvo na INTENT, o motor deve declarar 'A definir'\n%s", b)
	}
}

func TestSpecGenerator_RejectsUnresolvedPlaceholders(t *testing.T) {
	_, err := generate(t, &capturingLLM{output: strings.Replace(bodyOK, "Structs.", "{{CONTRATO}}", 1)}, "# INTENT")
	if err == nil || !strings.Contains(err.Error(), "placeholder") {
		t.Fatalf("esperava erro de placeholder sem preencher, obteve %v", err)
	}
}

func TestSpecGenerator_RequiresTitle(t *testing.T) {
	noTitle := strings.Replace(bodyOK, "# 📋 SPEC: Automação de Testes Go", "Sem título", 1)
	if _, err := generate(t, &capturingLLM{output: noTitle}, "# INTENT"); err == nil {
		t.Fatal("esperava erro: a SPEC precisa do título '# 📋 SPEC:'")
	}
}

func TestSpecGenerator_AllowsGoTemplateSyntaxInBody(t *testing.T) {
	body := strings.Replace(bodyOK, "Structs.", "Exemplo de template Go: {{.Nome}} e {{- range .Itens}}.", 1)
	if _, err := generate(t, &capturingLLM{output: body}, "# INTENT"); err != nil {
		t.Fatalf("sintaxe de template Go no corpo não é placeholder: %v", err)
	}
}

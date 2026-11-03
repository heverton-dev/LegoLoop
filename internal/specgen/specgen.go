package specgen

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"legoloop/internal/graph"
	"legoloop/internal/interview"
)

// Generator orquestra a criação determinística de especificações técnicas formais.
type Generator struct {
	llm   interview.LLMClient
	index *graph.Index
}

// NewGenerator inicializa o gerador híbrido de SPEC.
func NewGenerator(llm interview.LLMClient, idx *graph.Index) *Generator {
	return &Generator{
		llm:   llm,
		index: idx,
	}
}

// SpecResult contém os dados estruturados extraídos da especificação gerada.
type SpecResult struct {
	Theme       string   `json:"theme"`
	FilePath    string   `json:"file_path"`
	Criteria    []string `json:"criteria"`
	RawMarkdown string   `json:"raw_markdown"`
}

// PromptCanônico da SPEC: Fusão Addy Osmani (Premissas/Boundaries) + Pocock (Test Seams) + AIDD-Spec (BDD)
const PromptSpecSystem = `Você é o Arquiteto de Especificações Formais do LegoLoop (fusão Addy Osmani + Matt Pocock + AIDD-Spec).
Sua missão é transformar um documento de INTENÇÃO e o GRAFO de código existente em uma ESPECIFICAÇÃO TÉCNICA DE ENGENHARIA AUDITÁVEL.

ESTRUTURA OBRIGATÓRIA DA SPEC (NÃO altere os títulos dos cabeçalhos):
# 📋 SPEC: [Tema da Especificação]

> **Status:** Proposta para Revisão
> **Data:** %s
> **Referência da Intenção:** %s

## 1. Premissas Técnicas Assumidas (Addy Osmani)
[Liste explicitamente 3 a 5 premissas técnicas adotadas para não tomar decisões silenciosas]

## 2. Fronteiras de Ação & Regras de Governança (Boundaries — Addy Osmani)
* **Sempre Fazer:** [TDD estrito, funções pequenas < 35 linhas, validação rigorosa]
* **Perguntar Antes:** [Novas dependências, alterações de schemas persistidos]
* **Nunca Fazer:** [Zero stubs, zero mocks vazios, nunca commitar credenciais]

## 3. Visão Geral, Arquitetura & Costuras de Teste (Test Seams — Pocock)
[Descrição técnica com diagrama Mermaid e a identificação da costura primária mais alta para os testes baterem]

## 4. Contratos de Dados & Schemas Formais
[Tipagem estrita de structs Go e contratos JSON/REST]

## 5. Critérios Binários de Aceite por Nó (BDD)
[Para cada épico, liste os critérios AC-E{N}.{M} estritamente no formato Dado/Quando/Então]
Exemplo:
### E01 — Fundação
* [ ] AC-E01.1: [Título do Critério]
  * **Dado que:** [Estado inicial / fixtures]
  * **Quando:** [Ação disparada]
  * **Então:** [Comando com exit_code == 0 ou asserção binária]

## 6. Comandos Operacionais & Critério de Conclusão Global (Addy Osmani)
* **Build:** go build -o bin/legoloop.exe ./cmd/legoloop
* **Test:** go test -v ./...
* **Lint:** golangci-lint run ./...

---
🛑 STOP-GATE 02 DE AUDITORIA HUMANA: Aguardando aprovação explícita para gerar PLAN e ROUTE.
`

// Generate consome o caminho do INTENT, consulta o grafo e gera a SPEC.
func (g *Generator) Generate(ctx context.Context, baseDir, intentPath string) (*SpecResult, error) {
	intentContent, err := os.ReadFile(intentPath)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler artefato de intenção: %w", err)
	}

	today := time.Now().Format("2006-01-02")
	systemPrompt := fmt.Sprintf(PromptSpecSystem, today, filepath.ToSlash(intentPath))

	userPrompt := fmt.Sprintf(`Gere a SPEC formal para a seguinte INTENÇÃO:

%s

Atenção: Garanta que todas as structs propostas sejam compatíveis com a arquitetura Go e os critérios AC-XX usem asserções binárias verificáveis mecanicamente.`, string(intentContent))

	messages := []interview.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	rawMarkdown, err := g.llm.Chat(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("falha ao gerar spec via LLM: %w", err)
	}

	// Validação de formato (Linter da Spec)
	criteria, err := validateSpecMarkdown(rawMarkdown)
	if err != nil {
		return nil, fmt.Errorf("especificação gerada violou o contrato: %w", err)
	}

	// Persistência em docs/specs/
	specsDir := filepath.Join(baseDir, "docs", "specs")
	if err := os.MkdirAll(specsDir, 0755); err != nil {
		specsDir = baseDir
	}

	baseName := filepath.Base(intentPath)
	specFileName := strings.Replace(baseName, "INTENT", "SPEC", 1)
	if specFileName == baseName {
		specFileName = fmt.Sprintf("%s_SPEC_%s.md", today, "GERADA")
	}

	specFilePath := filepath.Join(specsDir, specFileName)
	if err := os.WriteFile(specFilePath, []byte(rawMarkdown), 0644); err != nil {
		return nil, fmt.Errorf("falha ao gravar spec no disco: %w", err)
	}

	return &SpecResult{
		Theme:       specFileName,
		FilePath:    specFilePath,
		Criteria:    criteria,
		RawMarkdown: rawMarkdown,
	}, nil
}

// validateSpecMarkdown assegura a integridade das 6 seções obrigatórias e extrai critérios AC-XX.
func validateSpecMarkdown(content string) ([]string, error) {
	requiredHeadings := []string{
		"1. Premissas Técnicas Assumidas",
		"2. Fronteiras de Ação & Regras de Governança",
		"3. Visão Geral, Arquitetura & Costuras de Teste",
		"4. Contratos de Dados & Schemas Formais",
		"5. Critérios Binários de Aceite por Nó",
		"6. Comandos Operacionais & Critério de Conclusão Global",
		"STOP-GATE 02",
	}

	for _, h := range requiredHeadings {
		if !strings.Contains(content, h) {
			return nil, fmt.Errorf("seção obrigatória ausente na SPEC: %q", h)
		}
	}

	reAC := regexp.MustCompile(`(?i)\*\s*\[\s*\]\s*(AC-[A-Za-z0-9_.-]+)`)
	matches := reAC.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return nil, errors.New("a SPEC não contém nenhum critério binário de aceite 'AC-XX'")
	}

	var criteria []string
	for _, m := range matches {
		if len(m) > 1 {
			criteria = append(criteria, m[1])
		}
	}

	return criteria, nil
}

package plangen

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"legoloop/internal/interview"
)

// SupervisionLevel define o nível de intervenção humana (Human-in-the-Loop).
type SupervisionLevel string

const (
	SupervisionNano       SupervisionLevel = "nano"       // Para a cada nano, micro e épico
	SupervisionMicro      SupervisionLevel = "micro"      // Para a cada micro e épico
	SupervisionEpic       SupervisionLevel = "epic"       // Para ao final de cada épico (padrão)
	SupervisionAutonomous SupervisionLevel = "autonomous" // Só para na homologação final
)

// Ticket representa uma unidade de trabalho atômica e verticalmente fatiada (Addy Osmani).
type Ticket struct {
	ID             string   `json:"id"`              // TICK-01.1
	Title          string   `json:"title"`           // Descrição imperativa
	EpicID         string   `json:"epic_id"`         // E01
	Nature         string   `json:"nature"`          // nano | micro
	ExecutionMode  string   `json:"execution_mode"`  // tdd_red_green | sequential | parallel
	DependsOn      []string `json:"depends_on"`      // IDs predecessores
	TargetFiles    []string `json:"target_files"`    // Arquivos tocados
	SpecCriteria   []string `json:"spec_criteria"`   // AC-E01.1
	ValidationCmd  string   `json:"validation_cmd"`  // Comando mecânico de teste
	HumanGateAfter bool     `json:"human_gate_after"`// Se há parada após este ticket
}

// PlanResult consolida os artefatos de planejamento gerados.
type PlanResult struct {
	Theme            string           `json:"theme"`
	SupervisionLevel SupervisionLevel `json:"supervision_level"`
	PlanMarkdownPath string           `json:"plan_markdown_path"`
	RouteJSONPath    string           `json:"route_json_path"`
	Tickets          []Ticket         `json:"tickets"`
	RawPlan          string           `json:"raw_plan"`
}

// Generator orquestra a decomposição socrática em tarefas fatiadas verticalmente.
type Generator struct {
	llm interview.LLMClient
}

// NewGenerator inicializa o gerador de planos.
func NewGenerator(llm interview.LLMClient) *Generator {
	return &Generator{llm: llm}
}

// PromptPlanSystem adaptado de Addy Osmani (planning-and-task-breakdown) para PT-BR
const PromptPlanSystem = `Você é o Arquiteto de Planejamento e Decomposição de Tarefas do LegoLoop (método Addy Osmani).
Sua missão é transformar a ESPECIFICAÇÃO TÉCNICA em um PLANO DE AÇÃO ATÔMICO com FATIAMENTO VERTICAL (Vertical Slicing) estrito.

REGRAS INVIOLÁVEIS DO ADDY OSMANI:
1. FATIAMENTO VERTICAL: Cada ticket deve entregar uma funcionalidade ponta a ponta e testável (contrato + teste falhando + implementação + validação). NUNCA crie tickets de "só criar arquivos vazios", "só tipos" ou "só testes".
2. GRAFO DE DEPENDÊNCIAS: Ordene os tickets bottom-up (fundações e modelos primeiro, endpoints/lógica depois, interfaces por último).
3. COBERTURA TOTAL: Todo critério AC-XX da SPEC DEVE ser coberto por pelo menos um ticket.
4. VERIFICAÇÃO TRIPLA: Todo ticket deve ter um comando de teste/linter automatizado verificável (exit code == 0).

FORMATO ESTRUTURADO DE CADA TICKET:
### [TICKET-E{N}.{M}] {Título Imperativo}
- **Épico:** E{N}
- **Natureza:** micro | nano
- **Modo:** tdd_red_green | sequential | parallel
- **Dependências:** [nenhuma | TICK-XX.X]
- **Arquivos Afetados:** {caminhos dos arquivos separados por vírgula}
- **Critérios Atendidos:** {AC-E{N}.{M}}
- **Comando de Validação:** {ex: go test -v ./internal/pkg/...}
`

// Generate cria o PLAN.md e o ROUTE.json em sincronia com o nível de supervisão.
func (g *Generator) Generate(ctx context.Context, baseDir, specPath string, level SupervisionLevel) (*PlanResult, error) {
	if level == "" {
		level = SupervisionEpic
	}

	specBytes, err := os.ReadFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler spec: %w", err)
	}

	userPrompt := fmt.Sprintf(`Gere o Plano de Execução do Addy Osmani para a seguinte SPEC:

%s

Decomponha todos os critérios AC-XX em tickets verticais sequenciados logicamente por dependências.`, string(specBytes))

	messages := []interview.Message{
		{Role: "system", Content: PromptPlanSystem},
		{Role: "user", Content: userPrompt},
	}

	rawMarkdown, err := g.llm.Chat(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("falha ao invocar LLM para planejamento: %w", err)
	}

	tickets, err := parseTickets(rawMarkdown, level)
	if err != nil {
		return nil, fmt.Errorf("falha ao interpretar tickets gerados: %w", err)
	}
	if len(tickets) == 0 {
		return nil, errors.New("nenhum ticket foi identificado no plano gerado")
	}

	baseName := filepath.Base(specPath)
	planFileName := strings.Replace(baseName, "SPEC", "PLAN", 1)
	routeFileName := strings.Replace(baseName, "SPEC", "ROUTE", 1)
	routeFileName = strings.Replace(routeFileName, ".md", ".json", 1)

	plansDir := filepath.Join(baseDir, "docs", "plans")
	routesDir := filepath.Join(baseDir, "docs", "routes")
	_ = os.MkdirAll(plansDir, 0755)
	_ = os.MkdirAll(routesDir, 0755)

	planPath := filepath.Join(plansDir, planFileName)
	routePath := filepath.Join(routesDir, routeFileName)

	// Materializar PLAN.md
	planMD := buildPlanMarkdown(planFileName, specPath, routePath, level, tickets, rawMarkdown)
	if err := os.WriteFile(planPath, []byte(planMD), 0644); err != nil {
		return nil, fmt.Errorf("falha ao gravar plan.md: %w", err)
	}

	// Materializar ROUTE.json
	routeJSON, err := buildRouteJSON(routeFileName, specPath, planPath, level, tickets)
	if err != nil {
		return nil, fmt.Errorf("falha ao construir route.json: %w", err)
	}
	if err := os.WriteFile(routePath, routeJSON, 0644); err != nil {
		return nil, fmt.Errorf("falha ao gravar route.json: %w", err)
	}

	return &PlanResult{
		Theme:            planFileName,
		SupervisionLevel: level,
		PlanMarkdownPath: planPath,
		RouteJSONPath:    routePath,
		Tickets:          tickets,
		RawPlan:          rawMarkdown,
	}, nil
}

func parseTickets(content string, level SupervisionLevel) ([]Ticket, error) {
	lines := strings.Split(content, "\n")
	type rawBlock struct {
		id    string
		title string
		lines []string
	}

	var blocks []rawBlock
	reHeader := regexp.MustCompile(`(?i)^###\s*\[(TICKET-[A-Za-z0-9_.-]+)\]\s*(.*)`)

	var cur *rawBlock
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if m := reHeader.FindStringSubmatch(trimmed); len(m) > 2 {
			if cur != nil {
				blocks = append(blocks, *cur)
			}
			cur = &rawBlock{
				id:    strings.TrimSpace(m[1]),
				title: strings.TrimSpace(m[2]),
			}
			continue
		}
		if cur != nil {
			cur.lines = append(cur.lines, l)
		}
	}
	if cur != nil {
		blocks = append(blocks, *cur)
	}

	var tickets []Ticket
	reField := func(blockLines []string, name string) string {
		re := regexp.MustCompile(`(?i)^\s*-\s*\*\*` + name + `:\*\*\s*(.*)`)
		for _, bl := range blockLines {
			if m := re.FindStringSubmatch(bl); len(m) > 1 {
				return strings.TrimSpace(m[1])
			}
		}
		return ""
	}

	for _, b := range blocks {
		epicID := reField(b.lines, "Épico")
		if epicID == "" {
			epicID = "E01"
		}
		nature := strings.ToLower(reField(b.lines, "Natureza"))
		if nature != "nano" {
			nature = "micro"
		}
		mode := strings.ToLower(reField(b.lines, "Modo"))
		if mode == "" {
			mode = "tdd_red_green"
		}

		depRaw := reField(b.lines, "Dependências")
		var deps []string
		if depRaw != "" && !strings.Contains(strings.ToLower(depRaw), "nenhuma") {
			for _, d := range strings.Split(depRaw, ",") {
				dClean := strings.TrimSpace(d)
				if dClean != "" {
					deps = append(deps, dClean)
				}
			}
		}

		filesRaw := reField(b.lines, "Arquivos Afetados")
		var targetFiles []string
		for _, f := range strings.Split(filesRaw, ",") {
			fClean := strings.Trim(strings.TrimSpace(f), "`")
			if fClean != "" {
				targetFiles = append(targetFiles, fClean)
			}
		}

		critRaw := reField(b.lines, "Critérios Atendidos")
		var criteria []string
		for _, c := range strings.Split(critRaw, ",") {
			cClean := strings.Trim(strings.TrimSpace(c), "`")
			if cClean != "" {
				criteria = append(criteria, cClean)
			}
		}

		valCmd := reField(b.lines, "Comando de Validação")
		if valCmd == "" {
			valCmd = "go test ./..."
		}

		gateAfter := false
		switch level {
		case SupervisionNano:
			gateAfter = true
		case SupervisionMicro:
			if nature == "micro" {
				gateAfter = true
			}
		case SupervisionEpic, SupervisionAutonomous:
			gateAfter = false
		}

		tickets = append(tickets, Ticket{
			ID:             b.id,
			Title:          b.title,
			EpicID:         epicID,
			Nature:         nature,
			ExecutionMode:  mode,
			DependsOn:      deps,
			TargetFiles:    targetFiles,
			SpecCriteria:   criteria,
			ValidationCmd:  valCmd,
			HumanGateAfter: gateAfter,
		})
	}

	return tickets, nil
}

func buildPlanMarkdown(fileName, specPath, routePath string, level SupervisionLevel, tickets []Ticket, rawPlan string) string {
	today := time.Now().Format("2006-01-02")
	return fmt.Sprintf(`# 🗺️ PLAN: %s

> **Status:** Proposta para Revisão Humana  
> **Data:** %s  
> **Referência da Spec:** %s  
> **Espelho Executável JSON:** %s  
> **Modo de Supervisão Humana:** %s  

---

## 1. Estratégia de Execução & Human-in-the-Loop

* **Modo de Supervisão:** %s
* **Isolamento Padrão:** Workspace: local | Executor: orchestrix_hermetic (Zero Context Bloat)
* **TDD Estrito:** Obrigatório (Red ➔ Green ➔ Refactor)
* **Limite de Workers:** 2

---

## 2. Decomposição Estruturada de Tickets (Addy Osmani)

%s

---

🛑 **STOP-GATE 03 DE AUDITORIA HUMANA:**  
*Plano de execução fatiado com %d tickets. Deseja aprovar o plano e liberar a execução autônoma do BUILD?*
`, fileName, today, filepath.ToSlash(specPath), filepath.ToSlash(routePath), string(level), string(level), rawPlan, len(tickets))
}

func buildRouteJSON(name, specPath, planPath string, level SupervisionLevel, tickets []Ticket) ([]byte, error) {
	type RouteNode struct {
		ID            string   `json:"id"`
		Label         string   `json:"label"`
		Nature        string   `json:"nature"`
		Type          string   `json:"type"`
		ExecutionMode string   `json:"execution_mode"`
		DependsOn     []string `json:"depends_on"`
		Isolation     map[string]any `json:"isolation"`
		Routing       map[string]any `json:"routing"`
		Inputs        map[string]any `json:"inputs"`
		Output        map[string]any `json:"output"`
	}

	var nodes []RouteNode
	for _, t := range tickets {
		effort := "medium"
		model := "code-fast"
		if t.Nature == "nano" {
			effort = "low"
			model = "code-nano"
		}

		nodes = append(nodes, RouteNode{
			ID:            t.ID,
			Label:         t.Title,
			Nature:        t.Nature,
			Type:          "unit_task",
			ExecutionMode: t.ExecutionMode,
			DependsOn:     t.DependsOn,
			Isolation: map[string]any{
				"workspace":         "local",
				"executor":          "orchestrix_hermetic",
				"subagents_allowed": false,
			},
			Routing: map[string]any{
				"mode":    "suggested",
				"harness": "9router",
				"model":   model,
				"effort":  effort,
				"fallback": map[string]string{
					"harness": "antigravity",
					"model":   "flash",
				},
			},
			Inputs: map[string]any{
				"source_files": t.TargetFiles,
				"spec_ref":     t.SpecCriteria,
			},
			Output: map[string]any{
				"artifacts": t.TargetFiles,
				"gate":      "exit_code == 0 && " + t.ValidationCmd,
			},
		})

		// Se há stop-gate humano após este ticket
		if t.HumanGateAfter {
			nodes = append(nodes, RouteNode{
				ID:            "SG-" + t.ID,
				Label:         "Stop-Gate Humano pós " + t.ID,
				Nature:        "nano",
				Type:          "stop_gate",
				ExecutionMode: "sequential",
				DependsOn:     []string{t.ID},
				Isolation: map[string]any{
					"workspace":         "local",
					"executor":          "local_process",
					"subagents_allowed": false,
				},
				Routing: map[string]any{"mode": "manual", "model": "human/operator"},
				Inputs:  map[string]any{},
				Output: map[string]any{
					"gate": "human_approved == true",
				},
			})
		}
	}

	routeObj := map[string]any{
		"name":               name,
		"version":            "1.0.0",
		"date":               time.Now().Format("2006-01-02"),
		"spec_ref":           filepath.ToSlash(specPath),
		"plan_ref":           filepath.ToSlash(planPath),
		"execution_strategy": map[string]any{
			"supervision_level":    string(level),
			"max_parallel_workers": 2,
			"tdd_enforced":         true,
		},
		"nodes": nodes,
	}

	return json.MarshalIndent(routeObj, "", "  ")
}

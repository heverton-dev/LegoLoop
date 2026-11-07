package plangen

import (
	"context"
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
	ID             string   `json:"id"`                      // TICK-01.1
	Title          string   `json:"title"`                   // Descrição imperativa
	EpicID         string   `json:"epic_id"`                 // E01
	InitiativeID   string   `json:"initiative_id,omitempty"` // I01 (opcional)
	Nature         string   `json:"nature"`                  // nano | micro
	ExecutionMode  string   `json:"execution_mode"`          // tdd_red_green | sequential | parallel
	DependsOn      []string `json:"depends_on"`              // IDs predecessores
	TargetFiles    []string `json:"target_files"`            // Arquivos tocados
	SpecCriteria   []string `json:"spec_criteria"`           // AC-E01.1
	ValidationCmd  string   `json:"validation_cmd"`          // Comando mecânico de teste
	Gate           string   `json:"gate,omitempty"`          // Gate DSL do ticket; vazio = exit code 0 e artefatos existentes
	HumanGateAfter bool     `json:"human_gate_after"`        // Se há parada após este ticket
}

// PlanResult consolida os artefatos de planejamento gerados.
type PlanResult struct {
	Theme            string           `json:"theme"`
	SupervisionLevel SupervisionLevel `json:"supervision_level"`
	PlanMarkdownPath string           `json:"plan_markdown_path"`
	RouteJSONPath    string           `json:"route_json_path"`
	Tickets          []Ticket         `json:"tickets"`
	RawPlan          string           `json:"raw_plan"`
	Warnings         []string         `json:"warnings,omitempty"`
}

// Generator orquestra a decomposição socrática em tarefas fatiadas verticalmente.
type Generator struct {
	llm interview.LLMClient
	// Now fornece a data dos documentos; nil usa time.Now. Injetável para geração determinística.
	Now func() time.Time
}

func (g *Generator) now() time.Time {
	if g.Now != nil {
		return g.Now()
	}
	return time.Now()
}

// NewGenerator inicializa o gerador de planos.
func NewGenerator(llm interview.LLMClient) *Generator {
	return &Generator{llm: llm}
}

// PromptPlanSystem adaptado de Addy Osmani (planning-and-task-breakdown) para PT-BR.
// O PLAN.md e o ROUTE.json são montados pelo motor a partir dos tickets, seguindo TEMPLATE_PLAN.md e
// TEMPLATE_ROUTE.json (docs/templates); o modelo só decompõe, no formato abaixo.
const PromptPlanSystem = `Você é o Arquiteto de Planejamento e Decomposição de Tarefas do LegoLoop (método Addy Osmani).
Sua missão é transformar a ESPECIFICAÇÃO TÉCNICA em um PLANO DE AÇÃO ATÔMICO com FATIAMENTO VERTICAL (Vertical Slicing) estrito.

REGRAS INVIOLÁVEIS DO ADDY OSMANI:
1. FATIAMENTO VERTICAL: Cada ticket deve entregar uma funcionalidade ponta a ponta e testável (contrato + teste falhando + implementação + validação). NUNCA crie tickets de "só criar arquivos vazios", "só tipos" ou "só testes".
2. GRAFO DE DEPENDÊNCIAS: Ordene os tickets bottom-up (fundações e modelos primeiro, endpoints/lógica depois, interfaces por último). Dependências só podem citar tickets que você criou.
3. COBERTURA TOTAL: Todo critério AC-XX da SPEC DEVE ser coberto por pelo menos um ticket.
4. VERIFICAÇÃO TRIPLA: Todo ticket deve ter um comando de teste/linter automatizado verificável (exit code == 0).
5. COMANDOS SEM SHELL: o Comando de Validação é executado como argv, sem shell. Use um comando simples ou vários encadeados com "&&". NÃO use pipe, redirecionamento, aspas, glob nem variáveis. Use somente comandos que existem no repositório (go test, go vet, go build, gofmt).
6. IDs no formato TICK-{NN}.{M}, com NN o número do épico com dois dígitos (ex.: TICK-01.1, TICK-02.3).

FORMATO ESTRUTURADO DE CADA TICKET:
### [TICK-{NN}.{M}] {Título Imperativo}
- **Épico:** E{NN}
- **Natureza:** micro | nano
- **Modo:** tdd_red_green | sequential | parallel
- **Dependências:** [nenhuma | TICK-XX.X]
- **Arquivos Afetados:** {caminhos relativos ao projeto, separados por vírgula}
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
	raw, err := g.llm.Chat(ctx, planMessages(string(specBytes)))
	if err != nil {
		return nil, fmt.Errorf("falha ao invocar LLM para planejamento: %w", err)
	}
	tickets, err := parseTickets(raw, level)
	if err != nil {
		return nil, fmt.Errorf("falha ao interpretar tickets gerados: %w", err)
	}
	if len(tickets) == 0 {
		return nil, errors.New("nenhum ticket foi identificado no plano gerado")
	}
	res, err := materialize(baseDir, specPath, string(specBytes), level, tickets, g.now())
	if err != nil {
		return nil, err
	}
	res.RawPlan = raw
	return res, nil
}

// Render monta o PLAN.md e o ROUTE.json a partir de tickets já prontos, sem LLM, com a data vinda de now.
// É o mesmo caminho de Generate depois da interpretação dos tickets: a saída é idêntica e determinística.
// Os tickets do chamador não são alterados; a saída traz IDs, épicos e dependências canônicos.
func Render(baseDir, specPath string, tickets []Ticket, level SupervisionLevel, now time.Time) (*PlanResult, error) {
	if level == "" {
		level = SupervisionEpic
	}
	specBytes, err := os.ReadFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler spec: %w", err)
	}
	norm, err := normalizeTickets(tickets, level)
	if err != nil {
		return nil, err
	}
	return materialize(baseDir, specPath, string(specBytes), level, norm, now)
}

// normalizeTickets copia os tickets, canoniza IDs, épicos e dependências, aplica os padrões e recalcula os gates do nível.
func normalizeTickets(in []Ticket, level SupervisionLevel) ([]Ticket, error) {
	if len(in) == 0 {
		return nil, errors.New("nenhum ticket informado")
	}
	seen := map[string]bool{}
	out := make([]Ticket, len(in))
	for i, t := range in {
		t.ID = canonID(t.ID)
		if t.ID == "" || seen[t.ID] {
			return nil, fmt.Errorf("ticket sem ID ou duplicado: %q", t.ID)
		}
		seen[t.ID] = true
		out[i] = completeTicket(t, level)
	}
	return out, nil
}

func completeTicket(t Ticket, level SupervisionLevel) Ticket {
	t.EpicID = canonEpic(t.EpicID, t.ID)
	if t.Nature != "nano" {
		t.Nature = "micro"
	}
	if t.ExecutionMode == "" {
		t.ExecutionMode = "tdd_red_green"
	}
	if t.ValidationCmd == "" {
		t.ValidationCmd = "go test ./..."
	}
	deps := make([]string, len(t.DependsOn))
	for i, d := range t.DependsOn {
		deps[i] = canonID(d)
	}
	t.DependsOn = deps
	t.HumanGateAfter = gateAfter(level, t.Nature)
	return t
}

func planMessages(spec string) []interview.Message {
	user := fmt.Sprintf(`Gere o Plano de Execução do Addy Osmani para a seguinte SPEC:

%s

Decomponha todos os critérios AC-XX em tickets verticais sequenciados logicamente por dependências.`, spec)
	return []interview.Message{
		{Role: "system", Content: PromptPlanSystem},
		{Role: "user", Content: user},
	}
}

// materialize monta a rota e o plano a partir da mesma lista de nós e grava os dois arquivos.
func materialize(baseDir, specPath, specText string, level SupervisionLevel, tickets []Ticket, now time.Time) (*PlanResult, error) {
	ps, err := newPathSet(baseDir, specPath)
	if err != nil {
		return nil, err
	}
	meta := newMeta(ps, specText, now.Format("2006-01-02"), level)
	warnings := resolveDeps(tickets)
	nodes, err := buildNodes(tickets, meta)
	if err != nil {
		return nil, fmt.Errorf("falha ao construir route.json: %w", err)
	}
	routeJSON, err := newRouteDoc(meta, nodes).marshal()
	if err != nil {
		return nil, fmt.Errorf("falha ao construir route.json: %w", err)
	}
	if err := os.WriteFile(ps.planPath, []byte(renderPlan(meta, nodes)), 0644); err != nil {
		return nil, fmt.Errorf("falha ao gravar plan.md: %w", err)
	}
	if err := os.WriteFile(ps.routePath, routeJSON, 0644); err != nil {
		return nil, fmt.Errorf("falha ao gravar route.json: %w", err)
	}
	return &PlanResult{
		Theme: filepath.Base(ps.planPath), SupervisionLevel: level, PlanMarkdownPath: ps.planPath,
		RouteJSONPath: ps.routePath, Tickets: tickets,
		Warnings: append(warnings, uncovered(tickets, meta.SpecCriteria)...),
	}, nil
}

// newPathSet deriva os nomes de PLAN e ROUTE do nome da SPEC e cria docs/plans e docs/routes.
func newPathSet(baseDir, specPath string) (pathSet, error) {
	base := filepath.Base(specPath)
	planName := strings.Replace(base, "SPEC", "PLAN", 1)
	routeName := strings.Replace(strings.Replace(base, "SPEC", "ROUTE", 1), ".md", ".json", 1)
	plansDir, routesDir := filepath.Join(baseDir, "docs", "plans"), filepath.Join(baseDir, "docs", "routes")
	for _, d := range []string{plansDir, routesDir} {
		if err := os.MkdirAll(d, 0755); err != nil {
			return pathSet{}, fmt.Errorf("falha ao criar %s: %w", d, err)
		}
	}
	return pathSet{baseDir: baseDir, specPath: specPath,
		planPath: filepath.Join(plansDir, planName), routePath: filepath.Join(routesDir, routeName)}, nil
}

var ticketHeaderRe = regexp.MustCompile(`(?i)^###\s*\[(TICK(?:ET)?-[A-Za-z0-9_.-]+)\]\s*(.*)`)

type rawBlock struct {
	id    string
	title string
	lines []string
}

// splitBlocks agrupa as linhas por cabeçalho de ticket.
func splitBlocks(content string) []rawBlock {
	var blocks []rawBlock
	var cur *rawBlock
	for l := range strings.SplitSeq(content, "\n") {
		if m := ticketHeaderRe.FindStringSubmatch(strings.TrimSpace(l)); len(m) > 2 {
			if cur != nil {
				blocks = append(blocks, *cur)
			}
			cur = &rawBlock{id: strings.TrimSpace(m[1]), title: strings.TrimSpace(m[2])}
			continue
		}
		if cur != nil {
			cur.lines = append(cur.lines, l)
		}
	}
	if cur != nil {
		blocks = append(blocks, *cur)
	}
	return blocks
}

func blockField(lines []string, name string) string {
	re := regexp.MustCompile(`(?i)^\s*-\s*\*\*` + name + `:\*\*\s*(.*)`)
	for _, l := range lines {
		if m := re.FindStringSubmatch(l); len(m) > 1 {
			return strings.TrimSpace(m[1])
		}
	}
	return ""
}

// splitList separa por vírgula, tira crases e descarta vazios.
func splitList(raw string) []string {
	var out []string
	for p := range strings.SplitSeq(raw, ",") {
		if c := strings.Trim(strings.TrimSpace(p), "`"); c != "" {
			out = append(out, c)
		}
	}
	return out
}

func parseDeps(raw string) []string {
	if raw == "" || strings.Contains(strings.ToLower(raw), "nenhuma") {
		return nil
	}
	var deps []string
	for _, d := range splitList(raw) {
		deps = append(deps, canonID(d))
	}
	return deps
}

func gateAfter(level SupervisionLevel, nature string) bool {
	switch level {
	case SupervisionNano:
		return true
	case SupervisionMicro:
		return nature == "micro"
	}
	return false
}

func parseTickets(content string, level SupervisionLevel) ([]Ticket, error) {
	var tickets []Ticket
	seen := map[string]bool{}
	for _, b := range splitBlocks(content) {
		id := canonID(b.id)
		if seen[id] {
			return nil, fmt.Errorf("ticket duplicado: %s", id)
		}
		seen[id] = true
		nature := strings.ToLower(blockField(b.lines, "Natureza"))
		if nature != "nano" {
			nature = "micro"
		}
		mode := strings.ToLower(blockField(b.lines, "Modo"))
		if mode == "" {
			mode = "tdd_red_green"
		}
		valCmd := blockField(b.lines, "Comando de Validação")
		if valCmd == "" {
			valCmd = "go test ./..."
		}
		tickets = append(tickets, Ticket{
			ID: id, Title: b.title, EpicID: canonEpic(blockField(b.lines, "Épico"), id), Nature: nature, ExecutionMode: mode,
			DependsOn: parseDeps(blockField(b.lines, "Dependências")), TargetFiles: splitList(blockField(b.lines, "Arquivos Afetados")),
			SpecCriteria: splitList(blockField(b.lines, "Critérios Atendidos")), ValidationCmd: valCmd,
			HumanGateAfter: gateAfter(level, nature),
		})
	}
	return tickets, nil
}

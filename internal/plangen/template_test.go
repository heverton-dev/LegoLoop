package plangen_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"legoloop/docs/templates"
	"legoloop/internal/gate"
	"legoloop/internal/interview"
	"legoloop/internal/plangen"
)

const specBody = `# 📋 SPEC: Serviço de Exemplo

> **Referência da Intenção:** [docs/intents/2026-11-02_INTENT_SERVICO.md](../intents/2026-11-02_INTENT_SERVICO.md)

## 5. Critérios Binários de Aceite por Nó / Ticket (BDD)

### Épico 01 — Fundação
* [ ] AC-E01.1: Contrato
* [ ] AC-E01.2: Lint
### Épico 02 — API
* [ ] AC-E02.1: Endpoint
* [ ] AC-E02.9: Critério que nenhum ticket cobre
`

const threeTickets = `
### [TICKET-E01.1] Criar contrato do serviço
- **Épico:** E01
- **Natureza:** micro
- **Modo:** tdd_red_green
- **Dependências:** nenhuma
- **Arquivos Afetados:** internal/service/types.go, internal/service/types_test.go
- **Critérios Atendidos:** AC-E01.1
- **Comando de Validação:** go test -v ./internal/service/...

### [TICKET-E01.2] Formatar o serviço
- **Épico:** E01
- **Natureza:** nano
- **Modo:** sequential
- **Dependências:** TICK-E01.1
- **Arquivos Afetados:** internal/service/types.go
- **Critérios Atendidos:** AC-E01.2
- **Comando de Validação:** go vet ./internal/service/... && gofmt -l internal/service

### [TICKET-E02.1] Expor o endpoint
- **Épico:** E02
- **Natureza:** micro
- **Modo:** tdd_red_green
- **Dependências:** TICKET-E01.2
- **Arquivos Afetados:** internal/api/handler.go, internal/api/handler_test.go
- **Critérios Atendidos:** AC-E02.1
- **Comando de Validação:** go test ./internal/api/... | tee out.txt
`

type fixedLLM struct{ out string }

func (f fixedLLM) Chat(context.Context, []interview.Message) (string, error) { return f.out, nil }

func run(t *testing.T, level plangen.SupervisionLevel) (*plangen.PlanResult, map[string]any, string) {
	t.Helper()
	base := t.TempDir()
	specDir := filepath.Join(base, "docs", "specs")
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatal(err)
	}
	spec := filepath.Join(specDir, "2026-11-02_SPEC_SERVICO.md")
	if err := os.WriteFile(spec, []byte(specBody), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := plangen.NewGenerator(fixedLLM{threeTickets}).Generate(context.Background(), base, spec, level)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	var route map[string]any
	b, _ := os.ReadFile(res.RouteJSONPath)
	if err := json.Unmarshal(b, &route); err != nil {
		t.Fatalf("ROUTE inválido: %v", err)
	}
	plan, _ := os.ReadFile(res.PlanMarkdownPath)
	return res, route, string(plan)
}

func keys(v any) []string {
	m, _ := v.(map[string]any)
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	slices.Sort(out)
	return out
}

func nodesOf(route map[string]any) []map[string]any {
	var out []map[string]any
	for _, n := range route["nodes"].([]any) {
		out = append(out, n.(map[string]any))
	}
	return out
}

func ids(route map[string]any) []string {
	var out []string
	for _, n := range nodesOf(route) {
		out = append(out, n["id"].(string))
	}
	return out
}

func TestRouteFollowsTemplateKeys(t *testing.T) {
	want, err := templates.RouteKeys()
	if err != nil {
		t.Fatal(err)
	}
	_, route, _ := run(t, plangen.SupervisionMicro)
	if got := keys(route); !slices.Equal(got, want.Top) {
		t.Errorf("chaves do topo divergem do template:\n got %v\nwant %v", got, want.Top)
	}
	if got := keys(route["execution_strategy"]); !slices.Equal(got, want.Strategy) {
		t.Errorf("execution_strategy diverge:\n got %v\nwant %v", got, want.Strategy)
	}
	for _, n := range nodesOf(route) {
		checks := map[string][]string{"": want.Node, "isolation": want.Isolation, "routing": want.Routing, "tooling": want.Tooling, "inputs": want.Inputs, "output": want.Output}
		for sub, w := range checks {
			var got []string
			if sub == "" {
				got = keys(n)
			} else {
				got = keys(n[sub])
			}
			if !slices.Equal(got, w) {
				t.Errorf("nó %v, nível %q diverge do template:\n got %v\nwant %v", n["id"], sub, got, w)
			}
		}
	}
	if route["template_version"] != templates.Version {
		t.Errorf("template_version deve ser %s", templates.Version)
	}
}

func TestEveryGeneratedGateCompilesAndNoSubagents(t *testing.T) {
	for _, lvl := range []plangen.SupervisionLevel{plangen.SupervisionNano, plangen.SupervisionMicro, plangen.SupervisionEpic, plangen.SupervisionAutonomous} {
		_, route, _ := run(t, lvl)
		for _, n := range nodesOf(route) {
			g := n["output"].(map[string]any)["gate"].(string)
			if _, err := gate.Compile(g); err != nil {
				t.Errorf("[%s] gate do nó %v não compila: %q: %v", lvl, n["id"], g, err)
			}
			if n["isolation"].(map[string]any)["subagents_allowed"] != false {
				t.Errorf("[%s] nó %v permite subagentes (AGENTS.md, regra 4)", lvl, n["id"])
			}
		}
	}
}

func TestTicketIDsAreCanonicalAndDependenciesResolve(t *testing.T) {
	res, route, _ := run(t, plangen.SupervisionAutonomous)
	got := []string{res.Tickets[0].ID, res.Tickets[1].ID, res.Tickets[2].ID}
	if !slices.Equal(got, []string{"TICK-01.1", "TICK-01.2", "TICK-02.1"}) {
		t.Fatalf("IDs não canônicos: %v", got)
	}
	if !slices.Equal(res.Tickets[1].DependsOn, []string{"TICK-01.1"}) || !slices.Equal(res.Tickets[2].DependsOn, []string{"TICK-01.2"}) {
		t.Errorf("aliases de dependência não foram resolvidos: %v %v", res.Tickets[1].DependsOn, res.Tickets[2].DependsOn)
	}
	known := ids(route)
	for _, n := range nodesOf(route) {
		for _, d := range n["depends_on"].([]any) {
			if !slices.Contains(known, d.(string)) {
				t.Errorf("nó %v depende de %v, que não existe", n["id"], d)
			}
		}
	}
}

func TestSupervisionLevelsInjectTheTemplateGates(t *testing.T) {
	cases := map[plangen.SupervisionLevel][]string{
		plangen.SupervisionAutonomous: {"TICK-01.1", "TICK-01.2", "TICK-02.1", "SG-FINAL"},
		plangen.SupervisionEpic:       {"TICK-01.1", "TICK-01.2", "SG-E01", "TICK-02.1", "SG-FINAL"},
		plangen.SupervisionMicro:      {"TICK-01.1", "SG-01.1", "TICK-01.2", "SG-E01", "TICK-02.1", "SG-02.1", "SG-FINAL"},
		plangen.SupervisionNano:       {"TICK-01.1", "SG-01.1", "TICK-01.2", "SG-01.2", "TICK-02.1", "SG-02.1", "SG-FINAL"},
	}
	for lvl, want := range cases {
		_, route, _ := run(t, lvl)
		if got := ids(route); !slices.Equal(got, want) {
			t.Errorf("[%s] nós divergem:\n got %v\nwant %v", lvl, got, want)
		}
	}
}

func TestHumanGatesBlockTheirDependents(t *testing.T) {
	_, route, _ := run(t, plangen.SupervisionMicro)
	deps := map[string][]any{}
	for _, n := range nodesOf(route) {
		deps[n["id"].(string)] = n["depends_on"].([]any)
	}
	if !slices.Contains(deps["TICK-01.2"], any("SG-01.1")) || slices.Contains(deps["TICK-01.2"], any("TICK-01.1")) {
		t.Errorf("TICK-01.2 deve esperar o gate SG-01.1, não só o ticket: %v", deps["TICK-01.2"])
	}
	if !slices.Contains(deps["TICK-02.1"], any("SG-E01")) {
		t.Errorf("o épico 02 deve esperar o gate do épico 01: %v", deps["TICK-02.1"])
	}
	if !slices.Contains(deps["SG-FINAL"], any("SG-02.1")) {
		t.Errorf("SG-FINAL deve depender do último nó: %v", deps["SG-FINAL"])
	}
}

func TestVerifyIsArgvWithoutShellAndFallsBackToHumanGate(t *testing.T) {
	_, route, _ := run(t, plangen.SupervisionAutonomous)
	out := func(id string) map[string]any {
		for _, n := range nodesOf(route) {
			if n["id"] == id {
				return n["output"].(map[string]any)
			}
		}
		t.Fatalf("nó %s ausente", id)
		return nil
	}
	v := out("TICK-01.2")["verify"].([]any)
	if len(v) != 2 || v[0].([]any)[0] != "go" || v[1].([]any)[0] != "gofmt" {
		t.Errorf("'&&' deve virar dois argv: %v", v)
	}
	if g := out("TICK-02.1")["gate"].(string); !strings.HasPrefix(g, "prompt(") || len(out("TICK-02.1")["verify"].([]any)) != 0 {
		t.Errorf("comando com pipe não pode virar argv: deve cair em gate humano, obteve %q / %v", g, out("TICK-02.1")["verify"])
	}
	if g := out("TICK-01.1")["gate"].(string); !strings.Contains(g, "file_exists('internal/service/types_test.go')") {
		t.Errorf("o gate deve provar os artefatos: %q", g)
	}
}

func TestFinalGateVerifiesGlobalBuild(t *testing.T) {
	_, route, _ := run(t, plangen.SupervisionAutonomous)
	nodes := nodesOf(route)
	last := nodes[len(nodes)-1]
	if last["id"] != "SG-FINAL" || last["type"] != "stop_gate" {
		t.Fatalf("o último nó deve ser SG-FINAL: %v", last["id"])
	}
	b, _ := json.Marshal(last["output"].(map[string]any)["verify"])
	if !strings.Contains(string(b), `["go","build","-o","bin/legoloop.exe","./cmd/legoloop"]`) || !strings.Contains(string(b), `["go","test","./..."]`) {
		t.Errorf("SG-FINAL deve rodar test e build globais: %s", b)
	}
}

func TestReferencesAreRelativeAndTraceable(t *testing.T) {
	_, route, _ := run(t, plangen.SupervisionAutonomous)
	if route["spec_ref"] != "docs/specs/2026-11-02_SPEC_SERVICO.md" || route["plan_ref"] != "docs/plans/2026-11-02_PLAN_SERVICO.md" {
		t.Errorf("referências devem ser relativas ao projeto: %v %v", route["spec_ref"], route["plan_ref"])
	}
	if route["intent_ref"] != "docs/intents/2026-11-02_INTENT_SERVICO.md" {
		t.Errorf("intent_ref deve vir do cabeçalho da SPEC: %v", route["intent_ref"])
	}
	n := nodesOf(route)[0]
	if n["spec_ref"] != "AC-E01.1" || n["inputs"].(map[string]any)["spec_anchor"] != "docs/specs/2026-11-02_SPEC_SERVICO.md#AC-E01.1" {
		t.Errorf("rastreabilidade AC quebrada: %v / %v", n["spec_ref"], n["inputs"])
	}
	if n["test_target"] != "internal/service/types_test.go" {
		t.Errorf("test_target deve ser o arquivo de teste: %v", n["test_target"])
	}
}

func TestPlanFollowsTemplateStructureAndMirrorsRoute(t *testing.T) {
	_, route, plan := run(t, plangen.SupervisionMicro)
	pos := 0
	for _, h := range templates.Headings(templates.Plan()) {
		i := strings.Index(plan[pos:], "## "+h)
		if i < 0 {
			t.Fatalf("seção do template ausente ou fora de ordem: %q\n%s", h, plan)
		}
		pos += i + len(h)
	}
	if strings.Contains(plan, "{{") || !strings.Contains(plan, "STOP-GATE 03") || !strings.Contains(plan, "TEMPLATE_PLAN.md v"+templates.Version) {
		t.Errorf("PLAN fora do template:\n%s", plan)
	}
	for _, n := range nodesOf(route) {
		if !strings.Contains(plan, "`"+n["id"].(string)+"`") {
			t.Errorf("o PLAN não espelha o nó %v", n["id"])
		}
	}
	for _, want := range []string{"Épico 01 — Fundação", "Épico 02 — API", "`180s`", "`600` tokens", "`300` tokens"} {
		if !strings.Contains(plan, want) {
			t.Errorf("PLAN sem %q", want)
		}
	}
}

func TestUncoveredCriteriaAreReported(t *testing.T) {
	res, _, _ := run(t, plangen.SupervisionAutonomous)
	if !slices.ContainsFunc(res.Warnings, func(w string) bool { return strings.Contains(w, "AC-E02.9") }) {
		t.Errorf("o AC sem ticket deve virar aviso: %v", res.Warnings)
	}
}

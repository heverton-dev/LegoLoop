package plangen_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"legoloop/docs/templates"
	"legoloop/internal/gate"
	"legoloop/internal/plangen"
)

var fixedNow = time.Date(2026, 11, 3, 10, 0, 0, 0, time.UTC)

func writeSpec(t *testing.T) (base, spec string) {
	t.Helper()
	base = t.TempDir()
	dir := filepath.Join(base, "docs", "specs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	spec = filepath.Join(dir, "2026-11-03_SPEC_SERVICO.md")
	if err := os.WriteFile(spec, []byte(specBody), 0o644); err != nil {
		t.Fatal(err)
	}
	return base, spec
}

func readBoth(t *testing.T, res *plangen.PlanResult) (string, string) {
	t.Helper()
	p, err := os.ReadFile(res.PlanMarkdownPath)
	if err != nil {
		t.Fatal(err)
	}
	r, err := os.ReadFile(res.RouteJSONPath)
	if err != nil {
		t.Fatal(err)
	}
	return string(p), string(r)
}

func generated(t *testing.T) *plangen.PlanResult {
	t.Helper()
	base, spec := writeSpec(t)
	gen := plangen.NewGenerator(fixedLLM{threeTickets})
	gen.Now = func() time.Time { return fixedNow }
	res, err := gen.Generate(context.Background(), base, spec, plangen.SupervisionMicro)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func TestRender_EqualsGenerateForTheSameTickets(t *testing.T) {
	viaGenerate := generated(t)
	wantPlan, wantRoute := readBoth(t, viaGenerate)

	base, spec := writeSpec(t)
	res, err := plangen.Render(base, spec, viaGenerate.Tickets, plangen.SupervisionMicro, fixedNow)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	gotPlan, gotRoute := readBoth(t, res)
	if gotPlan != wantPlan || gotRoute != wantRoute {
		t.Fatalf("Render diverge de Generate\nPLAN igual: %v | ROUTE igual: %v", gotPlan == wantPlan, gotRoute == wantRoute)
	}
}

func TestRender_IsDeterministicForTheSameNow(t *testing.T) {
	tickets := generated(t).Tickets
	var plans, routes []string
	for range 2 {
		base, spec := writeSpec(t)
		res, err := plangen.Render(base, spec, tickets, plangen.SupervisionEpic, fixedNow)
		if err != nil {
			t.Fatal(err)
		}
		p, r := readBoth(t, res)
		plans, routes = append(plans, p), append(routes, r)
	}
	if plans[0] != plans[1] || routes[0] != routes[1] {
		t.Fatal("duas chamadas com o mesmo now devem gerar bytes idênticos")
	}
	if !strings.Contains(routes[0], `"date": "2026-11-03"`) {
		t.Errorf("a data do ROUTE deve vir do parâmetro now:\n%s", routes[0])
	}
}

func TestRender_FollowsTemplateRouteAndCompilesGates(t *testing.T) {
	want, err := templates.RouteKeys()
	if err != nil {
		t.Fatal(err)
	}
	base, spec := writeSpec(t)
	res, err := plangen.Render(base, spec, generated(t).Tickets, plangen.SupervisionNano, fixedNow)
	if err != nil {
		t.Fatal(err)
	}
	_, routeRaw := readBoth(t, res)
	route := parseRoute(t, routeRaw)
	if got := keys(route); !slices.Equal(got, want.Top) {
		t.Errorf("chaves do topo divergem: %v", got)
	}
	for _, n := range nodesOf(route) {
		if !slices.Equal(keys(n), want.Node) {
			t.Errorf("nó %v diverge do template", n["id"])
		}
		g := n["output"].(map[string]any)["gate"].(string)
		if _, err := gate.Compile(g); err != nil {
			t.Errorf("gate de %v não compila: %v", n["id"], err)
		}
	}
}

func TestRender_DoesNotMutateTheCallersTickets(t *testing.T) {
	in := []plangen.Ticket{{ID: "TICKET-E01.1", Title: "A", EpicID: "E1", Nature: "micro", ValidationCmd: "go test ./..."}}
	base, spec := writeSpec(t)
	res, err := plangen.Render(base, spec, in, plangen.SupervisionAutonomous, fixedNow)
	if err != nil {
		t.Fatal(err)
	}
	if in[0].ID != "TICKET-E01.1" || in[0].EpicID != "E1" {
		t.Errorf("o slice do chamador foi alterado: %+v", in[0])
	}
	if res.Tickets[0].ID != "TICK-01.1" || res.Tickets[0].EpicID != "E01" {
		t.Errorf("a saída deve ter IDs canônicos: %+v", res.Tickets[0])
	}
}

func TestRender_RejectsEmptyTicketsAndMissingSpec(t *testing.T) {
	base, spec := writeSpec(t)
	if _, err := plangen.Render(base, spec, nil, plangen.SupervisionEpic, fixedNow); err == nil {
		t.Error("sem tickets deve falhar")
	}
	tk := []plangen.Ticket{{ID: "TICK-01.1", Title: "A", ValidationCmd: "go test ./..."}}
	if _, err := plangen.Render(base, filepath.Join(base, "nao-existe.md"), tk, plangen.SupervisionEpic, fixedNow); err == nil {
		t.Error("SPEC inexistente deve falhar")
	}
}

func parseRoute(t *testing.T, raw string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		t.Fatalf("ROUTE inválido: %v", err)
	}
	return m
}

// Mesmo formato do template: o ID do critério vem entre crases. Sem reconhecê-lo, o aviso de critério sem
// ticket nunca disparava para SPECs reais.
func TestRender_RecognizesCriteriaInTheTemplateFormat(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "docs", "specs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	spec := filepath.Join(dir, "2026-11-03_SPEC_SERVICO.md")
	body := strings.ReplaceAll(specBody, "* [ ] AC-", "* [ ] `AC-")
	body = regexp.MustCompile("`AC-([A-Za-z0-9_.-]+):").ReplaceAllString(body, "`AC-$1`:")
	if body == specBody || !strings.Contains(body, "`AC-") {
		t.Skip("a SPEC de teste não usa a lista de critérios esperada")
	}
	if err := os.WriteFile(spec, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	tickets := []plangen.Ticket{{ID: "TICK-01.1", Title: "Primeiro", EpicID: "E01", Nature: "micro", ExecutionMode: "sequential",
		TargetFiles: []string{"a.go"}, ValidationCmd: "go test ./...", SpecCriteria: []string{"AC-E01.1"}}}
	res, err := plangen.Render(base, spec, tickets, plangen.SupervisionEpic, fixedNow)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(res.Warnings, "\n"), "AC-E01.2") {
		t.Errorf("os critérios entre crases devem ser reconhecidos (o AC-E01.2 não tem ticket): %v", res.Warnings)
	}
}

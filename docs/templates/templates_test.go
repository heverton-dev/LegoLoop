package templates_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"legoloop/docs/templates"
	"legoloop/internal/assets"
	"legoloop/internal/dag"
	"legoloop/internal/gate"
	"legoloop/internal/loop"
)

func read(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("ler %s: %v", name, err)
	}
	return b
}

func TestJSONTemplatesAreValid(t *testing.T) {
	for _, n := range []string{"TEMPLATE_ROUTE.json", "TEMPLATE_PIECE_GATE.json", "TEMPLATE_PIECE_HOOK.json", "TEMPLATE_PIECE_MCP.json", "TEMPLATE_PIECE_LOOP.json"} {
		var v any
		if err := json.Unmarshal(read(t, n), &v); err != nil {
			t.Errorf("%s não é JSON válido: %v", n, err)
		}
	}
}

func TestEmbeddedTemplatesMatchFiles(t *testing.T) {
	if templates.Intent() != string(read(t, "TEMPLATE_INTENT.md")) ||
		templates.Spec() != string(read(t, "TEMPLATE_SPEC.md")) ||
		templates.Plan() != string(read(t, "TEMPLATE_PLAN.md")) ||
		templates.Report() != string(read(t, "TEMPLATE_REPORT.md")) ||
		string(templates.Route()) != string(read(t, "TEMPLATE_ROUTE.json")) {
		t.Fatal("o conteúdo embutido diverge dos arquivos em docs/templates")
	}
}

func TestDocumentHeadingsAreNumberedAndStable(t *testing.T) {
	want := map[string]int{"intent": 6, "spec": 6, "plan": 4, "report": 6}
	docs := map[string]string{
		"intent": templates.Intent(),
		"spec":   templates.Spec(),
		"plan":   templates.Plan(),
		"report": templates.Report(),
	}
	for name, md := range docs {
		hs := templates.Headings(md)
		if len(hs) != want[name] {
			t.Errorf("%s: esperava %d seções numeradas, obteve %d: %v", name, want[name], len(hs), hs)
		}
	}
}

func TestStripGuidanceRemovesHTMLComments(t *testing.T) {
	out := templates.StripGuidance("a\n<!--\nGUIA\n-->\nb")
	if strings.Contains(out, "GUIA") || !strings.Contains(out, "a") || !strings.Contains(out, "b") {
		t.Fatalf("StripGuidance incorreto: %q", out)
	}
}

func routeGates(t *testing.T) []string {
	t.Helper()
	var r struct {
		Nodes []struct {
			Output struct {
				Gate string `json:"gate"`
			} `json:"output"`
		} `json:"nodes"`
	}
	if err := json.Unmarshal(templates.Route(), &r); err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, n := range r.Nodes {
		out = append(out, n.Output.Gate)
	}
	return out
}

func TestEveryGateExpressionCompiles(t *testing.T) {
	exprs := routeGates(t)
	var rec struct {
		Recipes []struct {
			Expression string `json:"expression"`
		} `json:"recipes"`
	}
	_ = json.Unmarshal(read(t, "TEMPLATE_PIECE_GATE.json"), &rec)
	for _, r := range rec.Recipes {
		exprs = append(exprs, r.Expression)
	}
	var lp loop.Loop
	_ = json.Unmarshal(read(t, "TEMPLATE_PIECE_LOOP.json"), &lp)
	for _, n := range lp.Nodes {
		exprs = append(exprs, n.Output.Gate)
	}
	re := regexp.MustCompile("Gate \\(DSL\\):\\*\\* `([^`]+)`")
	for _, m := range re.FindAllStringSubmatch(templates.Plan(), -1) {
		exprs = append(exprs, m[1])
	}
	if len(exprs) < 15 {
		t.Fatalf("poucas expressões coletadas (%d): o extrator quebrou", len(exprs))
	}
	for _, e := range exprs {
		if _, err := gate.Compile(e); err != nil {
			t.Errorf("gate não compila: %q: %v", e, err)
		}
	}
}

func TestLoopTemplateIsValidDAGAndLoop(t *testing.T) {
	data := read(t, "TEMPLATE_PIECE_LOOP.json")
	if _, err := dag.Parse(data); err != nil {
		t.Fatalf("dag.Parse: %v", err)
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "exemplo-loop-spec-plano.loop.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
	st, err := loop.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	if lp, err := st.Get("exemplo-loop-spec-plano"); err != nil || len(lp.Nodes) != 2 {
		t.Fatalf("loop.Store.Get: %v %+v", err, lp)
	}
}

func put(t *testing.T, home, rel string, data []byte) {
	t.Helper()
	p := filepath.Join(home, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestPiecesAreRecognizedByInventory(t *testing.T) {
	home := t.TempDir()
	skill := strings.NewReplacer("{{nome-da-skill-em-kebab-case}}", "skill-teste", "{{ARTEFATO_FINAL.md}}", "OUT.md").Replace(string(read(t, "TEMPLATE_PIECE_SKILL.md")))
	put(t, home, ".claude/settings.json", read(t, "TEMPLATE_PIECE_HOOK.json"))
	put(t, home, ".cursor/mcp.json", read(t, "TEMPLATE_PIECE_MCP.json"))
	put(t, home, ".claude/skills/skill-teste/SKILL.md", []byte(skill))
	put(t, home, ".aidd/scripts/TEMPLATE_PIECE_SCRIPT.ps1", read(t, "TEMPLATE_PIECE_SCRIPT.ps1"))
	put(t, home, "loops/exemplo-loop-spec-plano.loop.json", read(t, "TEMPLATE_PIECE_LOOP.json"))

	sc := assets.NewScanner(assets.Config{Home: home, LoopsDir: filepath.Join(home, "loops"),
		LookPath: func(string) (string, error) { return "", os.ErrNotExist }, Getenv: func(string) string { return "" }})
	snap := sc.Scan(context.Background())
	for _, k := range []assets.Kind{assets.KindSkill, assets.KindHook, assets.KindMCP, assets.KindScript, assets.KindLoop} {
		if len(snap.Of(k)) == 0 {
			t.Errorf("o inventário não reconheceu o template de %s", k)
		}
	}
	for _, it := range snap.Of(assets.KindHook) {
		if assets.ContainsSecret(it.Description) {
			t.Errorf("hook com segredo literal no comando: %s", it.Description)
		}
	}
}

func TestPlanMirrorsRoute(t *testing.T) {
	var r struct {
		Nodes []struct {
			ID, Nature, Type string
			SpecRef          *string `json:"spec_ref"`
			TimeoutSec       int     `json:"timeout_sec"`
			TokenBudget      struct {
				Max int `json:"max_output_tokens"`
			} `json:"token_budget"`
			Routing struct{ Model string } `json:"routing"`
		} `json:"nodes"`
	}
	if err := json.Unmarshal(templates.Route(), &r); err != nil {
		t.Fatal(err)
	}
	plan := templates.Plan()
	for _, n := range r.Nodes {
		if !strings.Contains(plan, "`"+n.ID+"`") {
			t.Errorf("o PLAN não cita o nó %s", n.ID)
			continue
		}
		if n.Type == "stop_gate" {
			continue
		}
		checks := []string{"`" + n.Nature + "`", "`" + itoa(n.TimeoutSec) + "s`", "`" + itoa(n.TokenBudget.Max) + "` tokens", "`" + n.Routing.Model + "`"}
		if n.SpecRef != nil {
			checks = append(checks, "`"+*n.SpecRef+"`")
		}
		for _, c := range checks {
			if !strings.Contains(plan, c) {
				t.Errorf("nó %s: o PLAN não contém %s", n.ID, c)
			}
		}
	}
}

func itoa(n int) string {
	b, _ := json.Marshal(n)
	return string(b)
}

func TestNativeFragmentsHaveNoPlaceholdersAndScriptHasBOM(t *testing.T) {
	for _, n := range []string{"TEMPLATE_PIECE_HOOK.json", "TEMPLATE_PIECE_LOOP.json"} {
		if strings.Contains(string(read(t, n)), "{{") {
			t.Errorf("%s deve ser exemplo válido, sem placeholders", n)
		}
	}
	if b := read(t, "TEMPLATE_PIECE_SCRIPT.ps1"); !slices.Equal(b[:3], []byte{0xEF, 0xBB, 0xBF}) {
		t.Error("TEMPLATE_PIECE_SCRIPT.ps1 precisa de BOM UTF-8 (PowerShell 5.1)")
	}
}

func TestRouteKeysAreExposedForEngines(t *testing.T) {
	k, err := templates.RouteKeys()
	if err != nil {
		t.Fatal(err)
	}
	for name, set := range map[string][]string{"top": k.Top, "strategy": k.Strategy, "node": k.Node, "isolation": k.Isolation, "routing": k.Routing, "tooling": k.Tooling, "inputs": k.Inputs, "output": k.Output} {
		if len(set) == 0 || !slices.IsSorted(set) {
			t.Errorf("conjunto %s vazio ou não ordenado: %v", name, set)
		}
	}
	for _, must := range []string{"spec_ref", "token_budget", "rollback_strategy", "defects", "tooling", "initiative_id", "epic_id", "retry_policy", "telemetry"} {
		if !slices.Contains(k.Node, must) {
			t.Errorf("a chave de nó %q sumiu do template", must)
		}
	}
	if !slices.Contains(k.Isolation, "write_boundary") {
		t.Errorf("a chave isolation 'write_boundary' deve existir no template v2.2")
	}
	if !slices.Contains(k.Top, "initiative_ref") {
		t.Errorf("a chave top 'initiative_ref' sumiu do template")
	}
	if !slices.Contains(k.Top, "upstream_signatures") {
		t.Errorf("a chave top 'upstream_signatures' deve existir no template v2.1")
	}
	for _, must := range []string{"concurrency", "rollback_policy"} {
		if !slices.Contains(k.Strategy, must) {
			t.Errorf("a chave strategy %q deve existir no template v2.1", must)
		}
	}
	if templates.Version != "2.2.0" {
		t.Errorf("templates.Version deve ser 2.2.0, obteve %s", templates.Version)
	}
}

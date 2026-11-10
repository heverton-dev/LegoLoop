package templates_test

import (
	"context"
	"encoding/json"
	"legoloop/internal/caminhos"
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
	studio "legoloop/internal/loop"
	"legoloop/internal/taxonomy"
)

// Nomes reais em disco da Peça de Processo. Ficam em um só lugar.
const (
	arquivoProcesso = "TEMPLATE_PIECE_LOOP.json"
	idExemplo       = "exemplo-processo-spec-plano"
	sufixoStudio    = ".loop.json"
	dirDocumentos   = "loops"
	// limiteTaxonomia é o mesmo teto do gate do ticket (-tax-max).
	limiteTaxonomia = 38
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
	for _, n := range []string{"TEMPLATE_ROUTE.json", "TEMPLATE_PIECE_GATE.json", "TEMPLATE_PIECE_HOOK.json", "TEMPLATE_PIECE_MCP.json", arquivoProcesso, "TEMPLATE_DAG_VIDEOREA.json", "TEMPLATE_DAG.json"} {
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
		string(templates.Route()) != string(read(t, "TEMPLATE_ROUTE.json")) ||
		string(templates.DAGVideoREA()) != string(read(t, "TEMPLATE_DAG_VIDEOREA.json")) ||
		string(templates.DAG()) != string(read(t, "TEMPLATE_DAG.json")) {
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
	var doc struct {
		Nodes []struct {
			Output struct {
				Gate string `json:"gate"`
			} `json:"output"`
		} `json:"nodes"`
	}
	_ = json.Unmarshal(read(t, arquivoProcesso), &doc)
	for _, n := range doc.Nodes {
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

func TestProcessoTemplateIsValidDAGAndDocument(t *testing.T) {
	data := read(t, arquivoProcesso)
	if _, err := dag.Parse(data); err != nil {
		t.Fatalf("dag.Parse: %v", err)
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, idExemplo+sufixoStudio), data, 0o600); err != nil {
		t.Fatal(err)
	}
	st, err := studio.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	if doc, err := st.Get(idExemplo); err != nil || len(doc.Nodes) != 2 {
		t.Fatalf("Store.Get: %v %+v", err, doc)
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
	put(t, home, dirDocumentos+"/"+idExemplo+sufixoStudio, read(t, arquivoProcesso))

	sc := assets.NewScanner(assets.Config{Home: home, LoopsDir: caminhos.DirLoopsEm(home),
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
	for _, n := range []string{"TEMPLATE_PIECE_HOOK.json", arquivoProcesso} {
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

func TestQuintetoTemplatesConformidade(t *testing.T) {
	quinteto := []struct {
		name string
		data []byte
	}{
		{"TEMPLATE_INTENT.md", []byte(templates.Intent())},
		{"TEMPLATE_SPEC.md", []byte(templates.Spec())},
		{"TEMPLATE_PLAN.md", []byte(templates.Plan())},
		{"TEMPLATE_ROUTE.json", templates.Route()},
		{"TEMPLATE_DAG.json", templates.DAG()},
	}

	for _, item := range quinteto {
		if len(item.data) == 0 {
			t.Errorf("template do quinteto %s não pode estar vazio", item.name)
		}
	}

	var d struct {
		Name  string `json:"name"`
		Nodes []any  `json:"nodes"`
	}
	if err := json.Unmarshal(templates.DAG(), &d); err != nil {
		t.Fatalf("TEMPLATE_DAG.json não é um DAG válido: %v", err)
	}
	if len(d.Nodes) == 0 {
		t.Errorf("TEMPLATE_DAG.json deve conter nós de exemplo")
	}
}

// regra4 devolve o texto da regra 4 do AGENTS.md (até a regra 5).
func regra4(t *testing.T) string {
	t.Helper()
	md := string(read(t, filepath.Join("..", "..", "AGENTS.md")))
	i := strings.Index(md, "4. **Orchestrix")
	j := strings.Index(md, "5. **TDD")
	if i < 0 || j < i {
		t.Fatal("regra 4 do AGENTS.md não localizada")
	}
	return md[i:j]
}

func TestContratoSubagenteEfemero(t *testing.T) {
	r := regra4(t)
	if !strings.Contains(r, "efêmero") || !strings.Contains(r, `subagents_allowed: "ephemeral"`) {
		t.Error("a regra 4 deve declarar o subagente efêmero com subagents_allowed: \"ephemeral\"")
	}
	exc := r[strings.Index(r, "Exceção"):]
	for _, n := range []string{"1.", "2.", "3.", "4.", "5.", "6."} {
		if !regexp.MustCompile(`(?m)^\s+` + regexp.QuoteMeta(n) + ` `).MatchString(exc) {
			t.Errorf("a regra 4 deve listar a condição %s", n)
		}
	}
	if regexp.MustCompile(`(?m)^\s+7\. `).MatchString(exc) {
		t.Error("a regra 4 deve listar exatamente seis condições")
	}
	if !strings.Contains(exc, "`strict`") {
		t.Error("a sexta condição deve exigir write_boundary strict")
	}
	agents := string(read(t, filepath.Join("..", "..", "AGENTS.md")))
	if !strings.Contains(agents, "Exceção de marca") || !strings.Contains(agents, `"LegoLoop"`) {
		t.Error("o item 12 do AGENTS.md deve declarar a exceção de marca LegoLoop")
	}
}

func TestContratoSubagenteEfemeroNosTemplates(t *testing.T) {
	var r struct {
		Nodes []struct {
			ID        string         `json:"id"`
			Extreme   *bool          `json:"extreme"`
			Isolation map[string]any `json:"isolation"`
		} `json:"nodes"`
	}
	if err := json.Unmarshal(templates.Route(), &r); err != nil {
		t.Fatal(err)
	}
	for _, n := range r.Nodes {
		if n.Extreme == nil {
			t.Errorf("nó %s sem o campo extreme", n.ID)
		}
		switch v := n.Isolation["subagents_allowed"].(type) {
		case bool:
			if v {
				t.Errorf("nó %s: subagents_allowed true é proibido", n.ID)
			}
		case string:
			if v != "ephemeral" || n.Isolation["write_boundary"] != "strict" {
				t.Errorf("nó %s: ephemeral exige write_boundary strict", n.ID)
			}
		default:
			t.Errorf("nó %s: subagents_allowed ausente ou inválido", n.ID)
		}
	}
	if !strings.Contains(templates.Plan(), "`ephemeral`") {
		t.Error("TEMPLATE_PLAN.md deve documentar o valor ephemeral")
	}
	src := string(read(t, filepath.Join("..", "..", "internal", "plangen", "route.go")))
	if !strings.Contains(src, "SubagentsAllowed any") || strings.Contains(src, "code-nano") {
		t.Error("plangen/route.go deve aceitar ephemeral e não usar code-nano")
	}
}

// termoProibido diz se a palavra, ou o seu plural simples, é sinônimo proibido.
func termoProibido(w string) bool {
	w = strings.ToLower(w)
	if _, ok := taxonomy.CheckForbidden(w); ok {
		return true
	}
	if len(w) > 2 && strings.HasSuffix(w, "s") {
		_, ok := taxonomy.CheckForbidden(w[:len(w)-1])
		return ok
	}
	return false
}

// TestTemplatesSemTermoProibido cobre o AC-E02.2: o texto dos templates usa só os
// termos canônicos. O que sobra são nomes de contrato (tipos e arquivos em disco).
func TestTemplatesSemTermoProibido(t *testing.T) {
	palavra := regexp.MustCompile(`\p{L}+`)
	total := 0
	for _, padrao := range []string{"*.md", "*.json", "*.typ"} {
		arquivos, err := filepath.Glob(padrao)
		if err != nil {
			t.Fatal(err)
		}
		for _, a := range arquivos {
			n := 0
			for _, w := range palavra.FindAllString(string(read(t, a)), -1) {
				if termoProibido(w) {
					n++
				}
			}
			t.Logf("%s: %d termo(s) proibido(s)", a, n)
			total += n
		}
	}
	if total > limiteTaxonomia {
		t.Errorf("%d termos proibidos nos templates; o limite é %d", total, limiteTaxonomia)
	}
}

func TestGovernanceAC_E00_3(t *testing.T) {
	k, err := templates.RouteKeys()
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(k.Node, "extreme") {
		t.Error("o template de ROUTE deve ter o campo extreme no nó")
	}
	for _, n := range []string{"TEMPLATE_ROUTE.json", "TEMPLATE_ROUTE_LEGOPRESS.json", "TEMPLATE_DAG_LEGOPRESS.json", "TEMPLATE_PLAN.md"} {
		if strings.Contains(string(read(t, n)), "code-nano") {
			t.Errorf("%s ainda usa o combo code-nano, que não existe no 9Router", n)
		}
	}
	schema, err := os.ReadFile(filepath.Join("..", "..", "schemas", "route.v2.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`"ephemeral"`, `"extreme"`} {
		if !strings.Contains(string(schema), want) {
			t.Errorf("schemas/route.v2.json deve aceitar %s", want)
		}
	}
}

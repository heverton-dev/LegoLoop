package studio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/assets"
	"legoloop/internal/dag"
	"legoloop/internal/graph"
	"legoloop/internal/loop"
)

// piecesFixture monta HOME e projeto com uma skill, um script, dois hooks (um com credencial) e um MCP.
func piecesFixture(t *testing.T) (*Server, string, string) {
	t.Helper()
	s, home, project := newAssetsServer(t)
	writeFile(t, filepath.Join(home, ".claude", "skills", "meu-skill", "SKILL.md"), "---\nname: Meu Skill\ndescription: faz coisas\n---\n")
	writeFile(t, filepath.Join(project, "scripts", "limpa.sh"), "echo limpo")
	writeFile(t, filepath.Join(home, ".claude", "settings.json"), `{"hooks":{
	  "Stop":[{"hooks":[{"type":"command","command":"notifier done"}]}],
	  "PreToolUse":[{"matcher":"Bash","hooks":[{"type":"command","command":"rtk run --token=abc123"}]}]}}`)
	writeFile(t, filepath.Join(home, ".claude.json"), `{"mcpServers":{"rea":{"command":"npx","args":["-y","segredo-xyz"]}}}`)
	return s, home, project
}

type itemRef struct{ ID, Name, Kind string }

func listKind(t *testing.T, s *Server, kind string) []itemRef {
	t.Helper()
	var b struct {
		Items []itemRef `json:"items"`
	}
	decode(t, do(t, s.Handler(), "GET", "/api/assets/"+kind, ""), &b)
	return b.Items
}

func idOf(t *testing.T, s *Server, kind, name string) string {
	t.Helper()
	for _, it := range listKind(t, s, kind) {
		if strings.HasPrefix(it.Name, name) {
			return it.ID
		}
	}
	t.Fatalf("item %s/%s não encontrado", kind, name)
	return ""
}

func nodeJSON(id, block, extra string) string {
	return `{"id":"` + id + `","block":"` + block + `","routing":{"mode":"manual","model":"sonnet"}` + extra + `}`
}

func pieceJSON(id, piece, extra string) string {
	return `{"id":"` + id + `","piece":` + piece + `,"routing":{"mode":"manual","model":"sonnet"}` + extra + `}`
}

func comp(nodes []string, edges string) string {
	return `{"name":"pecas","nodes":[` + strings.Join(nodes, ",") + `],"edges":[` + edges + `]}`
}

func post(t *testing.T, s *Server, path, body string) (int, ValidationResult, string) {
	t.Helper()
	rec := do(t, s.Handler(), "POST", path, body)
	var res ValidationResult
	_ = json.Unmarshal(rec.Body.Bytes(), &res)
	return rec.Code, res, rec.Body.String()
}

func firstIssue(res ValidationResult) string {
	if len(res.Issues) == 0 {
		return ""
	}
	return res.Issues[0].Message
}

// ---------- /api/suggest ----------

type suggestBody struct {
	Gates   []struct{ Expr, Label string } `json:"gates"`
	Scripts []itemRef                      `json:"scripts"`
	Hooks   struct {
		Pre, Post []itemRef
		Skipped   int
	} `json:"hooks"`
	Neighbors []struct {
		Kind, ID, Name string
		Installed      bool
	} `json:"neighbors"`
}

func TestSuggestForCatalogBlock(t *testing.T) {
	s, home, project := piecesFixture(t)
	rec := do(t, s.Handler(), "GET", "/api/suggest?block=aidd-spec", "")
	if rec.Code != 200 {
		t.Fatalf("status %d: %s", rec.Code, rec.Body)
	}
	var b suggestBody
	decode(t, rec, &b)
	var exprs []string
	for _, g := range b.Gates {
		exprs = append(exprs, g.Expr)
	}
	if !strings.Contains(strings.Join(exprs, "|"), "has_acceptance_criteria('SPEC.md')") {
		t.Errorf("gates = %v", exprs)
	}
	if len(b.Scripts) != 1 || b.Scripts[0].Name != "limpa.sh" || b.Scripts[0].ID == "" {
		t.Errorf("scripts = %+v", b.Scripts)
	}
	if len(b.Hooks.Post) != 1 || b.Hooks.Post[0].Name != "Stop" {
		t.Errorf("hooks pós = %+v", b.Hooks.Post)
	}
	// o hook com credencial nunca é sugerido (iria parar no pipeline.dag.json)
	for _, h := range append(append([]itemRef{}, b.Hooks.Pre...), b.Hooks.Post...) {
		if strings.HasPrefix(h.Name, "PreToolUse") {
			t.Errorf("hook com credencial sugerido: %+v", h)
		}
	}
	if b.Hooks.Skipped != 1 {
		t.Errorf("Skipped = %d, quer 1 (o hook com credencial)", b.Hooks.Skipped)
	}
	body := rec.Body.String()
	for _, secret := range []string{"abc123", "segredo-xyz"} {
		if strings.Contains(body, secret) {
			t.Errorf("a resposta vaza %q", secret)
		}
	}
	if leaks(body, home) || leaks(body, project) {
		t.Errorf("a resposta vaza caminho absoluto: %s", body)
	}
	var kinds []string
	for _, n := range b.Neighbors {
		kinds = append(kinds, n.Kind)
	}
	if !strings.Contains(strings.Join(kinds, ","), "mcp") {
		t.Errorf("vizinhos sem MCP do inventário: %v", kinds)
	}
}

func TestSuggestForGenericPiece(t *testing.T) {
	s, _, _ := piecesFixture(t)
	rec := do(t, s.Handler(), "GET", "/api/suggest?kind=script&artifact="+url.QueryEscape("SAIDA.md"), "")
	if rec.Code != 200 {
		t.Fatalf("status %d: %s", rec.Code, rec.Body)
	}
	var b suggestBody
	decode(t, rec, &b)
	found := false
	for _, g := range b.Gates {
		found = found || g.Expr == "has_binary_criteria('SAIDA.md')"
	}
	if !found {
		t.Errorf("gates da peça genérica .md = %+v", b.Gates)
	}
}

func TestSuggestValidatesTheRequest(t *testing.T) {
	s, _, _ := piecesFixture(t)
	cases := map[string]int{
		"/api/suggest?block=nao-existe":                                404,
		"/api/suggest?kind=banana&artifact=A.md":                       400,
		"/api/suggest":                                                 400,
		"/api/suggest?kind=skill":                                      400,
		"/api/suggest?block=interview&kind=skill":                      400,
		"/api/suggest?kind=skill&artifact=" + strings.Repeat("a", 129): 400,
	}
	for path, want := range cases {
		if rec := do(t, s.Handler(), "GET", path, ""); rec.Code != want {
			t.Errorf("%s: status %d, quer %d", path, rec.Code, want)
		}
	}
	bare, _ := newServer(t)
	if rec := do(t, bare.Handler(), "GET", "/api/suggest?block=interview", ""); rec.Code != 503 {
		t.Errorf("sem inventário: status %d, quer 503", rec.Code)
	}
}

// ---------- composição: slots e execução ----------

func TestComposeResolvesSlotsOnTheServerAndMasksPathsInTheResponse(t *testing.T) {
	s, home, project := piecesFixture(t)
	script, stop := idOf(t, s, "script", "limpa.sh"), idOf(t, s, "hook", "Stop")
	slots := fmt.Sprintf(`,"slots":{"gates":["word_count('BRIEFING.md') > 100"],"scripts":[{"ref":%q}],
	  "hooks":{"pre":[{"command":"notifier start \"a b\""}],"post":[{"ref":%q}]}},"execution":{"mode":"isolated"}`, script, stop)
	body := comp([]string{nodeJSON("n1", "interview", slots)}, "")

	code, res, raw := post(t, s, "/api/validate", body)
	if code != 200 || !res.OK {
		t.Fatalf("validate: %d %s", code, raw)
	}
	if leaks(raw, home) || leaks(raw, project) {
		t.Errorf("a resposta expõe caminho absoluto: %s", raw)
	}
	if code, _, raw := post(t, s, "/api/materialize", body); code != 200 {
		t.Fatalf("materialize: %d %s", code, raw)
	}
	data, err := os.ReadFile(filepath.Join(s.ProjectDir, "pipeline.dag.json"))
	if err != nil {
		t.Fatal(err)
	}
	d, err := dag.Parse(data)
	if err != nil {
		t.Fatalf("pipeline gravado inválido: %v", err)
	}
	n := d.Node("n1")
	if n.Slots == nil || n.EffectiveMode() != dag.ExecIsolated {
		t.Fatalf("slots/execução perdidos: %+v", n)
	}
	if got := n.Slots.Scripts[0]; got[0] != "bash" || !strings.HasSuffix(filepath.ToSlash(got[1]), "scripts/limpa.sh") {
		t.Errorf("script resolvido = %q", got)
	}
	if got := n.Slots.Hooks.Pre[0]; len(got) != 3 || got[0] != "notifier" || got[2] != "a b" {
		t.Errorf("hook pré (comando livre) = %q", got)
	}
	if got := n.Slots.Hooks.Post[0]; len(got) != 2 || got[0] != "notifier" || got[1] != "done" {
		t.Errorf("hook pós (do inventário) = %q", got)
	}
	if len(n.Slots.Gates) != 1 {
		t.Errorf("gates = %v", n.Slots.Gates)
	}
}

func TestComposeSlotErrorsAreReportedOnTheNode(t *testing.T) {
	s, _, _ := piecesFixture(t)
	skill, script := idOf(t, s, "skill", "Meu Skill"), idOf(t, s, "script", "limpa.sh")
	secret := idOf(t, s, "hook", "PreToolUse")
	cases := map[string]string{
		"ref desconhecida":            `"scripts":[{"ref":"zzzzzzzzzzzz"}]`,
		"ref de tipo errado":          fmt.Sprintf(`"scripts":[{"ref":%q}]`, skill),
		"hook com credencial":         fmt.Sprintf(`"hooks":{"pre":[{"ref":%q}]}`, secret),
		"ref e comando juntos":        fmt.Sprintf(`"scripts":[{"ref":%q,"command":"notifier x"}]`, script),
		"slot vazio":                  `"scripts":[{}]`,
		"comando com pipe":            `"hooks":{"post":[{"command":"a | b"}]}`,
		"comando relativo":            `"hooks":{"post":[{"command":"tools/x.exe run"}]}`,
		"gate inválido":               `"gates":["exit_code =="]`,
		"gates demais":                `"gates":[` + strings.TrimSuffix(strings.Repeat(`"true",`, dag.MaxSlotEntries+1), ",") + `]`,
		"script ref em slot de hooks": fmt.Sprintf(`"hooks":{"pre":[{"ref":%q}]}`, script),
	}
	for name, slots := range cases {
		_, res, raw := post(t, s, "/api/validate", comp([]string{nodeJSON("n1", "interview", `,"slots":{`+slots+`}`)}, ""))
		if res.OK || len(res.Issues) == 0 {
			t.Errorf("%s: deveria reprovar: %s", name, raw)
			continue
		}
		if res.Issues[0].Node != "n1" && !strings.Contains(raw, "slot") {
			t.Errorf("%s: issue sem nó nem menção ao slot: %s", name, raw)
		}
	}
}

func TestComposeRejectsInvalidExecutionMode(t *testing.T) {
	s, _, _ := piecesFixture(t)
	_, res, _ := post(t, s, "/api/validate", comp([]string{nodeJSON("n1", "interview", `,"execution":{"mode":"yolo"}`)}, ""))
	if res.OK || !strings.Contains(firstIssue(res), "execution") {
		t.Errorf("modo inválido deveria reprovar: %+v", res)
	}
}

// ---------- composição: peças genéricas ----------

func TestComposeGenericSkillPieceBetweenTypedBlocks(t *testing.T) {
	s, _, _ := piecesFixture(t)
	skill := idOf(t, s, "skill", "Meu Skill")
	body := comp([]string{
		nodeJSON("n1", "interview", ""),
		pieceJSON("n2", fmt.Sprintf(`{"kind":"skill","ref":%q,"artifact":"ANALISE.md"}`, skill), ""),
		nodeJSON("n3", "aidd-spec", ""),
	}, `{"from":"n1","to":"n2"},{"from":"n2","to":"n3"}`)
	code, res, raw := post(t, s, "/api/materialize", body)
	if code != 200 || !res.OK {
		t.Fatalf("materialize: %d %s", code, raw)
	}
	data, _ := os.ReadFile(filepath.Join(s.ProjectDir, "pipeline.dag.json"))
	d, err := dag.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	n2 := d.Node("n2")
	if n2.Type != dag.TypeSkill || n2.Ref != "skills/meu-skill" {
		t.Errorf("o ref usa o NOME DA PASTA da skill, não o do frontmatter: %+v", n2)
	}
	if n2.Output.Artifact != "ANALISE.md" || n2.Output.Gate != "file_exists('ANALISE.md')" {
		t.Errorf("handoff/gate padrão = %+v", n2.Output)
	}
	if n2.Inputs["n1"] != "file('BRIEFING.md')" || len(n2.DependsOn) != 1 {
		t.Errorf("entradas = %+v dep=%v", n2.Inputs, n2.DependsOn)
	}
	if got := d.Node("n3").Inputs["n2"]; got != "file('ANALISE.md')" {
		t.Errorf("n3 deveria ler o artefato da peça genérica: %v", d.Node("n3").Inputs)
	}
}

func TestComposeGenericPieceUsesCustomGateAndRootInput(t *testing.T) {
	s, _, _ := piecesFixture(t)
	skill := idOf(t, s, "skill", "Meu Skill")
	body := comp([]string{pieceJSON("n1", fmt.Sprintf(
		`{"kind":"skill","ref":%q,"artifact":"X.md","gate":"word_count('X.md') > 5"}`, skill), "")}, "")
	if code, res, raw := post(t, s, "/api/materialize", body); code != 200 || !res.OK {
		t.Fatalf("materialize: %d %s", code, raw)
	}
	data, _ := os.ReadFile(filepath.Join(s.ProjectDir, "pipeline.dag.json"))
	d, _ := dag.Parse(data)
	n := d.Node("n1")
	if n.Output.Gate != "word_count('X.md') > 5" || n.Inputs["input"] != "var.input" {
		t.Errorf("gate/entrada raiz = %+v %+v", n.Output, n.Inputs)
	}
}

func TestComposeGenericScriptPieceBecomesAScriptNode(t *testing.T) {
	s, _, _ := piecesFixture(t)
	script := idOf(t, s, "script", "limpa.sh")
	body := comp([]string{pieceJSON("n1", fmt.Sprintf(`{"kind":"script","ref":%q,"artifact":"LIMPO.txt"}`, script), "")}, "")
	if code, res, raw := post(t, s, "/api/materialize", body); code != 200 || !res.OK {
		t.Fatalf("materialize: %d %s", code, raw)
	}
	data, _ := os.ReadFile(filepath.Join(s.ProjectDir, "pipeline.dag.json"))
	d, _ := dag.Parse(data)
	n := d.Node("n1")
	if n.Type != dag.TypeScript || len(n.Command) != 2 || n.Command[0] != "bash" {
		t.Errorf("nó script = %+v", n)
	}
}

func TestComposePieceErrors(t *testing.T) {
	s, _, _ := piecesFixture(t)
	skill, script := idOf(t, s, "skill", "Meu Skill"), idOf(t, s, "script", "limpa.sh")
	mcp := idOf(t, s, "mcp", "rea")
	one := func(piece string) string { return comp([]string{pieceJSON("n1", piece, "")}, "") }
	cases := map[string]struct{ body, want string }{
		"ref desconhecida":         {one(`{"kind":"skill","ref":"zzzzzzzzzzzz","artifact":"A.md"}`), "inventário"},
		"tipo errado (script)":     {one(fmt.Sprintf(`{"kind":"skill","ref":%q,"artifact":"A.md"}`, script)), "skill"},
		"tipo errado (skill)":      {one(fmt.Sprintf(`{"kind":"script","ref":%q,"artifact":"A.md"}`, skill)), "script"},
		"MCP não executa":          {one(fmt.Sprintf(`{"kind":"mcp","ref":%q,"artifact":"A.md"}`, mcp)), "MCP"},
		"subloop ainda não":        {one(`{"kind":"subloop","ref":"x","artifact":"A.md"}`), "acervo"},
		"tipo desconhecido":        {one(`{"kind":"banana","ref":"x","artifact":"A.md"}`), "banana"},
		"sem artefato":             {one(fmt.Sprintf(`{"kind":"skill","ref":%q}`, skill)), "artefato"},
		"artefato com aspas":       {one(fmt.Sprintf(`{"kind":"skill","ref":%q,"artifact":"a'b.md"}`, skill)), "artefato"},
		"artefato fora do projeto": {one(fmt.Sprintf(`{"kind":"skill","ref":%q,"artifact":"../x.md"}`, skill)), "artifact"},
		"bloco e peça juntos":      {comp([]string{`{"id":"n1","block":"interview","piece":{"kind":"skill","ref":"x","artifact":"A.md"},"routing":{"mode":"manual","model":"sonnet"}}`}, ""), "bloco"},
		"nem bloco nem peça":       {comp([]string{`{"id":"n1","routing":{"mode":"manual","model":"sonnet"}}`}, ""), "bloco"},
	}
	for name, c := range cases {
		code, res, raw := post(t, s, "/api/validate", c.body)
		if code != 200 || res.OK || !strings.Contains(strings.ToLower(raw), strings.ToLower(c.want)) {
			t.Errorf("%s: esperava reprovação citando %q: %d %s", name, c.want, code, raw)
		}
	}
}

func TestComposePieceNeedsTheInventory(t *testing.T) {
	bare, _ := newServer(t)
	body := comp([]string{pieceJSON("n1", `{"kind":"skill","ref":"abcdefabcdef","artifact":"A.md"}`, "")}, "")
	if _, res, _ := post(t, bare, "/api/validate", body); res.OK {
		t.Error("sem inventário, uma peça do inventário não pode validar")
	}
}

func TestComposeSubloopSuccess(t *testing.T) {
	s, home, _ := newAssetsServer(t)
	loopsDir := filepath.Join(home, ".aidd", "canonical", "loops")
	store, err := loop.NewStore(loopsDir)
	if err != nil {
		t.Fatal(err)
	}
	s.Loops = store

	lp := &loop.Loop{
		ID:      "sub-pipe",
		Name:    "Sub Pipe",
		Version: "2.0.0",
		Trigger: &dag.Trigger{Type: dag.TriggerManual},
		Nodes: []dag.Node{
			{
				ID:      "n1",
				Type:    dag.TypeScript,
				Command: []string{"go", "version"},
				Routing: dag.RoutingConfig{Mode: dag.ModeManual, Model: "local"},
				Output:  dag.HandoffContract{Artifact: "SUB_OUT.md", Gate: "true"},
			},
		},
	}
	if err := store.Save(lp); err != nil {
		t.Fatal(err)
	}
	s.Assets.Refresh(context.Background())

	// Acha o ID público no inventário
	var loopItem assets.Item
	for _, it := range s.Assets.Snapshot().Of(assets.KindLoop) {
		if it.Name == "Sub Pipe" {
			loopItem = it
			break
		}
	}
	if loopItem.ID == "" {
		t.Fatal("loop não encontrado no inventário de assets")
	}

	body := comp([]string{pieceJSON("n1", fmt.Sprintf(`{"kind":"subloop","ref":%q,"artifact":"SUB_OUT.md"}`, loopItem.ID), "")}, "")
	code, res, raw := post(t, s, "/api/validate", body)
	if code != 200 || !res.OK {
		t.Fatalf("esperava validação com sucesso do subloop: %d %s", code, raw)
	}
	if !strings.Contains(raw, `"type":"subloop"`) || !strings.Contains(raw, `"ref":"sub-pipe"`) {
		t.Fatalf("DAG gerado não contém nó subloop correto: %s", raw)
	}
}

func TestBoardTriggerPiece(t *testing.T) {
	s, _, _ := piecesFixture(t)
	// Inicializa o index com FullCatalog para conter o bloco github-board-trigger
	fullIdx, err := graph.NewIndex(graph.FullCatalog())
	if err != nil {
		t.Fatalf("falha ao criar index com FullCatalog: %v", err)
	}
	s.Index = fullIdx

	body := comp([]string{nodeJSON("trig1", "github-board-trigger", "")}, "")
	code, res, raw := post(t, s, "/api/validate", body)
	if code != 200 || !res.OK {
		t.Fatalf("esperava validação com sucesso do nó trigger: %d %s", code, raw)
	}
	if !strings.Contains(raw, `"artifact":"WAVE_PLAN.json"`) {
		t.Fatalf("DAG gerada não contém artefato WAVE_PLAN.json esperado: %s", raw)
	}
}

func TestTriggerBlocksCatalog(t *testing.T) {
	trigs := graph.TriggerBlocks()
	if len(trigs) < 5 {
		t.Fatalf("esperava ao menos 5 blocos de trigger, obteve %d", len(trigs))
	}
	expectedIDs := map[string]bool{
		"trigger-manual":       false,
		"trigger-cron":         false,
		"trigger-file-change":  false,
		"trigger-loop-end":     false,
		"github-board-trigger": false,
	}
	for _, b := range trigs {
		if _, ok := expectedIDs[b.ID]; ok {
			expectedIDs[b.ID] = true
		}
		if !strings.HasPrefix(b.Kind, "trigger_") {
			t.Errorf("bloco trigger %q deve ter Kind iniciando com 'trigger_', obteve %q", b.ID, b.Kind)
		}
	}
	for id, found := range expectedIDs {
		if !found {
			t.Errorf("trigger esperado %q não encontrado em TriggerBlocks()", id)
		}
	}

	// Verifica se estão no FullCatalog
	full := graph.FullCatalog()
	foundInFull := 0
	for _, b := range full {
		if strings.HasPrefix(b.Kind, "trigger_") {
			foundInFull++
		}
	}
	if foundInFull < 5 {
		t.Errorf("esperava ao menos 5 triggers no FullCatalog(), obteve %d", foundInFull)
	}
}

func TestComposeWithTriggerPiece(t *testing.T) {
	s, _, _ := piecesFixture(t)
	fullIdx, err := graph.NewIndex(graph.FullCatalog())
	if err != nil {
		t.Fatalf("falha ao criar index: %v", err)
	}
	s.Index = fullIdx

	// Teste com trigger-cron
	body := comp([]string{nodeJSON("trigCron", "trigger-cron", "")}, "")
	code, res, raw := post(t, s, "/api/validate", body)
	if code != 200 || !res.OK {
		t.Fatalf("esperava 200 OK na validação com trigger-cron: %d %s", code, raw)
	}
	if !strings.Contains(raw, `"trigger":{"type":"cron","cron":"0 * * * *"}`) {
		t.Fatalf("esperava DAG com trigger cron compilado, obteve: %s", raw)
	}
}




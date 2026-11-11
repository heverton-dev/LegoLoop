package studio

import (
	"encoding/json"
	"strings"
	"testing"
)

const paletteFixture = `
const blocks = [
  {id:"interview", label:"Interview", ref:"skills/interview", consumes:["Artifact:Idea"], produces:"Artifact:Briefing"},
  {id:"aidd-spec", label:"AIDD Spec", ref:"skills/aidd-spec", consumes:["Artifact:Briefing","Artifact:DesignTokens"], produces:"Artifact:Spec"}];
const inv = {
  skill: [{id:"s1",name:"interview"},{id:"s2",name:"outra-skill",description:"faz outra coisa"},{id:"s3",name:"AIDD-Spec"}],
  mcp: [{id:"m1",name:"rea",origin:"claude",meta:{transport:"stdio"}}],
  script: [{id:"c1",name:"limpa.sh",scope:"project"}],
  hook: [{id:"h1",name:"Stop",description:"notifier done",meta:{event:"Stop"}}],
  gate: [{id:"g1",name:"file_exists",description:"Arquivo existe",meta:{signature:"file_exists(caminho)"}}],
  loop: []};
const model = buildPaletteModel(blocks, inv);`

// R-UI-6/7: categorias com contadores, catálogo sem duplicata do inventário, busca por palavras.
func TestPaletteModelCategoriesCountsAndSearch(t *testing.T) {
	out := runNode(t, `(() => { `+paletteFixture+`
	  const names = (l) => l.filter((e) => !e.header).map((e) => e.name);
	  return {
	    counts: paletteCounts(model),
	    skills: names(flattenPalette(model, "skill", "")),
	    firstIsContract: model.skill[0].kind === "block" && model.skill[0].badge === "contrato",
	    all: flattenPalette(model, "all", "").map((e) => e.header ? "#" + e.label + ":" + e.count : e.name),
	    loop: flattenPalette(model, "loop", "").length,
	    outra: names(flattenPalette(model, "skill", "outra coisa")),
	    spec: names(flattenPalette(model, "all", "spec")),
	    none: flattenPalette(model, "all", "zzzxx").length,
	    hook: model.hook[0].event, gate: model.gate[0].gateName,
	    emptySearch: paletteEmptyMessage("skill", "zzz").title, emptyCat: paletteEmptyMessage("loop", "").text,
	  };
	})()`, "search.js", "virtual.js", "palette.js")
	var r struct {
		Counts          map[string]int `json:"counts"`
		Skills, All     []string
		Outra, Spec     []string
		FirstIsContract bool
		Loop, None      int
		Hook, Gate      string
		EmptySearch     string
		EmptyCat        string
	}
	if err := json.Unmarshal(out, &r); err != nil {
		t.Fatalf("saída inesperada: %s", out)
	}
	want := map[string]int{"all": 7, "skill": 3, "mcp": 1, "loop": 0, "script": 1, "hook": 1, "gate": 1}
	for k, n := range want {
		if r.Counts[k] != n {
			t.Errorf("contador %s = %d, quer %d (%v)", k, r.Counts[k], n, r.Counts)
		}
	}
	if strings.Join(r.Skills, ",") != "Interview,AIDD Spec,outra-skill" || !r.FirstIsContract {
		t.Errorf("skills (Peças do catálogo primeiro, sem duplicar o inventário) = %v", r.Skills)
	}
	if strings.Join(r.All, "|") != "#Skills:3|Interview|AIDD Spec|outra-skill|#MCPs:1|rea|#Scripts:1|limpa.sh|#Hooks:1|Stop|#Gates:1|file_exists" {
		t.Errorf("Todas (agrupada, sem grupos vazios) = %v", r.All)
	}
	if r.Loop != 0 || r.None != 0 {
		t.Errorf("processo=%d none=%d", r.Loop, r.None)
	}
	if strings.Join(r.Outra, ",") != "outra-skill" || strings.Join(r.Spec, ",") != "AIDD Spec" {
		t.Errorf("busca: outra=%v spec=%v", r.Outra, r.Spec)
	}
	if r.Hook != "Stop" || r.Gate != "file_exists" {
		t.Errorf("campos de encaixe: hook=%q gate=%q", r.Hook, r.Gate)
	}
	if !strings.Contains(r.EmptySearch, "zzz") || !strings.Contains(r.EmptyCat, "acervo") {
		t.Errorf("estados vazios: %q / %q", r.EmptySearch, r.EmptyCat)
	}
}

// Com 5.000 itens a paleta rende dezenas de linhas, não milhares, e a busca é rápida.
func TestPaletteHandlesFiveThousandItemsThroughTheVirtualWindow(t *testing.T) {
	out := runNode(t, `(() => {
	  const skill = Array.from({length: 5000}, (_, i) => ({id: "s" + i, name: "skill-" + String(i).padStart(4, "0")}));
	  const model = buildPaletteModel([], {skill});
	  const flat = flattenPalette(model, "skill", "");
	  const t0 = Date.now();
	  const hits = flattenPalette(model, "all", "skill 4999").length;
	  const window = virtualWindow(flat.length, PALETTE_ROW, 0, 700, 6);
	  return {total: flat.length, rendered: window.end - window.start, hits, ms: Date.now() - t0};
	})()`, "search.js", "virtual.js", "palette.js")
	var r struct{ Total, Rendered, Hits, Ms int }
	if err := json.Unmarshal(out, &r); err != nil {
		t.Fatalf("saída inesperada: %s", out)
	}
	if r.Total != 5000 || r.Rendered > 100 || r.Hits != 2 /* cabeçalho + 1 item */ {
		t.Errorf("5.000 itens: %+v", r)
	}
	if r.Ms > 100 {
		t.Errorf("a busca em 5.000 itens levou %d ms (limite 100)", r.Ms)
	}
}

// R-NODE-2/9: lógica pura dos slots (encaixe de gates, hooks por evento, serialização opcional).
func TestSlotsPureLogic(t *testing.T) {
	out := runNode(t, `(() => ({
	  gates: ["file_exists","word_count","line_count","has_binary_criteria","has_acceptance_criteria","has_pattern","desconhecido"]
	    .map((n) => instantiateGate(n, "SPEC.md")),
	  quote: instantiateGate("file_exists", "a'b.md"),
	  hooks: ["PreToolUse","SessionStart","UserPromptSubmit","PreCompact","PostToolUse","Stop","SubagentStop","SessionEnd","Elicitation",""]
	    .map((e) => { const s = hookSlotsFor(e); return (s.pre ? "P" : "-") + (s.post ? "S" : "-"); }),
	  empty: slotsToComposition(emptySlots()) === undefined,
	  full: slotsToComposition({gates: ["g"], scripts: [{ref: "r1", label: "x"}], pre: [{command: "a b", label: "c"}], post: [{ref: "r2"}]}),
	  onlyGate: slotsToComposition({gates: ["g"], scripts: [], pre: [], post: []}),
	  counts: slotCount({gates: ["g"], scripts: [{}], pre: [{}, {}], post: []}),
	}))()`, "slots.js")
	var r struct {
		Gates  []*string
		Quote  *string
		Hooks  []string
		Empty  bool
		Full   map[string]any
		OnlyG  map[string]any `json:"onlyGate"`
		Counts int
	}
	if err := json.Unmarshal(out, &r); err != nil {
		t.Fatalf("saída inesperada: %s", out)
	}
	wantGates := []string{"file_exists('SPEC.md')", "word_count('SPEC.md') > 100", "line_count('SPEC.md') > 10",
		"has_binary_criteria('SPEC.md')", "has_acceptance_criteria('SPEC.md')", "has_pattern('SPEC.md', '')"}
	for i, w := range wantGates {
		if r.Gates[i] == nil || *r.Gates[i] != w {
			t.Errorf("gate %d = %v, quer %q", i, r.Gates[i], w)
		}
	}
	if r.Gates[6] != nil || r.Quote != nil {
		t.Errorf("gate desconhecido ou artefato com aspas deveria dar null: %v %v", r.Gates[6], r.Quote)
	}
	if strings.Join(r.Hooks, ",") != "P-,P-,P-,P-,-S,-S,-S,-S,PS,PS" {
		t.Errorf("slots por evento = %v", r.Hooks)
	}
	if !r.Empty || r.OnlyG["gates"] == nil || r.OnlyG["scripts"] != nil || r.OnlyG["hooks"] != nil || r.Counts != 4 {
		t.Errorf("serialização/contagem: empty=%v onlyGate=%v counts=%d", r.Empty, r.OnlyG, r.Counts)
	}
	hooks, _ := r.Full["hooks"].(map[string]any)
	scripts, _ := r.Full["scripts"].([]any)
	if len(scripts) != 1 || hooks["pre"] == nil || hooks["post"] == nil {
		t.Fatalf("full = %v", r.Full)
	}
	if pre := hooks["pre"].([]any)[0].(map[string]any); pre["command"] != "a b" || pre["label"] != nil || pre["ref"] != nil {
		t.Errorf("pré (comando livre) = %v — só command vai ao servidor", pre)
	}
	if post := hooks["post"].([]any)[0].(map[string]any); post["ref"] != "r2" || post["command"] != nil {
		t.Errorf("pós (ref) = %v", post)
	}
}

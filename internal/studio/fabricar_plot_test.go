package studio

import (
	"encoding/json"
	"strings"
	"testing"
)

const plotFixture = `
const result = {canonical_name:"legoloop-extrator", slug:"extrator",
  artifacts:{intent_path:"<projeto>/docs/intents/2026-11-03_INTENT_EXTRATOR.md", spec_path:"<projeto>/docs/specs/2026-11-03_SPEC_EXTRATOR.md",
    plan_path:"<projeto>/docs/plans/2026-11-03_PLAN_EXTRATOR.md", route_path:"<projeto>/docs/routes/2026-11-03_ROUTE_EXTRATOR.json",
    template_path:"<projeto>/docs/templates/TEMPLATE_EXTRATOR.md", skill_path:"~/.aidd/canonical/skills/legoloop-extrator/SKILL.md",
    loop_path:"~/.aidd/canonical/loops/legoloop-extrator.loop.json", dag_path:"~/.aidd/canonical/loops/legoloop-extrator.dag.json"},
  curated_pieces:[{name:"alfa",kind:"skill",tier:"ouro"},{name:"beta",kind:"skill",tier:"prata"}],
  warnings:["limite de 5 passos atingido"]};
const loop = {id:"legoloop-extrator", name:"legoloop-extrator", nodes:[
  {id:"n01_alfa", type:"skill", ref:"skills/alfa", routing:{mode:"suggested", model:"code-fast", fallback:"code-pro"}, output:{artifact:"ALFA.md", gate:"file_exists('ALFA.md')"}},
  {id:"n02_beta", type:"skill", ref:"skills/beta", depends_on:["n01_alfa"], routing:{mode:"suggested", model:"code-fast"}, output:{artifact:"BETA.md", gate:"file_exists('BETA.md')"}},
  {id:"n03_gama", type:"skill", ref:"skills/gama", depends_on:["n02_beta","fantasma"], routing:{mode:"suggested", model:"code-fast"}, output:{artifact:"GAMA.md", gate:"file_exists('GAMA.md')"}}]};`

type plotNode struct {
	ID, Label string
	X, Y      float64
	Piece     struct{ Kind, Ref, Artifact, Gate string }
	Routing   struct{ Mode, Model, Fallback string }
	Execution string
}

type plotOut struct {
	Name  string
	Nodes []plotNode
	Edges []struct{ From, To string }
	Entry struct{ Cat, Name, ItemID string }
	Insp  struct {
		Title string
		Rows  []struct{ Label, Value string }
		Notes []string
	}
	Empty   struct{ Nodes, Edges int }
	NullLp  struct{ Nodes, Edges int }
	NodeW   float64
	Dup     int
	SameRef bool
}

func TestFabricarPlotModel(t *testing.T) {
	out := runNode(t, `(() => { `+plotFixture+`
	  const m = fabricarPlotModel(result, loop);
	  const e = fabricarPlotModel(result, {id:"x", nodes:[]});
	  const n = fabricarPlotModel(result, null);
	  return {name: m.name, nodes: m.nodes, edges: m.edges, entry: m.loopEntry,
	    insp: {title: m.inspector.title, rows: m.inspector.rows, notes: m.inspector.notes},
	    empty: {nodes: e.nodes.length, edges: e.edges.length}, nullLp: {nodes: n.nodes.length, edges: n.edges.length},
	    dup: new Set(m.nodes.map((x) => x.id)).size};
	})()`, "fabricar.js")
	var r plotOut
	if err := json.Unmarshal(out, &r); err != nil {
		t.Fatalf("saída inesperada: %s", out)
	}
	if r.Name != "legoloop-extrator" || len(r.Nodes) != 3 || r.Dup != 3 {
		t.Fatalf("o Canvas deve plotar os 3 nós do loop gerado, com ids únicos: %+v", r)
	}
	n0, n1 := r.Nodes[0], r.Nodes[1]
	if n0.ID != "n01_alfa" || n0.Label != "alfa" || n0.Piece.Kind != "skill" || n0.Piece.Ref != "skills/alfa" || n0.Piece.Artifact != "ALFA.md" || !strings.Contains(n0.Piece.Gate, "ALFA.md") {
		t.Errorf("nó 1 não reflete a peça do DAG: %+v", n0)
	}
	if n0.Routing.Model != "code-fast" || n0.Routing.Fallback != "code-pro" || n0.Execution != "inline" {
		t.Errorf("o roteamento do DAG deve ser preservado: %+v", n0)
	}
	if !(n1.X > n0.X) || n1.X-n0.X < 236 {
		t.Errorf("os nós devem seguir da esquerda para a direita sem se sobrepor (largura 236): x0=%v x1=%v", n0.X, n1.X)
	}
	if len(r.Edges) != 2 || r.Edges[0].From != "n01_alfa" || r.Edges[0].To != "n02_beta" || r.Edges[1].From != "n02_beta" || r.Edges[1].To != "n03_gama" {
		t.Errorf("arestas vêm de depends_on, ignorando dependência inexistente: %+v", r.Edges)
	}
	if r.Entry.Cat != "loop" || r.Entry.Name != "legoloop-extrator" || r.Entry.ItemID == "" {
		t.Errorf("a peça nova deve entrar na categoria Loops: %+v", r.Entry)
	}
	if r.Insp.Title != "legoloop-extrator" {
		t.Errorf("título do Inspetor: %q", r.Insp.Title)
	}
	labels := map[string]string{}
	for _, row := range r.Insp.Rows {
		labels[row.Label] = row.Value
	}
	for _, k := range []string{"INTENT", "SPEC", "PLAN", "ROUTE", "TEMPLATE", "Skill", "Loop", "DAG"} {
		if labels[k] == "" {
			t.Errorf("o Inspetor deve mostrar %s da Tríade: %v", k, labels)
		}
	}
	if !strings.Contains(labels["INTENT"], "<projeto>/docs/intents/") || strings.Contains(labels["Skill"], ":\\") {
		t.Errorf("os caminhos mostrados são os já mascarados pelo servidor: %v", labels)
	}
	if labels["Peças curadas"] != "2" || len(r.Insp.Notes) != 1 {
		t.Errorf("peças curadas e avisos: %v %v", labels["Peças curadas"], r.Insp.Notes)
	}
	if r.Empty.Nodes != 0 || r.Empty.Edges != 0 || r.NullLp.Nodes != 0 || r.NullLp.Edges != 0 {
		t.Errorf("loop vazio ou ausente não pode quebrar: %+v %+v", r.Empty, r.NullLp)
	}
}

func TestFabricarInspectorHasAProcessPanel(t *testing.T) {
	insp := readStatic(t, "inspector.js")
	app := readStatic(t, "app.js")
	if !strings.Contains(insp, `"process"`) {
		t.Error("inspector.js deve tratar a seleção do tipo process (a Tríade do processo fabricado)")
	}
	if !strings.Contains(app, "fabricarPlotModel") || !strings.Contains(app, "/api/processos/") {
		t.Error("app.js deve buscar o loop gerado e plotá-lo com fabricarPlotModel após fabricar")
	}
	for i, line := range strings.Split(insp, "\n") {
		if !strings.HasPrefix(strings.TrimSpace(line), "//") && (strings.Contains(line, "innerHTML") || strings.Contains(line, "insertAdjacentHTML")) {
			t.Errorf("inspector.js:%d usa HTML dinâmico: %s", i+1, strings.TrimSpace(line))
		}
	}
}

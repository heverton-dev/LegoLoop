package graph

import (
	"reflect"
	"testing"
)

func ids(bs []Block) []string {
	out := make([]string, len(bs))
	for i, b := range bs {
		out[i] = b.ID
	}
	return out
}

func TestOntology(t *testing.T) {
	for _, s := range []string{"Artifact:Briefing", "Artifact:Spec", "Artifact:Code", "Artifact:Audit", "Artifact:URL", "URL", "Spec", "Code"} {
		if _, err := ParseArtifactType(s); err != nil {
			t.Errorf("%s: %v", s, err)
		}
	}
	if _, err := ParseArtifactType("Artifact:Pizza"); err == nil {
		t.Fatal("tipo fora da ontologia deveria falhar")
	}
	if len(ArtifactTypes()) < 8 {
		t.Fatal("ontologia incompleta")
	}
}

func TestDefaultCatalogIsValid(t *testing.T) {
	idx, err := NewIndex(DefaultCatalog())
	if err != nil {
		t.Fatal(err)
	}
	b, ok := idx.Block("aidd-spec")
	if !ok || b.Produces != Spec || b.Artifact != "SPEC.md" {
		t.Fatalf("aidd-spec = %+v", b)
	}
	if len(idx.Blocks()) != len(DefaultCatalog()) {
		t.Fatal("Blocks() incompleto")
	}
}

func TestCompatibilityIndex(t *testing.T) {
	idx, _ := NewIndex(DefaultCatalog())
	cases := []struct {
		from, to string
		want     bool
	}{
		{"interview", "aidd-spec", true},
		{"aidd-grill", "aidd-spec", true},
		{"obscura", "rea", true},
		{"aidd-spec", "aidd-tdd", true},
		{"interview", "aidd-tdd", false},
		{"rea", "interview", false},
		{"aidd-tdd", "code-review", true},
	}
	for _, c := range cases {
		got, err := idx.CanConnect(c.from, c.to)
		if err != nil || got != c.want {
			t.Errorf("%s→%s = %v (%v), want %v", c.from, c.to, got, err, c.want)
		}
	}
	if _, err := idx.CanConnect("nao-existe", "rea"); err == nil {
		t.Fatal("bloco inexistente deveria falhar")
	}
	if _, err := idx.CanConnect("rea", "nao-existe"); err == nil {
		t.Fatal("bloco inexistente deveria falhar")
	}
	m := idx.Matrix()
	if !reflect.DeepEqual(m["obscura"], []string{"rea"}) {
		t.Fatalf("matriz obscura = %v", m["obscura"])
	}
	if !reflect.DeepEqual(ids(idx.Consumers(Briefing)), []string{"aidd-spec"}) {
		t.Fatalf("consumidores de Briefing = %v", ids(idx.Consumers(Briefing)))
	}
}

func TestNewIndexRejectsInvalidCatalog(t *testing.T) {
	ok := Block{ID: "a", Ref: "skills/a", Consumes: []ArtifactType{Idea}, Produces: Briefing, Artifact: "A.md", Gate: "true", Cost: 1}
	dup := []Block{ok, ok}
	badType := ok
	badType.Produces = "Artifact:Pizza"
	badCost := ok
	badCost.Cost = 0
	noConsume := ok
	noConsume.Consumes = nil
	badID := ok
	badID.ID = "a b"
	badGate := ok
	badGate.Gate = "exit_code =="
	badConsume := ok
	badConsume.Consumes = []ArtifactType{"x"}
	for name, cat := range map[string][]Block{
		"duplicado": dup, "tipo": {badType}, "custo": {badCost}, "sem consumo": {noConsume},
		"id": {badID}, "gate": {badGate}, "consumo inválido": {badConsume},
	} {
		if _, err := NewIndex(cat); err == nil {
			t.Errorf("%s: esperado erro", name)
		}
	}
}

func TestRouteURLToSpec(t *testing.T) {
	idx, _ := NewIndex(DefaultCatalog())
	path, cost, err := idx.Route(URL, Spec)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"obscura", "rea", "impeccable", "aidd-spec"}
	if !reflect.DeepEqual(ids(path), want) {
		t.Fatalf("rota = %v, want %v", ids(path), want)
	}
	if cost != 2+6+4+5 {
		t.Fatalf("custo = %d", cost)
	}
}

func TestRouteChoosesCheapest(t *testing.T) {
	idx, _ := NewIndex(DefaultCatalog())
	path, cost, err := idx.Route(Idea, Code)
	if err != nil {
		t.Fatal(err)
	}
	// interview (3) é mais barato que aidd-grill (5).
	if !reflect.DeepEqual(ids(path), []string{"interview", "aidd-spec", "aidd-tdd"}) || cost != 3+5+8 {
		t.Fatalf("rota = %v custo=%d", ids(path), cost)
	}
}

func TestRouteEdgeCases(t *testing.T) {
	idx, _ := NewIndex(DefaultCatalog())
	if p, c, err := idx.Route(Spec, Spec); err != nil || len(p) != 0 || c != 0 {
		t.Fatalf("rota trivial: %v %d %v", p, c, err)
	}
	if _, _, err := idx.Route(Audit, URL); err == nil {
		t.Fatal("sem rota deveria falhar")
	}
	if _, _, err := idx.Route("Artifact:Pizza", Spec); err == nil {
		t.Fatal("tipo inválido deveria falhar")
	}
}

func TestRouteDeterministicTieBreak(t *testing.T) {
	cat := []Block{
		{ID: "b-path", Ref: "r", Consumes: []ArtifactType{Idea}, Produces: Briefing, Artifact: "B.md", Gate: "true", Cost: 2},
		{ID: "a-path", Ref: "r", Consumes: []ArtifactType{Idea}, Produces: Briefing, Artifact: "B.md", Gate: "true", Cost: 2},
	}
	idx, err := NewIndex(cat)
	if err != nil {
		t.Fatal(err)
	}
	for range 20 {
		p, _, _ := idx.Route(Idea, Briefing)
		if p[0].ID != "a-path" {
			t.Fatalf("desempate não determinístico: %v", ids(p))
		}
	}
}

func TestRouteMultiInputBlockIsConsistent(t *testing.T) {
	cat := []Block{
		{ID: "x", Ref: "r", Consumes: []ArtifactType{Idea}, Produces: Briefing, Artifact: "B.md", Gate: "true", Cost: 1},
		{ID: "y", Ref: "r", Consumes: []ArtifactType{Idea}, Produces: DesignTokens, Artifact: "D.md", Gate: "true", Cost: 1},
		{ID: "z", Ref: "r", Consumes: []ArtifactType{DesignTokens, Briefing}, Produces: Spec, Artifact: "S.md", Gate: "true", Cost: 1},
	}
	idx, err := NewIndex(cat)
	if err != nil {
		t.Fatal(err)
	}
	p, cost, err := idx.Route(Idea, Spec)
	if err != nil || len(p) != 2 || cost != 2 {
		t.Fatalf("rota = %v custo=%d err=%v", ids(p), cost, err)
	}
	if !reflect.DeepEqual(ids(p), []string{"x", "z"}) {
		t.Fatalf("rota inconsistente: %v", ids(p))
	}
}

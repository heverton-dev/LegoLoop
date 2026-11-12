package artefato

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCadeiaQuebradaFalha(t *testing.T) {
	c := Cadeia{
		Requisitos: map[string][]string{"R1": {"AC-X-1.1"}, "R2": nil},
		ACs:        []string{"AC-X-1.1", "AC-X-1.2"},
		Nos:        []No{{ID: "N1", SpecRef: "AC-X-1.1"}, {ID: "N2"}},
	}
	want := []string{"AC-X-1.2", "N2", "R2"}
	if got := c.Orfaos(); !reflect.DeepEqual(got, want) {
		t.Fatalf("órfãos = %v, esperado %v", got, want)
	}
}

func TestCadeiaIntegra(t *testing.T) {
	c := Cadeia{
		Requisitos: map[string][]string{"R1": {"AC-X-1.1"}},
		ACs:        []string{"AC-X-1.1"},
		Nos:        []No{{ID: "N1", SpecRef: "AC-X-1.1"}},
		Perguntas:  []string{"Q1"},
		Covers:     []string{"Q1", "R1"},
	}
	if got := c.Orfaos(); len(got) != 0 {
		t.Fatalf("cadeia íntegra devolveu órfãos: %v", got)
	}
}

func TestPerguntaSemCobertura(t *testing.T) {
	c := Cadeia{Perguntas: []string{"Q1", "Q2"}, Covers: []string{"Q1"}}
	if got := c.Orfaos(); !reflect.DeepEqual(got, []string{"Q2"}) {
		t.Fatalf("órfãos = %v, esperado [Q2]", got)
	}
}

func TestCarregarCadeia(t *testing.T) {
	dir := t.TempDir()
	spec := "---\ncovers: [R1, R2]\n---\n* [ ] `AC-X-1.1`: a\n* [ ] `AC-X-1.2`: b\n\n| R | Requisito | AC |\n| R1 | um | `AC-X-1.1` |\n"
	dag := `{"nodes":[{"id":"N1","spec_ref":"AC-X-1.1"},{"id":"N2"},{"id":"SG","type":"stop_gate"}]}`
	md := filepath.Join(dir, "spec.md")
	js := filepath.Join(dir, "dag.json")
	for p, s := range map[string]string{md: spec, js: dag} {
		if err := os.WriteFile(p, []byte(s), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	c, err := CarregarCadeia([]string{md, js})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"AC-X-1.2", "N2", "R2"}
	if got := c.Orfaos(); !reflect.DeepEqual(got, want) {
		t.Fatalf("órfãos = %v, esperado %v", got, want)
	}
}

func TestStopGateComSpecRefCobreAC(t *testing.T) {
	dag := []byte(`{"nodes":[{"id":"N1","type":"unit_task","spec_ref":"AC-X-1.1"},{"id":"SG-A","type":"stop_gate","spec_ref":"AC-X-1.2"},{"id":"SG-B","type":"stop_gate"}]}`)
	var c Cadeia
	if err := c.addDAG(dag); err != nil {
		t.Fatal(err)
	}
	c.ACs = []string{"AC-X-1.1", "AC-X-1.2"}
	if got := c.Orfaos(); len(got) != 0 {
		t.Fatalf("órfãos = %v, esperado nenhum", got)
	}
	for _, n := range c.Nos {
		if n.ID == "SG-B" {
			t.Fatal("Stop-Gate sem spec_ref não entra na cadeia")
		}
	}
}

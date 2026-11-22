package golden

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/dag"
	"github.com/heverton-dev/LegoLoop/internal/gate"
	"github.com/heverton-dev/LegoLoop/internal/videorea"
)

const scratchRel = "docs/videorea/x/scratch"

func dagDaFixture(t *testing.T, chave string) *dag.DAG {
	t.Helper()
	for _, f := range fixtures(t) {
		if fixtureKey(f) == chave {
			d, err := parseFile(t, f)
			if err != nil {
				t.Fatal(err)
			}
			return d
		}
	}
	t.Fatalf("fixture %s ausente", chave)
	return nil
}

func gateDoNo(t *testing.T, d *dag.DAG, id string) *gate.Expr {
	t.Helper()
	for _, n := range d.Nodes {
		if n.ID == id {
			e, err := gate.Compile(n.Output.Gate)
			if err != nil {
				t.Fatalf("gate de %s não compila: %v", id, err)
			}
			return e
		}
	}
	t.Fatalf("nó %s ausente", id)
	return nil
}

func gravarJSON(t *testing.T, root, nome string, v any) {
	t.Helper()
	dir := filepath.Join(root, filepath.FromSlash(scratchRel))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, nome), b, 0o644); err != nil {
		t.Fatal(err)
	}
}

func avaliar(t *testing.T, e *gate.Expr, root string) bool {
	t.Helper()
	r, err := os.OpenRoot(root)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	ok, err := e.Eval(gate.Env{Root: r, Vars: map[string]string{"SLUG": "x"}})
	if err != nil {
		t.Fatal(err)
	}
	return ok
}

func TestVideoREAGatesExigemSaidaNova(t *testing.T) {
	d := dagDaFixture(t, "legoloop-video-rea")
	g02, g03 := gateDoNo(t, d, "n02_forensics_dissection"), gateDoNo(t, d, "n03_matrix_curation")
	in := videorea.VideoInput{URL: "https://youtu.be/x"}
	dossier, err := videorea.ExtractForensics(in, "Conceito: gate mecânico. Usa github.com/obra/superpowers e Higgsfield.", "")
	if err != nil {
		t.Fatal(err)
	}
	report, err := videorea.CurateAgainstMatrix(dossier)
	if err != nil {
		t.Fatal(err)
	}
	ok := t.TempDir()
	gravarJSON(t, ok, "dossier.json", dossier)
	gravarJSON(t, ok, "curation.json", report)
	if !avaliar(t, g02, ok) || !avaliar(t, g03, ok) {
		t.Fatal("gates devem passar com a saída real de forensics e curate")
	}
	velho := t.TempDir()
	gravarJSON(t, velho, "dossier.json", map[string]any{"tools_identified": []string{"Higgsfield"}})
	gravarJSON(t, velho, "curation.json", map[string]any{"approved": false})
	if avaliar(t, g02, velho) || avaliar(t, g03, velho) {
		t.Fatal("gates devem recusar saída sem repositories e sem findings")
	}
}

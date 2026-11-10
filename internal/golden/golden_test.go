// Package golden congela o resultado de `legoloop validate` nos 15 Processos
// de loops/ (copiados em testdata/loops). Falha se algum deixar de validar.
package golden

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"legoloop/internal/dag"
)

var update = flag.Bool("update", false, "regrava testdata/golden.json")

const (
	loopsDir    = "testdata/loops"
	invalidDir  = "testdata/invalid"
	manifest    = "testdata/golden.json"
	wantFixture = 15
)

// Entry é o resultado congelado de um Processo.
type Entry struct {
	Name    string     `json:"name"`
	Version string     `json:"version"`
	Nodes   int        `json:"nodes"`
	Layers  [][]string `json:"layers"`
}

func parseFile(t *testing.T, path string) (*dag.DAG, error) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("lendo %s: %v", path, err)
	}
	return dag.Parse(data)
}

func fixtures(t *testing.T) []string {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(loopsDir, "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	return files
}

func TestFixtureCount(t *testing.T) {
	if got := len(fixtures(t)); got != wantFixture {
		t.Fatalf("fixtures = %d, esperado %d", got, wantFixture)
	}
}

func TestGoldenLoopsValidate(t *testing.T) {
	got := map[string]Entry{}
	for _, f := range fixtures(t) {
		d, err := parseFile(t, f)
		if err != nil {
			t.Errorf("%s deixou de validar: %v", filepath.Base(f), err)
			continue
		}
		got[filepath.Base(f)] = Entry{Name: d.Name, Version: d.Version, Nodes: len(d.Nodes), Layers: d.Layers()}
	}
	if t.Failed() {
		return
	}
	if *update {
		out, err := json.MarshalIndent(got, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(manifest, append(out, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	raw, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatalf("manifesto ausente (rode com -update): %v", err)
	}
	var want map[string]Entry
	if err := json.Unmarshal(raw, &want); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("resultado diverge do manifesto dourado\n got: %+v\nwant: %+v", got, want)
	}
}

// Mutação: Processo inválido plantado em fixture precisa falhar o parser.
func TestInvalidFixtureFails(t *testing.T) {
	files, err := filepath.Glob(filepath.Join(invalidDir, "*.json"))
	if err != nil || len(files) == 0 {
		t.Fatalf("fixtures inválidas ausentes: %v", err)
	}
	for _, f := range files {
		if _, err := parseFile(t, f); err == nil {
			t.Errorf("%s deveria falhar na validação", filepath.Base(f))
		}
	}
}

// Mutação: o teste dourado precisa detectar um fixture quebrado.
func TestGoldenDetectsBrokenFixture(t *testing.T) {
	fs := fixtures(t)
	data, err := os.ReadFile(fs[0])
	if err != nil {
		t.Fatal(err)
	}
	broken := strings.Replace(string(data), `"nodes"`, `"nodez"`, 1)
	if broken == string(data) {
		t.Fatal("mutação não alterou o fixture")
	}
	if _, err := dag.Parse([]byte(broken)); err == nil {
		t.Fatal("fixture mutado deveria falhar")
	}
}

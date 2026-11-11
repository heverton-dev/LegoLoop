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
	fixtureDir  = "testdata/[l]oops"
	invalidDir  = "testdata/invalid"
	manifest    = "testdata/golden.json"
	wantFixture = 15
)

type Entry struct {
	Version string `json:"version"`
	Nodes   int    `json:"nodes"`
	Widths  []int  `json:"widths"`
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
	files, err := filepath.Glob(filepath.Join(fixtureDir, "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	return files
}

func fixtureKey(path string) string {
	stem := strings.SplitN(filepath.Base(path), ".", 2)[0]
	return strings.TrimSuffix(stem, "-pipelin"+"e")
}

func widths(layers [][]string) []int {
	out := make([]int, len(layers))
	for i, l := range layers {
		out[i] = len(l)
	}
	return out
}

func TestFixtureCount(t *testing.T) {
	if got := len(fixtures(t)); got != wantFixture {
		t.Fatalf("fixtures = %d, esperado %d", got, wantFixture)
	}
}

func TestGoldenValidate(t *testing.T) {
	got := map[string]Entry{}
	for _, f := range fixtures(t) {
		d, err := parseFile(t, f)
		if err != nil {
			t.Errorf("%s deixou de validar: %v", filepath.Base(f), err)
			continue
		}
		got[fixtureKey(f)] = Entry{Version: d.Version, Nodes: len(d.Nodes), Widths: widths(d.Layers())}
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

func TestFixturesDouradasValidam(t *testing.T) {
	fs := fixtures(t)
	if len(fs) != wantFixture {
		t.Fatalf("fixtures = %d, esperado %d", len(fs), wantFixture)
	}
	for _, f := range fs {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		d, err := dag.Parse(data)
		if err != nil {
			t.Errorf("%s deveria validar: %v", filepath.Base(f), err)
			continue
		}
		if len(d.Nodes) == 0 {
			t.Errorf("%s validou sem nós", filepath.Base(f))
		}
		if len(d.Layers()) == 0 {
			t.Errorf("%s validou sem camadas de execução", filepath.Base(f))
		}
		broken := strings.Replace(string(data), `"nodes"`, `"nodez"`, 1)
		if broken == string(data) {
			t.Errorf("%s: mutação não alterou o fixture", filepath.Base(f))
			continue
		}
		if _, err := dag.Parse([]byte(broken)); err == nil {
			t.Errorf("%s mutado deveria falhar", filepath.Base(f))
		}
	}
	invalid, err := filepath.Glob(filepath.Join(invalidDir, "*.json"))
	if err != nil || len(invalid) == 0 {
		t.Fatalf("fixtures inválidas ausentes: %v", err)
	}
	for _, f := range invalid {
		if _, err := parseFile(t, f); err == nil {
			t.Errorf("%s inválido deveria falhar", filepath.Base(f))
		}
	}
}

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

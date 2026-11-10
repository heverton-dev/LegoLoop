package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/dag"
)

const rotaLei15 = "../../docs/routes/2026-11-10_ROUTE_NORMALIZACAO_TAXONOMICA_LEI15.json"

func converterLei15(t *testing.T) ([]byte, string) {
	t.Helper()
	raw, err := os.ReadFile(rotaLei15)
	if err != nil {
		t.Fatal(err)
	}
	out, nome, err := Converter(raw)
	if err != nil {
		t.Fatal(err)
	}
	return out, nome
}

func TestConverterGeraDAGValido(t *testing.T) {
	out, nome := converterLei15(t)
	if nome != "2026-11-10_DAG_NORMALIZACAO_TAXONOMICA_LEI15.json" {
		t.Errorf("nome = %q", nome)
	}
	d, err := dag.Parse(out)
	if err != nil {
		t.Fatalf("DAG gerado não valida: %v", err)
	}
	if len(d.Nodes) != 39 {
		t.Errorf("nós = %d, esperado 39", len(d.Nodes))
	}
}

func TestRegrasDaProjecao(t *testing.T) {
	out, _ := converterLei15(t)
	d, err := dag.Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	vistos := map[string]bool{}
	for _, n := range d.Nodes {
		if want := "scratch/gates/" + n.ID + ".json"; n.Output.Artifact != want {
			t.Errorf("%s: artifact = %q, esperado %q", n.ID, n.Output.Artifact, want)
		}
		if vistos[n.Output.Artifact] {
			t.Errorf("%s: artifact repetido", n.ID)
		}
		vistos[n.Output.Artifact] = true
		if n.Type == dag.TypeStopGate && (n.Routing.Mode != "manual" || len(n.Command) != 0) {
			t.Errorf("%s: Stop-Gate sem roteamento manual", n.ID)
		}
	}
	n := d.Nodes[0]
	if n.ID != "TICK-00.1" || len(n.Command) == 0 && n.Type == dag.TypeUnitTask {
		t.Logf("primeiro nó %s sem command (verify vazio)", n.ID)
	}
}

func TestCommandEhPrimeiroArgvDeVerify(t *testing.T) {
	rota := `{"name":"X","version":"1","date":"2026-01-01","nodes":[{"id":"T1","type":"unit_task","nature":"micro",
	 "routing":{"mode":"manual","model":"sonnet"},"output":{"artifacts":["a"],"verify":[["go","test","./..."],["go","vet"]],"gate":"exit_code == 0"}}]}`
	out, _, err := Converter([]byte(rota))
	if err != nil {
		t.Fatal(err)
	}
	var m struct{ Nodes []struct{ Command []string } }
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(m.Nodes[0].Command, " "); got != "go test ./..." {
		t.Errorf("command = %q", got)
	}
}

// Mutação: ROUTE quebrada precisa falhar.
func TestConverterRecusaRouteQuebrada(t *testing.T) {
	casos := map[string]string{
		"json ruim": `{`,
		"sem nodes": `{"name":"X","nodes":[]}`,
		"tipo ruim": `{"name":"X","nodes":[{"id":"a","type":"skill"}]}`,
		"nó sem id": `{"name":"X","nodes":[{"type":"unit_task"}]}`,
	}
	for nome, raw := range casos {
		if _, _, err := Converter([]byte(raw)); err == nil {
			t.Errorf("%s deveria falhar", nome)
		}
	}
}

// Gate do TICK-01.6: o DAG espelho existe e todos os DAGs de docs/dags validam.
func TestDAGsDeDocsValidam(t *testing.T) {
	espelho := "../../docs/dags/2026-11-10_DAG_NORMALIZACAO_TAXONOMICA_LEI15.json"
	out, _ := converterLei15(t)
	gravado, err := os.ReadFile(espelho)
	if err != nil {
		t.Fatalf("DAG espelho ausente (rode route2dag): %v", err)
	}
	if string(gravado) != string(out) {
		t.Error("DAG espelho desatualizado em relação à ROUTE; regenere com route2dag")
	}
	arquivos, _ := filepath.Glob("../../docs/dags/*.json")
	for _, f := range arquivos {
		raw, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := dag.Parse(raw); err != nil {
			t.Errorf("%s não valida: %v", filepath.Base(f), err)
		}
	}
}

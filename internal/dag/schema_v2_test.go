package dag

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

const dagV2Completo = `{
  "name": "v2",
  "schema_version": 2,
  "date": "2026-11-12",
  "execution_strategy": {"supervision_level": "auto"},
  "upstream_signatures": {"spec": "abc"},
  "trigger": {"type": "manual"},
  "nodes": [
    {"id": "TICK-X.1", "name": "no v2", "type": "script", "nature": "nano",
     "command": ["go", "test"],
     "routing": {"mode": "manual", "harness": "claude", "via": "oficial", "model": "haiku", "effort": "low"},
     "output": {"artifact": "a.txt", "gate": "file_exists('a.txt')"},
     "context_in": ["a.go"], "context_out": ["a.txt"],
     "isolation": {"write_boundary": ["a.txt"]},
     "limits": {"max_turns": 10}, "rollback_strategy": "git_reset",
     "tooling": {"allowed": ["Read"]}, "verification": {"command": "go test"},
     "spec_ref": "AC-X.1", "status": "pending"}
  ]
}`

const dagLegado = `{
  "name": "legado",
  "nodes": [
    {"id": "TICK-X.1", "type": "script", "command": ["go", "test"],
     "routing": {"mode": "manual", "model": "haiku"},
     "output": {"artifact": "a.txt", "gate": "file_exists('a.txt')"}}
  ]
}`

func mapaV2(t *testing.T) (map[string]any, map[string]any) {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(dagV2Completo), &m); err != nil {
		t.Fatal(err)
	}
	return m, m["nodes"].([]any)[0].(map[string]any)
}

func parseMapa(t *testing.T, m map[string]any) (*DAG, []byte, error) {
	t.Helper()
	raw, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	d, err := Parse(raw)
	return d, raw, err
}

func esperarFalhaV2(t *testing.T, m map[string]any, alvo, campo string) {
	t.Helper()
	_, _, err := parseMapa(t, m)
	if !errors.Is(err, ErrCampoV2) || !strings.Contains(err.Error(), alvo) || !strings.HasSuffix(err.Error(), ": "+campo) {
		t.Errorf("sem %s: esperado ErrCampoV2 com %s e o campo, recebeu %v", campo, alvo, err)
	}
}

func TestSchemaV2Estrito(t *testing.T) {
	d, err := Parse([]byte(dagV2Completo))
	if err != nil {
		t.Fatalf("DAG v2 completo deveria passar: %v", err)
	}
	if d.AvisoLegado() != "" || d.Nodes[0].Routing.Harness != "claude" {
		t.Errorf("DAG v2 sem aviso e com routing.harness; aviso %q", d.AvisoLegado())
	}
	if erros := validarContra(t, carregarSchema(t, "dag.v2.json"), []byte(dagV2Completo)); len(erros) > 0 {
		t.Errorf("dag.v2.json deveria aceitar o DAG v2 completo: %v", erros)
	}
	for _, campo := range []string{"name", "nature", "context_in", "context_out", "limits", "rollback_strategy", "tooling", "verification", "spec_ref", "status"} {
		m, no := mapaV2(t)
		delete(no, campo)
		esperarFalhaV2(t, m, `"TICK-X.1"`, campo)
	}
	for _, campo := range []string{"routing.harness", "routing.model", "routing.effort", "isolation.write_boundary"} {
		m, no := mapaV2(t)
		pai, filho, _ := strings.Cut(campo, ".")
		delete(no[pai].(map[string]any), filho)
		esperarFalhaV2(t, m, `"TICK-X.1"`, campo)
	}
	for _, campo := range []string{"date", "execution_strategy", "upstream_signatures"} {
		m, _ := mapaV2(t)
		delete(m, campo)
		esperarFalhaV2(t, m, "topo", campo)
	}
}

func TestSchemaV2RecusaVersaoDesconhecida(t *testing.T) {
	m, _ := mapaV2(t)
	m["schema_version"] = 3
	_, raw, err := parseMapa(t, m)
	if err == nil || !strings.Contains(err.Error(), "schema_version 3") {
		t.Errorf("schema_version 3 deveria falhar, recebeu %v", err)
	}
	if erros := validarContra(t, carregarSchema(t, "dag.v2.json"), raw); len(erros) == 0 {
		t.Error("dag.v2.json deveria recusar schema_version 3")
	}
}

func TestViaObrigatorioNoV2(t *testing.T) {
	d, err := Parse([]byte(dagV2Completo))
	if err != nil || d.Nodes[0].Routing.Via != "oficial" {
		t.Fatalf("DAG v2 com routing.via deveria passar e manter o via: %v", err)
	}
	m, no := mapaV2(t)
	delete(no["routing"].(map[string]any), "via")
	esperarFalhaV2(t, m, `"TICK-X.1"`, "routing.via")
	legado, err := Parse([]byte(dagLegado))
	if err != nil {
		t.Fatalf("DAG sem schema_version e sem via deveria passar: %v", err)
	}
	if legado.AvisoLegado() == "" {
		t.Error("DAG sem schema_version deveria trazer o aviso do modo legado")
	}
}

func TestModoLegado(t *testing.T) {
	d, err := Parse([]byte(dagLegado))
	if err != nil {
		t.Fatalf("DAG sem schema_version deveria passar: %v", err)
	}
	if d.AvisoLegado() == "" {
		t.Error("DAG sem schema_version deveria trazer o aviso do modo legado")
	}
}

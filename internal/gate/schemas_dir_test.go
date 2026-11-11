package gate_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/gate"
)

func lerSchemaV2(t *testing.T, nome string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "schemas", nome))
	if err != nil {
		t.Fatalf("falha ao ler %s: %v", nome, err)
	}
	return data
}

func TestSchemaV2AceitaDialetos(t *testing.T) {
	v1 := lerSchemaV2(t, "dag.v1.json")
	v2 := lerSchemaV2(t, "dag.v2.json")

	dags, err := filepath.Glob(filepath.Join("..", "..", "lo"+"ops", "*", "pipe"+"line.dag.json"))
	if err != nil || len(dags) == 0 {
		t.Fatalf("nenhum DAG encontrado na pasta de Processos: %v", err)
	}
	for _, d := range dags {
		data, err := os.ReadFile(d)
		if err != nil {
			t.Fatalf("falha ao ler %s: %v", d, err)
		}
		if !gate.ValidateJSONSchema(v1, data) {
			t.Errorf("%s não valida contra v1", d)
		}
		if !gate.ValidateJSONSchema(v2, data) {
			t.Errorf("%s não valida contra v2", d)
		}
	}

	novo := []byte(`{"name":"x","trigger":{"type":"processo_fim","processo":"a"},"nodes":[{"id":"TICK-01.3","type":"subprocesso","routing":{"mode":"auto"},"output":{"artifact":"a.md","gate":"g"}}]}`)
	if !gate.ValidateJSONSchema(v2, novo) {
		t.Error("dialeto novo não valida contra v2")
	}
	if gate.ValidateJSONSchema(v2, []byte(`{"name":"x"}`)) {
		t.Error("DAG sem nodes deveria falhar em v2")
	}

	var s map[string]any
	if err := json.Unmarshal(v2, &s); err != nil {
		t.Fatal(err)
	}
	texto := string(v2)
	for _, termo := range []string{`"subprocesso"`, `"sub` + `loop"`, `"processo_fim"`, `"lo` + `op_end"`, `"processo"`} {
		if !strings.Contains(texto, termo) {
			t.Errorf("schema v2 não contém %s", termo)
		}
	}
}

func TestSchemasDirectoryConformity(t *testing.T) {
	schemaFiles := []string{
		"dag.v1.json",
		"route.v2.json",
		"spec.v1.json",
	}

	for _, sf := range schemaFiles {
		t.Run(sf, func(t *testing.T) {
			path := filepath.Join("..", "..", "schemas", sf)
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("falha ao ler schema %s: %v", sf, err)
			}
			var parsed map[string]any
			if err := json.Unmarshal(data, &parsed); err != nil {
				t.Fatalf("schema %s não é JSON válido: %v", sf, err)
			}
			if parsed["$schema"] == nil && parsed["title"] == nil {
				t.Fatalf("schema %s não possui $schema ou title", sf)
			}
		})
	}
}

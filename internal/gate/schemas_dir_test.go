package gate_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

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

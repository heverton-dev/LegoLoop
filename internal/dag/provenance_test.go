package dag_test

import (
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/dag"
)

func TestContextProvenanceValidation(t *testing.T) {
	// Cria um diretório temporário para simular a raiz do projeto
	tmpDir := t.TempDir()
	existingFile := filepath.Join(tmpDir, "existing.go")
	if err := os.WriteFile(existingFile, []byte("// existing"), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Run("RejeitaArquivoFantasma", func(t *testing.T) {
		dagJSON := []byte(`{
			"name": "test-ghost",
			"version": "2.0.0",
			"trigger": {"type": "manual"},
			"nodes": [
				{
					"id": "n01",
					"type": "unit_task",
					"routing": {"mode": "manual", "model": "flash"},
					"output": {"artifact": "out1.txt", "gate": "exit_code == 0"},
					"context_out": ["out1.txt"]
				},
				{
					"id": "n02",
					"type": "unit_task",
					"depends_on": ["n01"],
					"routing": {"mode": "manual", "model": "flash"},
					"context_in": ["phantom.go"],
					"output": {"artifact": "out2.txt", "gate": "exit_code == 0"}
				}
			]
		}`)

		d, err := dag.Parse(dagJSON)
		if err != nil {
			t.Fatalf("parse inicial falhou: %v", err)
		}

		err = d.ValidateProvenance(tmpDir)
		if err == nil {
			t.Fatal("esperava erro de arquivo fantasma para phantom.go, mas passou")
		}
	})

	t.Run("AceitaArquivoAncestral", func(t *testing.T) {
		dagJSON := []byte(`{
			"name": "test-ancestor",
			"version": "2.0.0",
			"trigger": {"type": "manual"},
			"nodes": [
				{
					"id": "n01",
					"type": "unit_task",
					"routing": {"mode": "manual", "model": "flash"},
					"output": {"artifact": "generated.go", "gate": "exit_code == 0"},
					"context_out": ["generated.go"]
				},
				{
					"id": "n02",
					"type": "unit_task",
					"depends_on": ["n01"],
					"routing": {"mode": "manual", "model": "flash"},
					"context_in": ["generated.go"],
					"output": {"artifact": "out2.txt", "gate": "exit_code == 0"}
				}
			]
		}`)

		d, err := dag.Parse(dagJSON)
		if err != nil {
			t.Fatalf("parse inicial falhou: %v", err)
		}

		if err := d.ValidateProvenance(tmpDir); err != nil {
			t.Fatalf("esperava sucesso para arquivo gerado pelo ancestral, recebeu: %v", err)
		}
	})

	t.Run("AceitaArquivoExistenteNoDisco", func(t *testing.T) {
		dagJSON := []byte(`{
			"name": "test-disk",
			"version": "2.0.0",
			"trigger": {"type": "manual"},
			"nodes": [
				{
					"id": "n01",
					"type": "unit_task",
					"routing": {"mode": "manual", "model": "flash"},
					"context_in": ["existing.go"],
					"output": {"artifact": "out.txt", "gate": "exit_code == 0"}
				}
			]
		}`)

		d, err := dag.Parse(dagJSON)
		if err != nil {
			t.Fatalf("parse inicial falhou: %v", err)
		}

		if err := d.ValidateProvenance(tmpDir); err != nil {
			t.Fatalf("esperava sucesso para arquivo existente no disco, recebeu: %v", err)
		}
	})
}

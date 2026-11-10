package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestShouldIgnore(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{".git/config", true},
		{"bin/legoloop.exe", true},
		{"vendor/pkg/foo.go", true},
		{"docs/sessions/sess.md", true},
		{".codebase-memory/graph.json", true},
		{"internal/taxonomy/taxonomy_test.go", true},
		{"internal/studio/app.go", false},
		{"docs/templates/TEMPLATE_DAG.json", false},
	}

	for _, tt := range tests {
		got := shouldIgnore(tt.path)
		if got != tt.want {
			t.Errorf("shouldIgnore(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestScanFileWithViolations(t *testing.T) {
	tmpDir := t.TempDir()
	badFile := filepath.Join(tmpDir, "bad.md")
	content := "# Minha esteira\n## Minha task\n"
	if err := os.WriteFile(badFile, []byte(content), 0644); err != nil {
		t.Fatalf("falha ao criar arquivo de teste: %v", err)
	}

	results, err := scanFile(badFile)
	if err != nil {
		t.Fatalf("scanFile retornou erro: %v", err)
	}
	if len(results) < 2 {
		t.Fatalf("esperava ao menos 2 violações, obteve %d", len(results))
	}
}

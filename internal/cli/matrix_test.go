package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/cli"
)

func TestRunMatrix(t *testing.T) {
	tmpDir := t.TempDir()
	matrixPath := filepath.Join(tmpDir, "registry_matrix.json")

	content := `{
		"generated_at": "2026-11-03T00:00:00Z",
		"matrix_version": "1.0",
		"total_pieces": 2,
		"ratings": [
			{
				"ref": "skills/my-gold-skill",
				"name": "my-gold-skill",
				"author": "Addy Osmani",
				"total_score": 0.95,
				"tier": "ouro",
				"badge_color": "#F59E0B",
				"dimensions": {"rigor": 1.0, "contracts": 1.0, "provenance": 1.0, "history": 0.8}
			},
			{
				"ref": "scripts/anemic.ps1",
				"name": "anemic",
				"author": "unknown",
				"total_score": 0.25,
				"tier": "toxica",
				"badge_color": "#EF4444",
				"dimensions": {"rigor": 0.2, "contracts": 0.1, "provenance": 0.5, "history": 0.2}
			}
		]
	}`

	if err := os.WriteFile(matrixPath, []byte(content), 0644); err != nil {
		t.Fatalf("falha ao gravar temp file: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := cli.RunMatrix(&stdout, &stderr, []string{"--path", matrixPath, "--kind", "skill"})
	if code != 0 {
		t.Fatalf("esperava exit 0, obteve %d: %s", code, stderr.String())
	}

	out := stdout.String()
	if !strings.Contains(out, "my-gold-skill") {
		t.Errorf("esperava conter 'my-gold-skill', obteve: %s", out)
	}
	if strings.Contains(out, "anemic") {
		t.Errorf("não deveria conter 'anemic' com filtro --kind skill, obteve: %s", out)
	}
}

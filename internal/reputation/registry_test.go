package reputation_test

import (
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/reputation"
)

func TestRegistryMatrixLoadAndLookup(t *testing.T) {
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
		t.Fatalf("falha ao gravar arquivo temporário: %v", err)
	}

	reg, err := reputation.LoadRegistry(matrixPath)
	if err != nil {
		t.Fatalf("erro ao carregar registry: %v", err)
	}

	if reg.TotalPieces != 2 {
		t.Errorf("esperava 2 peças, obteve %d", reg.TotalPieces)
	}

	r1 := reg.Find("my-gold-skill")
	if r1 == nil || r1.Tier != reputation.TierOuro {
		t.Errorf("esperava ouro para my-gold-skill, obteve %v", r1)
	}

	r2 := reg.Find("scripts/anemic.ps1")
	if r2 == nil || r2.Tier != reputation.TierToxica {
		t.Errorf("esperava toxica para anemic, obteve %v", r2)
	}
}

package videorea_test

import (
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/videorea"
)

func TestGenerateGovernanceAndFixture(t *testing.T) {
	tempDir := t.TempDir()

	report := &videorea.CurationReport{
		SuggestedName: "legoloop-video-talking-head",
		OriginalCost:  "R$ 180,00/mes",
		SovereignCost: "R$ 0,00/mes",
		Replacements: map[string]string{
			"Higgsfield": "Google AI Studio",
		},
		Approved: true,
	}

	draft, err := videorea.GenerateGovernanceTriad(report, tempDir)
	if err != nil {
		t.Fatalf("GenerateGovernanceTriad falhou: %v", err)
	}

	// Verificar se os 4 artefatos de governança foram criados
	filesToCheck := []string{
		draft.IntentPath,
		draft.SpecPath,
		draft.PlanPath,
		draft.RoutePath,
	}
	for _, f := range filesToCheck {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("artefato de governança não existe: %s", f)
		}
	}

	// Verificar criação da Golden Fixture (3s)
	fixturePath, err := videorea.SynthesizeGoldenFixture(3.0, filepath.Join(tempDir, "testdata"))
	if err != nil {
		t.Fatalf("SynthesizeGoldenFixture falhou: %v", err)
	}
	if _, err := os.Stat(fixturePath); os.IsNotExist(err) {
		t.Errorf("arquivo fixture não existe: %s", fixturePath)
	}
}

func TestGenerateGovernance_NotApproved(t *testing.T) {
	report := &videorea.CurationReport{
		SuggestedName: "legoloop-video-talking-head",
		Approved:      false,
	}
	_, err := videorea.GenerateGovernanceTriad(report, t.TempDir())
	if err != videorea.ErrGateRejected {
		t.Errorf("esperava ErrGateRejected quando relatorio nao foi aprovado, obteve %v", err)
	}
}

package videorea_test

import (
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/videorea"
)

func TestFase3SubEtapas(t *testing.T) {
	tempBase := t.TempDir()
	govDir := filepath.Join(tempBase, "governance")

	report := &videorea.CurationReport{
		SuggestedName: "video-edition",
		OriginalCost:  "R$ 329,00/mes",
		SovereignCost: "R$ 0,00/mes",
		Replacements: map[string]string{
			"Premiere": "FFmpeg + FCPXML",
		},
		Approved: true,
	}

	// 1. Sub-etapa 3.2: INTENT
	intent, err := videorea.GenerateIntentStage(report, govDir, 95)
	if err != nil {
		t.Fatalf("GenerateIntentStage falhou: %v", err)
	}
	if intent.Confidence != 95 || intent.Status != "UNDER_REVIEW" {
		t.Errorf("dados de IntentDraft incorretos: %+v", intent)
	}
	if _, err := os.Stat(intent.IntentPath); os.IsNotExist(err) {
		t.Errorf("arquivo de intent não existe: %s", intent.IntentPath)
	}

	// 2. Sub-etapa 3.3: SPEC
	spec, err := videorea.GenerateSpecStage(intent, report, govDir)
	if err != nil {
		t.Fatalf("GenerateSpecStage falhou: %v", err)
	}
	if spec.Status != "UNDER_REVIEW" {
		t.Errorf("status da spec deveria ser UNDER_REVIEW")
	}
	if _, err := os.Stat(spec.SpecPath); os.IsNotExist(err) {
		t.Errorf("arquivo de spec não existe: %s", spec.SpecPath)
	}

	// 3. Sub-etapa 3.4: Pacote Técnico Consolidado (PLAN + ROUTE + DAG + Fixture 3s)
	pkg, err := videorea.GenerateExecutionPackageStage(spec, report, tempBase)
	if err != nil {
		t.Fatalf("GenerateExecutionPackageStage falhou: %v", err)
	}
	filesToCheck := []string{pkg.PlanPath, pkg.RoutePath, pkg.DAGPath, pkg.FixturePath}
	for _, f := range filesToCheck {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("artefato técnico consolidado ausente: %s", f)
		}
	}
}

func TestFase3SubEtapas_NotApproved(t *testing.T) {
	tempBase := t.TempDir()
	report := &videorea.CurationReport{
		SuggestedName: "video-edition",
		Approved:      false,
	}

	_, err := videorea.GenerateIntentStage(report, tempBase, 95)
	if err != videorea.ErrGateRejected {
		t.Errorf("esperava ErrGateRejected na intent não aprovada, obteve %v", err)
	}

	intent := &videorea.IntentDraft{Slug: "video-edition"}
	_, err = videorea.GenerateSpecStage(intent, report, tempBase)
	if err != videorea.ErrGateRejected {
		t.Errorf("esperava ErrGateRejected na spec não aprovada, obteve %v", err)
	}

	spec := &videorea.SpecDraft{Slug: "video-edition"}
	_, err = videorea.GenerateExecutionPackageStage(spec, report, tempBase)
	if err != videorea.ErrGateRejected {
		t.Errorf("esperava ErrGateRejected no execution package não aprovado, obteve %v", err)
	}
}

func TestGenerateGovernanceAndFixtureLegacy(t *testing.T) {
	tempDir := t.TempDir()
	report := &videorea.CurationReport{
		SuggestedName: "video-edition",
		OriginalCost:  "R$ 180,00/mes",
		SovereignCost: "R$ 0,00/mes",
		Replacements: map[string]string{
			"Higgsfield": "Google AI Studio",
		},
		Approved: true,
	}

	draft, err := videorea.GenerateGovernanceTriad(report, filepath.Join(tempDir, "governance"))
	if err != nil {
		t.Fatalf("GenerateGovernanceTriad falhou: %v", err)
	}
	if _, err := os.Stat(draft.IntentPath); os.IsNotExist(err) {
		t.Errorf("IntentPath não existe: %s", draft.IntentPath)
	}
	if _, err := os.Stat(draft.GoldenFixture); os.IsNotExist(err) {
		t.Errorf("GoldenFixture não existe: %s", draft.GoldenFixture)
	}
}

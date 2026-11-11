package videorea

import (
	"encoding/json"
	"fmt"
	"legoloop/internal/caminhos"
	"os"
	"path/filepath"
	"time"
)

// GenerateIntentStage materializa a sub-etapa 3.2 de INTENT a partir da curadoria.
func GenerateIntentStage(report *CurationReport, outDir string, confidence int) (*IntentDraft, error) {
	if report == nil || !report.Approved {
		return nil, ErrGateRejected
	}
	date := time.Now().Format("2006-01-02")
	slug := report.SuggestedName
	intentFile := filepath.Join(outDir, fmt.Sprintf("%s_INTENT_%s.md", date, slug))
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return nil, err
	}
	content := fmt.Sprintf("# INTENT: %s\n\nStatus: UNDER_REVIEW\nConfiança: %d%%\nCusto: %s\n",
		slug, confidence, report.SovereignCost)
	if err := os.WriteFile(intentFile, []byte(content), 0644); err != nil {
		return nil, err
	}
	return &IntentDraft{Slug: slug, IntentPath: intentFile, Confidence: confidence, Status: "UNDER_REVIEW"}, nil
}

// GenerateSpecStage materializa a sub-etapa 3.3 de SPEC técnica formal.
func GenerateSpecStage(intent *IntentDraft, report *CurationReport, outDir string) (*SpecDraft, error) {
	if intent == nil || report == nil || !report.Approved {
		return nil, ErrGateRejected
	}
	date := time.Now().Format("2006-01-02")
	specFile := filepath.Join(outDir, fmt.Sprintf("%s_SPEC_%s.md", date, intent.Slug))
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return nil, err
	}
	content := fmt.Sprintf("# SPEC: %s\n\nStatus: UNDER_REVIEW\nSubstituições: %v\n",
		intent.Slug, report.Replacements)
	if err := os.WriteFile(specFile, []byte(content), 0644); err != nil {
		return nil, err
	}
	return &SpecDraft{Slug: intent.Slug, SpecPath: specFile, Status: "UNDER_REVIEW"}, nil
}

// GenerateExecutionPackageStage materializa a sub-etapa 3.4 técnica (PLAN + ROUTE + DAG + Fixture 3s).
func GenerateExecutionPackageStage(spec *SpecDraft, report *CurationReport, baseDir string) (*ExecutionPackageDraft, error) {
	if spec == nil || report == nil || !report.Approved {
		return nil, ErrGateRejected
	}
	date := time.Now().Format("2006-01-02")
	govDir := filepath.Join(baseDir, "governance")
	processoDir := filepath.Join(baseDir, "processo")
	if err := makeDirs(govDir, processoDir); err != nil {
		return nil, err
	}
	planFile := filepath.Join(govDir, fmt.Sprintf("%s_PLAN_%s.md", date, spec.Slug))
	routeFile := filepath.Join(govDir, fmt.Sprintf("%s_ROUTE_%s.json", date, spec.Slug))
	dagFile := caminhos.CaminhoDAG(processoDir)
	if err := writeExecutionFiles(planFile, routeFile, dagFile, spec.Slug); err != nil {
		return nil, err
	}
	fixPath, err := SynthesizeGoldenFixture(3.0, processoDir)
	if err != nil {
		return nil, err
	}
	return &ExecutionPackageDraft{
		Slug: spec.Slug, PlanPath: planFile, RoutePath: routeFile, DAGPath: dagFile, FixturePath: fixPath,
	}, nil
}

func makeDirs(dirs ...string) error {
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}
	return nil
}

func writeExecutionFiles(plan, route, dag, slug string) error {
	planData := fmt.Sprintf("# PLAN: %s\n\nStatus: UNDER_REVIEW\n", slug)
	routeData := fmt.Sprintf(`{"name": "%s", "version": "1.0.0", "status": "PENDING"}`, slug)
	dagData := fmt.Sprintf(`{"name": "%s", "version": "1.0.0", "nodes": []}`, slug)
	items := map[string]string{plan: planData, route: routeData, dag: dagData}
	for p, c := range items {
		if err := os.WriteFile(p, []byte(c), 0644); err != nil {
			return err
		}
	}
	return nil
}

// GenerateGovernanceTriad emite a governança ISO completa (mantém retrocompatibilidade).
func GenerateGovernanceTriad(report *CurationReport, outDir string) (*GovernanceDraft, error) {
	intent, err := GenerateIntentStage(report, outDir, 95)
	if err != nil {
		return nil, err
	}
	spec, err := GenerateSpecStage(intent, report, outDir)
	if err != nil {
		return nil, err
	}
	pkg, err := GenerateExecutionPackageStage(spec, report, filepath.Dir(outDir))
	if err != nil {
		return nil, err
	}
	return &GovernanceDraft{
		Slug: report.SuggestedName, IntentPath: intent.IntentPath, SpecPath: spec.SpecPath,
		PlanPath: pkg.PlanPath, RoutePath: pkg.RoutePath, BookDir: filepath.Join(outDir, "books", spec.Slug),
		GoldenFixture: pkg.FixturePath,
	}, nil
}

// SynthesizeGoldenFixture cria um arquivo sintético determinístico de teste de 3s.
func SynthesizeGoldenFixture(durationSec float64, outDir string) (string, error) {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return "", err
	}
	fixturePath := filepath.Join(outDir, "fixture_3s.mp4")
	meta := map[string]interface{}{
		"format": "mp4", "duration_sec": durationSec, "fps": 60, "audio_hz": 48000, "channels": 2, "punch_word": "LEGO",
	}
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return "", err
	}
	return fixturePath, os.WriteFile(fixturePath, data, 0644)
}

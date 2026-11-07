package videorea

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// GenerateGovernanceTriad emite os artefatos de governança ISO para o pipeline aprovado.
func GenerateGovernanceTriad(report *CurationReport, outDir string) (*GovernanceDraft, error) {
	if report == nil || !report.Approved {
		return nil, ErrGateRejected
	}
	date := time.Now().Format("2006-01-02")
	slug := report.SuggestedName
	draft := &GovernanceDraft{
		Slug:          slug,
		IntentPath:    filepath.Join(outDir, fmt.Sprintf("%s_INTENT_%s.md", date, slug)),
		SpecPath:      filepath.Join(outDir, fmt.Sprintf("%s_SPEC_%s.md", date, slug)),
		PlanPath:      filepath.Join(outDir, fmt.Sprintf("%s_PLAN_%s.md", date, slug)),
		RoutePath:     filepath.Join(outDir, fmt.Sprintf("%s_ROUTE_%s.json", date, slug)),
		BookDir:       filepath.Join(outDir, "books", slug),
		GoldenFixture: filepath.Join(outDir, "testdata", "fixture_3s.mp4"),
	}
	if err := writeDraftFiles(draft, report); err != nil {
		return nil, err
	}
	return draft, nil
}

// SynthesizeGoldenFixture cria um arquivo sintético determinístico de teste de 3s.
func SynthesizeGoldenFixture(durationSec float64, outDir string) (string, error) {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return "", err
	}
	fixturePath := filepath.Join(outDir, "fixture_3s.mp4")
	meta := map[string]interface{}{
		"format":       "mp4",
		"duration_sec": durationSec,
		"fps":          60,
		"audio_hz":     48000,
		"channels":     2,
		"punch_word":   "LEGO",
	}
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return "", err
	}
	return fixturePath, os.WriteFile(fixturePath, data, 0644)
}

func writeDraftFiles(d *GovernanceDraft, r *CurationReport) error {
	dir := filepath.Dir(d.IntentPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	intent := fmt.Sprintf("# INTENT: %s\n\nCusto Soberano: %s\n", r.SuggestedName, r.SovereignCost)
	spec := fmt.Sprintf("# SPEC: %s\n\nSubstituições: %v\n", r.SuggestedName, r.Replacements)
	plan := fmt.Sprintf("# PLAN: %s\n\nSlug: %s\n", r.SuggestedName, d.Slug)
	route := fmt.Sprintf(`{"name": "%s", "version": "1.0.0"}`, d.Slug)

	files := map[string]string{
		d.IntentPath: intent,
		d.SpecPath:   spec,
		d.PlanPath:   plan,
		d.RoutePath:  route,
	}
	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

package videorea

import (
	"fmt"
	"os"
	"path/filepath"
)

// MaterializeAndProject cria o bundle isolado e projeta pontes na raiz canônica.
func MaterializeAndProject(report *CurationReport, baseDir, repoRoot string) (*TargetBundle, error) {
	if report == nil || !report.Approved {
		return nil, ErrGateRejected
	}
	bundle := NewTargetBundle(baseDir, report.SuggestedName)
	if err := bundle.Materialize(); err != nil {
		return nil, err
	}
	if err := writeBundleArtifacts(bundle, report); err != nil {
		return nil, err
	}
	if err := projectCanonical(bundle, repoRoot); err != nil {
		return nil, err
	}
	return bundle, nil
}

func writeBundleArtifacts(b *TargetBundle, r *CurationReport) error {
	if err := writeLocalGovernance(b, r); err != nil {
		return err
	}
	if err := writeLocalBook(b, r); err != nil {
		return err
	}
	if err := writeLocalLoop(b, r); err != nil {
		return err
	}
	return writeLocalSkill(b, r)
}

func writeLocalGovernance(b *TargetBundle, r *CurationReport) error {
	intent := fmt.Sprintf("# INTENT: %s\n\nCusto Soberano: %s\n", r.SuggestedName, r.SovereignCost)
	spec := fmt.Sprintf("# SPEC: %s\n\nSubstituicoes: %v\n", r.SuggestedName, r.Replacements)
	plan := fmt.Sprintf("# PLAN: %s\n", r.SuggestedName)
	route := fmt.Sprintf(`{"name": %q, "version": "1.0.0"}`, r.SuggestedName)

	items := map[string]string{
		filepath.Join(b.GovernanceDir, "INTENT.md"):  intent,
		filepath.Join(b.GovernanceDir, "SPEC.md"):    spec,
		filepath.Join(b.GovernanceDir, "PLAN.md"):    plan,
		filepath.Join(b.GovernanceDir, "ROUTE.json"): route,
	}
	for p, content := range items {
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func writeLocalBook(b *TargetBundle, r *CurationReport) error {
	bookJSON := fmt.Sprintf(`{"nome_base": %q, "titulo": %q}`, b.Slug, r.SuggestedName)
	return os.WriteFile(filepath.Join(b.BookDir, "livro.json"), []byte(bookJSON), 0o644)
}

func writeLocalLoop(b *TargetBundle, r *CurationReport) error {
	dagJSON := fmt.Sprintf(`{"name": %q, "version": "1.0.0", "nodes": []}`, r.SuggestedName)
	if err := os.WriteFile(filepath.Join(b.LoopDir, "pipeline.dag.json"), []byte(dagJSON), 0o644); err != nil {
		return err
	}
	_, err := SynthesizeGoldenFixture(3.0, b.LoopDir)
	return err
}

func writeLocalSkill(b *TargetBundle, r *CurationReport) error {
	skillMD := fmt.Sprintf("---\nname: %s\ndescription: Skill de %s\n---\n# %s\n", b.Slug, r.SuggestedName, b.Slug)
	return os.WriteFile(filepath.Join(b.SkillDir, "SKILL.md"), []byte(skillMD), 0o644)
}

func projectCanonical(b *TargetBundle, repoRoot string) error {
	loopDest := filepath.Join(repoRoot, "loops", b.Slug)
	if err := os.MkdirAll(loopDest, 0o755); err != nil {
		return err
	}
	dagSrc, err := os.ReadFile(filepath.Join(b.LoopDir, "pipeline.dag.json"))
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(loopDest, "pipeline.dag.json"), dagSrc, 0o644); err != nil {
		return err
	}

	skillDest := filepath.Join(repoRoot, "skills", b.Slug)
	if err := os.MkdirAll(skillDest, 0o755); err != nil {
		return err
	}
	skillSrc, err := os.ReadFile(filepath.Join(b.SkillDir, "SKILL.md"))
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(skillDest, "SKILL.md"), skillSrc, 0o644)
}

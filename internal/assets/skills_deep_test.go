package assets

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestFindSkillFilesRecursiveInMonorepo(t *testing.T) {
	tmp := t.TempDir()
	// Monorepo macro-skill
	macro := filepath.Join(tmp, "aidd-osmani")
	if err := os.MkdirAll(macro, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(macro, "SKILL.md"), []byte("---\nname: macro\n---\nMacro"), 0600); err != nil {
		t.Fatal(err)
	}

	// Sub-skill aninhada em upstream/skills/spec
	specSkill := filepath.Join(macro, "upstream", "skills", "spec")
	if err := os.MkdirAll(specSkill, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(specSkill, "SKILL.md"), []byte("---\nname: spec\n---\nSpec"), 0600); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	files := findSkillFiles(ctx, tmp, 0)
	if len(files) < 2 {
		t.Fatalf("esperava pelo menos 2 skills encontradas no monorepo, obteve %d: %v", len(files), files)
	}
}

func TestSkillRootsContainsEcosystemDirs(t *testing.T) {
	s := NewScanner(Config{Home: "C:/fakehome"})
	roots := s.skillRoots()
	foundGeminiCli := false
	foundCline := false
	for _, r := range roots {
		if filepath.ToSlash(r.dir) == "C:/fakehome/.gemini/antigravity-cli/skills" {
			foundGeminiCli = true
		}
		if filepath.ToSlash(r.dir) == "C:/fakehome/.cline/skills" {
			foundCline = true
		}
	}
	if !foundGeminiCli {
		t.Error("skillRoots deve conter .gemini/antigravity-cli/skills")
	}
	if !foundCline {
		t.Error("skillRoots deve conter .cline/skills")
	}
}

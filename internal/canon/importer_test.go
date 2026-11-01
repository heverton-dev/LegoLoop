package canon_test

import (
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/canon"
)

func TestParseGitHubURL(t *testing.T) {
	cases := []struct {
		url       string
		wantRepo  string
		wantPath  string
		wantName  string
		shouldErr bool
	}{
		{
			url:      "https://github.com/owner/minha-skill",
			wantRepo: "minha-skill",
			wantPath: "",
			wantName: "minha-skill",
		},
		{
			url:      "https://github.com/owner/repo/tree/main/skills/sub-skill",
			wantRepo: "repo",
			wantPath: "skills/sub-skill",
			wantName: "sub-skill",
		},
		{
			url:       "https://evil.example.com/owner/repo",
			shouldErr: true,
		},
		{
			url:       "https://github.com/owner/repo/../trav",
			shouldErr: true,
		},
		{
			url:       "not-a-url",
			shouldErr: true,
		},
	}

	for _, tc := range cases {
		info, err := canon.ParseGitHubURL(tc.url)
		if tc.shouldErr {
			if err == nil {
				t.Errorf("URL %q esperava erro, mas passou", tc.url)
			}
			continue
		}
		if err != nil {
			t.Fatalf("URL %q falhou com erro: %v", tc.url, err)
		}
		if info.Repo != tc.wantRepo {
			t.Errorf("URL %q: Repo got %q, want %q", tc.url, info.Repo, tc.wantRepo)
		}
		if info.Subpath != tc.wantPath {
			t.Errorf("URL %q: Subpath got %q, want %q", tc.url, info.Subpath, tc.wantPath)
		}
		if info.SuggestedName != tc.wantName {
			t.Errorf("URL %q: SuggestedName got %q, want %q", tc.url, info.SuggestedName, tc.wantName)
		}
	}
}

func TestInstallAndCanonizeLocalSkill(t *testing.T) {
	srcDir, err := os.MkdirTemp("", "src_skill_*")
	if err != nil {
		t.Fatalf("falha ao criar srcDir: %v", err)
	}
	defer os.RemoveAll(srcDir)

	skillContent := "---\nname: imported-skill\ndescription: skill importada do github\nconsumes: Artifact:Idea\nproduces: Artifact:Briefing\n---\n# Imported Skill\n"
	if err := os.WriteFile(filepath.Join(srcDir, "SKILL.md"), []byte(skillContent), 0644); err != nil {
		t.Fatalf("falha ao criar SKILL.md: %v", err)
	}

	destRoot, err := os.MkdirTemp("", "dest_vault_*")
	if err != nil {
		t.Fatalf("falha ao criar destRoot: %v", err)
	}
	defer os.RemoveAll(destRoot)

	res, err := canon.InstallAndCanonizeSkill(srcDir, "imported-skill", destRoot)
	if err != nil {
		t.Fatalf("InstallAndCanonizeSkill falhou: %v", err)
	}

	if res.Name != "imported-skill" {
		t.Errorf("Nome esperado imported-skill, obteve %q", res.Name)
	}
	if res.SHA256 == "" {
		t.Errorf("SHA256 não deve ser vazio")
	}
	if res.FilesCount < 1 {
		t.Errorf("FilesCount deve ser >= 1, obteve %d", res.FilesCount)
	}

	targetSkillMD := filepath.Join(destRoot, "skills", "imported-skill", "SKILL.md")
	if _, err := os.Stat(targetSkillMD); err != nil {
		t.Fatalf("arquivo alvo %s não foi criado: %v", targetSkillMD, err)
	}
}

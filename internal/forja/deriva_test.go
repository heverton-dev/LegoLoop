package forja

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func criarSkillDerivaTeste(t *testing.T, dirSkill, nome string) {
	t.Helper()
	if err := os.MkdirAll(dirSkill, 0o755); err != nil {
		t.Fatal(err)
	}
	conteudo := "---\nname: " + nome + "\ndescription: skill de teste\n---\n# Corpo\n"
	if err := os.WriteFile(filepath.Join(dirSkill, "SKILL.md"), []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestDerivaReprova(t *testing.T) {
	home := t.TempDir()
	criarSkillDerivaTeste(t, filepath.Join(home, ".claude", "skills", "legoloop-skills"), "legoloop-skills")
	pastaEager := filepath.Join(home, ".claude", "skills", "minha-skill-eager")
	criarSkillDerivaTeste(t, pastaEager, "minha-skill-eager")
	pastaOpencode := filepath.Join(home, ".config", "opencode", "skills", "outra-eager")
	criarSkillDerivaTeste(t, pastaOpencode, "outra-eager")

	derivas, err := DetectarDeriva(home)
	if err != nil {
		t.Fatalf("DetectarDeriva falhou: %v", err)
	}
	if len(derivas) != 2 {
		t.Fatalf("esperava 2 derivas, obteve %d: %+v", len(derivas), derivas)
	}
	if derivas[0].Harness != "claude" || derivas[0].Skill != "minha-skill-eager" || derivas[0].Pasta != pastaEager {
		t.Fatalf("primeira deriva inesperada: %+v", derivas[0])
	}
	if derivas[1].Harness != "opencode" || derivas[1].Skill != "outra-eager" || derivas[1].Pasta != pastaOpencode {
		t.Fatalf("segunda deriva inesperada: %+v", derivas[1])
	}
	texto := FormatarDerivas(derivas)
	if !strings.Contains(texto, "claude") || !strings.Contains(texto, filepath.ToSlash(pastaEager)) {
		t.Fatalf("formatacao deve citar harness e pasta: %s", texto)
	}
}

func TestDerivaLimpaAprova(t *testing.T) {
	home := t.TempDir()
	criarSkillDerivaTeste(t, filepath.Join(home, ".claude", "skills", "legoloop-skills"), "legoloop-skills")
	criarSkillDerivaTeste(t, filepath.Join(home, ".gemini", "skills", "legoloop-skills"), "legoloop-skills")

	derivas, err := DetectarDeriva(home)
	if err != nil {
		t.Fatalf("DetectarDeriva falhou: %v", err)
	}
	if len(derivas) != 0 {
		t.Fatalf("esperava 0 derivas, obteve %+v", derivas)
	}
}

func TestDerivaHomeVazio(t *testing.T) {
	if _, err := DetectarDeriva("   "); !errors.Is(err, ErrSemHome) {
		t.Fatalf("esperava ErrSemHome, obteve %v", err)
	}
}

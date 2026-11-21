package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/skills"
)

// TestSkillsMigrarCLI prova o comando: sem flag recusa; --dry-run só conta; --apply move ao Cofre e
// deixa o harness só com o índice.
func TestSkillsMigrarCLI(t *testing.T) {
	d, out, errb := deps(t)
	vault := filepath.Join(d.Home, ".aidd", "canonical")
	pasta := filepath.Join(d.Home, ".cline", "skills")
	skillCusto(t, pasta, "alfa")
	skillCusto(t, pasta, "beta")

	if code := run(t, d, "skills-migrar", "--vault", vault); code != 2 {
		t.Fatalf("sem --dry-run/--apply deve sair com 2, saiu %d", code)
	}
	out.Reset()
	if code := run(t, d, "skills-migrar", "--dry-run", "--json", "--vault", vault); code != 0 {
		t.Fatalf("dry-run saiu %d: %s", code, errb.String())
	}
	var rel skills.MigrarRelatorio
	if err := json.Unmarshal(out.Bytes(), &rel); err != nil {
		t.Fatalf("saída não é JSON: %v\n%s", err, out.String())
	}
	if len(rel.Harnesses) == 0 || rel.Harnesses[0].Nome != "cline" || rel.Harnesses[0].Mover != 2 {
		t.Fatalf("dry-run = %+v, quer cline com 2 a mover", rel.Harnesses)
	}
	if _, err := os.Stat(filepath.Join(pasta, "alfa", "SKILL.md")); err != nil {
		t.Fatal("dry-run não pode mover")
	}
	out.Reset()
	if code := run(t, d, "skills-migrar", "--apply", "--vault", vault); code != 0 {
		t.Fatalf("apply saiu %d: %s", code, errb.String())
	}
	for _, nome := range []string{"alfa", "beta"} {
		if _, err := os.Stat(filepath.Join(vault, "skills", nome, "SKILL.md")); err != nil {
			t.Errorf("%s deveria estar no Cofre: %v", nome, err)
		}
	}
	if _, err := os.Stat(filepath.Join(pasta, "legoloop-skills", "SKILL.md")); err != nil {
		t.Error("o índice deveria estar no harness")
	}
}

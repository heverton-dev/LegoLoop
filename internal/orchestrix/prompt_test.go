package orchestrix

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/dag"
)

func setupPromptTestEnv(t *testing.T) (*Orchestrix, *os.Root) {
	t.Helper()
	skillsDir := t.TempDir()
	workDir := t.TempDir()

	mainSkillDir := filepath.Join(skillsDir, "peca-principal")
	if err := os.MkdirAll(mainSkillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(mainSkillDir, "SKILL.md"), []byte("Instrucao da peca principal."), 0o644); err != nil {
		t.Fatal(err)
	}

	extraSkillDir := filepath.Join(skillsDir, "regra-extra")
	if err := os.MkdirAll(extraSkillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(extraSkillDir, "SKILL.md"), []byte("Regra extra de tooling."), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(workDir, "spec.md"), []byte("Conteudo da spec de entrada."), 0o644); err != nil {
		t.Fatal(err)
	}

	root, err := os.OpenRoot(workDir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = root.Close() })

	o := &Orchestrix{
		SkillsDir: skillsDir,
		Vars:      map[string]string{},
	}
	return o, root
}

func TestBuildPromptInjectsMinimalCoreOnly(t *testing.T) {
	o, root := setupPromptTestEnv(t)
	n := dag.Node{
		ID:     "n1",
		Ref:    "peca-principal",
		Output: dag.HandoffContract{Artifact: "out.md", Gate: "exit_code == 0"},
	}

	got, err := o.buildPrompt(n, root)
	if err != nil {
		t.Fatalf("buildPrompt falhou: %v", err)
	}
	for _, want := range []string{"STE-PT", "Zero Stubs", "Path Jail"} {
		if !strings.Contains(got, want) {
			t.Errorf("prompt deveria conter núcleo mínimo %q, obtido:\n%s", want, got)
		}
	}
	if strings.Contains(got, "As 5 Prisões Mecânicas") {
		t.Errorf("prompt não deveria carregar regras globais completas do AGENTS.md")
	}
}

func TestBuildPromptInjectsContextIn(t *testing.T) {
	o, root := setupPromptTestEnv(t)
	n := dag.Node{
		ID:        "n2",
		Ref:       "peca-principal",
		ContextIn: []string{"spec.md"},
		Output:    dag.HandoffContract{Artifact: "out.md", Gate: "exit_code == 0"},
	}

	got, err := o.buildPrompt(n, root)
	if err != nil {
		t.Fatalf("buildPrompt falhou: %v", err)
	}
	if !strings.Contains(got, "Conteudo da spec de entrada.") {
		t.Errorf("prompt deveria conter conteúdo de context_in spec.md, obtido:\n%s", got)
	}
}

func TestBuildPromptInjectsToolingSkills(t *testing.T) {
	o, root := setupPromptTestEnv(t)
	n := dag.Node{
		ID:      "n3",
		Ref:     "peca-principal",
		Tooling: map[string]any{"skills": []string{"regra-extra"}},
		Output:  dag.HandoffContract{Artifact: "out.md", Gate: "exit_code == 0"},
	}

	got, err := o.buildPrompt(n, root)
	if err != nil {
		t.Fatalf("buildPrompt falhou: %v", err)
	}
	if !strings.Contains(got, "Regra extra de tooling.") {
		t.Errorf("prompt deveria conter skill referenciada em tooling, obtido:\n%s", got)
	}
}

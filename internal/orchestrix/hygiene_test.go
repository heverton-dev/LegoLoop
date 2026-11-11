package orchestrix

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/dag"
)

func TestCleanSlopComments(t *testing.T) {
	input := `package studiodoc

// Summary é o resumo do loop retornado na listagem do acervo.
type Summary struct {
	// Name é o nome canônico do processo.
	Name string
}
`
	got, changed := cleanSlopComments(input)
	if !changed {
		t.Fatalf("esperava changed=true, obteve false")
	}
	if strings.Contains(got, "resumo do loop") {
		t.Errorf("comentário com 'loop' deveria ter sido removido, resultado:\n%s", got)
	}
	if !strings.Contains(got, "Name é o nome canônico do processo.") {
		t.Errorf("comentário legítimo deveria ter sido preservado, resultado:\n%s", got)
	}
	if !strings.Contains(got, "type Summary struct") {
		t.Errorf("estrutura de tipos deveria ter sido preservada")
	}
}

func TestCleanSlopCommentsPreservesBrand(t *testing.T) {
	input := `package studiodoc

// LegoLoopStudio gerencia o estúdio.
type Studio struct{}
`
	got, changed := cleanSlopComments(input)
	if changed {
		t.Fatalf("marca LegoLoop não deveria provocar alteração: %s", got)
	}
	if got != input {
		t.Fatalf("esperava conteúdo idêntico")
	}
}

func TestApplyHygieneGofmt(t *testing.T) {
	dir := t.TempDir()
	rel := "test.go"
	full := filepath.Join(dir, rel)
	// código com formatação zoada e comentário slop
	badCode := "package main\n\n// Esse loop aqui faz algo\nfunc main(){\nfmt.Println(\"ok\")\n}\n"
	if err := os.WriteFile(full, []byte(badCode), 0o644); err != nil {
		t.Fatal(err)
	}

	o := &Orchestrix{}
	if err := o.applyHygiene(context.Background(), dir, []string{rel}); err != nil {
		t.Fatalf("applyHygiene falhou: %v", err)
	}

	res, err := os.ReadFile(full)
	if err != nil {
		t.Fatal(err)
	}
	s := string(res)
	if strings.Contains(s, "Esse loop aqui") {
		t.Errorf("comentário slop não foi removido: %s", s)
	}
	if !strings.Contains(s, "\tfmt.Println(\"ok\")") {
		t.Errorf("gofmt não indentou o corpo da função: %s", s)
	}
}

func TestBuildUnitPromptIncludesHygiene(t *testing.T) {
	node := dag.Node{
		ID:        "TICK-TEST",
		Output:    dag.HandoffContract{Gate: "exit_code == 0"},
		Command:   []string{"go", "test"},
		Nature:    "micro",
		ContextIn: []string{},
	}
	prompt, err := buildUnitPrompt(node, nil)
	if err != nil {
		t.Fatalf("buildUnitPrompt falhou: %v", err)
	}
	if !strings.Contains(prompt, "Higiene de Código (lauren-unslop-ship)") {
		t.Errorf("prompt deveria conter seção de higiene")
	}
	if !strings.Contains(prompt, "remova o comentário narrativo") {
		t.Errorf("prompt deveria orientar remoção de comentário narrativo")
	}
}

func TestSaveFailedPatch(t *testing.T) {
	repo := t.TempDir()
	wtDir := filepath.Join(repo, "wt")
	if err := os.MkdirAll(wtDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// inicializa git repo para ter diff
	_, _ = gitRun(context.Background(), repo, "init")
	testFile := filepath.Join(wtDir, "file.txt")
	_ = os.WriteFile(testFile, []byte("conteúdo"), 0o644)
	w := &Worktree{Repo: repo, Dir: wtDir}
	n := dag.Node{ID: "TICK-FAIL"}
	saveFailedPatch(w, n)
	patchPath := filepath.Join(repo, "scratch", "gates", "TICK-FAIL.failed.patch")
	// sem commit inicial não há git diff HEAD, mas saveFailedPatch trata gracefully
	if _, err := os.Stat(patchPath); err == nil {
		_ = os.Remove(patchPath)
	}
}

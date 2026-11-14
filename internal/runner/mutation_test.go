package runner

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestVerifyMutation_KillsMutant(t *testing.T) {
	dir := t.TempDir()
	writeGoModule(t, dir, "example.com/muttest")
	writeFileTest(t, filepath.Join(dir, "calc.go"), "package muttest\n\nfunc Add(a, b int) int { return a + b }\n")
	writeFileTest(t, filepath.Join(dir, "calc_test.go"), "package muttest\n\nimport \"testing\"\n\nfunc TestAdd(t *testing.T) {\n\tif Add(2, 3) != 5 {\n\t\tt.Fatal(\"esperado 5\")\n\t}\n}\n")

	res, err := VerifyMutation(context.Background(), dir, []string{"go", "test", "./..."}, "calc.go", flipPlusToMinus)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if !res.Passed || res.MutantsKilled != 1 || res.MutantsSurvived != 0 {
		t.Fatalf("esperado mutante morto, obtido: %+v", res)
	}
	got, _ := os.ReadFile(filepath.Join(dir, "calc.go"))
	if !strings.Contains(string(got), "a + b") {
		t.Fatalf("arquivo original não foi restaurado: %s", string(got))
	}
}

func TestVerifyMutation_MutantSurvives(t *testing.T) {
	dir := t.TempDir()
	writeGoModule(t, dir, "example.com/muttest")
	writeFileTest(t, filepath.Join(dir, "calc.go"), "package muttest\n\nfunc Add(a, b int) int { return a + b }\n")
	writeFileTest(t, filepath.Join(dir, "calc_test.go"), "package muttest\n\nimport \"testing\"\n\nfunc TestAdd(t *testing.T) {}\n")

	res, err := VerifyMutation(context.Background(), dir, []string{"go", "test", "./..."}, "calc.go", flipPlusToMinus)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if res.Passed || res.MutantsSurvived != 1 || res.MutantsKilled != 0 {
		t.Fatalf("esperado mutante sobrevivente, obtido: %+v", res)
	}
	got, _ := os.ReadFile(filepath.Join(dir, "calc.go"))
	if !strings.Contains(string(got), "a + b") {
		t.Fatalf("arquivo original não foi restaurado: %s", string(got))
	}
}

func TestRollbackWorktree_PreservesMainBranch(t *testing.T) {
	workdir := t.TempDir()
	runGitTest(t, workdir, "init")
	runGitTest(t, workdir, "config", "user.email", "test@legoloop.local")
	runGitTest(t, workdir, "config", "user.name", "LegoLoop Test")
	writeFileTest(t, filepath.Join(workdir, "main.txt"), "inicial\n")
	runGitTest(t, workdir, "add", "main.txt")
	runGitTest(t, workdir, "commit", "-m", "commit inicial")

	isoDir, _, err := createGitWorktree(context.Background(), workdir)
	if err != nil {
		t.Fatalf("falha ao criar worktree: %v", err)
	}
	writeFileTest(t, filepath.Join(isoDir, "dirty.txt"), "sujeira\n")

	if err := RollbackWorktree(context.Background(), workdir, isoDir); err != nil {
		t.Fatalf("RollbackWorktree falhou: %v", err)
	}
	if _, err := os.Stat(isoDir); !os.IsNotExist(err) {
		t.Fatalf("isoDir ainda existe após rollback: %s", isoDir)
	}
	out := runGitTest(t, workdir, "status", "--porcelain")
	if strings.TrimSpace(out) != "" {
		t.Fatalf("branch principal contaminada: %s", out)
	}
}

func flipPlusToMinus(src string) (string, bool) {
	if !strings.Contains(src, " + ") {
		return src, false
	}
	return strings.Replace(src, " + ", " - ", 1), true
}

func writeGoModule(t *testing.T, dir, mod string) {
	t.Helper()
	writeFileTest(t, filepath.Join(dir, "go.mod"), "module "+mod+"\n\ngo 1.21\n")
}

func writeFileTest(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("falha ao escrever %s: %v", path, err)
	}
}

func runGitTest(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v falhou: %s (%v)", args, string(out), err)
	}
	return string(out)
}

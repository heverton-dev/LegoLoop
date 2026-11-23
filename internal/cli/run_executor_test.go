package cli

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func gravarDAGInvalidoXYZ(t *testing.T, dir string) string {
	t.Helper()
	conteudo := `{
  "name": "dag-executor-xyz",
  "version": "1.0.0",
  "nodes": [
    {
      "id": "no-invalido-xyz",
      "type": "unit_no",
      "isolation": {"executor": "xyz"},
      "routing": {"mode": "manual", "model": "m1"},
      "output": {"artifact": "saida.md", "gate": "file_exists('saida.md')"}
    }
  ]
}`
	arq := filepath.Join(dir, "xyz.dag.json")
	if err := os.WriteFile(arq, []byte(conteudo), 0o600); err != nil {
		t.Fatalf("falha ao gravar DAG invalido: %v", err)
	}
	return arq
}

func contarWorktrees(t *testing.T, repo string) int {
	t.Helper()
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	cmd.Dir = repo
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git worktree list falhou: %v (%s)", err, out)
	}
	return bytes.Count(out, []byte("worktree "))
}

func TestRunRecusaExecutorSemWorktree(t *testing.T) {
	bin := buildBinarioTeste(t)
	repo := e2eGitRepo(t)
	dagPath := gravarDAGInvalidoXYZ(t, repo)

	antes := contarWorktrees(t, repo)
	code, saida := executarBinario(bin, "run", dagPath, "--workdir", repo)
	depois := contarWorktrees(t, repo)

	if code != 1 {
		t.Fatalf("esperava exit code 1, obteve %d\nsaida:\n%s", code, saida)
	}
	if !strings.Contains(saida, "no-invalido-xyz") {
		t.Fatalf("saida deve citar no invalido (no-invalido-xyz): %s", saida)
	}
	if !strings.Contains(saida, "xyz") {
		t.Fatalf("saida deve citar executor invalido (xyz): %s", saida)
	}
	if depois != antes {
		t.Fatalf("worktrees aumentaram de %d para %d", antes, depois)
	}
}

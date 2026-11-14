package runner_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"legoloop/internal/runner"
)

func initTestGitRepo(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "legoloop-git-test-*")
	if err != nil {
		t.Fatalf("falha ao criar temp dir: %v", err)
	}

	runGit := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v falhou: %s (%v)", args, string(out), err)
		}
	}

	runGit("init")
	runGit("config", "user.name", "LegoLoop Test")
	runGit("config", "user.email", "test@legoloop.local")

	readme := filepath.Join(dir, "README.md")
	if err := os.WriteFile(readme, []byte("# Test Repo\n"), 0644); err != nil {
		t.Fatalf("falha ao criar README: %v", err)
	}

	runGit("add", "README.md")
	runGit("commit", "-m", "initial commit")

	return dir
}

func TestTicketLifecycleWorktree(t *testing.T) {
	repoDir := initTestGitRepo(t)
	defer os.RemoveAll(repoDir)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	wtDir, branch, cleanup, err := runner.PrepareTicketWorktree(ctx, repoDir, 42, "Implement Wave Runner")
	if err != nil {
		t.Fatalf("falha ao criar worktree do Ticket: %v", err)
	}
	defer cleanup()

	if branch != "ticket-42-implement-wave-runner" {
		t.Fatalf("branch incorreta: %s", branch)
	}

	// Verifica se wtDir existe e contém o README inicial
	readmePath := filepath.Join(wtDir, "README.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		t.Fatalf("README.md não encontrado dentro do worktree: %v", err)
	}

	// Simula alteração e commit dentro do worktree
	newFile := filepath.Join(wtDir, "lifecycle.go")
	if err := os.WriteFile(newFile, []byte("package runner\n"), 0644); err != nil {
		t.Fatalf("falha ao criar arquivo no worktree: %v", err)
	}

	err = runner.CommitTicketWorktree(ctx, wtDir, "feat(runner): implement ticket lifecycle")
	if err != nil {
		t.Fatalf("falha ao commitar no worktree: %v", err)
	}
}

package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"legoloop/internal/board"
)

var nonAlphanumericSlug = regexp.MustCompile(`[^a-z0-9]+`)

// generateIssueBranchName gera uma branch padronizada como issue-42-implement-wave-runner.
func generateIssueBranchName(card board.Card) string {
	slug := strings.ToLower(card.Title)
	slug = nonAlphanumericSlug.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if len(slug) > 40 {
		slug = slug[:40]
		slug = strings.TrimRight(slug, "-")
	}
	if slug == "" {
		slug = "task"
	}
	return fmt.Sprintf("issue-%d-%s", card.Number, slug)
}

// PrepareIssueWorktree provisiona um worktree com branch dedicada para uma issue do GitHub.
// Retorna o diretório do worktree, o nome da branch gerada e uma função de cleanup seguro.
func PrepareIssueWorktree(ctx context.Context, repoDir string, card board.Card) (string, string, func(), error) {
	branch := generateIssueBranchName(card)

	tmpDir, err := os.MkdirTemp("", fmt.Sprintf("legoloop-issue-%d-*", card.Number))
	if err != nil {
		return "", "", nil, fmt.Errorf("falha ao criar pasta temporária para worktree: %w", err)
	}

	cmdCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Cria o worktree com a nova branch associada
	cmd := exec.CommandContext(cmdCtx, "git", "-C", repoDir, "worktree", "add", "-b", branch, tmpDir, "HEAD")
	if out, err := cmd.CombinedOutput(); err != nil {
		_ = os.RemoveAll(tmpDir)
		return "", "", nil, fmt.Errorf("git worktree add -b %s falhou: %s (%w)", branch, strings.TrimSpace(string(out)), err)
	}

	cleanup := func() {
		cleanCtx, cleanCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanCancel()
		_ = exec.CommandContext(cleanCtx, "git", "-C", repoDir, "worktree", "remove", "--force", tmpDir).Run()
		_ = os.RemoveAll(tmpDir)
	}

	return tmpDir, branch, cleanup, nil
}

// CommitIssueWorktree adiciona todas as alterações e executa o commit dentro do worktree.
func CommitIssueWorktree(ctx context.Context, worktreeDir, message string) error {
	cmdCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	addCmd := exec.CommandContext(cmdCtx, "git", "-C", worktreeDir, "add", "-A")
	if out, err := addCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git add falhou: %s (%w)", strings.TrimSpace(string(out)), err)
	}

	commitCmd := exec.CommandContext(cmdCtx, "git", "-C", worktreeDir, "commit", "-m", message)
	if out, err := commitCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git commit falhou: %s (%w)", strings.TrimSpace(string(out)), err)
	}

	return nil
}

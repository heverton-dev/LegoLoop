package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var nonAlphanumericSlug = regexp.MustCompile(`[^a-z0-9]+`)

// generateTicketBranchName gera uma branch padronizada como ticket-42-implement-wave-runner.
func generateTicketBranchName(number int, title string) string {
	slug := strings.ToLower(title)
	slug = nonAlphanumericSlug.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if len(slug) > 40 {
		slug = slug[:40]
		slug = strings.TrimRight(slug, "-")
	}
	if slug == "" {
		slug = "sem-titulo"
	}
	return fmt.Sprintf("ticket-%d-%s", number, slug)
}

// PrepareTicketWorktree provisiona um worktree com branch dedicada para um Ticket.
// Retorna o diretório do worktree, o nome da branch gerada e uma função de cleanup seguro.
func PrepareTicketWorktree(ctx context.Context, repoDir string, number int, title string) (string, string, func(), error) {
	branch := generateTicketBranchName(number, title)

	tmpDir, err := os.MkdirTemp("", fmt.Sprintf("legoloop-ticket-%d-*", number))
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

// CommitTicketWorktree adiciona todas as alterações e executa o commit dentro do worktree.
func CommitTicketWorktree(ctx context.Context, worktreeDir, message string) error {
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

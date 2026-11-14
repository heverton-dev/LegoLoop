package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// MutationResult registra o resultado da verificação de mutação na baia.
type MutationResult struct {
	Passed          bool
	MutantsTotal    int
	MutantsKilled   int
	MutantsSurvived int
}

// VerifyMutation aplica uma mutação controlada em targetFile dentro de isoDir,
// executa testCmd e garante a restauração do arquivo original.
func VerifyMutation(ctx context.Context, isoDir string, testCmd []string, targetFile string, mutateFn func(string) (string, bool)) (MutationResult, error) {
	if len(testCmd) == 0 || mutateFn == nil {
		return MutationResult{}, fmt.Errorf("comando de teste e função de mutação são obrigatórios")
	}
	absTarget, err := resolveSafePath(isoDir, targetFile)
	if err != nil {
		return MutationResult{}, err
	}
	origBytes, err := os.ReadFile(absTarget)
	if err != nil {
		return MutationResult{}, fmt.Errorf("falha ao ler arquivo alvo: %w", err)
	}
	mutated, changed := mutateFn(string(origBytes))
	if !changed {
		return MutationResult{Passed: true}, nil
	}
	if err := os.WriteFile(absTarget, []byte(mutated), 0o644); err != nil {
		return MutationResult{}, fmt.Errorf("falha ao aplicar mutação: %w", err)
	}
	defer func() { _ = os.WriteFile(absTarget, origBytes, 0o644) }()

	cmd := exec.CommandContext(ctx, testCmd[0], testCmd[1:]...)
	cmd.Dir = isoDir
	if err := cmd.Run(); err != nil {
		return MutationResult{Passed: true, MutantsTotal: 1, MutantsKilled: 1}, nil
	}
	return MutationResult{Passed: false, MutantsTotal: 1, MutantsSurvived: 1}, nil
}

// RollbackWorktree executa rollback atômico do worktree isolado em falha,
// removendo alterações sem contaminar a branch principal em workdir.
func RollbackWorktree(ctx context.Context, workdir, isoDir string) error {
	if isoDir == "" {
		return nil
	}
	if isGitRepo(workdir) {
		cleanCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		_ = exec.CommandContext(cleanCtx, "git", "-C", isoDir, "reset", "--hard", "HEAD").Run()
		_ = exec.CommandContext(cleanCtx, "git", "-C", isoDir, "clean", "-fd").Run()
		_ = exec.CommandContext(cleanCtx, "git", "-C", workdir, "worktree", "remove", "--force", isoDir).Run()
	}
	if err := os.RemoveAll(isoDir); err != nil {
		return fmt.Errorf("falha ao remover baia isolada: %w", err)
	}
	return nil
}

func resolveSafePath(baseDir, relPath string) (string, error) {
	clean := filepath.Clean(relPath)
	if filepath.IsAbs(clean) || strings.HasPrefix(clean, "..") {
		return "", fmt.Errorf("caminho fora do Path Jail: %s", relPath)
	}
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return "", err
	}
	absTarget := filepath.Join(absBase, clean)
	if !strings.HasPrefix(absTarget, absBase+string(filepath.Separator)) && absTarget != absBase {
		return "", fmt.Errorf("escape de Path Jail detectado: %s", relPath)
	}
	return absTarget, nil
}

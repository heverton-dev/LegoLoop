package runner

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"legoloop/internal/dag"
)

// prepareIsolation provisiona um ambiente hermético e efêmero para execução do nó:
// se o diretório do projeto for um repositório git, usa `git worktree add --detach`;
// caso contrário, cria uma cópia temporária (R-NODE-8).
func prepareIsolation(ctx context.Context, workdir string, n dag.Node) (string, func(), error) {
	if isGitRepo(workdir) {
		isolatedDir, cleanup, err := createGitWorktree(ctx, workdir)
		if err == nil {
			return isolatedDir, cleanup, nil
		}
		// Se git worktree falhou (ex: repo recém-iniciado sem commits), faz fallback para temp copy
	}

	return createTempCopy(workdir)
}

func isGitRepo(dir string) bool {
	gitDir := filepath.Join(dir, ".git")
	info, err := os.Stat(gitDir)
	return err == nil && (info.IsDir() || !info.IsDir()) // .git pode ser dir ou arquivo (submódulo/worktree)
}

func createGitWorktree(ctx context.Context, workdir string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "legoloop-wt-*")
	if err != nil {
		return "", nil, err
	}

	cmdCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, "git", "-C", workdir, "worktree", "add", "--detach", tmpDir, "HEAD")
	if out, err := cmd.CombinedOutput(); err != nil {
		_ = os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("git worktree add falhou: %s (%w)", strings.TrimSpace(string(out)), err)
	}

	cleanup := func() {
		cleanCtx, cleanCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanCancel()
		_ = exec.CommandContext(cleanCtx, "git", "-C", workdir, "worktree", "remove", "--force", tmpDir).Run()
		_ = os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup, nil
}

func createTempCopy(workdir string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "legoloop-iso-*")
	if err != nil {
		return "", nil, err
	}

	// Copia arquivos do workdir para tmpDir, pulando diretórios pesados e .git
	err = filepath.WalkDir(workdir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(workdir, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		// Pula diretórios ignorados
		name := d.Name()
		if d.IsDir() && (name == ".git" || name == "node_modules" || name == "bin" || name == ".gemini" || name == ".claude") {
			return filepath.SkipDir
		}

		target := filepath.Join(tmpDir, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		return copyFile(path, target)
	})

	if err != nil {
		_ = os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("falha ao copiar workspace para ambiente isolado: %w", err)
	}

	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

// promoteArtifact copia única e exclusivamente o artefato declarado do ambiente isolado
// para o workspace principal do projeto (R-NODE-8).
func promoteArtifact(isolatedDir, mainWorkdir, artifactRel string) error {
	if artifactRel == "" {
		return nil
	}

	cleanRel := filepath.Clean(artifactRel)
	if strings.HasPrefix(cleanRel, "..") || filepath.IsAbs(cleanRel) {
		return fmt.Errorf("caminho de artefato inválido: %q", artifactRel)
	}

	src := filepath.Join(isolatedDir, cleanRel)
	dst := filepath.Join(mainWorkdir, cleanRel)

	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("artefato declarado %q não foi gerado no ambiente isolado: %w", artifactRel, err)
	}
	if info.IsDir() {
		return fmt.Errorf("artefato declarado %q é um diretório, esperado arquivo", artifactRel)
	}

	return copyFile(src, dst)
}

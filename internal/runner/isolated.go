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
// caso contrário, cria uma cópia temporária (R-NODE-8). Em seguida, injeta context_in.
func prepareIsolation(ctx context.Context, workdir string, n dag.Node) (string, func(), error) {
	isoDir, cleanup, err := createIsolationEnv(ctx, workdir)
	if err != nil {
		return "", nil, err
	}
	if err := injectContextIn(workdir, isoDir, n.ContextIn); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("falha ao injetar context_in: %w", err)
	}
	return isoDir, cleanup, nil
}

func createIsolationEnv(ctx context.Context, workdir string) (string, func(), error) {
	if isGitRepo(workdir) {
		if isolatedDir, cleanup, err := createGitWorktree(ctx, workdir); err == nil {
			return isolatedDir, cleanup, nil
		}
	}
	return createTempCopy(workdir)
}

func isGitRepo(dir string) bool {
	gitDir := filepath.Join(dir, ".git")
	info, err := os.Stat(gitDir)
	return err == nil && (info.IsDir() || !info.IsDir())
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
	if err := copyWorkspaceFiles(workdir, tmpDir); err != nil {
		_ = os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("falha ao copiar workspace para ambiente isolado: %w", err)
	}
	return tmpDir, func() { _ = os.RemoveAll(tmpDir) }, nil
}

func copyWorkspaceFiles(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil || rel == "." {
			return err
		}
		if d.IsDir() && isIgnoredDir(d.Name()) {
			return filepath.SkipDir
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target)
	})
}

func isIgnoredDir(name string) bool {
	return name == ".git" || name == "node_modules" || name == "bin" || name == ".gemini" || name == ".claude"
}

func injectContextIn(srcWorkdir, dstWorkdir string, contextIn []string) error {
	for _, inRel := range contextIn {
		clean := filepath.Clean(filepath.FromSlash(strings.TrimSpace(inRel)))
		if clean == "" || !filepath.IsLocal(clean) {
			return fmt.Errorf("path jail: context_in %q fora da raiz permitida", inRel)
		}
		src := filepath.Join(srcWorkdir, clean)
		dst := filepath.Join(dstWorkdir, clean)
		if err := copyPath(src, dst); err != nil {
			return err
		}
	}
	return nil
}

func copyPath(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return nil
	}
	if info.IsDir() {
		return copyDir(src, dst)
	}
	return copyFile(src, dst)
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, p)
		if err != nil || rel == "." {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(p, target)
	})
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

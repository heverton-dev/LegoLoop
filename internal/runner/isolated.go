package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/dag"
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
	info, err := os.Lstat(workdir)
	if err != nil {
		return "", nil, fmt.Errorf("workdir inválido: %w", err)
	}
	if isLink(info.Mode()) {
		return "", nil, fmt.Errorf("symlink ou junction não permitido no workdir: %s", workdir)
	}
	gitOK, err := checkGitRepo(workdir)
	if err != nil {
		return "", nil, err
	}
	if gitOK {
		if isolatedDir, cleanup, err := createGitWorktree(ctx, workdir); err == nil {
			return isolatedDir, cleanup, nil
		}
	}
	return createTempCopy(workdir)
}

func isGitRepo(dir string) bool {
	ok, err := checkGitRepo(dir)
	return err == nil && ok
}

func checkGitRepo(dir string) (bool, error) {
	gitDir := filepath.Join(dir, ".git")
	info, err := os.Lstat(gitDir)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if isLink(info.Mode()) {
		return false, fmt.Errorf("symlink ou junction não permitido em .git: %s", gitDir)
	}
	return true, nil
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
	return copyTree(src, dst)
}

func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(path string, _ fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(src, path)
		if err != nil || rel == "." {
			return err
		}
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		if isLink(info.Mode()) {
			return fmt.Errorf("symlink ou junction não permitido: %s", rel)
		}
		name := info.Name()
		if info.IsDir() {
			if isIgnoredDir(name) || isSensitiveDir(name) {
				return filepath.SkipDir
			}
			return os.MkdirAll(filepath.Join(dst, rel), 0o755)
		}
		if isSensitiveFile(name) {
			return nil
		}
		return copyFile(path, filepath.Join(dst, rel))
	})
}

func isLink(m fs.FileMode) bool {
	return m&(os.ModeSymlink|os.ModeIrregular|os.ModeDevice|os.ModeNamedPipe|os.ModeSocket) != 0
}

func isIgnoredDir(name string) bool {
	return name == ".git" || name == "node_modules" || name == "bin" || name == ".gemini" || name == ".claude"
}

func isSensitiveDir(name string) bool {
	switch strings.ToLower(name) {
	case ".git", ".claude", ".gemini", ".vscode", ".idea", ".husky", ".cursor", ".agents", ".ssh", ".aws", ".gnupg", ".kube", ".docker":
		return true
	}
	return false
}

func isSensitiveFile(name string) bool {
	lower := strings.ToLower(name)
	if strings.HasPrefix(lower, ".env") || strings.HasPrefix(lower, "credential_") || strings.HasPrefix(lower, "credentials") {
		return true
	}
	if strings.HasPrefix(lower, "secret") || strings.HasPrefix(lower, ".secret") {
		return true
	}
	if lower == ".npmrc" || lower == ".netrc" || lower == ".pypirc" || lower == ".htpasswd" {
		return true
	}
	if strings.HasPrefix(lower, "id_rsa") || strings.HasPrefix(lower, "id_ed25519") || strings.HasPrefix(lower, "id_ecdsa") {
		return true
	}
	switch filepath.Ext(lower) {
	case ".pem", ".key", ".p12", ".pfx", ".jks", ".kdbx":
		return true
	}
	return false
}

func isSensitiveRel(rel string) bool {
	for part := range strings.SplitSeq(filepath.ToSlash(filepath.Clean(rel)), "/") {
		if isSensitiveDir(part) || isSensitiveFile(part) {
			return true
		}
	}
	return false
}

func validateLocalRel(raw, label string) (string, error) {
	clean := filepath.Clean(filepath.FromSlash(strings.TrimSpace(raw)))
	if clean == "" || clean == "." || strings.Contains(raw, ":") || !filepath.IsLocal(clean) {
		return "", fmt.Errorf("path jail: %s %q fora da raiz permitida", label, raw)
	}
	if isSensitiveRel(clean) {
		return "", fmt.Errorf("credencial ou caminho sensível bloqueado em %s: %q", label, raw)
	}
	return clean, nil
}

func checkNoLinksRoot(root *os.Root, rel string, mustExist bool) error {
	cur := ""
	for part := range strings.SplitSeq(rel, string(filepath.Separator)) {
		if cur == "" {
			cur = part
		} else {
			cur = filepath.Join(cur, part)
		}
		info, err := root.Lstat(cur)
		if errors.Is(err, os.ErrNotExist) && !mustExist {
			return nil
		}
		if err != nil {
			return err
		}
		if isLink(info.Mode()) {
			return fmt.Errorf("symlink ou junction não permitido: %q", filepath.ToSlash(cur))
		}
	}
	return nil
}

func checkNoLinks(baseDir, rel string, mustExist bool) error {
	root, err := os.OpenRoot(baseDir)
	if err != nil {
		return err
	}
	defer root.Close()
	return checkNoLinksRoot(root, rel, mustExist)
}

func injectContextIn(srcWorkdir, dstWorkdir string, contextIn []string) error {
	for _, inRel := range contextIn {
		clean, err := validateLocalRel(inRel, "context_in")
		if err != nil {
			return err
		}
		if err := checkNoLinks(srcWorkdir, clean, true); err != nil {
			return err
		}
		if err := checkNoLinks(dstWorkdir, clean, false); err != nil {
			return err
		}
		if err := copyPath(filepath.Join(srcWorkdir, clean), filepath.Join(dstWorkdir, clean)); err != nil {
			return err
		}
	}
	return nil
}

func copyPath(src, dst string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	if isLink(info.Mode()) {
		return fmt.Errorf("symlink ou junction não permitido: %s", src)
	}
	if info.IsDir() {
		return copyTree(src, dst)
	}
	return copyFile(src, dst)
}

func copyFile(src, dst string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	if isLink(info.Mode()) || !info.Mode().IsRegular() {
		return fmt.Errorf("origem não é arquivo regular (symlink/especial rejeitado): %s", src)
	}
	if dinfo, lerr := os.Lstat(dst); lerr == nil && isLink(dinfo.Mode()) {
		return fmt.Errorf("symlink ou junction não permitido no destino: %s", dst)
	}
	return copyRegularBytes(src, dst)
}

func copyRegularBytes(src, dst string) error {
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
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}

// promoteArtifact copia única e exclusivamente o artefato declarado do ambiente isolado
// para o workspace principal do projeto (R-NODE-8), validando Path Jail, symlinks e credenciais.
func promoteArtifact(isolatedDir, mainWorkdir, artifactRel string) error {
	if artifactRel == "" {
		return nil
	}
	cleanRel, err := validateLocalRel(artifactRel, "artefato")
	if err != nil {
		return err
	}
	srcRoot, err := os.OpenRoot(isolatedDir)
	if err != nil {
		return fmt.Errorf("falha ao abrir origem isolada: %w", err)
	}
	defer srcRoot.Close()
	dstRoot, err := os.OpenRoot(mainWorkdir)
	if err != nil {
		return fmt.Errorf("falha ao abrir workspace principal: %w", err)
	}
	defer dstRoot.Close()
	return copyArtifactRoot(srcRoot, dstRoot, cleanRel, artifactRel)
}

func copyArtifactRoot(srcRoot, dstRoot *os.Root, cleanRel, artifactRel string) error {
	if err := checkNoLinksRoot(srcRoot, cleanRel, true); err != nil {
		return fmt.Errorf("artefato %q inválido na origem: %w", artifactRel, err)
	}
	info, err := srcRoot.Lstat(cleanRel)
	if err != nil {
		return fmt.Errorf("artefato declarado %q não foi gerado no ambiente isolado: %w", artifactRel, err)
	}
	if info.IsDir() {
		return fmt.Errorf("artefato declarado %q é um diretório, esperado arquivo", artifactRel)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("artefato declarado %q não é arquivo regular", artifactRel)
	}
	if err := checkNoLinksRoot(dstRoot, cleanRel, false); err != nil {
		return fmt.Errorf("destino de artefato %q inválido: %w", artifactRel, err)
	}
	return transferRootFile(srcRoot, dstRoot, cleanRel)
}

func transferRootFile(srcRoot, dstRoot *os.Root, cleanRel string) error {
	if parent := filepath.Dir(cleanRel); parent != "." {
		if err := dstRoot.MkdirAll(parent, 0o755); err != nil {
			return err
		}
	}
	in, err := srcRoot.Open(cleanRel)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := dstRoot.Create(cleanRel)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}

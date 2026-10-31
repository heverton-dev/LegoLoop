package canon

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// backupSuffix marca o original preservado durante a troca por link.
const backupSuffix = ".legoloop-bak"

// Method é a técnica de link efetivamente aplicada.
type Method string

const (
	MethodSymlink  Method = "symlink"
	MethodHardlink Method = "hardlink"
)

// Linker abstrai as chamadas de SO. No Windows, os.Symlink ≡ `mklink` (exige Developer
// Mode/admin) e os.Link ≡ `mklink /H` — APIs nativas, sem shell, imunes a injeção.
type Linker struct {
	Symlink  func(oldname, newname string) error
	Hardlink func(oldname, newname string) error
}

// DefaultLinker usa as syscalls nativas.
func DefaultLinker() Linker { return Linker{Symlink: os.Symlink, Hardlink: os.Link} }

// DefaultVault é ~/.aidd/canonical.
func DefaultVault() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".aidd", "canonical"), nil
}

// LinkResult registra um arquivo substituído por link.
type LinkResult struct {
	Path   string `json:"path"`
	Method Method `json:"method"`
}

// Failure registra um arquivo preservado (rollback) e o motivo.
type Failure struct {
	Path  string `json:"path"`
	Error string `json:"error"`
}

// ApplyReport é o resultado do --apply.
type ApplyReport struct {
	Linked     []LinkResult `json:"linked"`
	Failed     []Failure    `json:"failed"`
	SavedBytes int64        `json:"saved_bytes"`
}

// Apply executa o plano grupo a grupo; qualquer falha restaura o original.
func Apply(ctx context.Context, p Plan, l Linker) ApplyReport {
	var rep ApplyReport
	for _, g := range p.Groups {
		applyGroup(ctx, g, l, &rep)
	}
	return rep
}

func applyGroup(ctx context.Context, g GroupPlan, l Linker, rep *ApplyReport) {
	created, err := ensureVault(g)
	if err != nil {
		for _, path := range g.Links {
			rep.Failed = append(rep.Failed, Failure{Path: path, Error: err.Error()})
		}
		return
	}
	linked := 0
	for _, path := range g.Links {
		m, err := replaceWithLink(ctx, g.VaultPath, path, g.SHA256, l)
		if err != nil {
			rep.Failed = append(rep.Failed, Failure{Path: path, Error: err.Error()})
			continue
		}
		rep.Linked = append(rep.Linked, LinkResult{Path: path, Method: m})
		linked++
	}
	saved := g.Canonical.Size * int64(linked)
	if created {
		saved -= g.Canonical.Size
	}
	if created && linked == 0 {
		// Nenhum link aproveitou a cópia: remove do cofre; se falhar, fica um órfão inofensivo.
		_ = os.Remove(g.VaultPath)
	}
	rep.SavedBytes += max(saved, 0)
}

// ensureVault copia o canônico para o cofre (via arquivo temporário + rename) e valida o hash.
func ensureVault(g GroupPlan) (bool, error) {
	if h, err := hashOf(g.VaultPath); err == nil {
		if h != g.SHA256 {
			return false, fmt.Errorf("cofre corrompido em %s", g.VaultPath)
		}
		return false, nil
	}
	if err := os.MkdirAll(filepath.Dir(g.VaultPath), 0o755); err != nil {
		return false, err
	}
	tmp := g.VaultPath + ".tmp"
	if err := copyFile(g.Canonical.Path, tmp); err != nil {
		return false, errors.Join(err, removeIfExists(tmp))
	}
	if h, err := hashOf(tmp); err != nil || h != g.SHA256 {
		return false, errors.Join(fmt.Errorf("canônico %s mudou desde o plano", g.Canonical.Path), removeIfExists(tmp))
	}
	return true, os.Rename(tmp, g.VaultPath)
}

// replaceWithLink: valida hash (anti-TOCTOU) → renomeia para backup → link (symlink,
// fallback hardlink) → revalida link E backup → remove backup. Qualquer erro ⇒ rollback.
func replaceWithLink(ctx context.Context, vaultPath, path, sha string, l Linker) (Method, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if err := checkRegularWithHash(path, sha); err != nil {
		return "", err
	}
	backup := path + backupSuffix
	if _, err := os.Lstat(backup); err == nil {
		return "", fmt.Errorf("backup %s já existe", backup)
	}
	if err := os.Rename(path, backup); err != nil {
		return "", err
	}
	m, err := link(vaultPath, path, l)
	if err == nil {
		err = verifySwap(path, backup, vaultPath, sha)
	}
	if err != nil {
		return "", errors.Join(err, restore(backup, path))
	}
	return m, os.Remove(backup)
}

// restore devolve o backup, exceto se algo recriou o caminho durante a troca
// (salvamento atômico de editor): nesse caso nada é sobrescrito e o backup fica.
func restore(backup, path string) error {
	if _, err := os.Lstat(path); err == nil {
		return fmt.Errorf("%s recriado durante a troca; original preservado em %s", path, backup)
	}
	return os.Rename(backup, path)
}

// verifySwap confere que o link reproduz o original e que ninguém editou o original
// (agora no backup) durante a troca; se falhar, remove o link — mas só se ainda for
// o NOSSO link (um editor pode tê-lo trocado por um arquivo novo).
func verifySwap(path, backup, vaultPath, sha string) error {
	var err error
	if h, herr := hashOf(path); herr != nil || h != sha {
		err = errors.New("link criado não reproduz o conteúdo original")
	} else if h, herr := hashOf(backup); herr != nil || h != sha {
		err = errors.New("arquivo editado durante a troca; edição preservada")
	}
	if err == nil {
		return nil
	}
	if !isOurLink(path, vaultPath) {
		return errors.Join(err, fmt.Errorf("%s foi substituído durante a troca; mantido", path))
	}
	return errors.Join(err, os.Remove(path))
}

// isOurLink: symlink apontando para o cofre ou hardlink do mesmo arquivo físico.
func isOurLink(path, vaultPath string) bool {
	st, err := os.Lstat(path)
	if err != nil {
		return false
	}
	if st.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(path)
		return err == nil && target == vaultPath
	}
	vi, err := os.Stat(vaultPath)
	return err == nil && os.SameFile(st, vi)
}

func link(target, path string, l Linker) (Method, error) {
	symErr := l.Symlink(target, path)
	if symErr == nil {
		return MethodSymlink, nil
	}
	hardErr := l.Hardlink(target, path)
	if hardErr == nil {
		return MethodHardlink, nil
	}
	return "", fmt.Errorf("symlink: %v; hardlink: %w", symErr, hardErr)
}

func checkRegularWithHash(path, sha string) error {
	st, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if !st.Mode().IsRegular() {
		return fmt.Errorf("%s não é arquivo regular", path)
	}
	h, err := hashOf(path)
	if err != nil {
		return err
	}
	if h != sha {
		return fmt.Errorf("%s mudou desde o plano; preservado", path)
	}
	return nil
}

func hashOf(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, in)
	return errors.Join(err, out.Close())
}

func removeIfExists(path string) error {
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

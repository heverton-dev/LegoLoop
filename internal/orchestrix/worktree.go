package orchestrix

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// maxApplyBytes limita o tamanho de um arquivo copiado da worktree para o principal.
const maxApplyBytes = 16 << 20

// Worktree é uma git worktree efêmera, criada a partir de HEAD fora do repositório principal.
type Worktree struct {
	Repo string // repositório principal
	Dir  string // diretório da worktree
	base string // diretório temporário que contém Dir
}

// gitRun executa git sem hooks e sem configuração que altere a leitura de caminhos.
func gitRun(ctx context.Context, dir string, args ...string) ([]byte, error) {
	full := append([]string{"-c", "core.hooksPath=" + os.DevNull, "-c", "core.quotepath=false", "-c", "core.fsmonitor=false"}, args...)
	cmd := exec.CommandContext(ctx, "git", full...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("git %s: %w: %s", args[0], err, strings.TrimSpace(stderr.String()))
	}
	return stdout.Bytes(), nil
}

// CreateWorktree cria uma worktree destacada a partir de HEAD de repo.
func CreateWorktree(ctx context.Context, repo string) (*Worktree, error) {
	abs, err := filepath.Abs(repo)
	if err != nil {
		return nil, err
	}
	base, err := os.MkdirTemp("", "legoloop-wt-")
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(base, "wt")
	if _, err := gitRun(ctx, abs, "worktree", "add", "--detach", dir, "HEAD"); err != nil {
		os.RemoveAll(base)
		return nil, err
	}
	return &Worktree{Repo: abs, Dir: dir, base: base}, nil
}

// RunInWorktree cria a worktree, executa fn e descarta a worktree sempre, inclusive em erro.
func RunInWorktree(ctx context.Context, repo string, fn func(*Worktree) error) (err error) {
	w, err := CreateWorktree(ctx, repo)
	if err != nil {
		return err
	}
	defer func() {
		// O contexto de fn pode já estar cancelado: o descarte usa contexto próprio.
		err = errors.Join(err, w.Discard(context.WithoutCancel(ctx)))
	}()
	return fn(w)
}

// Changed lista os caminhos alterados, criados, removidos ou renomeados na worktree
// (git status --porcelain -z, com arquivos novos expandidos). Caminhos usam "/".
func (w *Worktree) Changed(ctx context.Context) ([]string, error) {
	out, err := gitRun(ctx, w.Dir, "status", "--porcelain", "-z", "--untracked-files=all")
	if err != nil {
		return nil, err
	}
	return parsePorcelainZ(out)
}

// IgnoredAmong lista os arquivos ignorados pelo git que existem na worktree e casam com
// os caminhos permitidos. O status comum não os mostra. Só caminhos literais, nunca padrão.
func (w *Worktree) IgnoredAmong(ctx context.Context, allowed []string) ([]string, error) {
	args := []string{"ls-files", "--others", "--ignored", "--exclude-standard", "-z", "--"}
	for _, a := range allowed {
		rel, err := normRel(a)
		if err != nil {
			return nil, err
		}
		args = append(args, ":(literal)"+rel)
	}
	if len(args) == 6 {
		return nil, nil
	}
	out, err := gitRun(ctx, w.Dir, args...)
	if err != nil {
		return nil, err
	}
	var paths []string
	for p := range strings.SplitSeq(string(out), "\x00") {
		if p != "" {
			paths = append(paths, p)
		}
	}
	return paths, nil
}

// parsePorcelainZ lê a saída de `git status --porcelain -z`.
// Cada registro é "XY caminho\0"; renomeação e cópia acrescentam "origem\0".
func parsePorcelainZ(out []byte) ([]string, error) {
	var paths []string
	parts := strings.Split(string(out), "\x00")
	for i := 0; i < len(parts); i++ {
		rec := parts[i]
		if rec == "" {
			continue
		}
		if len(rec) < 4 || rec[2] != ' ' {
			return nil, fmt.Errorf("registro de status inválido: %q", rec)
		}
		paths = append(paths, rec[3:])
		if rec[0] == 'R' || rec[0] == 'C' || rec[1] == 'R' || rec[1] == 'C' {
			i++
			if i >= len(parts) || parts[i] == "" {
				return nil, fmt.Errorf("renomeação sem origem: %q", rec)
			}
			paths = append(paths, parts[i])
		}
	}
	return paths, nil
}

// isLink informa se o modo é link simbólico, junction ou ponto de montagem.
func isLink(m os.FileMode) bool {
	return m&(os.ModeSymlink|os.ModeIrregular|os.ModeDevice|os.ModeNamedPipe|os.ModeSocket) != 0
}

// checkNoLinks recusa o caminho se ele ou qualquer diretório-pai for link simbólico ou
// junction em root. Componente ausente encerra a verificação sem erro.
func checkNoLinks(root *os.Root, rel string) error {
	cur := ""
	for part := range strings.SplitSeq(filepath.FromSlash(rel), string(filepath.Separator)) {
		cur = filepath.Join(cur, part)
		info, err := root.Lstat(cur)
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		if err != nil {
			return err
		}
		if isLink(info.Mode()) {
			return fmt.Errorf("link simbólico ou junction em caminho alterado: %q", filepath.ToSlash(cur))
		}
	}
	return nil
}

// CheckLinks reprova qualquer caminho alterado que atravesse link simbólico ou junction.
func (w *Worktree) CheckLinks(changed []string) error {
	root, err := os.OpenRoot(w.Dir)
	if err != nil {
		return err
	}
	defer root.Close()
	var errs []error
	for _, c := range changed {
		rel, err := normRel(c)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if err := checkNoLinks(root, rel); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// Apply aplica no repositório principal os caminhos alterados da worktree:
// copia arquivo criado ou editado e remove arquivo apagado. Valida tudo antes de escrever.
func (w *Worktree) Apply(changed []string) error {
	src, err := os.OpenRoot(w.Dir)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.OpenRoot(w.Repo)
	if err != nil {
		return err
	}
	defer dst.Close()

	ops, err := planApply(src, dst, changed)
	if err != nil {
		return err
	}
	for _, o := range ops {
		if o.remove {
			if err := dst.Remove(filepath.FromSlash(o.rel)); err != nil && !errors.Is(err, os.ErrNotExist) {
				return err
			}
			continue
		}
		if err := writeArtifact(dst, o.rel, o.data); err != nil {
			return err
		}
	}
	return nil
}

// applyOp é uma escrita ou remoção já validada no repositório principal.
type applyOp struct {
	rel    string
	remove bool
	data   []byte
}

// planApply valida todos os caminhos e lê o conteúdo antes de qualquer escrita.
func planApply(src, dst *os.Root, changed []string) ([]applyOp, error) {
	var ops []applyOp
	for _, c := range changed {
		rel, err := normRel(c)
		if err != nil {
			return nil, err
		}
		if err := checkNoLinks(src, rel); err != nil {
			return nil, err
		}
		if err := checkNoLinks(dst, rel); err != nil {
			return nil, fmt.Errorf("destino: %w", err)
		}
		native := filepath.FromSlash(rel)
		info, err := src.Lstat(native)
		if errors.Is(err, os.ErrNotExist) {
			ops = append(ops, applyOp{rel: rel, remove: true})
			continue
		}
		if err != nil {
			return nil, err
		}
		if !info.Mode().IsRegular() {
			return nil, fmt.Errorf("%q não é arquivo regular", rel)
		}
		data, err := readAll(src, native)
		if err != nil {
			return nil, err
		}
		ops = append(ops, applyOp{rel: rel, data: data})
	}
	return ops, nil
}

func readAll(root *os.Root, native string) ([]byte, error) {
	f, err := root.Open(native)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, maxApplyBytes+1))
	if err == nil && len(data) > maxApplyBytes {
		err = fmt.Errorf("%q excede %d bytes", native, maxApplyBytes)
	}
	return data, err
}

// Discard remove a worktree e o diretório temporário. É idempotente.
func (w *Worktree) Discard(ctx context.Context) error {
	if w.base == "" {
		return nil
	}
	var errs []error
	// Links (symlink/junction) saem primeiro: nem git nem RemoveAll podem seguir um deles
	// e apagar o conteúdo do alvo fora da worktree.
	if err := unlinkAll(w.Dir); err != nil {
		errs = append(errs, err)
	}
	if _, err := gitRun(ctx, w.Repo, "worktree", "remove", "--force", w.Dir); err != nil {
		errs = append(errs, err)
	}
	if err := os.RemoveAll(w.base); err != nil {
		errs = append(errs, err)
	}
	if _, err := gitRun(ctx, w.Repo, "worktree", "prune"); err != nil {
		errs = append(errs, err)
	}
	// A remoção direta cobre falha do git: o estado final do disco é o que vale.
	if !dirExists(w.base) {
		w.base = ""
		return nil
	}
	return errors.Join(errs...)
}

func dirExists(p string) bool {
	_, err := os.Lstat(p)
	return err == nil
}

// unlinkAll remove, sem seguir, todo link simbólico e junction dentro de dir.
func unlinkAll(dir string) error {
	var errs []error
	var walk func(string)
	walk = func(d string) {
		entries, err := os.ReadDir(d)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				errs = append(errs, err)
			}
			return
		}
		for _, e := range entries {
			p := filepath.Join(d, e.Name())
			info, err := os.Lstat(p)
			if err != nil {
				continue
			}
			switch {
			case isLink(info.Mode()):
				if err := os.Remove(p); err != nil {
					errs = append(errs, err)
				}
			case info.IsDir():
				walk(p)
			}
		}
	}
	walk(dir)
	return errors.Join(errs...)
}

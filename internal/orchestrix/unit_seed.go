package orchestrix

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"legoloop/internal/dag"
)

// seedContext copia para a worktree os arquivos de context_in que existem no repositório
// principal. A worktree nasce de HEAD e não contém artefatos de nós anteriores ainda não
// commitados. Devolve o conteúdo semeado por caminho, para excluir os semeados de Changed.
// Arquivo ausente é ignorado. Link, arquivo irregular e ".git" reprovam.
func seedContext(w *Worktree, n dag.Node) (map[string][]byte, error) {
	src, err := os.OpenRoot(w.Repo)
	if err != nil {
		return nil, err
	}
	defer src.Close()
	dst, err := os.OpenRoot(w.Dir)
	if err != nil {
		return nil, err
	}
	defer dst.Close()
	seeded := map[string][]byte{}
	for _, f := range n.ContextIn {
		rel, data, err := readSeed(src, f)
		if errors.Is(err, errSeedDir) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("context_in %q: %w", f, err)
		}
		if data == nil {
			continue
		}
		if err := writeArtifact(dst, rel, data); err != nil {
			return nil, fmt.Errorf("context_in %q: %w", f, err)
		}
		seeded[rel] = data
	}
	return seeded, nil
}

// errSeedDir sinaliza que o context_in é uma pasta: a worktree já traz o conteúdo rastreado.
var errSeedDir = errors.New("context_in é pasta")

// readSeed valida um caminho de context_in e lê o arquivo; data nulo = arquivo ausente.
func readSeed(src *os.Root, f string) (string, []byte, error) {
	rel, err := normRel(f)
	if err != nil {
		return "", nil, err
	}
	if first, _, _ := strings.Cut(rel, "/"); strings.EqualFold(first, ".git") {
		return "", nil, errors.New("caminho dentro de .git")
	}
	if err := checkNoLinks(src, rel); err != nil {
		return "", nil, err
	}
	native := filepath.FromSlash(rel)
	info, err := src.Lstat(native)
	if errors.Is(err, os.ErrNotExist) {
		return rel, nil, nil
	}
	if err != nil {
		return "", nil, err
	}
	if info.IsDir() {
		return rel, nil, errSeedDir
	}
	if !info.Mode().IsRegular() {
		return "", nil, errors.New("não é arquivo regular")
	}
	data, err := readAll(src, native)
	if data == nil {
		data = []byte{}
	}
	return rel, data, err
}

// unitChanged lista os alterados da worktree sem os arquivos semeados e inalterados,
// e aplica a checagem de links e o write_boundary.
func unitChanged(ctx context.Context, n dag.Node, w *Worktree, seeded map[string][]byte) ([]string, error) {
	all, err := w.Changed(ctx)
	if err != nil {
		return nil, err
	}
	ignored, err := w.IgnoredAmong(ctx, allowedWrites(n))
	if err != nil {
		return nil, err
	}
	var changed []string
	for _, c := range dedupe(append(all, ignored...)) {
		if !seededUnchanged(w.Dir, c, seeded) {
			changed = append(changed, c)
		}
	}
	sort.Strings(changed)
	return changed, errors.Join(w.CheckLinks(changed), checkBoundary(changed, allowedWrites(n)))
}

// dedupe remove caminhos repetidos e mantém a ordem.
func dedupe(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, p := range in {
		if !seen[p] {
			seen[p] = true
			out = append(out, p)
		}
	}
	return out
}

// seededUnchanged informa se o caminho foi semeado e o conteúdo na worktree continua idêntico.
func seededUnchanged(dir, changed string, seeded map[string][]byte) bool {
	rel, err := normRel(changed)
	if err != nil {
		return false
	}
	want, ok := seeded[rel]
	if !ok {
		return false
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return false
	}
	defer root.Close()
	if checkNoLinks(root, rel) != nil {
		return false
	}
	info, err := root.Lstat(filepath.FromSlash(rel))
	if err != nil || !info.Mode().IsRegular() {
		return false
	}
	got, err := readAll(root, filepath.FromSlash(rel))
	return err == nil && bytes.Equal(got, want)
}

package processo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	tmpSuffix = ".legoloop-tmp"
	bakSuffix = ".legoloop-bak"
)

// ErrExiste: já existe um arquivo do processo e Force=false.
var ErrExiste = errors.New("processo: já existe")

// ErrPlano: o plano de gravação é inválido (vazio, duplicado ou fora das raízes permitidas).
var ErrPlano = errors.New("processo: plano de gravação inválido")

// Roots são as raízes de gravação; em teste, diretórios temporários.
type Roots struct {
	Vault string // ~/.aidd/canonical
	Docs  string // raiz do projeto: contém docs/intents, docs/specs, docs/plans, docs/routes, docs/templates
}

// FS é a borda de gravação, injetável para provar atomicidade.
type FS interface {
	MkdirAll(path string, perm os.FileMode) error
	WriteFile(name string, data []byte, perm os.FileMode) error
	Rename(oldpath, newpath string) error
	RemoveAll(path string) error
	Stat(name string) (os.FileInfo, error)
}

type osFS struct{}

// OSFS devolve o FS real, sobre o disco.
func OSFS() FS { return osFS{} }

func (osFS) MkdirAll(p string, perm os.FileMode) error { return os.MkdirAll(p, perm) }
func (osFS) WriteFile(n string, d []byte, perm os.FileMode) error {
	return os.WriteFile(n, d, perm)
}
func (osFS) Rename(a, b string) error           { return os.Rename(a, b) }
func (osFS) RemoveAll(p string) error           { return os.RemoveAll(p) }
func (osFS) Stat(n string) (os.FileInfo, error) { return os.Stat(n) }

// File é um arquivo a gravar. Kind é livre: intent, spec, plan, route, template, skill, loop, dag...
type File struct {
	Kind string
	Path string
	Data []byte
}

// WriteOpts controla a gravação.
type WriteOpts struct {
	Force  bool // sobrescreve após preservar o original como .legoloop-bak
	DryRun bool // valida e antecipa conflitos; não grava nada
	// Depois roda com tudo já promovido (ex.: criar os vínculos aos harnesses). Se falhar, a gravação inteira é desfeita.
	// A própria Depois desfaz o que ela criou antes de devolver o erro.
	Depois func() error
}

// Store é o acesso ao cofre e ao projeto. É a única porta de gravação das funções de processo.
type Store struct {
	Roots Roots
	FS    FS // nil = disco real
}

func (s Store) fs() FS {
	if s.FS == nil {
		return OSFS()
	}
	return s.FS
}

// Write grava todos os arquivos de forma atômica: monta tudo em temporários no mesmo diretório e só então
// promove com Rename. Qualquer falha desfaz o que foi feito e deixa o estado anterior.
func (s Store) Write(ctx context.Context, files []File, opts WriteOpts) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := s.validate(files); err != nil {
		return err
	}
	fsys := s.fs()
	if found := Existentes(fsys, files); len(found) > 0 && !opts.Force {
		return fmt.Errorf("%w: %s", ErrExiste, strings.Join(found, ", "))
	}
	if opts.DryRun {
		return nil
	}
	tx := &txn{fsys: fsys, backups: map[string]string{}}
	err := tx.run(ctx, files)
	if err == nil && opts.Depois != nil {
		err = opts.Depois()
	}
	if err != nil {
		tx.rollback(files)
		return err
	}
	return nil
}

func (s Store) validate(files []File) error {
	if len(files) == 0 {
		return fmt.Errorf("%w: sem arquivos", ErrPlano)
	}
	seen := map[string]bool{}
	for _, f := range files {
		if f.Path == "" || (!inside(f.Path, s.Roots.Vault) && !inside(f.Path, s.Roots.Docs)) {
			return fmt.Errorf("%w: caminho fora das raízes permitidas: %q", ErrPlano, f.Path)
		}
		key := filepath.Clean(f.Path)
		if seen[key] {
			return fmt.Errorf("%w: caminho duplicado: %q", ErrPlano, f.Path)
		}
		seen[key] = true
	}
	return nil
}

// inside confere, depois de limpar o caminho, que ele está sob root (sem ".." que escape).
func inside(path, root string) bool {
	if root == "" {
		return false
	}
	rel, err := filepath.Rel(filepath.Clean(root), filepath.Clean(path))
	return err == nil && rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// Existentes devolve o nome dos arquivos do plano que já existem em disco.
func Existentes(fsys FS, files []File) []string {
	var out []string
	for _, f := range files {
		if _, err := fsys.Stat(f.Path); err == nil {
			out = append(out, filepath.Base(f.Path))
		}
	}
	return out
}

// txn guarda o que já foi feito, para desfazer na ordem inversa.
type txn struct {
	fsys     FS
	backups  map[string]string // destino -> backup criado
	promoted []string
	newDirs  []string // topo de cada árvore de diretórios criada por nós
}

func (tx *txn) run(ctx context.Context, files []File) error {
	for _, f := range files {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := tx.stage(f); err != nil {
			return err
		}
	}
	for _, f := range files {
		if err := tx.promote(f.Path); err != nil {
			return err
		}
	}
	return nil
}

// stage grava o temporário, registrando a pasta mais alta que precisou ser criada.
func (tx *txn) stage(f File) error {
	dir := filepath.Dir(f.Path)
	if top := tx.firstMissing(dir); top != "" {
		tx.newDirs = append(tx.newDirs, top)
	}
	if err := tx.fsys.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return tx.fsys.WriteFile(f.Path+tmpSuffix, f.Data, 0o644)
}

func (tx *txn) firstMissing(dir string) string {
	top := ""
	for d := filepath.Clean(dir); ; d = filepath.Dir(d) {
		if _, err := tx.fsys.Stat(d); err == nil {
			return top
		}
		top = d
		if filepath.Dir(d) == d {
			return top
		}
	}
}

func (tx *txn) promote(path string) error {
	if _, err := tx.fsys.Stat(path); err == nil {
		bak := path + bakSuffix
		if err := tx.fsys.RemoveAll(bak); err != nil {
			return err
		}
		if err := tx.fsys.Rename(path, bak); err != nil {
			return err
		}
		tx.backups[path] = bak
	}
	if err := tx.fsys.Rename(path+tmpSuffix, path); err != nil {
		return err
	}
	tx.promoted = append(tx.promoted, path)
	return nil
}

// rollback desfaz em ordem inversa: remove o que foi promovido, devolve os originais, apaga os
// temporários e as pastas que criamos. Erros de limpeza não mascaram o erro original.
func (tx *txn) rollback(files []File) {
	for i := len(tx.promoted) - 1; i >= 0; i-- {
		_ = tx.fsys.RemoveAll(tx.promoted[i])
	}
	for path, bak := range tx.backups {
		_ = tx.fsys.RemoveAll(path)
		_ = tx.fsys.Rename(bak, path)
	}
	for _, f := range files {
		_ = tx.fsys.RemoveAll(f.Path + tmpSuffix)
	}
	for _, d := range tx.newDirs {
		_ = tx.fsys.RemoveAll(d)
	}
}

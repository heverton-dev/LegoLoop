package fabricar

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

// ErrExists: já existe um processo com o mesmo nome canônico e Force=false.
var ErrExists = errors.New("fabricar: processo já existe")

// Roots são as raízes de gravação; em teste, diretórios temporários.
type Roots struct {
	Vault string // ~/.aidd/canonical
	Docs  string // raiz que contém docs/intents, docs/specs, docs/plans, docs/routes, docs/templates
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

// PlannedFile é um arquivo a gravar. Kind: intent, spec, plan, route, template, skill, loop ou dag.
type PlannedFile struct {
	Kind string
	Path string
	Data []byte
}

// MaterializePlan descreve a gravação de um processo.
type MaterializePlan struct {
	Name   string
	Roots  Roots
	Files  []PlannedFile
	Force  bool // sobrescreve após preservar o original como .legoloop-bak
	DryRun bool // calcula tudo e não grava nada
}

// CanonicalArtifacts consolida os caminhos gerados.
type CanonicalArtifacts struct {
	IntentPath   string   `json:"intent_path"`
	SpecPath     string   `json:"spec_path"`
	PlanPath     string   `json:"plan_path"`
	RoutePath    string   `json:"route_path"`
	TemplatePath string   `json:"template_path"`
	SkillPath    string   `json:"skill_path"`
	LoopPath     string   `json:"loop_path"`
	DagPath      string   `json:"dag_path"`
	SkillLinks   []string `json:"skill_links"` // links nas pastas dos harnesses (vazio sem --link-harness)
}

// Materialize grava todos os arquivos do plano de forma atômica: monta tudo em temporários no mesmo
// diretório e só então promove com os.Rename. Qualquer falha desfaz o que foi feito e deixa o estado anterior.
func Materialize(ctx context.Context, plan MaterializePlan, fsys FS) (CanonicalArtifacts, error) {
	if err := ctx.Err(); err != nil {
		return CanonicalArtifacts{}, err
	}
	if err := validatePlan(plan); err != nil {
		return CanonicalArtifacts{}, err
	}
	art := artifactsOf(plan.Files)
	conflicts := existing(fsys, plan.Files)
	if len(conflicts) > 0 && !plan.Force {
		return art, fmt.Errorf("%w: %s", ErrExists, strings.Join(conflicts, ", "))
	}
	if plan.DryRun {
		return art, nil
	}
	tx := &txn{fsys: fsys, force: plan.Force, backups: map[string]string{}}
	if err := tx.run(ctx, plan.Files); err != nil {
		tx.rollback(plan.Files)
		return CanonicalArtifacts{}, err
	}
	return art, nil
}

func validatePlan(plan MaterializePlan) error {
	if len(plan.Files) == 0 {
		return errors.New("fabricar: plano sem arquivos")
	}
	seen := map[string]bool{}
	for _, f := range plan.Files {
		if f.Path == "" || (!inside(f.Path, plan.Roots.Vault) && !inside(f.Path, plan.Roots.Docs)) {
			return fmt.Errorf("fabricar: caminho fora das raízes permitidas: %q", f.Path)
		}
		key := filepath.Clean(f.Path)
		if seen[key] {
			return fmt.Errorf("fabricar: caminho duplicado no plano: %q", f.Path)
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

func existing(fsys FS, files []PlannedFile) []string {
	var out []string
	for _, f := range files {
		if _, err := fsys.Stat(f.Path); err == nil {
			out = append(out, filepath.Base(f.Path))
		}
	}
	return out
}

func artifactsOf(files []PlannedFile) CanonicalArtifacts {
	var a CanonicalArtifacts
	for _, f := range files {
		switch f.Kind {
		case "intent":
			a.IntentPath = f.Path
		case "spec":
			a.SpecPath = f.Path
		case "plan":
			a.PlanPath = f.Path
		case "route":
			a.RoutePath = f.Path
		case "template":
			a.TemplatePath = f.Path
		case "skill":
			a.SkillPath = f.Path
		case "loop":
			a.LoopPath = f.Path
		case "dag":
			a.DagPath = f.Path
		}
	}
	return a
}

// txn guarda o que já foi feito, para desfazer na ordem inversa.
type txn struct {
	fsys     FS
	force    bool
	backups  map[string]string // destino -> backup criado
	promoted []string
	newDirs  []string // topo de cada árvore de diretórios criada por nós
}

func (tx *txn) run(ctx context.Context, files []PlannedFile) error {
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
func (tx *txn) stage(f PlannedFile) error {
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
func (tx *txn) rollback(files []PlannedFile) {
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

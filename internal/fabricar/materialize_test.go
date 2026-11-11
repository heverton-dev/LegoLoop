package fabricar_test

import (
	"context"
	"errors"
	"io/fs"
	"legoloop/internal/paths"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/fabricar"
)

// failFS delega ao disco real e falha na N-ésima gravação ou renomeação (1-based; 0 = nunca).
type failFS struct {
	fabricar.FS
	failWriteAt, failRenameAt int
	writes, renames           int
}

var errInjected = errors.New("falha injetada")

func (f *failFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	f.writes++
	if f.writes == f.failWriteAt {
		return errInjected
	}
	return f.FS.WriteFile(name, data, perm)
}

func (f *failFS) Rename(oldpath, newpath string) error {
	f.renames++
	if f.renames == f.failRenameAt {
		return errInjected
	}
	return f.FS.Rename(oldpath, newpath)
}

type env struct {
	vault, docs string
	plan        fabricar.MaterializePlan
}

func newEnv(t *testing.T) env {
	t.Helper()
	vault, docs := t.TempDir(), t.TempDir()
	roots := fabricar.Roots{Vault: vault, Docs: docs}
	return env{vault: vault, docs: docs, plan: fabricar.MaterializePlan{Name: "legoloop-x", Roots: roots, Files: []fabricar.PlannedFile{
		{Kind: "skill", Path: filepath.Join(vault, "skills", "legoloop-x", "SKILL.md"), Data: []byte("skill")},
		{Kind: "loop", Path: filepath.Join(vault, paths.DirAntigo, "legoloop-x"+paths.SufixoAntigo), Data: []byte("loop")},
		{Kind: "dag", Path: filepath.Join(vault, paths.DirAntigo, "legoloop-x.dag.json"), Data: []byte("dag")},
		{Kind: "intent", Path: filepath.Join(docs, "docs", "intents", "2026-11-03_INTENT_X.md"), Data: []byte("intent")},
		{Kind: "route", Path: filepath.Join(docs, "docs", "routes", "2026-11-03_ROUTE_X.json"), Data: []byte("route")},
	}}}
}

// tree lista tudo que existe sob dir (arquivos e pastas), para provar que nada sobrou.
func tree(t *testing.T, dir string) []string {
	t.Helper()
	var out []string
	_ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err == nil && p != dir {
			out = append(out, p)
		}
		return nil
	})
	return out
}

func TestMaterialize_WritesEverythingAndFillsArtifacts(t *testing.T) {
	e := newEnv(t)
	art, err := fabricar.Materialize(context.Background(), e.plan, fabricar.OSFS())
	if err != nil {
		t.Fatal(err)
	}
	if art.SkillPath != e.plan.Files[0].Path || art.LoopPath != e.plan.Files[1].Path || art.DagPath != e.plan.Files[2].Path ||
		art.IntentPath != e.plan.Files[3].Path || art.RoutePath != e.plan.Files[4].Path {
		t.Errorf("artefatos incorretos: %+v", art)
	}
	for _, f := range e.plan.Files {
		if b, err := os.ReadFile(f.Path); err != nil || string(b) != string(f.Data) {
			t.Errorf("%s: %q %v", f.Path, b, err)
		}
	}
	for _, p := range append(tree(t, e.vault), tree(t, e.docs)...) {
		if strings.HasSuffix(p, ".legoloop-tmp") {
			t.Errorf("sobrou temporário: %s", p)
		}
	}
}

func TestMaterialize_RejectsPathsOutsideTheRoots(t *testing.T) {
	e := newEnv(t)
	e.plan.Files[1].Path = filepath.Join(e.vault, paths.DirAntigo, "..", "..", "fora.json")
	if _, err := fabricar.Materialize(context.Background(), e.plan, fabricar.OSFS()); err == nil {
		t.Fatal("caminho fora das raízes deve ser recusado")
	}
	if got := tree(t, e.vault); len(got) != 0 {
		t.Errorf("nada deve ser gravado quando o plano é inválido: %v", got)
	}
}

func TestMaterialize_FailureOnThirdWriteLeavesNothingBehind(t *testing.T) {
	e := newEnv(t)
	f := &failFS{FS: fabricar.OSFS(), failWriteAt: 3}
	if _, err := fabricar.Materialize(context.Background(), e.plan, f); !errors.Is(err, errInjected) {
		t.Fatalf("esperava a falha injetada: %v", err)
	}
	if v, d := tree(t, e.vault), tree(t, e.docs); len(v)+len(d) != 0 {
		t.Errorf("sobrou algo após a falha: %v %v", v, d)
	}
}

func TestMaterialize_ExistingWithoutForceReturnsErrExistsAndWritesNothing(t *testing.T) {
	e := newEnv(t)
	if _, err := fabricar.Materialize(context.Background(), e.plan, fabricar.OSFS()); err != nil {
		t.Fatal(err)
	}
	e.plan.Files[0].Data = []byte("OUTRO")
	_, err := fabricar.Materialize(context.Background(), e.plan, fabricar.OSFS())
	if !errors.Is(err, fabricar.ErrExists) {
		t.Fatalf("esperava ErrExists: %v", err)
	}
	if b, _ := os.ReadFile(e.plan.Files[0].Path); string(b) != "skill" {
		t.Errorf("o original deve ficar intacto: %q", b)
	}
}

func TestMaterialize_ForceKeepsTheOriginalAsBak(t *testing.T) {
	e := newEnv(t)
	if _, err := fabricar.Materialize(context.Background(), e.plan, fabricar.OSFS()); err != nil {
		t.Fatal(err)
	}
	e.plan.Files[0].Data, e.plan.Force = []byte("NOVO"), true
	if _, err := fabricar.Materialize(context.Background(), e.plan, fabricar.OSFS()); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(e.plan.Files[0].Path); string(b) != "NOVO" {
		t.Errorf("o novo conteúdo deve valer: %q", b)
	}
	if b, _ := os.ReadFile(e.plan.Files[0].Path + ".legoloop-bak"); string(b) != "skill" {
		t.Errorf("o original deve virar .legoloop-bak: %q", b)
	}
}

func TestMaterialize_DryRunWritesNothingButDescribesAndPredictsConflicts(t *testing.T) {
	e := newEnv(t)
	e.plan.DryRun = true
	art, err := fabricar.Materialize(context.Background(), e.plan, fabricar.OSFS())
	if err != nil || art.SkillPath == "" {
		t.Fatalf("dry-run deve descrever o que seria criado: %v %+v", err, art)
	}
	if v, d := tree(t, e.vault), tree(t, e.docs); len(v)+len(d) != 0 {
		t.Errorf("dry-run não grava: %v %v", v, d)
	}
	e.plan.DryRun = false
	if _, err := fabricar.Materialize(context.Background(), e.plan, fabricar.OSFS()); err != nil {
		t.Fatal(err)
	}
	e.plan.DryRun = true
	if _, err := fabricar.Materialize(context.Background(), e.plan, fabricar.OSFS()); !errors.Is(err, fabricar.ErrExists) {
		t.Errorf("dry-run deve antecipar ErrExists: %v", err)
	}
}

func TestMaterialize_RenameFailureRollsBackAndRestoresBackups(t *testing.T) {
	e := newEnv(t)
	if _, err := fabricar.Materialize(context.Background(), e.plan, fabricar.OSFS()); err != nil {
		t.Fatal(err)
	}
	e.plan.Force = true
	for i := range e.plan.Files {
		e.plan.Files[i].Data = []byte("v2")
	}
	// com Force cada arquivo faz 2 renomeações (backup e promoção): falha na 5ª, no meio do 3º arquivo
	f := &failFS{FS: fabricar.OSFS(), failRenameAt: 5}
	if _, err := fabricar.Materialize(context.Background(), e.plan, f); !errors.Is(err, errInjected) {
		t.Fatalf("esperava a falha injetada: %v", err)
	}
	for _, file := range e.plan.Files {
		if b, err := os.ReadFile(file.Path); err != nil || string(b) != originalOf(file.Kind) {
			t.Errorf("%s deve voltar ao original %q, está %q (%v)", file.Path, originalOf(file.Kind), b, err)
		}
	}
	for _, p := range append(tree(t, e.vault), tree(t, e.docs)...) {
		if strings.HasSuffix(p, ".legoloop-tmp") {
			t.Errorf("sobrou temporário: %s", p)
		}
	}
}

// originalOf: no plano de teste o conteúdo original de cada arquivo é o próprio Kind.
func originalOf(kind string) string { return kind }

func TestMaterialize_RejectsEmptyDuplicateAndCanceled(t *testing.T) {
	e := newEnv(t)
	empty := e.plan
	empty.Files = nil
	if _, err := fabricar.Materialize(context.Background(), empty, fabricar.OSFS()); err == nil {
		t.Error("plano vazio deve falhar")
	}
	dup := e.plan
	dup.Files = append(append([]fabricar.PlannedFile{}, e.plan.Files...), e.plan.Files[0])
	if _, err := fabricar.Materialize(context.Background(), dup, fabricar.OSFS()); err == nil {
		t.Error("caminho duplicado deve falhar")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := fabricar.Materialize(ctx, e.plan, fabricar.OSFS()); !errors.Is(err, context.Canceled) {
		t.Errorf("contexto cancelado deve abortar: %v", err)
	}
	if v := tree(t, e.vault); len(v) != 0 {
		t.Errorf("nada deve ser gravado: %v", v)
	}
}

func TestMaterialize_PersisteDAGIso(t *testing.T) {
	e := newEnv(t)
	dagIsoPath := filepath.Join(e.docs, "docs", "dags", "2026-11-09_DAG_X.json")
	e.plan.Files = append(e.plan.Files, fabricar.PlannedFile{
		Kind: "dag_iso",
		Path: dagIsoPath,
		Data: []byte(`{"name":"dag-iso"}`),
	})
	art, err := fabricar.Materialize(context.Background(), e.plan, fabricar.OSFS())
	if err != nil {
		t.Fatal(err)
	}
	if art.DagIsoPath != dagIsoPath {
		t.Errorf("esperava DagIsoPath %s, obteve %s", dagIsoPath, art.DagIsoPath)
	}
	b, err := os.ReadFile(dagIsoPath)
	if err != nil || string(b) != `{"name":"dag-iso"}` {
		t.Errorf("conteúdo de DAG ISO incorreto: %q err=%v", b, err)
	}
}

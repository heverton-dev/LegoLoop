package processo_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"legoloop/internal/processo"
)

var errInjected = errors.New("falha injetada")

// failFS delega ao disco real e falha na N-ésima gravação ou renomeação (1-based; 0 = nunca).
type failFS struct {
	processo.FS
	failWriteAt, failRenameAt int
	writes, renames           int
}

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
	store       processo.Store
	files       []processo.File
}

// newEnv cria um plano de 8 arquivos (3 no cofre, 5 no projeto), como a gravação de um processo real.
func newEnv(t *testing.T) env {
	t.Helper()
	vault, docs := t.TempDir(), t.TempDir()
	d := func(rel string) string { return filepath.Join(docs, "docs", rel) }
	return env{vault: vault, docs: docs,
		store: processo.Store{Roots: processo.Roots{Vault: vault, Docs: docs}, FS: processo.OSFS()},
		files: []processo.File{
			{Kind: "skill", Path: filepath.Join(vault, "skills", "legoloop-x", "SKILL.md"), Data: []byte("skill")},
			{Kind: "loop", Path: filepath.Join(vault, "loops", "legoloop-x.loop.json"), Data: []byte("loop")},
			{Kind: "dag", Path: filepath.Join(vault, "loops", "legoloop-x.dag.json"), Data: []byte("dag")},
			{Kind: "intent", Path: d("intents/2026-11-04_INTENT_X.md"), Data: []byte("intent")},
			{Kind: "spec", Path: d("specs/2026-11-04_SPEC_X.md"), Data: []byte("spec")},
			{Kind: "plan", Path: d("plans/2026-11-04_PLAN_X.md"), Data: []byte("plan")},
			{Kind: "route", Path: d("routes/2026-11-04_ROUTE_X.json"), Data: []byte("route")},
			{Kind: "template", Path: d("templates/TEMPLATE_X.md"), Data: []byte("template")},
		}}
}

// snapshot devolve o hash de toda a árvore (caminhos relativos e conteúdo): igual antes e depois = nada mudou.
func snapshot(t *testing.T, dirs ...string) string {
	t.Helper()
	h := sha256.New()
	for _, dir := range dirs {
		var linhas []string
		_ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
			if err != nil || p == dir {
				return nil
			}
			rel, _ := filepath.Rel(dir, p)
			if d.IsDir() {
				linhas = append(linhas, "D "+rel)
				return nil
			}
			b, _ := os.ReadFile(p)
			linhas = append(linhas, "F "+rel+" "+string(b))
			return nil
		})
		sort.Strings(linhas)
		h.Write([]byte(strings.Join(linhas, "\n") + "\x00"))
	}
	return hex.EncodeToString(h.Sum(nil))
}

func semTemporarios(t *testing.T, dirs ...string) {
	t.Helper()
	for _, dir := range dirs {
		_ = filepath.WalkDir(dir, func(p string, _ fs.DirEntry, _ error) error {
			if strings.HasSuffix(p, ".legoloop-tmp") {
				t.Errorf("sobrou temporário: %s", p)
			}
			return nil
		})
	}
}

// AC-E01.8: falha na 3ª escrita de 8 arquivos desfaz tudo, sem resíduo.
func TestStore_FalhaNaTerceiraEscritaDesfazTudo(t *testing.T) {
	e := newEnv(t)
	// estado anterior não vazio: um arquivo alheio no cofre e no projeto não pode ser tocado
	_ = os.MkdirAll(filepath.Join(e.vault, "skills", "outra"), 0o755)
	_ = os.WriteFile(filepath.Join(e.vault, "skills", "outra", "SKILL.md"), []byte("alheia"), 0o644)
	_ = os.WriteFile(filepath.Join(e.docs, "README.md"), []byte("projeto"), 0o644)
	antes := snapshot(t, e.vault, e.docs)

	e.store.FS = &failFS{FS: processo.OSFS(), failWriteAt: 3}
	err := e.store.Write(context.Background(), e.files, processo.WriteOpts{})
	if !errors.Is(err, errInjected) {
		t.Fatalf("esperava a falha injetada: %v", err)
	}
	if depois := snapshot(t, e.vault, e.docs); depois != antes {
		t.Error("a árvore mudou depois da falha")
	}
	semTemporarios(t, e.vault, e.docs)
}

func TestStore_GravaTudoEPermiteLer(t *testing.T) {
	e := newEnv(t)
	if err := e.store.Write(context.Background(), e.files, processo.WriteOpts{}); err != nil {
		t.Fatal(err)
	}
	for _, f := range e.files {
		if b, err := os.ReadFile(f.Path); err != nil || string(b) != string(f.Data) {
			t.Errorf("%s = %q (%v), quer %q", f.Path, b, err, f.Data)
		}
	}
	semTemporarios(t, e.vault, e.docs)
}

func TestStore_ExistenteSemForceDevolveErrExisteENaoGravaNada(t *testing.T) {
	e := newEnv(t)
	if err := e.store.Write(context.Background(), e.files, processo.WriteOpts{}); err != nil {
		t.Fatal(err)
	}
	antes := snapshot(t, e.vault, e.docs)
	err := e.store.Write(context.Background(), e.files, processo.WriteOpts{})
	if !errors.Is(err, processo.ErrExiste) {
		t.Fatalf("esperava ErrExiste: %v", err)
	}
	if snapshot(t, e.vault, e.docs) != antes {
		t.Error("o erro não pode alterar o disco")
	}
}

func TestStore_ForcePreservaOOriginalComoBak(t *testing.T) {
	e := newEnv(t)
	if err := e.store.Write(context.Background(), e.files, processo.WriteOpts{}); err != nil {
		t.Fatal(err)
	}
	for i := range e.files {
		e.files[i].Data = []byte("v2")
	}
	if err := e.store.Write(context.Background(), e.files, processo.WriteOpts{Force: true}); err != nil {
		t.Fatal(err)
	}
	for _, f := range e.files {
		if b, _ := os.ReadFile(f.Path); string(b) != "v2" {
			t.Errorf("%s deveria ter a v2, tem %q", f.Path, b)
		}
		if b, err := os.ReadFile(f.Path + ".legoloop-bak"); err != nil || string(b) != f.Kind {
			t.Errorf("o backup de %s deveria guardar %q, guarda %q (%v)", f.Path, f.Kind, b, err)
		}
	}
}

func TestStore_FalhaNaQuintaRenomeacaoComForceVoltaAoOriginal(t *testing.T) {
	e := newEnv(t)
	if err := e.store.Write(context.Background(), e.files, processo.WriteOpts{}); err != nil {
		t.Fatal(err)
	}
	antes := snapshot(t, e.vault, e.docs)
	for i := range e.files {
		e.files[i].Data = []byte("v2")
	}
	// com Force cada arquivo faz 2 renomeações (backup e promoção): a 5ª cai no meio do 3º arquivo
	e.store.FS = &failFS{FS: processo.OSFS(), failRenameAt: 5}
	err := e.store.Write(context.Background(), e.files, processo.WriteOpts{Force: true})
	if !errors.Is(err, errInjected) {
		t.Fatalf("esperava a falha injetada: %v", err)
	}
	if snapshot(t, e.vault, e.docs) != antes {
		t.Error("a árvore não voltou ao estado anterior")
	}
}

func TestStore_DryRunNaoGravaMasAntecipaConflito(t *testing.T) {
	e := newEnv(t)
	antes := snapshot(t, e.vault, e.docs)
	if err := e.store.Write(context.Background(), e.files, processo.WriteOpts{DryRun: true}); err != nil {
		t.Fatal(err)
	}
	if snapshot(t, e.vault, e.docs) != antes {
		t.Error("dry-run gravou algo")
	}
	if err := e.store.Write(context.Background(), e.files, processo.WriteOpts{}); err != nil {
		t.Fatal(err)
	}
	if err := e.store.Write(context.Background(), e.files, processo.WriteOpts{DryRun: true}); !errors.Is(err, processo.ErrExiste) {
		t.Errorf("dry-run deve antecipar ErrExiste: %v", err)
	}
}

func TestStore_RecusaPlanoInvalido(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	casos := map[string][]processo.File{
		"vazio":              nil,
		"duplicado":          append(append([]processo.File{}, e.files...), e.files[0]),
		"fora das raízes":    {{Kind: "skill", Path: filepath.Join(t.TempDir(), "x.md"), Data: []byte("x")}},
		"escapa com ..":      {{Kind: "skill", Path: filepath.Join(e.vault, "..", "fora.md"), Data: []byte("x")}},
		"caminho vazio":      {{Kind: "skill", Path: "", Data: []byte("x")}},
		"a própria raiz":     {{Kind: "skill", Path: e.vault, Data: []byte("x")}},
		"raiz de docs vazia": nil,
	}
	delete(casos, "raiz de docs vazia")
	for nome, files := range casos {
		t.Run(nome, func(t *testing.T) {
			antes := snapshot(t, e.vault, e.docs)
			if err := e.store.Write(ctx, files, processo.WriteOpts{}); !errors.Is(err, processo.ErrPlano) {
				t.Fatalf("esperava ErrPlano, veio %v", err)
			}
			if snapshot(t, e.vault, e.docs) != antes {
				t.Error("plano inválido alterou o disco")
			}
		})
	}
}

func TestStore_ContextoCancelado(t *testing.T) {
	e := newEnv(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := e.store.Write(ctx, e.files, processo.WriteOpts{}); !errors.Is(err, context.Canceled) {
		t.Fatalf("esperava context.Canceled: %v", err)
	}
}

func TestStore_RaizVaziaNaoPermiteNada(t *testing.T) {
	st := processo.Store{Roots: processo.Roots{}, FS: processo.OSFS()}
	err := st.Write(context.Background(), []processo.File{{Kind: "skill", Path: filepath.Join(t.TempDir(), "x"), Data: []byte("x")}}, processo.WriteOpts{})
	if !errors.Is(err, processo.ErrPlano) {
		t.Fatalf("sem raízes nada pode ser gravado (ErrPlano): %v", err)
	}
}

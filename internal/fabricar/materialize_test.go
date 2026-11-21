package fabricar_test

import (
	"context"
	"errors"
	"io/fs"
	"legoloop/internal/caminhos"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"legoloop/internal/assets"
	"legoloop/internal/fabricar"
	"legoloop/internal/listar"
	"legoloop/internal/matriz"
	"legoloop/internal/processo"
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
		{Kind: "processo", Path: filepath.Join(vault, caminhos.DirProcessos, "legoloop-x"+caminhos.SufixoAntigo), Data: []byte("processo")},
		{Kind: "dag", Path: filepath.Join(vault, caminhos.DirProcessos, "legoloop-x.dag.json"), Data: []byte("dag")},
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
	if art.SkillPath != e.plan.Files[0].Path || art.ProcessoPath != e.plan.Files[1].Path || art.DagPath != e.plan.Files[2].Path ||
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
	e.plan.Files[1].Path = filepath.Join(e.vault, caminhos.DirProcessos, "..", "..", "fora.json")
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

// pecaContratada cria uma skill de mídia com manifesto contrato.json ao lado do SKILL.md.
func pecaContratada(t *testing.T) matriz.Entrada {
	t.Helper()
	dir := t.TempDir()
	ref := filepath.Join(dir, "midia-x", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(ref), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(ref, []byte("# skill"), 0o644); err != nil {
		t.Fatal(err)
	}
	contrato := `{"nome":"midia-x","tipo":"skill","naturezas":["midia"],"sai":[{"tipo":"texto","caminho_padrao":"out/video.md","promete":true}]}`
	if err := os.WriteFile(filepath.Join(dir, "midia-x", "contrato.json"), []byte(contrato), 0o644); err != nil {
		t.Fatal(err)
	}
	p := skillApta("midia-x", ref, 0.9, "construir")
	p.Perfil.Naturezas = []string{"midia"}
	p.Perfil.CasosDeUso = []string{"geracao-de-midia"}
	return p
}

// TestSextetoCompleto prova a AC-E03.1: dry-run de natureza midia com peça contratada lista os 6
// artefatos do Sexteto com caminhos ISO YYYY-MM-DD_* e o contrato da peça dirige o gate da etapa.
func TestSextetoCompleto(t *testing.T) {
	r := newRig(t, []matriz.Entrada{pecaContratada(t)})
	r.llm.needs = `{"natureza": "midia", "etapas": ["construir"], "casos": ["geracao-de-midia"]}`
	r.req.DryRun = true
	res := r.run(t)
	familias := map[string]bool{}
	re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}_(INTERVIEW|INTENT|SPEC|PLAN|ROUTE|DAG|REPORT|TEMPLATE)_`)
	for _, a := range res.Report.Gravacao.Arquivos {
		if m := re.FindStringSubmatch(filepath.Base(a)); m != nil {
			familias[m[1]] = true
		}
	}
	slots := map[string]bool{
		"INTERVIEW":       familias["INTERVIEW"],
		"INTENT":          familias["INTENT"],
		"SPEC":            familias["SPEC"],
		"PLAN":            familias["PLAN"],
		"ROUTE/DAG":       familias["ROUTE"] || familias["DAG"],
		"TEMPLATE/REPORT": familias["TEMPLATE"] || familias["REPORT"],
	}
	for nome, ok := range slots {
		if !ok {
			t.Errorf("falta o artefato %s do Sexteto com caminho ISO YYYY-MM-DD_*: %v", nome, res.Report.Gravacao.Arquivos)
		}
	}
	if len(res.Report.Triade.Gates) == 0 || !strings.Contains(res.Report.Triade.Gates[0], "out/video.md") {
		t.Errorf("o contrato da peça contratada deve dirigir o gate da etapa: %v", res.Report.Triade.Gates)
	}
	verboSemPedidoNaoGeraProcesso(t)
}

// verboSemPedidoNaoGeraProcesso prova a premissa anti-reimplementação na camada de teste:
// verbos sem pedido do usuário (matriz, listar) não geram processo novo nem etapa que
// produz artefato. Só o fabricar, que recebe pedido, gera processo.
func verboSemPedidoNaoGeraProcesso(t *testing.T) {
	t.Helper()
	vault := t.TempDir()
	ctx := context.Background()
	if _, err := matriz.Generate(ctx, matriz.Options{Vault: vault}, matriz.Deps{
		Scan: func(context.Context) ([]assets.Item, error) { return nil, nil },
		Home: vault,
	}); err != nil {
		t.Fatalf("matriz gerar não deve falhar num cofre vazio: %v", err)
	}
	if _, err := listar.Run(ctx, listar.Deps{
		Store: processo.Store{Roots: processo.Roots{Vault: vault}},
		Home:  vault,
	}); err != nil {
		t.Fatalf("listar não deve falhar num cofre vazio: %v", err)
	}
	err := filepath.WalkDir(vault, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if strings.HasSuffix(path, ".processo.json") {
			t.Errorf("verbo sem pedido gerou processo novo com etapa que produz artefato: %s",
				strings.TrimPrefix(path, vault+string(os.PathSeparator)))
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

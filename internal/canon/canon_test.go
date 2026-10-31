package canon

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func write(t *testing.T, path, body string, mod time.Time) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(path, mod, mod); err != nil {
		t.Fatal(err)
	}
}

const skillBody = "---\nname: interview\ndescription: entrevista\n---\n# Interview\n"

// fixture: 3 cópias idênticas do SKILL.md em 2 raízes, 1 arquivo único, 1 .txt ignorado.
func fixture(t *testing.T) (rootA, rootB string) {
	t.Helper()
	base := t.TempDir()
	rootA, rootB = filepath.Join(base, "A"), filepath.Join(base, "B")
	old := time.Now().Add(-48 * time.Hour)
	write(t, filepath.Join(rootA, "skills", "interview", "SKILL.md"), skillBody, old)
	write(t, filepath.Join(rootA, "proj", ".agents", "interview", "SKILL.md"), skillBody, time.Now())
	write(t, filepath.Join(rootB, "interview", "SKILL.md"), skillBody, old)
	write(t, filepath.Join(rootB, "unico.json"), `{"a":1,"b":2}`, old)
	write(t, filepath.Join(rootB, "ignorado.txt"), skillBody, old)
	write(t, filepath.Join(rootB, "node_modules", "x", "SKILL.md"), skillBody, old)
	return rootA, rootB
}

func TestHashTree(t *testing.T) {
	a, b := fixture(t)
	files, err := HashTree(context.Background(), ScanOptions{Roots: []string{a, b}, Workers: 4})
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 4 {
		t.Fatalf("esperado 4 arquivos .md/.json fora de node_modules, veio %d: %+v", len(files), files)
	}
	for _, f := range files {
		if len(f.SHA256) != 64 || f.Size == 0 {
			t.Fatalf("metadados incompletos: %+v", f)
		}
		if strings.HasSuffix(f.Path, "SKILL.md") && f.FrontmatterKeys != 2 {
			t.Fatalf("frontmatter keys = %d", f.FrontmatterKeys)
		}
		if strings.HasSuffix(f.Path, "unico.json") && f.FrontmatterKeys != 2 {
			t.Fatalf("json keys = %d", f.FrontmatterKeys)
		}
	}
}

func TestHashTreeExcludeAndErrors(t *testing.T) {
	a, b := fixture(t)
	files, err := HashTree(context.Background(), ScanOptions{Roots: []string{a, b}, Exclude: []string{b}})
	if err != nil || len(files) != 2 {
		t.Fatalf("exclude falhou: %d %v", len(files), err)
	}
	if _, err := HashTree(context.Background(), ScanOptions{Roots: []string{filepath.Join(a, "nao-existe")}}); err == nil {
		t.Fatal("raiz inexistente deveria falhar")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := HashTree(ctx, ScanOptions{Roots: []string{a}}); err == nil {
		t.Fatal("contexto cancelado deveria falhar")
	}
}

func TestElectAndGroup(t *testing.T) {
	now := time.Now()
	fs := []FileInfo{
		{Path: "/z/longo/caminho/SKILL.md", SHA256: "h1", ModTime: now, FrontmatterKeys: 1, Size: 10},
		{Path: "/a/SKILL.md", SHA256: "h1", ModTime: now, FrontmatterKeys: 3, Size: 10},
		{Path: "/b/SKILL.md", SHA256: "h1", ModTime: now.Add(-time.Hour), FrontmatterKeys: 9, Size: 10},
		{Path: "/c/unico.md", SHA256: "h2", ModTime: now, Size: 5},
	}
	if got := Elect(fs[:3]); got.Path != "/a/SKILL.md" {
		t.Fatalf("eleito = %s (critério: mais recente, depois mais frontmatter)", got.Path)
	}
	tie := []FileInfo{{Path: "/bb/x.md"}, {Path: "/a/x.md"}, {Path: "/ab/x.md"}}
	if got := Elect(tie); got.Path != "/a/x.md" {
		t.Fatalf("desempate por caminho mais curto/lexical = %s", got.Path)
	}
	groups := GroupDuplicates(fs)
	if len(groups) != 1 || groups[0].SHA256 != "h1" || len(groups[0].Files) != 3 {
		t.Fatalf("grupos = %+v", groups)
	}
}

func TestPlanDryRun(t *testing.T) {
	a, b := fixture(t)
	vault := filepath.Join(t.TempDir(), "canonical")
	files, _ := HashTree(context.Background(), ScanOptions{Roots: []string{a, b}})
	p := BuildPlan(files, vault)
	if len(p.Groups) != 1 {
		t.Fatalf("grupos = %d", len(p.Groups))
	}
	g := p.Groups[0]
	if !strings.Contains(g.Canonical.Path, ".agents") {
		t.Fatalf("canônico deveria ser o mais recente: %s", g.Canonical.Path)
	}
	if !strings.HasPrefix(g.VaultPath, vault) || filepath.Base(g.VaultPath) != "SKILL.md" {
		t.Fatalf("vault path = %s", g.VaultPath)
	}
	if want := int64(2 * len(skillBody)); p.SavedBytes != want || g.SavedBytes != want {
		t.Fatalf("economia = %d, want %d", p.SavedBytes, want)
	}
	if _, err := os.Stat(vault); !os.IsNotExist(err) {
		t.Fatal("dry-run não pode tocar o disco")
	}
}

func assertLinked(t *testing.T, path, vaultPath string) {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil || string(b) != skillBody {
		t.Fatalf("%s conteúdo quebrado: %q %v", path, b, err)
	}
	fi, _ := os.Stat(path)
	vi, _ := os.Stat(vaultPath)
	if !os.SameFile(fi, vi) {
		t.Fatalf("%s não aponta para o cofre", path)
	}
	if _, err := os.Stat(path + backupSuffix); !os.IsNotExist(err) {
		t.Fatalf("backup residual em %s", path)
	}
}

func TestApplyHardlinkFallback(t *testing.T) {
	a, b := fixture(t)
	vault := filepath.Join(t.TempDir(), "canonical")
	files, _ := HashTree(context.Background(), ScanOptions{Roots: []string{a, b}})
	p := BuildPlan(files, vault)
	l := Linker{Symlink: func(string, string) error { return errors.New("privilégio ausente") }, Hardlink: os.Link}
	rep := Apply(context.Background(), p, l)
	if len(rep.Failed) != 0 || len(rep.Linked) != 3 || rep.SavedBytes != p.SavedBytes {
		t.Fatalf("relatório = %+v", rep)
	}
	for _, r := range rep.Linked {
		if r.Method != MethodHardlink {
			t.Fatalf("método = %s", r.Method)
		}
		assertLinked(t, r.Path, p.Groups[0].VaultPath)
	}
	// Idempotência: segunda execução não encontra nada a fazer.
	files2, _ := HashTree(context.Background(), ScanOptions{Roots: []string{a, b}})
	rep2 := Apply(context.Background(), BuildPlan(files2, vault), l)
	if len(rep2.Linked) != 0 || len(rep2.Failed) != 0 || rep2.SavedBytes != 0 {
		t.Fatalf("segunda execução não é idempotente: %+v", rep2)
	}
}

func TestApplySymlinkPreferred(t *testing.T) {
	a, b := fixture(t)
	vault := filepath.Join(t.TempDir(), "canonical")
	files, _ := HashTree(context.Background(), ScanOptions{Roots: []string{a, b}})
	var symlinkTargets []string
	l := Linker{
		Symlink:  func(old, new string) error { symlinkTargets = append(symlinkTargets, old); return os.Link(old, new) },
		Hardlink: func(string, string) error { t.Fatal("hardlink não deveria ser usado"); return nil },
	}
	rep := Apply(context.Background(), BuildPlan(files, vault), l)
	if len(rep.Linked) != 3 || rep.Linked[0].Method != MethodSymlink {
		t.Fatalf("relatório = %+v", rep)
	}
	for _, tgt := range symlinkTargets {
		if !filepath.IsAbs(tgt) {
			t.Fatalf("alvo do symlink deve ser absoluto: %s", tgt)
		}
	}
}

func TestApplyRollbackWhenLinkingFails(t *testing.T) {
	a, b := fixture(t)
	vault := filepath.Join(t.TempDir(), "canonical")
	files, _ := HashTree(context.Background(), ScanOptions{Roots: []string{a, b}})
	fail := func(string, string) error { return errors.New("volume diferente") }
	rep := Apply(context.Background(), BuildPlan(files, vault), Linker{Symlink: fail, Hardlink: fail})
	if len(rep.Linked) != 0 || len(rep.Failed) != 3 || rep.SavedBytes != 0 {
		t.Fatalf("relatório = %+v", rep)
	}
	for _, f := range rep.Failed {
		b, err := os.ReadFile(f.Path)
		if err != nil || string(b) != skillBody {
			t.Fatalf("original não restaurado em %s", f.Path)
		}
		if _, err := os.Stat(f.Path + backupSuffix); !os.IsNotExist(err) {
			t.Fatal("backup residual após rollback")
		}
	}
}

func TestApplySkipsFilesChangedSincePlan(t *testing.T) {
	a, b := fixture(t)
	vault := filepath.Join(t.TempDir(), "canonical")
	files, _ := HashTree(context.Background(), ScanOptions{Roots: []string{a, b}})
	p := BuildPlan(files, vault)
	changed := filepath.Join(b, "interview", "SKILL.md")
	write(t, changed, "editado pelo usuário depois do plano", time.Now())
	rep := Apply(context.Background(), p, DefaultLinker())
	if len(rep.Failed) != 1 || rep.Failed[0].Path != changed {
		t.Fatalf("arquivo alterado deveria ser preservado: %+v", rep)
	}
	if body, _ := os.ReadFile(changed); string(body) != "editado pelo usuário depois do plano" {
		t.Fatal("edição do usuário foi destruída")
	}
}

func TestApplyCancelled(t *testing.T) {
	a, b := fixture(t)
	files, _ := HashTree(context.Background(), ScanOptions{Roots: []string{a, b}})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	rep := Apply(ctx, BuildPlan(files, filepath.Join(t.TempDir(), "v")), DefaultLinker())
	if len(rep.Linked) != 0 || len(rep.Failed) == 0 {
		t.Fatalf("cancelamento ignorado: %+v", rep)
	}
}

func TestDefaultVault(t *testing.T) {
	v, err := DefaultVault()
	if err != nil || !strings.HasSuffix(filepath.ToSlash(v), ".aidd/canonical") {
		t.Fatalf("vault = %s %v", v, err)
	}
}

func planFor(t *testing.T) (Plan, string, string) {
	t.Helper()
	a, b := fixture(t)
	files, _ := HashTree(context.Background(), ScanOptions{Roots: []string{a, b}})
	return BuildPlan(files, filepath.Join(t.TempDir(), "canonical")), a, b
}

func TestApplyRefusesCorruptVault(t *testing.T) {
	p, _, _ := planFor(t)
	write(t, p.Groups[0].VaultPath, "adulterado", time.Now())
	rep := Apply(context.Background(), p, DefaultLinker())
	if len(rep.Linked) != 0 || len(rep.Failed) != 3 || !strings.Contains(rep.Failed[0].Error, "corrompido") {
		t.Fatalf("cofre adulterado deveria bloquear: %+v", rep)
	}
}

func TestApplyRefusesExistingBackupAndCanonicalChange(t *testing.T) {
	p, _, _ := planFor(t)
	victim := nonCanonical(p.Groups[0])
	write(t, victim+backupSuffix, "backup do usuário", time.Now())
	rep := Apply(context.Background(), p, Linker{Symlink: os.Link, Hardlink: os.Link})
	if len(rep.Failed) != 1 || rep.Failed[0].Path != victim {
		t.Fatalf("backup pré-existente deveria ser preservado: %+v", rep)
	}
	if b, _ := os.ReadFile(victim + backupSuffix); string(b) != "backup do usuário" {
		t.Fatal("backup do usuário sobrescrito")
	}

	p2, _, _ := planFor(t)
	write(t, p2.Groups[0].Canonical.Path, "canônico alterado", time.Now())
	rep2 := Apply(context.Background(), p2, DefaultLinker())
	if len(rep2.Linked) != 0 || len(rep2.Failed) != 3 {
		t.Fatalf("canônico alterado deveria abortar o grupo: %+v", rep2)
	}
}

func TestApplyRejectsNonRegularTarget(t *testing.T) {
	p, _, _ := planFor(t)
	victim := nonCanonical(p.Groups[0])
	os.Remove(victim)
	if err := os.Mkdir(victim, 0o755); err != nil {
		t.Fatal(err)
	}
	rep := Apply(context.Background(), p, DefaultLinker())
	if len(rep.Failed) != 1 || !strings.Contains(rep.Failed[0].Error, "regular") {
		t.Fatalf("diretório no lugar do arquivo deveria ser rejeitado: %+v", rep)
	}
}

// Um "link" que não aponta para o cofre é indistinguível de um arquivo de terceiros:
// nada é apagado, o original fica íntegro no backup e a falha é reportada.
func TestForeignFileInPlaceOfLinkIsKept(t *testing.T) {
	p, _, _ := planFor(t)
	bogus := func(_, newname string) error { return os.WriteFile(newname, []byte("lixo"), 0o644) }
	rep := Apply(context.Background(), p, Linker{Symlink: bogus, Hardlink: bogus})
	if len(rep.Linked) != 0 || len(rep.Failed) != 3 {
		t.Fatalf("link inválido não pode ser aceito: %+v", rep)
	}
	for _, f := range rep.Failed {
		if b, _ := os.ReadFile(f.Path + backupSuffix); string(b) != skillBody {
			t.Fatalf("original não preservado para %s", f.Path)
		}
		if !strings.Contains(f.Error, "substituído") {
			t.Fatalf("motivo não reportado: %s", f.Error)
		}
	}
}

func TestMetaKeysEdgeCases(t *testing.T) {
	if metaKeys("a.json", []byte("[1,2]")) != 0 || metaKeys("a.md", []byte("# sem frontmatter")) != 0 {
		t.Fatal("conteúdo sem metadados deve contar 0")
	}
}

func nonCanonical(g GroupPlan) string {
	for _, l := range g.Links {
		if l != g.Canonical.Path {
			return l
		}
	}
	return ""
}

// SEC-04: arquivos dentro de repositórios git ficam fora da canonização por padrão
// (um link commitado quebraria o repositório para colaboradores).
func TestHashTreeSkipsGitRepos(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "repo", "SKILL.md"), skillBody, time.Now())
	os.MkdirAll(filepath.Join(root, "repo", ".git"), 0o755)
	write(t, filepath.Join(root, "solto", "SKILL.md"), skillBody, time.Now())
	files, err := HashTree(context.Background(), ScanOptions{Roots: []string{root}})
	if err != nil || len(files) != 1 || strings.Contains(files[0].Path, "repo") {
		t.Fatalf("repositório git deveria ser ignorado: %+v %v", files, err)
	}
	files, _ = HashTree(context.Background(), ScanOptions{Roots: []string{root}, IncludeGit: true})
	if len(files) != 2 {
		t.Fatalf("--include-git deveria incluir o repositório: %d", len(files))
	}
}

// SEC-05: edição gravada no original entre a verificação e o link não pode se perder.
func TestEditDuringLinkIsPreserved(t *testing.T) {
	p, _, _ := planFor(t)
	editor := func(oldname, newname string) error {
		// o "editor" grava no inode original, que agora está no backup
		if err := os.WriteFile(newname+backupSuffix, []byte("edição concorrente"), 0o644); err != nil {
			return err
		}
		return os.Link(oldname, newname)
	}
	rep := Apply(context.Background(), p, Linker{Symlink: editor, Hardlink: os.Link})
	if len(rep.Linked) != 0 {
		t.Fatalf("links criados apesar da edição concorrente: %+v", rep.Linked)
	}
	for _, f := range rep.Failed {
		if f.Path == p.Groups[0].Canonical.Path {
			continue
		}
		b, _ := os.ReadFile(f.Path)
		if string(b) != "edição concorrente" {
			t.Fatalf("edição perdida em %s: %q", f.Path, b)
		}
		if _, err := os.Stat(f.Path + backupSuffix); !os.IsNotExist(err) {
			t.Fatal("backup residual")
		}
	}
}

// R2-2: raiz DENTRO de um repositório git (".git" num ancestral) também fica fora por padrão.
func TestHashTreeSkipsRootInsideGitRepo(t *testing.T) {
	repo := t.TempDir()
	os.MkdirAll(filepath.Join(repo, ".git"), 0o755)
	docs := filepath.Join(repo, "docs")
	write(t, filepath.Join(docs, "a", "SKILL.md"), skillBody, time.Now())
	files, err := HashTree(context.Background(), ScanOptions{Roots: []string{docs}})
	if err != nil || len(files) != 0 {
		t.Fatalf("raiz dentro de repo git deveria ser ignorada: %d %v", len(files), err)
	}
	files, _ = HashTree(context.Background(), ScanOptions{Roots: []string{docs}, IncludeGit: true})
	if len(files) != 1 {
		t.Fatal("--include-git deveria incluir")
	}
}

// R2-3: se o arquivo foi recriado (salvamento atômico) durante a troca, o rollback não o sobrescreve.
func TestRollbackDoesNotClobberRecreatedFile(t *testing.T) {
	p, _, _ := planFor(t)
	victim := nonCanonical(p.Groups[0])
	recreate := func(_, newname string) error {
		if newname == victim {
			os.WriteFile(newname, []byte("versão nova do editor"), 0o644)
		}
		return errors.New("falhou")
	}
	rep := Apply(context.Background(), p, Linker{Symlink: recreate, Hardlink: func(string, string) error { return errors.New("falhou") }})
	if b, _ := os.ReadFile(victim); string(b) != "versão nova do editor" {
		t.Fatalf("rollback sobrescreveu a versão nova: %q", b)
	}
	if b, _ := os.ReadFile(victim + backupSuffix); string(b) != skillBody {
		t.Fatal("original deveria ficar preservado no backup")
	}
	found := false
	for _, f := range rep.Failed {
		found = found || (f.Path == victim && strings.Contains(f.Error, "recriado"))
	}
	if !found {
		t.Fatalf("conflito não reportado: %+v", rep.Failed)
	}
}

// R3-2: se o editor substituir o link por arquivo novo antes da verificação, nada é apagado.
func TestReplacedLinkIsNotDeleted(t *testing.T) {
	p, _, _ := planFor(t)
	victim := nonCanonical(p.Groups[0])
	editor := func(oldname, newname string) error {
		if err := os.Link(oldname, newname); err != nil {
			return err
		}
		if newname == victim { // salvamento atômico: troca a entrada por um arquivo novo
			tmp := newname + ".tmp-editor"
			os.WriteFile(tmp, []byte("salvo pelo editor"), 0o644)
			return os.Rename(tmp, newname)
		}
		return nil
	}
	rep := Apply(context.Background(), p, Linker{Symlink: editor, Hardlink: os.Link})
	if b, _ := os.ReadFile(victim); string(b) != "salvo pelo editor" {
		t.Fatalf("edição do editor apagada: %q", b)
	}
	if b, _ := os.ReadFile(victim + backupSuffix); string(b) != skillBody {
		t.Fatal("original deveria ficar no backup")
	}
	if len(rep.Linked) != 2 {
		t.Fatalf("demais cópias deveriam ser linkadas: %+v", rep)
	}
}

// R3-1: a busca por repositório git para no HOME (dotfiles em ~/.git não desativam tudo).
func TestGitRootStopsAtHome(t *testing.T) {
	home := t.TempDir()
	os.MkdirAll(filepath.Join(home, ".git"), 0o755)
	dir := filepath.Join(home, "skills", "x")
	os.MkdirAll(dir, 0o755)
	if repo, ok := GitRoot(dir, home); ok {
		t.Fatalf("~/.git não deveria contar: %s", repo)
	}
	os.MkdirAll(filepath.Join(home, "skills", ".git"), 0o755)
	if repo, ok := GitRoot(dir, home); !ok || repo != filepath.Join(home, "skills") {
		t.Fatalf("repo intermediário não detectado: %s %v", repo, ok)
	}
}

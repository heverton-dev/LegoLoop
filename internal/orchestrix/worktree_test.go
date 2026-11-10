package orchestrix

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// gitT roda git num repositório temporário, isolado de qualquer configuração do usuário.
func gitT(t *testing.T, dir string, args ...string) string {
	t.Helper()
	full := append([]string{"-c", "user.name=t", "-c", "user.email=t@t", "-c", "core.autocrlf=false", "-c", "commit.gpgsign=false"}, args...)
	cmd := exec.Command("git", full...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL="+os.DevNull, "GIT_CONFIG_NOSYSTEM=1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}

func writeT(t *testing.T, dir, rel, content string) {
	t.Helper()
	p := filepath.Join(dir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readT(t *testing.T, dir, rel string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(rel)))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// newRepoT cria um repositório real em t.TempDir() com um commit inicial.
func newRepoT(t *testing.T) string {
	t.Helper()
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	repo := t.TempDir()
	gitT(t, repo, "init", "-q")
	writeT(t, repo, "a.txt", "a\n")
	writeT(t, repo, "docs/b.txt", "b\n")
	writeT(t, repo, "del.txt", "del\n")
	gitT(t, repo, "add", "-A")
	gitT(t, repo, "commit", "-q", "-m", "base")
	return repo
}

// linkT cria link de diretório: symlink ou, sem privilégio, junction via mklink /J.
func linkT(t *testing.T, link, target string) {
	t.Helper()
	if err := os.Symlink(target, link); err == nil {
		return
	}
	if out, err := exec.Command("cmd", "/c", "mklink", "/J", link, target).CombinedOutput(); err != nil {
		t.Skipf("sem como criar link de diretório: %v %s", err, out)
	}
}

func sorted(s []string) []string {
	out := append([]string{}, s...)
	sort.Strings(out)
	return out
}

func TestParsePorcelainZ(t *testing.T) {
	in := " M a.txt\x00?? novo/x.txt\x00 D del.txt\x00R  novo.txt\x00velho.txt\x00"
	got, err := parsePorcelainZ([]byte(in))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"a.txt", "novo/x.txt", "del.txt", "novo.txt", "velho.txt"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if _, err := parsePorcelainZ([]byte("R  so.txt\x00")); err == nil {
		t.Fatal("renomeação sem origem aceita")
	}
	if _, err := parsePorcelainZ([]byte("lixo\x00")); err == nil {
		t.Fatal("registro inválido aceito")
	}
}

func TestWorktreeCreatesFromHEAD(t *testing.T) {
	repo := newRepoT(t)
	writeT(t, repo, "a.txt", "sujo no principal\n") // alteração não commitada não entra
	writeT(t, repo, "solto.txt", "x")
	w, err := CreateWorktree(context.Background(), repo)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Discard(context.Background())
	if got := readT(t, w.Dir, "a.txt"); got != "a\n" {
		t.Fatalf("worktree não veio de HEAD: %q", got)
	}
	if _, err := os.Stat(filepath.Join(w.Dir, "solto.txt")); !os.IsNotExist(err) {
		t.Fatal("arquivo não rastreado vazou para a worktree")
	}
	if strings.HasPrefix(w.Dir, repo) {
		t.Fatal("worktree criada dentro do repositório principal")
	}
	changed, err := w.Changed(context.Background())
	if err != nil || len(changed) != 0 {
		t.Fatalf("worktree recém-criada deve estar limpa: %v %v", changed, err)
	}
}

func TestWorktreeChangedListsEveryKind(t *testing.T) {
	repo := newRepoT(t)
	ctx := context.Background()
	w, err := CreateWorktree(ctx, repo)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Discard(ctx)
	writeT(t, w.Dir, "a.txt", "editado\n")
	writeT(t, w.Dir, "novo/dir/c.txt", "novo\n")
	os.Remove(filepath.Join(w.Dir, "del.txt"))
	gitT(t, w.Dir, "mv", "docs/b.txt", "docs/renomeado.txt")
	got, err := w.Changed(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"a.txt", "del.txt", "docs/b.txt", "docs/renomeado.txt", "novo/dir/c.txt"}
	if !reflect.DeepEqual(sorted(got), want) {
		t.Fatalf("alterados %v, esperado %v", sorted(got), want)
	}
}

func TestWorktreeApplyCopiesDiffToMain(t *testing.T) {
	repo := newRepoT(t)
	ctx := context.Background()
	w, err := CreateWorktree(ctx, repo)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Discard(ctx)
	writeT(t, w.Dir, "a.txt", "editado\n")
	writeT(t, w.Dir, "novo/dir/c.txt", "novo\n")
	os.Remove(filepath.Join(w.Dir, "del.txt"))
	changed, err := w.Changed(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.CheckLinks(changed); err != nil {
		t.Fatal(err)
	}
	if err := w.Apply(changed); err != nil {
		t.Fatal(err)
	}
	if got := readT(t, repo, "a.txt"); got != "editado\n" {
		t.Fatalf("a.txt = %q", got)
	}
	if got := readT(t, repo, "novo/dir/c.txt"); got != "novo\n" {
		t.Fatalf("novo/dir/c.txt = %q", got)
	}
	if _, err := os.Stat(filepath.Join(repo, "del.txt")); !os.IsNotExist(err) {
		t.Fatal("del.txt não foi removido no principal")
	}
	if got := readT(t, repo, "docs/b.txt"); got != "b\n" {
		t.Fatal("arquivo não alterado mudou no principal")
	}
	if left, _ := filepath.Glob(filepath.Join(repo, "*.legoloop-tmp")); len(left) != 0 {
		t.Fatalf("temporário sobrou: %v", left)
	}
	if st := gitT(t, repo, "log", "--oneline"); strings.Count(st, "\n") != 1 {
		t.Fatalf("a Orchestrix fez commit: %s", st)
	}
}

func TestWorktreeApplyRefusesEscapingPaths(t *testing.T) {
	repo := newRepoT(t)
	ctx := context.Background()
	w, err := CreateWorktree(ctx, repo)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Discard(ctx)
	for _, bad := range []string{"../fora.txt", "/abs.txt", `C:\x.txt`, "."} {
		if err := w.Apply([]string{bad}); err == nil {
			t.Errorf("Apply aceitou %q", bad)
		}
		if err := w.CheckLinks([]string{bad}); err == nil {
			t.Errorf("CheckLinks aceitou %q", bad)
		}
	}
}

// Link em caminho alterado reprova: o próprio caminho e um diretório-pai.
func TestWorktreeLinkInChangedPathFails(t *testing.T) {
	repo := newRepoT(t)
	ctx := context.Background()
	w, err := CreateWorktree(ctx, repo)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Discard(ctx)
	outside := t.TempDir()
	writeT(t, outside, "alvo.txt", "segredo\n")
	linkT(t, filepath.Join(w.Dir, "out"), outside)

	if err := w.CheckLinks([]string{"out"}); err == nil {
		t.Fatal("link em caminho alterado aceito")
	}
	if err := w.CheckLinks([]string{"out/alvo.txt"}); err == nil {
		t.Fatal("caminho que atravessa link aceito")
	}
	if err := w.Apply([]string{"out/alvo.txt"}); err == nil {
		t.Fatal("Apply atravessou link")
	}
	if err := w.Apply([]string{"out"}); err == nil {
		t.Fatal("Apply aceitou o próprio link")
	}
	if _, err := os.Stat(filepath.Join(repo, "out")); !os.IsNotExist(err) {
		t.Fatal("algo foi escrito no principal apesar do link")
	}
	// Caminho sem link continua aceito.
	if err := w.CheckLinks([]string{"a.txt"}); err != nil {
		t.Fatalf("caminho comum recusado: %v", err)
	}
	// Descartar não pode seguir o link e apagar o alvo.
	if err := w.Discard(ctx); err != nil {
		t.Fatal(err)
	}
	if got := readT(t, outside, "alvo.txt"); got != "segredo\n" {
		t.Fatal("Discard apagou o conteúdo do alvo do link")
	}
}

// Link no destino (principal) também reprova: o principal não pode desviar a escrita.
func TestWorktreeApplyRefusesLinkedDestination(t *testing.T) {
	repo := newRepoT(t)
	ctx := context.Background()
	w, err := CreateWorktree(ctx, repo)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Discard(ctx)
	outside := t.TempDir()
	linkT(t, filepath.Join(repo, "ponte"), outside)
	writeT(t, w.Dir, "ponte/x.txt", "x")
	if err := w.Apply([]string{"ponte/x.txt"}); err == nil {
		t.Fatal("Apply escreveu através de link do principal")
	}
	if _, err := os.Stat(filepath.Join(outside, "x.txt")); !os.IsNotExist(err) {
		t.Fatal("arquivo vazou para fora do repositório")
	}
}

func worktreeCount(t *testing.T, repo string) int {
	t.Helper()
	out := gitT(t, repo, "worktree", "list", "--porcelain")
	return strings.Count(out, "worktree ")
}

func TestRunInWorktreeDiscardsOnError(t *testing.T) {
	repo := newRepoT(t)
	boom := errors.New("falha do nó")
	var dir string
	err := RunInWorktree(context.Background(), repo, func(w *Worktree) error {
		dir = w.Dir
		writeT(t, w.Dir, "a.txt", "meio caminho\n")
		if worktreeCount(t, repo) != 2 {
			t.Error("worktree não registrada durante a execução")
		}
		return boom
	})
	if !errors.Is(err, boom) {
		t.Fatalf("erro original perdido: %v", err)
	}
	if _, statErr := os.Stat(dir); !os.IsNotExist(statErr) {
		t.Fatal("diretório da worktree sobrou após erro")
	}
	if n := worktreeCount(t, repo); n != 1 {
		t.Fatalf("worktree ainda registrada no git: %d", n)
	}
	if got := readT(t, repo, "a.txt"); got != "a\n" {
		t.Fatal("alteração da worktree vazou para o principal")
	}
}

func TestRunInWorktreeDiscardsOnSuccessAndCancelledContext(t *testing.T) {
	repo := newRepoT(t)
	ctx, cancel := context.WithCancel(context.Background())
	var dir string
	err := RunInWorktree(ctx, repo, func(w *Worktree) error {
		dir = w.Dir
		cancel() // contexto cancelado durante fn não pode impedir o descarte
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, statErr := os.Stat(dir); !os.IsNotExist(statErr) {
		t.Fatal("worktree sobrou com contexto cancelado")
	}
	if n := worktreeCount(t, repo); n != 1 {
		t.Fatalf("worktree ainda registrada: %d", n)
	}
}

func TestWorktreeDiscardIsIdempotent(t *testing.T) {
	repo := newRepoT(t)
	w, err := CreateWorktree(context.Background(), repo)
	if err != nil {
		t.Fatal(err)
	}
	for i := range 2 {
		if err := w.Discard(context.Background()); err != nil {
			t.Fatalf("Discard %d: %v", i, err)
		}
	}
}

func TestCreateWorktreeFailsOutsideRepoAndLeavesNothing(t *testing.T) {
	notRepo := t.TempDir()
	t.Setenv("GIT_CEILING_DIRECTORIES", filepath.Dir(notRepo))
	before, _ := filepath.Glob(filepath.Join(os.TempDir(), "legoloop-wt-*"))
	if _, err := CreateWorktree(context.Background(), notRepo); err == nil {
		t.Fatal("worktree criada fora de repositório")
	}
	after, _ := filepath.Glob(filepath.Join(os.TempDir(), "legoloop-wt-*"))
	if len(after) != len(before) {
		t.Fatalf("diretório temporário vazou: %d -> %d", len(before), len(after))
	}
}

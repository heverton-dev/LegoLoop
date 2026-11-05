package cli

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/processo"
)

// fotoRemover guarda caminho e hash de tudo sob root: duas fotos iguais provam que nada mudou.
func fotoRemover(t *testing.T, root string) string {
	t.Helper()
	var b strings.Builder
	err := filepath.WalkDir(root, func(p string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if e.IsDir() {
			fmt.Fprintf(&b, "%s dir\n", p)
			return nil
		}
		data, err := os.ReadFile(p)
		fmt.Fprintf(&b, "%s %x\n", p, sha256.Sum256(data))
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func existeRem(p string) bool { _, err := os.Lstat(p); return err == nil }

// completarCofreRemover põe no cofre de cofreListar o que o fabricar grava além do processo: a skill e o DAG.
func completarCofreRemover(t *testing.T, vault string) {
	t.Helper()
	for _, n := range []string{"legoloop-a", "legoloop-b"} {
		gravarArq(t, filepath.Join(vault, "skills", n, "SKILL.md"), "# "+n+"\n")
		gravarArq(t, filepath.Join(vault, "loops", n+".loop.json"), "{}")
		gravarArq(t, filepath.Join(vault, "loops", n+".dag.json"), "{}")
	}
}

// TestRemoverBinario prova o AC-E03.3 pelo executável: a falha com ErrEmUso (exit 1, nada muda); b sem
// confirmação não muda nada; b com confirmação sai; a com --yes sai com os 2 links.
func TestRemoverBinario(t *testing.T) {
	bin := buildBinary(t)
	home, vault := cofreListar(t)
	completarCofreRemover(t, vault)
	router := binRouter(t, 90)
	rodar := func(stdin string, args ...string) binRun {
		return runBinaryStdin(t, bin, router, home, t.TempDir(), stdin, append(args, "--vault", vault)...)
	}
	antes := fotoRemover(t, home)

	r := rodar("", "remover", "legoloop-a", "--yes")
	if r.code != 1 || !strings.Contains(r.errb, "legoloop-b") || !strings.Contains(r.errb, "usa ele") {
		t.Fatalf("a em uso: exit %d; esperava 1 com o nome de quem usa\n%s\n%s", r.code, r.out, r.errb)
	}
	if strings.Contains(r.out, "AÇÃO:") || fotoRemover(t, home) != antes {
		t.Fatal("a em uso não pode mostrar cartão nem mudar nada")
	}

	for _, stdin := range []string{"n\n", "\n", ""} {
		r = rodar(stdin, "remover", "legoloop-b")
		if r.code != 1 || fotoRemover(t, home) != antes {
			t.Fatalf("b sem confirmação (stdin %q): exit %d; esperava 1 e nada mudado\n%s\n%s", stdin, r.code, r.out, r.errb)
		}
	}
	if !strings.Contains(r.out, "AVISO: Esta ação não tem volta") || strings.Index(r.out, "AVISO:") > strings.Index(r.out, "O QUE VAI ACONTECER") {
		t.Errorf("o aviso de ação sem volta vem antes dos passos:\n%s", r.out)
	}

	if strings.Contains(r.out, "vínculos") {
		t.Errorf("b não tem vínculo, e o cartão não deve falar de vínculos:\n%s", r.out)
	}

	r = rodar("s\n", "remover", "legoloop-b")
	if r.code != 0 || !strings.Contains(r.out, "Processo legoloop-b removido") {
		t.Fatalf("b confirmado: exit %d\n%s\n%s", r.code, r.out, r.errb)
	}
	if existeRem(processo.CaminhoNoCofre(vault, "legoloop-b")) || existeRem(filepath.Join(vault, "skills", "legoloop-b")) {
		t.Error("b sobrou no cofre")
	}
	if !existeRem(processo.CaminhoNoCofre(vault, "legoloop-a")) || !existeRem(filepath.Join(home, ".claude", "skills", "legoloop-a")) {
		t.Error("a deve ficar até ser removido")
	}

	r = rodar("", "remover", "legoloop-a", "--yes")
	if r.code != 0 || !strings.Contains(r.out, "Vínculos desfeitos: 2") {
		t.Fatalf("a confirmado: exit %d\n%s\n%s", r.code, r.out, r.errb)
	}
	for _, p := range []string{processo.CaminhoNoCofre(vault, "legoloop-a"), filepath.Join(vault, "skills", "legoloop-a"),
		filepath.Join(home, ".claude", "skills", "legoloop-a"), filepath.Join(home, ".codex", "skills", "legoloop-a")} {
		if existeRem(p) {
			t.Errorf("sobrou %s", p)
		}
	}
	if !existeRem(filepath.Join(vault, "loops", "legoloop-velho.loop.json")) {
		t.Error("remover apagou um processo que não era o pedido")
	}
}

// TestRemoverUso cobre o uso incorreto e o processo que não existe, em processo, com HOME descartável.
func TestRemoverUso(t *testing.T) {
	home, vault := cofreListar(t)
	antes := fotoRemover(t, home)
	casos := []struct {
		nome string
		args []string
		code int
		erro string
	}{
		{"sem nome", []string{"remover"}, 2, "Uso: legoloop remover"},
		{"argumento a mais", []string{"remover", "legoloop-b", "extra"}, 2, "Uso: legoloop remover"},
		{"nome fora do padrão", []string{"remover", "../x", "--yes"}, 1, "não é de um processo"},
		{"não existe", []string{"remover", "legoloop-zzz", "--yes"}, 1, "não tem o processo legoloop-zzz"},
		{"sem terminal e sem --yes", []string{"remover", "legoloop-b"}, 1, "Use --yes"},
	}
	for _, c := range casos {
		var out, errb bytes.Buffer
		code := Run(context.Background(), append(c.args, "--vault", vault), Deps{Stdout: &out, Stderr: &errb, Home: home})
		if code != c.code || !strings.Contains(errb.String(), c.erro) {
			t.Errorf("%s: exit %d (quer %d); stderr %q", c.nome, code, c.code, errb.String())
		}
		if fotoRemover(t, home) != antes {
			t.Errorf("%s: mudou o disco", c.nome)
		}
	}
}

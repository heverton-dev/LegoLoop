package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// comHarnesses cria as pastas dos harnesses no HOME descartável do rig.
func comHarnesses(t *testing.T, r *interRig, nomes ...string) {
	t.Helper()
	for _, n := range nomes {
		if err := os.MkdirAll(filepath.Join(r.d.Home, "."+n), 0o755); err != nil {
			t.Fatal(err)
		}
	}
}

func existe(p string) bool { _, err := os.Lstat(p); return err == nil }

// AC-E02.13: o cartão único, as três decisões, as pausas extras e o --yes.
func TestCartaoPlanoUnicoCLI(t *testing.T) {
	t.Run("o cartão lista as ações e as três decisões", func(t *testing.T) {
		r := newInterRig(t, "\n1\n0\n", lowReply, doneReply)
		comHarnesses(t, r, "claude", "gemini")
		r.fabricar(t)
		for _, want := range []string{"AÇÃO:", "grava ", "cria o link", "claude, gemini", "SE DER ERRADO:", "[1] seguir  [2] personalizar  [0] cancelar"} {
			if !strings.Contains(r.out.String(), want) {
				t.Errorf("o cartão não tem %q:\n%s", want, r.out)
			}
		}
	})

	t.Run("cancelar deixa a árvore idêntica", func(t *testing.T) {
		r := newInterRig(t, "\n1\n0\n", lowReply, doneReply)
		comHarnesses(t, r, "claude")
		if code := r.fabricar(t); code != 1 || r.wrote(t) != 0 || existe(filepath.Join(r.d.Home, ".claude", "skills")) {
			t.Fatalf("exit %d, gravou %d", code, r.wrote(t))
		}
	})

	t.Run("Enter vazio no cartão cancela", func(t *testing.T) {
		r := newInterRig(t, "\n1\n\n", lowReply, doneReply)
		if code := r.fabricar(t); code != 1 || r.wrote(t) != 0 {
			t.Fatalf("exit %d, gravou %d", code, r.wrote(t))
		}
	})

	t.Run("personalizar troca os harnesses e volta ao cartão", func(t *testing.T) {
		r := newInterRig(t, "\n1\n2\ngemini\n1\ns\n", lowReply, doneReply)
		comHarnesses(t, r, "claude", "gemini")
		if code := r.fabricar(t); code != 0 {
			t.Fatalf("exit %d: %s", code, r.errb)
		}
		if n := strings.Count(r.out.String(), "AÇÃO:"); n != 2 {
			t.Errorf("o cartão deve aparecer de novo depois de personalizar (%d vezes)", n)
		}
		if !existe(filepath.Join(r.d.Home, ".gemini", "skills", "legoloop-novo-fluxo")) || existe(filepath.Join(r.d.Home, ".claude", "skills")) {
			t.Error("só o gemini devia ser ligado")
		}
	})

	t.Run("recusar a pausa extra dos harnesses não grava nada", func(t *testing.T) {
		r := newInterRig(t, "\n1\n1\nn\n", lowReply, doneReply)
		comHarnesses(t, r, "claude")
		if code := r.fabricar(t); code != 1 || r.wrote(t) != 0 || existe(filepath.Join(r.d.Home, ".claude", "skills")) {
			t.Fatalf("exit %d, gravou %d", code, r.wrote(t))
		}
		if !strings.Contains(r.out.String(), "Ligar a skill a 1 harnesses?") {
			t.Errorf("faltou a pausa extra:\n%s", r.out)
		}
	})

	t.Run("sem harness a pausa dos harnesses não existe", func(t *testing.T) {
		r := newInterRig(t, "\n1\n1\n", lowReply, doneReply)
		if code := r.fabricar(t); code != 0 || strings.Contains(r.out.String(), "Ligar a skill") {
			t.Fatalf("exit %d:\n%s", code, r.out)
		}
	})

	t.Run("--yes imprime o cartão e segue só no que ele mostrou", func(t *testing.T) {
		r := newFabRig(t)
		if err := os.MkdirAll(filepath.Join(r.d.Home, ".claude"), 0o755); err != nil {
			t.Fatal(err)
		}
		code := run(t, r.d, r.args("novo processo", "--yes", "--sem-entrevista")...)
		out := r.out.String()
		if code != 0 || !strings.Contains(out, "AÇÃO:") || !strings.Contains(out, "Segue sem perguntar") || strings.Contains(out, "Escolha:") {
			t.Fatalf("exit %d:\n%s\n%s", code, out, r.errb)
		}
		if !existe(filepath.Join(r.d.Home, ".claude", "skills", "legoloop-novo-processo")) {
			t.Error("o link mostrado no cartão devia existir")
		}
	})

	t.Run("--yes --no-link mostra no cartão que não liga", func(t *testing.T) {
		r := newFabRig(t)
		if err := os.MkdirAll(filepath.Join(r.d.Home, ".claude"), 0o755); err != nil {
			t.Fatal(err)
		}
		code := run(t, r.d, r.args("novo processo", "--yes", "--sem-entrevista", "--no-link")...)
		if code != 0 || !strings.Contains(r.out.String(), "não liga a skill a nenhum harness") || existe(filepath.Join(r.d.Home, ".claude", "skills")) {
			t.Fatalf("exit %d:\n%s", code, r.out)
		}
	})
}

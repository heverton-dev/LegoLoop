package cli

import (
	"strings"
	"testing"
)

func TestDispatchFabricar(t *testing.T) {
	t.Run("help lista fabricar e canonize como comandos distintos", func(t *testing.T) {
		d, out, _ := deps(t)
		if code := run(t, d, "help"); code != 0 {
			t.Fatalf("exit %d", code)
		}
		for _, want := range []string{"fabricar", "canonize"} {
			if !strings.Contains(out.String(), want) {
				t.Errorf("help sem %q:\n%s", want, out)
			}
		}
	})

	t.Run("fabricar --help lista todas as flags e sai com 0", func(t *testing.T) {
		d, out, _ := deps(t)
		if code := run(t, d, "fabricar", "--help"); code != 0 {
			t.Fatalf("exit %d", code)
		}
		for _, flagName := range []string{"yes", "goal", "model", "category", "vault", "docs", "force", "dry-run", "trust-removable", "link-harness", "json"} {
			if !strings.Contains(out.String(), "-"+flagName) {
				t.Errorf("--help sem a flag -%s:\n%s", flagName, out)
			}
		}
		// o pacote flag toma o texto entre crases na descrição como o NOME do valor: não pode vazar para o --help
		if strings.Contains(out.String(), "-model legoloop") || !strings.Contains(out.String(), "-model string") {
			t.Errorf("--help com placeholder estranho em -model:\n%s", out)
		}
	})

	t.Run("fabricar sem argumento imprime o uso e sai com 2", func(t *testing.T) {
		d, out, errb := deps(t)
		if code := run(t, d, "fabricar"); code != 2 {
			t.Fatalf("exit %d, esperado 2", code)
		}
		if !strings.Contains(errb.String(), "Uso: legoloop fabricar") || out.Len() != 0 {
			t.Errorf("uso deve ir ao stderr: stdout=%q stderr=%q", out, errb)
		}
	})

	t.Run("flag desconhecida em fabricar sai com 2", func(t *testing.T) {
		d, _, _ := deps(t)
		if code := run(t, d, "fabricar", "x", "--explodir"); code != 2 {
			t.Fatalf("exit %d, esperado 2", code)
		}
	})

	t.Run("canonize continua exigindo --root e a mesma mensagem", func(t *testing.T) {
		d, _, errb := deps(t)
		if code := run(t, d, "canonize"); code != 2 {
			t.Fatalf("exit %d, esperado 2", code)
		}
		if !strings.Contains(errb.String(), "erro: informe --root e exatamente um de --dry-run ou --apply") {
			t.Errorf("mensagem do legado mudou:\n%s", errb)
		}
	})

	t.Run("canonize nao herda as flags do fabricar (sem despacho por flags)", func(t *testing.T) {
		for _, f := range []string{"--yes", "--goal=x", "--model=m", "--link-harness=claude"} {
			d, _, _ := deps(t)
			if code := run(t, d, "canonize", "--root", t.TempDir(), "--dry-run", f); code != 2 {
				t.Errorf("canonize aceitou %s (exit %d): o legado não pode despachar por flags", f, code)
			}
		}
	})

	t.Run("fabricar nao herda as flags do canonize", func(t *testing.T) {
		d, _, _ := deps(t)
		if code := run(t, d, "fabricar", "x", "--root", "C:\\", "--apply"); code != 2 {
			t.Errorf("fabricar aceitou flags do canonize (exit %d)", code)
		}
	})
}

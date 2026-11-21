package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/caminhos"
)

func TestRenomearCLI(t *testing.T) {
	t.Run("colisao sai com exit 2", func(t *testing.T) {
		raiz := t.TempDir()
		for _, nome := range []string{"velho", "legoloop-velho"} {
			dir := filepath.Join(raiz, caminhos.DirProcessos, nome)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, caminhos.ArquivoDAG), []byte(`{"name":"`+nome+`","nodes":[]}`), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		d, _, errb := deps(t)
		code := run(t, d, "renomear", "velho", "legoloop-velho", "--raiz", raiz, "--yes")
		if code != 2 || !strings.Contains(errb.String(), "já existe") {
			t.Fatalf("colisão: esperava exit 2 com 'já existe', obteve %d (stderr=%s)", code, errb.String())
		}
	})

	t.Run("nome hostil sai com exit 2", func(t *testing.T) {
		raiz := t.TempDir()
		for _, hostil := range []string{"../x", "a/b", `..\x`} {
			d, _, errb := deps(t)
			code := run(t, d, "renomear", hostil, "legoloop-novo", "--raiz", raiz, "--yes")
			if code != 2 || !strings.Contains(errb.String(), "jail") {
				t.Fatalf("hostil %q: esperava exit 2 com erro de jail, obteve %d (stderr=%s)", hostil, code, errb.String())
			}
		}
	})

	t.Run("inexistente sai com exit 2", func(t *testing.T) {
		raiz := t.TempDir()
		d, _, errb := deps(t)
		code := run(t, d, "renomear", "nao-existe", "legoloop-novo", "--raiz", raiz, "--yes")
		if code != 2 || !strings.Contains(errb.String(), "não encontrado") {
			t.Fatalf("inexistente: esperava exit 2 com 'não encontrado', obteve %d (stderr=%s)", code, errb.String())
		}
	})

	t.Run("exec de processo renomeado sai com exit 2 e aponta novo nome", func(t *testing.T) {
		raiz := t.TempDir()
		procDir := filepath.Join(raiz, caminhos.DirProcessos, "velho")
		if err := os.MkdirAll(procDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(procDir, caminhos.ArquivoDAG), []byte(processoScript), 0o644); err != nil {
			t.Fatal(err)
		}
		d, out, errb := deps(t)
		code := run(t, d, "renomear", "velho", "legoloop-novo", "--raiz", raiz, "--yes", "--json")
		if code != 0 {
			t.Fatalf("renomear falhou: code=%d out=%s err=%s", code, out.String(), errb.String())
		}
		d2, _, errb2 := deps(t)
		codeExec := run(t, d2, "exec", "velho",
			"--processos-dir", filepath.Join(raiz, caminhos.DirProcessos),
			"--workdir", t.TempDir())
		if codeExec != 2 {
			t.Fatalf("exec renomeado: esperava exit 2, obteve %d (stderr=%s)", codeExec, errb2.String())
		}
		if !strings.Contains(errb2.String(), "renomeado para legoloop-novo") {
			t.Fatalf("mensagem sem 'renomeado para legoloop-novo': %s", errb2.String())
		}
	})
}

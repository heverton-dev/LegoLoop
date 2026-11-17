package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/caminhos"
)

// TestExecProcesso (AC-E07.1): `legoloop exec <slug>` executa o processo pelo
// slug com exit 0 e artefatos de saída no disco, sem passo manual.
func TestExecProcesso(t *testing.T) {
	base := t.TempDir()
	procDir := filepath.Join(base, caminhos.DirProcessos, "meu-proc")
	if err := os.MkdirAll(procDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(caminhos.CaminhoDAG(procDir), []byte(processoScript), 0o644); err != nil {
		t.Fatal(err)
	}
	work := t.TempDir()
	d, out, errb := deps(t)
	code := run(t, d, "exec", "meu-proc",
		"--processos-dir", filepath.Join(base, caminhos.DirProcessos),
		"--workdir", work, "--allow-scripts")
	if code != 0 {
		t.Fatalf("exec code=%d out=%s err=%s", code, out.String(), errb.String())
	}
	if !strings.Contains(out.String(), "passed") {
		t.Fatalf("relatório sem status: %s", out.String())
	}
	if b, _ := os.ReadFile(filepath.Join(work, "GO_VERSION.txt")); !strings.Contains(string(b), "go1.") {
		t.Fatal("artefato de saída não gravado no disco")
	}
}

// TestExecProcessoErros: slug ausente é uso incorreto (exit 2); slug
// inexistente falha (exit 1) apontando o slug.
func TestExecProcessoErros(t *testing.T) {
	d, _, errb := deps(t)
	if code := run(t, d, "exec"); code != 2 {
		t.Fatalf("sem slug: esperava exit 2, obteve %d", code)
	}
	d2, _, errb2 := deps(t)
	if code := run(t, d2, "exec", "nao-existe", "--processos-dir", t.TempDir()); code != 1 {
		t.Fatalf("slug inexistente: esperava exit 1, obteve %d", code)
	}
	if !strings.Contains(errb2.String(), "nao-existe") {
		t.Fatalf("erro não aponta o slug: %s", errb2.String())
	}
	_ = errb
}

// TestExecProcessoJSON: a saída --json do exec traz o resultado do processo.
func TestExecProcessoJSON(t *testing.T) {
	base := t.TempDir()
	procDir := filepath.Join(base, caminhos.DirProcessos, "meu-proc")
	if err := os.MkdirAll(procDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(caminhos.CaminhoDAG(procDir), []byte(processoScript), 0o644); err != nil {
		t.Fatal(err)
	}
	d, out, _ := deps(t)
	code := run(t, d, "exec", "meu-proc",
		"--processos-dir", filepath.Join(base, caminhos.DirProcessos),
		"--workdir", t.TempDir(), "--allow-scripts", "--json")
	if code != 0 {
		t.Fatalf("exec --json code=%d out=%s", code, out.String())
	}
	if !strings.Contains(out.String(), `"ok"`) {
		t.Fatalf("json sem resultado: %s", out.String())
	}
}

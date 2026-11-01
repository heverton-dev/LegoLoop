package assets

import (
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/gate"
)

func TestScanScripts(t *testing.T) {
	cfg, home, project := testCfg(t)
	write(t, filepath.Join(home, ".aidd", "scripts", "limpa.ps1"), "Write-Host ok")
	write(t, filepath.Join(project, "scripts", "build.sh"), "echo ok")
	write(t, filepath.Join(project, "tools", "gera.py"), "print(1)")
	write(t, filepath.Join(project, "scripts", "leia.txt"), "não é script")
	write(t, filepath.Join(project, "scripts", "node_modules", "x", "z.js"), "1")
	write(t, filepath.Join(project, "scripts", "a", "b", "c", "d", "fundo.sh"), "1")
	items := scanKind(t, cfg, KindScript)
	if s := mustFind(t, items, "limpa.ps1"); s.Scope != "global" {
		t.Errorf("limpa.ps1 = %+v", s)
	}
	if s := mustFind(t, items, "build.sh"); s.Scope != "project" {
		t.Errorf("build.sh = %+v", s)
	}
	mustFind(t, items, "gera.py")
	for _, n := range []string{"leia.txt", "z.js", "fundo.sh"} {
		if _, ok := find(items, n); ok {
			t.Errorf("%s não deveria ser inventariado", n)
		}
	}
}

func TestScanGatesFromBuiltins(t *testing.T) {
	cfg, _, _ := testCfg(t)
	items := scanKind(t, cfg, KindGate)
	if len(items) != len(gate.Catalog()) {
		t.Fatalf("gates = %d, quer %d", len(items), len(gate.Catalog()))
	}
	g := mustFind(t, items, "file_exists")
	if g.Scope != "builtin" || g.Meta["signature"] != "file_exists(caminho)" || g.Description == "" {
		t.Errorf("file_exists = %+v", g)
	}
}

func TestScanHarnessesAndBinaries(t *testing.T) {
	cfg, home, _ := testCfg(t, "claude", "git", "go")
	if err := os.MkdirAll(filepath.Join(home, ".codex"), 0o755); err != nil {
		t.Fatal(err)
	}
	harnesses := scanKind(t, cfg, KindHarness)
	c := mustFind(t, harnesses, "claude")
	if c.Meta["cli"] != "sim" || c.Path != filepath.Join(home, "bin", "claude") {
		t.Errorf("claude = %+v", c)
	}
	x := mustFind(t, harnesses, "codex")
	if x.Meta["cli"] != "não" || x.Meta["config"] != "sim" {
		t.Errorf("codex (só config) = %+v", x)
	}
	if _, ok := find(harnesses, "gemini"); ok {
		t.Error("gemini não existe nem no PATH nem em disco")
	}
	bins := scanKind(t, cfg, KindBinary)
	for _, n := range []string{"git", "go"} {
		mustFind(t, bins, n)
	}
	if _, ok := find(bins, "docker"); ok {
		t.Error("docker não está no PATH falso")
	}
}

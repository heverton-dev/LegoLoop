package main

import (
	"legoloop/internal/caminhos"
	"os"
	"path/filepath"
	"testing"
)

func TestShouldIgnore(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{".git/config", true},
		{"bin/legoloop.exe", true},
		{"vendor/pkg/foo.go", true},
		{"docs/sessions/sess.md", true},
		{".codebase-memory/graph.json", true},
		{"internal/taxonomy/taxonomy_test.go", true},
		{"internal/studio/app.go", false},
		{"docs/templates/TEMPLATE_DAG.json", false},
	}

	for _, tt := range tests {
		got := shouldIgnore(tt.path)
		if got != tt.want {
			t.Errorf("shouldIgnore(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestScanFileWithViolations(t *testing.T) {
	tmpDir := t.TempDir()
	badFile := filepath.Join(tmpDir, "bad.md")
	content := "# Minha esteira\n## Minha task\n"
	if err := os.WriteFile(badFile, []byte(content), 0644); err != nil {
		t.Fatalf("falha ao criar arquivo de teste: %v", err)
	}

	results, err := scanFile(badFile)
	if err != nil {
		t.Fatalf("scanFile retornou erro: %v", err)
	}
	if len(results) < 2 {
		t.Fatalf("esperava ao menos 2 violações, obteve %d", len(results))
	}
}

func writeFixture(t *testing.T, name, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("falha ao criar fixture: %v", err)
	}
	return p
}

func termsOf(t *testing.T, path string) map[string]int {
	t.Helper()
	res, err := scanFile(path)
	if err != nil {
		t.Fatalf("scanFile(%s): %v", path, err)
	}
	m := map[string]int{}
	for _, v := range res {
		m[v.Violation.ForbiddenTerm]++
	}
	return m
}

func TestIdentificadoresGo(t *testing.T) {
	src := "package x\n\ntype Loop struct{}\n\ntype LoopID string\n\nvar loopCount int\n\nvar loop_end int\n"
	got := termsOf(t, writeFixture(t, "id.go", src))
	if got["loop"] != 4 {
		t.Fatalf("esperava 4 violações de loop em identificadores, obteve %v", got)
	}
}

func TestChavesCompostasELiteraisGo(t *testing.T) {
	src := "package x\n\nvar a = \"loop_end\"\nvar b = \"subloop\"\nvar c = \"loops/x.loop.json\"\nvar d = \"" + caminhos.ArquivoDAG + "\"\n"
	got := termsOf(t, writeFixture(t, "lit.go", src))
	if got["loop"] != 4 || got["pipeline"] != 1 {
		t.Fatalf("esperava loop=4 e pipeline=1, obteve %v", got)
	}
}

func TestValorJSONEChaveComposta(t *testing.T) {
	got := termsOf(t, writeFixture(t, "v.json", "{\n \"type\": \"loop_end\",\n \"nome\": \"minha esteira\"\n}\n"))
	if got["loop"] != 1 || got["esteira"] != 1 {
		t.Fatalf("esperava loop=1 e esteira=1, obteve %v", got)
	}
}

func TestCorpoMarkdownETyp(t *testing.T) {
	md := termsOf(t, writeFixture(t, "c.md", "# Título\n\nTexto com uma esteira no corpo.\n\n```\nloop\n```\n"))
	if md["esteira"] != 1 || md["loop"] != 1 {
		t.Fatalf("markdown: obteve %v", md)
	}
	typ := termsOf(t, writeFixture(t, "c.typ", "= Título\nUm pipeline aqui.\n"))
	if typ["pipeline"] != 1 {
		t.Fatalf("typ: obteve %v", typ)
	}
}

func TestMarcaLegoLoopNaoViola(t *testing.T) {
	go1 := "package x\n\n// LegoLoop é a marca.\nvar LegoLoop = \"legoloop\"\nvar legoLoopID = \"LegoLoop\"\n"
	if got := termsOf(t, writeFixture(t, "m.go", go1)); len(got) != 0 {
		t.Fatalf("go: marca reportada: %v", got)
	}
	if got := termsOf(t, writeFixture(t, "m.md", "# LegoLoop\n\nO LegoLoop é a marca.\n")); len(got) != 0 {
		t.Fatalf("md: marca reportada: %v", got)
	}
	if got := termsOf(t, writeFixture(t, "m.json", "{\"LegoLoop\": \"legoloop-x\"}\n")); len(got) != 0 {
		t.Fatalf("json: marca reportada: %v", got)
	}
}

func TestPalavrasLegitimasNaoViolam(t *testing.T) {
	src := "package x\n\nvar discard, cardinal, multitask = 1, 2, 3\nvar phase = \"phase\"\n"
	if got := termsOf(t, writeFixture(t, "ok.go", src)); len(got) != 0 {
		t.Fatalf("falso positivo: %v", got)
	}
}

func TestParseArgsMax(t *testing.T) {
	targets, max, err := parseArgs([]string{"-max", "38", "docs/templates", "skills"})
	if err != nil || max != 38 || len(targets) != 2 || targets[0] != "docs/templates" {
		t.Fatalf("targets = %v, max = %d, err = %v", targets, max, err)
	}
	targets, max, err = parseArgs([]string{"internal"})
	if err != nil || max != 0 || len(targets) != 1 {
		t.Fatalf("sem -max: targets = %v, max = %d, err = %v", targets, max, err)
	}
	if _, _, err := parseArgs([]string{"-max", "-3", "x"}); err == nil {
		t.Fatal("limite negativo deveria falhar")
	}
}

func TestExitCode(t *testing.T) {
	cases := []struct{ n, max, want int }{{0, 0, 0}, {1, 0, 1}, {38, 38, 0}, {39, 38, 1}}
	for _, c := range cases {
		if got := exitCode(c.n, c.max); got != c.want {
			t.Errorf("exitCode(%d,%d) = %d, esperado %d", c.n, c.max, got, c.want)
		}
	}
}

func TestAllowedExceptions(t *testing.T) {
	if !isAllowedException("internal/orchestrix/unit_claude.go", "task") {
		t.Fatal("esperava true para unit_claude.go e task")
	}
	if isAllowedException("internal/orchestrix/unit_claude.go", "loop") {
		t.Fatal("esperava false para unit_claude.go e loop")
	}
	if isAllowedException("internal/cli/outro.go", "task") {
		t.Fatal("esperava false para internal/cli/outro.go e task")
	}
}

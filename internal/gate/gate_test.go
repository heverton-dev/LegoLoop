package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newEnv(t *testing.T, files map[string]string, exitCode int) Env {
	t.Helper()
	dir := t.TempDir()
	for name, body := range files {
		p := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { root.Close() })
	return Env{Root: root, ExitCode: exitCode}
}

func evalExpr(t *testing.T, expr string, env Env) bool {
	t.Helper()
	e, err := Compile(expr)
	if err != nil {
		t.Fatalf("Compile(%q): %v", expr, err)
	}
	ok, err := e.Eval(env)
	if err != nil {
		t.Fatalf("Eval(%q): %v", expr, err)
	}
	return ok
}

func TestEvalTruthTable(t *testing.T) {
	words := strings.Repeat("palavra ", 120)
	env := newEnv(t, map[string]string{
		"BRIEFING.md": words,
		"SPEC.md":     "# Spec\n\n## Critérios de Aceitação\n- [ ] login retorna 200\n- [x] logout\n",
		"NOTES.md":    "sem criterios\n",
		"sub/a.json":  `{"ok":true}`,
	}, 0)
	cases := []struct {
		expr string
		want bool
	}{
		{"file_exists('BRIEFING.md')", true},
		{"file_exists('NADA.md')", false},
		{"file_exists('sub/a.json')", true},
		{"word_count('BRIEFING.md') > 100", true},
		{"word_count('BRIEFING.md') >= 121", false},
		{"line_count('SPEC.md') == 5", true},
		{"file_exists('BRIEFING.md') && word_count('BRIEFING.md') > 100", true},
		{"file_exists('NADA.md') || exit_code == 0", true},
		{"!file_exists('NADA.md')", true},
		{"has_pattern('SPEC.md', 'login retorna \\d+')", true},
		{"has_pattern('NOTES.md', '^#')", false},
		{"has_binary_criteria('SPEC.md')", true},
		{"has_binary_criteria('NOTES.md')", false},
		{"has_acceptance_criteria('SPEC.md')", true},
		{"has_acceptance_criteria('NOTES.md')", false},
		{"exit_code == 0", true},
		{"exit_code != 0", false},
		{"(exit_code < 1) && !(word_count('NOTES.md') <= 1)", true},
		{"true && !false", true},
	}
	for _, c := range cases {
		if got := evalExpr(t, c.expr, env); got != c.want {
			t.Errorf("%q = %v, want %v", c.expr, got, c.want)
		}
	}
}

func TestExitCodeFromEnv(t *testing.T) {
	env := newEnv(t, nil, 3)
	if evalExpr(t, "exit_code == 0", env) {
		t.Fatal("exit_code 3 must fail exit_code == 0")
	}
}

func TestMissingFileCountsAreErrors(t *testing.T) {
	env := newEnv(t, nil, 0)
	e, err := Compile("word_count('NADA.md') > 1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := e.Eval(env); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestPathTraversalBlocked(t *testing.T) {
	outer := t.TempDir()
	secret := filepath.Join(outer, "secret.txt")
	if err := os.WriteFile(secret, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	env := newEnv(t, nil, 0)
	for _, p := range []string{"../secret.txt", secret, "..\\secret.txt"} {
		expr := "file_exists('" + strings.ReplaceAll(p, "\\", "\\\\") + "')"
		e, err := Compile(expr)
		if err != nil {
			continue // rejeitado na compilação também é aceitável
		}
		ok, err := e.Eval(env)
		// No Linux "..\\x" é um nome literal dentro da raiz; só exige erro quando escapa de fato.
		escapes := !filepath.IsLocal(filepath.FromSlash(p))
		if ok || (escapes && err == nil) {
			t.Errorf("%q escaped root: ok=%v err=%v", p, ok, err)
		}
	}
}

func TestCompileErrors(t *testing.T) {
	bad := []string{
		"",
		"file_exists(",
		"unknown_fn('a')",
		"file_exists('a') &&",
		"exit_code ==",
		"word_count('a', 'b') > 1",
		"has_pattern('a', '[')",
		"'texto solto'",
		"exit_code > 'a'",
		"1 2",
		"file_exists('a') @ 1",
		"'sem fim",
		"word_count('a')",
		"exit_code",
		"(true",
	}
	for _, expr := range bad {
		if _, err := Compile(expr); err == nil {
			t.Errorf("Compile(%q) expected error", expr)
		}
	}
}

func TestExprString(t *testing.T) {
	e, err := Compile("  exit_code == 0 ")
	if err != nil {
		t.Fatal(err)
	}
	if e.String() != "exit_code == 0" {
		t.Fatalf("String() = %q", e.String())
	}
}

// SEC-08: expressões gigantes ou muito aninhadas são recusadas antes do parse recursivo.
func TestCompileLimits(t *testing.T) {
	deep := strings.Repeat("(", MaxDepth+1) + "true" + strings.Repeat(")", MaxDepth+1)
	if _, err := Compile(deep); err == nil {
		t.Fatal("aninhamento excessivo aceito")
	}
	ok := strings.Repeat("(", MaxDepth-1) + "true" + strings.Repeat(")", MaxDepth-1)
	if _, err := Compile(ok); err != nil {
		t.Fatalf("aninhamento no limite rejeitado: %v", err)
	}
	if _, err := Compile(strings.Repeat("!", MaxDepth+1) + "true"); err == nil {
		t.Fatal("cadeia de '!' excessiva aceita")
	}
	if _, err := Compile("true" + strings.Repeat(" ", MaxExprBytes)); err == nil {
		t.Fatal("expressão acima do limite de bytes aceita")
	}
}

// SEC-09: diretório não é lido como arquivo; arquivo acima do limite é recusado mesmo após Stat.
func TestReadOnlyRegularFiles(t *testing.T) {
	env := newEnv(t, map[string]string{"d/x.md": "a"}, 0)
	e, _ := Compile("word_count('d') > 0")
	if _, err := e.Eval(env); err == nil {
		t.Fatal("diretório lido como arquivo")
	}
}

func TestSchemaValidBuiltin(t *testing.T) {
	schema := `{
		"type": "object",
		"required": ["id", "status"],
		"properties": {
			"id": {"type": "string"},
			"status": {"type": "string"},
			"count": {"type": "number"}
		}
	}`
	env := newEnv(t, map[string]string{
		"schema.json":           schema,
		"corrupted_schema.json": `{"type": "object",`,
		"empty_schema.json":     ``,
		"valid.json":            `{"id": "node-1", "status": "passed", "count": 10}`,
		"missing.json":          `{"id": "node-1", "count": 10}`,
		"wrong_type.json":       `{"id": 999, "status": "passed"}`,
		"bad_json.json":         `{"id": "corrompido",`,
	}, 0)

	cases := []struct {
		expr string
		want bool
	}{
		{"schema_valid('schema.json', 'valid.json')", true},
		{"schema_valid('schema.json', 'missing.json')", false},
		{"schema_valid('schema.json', 'wrong_type.json')", false},
		{"schema_valid('schema.json', 'bad_json.json')", false},
		{"schema_valid('schema.json', 'non_existent.json')", false},
		{"schema_valid('corrupted_schema.json', 'valid.json')", false},
		{"schema_valid('empty_schema.json', 'valid.json')", false},
	}
	for _, tc := range cases {
		e, err := Compile(tc.expr)
		if err != nil {
			t.Fatalf("Compile(%q): %v", tc.expr, err)
		}
		got, err := e.Eval(env)
		if err != nil {
			t.Fatalf("Eval(%q): %v", tc.expr, err)
		}
		if got != tc.want {
			t.Errorf("Eval(%q) = %v, want %v", tc.expr, got, tc.want)
		}
	}

	badCompiles := []string{
		"schema_valid('../outside.json', 'valid.json')",
		"schema_valid('schema.json', '../outside.json')",
	}
	for _, expr := range badCompiles {
		if _, err := Compile(expr); err == nil {
			t.Errorf("Compile(%q) deveria falhar por caminho fora do workdir", expr)
		}
	}
}


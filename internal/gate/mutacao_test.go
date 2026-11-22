package gate

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type stubNode struct {
	k   valKind
	err error
}

func (s stubNode) kind() valKind           { return s.k }
func (s stubNode) eval(Env) (value, error) { return value{}, s.err }

type errPrompter struct{ err error }

func (p errPrompter) Prompt(string) (bool, error) { return false, p.err }

func abrirRootTemp(t *testing.T) (*os.Root, string) {
	t.Helper()
	dir := t.TempDir()
	root, err := os.OpenRoot(dir)
	if err != nil {
		t.Fatalf("OpenRoot: %v", err)
	}
	t.Cleanup(func() { _ = root.Close() })
	return root, dir
}

func gravarArq(t *testing.T, dir, nome string, data []byte) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, nome), data, 0o600); err != nil {
		t.Fatalf("WriteFile(%s): %v", nome, err)
	}
}

func TestMutacaoBuiltinsErrosEArgIndex(t *testing.T) {
	_, err := Compile("file_exists(1)")
	if err == nil || !strings.Contains(err.Error(), "argumento 1") {
		t.Fatalf("esperava erro com 'argumento 1', veio %v", err)
	}
	cn := callNode{name: "file_exists", b: builtins()["file_exists"], args: []node{stubNode{k: kStr, err: errors.New("boom")}}}
	if _, err := cn.eval(Env{}); err == nil {
		t.Fatal("esperava erro propagado de argumento")
	}
	root, _ := abrirRootTemp(t)
	for _, expr := range []string{"has_pattern('falta.txt', 'a')", "has_binary_criteria('falta.txt')", "has_acceptance_criteria('falta.md')"} {
		e, err := Compile(expr)
		if err != nil {
			t.Fatalf("Compile(%s): %v", expr, err)
		}
		if ok, err := e.Eval(Env{Root: root}); err == nil || ok {
			t.Fatalf("%s: esperava (false, err), veio (%v, %v)", expr, ok, err)
		}
	}
}

func TestMutacaoReadGateFileLimitesEDir(t *testing.T) {
	root, dir := abrirRootTemp(t)
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o700); err != nil {
		t.Fatal(err)
	}
	if _, err := readGateFile(Env{Root: root}, "sub"); err == nil || !strings.Contains(err.Error(), "não é arquivo regular") {
		t.Fatalf("esperava erro de não regular, veio %v", err)
	}
	exato := bytes.Repeat([]byte("a"), maxGateFileBytes)
	gravarArq(t, dir, "exato.bin", exato)
	got, err := readGateFile(Env{Root: root}, "exato.bin")
	if err != nil || len(got) != maxGateFileBytes {
		t.Fatalf("esperava %d bytes sem erro, veio len=%d err=%v", maxGateFileBytes, len(got), err)
	}
	gravarArq(t, dir, "grande.bin", append(exato, 'b'))
	if _, err := readGateFile(Env{Root: root}, "grande.bin"); err == nil || !strings.Contains(err.Error(), "excede") {
		t.Fatalf("esperava erro de excesso de tamanho, veio %v", err)
	}
}

func TestMutacaoSchemaValidEInteractive(t *testing.T) {
	root, dir := abrirRootTemp(t)
	gravarArq(t, dir, "alvo.json", []byte(`{}`))
	e, err := Compile("schema_valid('nao_existe.json', 'alvo.json')")
	if err != nil {
		t.Fatal(err)
	}
	if ok, err := e.Eval(Env{Root: root}); err != nil || ok {
		t.Fatalf("esperava (false, nil), veio (%v, %v)", ok, err)
	}
	v, err := interactivePrompt(Env{}, []value{{kind: kStr, s: "ok?"}}, nil)
	if err == nil || v.b {
		t.Fatalf("sem prompter: esperava (false, err), veio (%v, %v)", v.b, err)
	}
	v, err = interactivePrompt(Env{Prompter: errPrompter{err: errors.New("falha")}}, []value{{kind: kStr, s: "ok?"}}, nil)
	if err == nil || v.b {
		t.Fatalf("erro prompter: esperava (false, err), veio (%v, %v)", v.b, err)
	}
}

func TestMutacaoNoStubs(t *testing.T) {
	root, dir := abrirRootTemp(t)
	gravarArq(t, dir, "vazio.go", []byte("   \n\t "))
	gravarArq(t, dir, "todo.go", []byte("func f() { // TODO\n}"))
	gravarArq(t, dir, "limpo.go", []byte("package x\nfunc ok() int { return 1 }\n"))
	casos := []struct {
		arq  string
		want bool
	}{
		{"falta.go", false},
		{"vazio.go", false},
		{"todo.go", false},
		{"limpo.go", true},
	}
	for _, c := range casos {
		e, err := Compile("no_stubs('" + c.arq + "')")
		if err != nil {
			t.Fatalf("Compile(%s): %v", c.arq, err)
		}
		got, err := e.Eval(Env{Root: root})
		if err != nil || got != c.want {
			t.Fatalf("no_stubs(%s) = (%v, %v), quer (%v, nil)", c.arq, got, err, c.want)
		}
	}
}

func TestMutacaoCitesBordas(t *testing.T) {
	root, dir := abrirRootTemp(t)
	gravarArq(t, dir, "vazio.go", nil)
	gravarArq(t, dir, "doc.md", []byte("ver vazio.go:1\n"))
	verdicts, err := CiteVerdicts(root, "doc.md", ".")
	if err != nil || len(verdicts) != 1 || verdicts[0].InRange {
		t.Fatalf("arquivo vazio citado: esperava 1 verdict fora do intervalo, veio %+v err=%v", verdicts, err)
	}
	if ids := lineIdents("`Alpha` e `Alpha`"); len(ids) != 1 || ids[0] != "Alpha" {
		t.Fatalf("lineIdents duplicado: veio %v", ids)
	}
	if _, err := CiteVerdicts(root, "inexistente.md", "."); err == nil {
		t.Fatal("esperava erro ao ler artefato inexistente")
	}
	e, err := Compile("cites_exist('inexistente.md', '.')")
	if err != nil {
		t.Fatal(err)
	}
	if ok, err := e.Eval(Env{Root: root}); err == nil || ok {
		t.Fatalf("cites_exist inexistente: esperava (false, err), veio (%v, %v)", ok, err)
	}
}

func TestMutacaoGateLimitesEErros(t *testing.T) {
	root, _ := abrirRootTemp(t)
	exprExata := strings.Repeat(" ", MaxExprBytes-4) + "true"
	e, err := Compile(exprExata)
	if err != nil {
		t.Fatalf("expressão com exatamente MaxExprBytes deve compilar: %v", err)
	}
	if ok, err := e.Eval(Env{}); err == nil || ok {
		t.Fatalf("Env.Root nil: esperava (false, err), veio (%v, %v)", ok, err)
	}
	erros := []string{
		"!has_acceptance_criteria('falta.md')",
		"has_acceptance_criteria('falta.md') && true",
		"0 == word_count('falta.txt')",
	}
	for _, src := range erros {
		ex, err := Compile(src)
		if err != nil {
			t.Fatalf("Compile(%s): %v", src, err)
		}
		if ok, err := ex.Eval(Env{Root: root}); err == nil || ok {
			t.Fatalf("Eval(%s): esperava (false, err), veio (%v, %v)", src, ok, err)
		}
	}
}

func TestMutacaoGateLogicaEComparacoes(t *testing.T) {
	root, _ := abrirRootTemp(t)
	casos := []struct {
		src  string
		code int
		want bool
	}{
		{"true && false", 0, false},
		{"false || true", 0, true},
		{"false && has_acceptance_criteria('falta.md')", 0, false},
		{"true || has_acceptance_criteria('falta.md')", 0, true},
		{"exit_code > 5", 5, false},
		{"exit_code >= 5", 5, true},
		{"exit_code < 5", 5, false},
		{"exit_code <= 5", 5, true},
	}
	for _, c := range casos {
		e, err := Compile(c.src)
		if err != nil {
			t.Fatalf("Compile(%s): %v", c.src, err)
		}
		got, err := e.Eval(Env{Root: root, ExitCode: c.code})
		if err != nil || got != c.want {
			t.Fatalf("Eval(%s, %d) = (%v, %v), quer (%v, nil)", c.src, c.code, got, err, c.want)
		}
	}
}

func TestMutacaoLexerEscapes(t *testing.T) {
	if _, err := lex("'\\"); err == nil {
		t.Fatal("string sem fechamento terminando em barra deve retornar erro")
	}
	toks, err := lex(`'a\'b'`)
	if err != nil || len(toks) < 1 || toks[0].text != "a'b" {
		t.Fatalf("escape aspa: toks=%+v err=%v", toks, err)
	}
	toks, err = lex(`'\\'`)
	if err != nil || len(toks) < 1 || toks[0].text != `\` {
		t.Fatalf("escape barra: toks=%+v err=%v", toks, err)
	}
}

func TestMutacaoParserProfundidadeEErros(t *testing.T) {
	okExpr := strings.Repeat("(", MaxDepth) + "true" + strings.Repeat(")", MaxDepth)
	if _, err := Compile(okExpr); err != nil {
		t.Fatalf("profundidade exata MaxDepth deve compilar: %v", err)
	}
	casos := []struct {
		src  string
		want string
	}{
		{strings.Repeat("(", MaxDepth+1) + "true" + strings.Repeat(")", MaxDepth+1), "aninhamento"},
		{"true && 1", "exige operandos bool"},
		{"1 && true", "exige operandos bool"},
		{"(true && 1)", "exige operandos bool"},
		{"file_exists('a'", "esperado nos argumentos"},
		{"file_exists(true && 1)", "exige operandos bool"},
	}
	for _, c := range casos {
		if _, err := Compile(c.src); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Fatalf("Compile(%q): esperava %q, veio %v", c.src, c.want, err)
		}
	}
}

func TestMutacaoSchemaRefsMinimumEItems(t *testing.T) {
	casos := []struct {
		nome   string
		schema string
		alvo   string
		want   bool
	}{
		{"ref inexistente", `{"$defs":{"a":{"type":"string"}},"$ref":"#/$defs/b"}`, `"x"`, false},
		{"ref prefixo invalido", `{"$ref":"#/definitions/a"}`, `"x"`, false},
		{"ref sem defs", `{"$ref":"#/$defs/a"}`, `"x"`, false},
		{"minimum igual", `{"minimum":5}`, `5`, true},
		{"minimum nao numero", `{"minimum":-1}`, `"texto"`, false},
		{"items com null", `{"items":{"type":"string"}}`, `null`, true},
	}
	for _, c := range casos {
		if got := ValidateJSONSchema([]byte(c.schema), []byte(c.alvo)); got != c.want {
			t.Fatalf("%s: got=%v want=%v", c.nome, got, c.want)
		}
	}
	rec := []byte(`{"$defs":{"n":{"type":"array","items":{"$ref":"#/$defs/n"}}},"$ref":"#/$defs/n"}`)
	alvo := []byte(strings.Repeat("[", MaxDepth+2) + strings.Repeat("]", MaxDepth+2))
	if ValidateJSONSchema(rec, alvo) {
		t.Fatal("recursão acima de MaxDepth em items deveria falhar")
	}
}

func TestMutacaoSTE100LimitesESeparadores(t *testing.T) {
	vinte := strings.TrimSpace(strings.Repeat("palavra ", MaxSTE100WordsPerSentence))
	if v, ok := ValidateSTE100(vinte); !ok || len(v) != 0 {
		t.Fatalf("20 palavras exatas devem passar: ok=%v v=%v", ok, v)
	}
	frase := strings.TrimSpace(strings.Repeat("termo ", 12))
	for _, sep := range []string{". ", "! ", "? ", "; "} {
		if v, ok := ValidateSTE100(frase + sep + frase); !ok || len(v) != 0 {
			t.Fatalf("separador %q deveria dividir frases: ok=%v v=%v", sep, ok, v)
		}
	}
}

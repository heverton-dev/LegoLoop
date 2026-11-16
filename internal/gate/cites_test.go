package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// repoCites monta um repositorio de mentira: a.go (3 linhas, termina em \n), sub/b.go (5 linhas, sem \n final).
func repoCites(t *testing.T) (RootFS, string) {
	t.Helper()
	dir := t.TempDir()
	must := func(rel, body string) {
		p := filepath.Join(dir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	must("a.go", "package a\nfunc Bar() {}\nl3\n")
	must("sub/b.go", "1\n2\n3\nfunc Foo() {}\n5")
	root, err := os.OpenRoot(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = root.Close() })
	return root, dir
}

func evalCites(t *testing.T, root RootFS, artifact, base string) bool {
	t.Helper()
	x, err := Compile("cites_exist('out.md', '" + base + "')")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root.(*os.Root).Name(), "out.md"), []byte(artifact), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err := x.Eval(Env{Root: root})
	if err != nil {
		t.Fatalf("Eval: %v", err)
	}
	return ok
}

func TestCitesExist(t *testing.T) {
	root, _ := repoCites(t)
	casos := []struct {
		nome, texto string
		quer        bool
	}{
		{"dois pontos validos", "Prova em `a.go:3` e `sub/b.go:5`.", true},
		{"ultima linha sem newline final", "sub/b.go:5", true},
		{"linha alem do fim", "`a.go:4`", false},
		{"arquivo inexistente", "`zzz.go:1`", false},
		{"forma 'linha N' valida", "Arquivo: `a.go` (linha 3)", true},
		{"forma 'linha N' fora do arquivo", "Arquivo: `a.go` (linha 9)", false},
		{"funcao entre o caminho e a linha", "sub/b.go (funcao `Foo`, linha 4)", true},
		{"intervalo valido", "a.go:2-3", true},
		{"intervalo alem do fim", "a.go:2-9", false},
		{"duas citacoes na mesma linha, a segunda ruim", "a.go:2 e sub/b.go:99", false},
		{"duas citacoes na mesma linha, ambas boas", "a.go:2 e sub/b.go:5", true},
		{"uma ruim entre boas reprova tudo", "a.go:1\nsub/b.go:2\na.go:40\n", false},
		{"linha zero", "a.go:0", false},
		{"sem nenhuma citacao", "Texto sem prova nenhuma.", false},
		{"caminho sem numero nao e citacao", "Veja a.go para detalhes.", false},
		{"travessia de diretorio", "../x.go:1", false},
		{"diretorio no lugar de arquivo", "sub/b.go:1\nsub.go:1", false},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			if got := evalCites(t, root, c.texto, "."); got != c.quer {
				t.Errorf("cites_exist(%q) = %v, quer %v", c.texto, got, c.quer)
			}
		})
	}
}

// Coerencia de nomes: identificador entre crases na mesma linha da citacao precisa existir no arquivo citado.
func TestCitesExistIdentificadores(t *testing.T) {
	root, _ := repoCites(t)
	casos := []struct {
		nome, texto string
		quer        bool
	}{
		{"identificador presente", "`a.go:3` chama `Bar`", true},
		{"identificador ausente reprova", "`a.go:3` usa `Baz`", false},
		{"identificador em outra linha do artefato e ignorado", "`a.go:3`\nCorrecao: use `Baz`", true},
		{"com parenteses e trecho pontuado", "`a.go:3` com `Bar()` e `x.Bar`", true},
		{"pontuado usa o ultimo trecho", "`a.go:3` com `sync.Mutex`", false},
		{"identificador curto e ignorado", "`a.go:3` e `id`", true},
		{"caminho entre crases nao e identificador", "`sub/b.go:4` e `a.go`", true},
		{"vale se estiver em um dos arquivos citados na linha", "a.go:3 e sub/b.go:4 usam `Foo`", true},
		{"palavra inteira, nao pedaco", "`a.go:3` usa `Ba`", true}, // "Ba" tem menos de 3 letras: ignorado
		{"pedaco de palavra nao vale", "`a.go:3` usa `Barr`", false},
		{"texto que nao e identificador e ignorado", "`a.go:3` faz `x + y`", true},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			if got := evalCites(t, root, c.texto, "."); got != c.quer {
				t.Errorf("cites_exist(%q) = %v, quer %v", c.texto, got, c.quer)
			}
		})
	}
}

func TestCiteVerdictsApontaOQueFalhou(t *testing.T) {
	root, dir := repoCites(t)
	if err := os.WriteFile(filepath.Join(dir, "out.md"), []byte("`a.go:3` usa `Baz`\n`a.go:9`\n`sub/b.go:4` com `Foo`\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	vs, err := CiteVerdicts(root, "out.md", ".")
	if err != nil || len(vs) != 3 {
		t.Fatalf("CiteVerdicts: %v %+v", err, vs)
	}
	if !vs[0].InRange || len(vs[0].MissingIdent) != 1 || vs[0].MissingIdent[0] != "Baz" || vs[0].OK() {
		t.Errorf("primeira: existe mas Baz ausente: %+v", vs[0])
	}
	if vs[1].InRange || vs[1].OK() {
		t.Errorf("segunda: linha 9 fora do arquivo: %+v", vs[1])
	}
	if !vs[2].OK() {
		t.Errorf("terceira deveria passar: %+v", vs[2])
	}
}

func TestCitesExistBase(t *testing.T) {
	root, _ := repoCites(t)
	if !evalCites(t, root, "b.go:5", "sub") {
		t.Error("com base=sub, b.go:5 deveria existir")
	}
	if evalCites(t, root, "a.go:1", "sub") {
		t.Error("com base=sub, a.go nao existe")
	}
}

func TestCitesExistLimiteDeCitacoes(t *testing.T) {
	root, _ := repoCites(t)
	if evalCites(t, root, strings.Repeat("a.go:1\n", maxCitacoes+1), ".") {
		t.Error("acima do limite de citacoes o gate deve reprovar")
	}
	if !evalCites(t, root, strings.Repeat("a.go:1\n", maxCitacoes), ".") {
		t.Error("no limite o gate deve aprovar")
	}
}

func TestCitesExistCompilacao(t *testing.T) {
	for _, src := range []string{
		"cites_exist('out.md')",
		"cites_exist('../fora.md', '.')",
		"cites_exist('out.md', '../fora')",
		"cites_exist('out.md', '/abs')",
	} {
		if _, err := Compile(src); err == nil {
			t.Errorf("Compile(%q) deveria falhar", src)
		}
	}
	x, err := Compile("cites_exist('out.md', '.')")
	if err != nil || !x.ChecksContent() {
		t.Errorf("cites_exist deve compilar e provar conteudo: err=%v", err)
	}
}

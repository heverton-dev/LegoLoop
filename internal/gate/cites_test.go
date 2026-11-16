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
	must("a.go", "l1\nl2\nl3\n")
	must("sub/b.go", "1\n2\n3\n4\n5")
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

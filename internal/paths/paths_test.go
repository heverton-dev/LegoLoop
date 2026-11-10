package paths

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// TestResolucaoUnica garante o AC-E01.1: nenhum arquivo Go de produção em
// internal/ declara os literais de nome fora do pacote central, e a fachada
// resolve os mesmos caminhos do pacote central.
func TestResolucaoUnica(t *testing.T) {
	if got, want := DirLoopsEm("v"), filepath.Join("v", DirLoops); got != want {
		t.Errorf("DirLoopsEm = %q, esperado %q", got, want)
	}
	if got, want := CaminhoDAG("d"), filepath.Join("d", ArquivoDAG); got != want {
		t.Errorf("CaminhoDAG = %q, esperado %q", got, want)
	}
	verificaLeituraDual(t)
	achados, err := literaisFora(filepath.Join(".."), PacoteCentral)
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range achados {
		t.Errorf("literal de nome fora de internal/%s: %s", PacoteCentral, a)
	}
}

// verificaLeituraDual confere que a fachada prefere o dialeto novo na leitura
// e cai no antigo quando o novo não existe.
func verificaLeituraDual(t *testing.T) {
	t.Helper()
	base := t.TempDir()
	if got := CaminhoDAGLeitura(base); got != CaminhoDAG(base) {
		t.Errorf("CaminhoDAGLeitura sem arquivo novo = %q, esperado %q", got, CaminhoDAG(base))
	}
	novo := filepath.Join(base, ArquivoProcessoDAG)
	if err := os.WriteFile(novo, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := CaminhoDAGLeitura(base); got != novo {
		t.Errorf("CaminhoDAGLeitura com arquivo novo = %q, esperado %q", got, novo)
	}
	for _, d := range []string{DirLoops, DirProcessos} {
		if err := os.MkdirAll(filepath.Join(base, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	want := []string{filepath.Join(base, DirProcessos), DirLoopsEm(base)}
	got := DirsLeitura(base)
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("DirsLeitura = %v, esperado %v", got, want)
	}
}

// TestResolucaoUnicaDetectaMutacao planta um literal proibido e exige que a
// varredura o reporte. Sem isso, o teste principal passaria vazio.
func TestResolucaoUnicaDetectaMutacao(t *testing.T) {
	raiz := t.TempDir()
	for i, lit := range Literais() {
		dir := filepath.Join(raiz, "plantado"+strconv.Itoa(i))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		src := "package p\n\nconst x = " + strconv.Quote(lit) + "\n"
		if err := os.WriteFile(filepath.Join(dir, "p.go"), []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	achados, err := literaisFora(raiz, PacoteCentral)
	if err != nil {
		t.Fatal(err)
	}
	if len(achados) != len(Literais()) {
		t.Fatalf("esperado %d achados, veio %v", len(Literais()), achados)
	}
}

// literaisFora varre os arquivos Go de produção sob raiz, exceto a pasta
// central, e devolve "arquivo:linha" de cada literal de nome encontrado.
func literaisFora(raiz, central string) ([]string, error) {
	proibidos := map[string]bool{}
	for _, l := range Literais() {
		proibidos[l] = true
	}
	var achados []string
	err := filepath.WalkDir(raiz, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() == central && filepath.Dir(p) == filepath.Clean(raiz) {
			return filepath.SkipDir
		}
		if d.IsDir() || !ehGoProducao(p) {
			return nil
		}
		a, err := literaisNoArquivo(p, proibidos)
		achados = append(achados, a...)
		return err
	})
	return achados, err
}

func ehGoProducao(p string) bool {
	return strings.HasSuffix(p, ".go") && !strings.HasSuffix(p, "_test.go")
}

func literaisNoArquivo(p string, proibidos map[string]bool) ([]string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, p, nil, parser.SkipObjectResolution)
	if err != nil {
		return nil, err
	}
	var achados []string
	ast.Inspect(f, func(n ast.Node) bool {
		lit, ok := n.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return true
		}
		if v, err := strconv.Unquote(lit.Value); err == nil && proibidos[v] {
			achados = append(achados, fset.Position(lit.Pos()).String())
		}
		return true
	})
	return achados, nil
}

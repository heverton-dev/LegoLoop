package caminhos

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

// TestConstantesEFuncoesNoMesmoPacote prova AC-E03.1: constantes de nome e
// funções de caminho convivem no pacote caminhos, com valores e comportamentos
// esperados.
func TestConstantesEFuncoesNoMesmoPacote(t *testing.T) {
	pares := [][2]string{
		{DirProcessos, "processos"},
		{SufixoStudio, ".studio.json"},
		{SufixoAntigo, SufixoLoop},
		{DirAntigo, DirLoops},
		{ArquivoDAG, "pipeline.dag.json"},
		{ArquivoProcessoDAG, "processo.dag.json"},
		{TemplatePecaAntigo, "TEMPLATE_PIECE_LOOP.json"},
	}
	for _, p := range pares {
		if p[0] != p[1] {
			t.Errorf("constante = %q, esperado %q", p[0], p[1])
		}
	}
	if got, want := DirAntigoEm("v"), filepath.Join("v", DirAntigo); got != want {
		t.Errorf("DirAntigoEm = %q, esperado %q", got, want)
	}
	if got, want := CaminhoDAG("d"), filepath.Join("d", ArquivoDAG); got != want {
		t.Errorf("CaminhoDAG = %q, esperado %q", got, want)
	}
	verificaLeituraDual(t)
}

func verificaLeituraDual(t *testing.T) {
	t.Helper()
	base := t.TempDir()
	if got := CaminhoDAGLeitura(base); got != CaminhoDAG(base) {
		t.Errorf("CaminhoDAGLeitura sem novo = %q, esperado %q", got, CaminhoDAG(base))
	}
	novo := filepath.Join(base, ArquivoProcessoDAG)
	if err := os.WriteFile(novo, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := CaminhoDAGLeitura(base); got != novo {
		t.Errorf("CaminhoDAGLeitura com novo = %q, esperado %q", got, novo)
	}
	for _, d := range []string{DirAntigo, DirProcessos} {
		if err := os.MkdirAll(filepath.Join(base, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	want := []string{filepath.Join(base, DirProcessos), DirAntigoEm(base)}
	got := DirsLeitura(base)
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("DirsLeitura = %v, esperado %v", got, want)
	}
}

// TestResolucaoUnica garante que nenhum arquivo Go de produção em internal/
// declara os literais de nome fora de internal/caminhos.
func TestResolucaoUnica(t *testing.T) {
	achados, err := literaisFora(filepath.Join(".."), PacoteCentral)
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range achados {
		t.Errorf("literal de nome fora de internal/%s: %s", PacoteCentral, a)
	}
}

// TestResolucaoUnicaDetectaMutacao planta literais proibidos e exige detecção.
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

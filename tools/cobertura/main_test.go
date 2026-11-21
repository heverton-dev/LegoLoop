package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func criarModuloTeste(t *testing.T, coberto bool) string {
	t.Helper()
	dir := t.TempDir()
	gomod := "module exemplo.local/teste\n\ngo 1.22\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod), 0o600); err != nil {
		t.Fatal(err)
	}
	pkgDir := filepath.Join(dir, "pkg")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	src := `package pkg

func Soma(a, b int) int {
	return a + b
}

func Sub(a, b int) int {
	return a - b
}
`
	if err := os.WriteFile(filepath.Join(pkgDir, "calc.go"), []byte(src), 0o600); err != nil {
		t.Fatal(err)
	}
	teste := `package pkg

import "testing"

func TestSoma(t *testing.T) {
	if Soma(1, 2) != 3 {
		t.Fatal("esperava 3")
	}
}
`
	if coberto {
		teste += `
func TestSub(t *testing.T) {
	if Sub(3, 1) != 2 {
		t.Fatal("esperava 2")
	}
}
`
	}
	if err := os.WriteFile(filepath.Join(pkgDir, "calc_test.go"), []byte(teste), 0o600); err != nil {
		t.Fatal(err)
	}
	limites := `{"./pkg": 80.0}`
	if err := os.WriteFile(filepath.Join(dir, "limites.json"), []byte(limites), 0o600); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestExecutarAbaixoDoMinimoRetornaExit2(t *testing.T) {
	dir := criarModuloTeste(t, false)
	var stdout, stderr bytes.Buffer
	code := executar([]string{"--limites", "limites.json"}, &stdout, &stderr, dir)
	if code != 2 {
		t.Fatalf("esperava exit 2, obtive %d (stdout=%q stderr=%q)", code, stdout.String(), stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "ABAIXO") || !strings.Contains(out, "./pkg") || !strings.Contains(out, "obtida=50.0%") || !strings.Contains(out, "minima=80.0%") {
		t.Fatalf("saida inesperada: %q", out)
	}
}

func TestExecutarAcimaDoMinimoRetornaExit0(t *testing.T) {
	dir := criarModuloTeste(t, true)
	var stdout, stderr bytes.Buffer
	code := executar([]string{"--limites", "limites.json"}, &stdout, &stderr, dir)
	if code != 0 {
		t.Fatalf("esperava exit 0, obtive %d (stdout=%q stderr=%q)", code, stdout.String(), stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "OK") || !strings.Contains(out, "./pkg") || !strings.Contains(out, "obtida=100.0%") {
		t.Fatalf("saida inesperada: %q", out)
	}
}

func TestExecutarArquivoInexistenteRetornaExit2(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := executar([]string{"--limites", "nao_existe.json"}, &stdout, &stderr, t.TempDir())
	if code != 2 {
		t.Fatalf("esperava exit 2, obtive %d", code)
	}
}

func TestExtrairCoberturas(t *testing.T) {
	saida := "ok  \tlegoloop/internal/board\t1.2s\tcoverage: 65.5% of statements\n" +
		"?   \tlegoloop/internal/vazio\t[no test files]\n"
	got := extrairCoberturas(saida)
	if got["legoloop/internal/board"] != 65.5 {
		t.Fatalf("esperava 65.5, obtive %v", got["legoloop/internal/board"])
	}
	if v, ok := got["legoloop/internal/vazio"]; !ok || v != 0 {
		t.Fatalf("esperava 0 para vazio, obtive %v (%v)", v, ok)
	}
}

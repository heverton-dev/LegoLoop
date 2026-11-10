package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func write(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestGofmtFalhaComArquivoSujo(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "ruim.go", "package x\nfunc  F( ){}\n")
	out, err := gofmtStep(config{Root: dir}).Run()
	if err == nil || !strings.Contains(out, "ruim.go") {
		t.Fatalf("esperava falha citando ruim.go, obteve err=%v out=%q", err, out)
	}
}

// soGofmtReal mantém gofmt real e troca os demais subprocessos por sucesso.
func soGofmtReal(t *testing.T) {
	t.Helper()
	orig := execCmd
	t.Cleanup(func() { execCmd = orig })
	execCmd = func(dir, name string, args ...string) (string, error) {
		if name == "gofmt" {
			return orig(dir, name, args...)
		}
		return "", nil
	}
}

func pacote(t *testing.T, arquivo, conteudo string) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "pkg"), 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(root, "pkg"), arquivo, conteudo)
	return root
}

// AC-E00.5: pacote sujo dá exit 1 citando o arquivo; pacote limpo dá exit 0.
func TestGofmtSujoFalha(t *testing.T) {
	soGofmtReal(t)
	sujo := pacote(t, "ruim.go", "package pkg\nfunc  F( ){}\n")
	var out bytes.Buffer
	if code := run([]string{"-root", sujo, "-pkgs", "./pkg/..."}, &out, &out); code != 1 {
		t.Fatalf("pacote sujo: exit=%d, queria 1; saída=%q", code, out.String())
	}
	if !strings.Contains(out.String(), "ruim.go") {
		t.Fatalf("saída deveria citar ruim.go: %q", out.String())
	}
	limpo := pacote(t, "bom.go", "package pkg\n\nfunc F() {}\n")
	out.Reset()
	if code := run([]string{"-root", limpo, "-pkgs", "./pkg/..."}, &out, &out); code != 0 {
		t.Fatalf("pacote limpo: exit=%d, queria 0; saída=%q", code, out.String())
	}
	if code := run([]string{"-inexistente"}, &out, &out); code != 2 {
		t.Fatalf("erro de uso: exit=%d, queria 2", code)
	}
}

func TestGofmtPassaComArquivoLimpo(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "bom.go", "package x\n\nfunc F() {}\n")
	if out, err := gofmtStep(config{Root: dir}).Run(); err != nil {
		t.Fatalf("esperava passar, obteve err=%v out=%q", err, out)
	}
}

func TestGofmtPacotePedidoEntraNaVerificacao(t *testing.T) {
	paths := gofmtPaths(config{Pkgs: []string{"./internal/dag/...", "./cmd"}})
	want := []string{".", "./internal/dag", "./cmd"}
	if strings.Join(paths, "|") != strings.Join(want, "|") {
		t.Fatalf("obteve %v, queria %v", paths, want)
	}
}

func TestTaxComViolacaoFalha(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "a.md", "# Título\n\nUma est"+"eira aqui.\n")
	out, err := taxStep(config{Root: "../..", Tax: []string{dir}}).Run()
	if err == nil {
		t.Fatalf("esperava falha do taxlinter, out=%q", out)
	}
}

func TestTaxSemViolacaoPassa(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "a.md", "# Título\n\nUm Processo limpo.\n")
	if out, err := taxStep(config{Root: "../..", Tax: []string{dir}}).Run(); err != nil {
		t.Fatalf("esperava passar, err=%v out=%q", err, out)
	}
}

func TestLivrosPDFInexistenteFalha(t *testing.T) {
	dir := t.TempDir()
	if _, err := livrosStep(config{Root: dir, Livros: []string{"nao-existe"}}).Run(); err == nil {
		t.Fatal("pasta inexistente deveria falhar")
	}
	if _, err := livrosStep(config{Root: dir, Livros: []string{"."}}).Run(); err == nil {
		t.Fatal("pasta sem PDF deveria falhar")
	}
}

func TestLivrosAuditaCadaPDF(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "a.pdf", "x")
	write(t, dir, "b.pdf", "x")
	orig := execCmd
	defer func() { execCmd = orig }()
	var vistos []string
	execCmd = func(_, name string, args ...string) (string, error) {
		vistos = append(vistos, filepath.Base(args[len(args)-1]))
		if strings.HasSuffix(args[len(args)-1], "b.pdf") {
			return "página órfã", errors.New("exit 1")
		}
		return "", nil
	}
	out, err := livrosStep(config{Root: dir, Livros: []string{"."}, Bin: "x"}).Run()
	if len(vistos) != 2 || err == nil || !strings.Contains(out, "b.pdf") {
		t.Fatalf("vistos=%v err=%v out=%q", vistos, err, out)
	}
}

func TestBuildStepsMontaSoOPedido(t *testing.T) {
	base := len(buildSteps(config{}))
	cheio := len(buildSteps(config{Pkgs: []string{"./x"}, Tax: []string{"x"}, Gold: true, Build: true, Livros: []string{"l"}}))
	if cheio != base+5 {
		t.Fatalf("base=%d cheio=%d", base, cheio)
	}
}

func TestRunStepsContaFalhas(t *testing.T) {
	ok := step{Name: "a", Run: func() (string, error) { return "", nil }}
	ruim := step{Name: "b", Run: func() (string, error) { return "saida", errors.New("x") }}
	var buf bytes.Buffer
	if n := runSteps([]step{ok, ruim, ruim}, &buf); n != 2 {
		t.Fatalf("falhas=%d", n)
	}
	if !strings.Contains(buf.String(), "[FALHA] b") || !strings.Contains(buf.String(), "[ok]    a") {
		t.Fatalf("saída inesperada: %s", buf.String())
	}
}

func TestParseArgs(t *testing.T) {
	c, err := parseArgs([]string{"-pkgs", "./a/...,./b", "-gold", "-livros", "docs/books"}, &bytes.Buffer{})
	if err != nil || len(c.Pkgs) != 2 || !c.Gold || c.Livros[0] != "docs/books" {
		t.Fatalf("c=%+v err=%v", c, err)
	}
	if _, err := parseArgs([]string{"-inexistente"}, &bytes.Buffer{}); err == nil {
		t.Fatal("flag desconhecida deveria falhar")
	}
}

func TestTaxMaxToleraAteOLimite(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "a.md", "# Título\n\nUma est"+"eira aqui.\n")
	if out, err := taxStep(config{Root: "../..", Tax: []string{dir}, TaxMax: 1}).Run(); err != nil {
		t.Fatalf("1 violação com limite 1 deveria passar: %v %q", err, out)
	}
	if out, err := taxStep(config{Root: "../..", Tax: []string{dir}, TaxMax: 0}).Run(); err == nil {
		t.Fatalf("1 violação com limite 0 deveria falhar: %q", out)
	}
}

func TestParseArgsTaxMax(t *testing.T) {
	c, err := parseArgs([]string{"-tax", "x", "-tax-max", "38"}, io.Discard)
	if err != nil || c.TaxMax != 38 {
		t.Fatalf("TaxMax = %d, err = %v", c.TaxMax, err)
	}
	if _, err := parseArgs([]string{"-tax-max", "-1"}, io.Discard); err == nil {
		t.Fatal("limite negativo deveria falhar")
	}
}

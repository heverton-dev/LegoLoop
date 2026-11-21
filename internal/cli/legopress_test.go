package cli_test

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/cli"
	"legoloop/internal/legopress"
)

func TestLegoPressCLI_Help(t *testing.T) {
	var stdout, stderr bytes.Buffer
	deps := cli.Deps{Stdout: &stdout, Stderr: &stderr}

	code := cli.Run(context.Background(), []string{"legopress", "--help"}, deps)
	if code != 0 {
		t.Fatalf("esperava exit code 0 para --help, obteve: %d", code)
	}

	out := stdout.String()
	for _, sub := range []string{"criar", "editar", "atualizar", "adequar", "auditar", "doctor", "preview"} {
		if !strings.Contains(out, sub) {
			t.Errorf("esperava subcomando %q no help do legopress, obteve: %s", sub, out)
		}
	}
}

func TestLegoPressCLI_Doctor(t *testing.T) {
	var stdout, stderr bytes.Buffer
	fakeCompiler := legopress.NewCompilerWithRunner("pandoc", "typst", func(bin string, _ ...string) ([]byte, error) {
		return []byte(bin + " 3.1.9\n"), nil
	})
	deps := cli.Deps{
		Stdout:    &stdout,
		Stderr:    &stderr,
		LegoPress: &cli.LegoPressDeps{Compiler: fakeCompiler},
	}

	code := cli.Run(context.Background(), []string{"legopress", "doctor"}, deps)
	if code != 0 {
		t.Fatalf("esperava exit code 0 para doctor, obteve: %d", code)
	}

	out := stdout.String()
	if !strings.Contains(out, "pandoc") {
		t.Errorf("esperava mencao ao pandoc no doctor do legopress, obteve: %s", out)
	}
}

func TestLegoPressCLI_Auditar(t *testing.T) {
	pdfs, err := filepath.Glob(filepath.Join("..", "..", "docs", "books", "legoloop-video-rea", "*.pdf"))
	if err != nil || len(pdfs) != 1 {
		t.Fatalf("esperava exatamente um PDF homologado, obteve %v (erro: %v)", pdfs, err)
	}
	pdfPath := pdfs[0]

	var stdout, stderr bytes.Buffer
	fakeAuditor := legopress.NewAuditorWithRunner(func(string) ([]byte, error) {
		return []byte(`{"total_pages":14,"orphan_pages":[],"broken_headers":[],"valid":true}`), nil
	})
	deps := cli.Deps{
		Stdout:    &stdout,
		Stderr:    &stderr,
		LegoPress: &cli.LegoPressDeps{Auditor: fakeAuditor},
	}

	code := cli.Run(context.Background(), []string{"legopress", "auditar", pdfPath}, deps)
	if code != 0 {
		t.Fatalf("esperava exit code 0 para auditar PDF homologado, obteve: %d (stderr: %s)", code, stderr.String())
	}

	out := stdout.String()
	if !strings.Contains(out, "14 páginas") && !strings.Contains(out, "14 paginas") {
		t.Errorf("esperava mencao a 14 paginas na auditoria, obteve: %s", out)
	}
}

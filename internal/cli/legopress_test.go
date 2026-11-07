package cli_test

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/cli"
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
	deps := cli.Deps{Stdout: &stdout, Stderr: &stderr}

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
	pdfPath := filepath.Join("..", "..", "docs", "books", "legoloop", "06-11-2026_LEGOLOOP-A-REVOLU-O-DA-IA-MODULAR-EM-PE-AS-SIMPLES.pdf")

	var stdout, stderr bytes.Buffer
	deps := cli.Deps{Stdout: &stdout, Stderr: &stderr}

	code := cli.Run(context.Background(), []string{"legopress", "auditar", pdfPath}, deps)
	if code != 0 {
		t.Fatalf("esperava exit code 0 para auditar PDF homologado, obteve: %d (stderr: %s)", code, stderr.String())
	}

	out := stdout.String()
	if !strings.Contains(out, "57 páginas") && !strings.Contains(out, "57 paginas") {
		t.Errorf("esperava mencao a 57 paginas na auditoria, obteve: %s", out)
	}
}

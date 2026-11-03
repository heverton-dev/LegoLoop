package cli_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/cli"
)

func TestCLI_Upgrade_RecommendsTheTemplateOfThePieceType(t *testing.T) {
	cases := map[string]string{
		"SKILL.md":          "TEMPLATE_PIECE_SKILL.md",
		"build.ps1":         "TEMPLATE_PIECE_SCRIPT.ps1",
		"fluxo.loop.json":   "TEMPLATE_PIECE_LOOP.json",
		"meu-mcp.json":      "TEMPLATE_PIECE_MCP.json",
		"guard-hook.json":   "TEMPLATE_PIECE_HOOK.json",
		"gate-receita.json": "TEMPLATE_PIECE_GATE.json",
	}
	dir := t.TempDir()
	for name, want := range cases {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte("conteúdo anêmico"), 0o644); err != nil {
			t.Fatal(err)
		}
		var stdout, stderr bytes.Buffer
		code := cli.Run(context.Background(), []string{"upgrade", p, "--diagnose"}, cli.Deps{Stdout: &stdout, Stderr: &stderr})
		if code != 0 {
			t.Fatalf("%s: exit %d: %s", name, code, stderr.String())
		}
		if !strings.Contains(stdout.String(), "docs/templates/"+want) {
			t.Errorf("%s: a recomendação deve apontar docs/templates/%s\n%s", name, want, stdout.String())
		}
	}
}

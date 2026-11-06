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

func TestCLIStept(t *testing.T) {
	t.Parallel()

	tmpHome := t.TempDir()

	// Simula presença de Gemini e Claude
	geminiDir := filepath.Join(tmpHome, ".gemini", "config", "rules")
	claudeDir := filepath.Join(tmpHome, ".claude", "rules")
	if err := os.MkdirAll(geminiDir, 0755); err != nil {
		t.Fatalf("falha ao criar gemini dir: %v", err)
	}
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatalf("falha ao criar claude dir: %v", err)
	}

	execCmd := func(args ...string) (int, string, string) {
		var outBuf, errBuf bytes.Buffer
		deps := cli.Deps{
			Stdout: &outBuf,
			Stderr: &errBuf,
			Home:   tmpHome,
		}
		cmdArgs := append([]string{"stept"}, args...)
		code := cli.Run(context.Background(), cmdArgs, deps)
		return code, outBuf.String(), errBuf.String()
	}

	t.Run("stept status exibe tabela formatada", func(t *testing.T) {
		code, out, errOut := execCmd("status")
		if code != 0 {
			t.Fatalf("esperava exit code 0, obteve %d. Err: %s", code, errOut)
		}
		cabecalhos := []string{"HARNESS", "INSTALADO", "STATUS", "MODO", "PAUSA", "ORIGEM"}
		for _, h := range cabecalhos {
			if !strings.Contains(out, h) {
				t.Errorf("tabela nao contem cabecalho %q:\n%s", h, out)
			}
		}
		if !strings.Contains(out, "Gemini") || !strings.Contains(out, "Claude") {
			t.Errorf("esperava presenca de Gemini e Claude na tabela:\n%s", out)
		}
	})

	t.Run("stept set atualiza override de harness", func(t *testing.T) {
		code, out, errOut := execCmd("set", "--harness", "claude", "--modo", "caveman", "--pausa", "sempre")
		if code != 0 {
			t.Fatalf("set falhou com code %d: %s", code, errOut)
		}
		if !strings.Contains(out, "caveman") || !strings.Contains(out, "override") {
			t.Errorf("saida do set nao refletiu override caveman:\n%s", out)
		}
	})

	t.Run("stept scan lista harnesses no sistema", func(t *testing.T) {
		code, out, errOut := execCmd("scan")
		if code != 0 {
			t.Fatalf("scan falhou com code %d: %s", code, errOut)
		}
		if !strings.Contains(out, "INSTALADO") || !strings.Contains(out, "Claude") {
			t.Errorf("saida do scan incompleta:\n%s", out)
		}
	})

	t.Run("stept off desativa governanca", func(t *testing.T) {
		code, out, errOut := execCmd("off", "--harness", "claude")
		if code != 0 {
			t.Fatalf("off falhou com code %d: %s", code, errOut)
		}
		if !strings.Contains(out, "off") {
			t.Errorf("saida do off nao refletiu modo off:\n%s", out)
		}
		// Confirma que stept-rules.md do claude foi removido
		claudeRule := filepath.Join(claudeDir, "stept-rules.md")
		if _, err := os.Stat(claudeRule); !os.IsNotExist(err) {
			t.Errorf("stept-rules.md ainda existe apos comando off: %s", claudeRule)
		}
	})

	t.Run("stept manual exibe guia formatado com moldura", func(t *testing.T) {
		code, out, errOut := execCmd("manual")
		if code != 0 {
			t.Fatalf("manual falhou com code %d: %s", code, errOut)
		}
		if !strings.Contains(out, "MANUAL DE USO DO ALTERNADOR UNIVERSAL STE-PT") {
			t.Errorf("manual nao contem titulo principal:\n%s", out)
		}
		if !strings.Contains(out, "┌") || !strings.Contains(out, "└") {
			t.Errorf("manual deve conter moldura formatada estilo gw:\n%s", out)
		}
	})
}

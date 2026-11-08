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

func TestCanalCLI_RedcAndGerar(t *testing.T) {
	tempDir := t.TempDir()

	// Criar fixture de métricas para o canal
	canalFixtureDir := filepath.Join(tempDir, "canal_raw")
	if err := os.MkdirAll(canalFixtureDir, 0755); err != nil {
		t.Fatalf("erro ao criar pasta fixture: %v", err)
	}

	metricsJSON := `[
		{
			"id": "v1",
			"titulo": "Metodo Viral 1",
			"views": 200000,
			"likes": 10000,
			"comentarios": 1000,
			"duracao_seg": 50,
			"gancho_3s": "Ninguem te contou isso antes",
			"metodo": "Metodo Dominador",
			"cluster_tema": "estrategia"
		}
	]`
	if err := os.WriteFile(filepath.Join(canalFixtureDir, "metricas.json"), []byte(metricsJSON), 0644); err != nil {
		t.Fatalf("erro ao gravar metricas.json: %v", err)
	}

	ctx := context.Background()
	var stdout, stderr bytes.Buffer
	deps := cli.Deps{
		Stdout: &stdout,
		Stderr: &stderr,
	}

	// 1. Testar legoloop canal redc
	code := cli.Run(ctx, []string{
		"canal", "redc", canalFixtureDir,
		"--canal", "canal-teste",
		"--root", tempDir,
		"--json",
	}, deps)

	if code != 0 {
		t.Fatalf("esperava exit code 0 para canal redc, obteve %d. Stderr: %s", code, stderr.String())
	}

	outStr := stdout.String()
	if !strings.Contains(outStr, "canal-teste") {
		t.Errorf("stdout deve conter canal-teste: %s", outStr)
	}

	// 2. Testar legoloop canal gerar
	stdout.Reset()
	stderr.Reset()
	codeGerar := cli.Run(ctx, []string{
		"canal", "gerar", "canal-teste",
		"--tema", "Agentes Autonomos",
		"--nome-marca", "MinhaMarca",
		"--root", tempDir,
		"--json",
	}, deps)

	if codeGerar != 0 {
		t.Fatalf("esperava exit code 0 para canal gerar, obteve %d. Stderr: %s", codeGerar, stderr.String())
	}

	gerarOut := stdout.String()
	if !strings.Contains(gerarOut, "roteiro") {
		t.Errorf("stdout deve conter indicacao de roteiro gerado: %s", gerarOut)
	}
}

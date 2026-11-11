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

func TestCLI_Upgrade_DiagnoseAndApply(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Skill anêmica (falta rigor e contratos)
	anemicSkill := filepath.Join(tempDir, "SKILL.md")
	os.WriteFile(anemicSkill, []byte("# Anemic Skill\nFaça um bom trabalho."), 0644)

	var stdout, stderr bytes.Buffer
	deps := cli.Deps{Stdout: &stdout, Stderr: &stderr}

	// Rodar diagnóstico
	code := cli.Run(context.Background(), []string{"upgrade", anemicSkill, "--diagnose"}, deps)
	if code != 0 {
		t.Fatalf("esperava exit code 0 no diagnose, obteve %d: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "GAPS PARA ATINGIR O PADRÃO OURO") {
		t.Errorf("saída não contém lista de gaps: %s", stdout.String())
	}

	// Tentar aplicar elevação sem corrigir (deve dar exit 1)
	stdout.Reset()
	stderr.Reset()
	code = cli.Run(context.Background(), []string{"upgrade", anemicSkill, "--apply"}, deps)
	if code != 1 {
		t.Errorf("esperava exit code 1 no apply para skill anêmica, obteve %d", code)
	}

	// 2. Agora atualizar para Padrão Ouro
	goldSkill := `# Gold Skill
by Addy Osmani

## Negative Guardrails
- NEVER commit stubs.
- Zero Stubs: no mocks.

## Contratos
- Consumes: A.md
- Produces: B.md

## Schemas
` + "```go\ntype Item struct { ID string }\n```" + `

## Testes e Verificação
- go test -v ./...
- exit code == 0

## Stopping Checklist
- [ ] Human reviews and approves.
`
	os.WriteFile(anemicSkill, []byte(goldSkill), 0644)

	stdout.Reset()
	stderr.Reset()
	code = cli.Run(context.Background(), []string{"upgrade", anemicSkill, "--apply"}, deps)
	if code != 0 {
		t.Fatalf("esperava exit code 0 no apply para skill ouro, obteve %d: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "SUCESSO: A peça SKILL.md foi elevada ao PADRÃO OURO") {
		t.Errorf("saída inesperada: %s", stdout.String())
	}
}

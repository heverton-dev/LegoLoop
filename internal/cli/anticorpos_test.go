package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCmdAnticorposSync(t *testing.T) {
	dir := t.TempDir()
	procDir := filepath.Join(dir, "processos", "teste-proc")
	if err := os.MkdirAll(procDir, 0o755); err != nil {
		t.Fatalf("criando pasta: %v", err)
	}
	jsonDummy := `{
  "schema_version": "anticorpos.v1",
  "processo": "teste-proc",
  "globais": [],
  "dominio": [
    {
      "assinatura": "sha256:123",
      "erro": "erro teste",
      "causa_raiz": "causa teste",
      "antidoto": "antidoto teste",
      "recorrencias": 1,
      "pacote": "pkg/teste",
      "timestamp": "2026-11-19T10:00:00Z",
      "antidoto_aplicado": true
    }
  ]
}`
	if err := os.WriteFile(filepath.Join(procDir, "anticorpos.json"), []byte(jsonDummy), 0o644); err != nil {
		t.Fatalf("gravando dummy: %v", err)
	}

	var stdout, stderr bytes.Buffer
	d := Deps{Stdout: &stdout, Stderr: &stderr}

	// Executa sync de todos apontando a raiz temporária
	args := []string{"sync", "--raiz", dir}
	code := cmdAnticorpos(context.Background(), args, d)
	if code != 0 {
		t.Fatalf("esperava exit code 0, obteve %d. Stderr: %s", code, stderr.String())
	}

	caminhoMD := filepath.Join(procDir, "ANTICORPOS.md")
	raw, err := os.ReadFile(caminhoMD)
	if err != nil {
		t.Fatalf("arquivo ANTICORPOS.md não foi gerado: %v", err)
	}
	if !strings.Contains(string(raw), "schema: anticorpos.v1") {
		t.Errorf("conteúdo inválido em ANTICORPOS.md:\n%s", string(raw))
	}
}

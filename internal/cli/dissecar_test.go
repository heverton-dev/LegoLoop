package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/dissecar"
)

func TestDissecarCommand_Run(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("falha ao obter cwd: %v", err)
	}

	subDir := filepath.Join(cwd, "test_sample")
	_ = os.MkdirAll(subDir, 0755)
	defer os.RemoveAll(subDir)

	sampleFile := filepath.Join(subDir, "main.go")
	err = os.WriteFile(sampleFile, []byte("package main\nfunc main() {}"), 0644)
	if err != nil {
		t.Fatalf("falha ao criar arquivo de teste: %v", err)
	}

	var buf bytes.Buffer
	cmd := NewDissecarCommand()
	cmd.output = &buf

	err = cmd.Run([]string{"--target", "test_sample", "--json"})
	if err != nil {
		t.Fatalf("erro inesperado no comando dissecar: %v", err)
	}

	var res dissecar.Result
	if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
		t.Fatalf("saida nao e um JSON valido: %v\nSaida: %s", err, buf.String())
	}

	if res.TargetID == "" {
		t.Errorf("esperado target_id preenchido")
	}
}

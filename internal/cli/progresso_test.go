package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIProgressoCicloCompleto(t *testing.T) {
	dir := t.TempDir()
	ticketsArq := filepath.Join(dir, "tickets.json")
	ticketsJSON := `[
  {"id":"TICK-01.1","title":"Passo 1","epic_id":"E01","initiative_id":"I04","validation_cmd":"echo ok"}
]`
	if err := os.WriteFile(ticketsArq, []byte(ticketsJSON), 0o644); err != nil {
		t.Fatal(err)
	}
	progArq := filepath.Join(dir, "PROGRESSO.json")
	var out, errBuf bytes.Buffer
	deps := Deps{Stdout: &out, Stderr: &errBuf}

	if code := Run(context.Background(), []string{"progresso", "gerar", "--tickets", ticketsArq, "--out", progArq, "--plano", "PLAN.md"}, deps); code != 0 {
		t.Fatalf("gerar exit=%d err=%s", code, errBuf.String())
	}
	if code := Run(context.Background(), []string{"progresso", "iniciar", "--arquivo", progArq, "--ticket", "TICK-01.1", "--responsavel", "worker-1", "--nota", "red"}, deps); code != 0 {
		t.Fatalf("iniciar exit=%d err=%s", code, errBuf.String())
	}
	if code := Run(context.Background(), []string{"progresso", "validar", "--arquivo", progArq, "--ticket", "TICK-01.1", "--raiz", dir}, deps); code != 0 {
		t.Fatalf("validar exit=%d err=%s", code, errBuf.String())
	}
	out.Reset()
	if code := Run(context.Background(), []string{"progresso", "status", "--arquivo", progArq}, deps); code != 0 {
		t.Fatalf("status exit=%d err=%s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "1/1 PASSED") {
		t.Fatalf("status inesperado: %s", out.String())
	}
}

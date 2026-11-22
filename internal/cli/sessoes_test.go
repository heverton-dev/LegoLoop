package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/triagem"
)

func TestCmdSessoesAtivasListaDuas(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "p")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"s1", "s2"} {
		mem := triagem.MemoriaExecucao{SessaoID: id, PID: os.Getpid(), EtapaAtual: "N-" + id}
		if err := triagem.SalvarMemoria(dir, mem); err != nil {
			t.Fatal(err)
		}
	}
	var out, errBuf bytes.Buffer
	code := cmdSessoes(context.Background(), []string{"p", "--json", "--processos-dir", base}, Deps{Stdout: &out, Stderr: &errBuf})
	if code != 0 {
		t.Fatalf("exit code = %d, stderr=%s", code, errBuf.String())
	}
	var lista []map[string]any
	if err := json.Unmarshal(out.Bytes(), &lista); err != nil {
		t.Fatalf("JSON inválido: %v\n%s", err, out.String())
	}
	if len(lista) != 2 {
		t.Fatalf("sessões ativas = %d, esperado 2", len(lista))
	}
	for _, s := range lista {
		if _, ok := s["sessao_id"]; !ok {
			t.Errorf("sessao_id ausente: %v", s)
		}
		if _, ok := s["pid"]; !ok {
			t.Errorf("pid ausente: %v", s)
		}
		if _, ok := s["etapa_atual"]; !ok {
			t.Errorf("etapa_atual ausente: %v", s)
		}
	}
	if strings.Contains(out.String(), base) {
		t.Errorf("saída expõe caminho absoluto: %s", out.String())
	}
}

func TestCmdSessoesVazio(t *testing.T) {
	base := t.TempDir()
	if err := os.MkdirAll(filepath.Join(base, "p"), 0o755); err != nil {
		t.Fatal(err)
	}
	var out, errBuf bytes.Buffer
	code := cmdSessoes(context.Background(), []string{"p", "--json", "--processos-dir", base}, Deps{Stdout: &out, Stderr: &errBuf})
	if code != 0 {
		t.Fatalf("exit code = %d, stderr=%s", code, errBuf.String())
	}
	if strings.TrimSpace(out.String()) != "[]" {
		t.Errorf("saída = %q, esperado []", strings.TrimSpace(out.String()))
	}
}

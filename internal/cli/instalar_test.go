package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/instalar"
)

func TestCmdInstalarDryRunEApply(t *testing.T) {
	home := t.TempDir()
	var out, errBuf bytes.Buffer
	d := Deps{Stdout: &out, Stderr: &errBuf, Home: home}

	code := Run(context.Background(), []string{"instalar", "--dry-run", "--home", home, "--json"}, d)
	if code != 0 {
		t.Fatalf("instalar --dry-run falhou (code=%d): %s", code, errBuf.String())
	}
	var resDry instalar.Resultado
	if err := json.Unmarshal(out.Bytes(), &resDry); err != nil {
		t.Fatalf("JSON inválido no dry-run: %v\nsaída: %s", err, out.String())
	}
	if resDry.Aplicado {
		t.Fatal("esperava Aplicado=false no dry-run")
	}

	out.Reset()
	errBuf.Reset()
	code = Run(context.Background(), []string{"instalar", "--apply", "--home", home, "--perfil", "padrao", "--json"}, d)
	if code != 0 {
		t.Fatalf("instalar --apply falhou (code=%d): %s", code, errBuf.String())
	}
	var resApply instalar.Resultado
	if err := json.Unmarshal(out.Bytes(), &resApply); err != nil {
		t.Fatalf("JSON inválido no apply: %v", err)
	}
	if !resApply.Aplicado || resApply.CaminhoJSON == "" {
		t.Fatalf("resultado inesperado no apply: %+v", resApply)
	}
	if _, err := os.Stat(filepath.Clean(resApply.CaminhoJSON)); err != nil {
		t.Fatalf("RELATORIO.json não encontrado: %v", err)
	}
}

func TestCmdInstalarPerfilInvalido(t *testing.T) {
	var out, errBuf bytes.Buffer
	d := Deps{Stdout: &out, Stderr: &errBuf, Home: t.TempDir()}
	code := Run(context.Background(), []string{"instalar", "--perfil", "xyz"}, d)
	if code != 2 {
		t.Fatalf("esperava exit 2 para perfil inválido, obteve %d (%s)", code, errBuf.String())
	}
	if !strings.Contains(errBuf.String(), "perfil") {
		t.Fatalf("mensagem esperada com 'perfil', obteve: %s", errBuf.String())
	}
}

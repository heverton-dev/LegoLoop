package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func registrarAnticorpo(t *testing.T, raiz string, args ...string) (int, string) {
	t.Helper()
	var out, errb bytes.Buffer
	code := cmdAnticorpos(context.Background(), append([]string{"registrar"}, args...), Deps{Stdout: &out, Stderr: &errb})
	return code, out.String() + errb.String()
}

func TestAnticorposRegistrarGravaERecorre(t *testing.T) {
	raiz := t.TempDir()
	proc := filepath.Join(raiz, "processos", "demo")
	if err := os.MkdirAll(proc, 0o755); err != nil {
		t.Fatal(err)
	}
	args := []string{"demo", "--raiz", raiz, "--pacote", "internal/gate", "--erro", "titulo de criterios nao reconhecido", "--causa", "regex estreita", "--antidoto", "aceitar o titulo canonico"}
	if code, saida := registrarAnticorpo(t, raiz, args...); code != 0 {
		t.Fatalf("registrar saiu %d: %s", code, saida)
	}
	js, _ := os.ReadFile(filepath.Join(proc, "anticorpos.json"))
	md, _ := os.ReadFile(filepath.Join(proc, "ANTICORPOS.md"))
	if !strings.Contains(string(js), "titulo de criterios nao reconhecido") || !strings.Contains(string(md), "regex estreita") {
		t.Fatalf("registro ausente.\njson=%s\nmd=%s", js, md)
	}
	if code, _ := registrarAnticorpo(t, raiz, args...); code != 0 || !strings.Contains(readFile(t, filepath.Join(proc, "anticorpos.json")), `"recorrencias": 2`) {
		t.Fatal("a repetição deve somar uma recorrência")
	}
}

func TestAnticorposRegistrarExigePacoteEErro(t *testing.T) {
	raiz := t.TempDir()
	if code, _ := registrarAnticorpo(t, raiz, "demo", "--raiz", raiz, "--pacote", "x"); code != 2 {
		t.Fatalf("sem --erro deveria sair com 2, saiu %d", code)
	}
}

func readFile(t *testing.T, p string) string {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

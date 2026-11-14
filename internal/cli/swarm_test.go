package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestCmdSwarm_Errors(t *testing.T) {
	var out, errBuf bytes.Buffer
	d := Deps{Stdout: &out, Stderr: &errBuf}

	if code := cmdSwarm(context.Background(), nil, d); code != 2 {
		t.Fatalf("esperava exit 2 sem args, obteve %d", code)
	}
	if !strings.Contains(errBuf.String(), "uso: legoloop swarm") {
		t.Fatalf("mensagem de uso ausente: %q", errBuf.String())
	}

	errBuf.Reset()
	if code := cmdSwarm(context.Background(), []string{"invalido"}, d); code != 2 {
		t.Fatalf("esperava exit 2 para subcomando desconhecido, obteve %d", code)
	}
	if !strings.Contains(errBuf.String(), "subcomando swarm desconhecido") {
		t.Fatalf("mensagem de erro ausente: %q", errBuf.String())
	}
}

func TestCmdSwarm_AddPostReadStatus(t *testing.T) {
	dir := t.TempDir()
	var out, errBuf bytes.Buffer
	d := Deps{Stdout: &out, Stderr: &errBuf}
	ctx := context.Background()

	if code := cmdSwarm(ctx, []string{"add", "-dir", dir, "-name", "w1", "-role", "engenheiro-tdd"}, d); code != 0 {
		t.Fatalf("add falhou (%d): %s", code, errBuf.String())
	}

	out.Reset()
	if code := cmdSwarm(ctx, []string{"post", "-dir", dir, "-from", "w1", "-body", "ola"}, d); code != 0 {
		t.Fatalf("post falhou (%d): %s", code, errBuf.String())
	}

	out.Reset()
	if code := cmdSwarm(ctx, []string{"read", "-dir", dir, "-since", "0"}, d); code != 0 {
		t.Fatalf("read falhou (%d): %s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "ola") {
		t.Fatalf("mensagem esperada nao encontrada em read: %q", out.String())
	}

	out.Reset()
	if code := cmdSwarm(ctx, []string{"status", "-dir", dir}, d); code != 0 {
		t.Fatalf("status falhou (%d): %s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "w1") {
		t.Fatalf("worker w1 ausente em status: %q", out.String())
	}
}

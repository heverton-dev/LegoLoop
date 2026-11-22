package a2ui

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"
)

func TestAwaitActionUIVence(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pr, pw := io.Pipe()
	defer pr.Close()
	defer pw.Close()

	uiCh := make(chan Action, 1)
	esperada := Action{
		SessionID: "sess_1",
		SurfaceID: "proc:concepcao",
		Name:      "stop_gate_decision",
		Payload:   map[string]any{"decision": "1"},
	}
	uiCh <- esperada

	got, err := AwaitAction(ctx, "sess_1", "proc:concepcao", uiCh, pr)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if got.Origin != "ui" {
		t.Fatalf("Origin = %q, esperado %q", got.Origin, "ui")
	}
	if got.Name != esperada.Name || got.Payload["decision"] != "1" {
		t.Fatalf("ação divergente: %+v", got)
	}
}

func TestAwaitActionStdinVence(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	uiCh := make(chan Action)
	stdin := strings.NewReader("  aprovar etapa 1  \r\n")

	got, err := AwaitAction(ctx, "sess_2", "proc:fabricar", uiCh, stdin)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if got.Origin != "stdin" {
		t.Fatalf("Origin = %q, esperado %q", got.Origin, "stdin")
	}
	if got.Name != "stdin_reply" {
		t.Fatalf("Name = %q, esperado %q", got.Name, "stdin_reply")
	}
	if got.SessionID != "sess_2" || got.SurfaceID != "proc:fabricar" {
		t.Fatalf("identificadores incorretos: %+v", got)
	}
	if got.Payload["text"] != "aprovar etapa 1" {
		t.Fatalf("Payload[text] = %v, esperado %q", got.Payload["text"], "aprovar etapa 1")
	}
}

func TestAwaitActionContextoCancelado(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	pr, pw := io.Pipe()
	defer pr.Close()
	defer pw.Close()

	uiCh := make(chan Action)
	done := make(chan error, 1)
	go func() {
		_, err := AwaitAction(ctx, "sess_3", "proc:raiz", uiCh, pr)
		done <- err
	}()

	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("erro = %v, esperado context.Canceled", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("AwaitAction travou após cancelamento do contexto")
	}
}

func TestAwaitActionSemStdinEComUIChFechado(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	uiCh := make(chan Action, 1)
	uiCh <- Action{Name: "submit_answer"}
	got, err := AwaitAction(ctx, "sess_4", "proc:sub", uiCh, nil)
	if err != nil {
		t.Fatalf("erro inesperado com stdin nil: %v", err)
	}
	if got.Origin != "ui" || got.SessionID != "sess_4" || got.SurfaceID != "proc:sub" {
		t.Fatalf("ação normalizada incorreta: %+v", got)
	}

	fechado := make(chan Action)
	close(fechado)
	got2, err := AwaitAction(ctx, "sess_5", "proc:sub", fechado, strings.NewReader("texto livre\n"))
	if err != nil {
		t.Fatalf("erro inesperado com uiCh fechado: %v", err)
	}
	if got2.Origin != "stdin" || got2.Payload["text"] != "texto livre" {
		t.Fatalf("ação stdin incorreta: %+v", got2)
	}
}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/triagem"
)

func TestRunTriagemDeterministicoJSON(t *testing.T) {
	dir := t.TempDir()
	loopDir := filepath.Join(dir, "meu-processo")
	if err := os.MkdirAll(loopDir, 0o755); err != nil {
		t.Fatal(err)
	}
	spec := triagem.GerarTriagemPadrao("meu-processo", []string{"etapa-1", "etapa-2"})
	dados, err := json.Marshal(spec)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(loopDir, "triagem.json"), dados, 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := runTriagem(context.Background(), []string{"--loops-dir", dir, "--json", "/legoloop-meu-processo"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("esperado exit 0, obtido %d: %s", code, stderr.String())
	}

	var dec triagem.Decisao
	if err := json.Unmarshal(stdout.Bytes(), &dec); err != nil {
		t.Fatalf("JSON inválido: %v", err)
	}
	if dec.Processo != "meu-processo" || dec.Camada != "deterministica" || dec.EtapaInicial != "etapa-1" {
		t.Fatalf("decisão inesperada: %+v", dec)
	}
}

func TestRunTriagemAmbiguidadeSemLLM(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer
	code := runTriagem(context.Background(), []string{"--loops-dir", dir, "--json", "pedido ambiguo qualquer"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("esperado exit 0 na ambiguidade estruturada, obtido %d: %s", code, stderr.String())
	}

	var dec triagem.Decisao
	if err := json.Unmarshal(stdout.Bytes(), &dec); err != nil {
		t.Fatalf("JSON inválido: %v", err)
	}
	if !dec.PrecisaConfirmacao || dec.Camada != "ambiguidade" {
		t.Fatalf("esperada ambiguidade estruturada, obtido %+v", dec)
	}
}

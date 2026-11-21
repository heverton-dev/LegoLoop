package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/caminhos"
	"legoloop/internal/triagem"
)

func TestRunTriagemDeterministicoJSON(t *testing.T) {
	dir := t.TempDir()
	procDir := filepath.Join(dir, "meu-processo")
	if err := os.MkdirAll(procDir, 0o755); err != nil {
		t.Fatal(err)
	}
	spec := triagem.GerarTriagemPadrao("meu-processo", []string{"etapa-1", "etapa-2"})
	dados, err := json.Marshal(spec)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(procDir, "triagem.json"), dados, 0o644); err != nil {
		t.Fatal(err)
	}

	for _, flagDir := range []string{"--processos-dir", "--" + caminhos.DirAntigo + "-dir"} {
		var stdout, stderr bytes.Buffer
		code := runTriagem(context.Background(), []string{flagDir, dir, "--json", "/legoloop-meu-processo"}, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("%s: esperado exit 0, obtido %d: %s", flagDir, code, stderr.String())
		}
		var dec triagem.Decisao
		if err := json.Unmarshal(stdout.Bytes(), &dec); err != nil || dec.Processo != "meu-processo" || dec.EtapaInicial != "etapa-1" {
			t.Fatalf("%s: decisão inesperada (%v): %+v", flagDir, err, dec)
		}
	}
}

func TestRunTriagemAmbiguidadeSemLLM(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer
	code := runTriagem(context.Background(), []string{"--processos-dir", dir, "--json", "pedido ambiguo qualquer"}, &stdout, &stderr)
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

func TestRunMigrarProcessosCompatibilidade(t *testing.T) {
	dir := t.TempDir()
	for _, sub := range []string{"migrar-processos", "migrar-" + caminhos.DirAntigo} {
		var stdout, stderr bytes.Buffer
		code := runMigrarProcessos(sub, []string{"--dir", dir, "--dry-run"}, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("%s: esperado exit 0, obtido %d: %s", sub, code, stderr.String())
		}
	}
}

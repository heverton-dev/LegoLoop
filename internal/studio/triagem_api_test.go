package studio

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/triagem"
)

func TestHandleTriagemDeterministico(t *testing.T) {
	projDir := t.TempDir()
	loopDir := filepath.Join(projDir, "loops", "demo-proc")
	if err := os.MkdirAll(loopDir, 0o755); err != nil {
		t.Fatal(err)
	}
	spec := triagem.GerarTriagemPadrao("demo-proc", []string{"inicio"})
	dados, err := json.Marshal(spec)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(loopDir, "triagem.json"), dados, 0o644); err != nil {
		t.Fatal(err)
	}

	srv := &Server{ProjectDir: projDir}
	body := bytes.NewBufferString(`{"entrada":"legoloop demo-proc","workspace":""}`)
	req := httptest.NewRequest(http.MethodPost, "/api/triagem", body)
	rec := httptest.NewRecorder()

	srv.handleTriagem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("esperado 200, obtido %d: %s", rec.Code, rec.Body.String())
	}

	var dec triagem.Decisao
	if err := json.Unmarshal(rec.Body.Bytes(), &dec); err != nil {
		t.Fatalf("JSON inválido: %v", err)
	}
	if dec.Processo != "demo-proc" || dec.Camada != "deterministica" {
		t.Fatalf("decisão inesperada: %+v", dec)
	}
}

func TestHandleTriagemAmbiguidade(t *testing.T) {
	srv := &Server{ProjectDir: t.TempDir()}
	body := bytes.NewBufferString(`{"entrada":"algo ambiguo"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/triagem", body)
	rec := httptest.NewRecorder()

	srv.handleTriagem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("esperado 200, obtido %d: %s", rec.Code, rec.Body.String())
	}

	var dec triagem.Decisao
	if err := json.Unmarshal(rec.Body.Bytes(), &dec); err != nil {
		t.Fatalf("JSON inválido: %v", err)
	}
	if !dec.PrecisaConfirmacao || dec.Camada != "ambiguidade" {
		t.Fatalf("esperada ambiguidade, obtido %+v", dec)
	}
}

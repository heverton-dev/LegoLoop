package kraften

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestReceiptGravaJSONValido(t *testing.T) {
	dest := t.TempDir()
	r := DeployReceipt{AppID: "app01", Status: "done", Timestamp: "2026-11-26T12:00:00Z", CFStatus: 1, HTTPCode: 200}
	caminho, err := GravarReceipt(r, dest)
	if err != nil {
		t.Fatalf("GravarReceipt: %v", err)
	}
	if filepath.Base(caminho) != "DEPLOY_RECEIPT.json" {
		t.Fatalf("nome do recibo inesperado: %s", caminho)
	}
	corpo, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatalf("ler recibo: %v", err)
	}
	var volta DeployReceipt
	if err := json.Unmarshal(corpo, &volta); err != nil {
		t.Fatalf("recibo não é JSON válido: %v", err)
	}
	if volta != r {
		t.Fatalf("recibo divergente: %+v != %+v", volta, r)
	}
}

func TestReceiptSobrescreveAnterior(t *testing.T) {
	dest := t.TempDir()
	antigo := DeployReceipt{AppID: "app01", Status: "failed", Timestamp: "t0", CFStatus: 0, HTTPCode: 500}
	if _, err := GravarReceipt(antigo, dest); err != nil {
		t.Fatalf("primeiro recibo: %v", err)
	}
	novo := DeployReceipt{AppID: "app01", Status: "done", Timestamp: "t1", CFStatus: 1, HTTPCode: 200}
	if _, err := GravarReceipt(novo, dest); err != nil {
		t.Fatalf("segundo recibo: %v", err)
	}
	corpo, _ := os.ReadFile(filepath.Join(dest, "DEPLOY_RECEIPT.json"))
	var volta DeployReceipt
	json.Unmarshal(corpo, &volta)
	if volta.Status != "done" {
		t.Fatalf("recibo antigo não foi substituído: %+v", volta)
	}
}

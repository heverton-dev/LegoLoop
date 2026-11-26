package kraften

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestDeployConfigValidateOK(t *testing.T) {
	cfg := DeployConfig{AppID: "app01", ZipPath: "/data/apps/app01/pacote.zip", Subdomain: "sub01", TargetHost: "server.vpsconexao.org"}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("configuração válida rejeitada: %v", err)
	}
}

func TestDeployConfigValidateVazio(t *testing.T) {
	for _, cfg := range []DeployConfig{
		{AppID: "", ZipPath: "pacote.zip"},
		{AppID: "  ", ZipPath: "pacote.zip"},
		{AppID: "app01", ZipPath: ""},
	} {
		if err := cfg.Validate(); !errors.Is(err, ErrInvalidDeployInput) {
			t.Fatalf("esperava ErrInvalidDeployInput para %+v, obteve %v", cfg, err)
		}
	}
}

func TestDeployConfigValidateZipSlip(t *testing.T) {
	cfg := DeployConfig{AppID: "app01", ZipPath: "/data/apps/../../etc/pacote.zip"}
	if err := cfg.Validate(); !errors.Is(err, ErrInvalidDeployInput) {
		t.Fatalf("caminho com .. deve ser rejeitado, obteve %v", err)
	}
}

func TestDeployReceiptJSON(t *testing.T) {
	r := DeployReceipt{AppID: "app01", Status: "done", Timestamp: "2026-11-26T00:00:00Z", CFStatus: 1, HTTPCode: 200}
	raw, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var volta DeployReceipt
	if err := json.Unmarshal(raw, &volta); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if volta != r {
		t.Fatalf("round-trip alterou o recibo: %+v != %+v", volta, r)
	}
}

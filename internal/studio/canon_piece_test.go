package studio

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"legoloop/internal/canon"
)

func TestFilterPlanBySHA(t *testing.T) {
	full := canon.Plan{
		Vault: "/tmp/vault",
		Groups: []canon.GroupPlan{
			{SHA256: "hash1", SavedBytes: 100},
			{SHA256: "hash2", SavedBytes: 200},
		},
		SavedBytes: 300,
	}
	f1 := filterPlanBySHA(full, "hash1")
	if len(f1.Groups) != 1 || f1.Groups[0].SHA256 != "hash1" || f1.SavedBytes != 100 {
		t.Fatalf("filtro falhou: %+v", f1)
	}
	fNone := filterPlanBySHA(full, "hash3")
	if len(fNone.Groups) != 0 || fNone.SavedBytes != 0 {
		t.Fatalf("esperava vazio: %+v", fNone)
	}
}

func TestHandleCanonizePieceBadRequest(t *testing.T) {
	s := &Server{}
	req := httptest.NewRequest("POST", "/api/canonize/piece", bytes.NewBufferString(`{}`))
	rec := httptest.NewRecorder()
	s.handleCanonizePiece(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperava 400, obteve %d", rec.Code)
	}
}

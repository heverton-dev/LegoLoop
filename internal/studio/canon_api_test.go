package studio

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestCanonizeStatus(t *testing.T) {
	s, _, _ := newAssetsServer(t)
	rec := do(t, s.Handler(), "GET", "/api/canonize/status", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; quer %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var body struct {
		Vault  string   `json:"vault"`
		Roots  []string `json:"roots"`
		Drives []string `json:"drives"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Vault == "" {
		t.Errorf("vault veio vazio")
	}
	if len(body.Drives) == 0 {
		t.Errorf("drives veio vazio")
	}
}

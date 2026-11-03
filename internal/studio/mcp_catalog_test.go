package studio

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestMCPCatalogList(t *testing.T) {
	s, _ := newServer(t)
	h := s.Handler()

	rec := do(t, h, "GET", "/api/mcps/catalog", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body = %s", rec.Code, rec.Body.String())
	}

	var res struct {
		Catalog []MCPCatalogEntry `json:"catalog"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("falha ao parsear json: %v", err)
	}

	if len(res.Catalog) == 0 {
		t.Fatal("catálogo MCP embutido não deve ser vazio")
	}

	for _, entry := range res.Catalog {
		if entry.Name == "" || entry.Origin == "" || entry.Command == "" {
			t.Errorf("mcp com campos vazios: %+v", entry)
		}
	}
}

func TestMCPInstallRequiresConsent(t *testing.T) {
	s, _ := newServer(t)
	h := s.Handler()

	body := `{"name":"playwright","consent":false}`
	rec := do(t, h, "POST", "/api/mcps/install", body)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperava 400 sem consentimento, obteve %d (%s)", rec.Code, rec.Body.String())
	}
}

func TestMCPInstallWithConsentCreatesConfig(t *testing.T) {
	s, _ := newServer(t)
	h := s.Handler()

	body := `{"name":"playwright","consent":true}`
	rec := do(t, h, "POST", "/api/mcps/install", body)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperava 200 com consentimento, obteve %d (%s)", rec.Code, rec.Body.String())
	}
}

package studio

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestTemplatesEndpointReturnsAtLeastThreeValidTemplates(t *testing.T) {
	s, _ := newServer(t)
	h := s.Handler()

	rec := do(t, h, "GET", "/api/templates", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body = %s", rec.Code, rec.Body.String())
	}

	var res struct {
		Templates []Template `json:"templates"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("falha ao parsear json de templates: %v", err)
	}

	if len(res.Templates) < 3 {
		t.Fatalf("esperava pelo menos 3 templates, obteve %d", len(res.Templates))
	}

	expectedIDs := map[string]bool{
		"greenfield":          false,
		"reverse-engineering": false,
		"security-audit":      false,
		"mega-pipeline":       false,
		"impeccable-design":   false,
	}

	for _, tmpl := range res.Templates {
		if _, ok := expectedIDs[tmpl.ID]; ok {
			expectedIDs[tmpl.ID] = true
		}
		if tmpl.Name == "" || tmpl.Description == "" {
			t.Errorf("template %s com campos vazios", tmpl.ID)
		}
		if len(tmpl.Composition.Nodes) == 0 {
			t.Errorf("template %s sem nós", tmpl.ID)
		}
	}

	for id, found := range expectedIDs {
		if !found {
			t.Errorf("template esperado %q não encontrado", id)
		}
	}
}

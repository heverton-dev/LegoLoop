package plangen

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestInitiativeRecognitionAndRouteEmission(t *testing.T) {
	spec := `# 📋 SPEC: Teste com Iniciativa
Referência da Intenção: docs/intents/2026-11-07_INTENT_TESTE.md
Iniciativa: I01 — Plataforma Central

## 5. Critérios Binários de Aceite por Nó / Ticket (BDD)
### Iniciativa I01 — Plataforma Central
#### Épico 01 — Núcleo
* [ ] AC-E01.1: Critério A
`
	ps := pathSet{baseDir: t.TempDir(), specPath: "docs/specs/test.md", planPath: "docs/plans/test.md", routePath: "docs/routes/test.json"}
	meta := newMeta(ps, spec, "2026-11-07", SupervisionMicro)

	if meta.InitiativeRef != "I01 — Plataforma Central" {
		t.Fatalf("esperava InitiativeRef capturada, obteve: %q", meta.InitiativeRef)
	}
	if meta.Initiatives[1] != "Plataforma Central" {
		t.Fatalf("esperava Initiatives[1] = 'Plataforma Central', obteve: %q", meta.Initiatives[1])
	}

	tickets := []Ticket{
		{
			ID:            "TICK-01.1",
			Title:         "Ticket com Iniciativa",
			EpicID:        "E01",
			InitiativeID:  "I01",
			Nature:        "micro",
			ExecutionMode: "tdd_red_green",
			TargetFiles:   []string{"pkg/core.go"},
			SpecCriteria:  []string{"AC-E01.1"},
		},
	}

	nodes, err := buildNodes(tickets, meta)
	if err != nil {
		t.Fatalf("buildNodes falhou: %v", err)
	}

	doc := newRouteDoc(meta, nodes)
	data, err := doc.marshal()
	if err != nil {
		t.Fatalf("marshal da rota falhou: %v", err)
	}

	var parsed struct {
		InitiativeRef *string `json:"initiative_ref"`
		Nodes         []struct {
			ID           string  `json:"id"`
			InitiativeID *string `json:"initiative_id"`
		} `json:"nodes"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal falhou: %v", err)
	}

	if parsed.InitiativeRef == nil || !strings.Contains(*parsed.InitiativeRef, "I01") {
		t.Errorf("initiative_ref não propagada para o topo: %v", parsed.InitiativeRef)
	}
	if len(parsed.Nodes) == 0 || parsed.Nodes[0].InitiativeID == nil || *parsed.Nodes[0].InitiativeID != "I01" {
		t.Errorf("initiative_id não propagada para o nó: %+v", parsed.Nodes[0])
	}
}

func TestInitiativeFallbackWhenAbsent(t *testing.T) {
	spec := `# 📋 SPEC: Teste sem Iniciativa
Referência da Intenção: docs/intents/2026-11-07_INTENT_LEGADO.md

## 5. Critérios
#### Épico 01 — Legado
* [ ] AC-E01.1: Legado
`
	ps := pathSet{baseDir: t.TempDir(), specPath: "docs/specs/test.md", planPath: "docs/plans/test.md", routePath: "docs/routes/test.json"}
	meta := newMeta(ps, spec, "2026-11-07", SupervisionMicro)

	if meta.InitiativeRef != "" {
		t.Errorf("esperava InitiativeRef vazia, obteve %q", meta.InitiativeRef)
	}

	tickets := []Ticket{
		{
			ID:            "TICK-01.1",
			Title:         "Ticket Legado",
			EpicID:        "E01",
			Nature:        "micro",
			ExecutionMode: "sequential",
			TargetFiles:   []string{"pkg/legacy.go"},
			SpecCriteria:  []string{"AC-E01.1"},
		},
	}

	nodes, err := buildNodes(tickets, meta)
	if err != nil {
		t.Fatalf("buildNodes: %v", err)
	}

	doc := newRouteDoc(meta, nodes)
	if doc.InitiativeRef != nil {
		t.Errorf("esperava initiative_ref nil no plano legado, obteve: %v", *doc.InitiativeRef)
	}
	if nodes[0].InitiativeID != nil {
		t.Errorf("esperava initiative_id nil no ticket legado, obteve: %v", *nodes[0].InitiativeID)
	}
}

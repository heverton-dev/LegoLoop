package studio

import (
	"encoding/json"
	"legoloop/internal/caminhos"
	"net/http"
	"path/filepath"
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/studiodoc"
)

func TestProcessosAndDAGCanonicalRoutes(t *testing.T) {
	s, home, _ := newAssetsServer(t)
	loopsDir := filepath.Join(home, ".aidd", "canonical", caminhos.DirLoops)
	store, err := studiodoc.NewStore(loopsDir)
	if err != nil {
		t.Fatal(err)
	}
	s.Loops = store

	// 1. GET /api/dag status
	rec := do(t, s.Handler(), "GET", "/api/dag", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/dag esperado 200, obteve %d", rec.Code)
	}
	var dagStatus map[string]any
	decode(t, rec, &dagStatus)
	if dagStatus["status"] != "idle" {
		t.Errorf("status do DAG esperado idle, obteve %v", dagStatus["status"])
	}

	// 2. POST /api/processos para criar novo processo
	lp := studiodoc.Processo{
		ID:          "proc-canonico",
		Name:        "Processo Canônico",
		Description: "Um processo canônico ASD-STE100",
		Version:     "1.0.0",
		Trigger:     &dag.Trigger{Type: dag.TriggerManual},
		Nodes: []dag.Node{
			{
				ID:      "n1",
				Type:    dag.TypeScript,
				Command: []string{"go", "version"},
				Routing: dag.RoutingConfig{Mode: dag.ModeManual, Model: "local"},
				Output:  dag.HandoffContract{Artifact: "RES.md", Gate: "file_exists('RES.md')"},
			},
		},
	}
	payload, _ := json.Marshal(lp)
	rec = do(t, s.Handler(), "POST", "/api/processos", string(payload))
	if rec.Code != http.StatusOK {
		t.Fatalf("POST /api/processos esperado 200, obteve %d: %s", rec.Code, rec.Body.String())
	}

	// 3. GET /api/processos/{id}
	rec = do(t, s.Handler(), "GET", "/api/processos/proc-canonico", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/processos/proc-canonico esperado 200, obteve %d", rec.Code)
	}

	// 4. GET /api/dag/{id}
	rec = do(t, s.Handler(), "GET", "/api/dag/proc-canonico", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/dag/proc-canonico esperado 200, obteve %d", rec.Code)
	}
}

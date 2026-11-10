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

func TestLoopsAPICRUD(t *testing.T) {
	s, home, _ := newAssetsServer(t)
	loopsDir := filepath.Join(home, ".aidd", "canonical", caminhos.DirLoops)
	store, err := studiodoc.NewStore(loopsDir)
	if err != nil {
		t.Fatal(err)
	}
	s.Loops = store

	// 1. Listar inicialmente vazio
	rec := do(t, s.Handler(), "GET", "/api/loops", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/loops esperado 200, obteve %d", rec.Code)
	}
	var res map[string][]studiodoc.Summary
	decode(t, rec, &res)
	if got := res[caminhos.ChaveAPILoops]; len(got) != 0 {
		t.Fatalf("esperava 0 loops, obteve %d", len(got))
	}

	// 2. Salvar novo loop
	lp := studiodoc.Loop{
		ID:          "meu-loop",
		Name:        "Meu Loop",
		Description: "Um loop de teste",
		Version:     "2.0.0",
		Trigger:     &dag.Trigger{Type: dag.TriggerManual},
		Nodes: []dag.Node{
			{
				ID:      "n1",
				Type:    dag.TypeScript,
				Command: []string{"python", "run.py"},
				Routing: dag.RoutingConfig{Mode: dag.ModeManual, Model: "local"},
				Output:  dag.HandoffContract{Artifact: "RES.md", Gate: "file_exists('RES.md')"},
			},
		},
	}
	payload, _ := json.Marshal(lp)
	rec = do(t, s.Handler(), "POST", "/api/loops", string(payload))
	if rec.Code != http.StatusOK {
		t.Fatalf("POST /api/loops esperado 200, obteve %d: %s", rec.Code, rec.Body.String())
	}

	// 3. Obter loop salvo por ID
	rec = do(t, s.Handler(), "GET", "/api/loops/meu-loop", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/loops/meu-loop esperado 200, obteve %d", rec.Code)
	}
	var got studiodoc.Loop
	decode(t, rec, &got)
	if got.ID != "meu-loop" || got.Name != "Meu Loop" {
		t.Fatalf("loop divergente: %+v", got)
	}

	// 4. Duplicar loop
	rec = do(t, s.Handler(), "POST", "/api/loops/meu-loop/duplicate", `{"id":"meu-loop-copia","name":"Meu Loop Cópia"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("POST /api/loops/meu-loop/duplicate esperado 200, obteve %d: %s", rec.Code, rec.Body.String())
	}

	// 5. Listar deve ter 2 loops
	rec = do(t, s.Handler(), "GET", "/api/loops", "")
	res = nil
	decode(t, rec, &res)
	if got := res[caminhos.ChaveAPILoops]; len(got) != 2 {
		t.Fatalf("esperava 2 loops após duplicação, obteve %d", len(got))
	}

	// 6. Deletar cópia
	rec = do(t, s.Handler(), "DELETE", "/api/loops/meu-loop-copia", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("DELETE esperado 200, obteve %d", rec.Code)
	}

	// 7. Obter cópia deletada dá 404
	rec = do(t, s.Handler(), "GET", "/api/loops/meu-loop-copia", "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("GET de loop deletado esperado 404, obteve %d", rec.Code)
	}

	// 8. Fork Copy-on-Write via /api/loops/{id}/fork
	recFork := do(t, s.Handler(), "POST", "/api/loops/meu-loop/fork", `{"id":"meu-loop-fork","name":"Meu Loop Fork"}`)
	if recFork.Code != http.StatusOK {
		t.Fatalf("POST /api/loops/meu-loop/fork esperado 200, obteve %d: %s", recFork.Code, recFork.Body.String())
	}
	var forked studiodoc.Loop
	decode(t, recFork, &forked)
	if forked.ID != "meu-loop-fork" {
		t.Fatalf("ID forked incorreto: %s", forked.ID)
	}
}

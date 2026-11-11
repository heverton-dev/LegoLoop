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

func TestProcessosAPICRUD(t *testing.T) {
	s, home, _ := newAssetsServer(t)
	dir := filepath.Join(home, ".aidd", "canonical", caminhos.DirLoops)
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	s.Processos = store

	if got, _ := store.List(); len(got) != 0 {
		t.Fatalf("esperava 0 Processos, obteve %d", len(got))
	}

	doc := studiodoc.Processo{
		ID:          "meu-processo",
		Name:        "Meu Processo",
		Description: "Um Processo de teste",
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
	payload, _ := json.Marshal(doc)
	rec := do(t, s.Handler(), "POST", "/api/processos", string(payload))
	if rec.Code != http.StatusOK {
		t.Fatalf("POST /api/processos esperado 200, obteve %d: %s", rec.Code, rec.Body.String())
	}

	rec = do(t, s.Handler(), "GET", "/api/processos/meu-processo", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/processos/meu-processo esperado 200, obteve %d", rec.Code)
	}
	var got studiodoc.Processo
	decode(t, rec, &got)
	if got.ID != "meu-processo" || got.Name != "Meu Processo" {
		t.Fatalf("Processo divergente: %+v", got)
	}

	rec = do(t, s.Handler(), "POST", "/api/processos/meu-processo/duplicate", `{"id":"meu-processo-copia","name":"Meu Processo Cópia"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("POST duplicate esperado 200, obteve %d: %s", rec.Code, rec.Body.String())
	}

	if got, _ := store.List(); len(got) != 2 {
		t.Fatalf("esperava 2 Processos após duplicação, obteve %d", len(got))
	}

	rec = do(t, s.Handler(), "DELETE", "/api/processos/meu-processo-copia", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("DELETE esperado 200, obteve %d", rec.Code)
	}

	rec = do(t, s.Handler(), "GET", "/api/processos/meu-processo-copia", "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("GET de Processo removido esperado 404, obteve %d", rec.Code)
	}

	recFork := do(t, s.Handler(), "POST", "/api/processos/meu-processo/fork", `{"id":"meu-processo-fork","name":"Meu Processo Fork"}`)
	if recFork.Code != http.StatusOK {
		t.Fatalf("POST fork esperado 200, obteve %d: %s", recFork.Code, recFork.Body.String())
	}
	var forked studiodoc.Processo
	decode(t, recFork, &forked)
	if forked.ID != "meu-processo-fork" {
		t.Fatalf("ID forked incorreto: %s", forked.ID)
	}
}

func TestRotaAntigaResponveIgualComDeprecation(t *testing.T) {
	s, home, _ := newAssetsServer(t)
	store, err := studiodoc.NewStore(filepath.Join(home, ".aidd", "canonical", caminhos.DirLoops))
	if err != nil {
		t.Fatal(err)
	}
	s.Processos = store

	doc := studiodoc.Processo{
		ID:      "alias",
		Name:    "Alias",
		Version: "2.0.0",
		Trigger: &dag.Trigger{Type: dag.TriggerManual},
		Nodes: []dag.Node{{
			ID:      "n1",
			Type:    dag.TypeScript,
			Command: []string{"python", "run.py"},
			Routing: dag.RoutingConfig{Mode: dag.ModeManual, Model: "local"},
			Output:  dag.HandoffContract{Artifact: "RES.md", Gate: "file_exists('RES.md')"},
		}},
	}
	payload, _ := json.Marshal(doc)
	if rec := do(t, s.Handler(), "POST", "/api/processos", string(payload)); rec.Code != http.StatusOK {
		t.Fatalf("POST esperado 200, obteve %d: %s", rec.Code, rec.Body.String())
	}

	novo := do(t, s.Handler(), "GET", "/api/processos/alias", "")
	antigo := do(t, s.Handler(), "GET", "/api/loops/alias", "")
	if antigo.Code != novo.Code || antigo.Body.String() != novo.Body.String() {
		t.Fatalf("rota antiga diverge: %d %q vs %d %q", antigo.Code, antigo.Body.String(), novo.Code, novo.Body.String())
	}
	if antigo.Header().Get("Deprecation") == "" {
		t.Error("rota antiga deve enviar o cabeçalho Deprecation")
	}
	if novo.Header().Get("Deprecation") != "" {
		t.Error("rota nova não deve enviar o cabeçalho Deprecation")
	}
}

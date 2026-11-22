package studio

import (
	"bufio"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestRunAPIRequiresConsentForScripts(t *testing.T) {
	s, _ := newServer(t)
	h := s.Handler()

	body := `{
		"dag": {
			"name": "script-dag",
			"version": "2.0.0",
			"trigger": {"type": "manual"},
			"nodes": [
				{
					"id": "a",
					"type": "script",
					"command": ["echo", "oi"],
					"routing": {"mode": "manual", "model": "none"},
					"output": {"artifact": "a.txt", "gate": "true"}
				}
			]
		},
		"allow_scripts": false
	}`

	rec := do(t, h, "POST", "/api/run", body)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperava status 400 sem consentimento de script, obteve %d (%s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "consentimento") {
		t.Fatalf("mensagem deve citar consentimento: %s", rec.Body.String())
	}
}

func TestRunAPISSEStreaming(t *testing.T) {
	s, _ := newServer(t)
	h := s.Handler()

	body := `{
		"dag": {
			"name": "stream-dag",
			"version": "2.0.0",
			"trigger": {"type": "manual"},
			"nodes": [
				{
					"id": "node1",
					"type": "script",
					"command": ["go", "version"],
					"routing": {"mode": "manual", "model": "none"},
					"output": {"artifact": "res.txt", "gate": "true"}
				}
			]
		},
		"allow_scripts": true
	}`

	rec := do(t, h, "POST", "/api/run", body)

	if rec.Code != http.StatusOK {
		t.Fatalf("status inesperado: %d (%s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Header().Get("Content-Type"), "text/event-stream") {
		t.Fatalf("esperava text/event-stream, veio %q", rec.Header().Get("Content-Type"))
	}

	scanner := bufio.NewScanner(rec.Body)
	var events []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			events = append(events, strings.TrimPrefix(line, "data: "))
		}
	}

	if len(events) < 2 {
		t.Fatalf("esperava eventos SSE emitidos, veio %d: %v", len(events), events)
	}

	var firstEv map[string]any
	if err := json.Unmarshal([]byte(events[0]), &firstEv); err != nil {
		t.Fatalf("falha ao parsear evento json: %v", err)
	}
	if firstEv["type"] != "node_start" && firstEv["type"] != "dag_start" {
		t.Fatalf("evento inicial inválido: %v", firstEv)
	}

	var finishEv map[string]any
	lastRaw := events[len(events)-1]
	if err := json.Unmarshal([]byte(lastRaw), &finishEv); err != nil {
		t.Fatalf("falha ao parsear último evento json: %v", err)
	}
	if finishEv["type"] != "dag_finish" {
		t.Fatalf("esperava último evento dag_finish, obteve %v", finishEv["type"])
	}
	if _, ok := finishEv["duration_ms"]; !ok {
		t.Fatalf("dag_finish deve conter duration_ms")
	}
	if _, ok := finishEv["cost_est"]; !ok {
		t.Fatalf("dag_finish deve conter cost_est")
	}
}

func TestRunAPICoexistsWithSessaoA2UI(t *testing.T) {
	s, _ := newServer(t)
	h := s.Handler()

	recSessao := do(t, h, "POST", "/api/sessao", `{"processo":"dag-canvas"}`)
	if recSessao.Code != http.StatusOK && recSessao.Code != http.StatusCreated {
		t.Fatalf("POST /api/sessao falhou: %d (%s)", recSessao.Code, recSessao.Body.String())
	}
	var sessResp map[string]string
	if err := json.Unmarshal(recSessao.Body.Bytes(), &sessResp); err != nil {
		t.Fatalf("json inválido em /api/sessao: %v", err)
	}
	if sessResp["session_id"] == "" || !strings.Contains(sessResp["url"], "/sessao/") {
		t.Fatalf("resposta de sessão incompleta: %v", sessResp)
	}

	runBody := `{"dag":{"name":"coexist","version":"2.0.0","trigger":{"type":"manual"},"nodes":[{"id":"n1","type":"script","command":["go","version"],"routing":{"mode":"manual","model":"none"},"output":{"artifact":"a.txt","gate":"true"}}]},"allow_scripts":true}`
	recRun := do(t, h, "POST", "/api/run", runBody)
	if recRun.Code != http.StatusOK {
		t.Fatalf("POST /api/run deve coexistir com /api/sessao: status %d (%s)", recRun.Code, recRun.Body.String())
	}
}

func TestAppJSSessaoA2UIPreservesRunAbortController(t *testing.T) {
	js := readStatic(t, "app.js")
	checks := []string{
		"function criarEAbrirSessaoA2UI(",
		"/api/sessao",
		"/sessao/",
		"criarEAbrirSessaoA2UI(",
		"new AbortController()",
		`fetch("/api/run"`,
		"signal: runAbortController.signal",
	}
	for _, needle := range checks {
		if !strings.Contains(js, needle) {
			t.Errorf("app.js deve conter %q para AC-E04.3 e BUG-RUN-02", needle)
		}
	}
	start := strings.Index(js, "function criarEAbrirSessaoA2UI(")
	if start >= 0 {
		end := start + 1200
		if end > len(js) {
			end = len(js)
		}
		trecho := js[start:end]
		for _, proibido := range []string{"innerHTML", "alert(", "confirm("} {
			if strings.Contains(trecho, proibido) {
				t.Errorf("criarEAbrirSessaoA2UI não pode usar %q", proibido)
			}
		}
	}
}

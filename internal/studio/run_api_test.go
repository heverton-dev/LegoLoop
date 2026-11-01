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
}

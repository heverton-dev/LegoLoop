package interview

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// O 9Router responde em streaming (SSE) por padrão; o cliente só sabe ler JSON. Sem "stream": false o corpo
// chega como "data: {...}" e a decodificação falha com "invalid character 'd'".
func TestRouterClient_AsksForANonStreamingResponse(t *testing.T) {
	var body map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"ok"}}]}`))
	}))
	defer srv.Close()

	got, err := NewRouterClient(srv.URL, "k", "code-pro").Chat(context.Background(), []Message{{Role: "user", Content: "oi"}})
	if err != nil || got != "ok" {
		t.Fatalf("Chat = %q, %v", got, err)
	}
	if v, ok := body["stream"]; !ok || v != false {
		t.Errorf(`a requisição deve levar "stream": false, veio %v (presente=%v)`, v, ok)
	}
}

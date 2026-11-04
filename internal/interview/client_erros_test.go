package interview

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func chatVia(t *testing.T, h http.HandlerFunc, key string) (string, error) {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return NewRouterClient(srv.URL, key, "m").Chat(context.Background(), []Message{{Role: "user", Content: "oi"}})
}

func TestChatEnviaChaveEPedeJSON(t *testing.T) {
	var auth string
	var body map[string]any
	_, _ = chatVia(t, func(w http.ResponseWriter, r *http.Request) {
		auth = r.Header.Get("Authorization")
		_ = json.NewDecoder(r.Body).Decode(&body)
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"ok"}}]}`))
	}, "segredo")
	if auth != "Bearer segredo" {
		t.Errorf("Authorization = %q", auth)
	}
	if body["stream"] != false || body["model"] != "m" {
		t.Errorf("corpo inesperado: %v", body)
	}
}

func TestChatSemChaveNaoEnviaAuthorization(t *testing.T) {
	var has bool
	_, _ = chatVia(t, func(w http.ResponseWriter, r *http.Request) {
		_, has = r.Header["Authorization"]
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"ok"}}]}`))
	}, "")
	if has {
		t.Error("sem chave não deve haver cabeçalho Authorization")
	}
}

func TestChatErros(t *testing.T) {
	for _, c := range []struct {
		nome, want string
		h          http.HandlerFunc
	}{
		{"status 500", "status 500", func(w http.ResponseWriter, _ *http.Request) { http.Error(w, "quebrou", 500) }},
		{"choices vazias", "choices vazias", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(`{"choices":[]}`)) }},
		{"JSON inválido", "decodificar", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(`não é json`)) }},
	} {
		t.Run(c.nome, func(t *testing.T) {
			if _, err := chatVia(t, c.h, ""); err == nil || !strings.Contains(err.Error(), c.want) {
				t.Errorf("erro = %v, queria conter %q", err, c.want)
			}
		})
	}
	if _, err := (&RouterLLMClient{}).Chat(context.Background(), nil); err == nil {
		t.Error("sem BaseURL deveria falhar")
	}
}

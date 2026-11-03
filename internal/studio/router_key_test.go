package studio

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"legoloop/internal/provider"
)

// As rotas /api/interview, /api/spec e /api/plan montavam o cliente do 9Router sem a chave, e um router
// remoto responde 401. Elas devem seguir a mesma regra do resto do projeto: chave só para a base confiável.
func TestDefaultRouterClientsSendTheTrustedKey(t *testing.T) {
	var auth atomic.Value
	auth.Store("")
	router := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth.Store(r.Header.Get("Authorization"))
		http.Error(w, "irrelevante para este teste", http.StatusTeapot)
	}))
	defer router.Close()
	t.Setenv("NINEROUTER_URL", router.URL)
	t.Setenv("NINEROUTER_API_KEY", "segredo-studio")

	cases := map[string]string{
		"/api/interview/start": `{"idea":"extrator"}`,
	}
	for path, body := range cases {
		s, _ := newServer(t)
		s.Discover = func(context.Context) provider.Inventory { return provider.Inventory{RouterURL: router.URL} }
		auth.Store("")
		do(t, s.Handler(), "POST", path, body)
		if got := auth.Load().(string); got != "Bearer segredo-studio" {
			t.Errorf("%s enviou Authorization %q; sem a chave o router remoto responde 401", path, got)
		}
	}
}

package cli

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"legoloop/internal/brain"
	"legoloop/internal/fabricar"
	"legoloop/internal/provider"
)

// fakeRouter é um 9Router mínimo: registra o cabeçalho Authorization e devolve a SPEC no formato OpenAI.
type fakeRouter struct {
	*httptest.Server
	auth  atomic.Value // último Authorization recebido
	model atomic.Value // modelo pedido na última requisição
	calls atomic.Int32
}

func newFakeRouter(t *testing.T) *fakeRouter {
	t.Helper()
	fr := &fakeRouter{}
	fr.auth.Store("")
	fr.model.Store("")
	fr.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fr.calls.Add(1)
		fr.auth.Store(r.Header.Get("Authorization"))
		var req struct{ Model string }
		_ = json.NewDecoder(r.Body).Decode(&req)
		fr.model.Store(req.Model)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"choices": []map[string]any{{"message": map[string]string{"content": specFake}}}})
	}))
	t.Cleanup(fr.Close)
	return fr
}

// routerRig roda o fabricar sem LLM injetado, então o cliente do 9Router é montado pela CLI.
func routerRig(t *testing.T, inventoryURL string) (*fabRig, int) {
	t.Helper()
	r := newFabRig(t)
	r.d.Discover = func(context.Context) provider.Inventory {
		return provider.Inventory{RouterURL: inventoryURL, RouterModels: []string{"code-pro"}}
	}
	r.d.Fabricar = &FabricarDeps{Matriz: matrizCLI(), Clock: func() time.Time { return fabCLINow }}
	return r, run(t, r.d, r.args("extrator", "--yes", "--sem-entrevista", "--model", "code-pro")...)
}

func TestFabricarSendsTheRouterKeyOnlyToTheTrustedRouter(t *testing.T) {
	t.Run("a chave vai ao router confiavel", func(t *testing.T) {
		router := newFakeRouter(t)
		t.Setenv("NINEROUTER_URL", router.URL)
		t.Setenv("NINEROUTER_API_KEY", "segredo-123")
		r, code := routerRig(t, router.URL)
		if code != 0 {
			t.Fatalf("exit %d: %s", code, r.errb)
		}
		if got := router.auth.Load().(string); got != "Bearer segredo-123" {
			t.Errorf("Authorization = %q; sem a chave o router remoto responde 401", got)
		}
	})

	t.Run("a chave nunca vai a uma base que nao e a confiavel", func(t *testing.T) {
		trusted, other := newFakeRouter(t), newFakeRouter(t)
		t.Setenv("NINEROUTER_URL", trusted.URL)
		t.Setenv("NINEROUTER_API_KEY", "segredo-123")
		r, code := routerRig(t, other.URL) // o inventário apontou para OUTRA base
		if code != 0 {
			t.Fatalf("exit %d: %s", code, r.errb)
		}
		if got := other.auth.Load().(string); got != "" {
			t.Errorf("a chave vazou para uma base não confiável: %q", got)
		}
		if strings.Contains(r.out.String()+r.errb.String(), "segredo-123") {
			t.Error("a chave apareceu na saída")
		}
	})
}

// Gerar a SPEC é uma chamada longa (prompt grande, modelo de raciocínio): o prazo padrão de 45 s do cliente
// do router não basta. Sem isto o fabricar morre em "Client.Timeout exceeded" no 9Router remoto.
func TestFabricarLLMHasALongEnoughTimeout(t *testing.T) {
	router := newFakeRouter(t)
	d, _, _ := deps(t)
	d.Home = t.TempDir()
	b, _, err := brain.Pick(provider.Inventory{RouterURL: router.URL, RouterModels: []string{"code-pro"}}, "code-pro")
	if err != nil {
		t.Fatal(err)
	}
	fd, err := buildFabricarDeps(fabricarOpts{vault: t.TempDir(), docs: t.TempDir()}, fabricarReq(), d, b)
	if err != nil {
		t.Fatal(err)
	}
	br, ok := fd.LLM.(brain.Brain)
	if !ok {
		t.Fatalf("o LLM padrão deve ser um cérebro (brain.Brain), veio %T", fd.LLM)
	}
	if br.Timeout() < 2*time.Minute {
		t.Errorf("timeout de %s é curto para gerar a SPEC; esperado >= 2m", br.Timeout())
	}
}

func fabricarReq() fabricar.Request {
	return fabricar.Request{RawName: "x", Category: fabricar.CategoryBuild, ModelChoice: "code-pro"}
}

// O seletor de modelo é o "cérebro" que conduz a fabricação (INTENT, Objetivo 5): o modelo escolhido é o que
// o 9Router recebe, e não um combo fixo. Sem isto a escolha do usuário seria só um rótulo.
func TestFabricarUsesTheChosenModelAsTheBrain(t *testing.T) {
	router := newFakeRouter(t)
	t.Setenv("NINEROUTER_URL", router.URL)
	r := newFabRig(t)
	r.d.Discover = func(context.Context) provider.Inventory {
		return provider.Inventory{RouterURL: router.URL, RouterModels: []string{"code-fast", "code-pro"}}
	}
	r.d.Fabricar = &FabricarDeps{Matriz: matrizCLI(), Clock: func() time.Time { return fabCLINow }}
	if code := run(t, r.d, r.args("extrator", "--yes", "--sem-entrevista", "--model", "code-fast")...); code != 0 {
		t.Fatalf("exit %d: %s", code, r.errb)
	}
	if got := router.model.Load().(string); got != "code-fast" {
		t.Errorf("o router recebeu o modelo %q; esperado o escolhido (code-fast)", got)
	}
}

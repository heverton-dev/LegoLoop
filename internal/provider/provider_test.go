package provider

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func ollamaServer(t *testing.T, delay time.Duration) *httptest.Server {
	t.Helper()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/tags" || r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}
		time.Sleep(delay)
		w.Write([]byte(`{"models":[{"name":"qwen2.5-coder:7b"},{"name":"llama3:latest"}]}`))
	}))
	t.Cleanup(s.Close)
	return s
}

func routerServer(t *testing.T, wantKey string) *httptest.Server {
	t.Helper()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/models" {
			http.NotFound(w, r)
			return
		}
		if wantKey != "" && r.Header.Get("Authorization") != "Bearer "+wantKey {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Write([]byte(`{"object":"list","data":[{"id":"code-pro"},{"id":"code-fast"}]}`))
	}))
	t.Cleanup(s.Close)
	return s
}

func TestScanOllama(t *testing.T) {
	s := ollamaServer(t, 0)
	got, err := ScanOllama(context.Background(), s.Client(), s.URL)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []string{"llama3:latest", "qwen2.5-coder:7b"}) {
		t.Fatalf("got %v", got)
	}
}

func TestScanOllamaTimeout200ms(t *testing.T) {
	s := ollamaServer(t, 600*time.Millisecond)
	start := time.Now()
	if _, err := ScanOllama(context.Background(), s.Client(), s.URL); err == nil {
		t.Fatal("esperado timeout")
	}
	if el := time.Since(start); el > 450*time.Millisecond {
		t.Fatalf("timeout de 200ms não respeitado: %v", el)
	}
}

func TestScanOllamaErrors(t *testing.T) {
	bad := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`não-json`))
	}))
	defer bad.Close()
	if _, err := ScanOllama(context.Background(), bad.Client(), bad.URL); err == nil {
		t.Fatal("esperado erro de json")
	}
	status := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer status.Close()
	if _, err := ScanOllama(context.Background(), status.Client(), status.URL); err == nil {
		t.Fatal("esperado erro de status")
	}
	if _, err := ScanOllama(context.Background(), http.DefaultClient, "http://127.0.0.1:1"); err == nil {
		t.Fatal("esperado erro de conexão")
	}
	if _, err := ScanOllama(context.Background(), http.DefaultClient, "::bad"); err == nil {
		t.Fatal("esperado erro de URL")
	}
}

func TestDetectCLIs(t *testing.T) {
	look := func(name string) (string, error) {
		if name == "claude" || name == "opencode" {
			return `C:\bin\` + name + ".exe", nil
		}
		return "", errors.New("not found")
	}
	got := DetectCLIs(look, []string{"claude", "agy", "opencode"})
	want := []CLI{{Name: "claude", Path: `C:\bin\claude.exe`}, {Name: "opencode", Path: `C:\bin\opencode.exe`}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v", got)
	}
}

func TestScanNineRouterFirstReachableBase(t *testing.T) {
	open := routerServer(t, "")
	base, models, err := ScanNineRouter(context.Background(), open.Client(), []string{"http://127.0.0.1:1", open.URL}, "")
	if err != nil {
		t.Fatal(err)
	}
	if base != open.URL || !reflect.DeepEqual(models, []string{"code-fast", "code-pro"}) {
		t.Fatalf("base=%s models=%v", base, models)
	}
	s := routerServer(t, "segredo")
	if _, _, err := ScanNineRouter(context.Background(), s.Client(), []string{s.URL}, "segredo"); err != nil {
		t.Fatalf("base primária deve receber a credencial: %v", err)
	}
	// SEC-06: bases alternativas nunca recebem a credencial.
	if _, _, err := ScanNineRouter(context.Background(), s.Client(), []string{"http://127.0.0.1:1", s.URL}, "segredo"); err == nil {
		t.Fatal("credencial vazou para base alternativa")
	}
	if _, _, err := ScanNineRouter(context.Background(), s.Client(), []string{s.URL}, "errada"); err == nil {
		t.Fatal("esperado erro 401")
	}
	if _, _, err := ScanNineRouter(context.Background(), s.Client(), nil, ""); err == nil {
		t.Fatal("esperado erro sem bases")
	}
}

func TestDetectAPIKeysNeverLeaksValues(t *testing.T) {
	env := map[string]string{"ANTHROPIC_API_KEY": "sk-ant-xyz", "GEMINI_API_KEY": "AIza", "OPENAI_API_KEY": "  "}
	got := DetectAPIKeys(func(k string) string { return env[k] })
	if !reflect.DeepEqual(got, []string{"anthropic", "gemini"}) {
		t.Fatalf("got %v", got)
	}
}

func TestDiscoverInventory(t *testing.T) {
	o := ollamaServer(t, 0)
	r := routerServer(t, "")
	inv := Discover(context.Background(), Config{
		HTTP:       http.DefaultClient,
		OllamaURL:  o.URL,
		RouterURLs: []string{r.URL},
		LookPath: func(n string) (string, error) {
			if n == "claude" {
				return "/usr/bin/claude", nil
			}
			return "", errors.New("x")
		},
		Getenv: func(k string) string { return map[string]string{"ANTHROPIC_API_KEY": "k"}[k] },
	})
	if inv.RouterURL != r.URL || len(inv.Ollama) != 2 || len(inv.CLIs) != 1 || inv.APIKeys[0] != "anthropic" {
		t.Fatalf("inventário = %+v", inv)
	}
	for _, m := range []string{"code-pro", "qwen2.5-coder:7b", "sonnet", "haiku", "opus"} {
		if !inv.Has(m) {
			t.Errorf("modelo %q deveria estar disponível", m)
		}
	}
	if inv.Has("gpt-5") {
		t.Error("modelo inexistente reportado como disponível")
	}
	if src, ok := inv.Source("code-pro"); !ok || src != SourceRouter {
		t.Errorf("source code-pro = %v", src)
	}
	models := inv.Models()
	if len(models) != 7 || models[0].Source != SourceRouter {
		t.Fatalf("Models() deve priorizar combos do 9Router: %+v", models)
	}
}

func TestDiscoverOfflineRegraDeOuro(t *testing.T) {
	inv := Discover(context.Background(), Config{
		HTTP:       http.DefaultClient,
		OllamaURL:  "http://127.0.0.1:1",
		RouterURLs: []string{"http://127.0.0.1:1"},
		LookPath:   func(string) (string, error) { return "", errors.New("x") },
		Getenv:     func(string) string { return "" },
	})
	if len(inv.Models()) != 0 || inv.RouterURL != "" {
		t.Fatalf("offline deveria ter inventário vazio: %+v", inv)
	}
	if len(inv.Errors) == 0 || !strings.Contains(strings.Join(inv.Errors, " "), "ollama") {
		t.Fatalf("erros de scan devem ser reportados: %v", inv.Errors)
	}
}

func TestDefaultConfig(t *testing.T) {
	t.Setenv("USERPROFILE", t.TempDir())
	t.Setenv("HOME", t.TempDir())
	t.Setenv("NINEROUTER_URL", "")
	c := DefaultConfig()
	if c.HTTP == nil || c.LookPath == nil || c.Getenv == nil || c.OllamaURL == "" || len(c.RouterURLs) != 2 {
		t.Fatalf("DefaultConfig incompleto: %+v", c)
	}
	if c.TrustedRouterURL() != "http://127.0.0.1:20128" {
		t.Fatalf("base confiável = %s", c.TrustedRouterURL())
	}
	t.Setenv("NINEROUTER_URL", "http://127.0.0.1:20999/")
	c = DefaultConfig()
	if len(c.RouterURLs) != 1 || c.TrustedRouterURL() != "http://127.0.0.1:20999" {
		t.Fatalf("NINEROUTER_URL ignorada: %v", c.RouterURLs)
	}
	if (Config{}).TrustedRouterURL() != "" {
		t.Fatal("config vazia não tem base confiável")
	}
}

// R2-4: NINEROUTER_URL só vale para https ou loopback (credencial não trafega em claro pela rede).
func TestRouterURLMustBeSafe(t *testing.T) {
	for url, ok := range map[string]bool{
		"http://127.0.0.1:20999": true, "http://localhost:20128": true, "https://router.exemplo.com": true,
		"http://router.exemplo.com": false, "ftp://127.0.0.1": false, "lixo": false,
	} {
		t.Setenv("NINEROUTER_URL", url)
		got := DefaultConfig().TrustedRouterURL() == url
		if got != ok {
			t.Errorf("%s aceito=%v, esperado %v", url, got, ok)
		}
	}
}

// Harnesses no PATH viram modelos com prefixo (agy:…, omp:…, mimo:…, opencode:…).
func TestHarnessModelsArePrefixed(t *testing.T) {
	inv := Inventory{CLIs: []CLI{{Name: "agy"}, {Name: "omp"}}, Ollama: []string{"qwen2.5-coder:7b"}}
	for model, want := range map[string]Source{
		"agy:gemini-3.8-flash-low": SourceAgy, "agy:default": SourceAgy, "omp:code-fast": SourceOmp,
		"qwen2.5-coder:7b": SourceOllama,
	} {
		if got, ok := inv.Source(model); !ok || got != want {
			t.Errorf("%s ⇒ %v (%v), esperado %v", model, got, ok, want)
		}
	}
	for _, absent := range []string{"mimo:x", "opencode:y", "agy:", "zzz:m"} {
		if inv.Has(absent) {
			t.Errorf("%s não deveria estar disponível", absent)
		}
	}
	ids := map[string]bool{}
	for _, m := range inv.Models() {
		ids[m.ID] = true
	}
	if !ids["agy:default"] || !ids["omp:default"] || ids["mimo:default"] {
		t.Fatalf("Models() sem as entradas padrão dos harnesses: %v", ids)
	}
}

// A convenção do ambiente do usuário (NINEROUTER_KEY) é aceita, e "/v1" no fim da URL é tolerado.
func TestRouterEnvConventions(t *testing.T) {
	t.Setenv("NINEROUTER_API_KEY", "")
	t.Setenv("NINEROUTER_KEY", "k-teste")
	t.Setenv("NINEROUTER_URL", "https://9router.exemplo.org/v1/")
	c := DefaultConfig()
	if c.RouterKey != "k-teste" || c.TrustedRouterURL() != "https://9router.exemplo.org" {
		t.Fatalf("convenções ignoradas: base=%s", c.TrustedRouterURL())
	}
}

// Bases remotas (9Router na VPS, ~1,4 s via HTTPS) não podem usar o orçamento local de 200 ms.
func TestScanTimeoutByHost(t *testing.T) {
	cases := map[string]time.Duration{
		"http://127.0.0.1:11434":         ScanTimeout,
		"http://localhost:20128":         ScanTimeout,
		"http://[::1]:9000":              ScanTimeout,
		"https://9router.vpsconexao.org": RemoteScanTimeout,
	}
	for base, want := range cases {
		if got := scanTimeoutFor(base); got != want {
			t.Errorf("%s ⇒ %v, esperado %v", base, got, want)
		}
	}
	if c := DefaultConfig(); c.HTTP.Timeout < RemoteScanTimeout {
		t.Fatalf("cliente HTTP corta antes do orçamento remoto: %v", c.HTTP.Timeout)
	}
}

// SPEC_ORCHESTRIX_UNIVERSAL_MODELS: combos do 9Router (IDs sem "/") vêm no topo, como recomendados.
func TestRouterCombosFirst(t *testing.T) {
	got := sortRouterModels([]string{"ag/gemini-3-flash", "code-pro", "af/gpt-oss-20b", "code-fast", "megacombo-coding"})
	want := []string{"code-fast", "code-pro", "megacombo-coding", "af/gpt-oss-20b", "ag/gemini-3-flash"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ordem = %v, esperado %v", got, want)
	}
}

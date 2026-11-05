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
	"legoloop/internal/provider"
)

func invOllama() provider.Inventory { return provider.Inventory{Ollama: []string{"qwen2.5"}} }

// fakeOllama devolve a SPEC no formato do Ollama e registra o modelo pedido.
func fakeOllama(t *testing.T) (*httptest.Server, *atomic.Value, *atomic.Int32) {
	t.Helper()
	modelo, calls := &atomic.Value{}, &atomic.Int32{}
	modelo.Store("")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		var req struct{ Model string }
		_ = json.NewDecoder(r.Body).Decode(&req)
		modelo.Store(req.Model)
		_ = json.NewEncoder(w).Encode(map[string]any{"response": specFake})
	}))
	t.Cleanup(srv.Close)
	return srv, modelo, calls
}

// semModelo monta os argumentos sem o --model que o rig padrão acrescenta.
func (r *fabRig) semModelo(name string, extra ...string) []string {
	return append([]string{"fabricar", name, "--vault", r.vault, "--docs", r.doc}, extra...)
}

// AC-E01.7: modelo fora do inventário sai com 2, lista os disponíveis e não grava nada.
func TestFabricarModelInvalido_Sai2ListaOsModelosENaoGrava(t *testing.T) {
	for _, extra := range [][]string{{"--yes", "--sem-entrevista", "--dry-run"}, {"--yes", "--sem-entrevista"}} {
		r := newFabRig(t)
		r.d.Discover = func(context.Context) provider.Inventory {
			return provider.Inventory{RouterURL: "http://127.0.0.1:20128", RouterModels: []string{"code-fast", "code-pro"}, Ollama: []string{"qwen2.5"}}
		}
		code := run(t, r.d, append(r.semModelo("extrator", "--model", "inexistente"), extra...)...)
		if code != 2 {
			t.Fatalf("%v: exit %d, esperado 2", extra, code)
		}
		for _, quer := range []string{"inexistente", "code-fast", "code-pro", "qwen2.5"} {
			if !strings.Contains(r.errb.String(), quer) {
				t.Errorf("a mensagem deve citar %q:\n%s", quer, r.errb)
			}
		}
		if len(r.files(t)) != 0 || r.llm.calls != 0 {
			t.Errorf("nada pode ser gravado nem consultado: arquivos=%v chamadas=%d", r.files(t), r.llm.calls)
		}
	}
}

func TestFabricarModelInvalido_SemNenhumCerebroExplicaOQueFalta(t *testing.T) {
	r := newFabRig(t)
	r.d.Discover = func(context.Context) provider.Inventory { return provider.Inventory{} }
	if code := run(t, r.d, r.semModelo("extrator", "--yes", "--sem-entrevista")...); code != 1 {
		t.Fatalf("exit %d, esperado 1", code)
	}
	for _, quer := range []string{"9Router", "Ollama", "harness"} {
		if !strings.Contains(r.errb.String(), quer) {
			t.Errorf("deve dizer como resolver (%s):\n%s", quer, r.errb)
		}
	}
	if len(r.files(t)) != 0 {
		t.Errorf("nada pode ser gravado: %v", r.files(t))
	}
}

// O 9Router deixa de ser requisito: só com Ollama a fabricação roda de ponta a ponta.
func TestFabricarCerebro_OllamaSemRouterFabricaDePontaAPonta(t *testing.T) {
	srv, modelo, calls := fakeOllama(t)
	r := newFabRig(t)
	r.d.Discover = func(context.Context) provider.Inventory { return invOllama() }
	r.d.Fabricar = &FabricarDeps{Matriz: matrizCLI(), Clock: func() time.Time { return fabCLINow }, BrainOpts: []brain.Option{brain.WithOllamaURL(srv.URL)}}
	if code := run(t, r.d, r.semModelo("extrator", "--yes", "--sem-entrevista", "--model", "qwen2.5")...); code != 0 {
		t.Fatalf("exit %d: %s", code, r.errb)
	}
	if modelo.Load().(string) != "qwen2.5" || calls.Load() == 0 {
		t.Errorf("o Ollama deveria conduzir a fabricação (modelo=%v chamadas=%d)", modelo.Load(), calls.Load())
	}
	if len(r.files(t)) != 9 {
		t.Errorf("devia gravar 9 arquivos: %v", r.files(t))
	}
}

func TestFabricarCerebro_SemModeloEscolheOPrimeiroEExplica(t *testing.T) {
	srv, modelo, _ := fakeOllama(t)
	r := newFabRig(t)
	r.d.Discover = func(context.Context) provider.Inventory { return invOllama() }
	r.d.Fabricar = &FabricarDeps{Matriz: matrizCLI(), Clock: func() time.Time { return fabCLINow }, BrainOpts: []brain.Option{brain.WithOllamaURL(srv.URL)}}
	if code := run(t, r.d, r.semModelo("extrator", "--yes", "--sem-entrevista")...); code != 0 {
		t.Fatalf("exit %d: %s", code, r.errb)
	}
	if !strings.Contains(r.out.String(), "Cérebro escolhido: ollama/qwen2.5.") || modelo.Load().(string) != "qwen2.5" {
		t.Errorf("deve escolher o primeiro e dizer qual:\n%s", r.out)
	}
}

func TestFabricarCerebro_JSONMantemOStdoutPuroEExplicaNoStderr(t *testing.T) {
	srv, _, _ := fakeOllama(t)
	r := newFabRig(t)
	r.d.Discover = func(context.Context) provider.Inventory { return invOllama() }
	r.d.Fabricar = &FabricarDeps{Matriz: matrizCLI(), Clock: func() time.Time { return fabCLINow }, BrainOpts: []brain.Option{brain.WithOllamaURL(srv.URL)}}
	if code := run(t, r.d, r.semModelo("extrator", "--yes", "--sem-entrevista", "--json")...); code != 0 {
		t.Fatalf("exit %d: %s", code, r.errb)
	}
	var doc map[string]any
	if err := json.Unmarshal(r.out.Bytes(), &doc); err != nil {
		t.Fatalf("o stdout do --json tem que ser um JSON só: %v\n%s", err, r.out)
	}
	if !strings.Contains(r.errb.String(), "Cérebro escolhido: ollama/qwen2.5.") {
		t.Errorf("a explicação vai ao stderr:\n%s", r.errb)
	}
}

// Contrato novo: o modelo do claude CLI (ou de qualquer executor do inventário) é aceito; o 9Router não é requisito.
func TestFabricarCerebro_ModeloDoClaudeCLIEAceitoSemRouter(t *testing.T) {
	r := newFabRig(t) // inventário só com o claude CLI e LLM injetado
	if code := run(t, r.d, r.args("extrator", "--yes", "--sem-entrevista", "--model", "sonnet")...); code != 0 {
		t.Fatalf("exit %d: %s", code, r.errb)
	}
	if strings.Contains(r.errb.String(), "9Router") || !strings.Contains(r.out.String(), "Cérebro escolhido: claude/sonnet, por pedido seu.") {
		t.Errorf("stdout:\n%s\nstderr:\n%s", r.out, r.errb)
	}
}

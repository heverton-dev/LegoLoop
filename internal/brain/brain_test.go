package brain_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"legoloop/internal/brain"
	"legoloop/internal/interview"
	"legoloop/internal/provider"
	"legoloop/internal/ste"
)

func cli(name string) provider.CLI { return provider.CLI{Name: name, Path: "/bin/" + name} }

// AC-E01.6: o cérebro é escolhido pelo inventário e a resposta explica o porquê.
func TestPick_EscolhaAutomaticaPorInventario(t *testing.T) {
	cases := []struct {
		nome           string
		inv            provider.Inventory
		nome1, explica string
	}{
		{"só Ollama", provider.Inventory{Ollama: []string{"qwen2.5"}}, "ollama/qwen2.5", "Ollama"},
		{"9Router e Ollama", provider.Inventory{RouterURL: "http://127.0.0.1:20128", RouterModels: []string{"code-fast", "code-pro"}, Ollama: []string{"qwen2.5"}}, "9router/code-fast", "9Router"},
		{"só claude CLI", provider.Inventory{CLIs: []provider.CLI{cli("claude")}}, "claude/sonnet", "claude"},
		{"só harness agy", provider.Inventory{CLIs: []provider.CLI{cli("agy")}}, "agy:default", "agy"},
		{"Ollama antes do claude CLI", provider.Inventory{Ollama: []string{"qwen2.5"}, CLIs: []provider.CLI{cli("claude")}}, "ollama/qwen2.5", "Ollama"},
	}
	for _, c := range cases {
		t.Run(c.nome, func(t *testing.T) {
			b, why, err := brain.Pick(c.inv, "")
			if err != nil {
				t.Fatal(err)
			}
			if b.Name() != c.nome1 {
				t.Errorf("Name() = %q, quer %q", b.Name(), c.nome1)
			}
			if !strings.Contains(why, c.explica) || !strings.Contains(why, c.nome1) {
				t.Errorf("a explicação deve citar o executor e o modelo: %q", why)
			}
			if vs := ste.Lint(why); len(vs) != 0 {
				t.Errorf("a explicação deve seguir o perfil STE-PT: %v", vs)
			}
		})
	}
}

func TestPick_SemCerebroExplicaOQueFalta(t *testing.T) {
	_, _, err := brain.Pick(provider.Inventory{}, "")
	if !errors.Is(err, brain.ErrSemCerebro) {
		t.Fatalf("esperava ErrSemCerebro: %v", err)
	}
	msg := err.Error()
	for _, quer := range []string{"9Router", "Ollama", "harness"} {
		if !strings.Contains(msg, quer) {
			t.Errorf("a mensagem deve dizer como resolver (%s): %q", quer, msg)
		}
	}
	if vs := ste.Lint(msg); len(vs) != 0 {
		t.Errorf("a mensagem deve seguir o perfil STE-PT: %v", vs)
	}
}

func TestPick_ModeloPedidoPeloUsuario(t *testing.T) {
	inv := provider.Inventory{RouterURL: "http://127.0.0.1:20128", RouterModels: []string{"code-fast"}, Ollama: []string{"qwen2.5"}, CLIs: []provider.CLI{cli("claude")}}
	b, why, err := brain.Pick(inv, "qwen2.5")
	if err != nil || b.Name() != "ollama/qwen2.5" {
		t.Fatalf("o pedido do usuário vale sobre a ordem automática: %v %v", b, err)
	}
	if !strings.Contains(why, "pedido") {
		t.Errorf("a explicação deve dizer que foi o usuário quem pediu: %q", why)
	}
	b, _, err = brain.Pick(inv, "opus")
	if err != nil || b.Name() != "claude/opus" {
		t.Fatalf("modelo do claude CLI: %v %v", b, err)
	}
}

func TestPick_ModeloForaDoInventarioListaOsDisponiveis(t *testing.T) {
	inv := provider.Inventory{RouterURL: "http://127.0.0.1:20128", RouterModels: []string{"code-fast", "code-pro"}, Ollama: []string{"qwen2.5"}}
	_, _, err := brain.Pick(inv, "inexistente")
	if !errors.Is(err, brain.ErrModeloDesconhecido) {
		t.Fatalf("esperava ErrModeloDesconhecido: %v", err)
	}
	for _, quer := range []string{"inexistente", "code-fast", "code-pro", "qwen2.5"} {
		if !strings.Contains(err.Error(), quer) {
			t.Errorf("o erro deve citar %q: %v", quer, err)
		}
	}
	// com inventário vazio o modelo pedido também falha, mas por falta de cérebro
	if _, _, err := brain.Pick(provider.Inventory{}, "x"); !errors.Is(err, brain.ErrSemCerebro) {
		t.Errorf("sem inventário o erro certo é ErrSemCerebro: %v", err)
	}
}

func TestPick_OllamaConversaPeloOrchestrix(t *testing.T) {
	var corpo map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &corpo)
		_, _ = w.Write([]byte(`{"response":"resposta do ollama","prompt_eval_count":3,"eval_count":4}`))
	}))
	defer srv.Close()
	b, _, err := brain.Pick(provider.Inventory{Ollama: []string{"qwen2.5"}}, "", brain.WithOllamaURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	got, err := b.Chat(context.Background(), []interview.Message{{Role: "system", Content: "Seja breve."}, {Role: "user", Content: "Olá"}})
	if err != nil || got != "resposta do ollama" {
		t.Fatalf("Chat = %q, %v", got, err)
	}
	prompt, _ := corpo["prompt"].(string)
	if corpo["model"] != "qwen2.5" || !strings.Contains(prompt, "Seja breve.") || !strings.Contains(prompt, "Olá") {
		t.Errorf("o pedido ao Ollama perdeu conteúdo: %v", corpo)
	}
	if strings.Index(prompt, "Seja breve.") > strings.Index(prompt, "Olá") {
		t.Errorf("a ordem das mensagens foi trocada: %q", prompt)
	}
}

func TestPick_RouterSemVazarAChaveEComStreamFalso(t *testing.T) {
	var auth string
	var corpo map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth = r.Header.Get("Authorization")
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &corpo)
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"resposta do router"}}]}`))
	}))
	defer srv.Close()
	inv := provider.Inventory{RouterURL: srv.URL, RouterModels: []string{"code-fast"}}
	b, _, err := brain.Pick(inv, "")
	if err != nil {
		t.Fatal(err)
	}
	got, err := b.Chat(context.Background(), []interview.Message{{Role: "system", Content: "S"}, {Role: "user", Content: "U"}})
	if err != nil || got != "resposta do router" {
		t.Fatalf("Chat = %q, %v", got, err)
	}
	if auth != "" {
		t.Errorf("a chave do roteador só vai para a URL confiável; este servidor recebeu %q", auth)
	}
	if corpo["stream"] != false || corpo["model"] != "code-fast" {
		t.Errorf("corpo do pedido: %v", corpo)
	}
	msgs, _ := corpo["messages"].([]any)
	if len(msgs) != 2 {
		t.Errorf("o router recebe as mensagens com os papéis originais: %v", corpo["messages"])
	}
}

func TestPick_ErroDoExecutorChegaAoChamador(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "quebrou", http.StatusInternalServerError)
	}))
	defer srv.Close()
	b, _, _ := brain.Pick(provider.Inventory{Ollama: []string{"qwen2.5"}}, "", brain.WithOllamaURL(srv.URL))
	if _, err := b.Chat(context.Background(), []interview.Message{{Role: "user", Content: "x"}}); err == nil {
		t.Fatal("falha do executor não pode virar resposta vazia")
	}
}

func TestPick_CatalogoDeMensagensSemViolacao(t *testing.T) {
	if vs := ste.Default.Violations(); len(vs) != 0 {
		t.Fatalf("catálogo STE-PT com violações: %v", vs)
	}
}

func TestPick_ModelETimeoutDeCadaExecutor(t *testing.T) {
	cases := []struct {
		nome  string
		inv   provider.Inventory
		model string
	}{
		{"router", provider.Inventory{RouterURL: "http://127.0.0.1:20128", RouterModels: []string{"code-fast"}}, "code-fast"},
		{"ollama", provider.Inventory{Ollama: []string{"qwen2.5"}}, "qwen2.5"},
		{"claude", provider.Inventory{CLIs: []provider.CLI{cli("claude")}}, "sonnet"},
		{"harness", provider.Inventory{CLIs: []provider.CLI{cli("agy")}}, "agy:default"},
	}
	for _, c := range cases {
		t.Run(c.nome, func(t *testing.T) {
			b, _, err := brain.Pick(c.inv, "")
			if err != nil {
				t.Fatal(err)
			}
			if b.Model() != c.model {
				t.Errorf("Model() = %q, quer %q (o id do modelo, sem o executor)", b.Model(), c.model)
			}
			if b.Timeout() != brain.DefaultTimeout || brain.DefaultTimeout < 2*time.Minute {
				t.Errorf("Timeout() = %v; gerar a SPEC é uma chamada longa (padrão %v)", b.Timeout(), brain.DefaultTimeout)
			}
			b2, _, _ := brain.Pick(c.inv, "", brain.WithTimeout(7*time.Second))
			if b2.Timeout() != 7*time.Second {
				t.Errorf("WithTimeout não chegou ao cérebro: %v", b2.Timeout())
			}
		})
	}
}

func TestPick_ClienteHTTPSemTimeoutNaoCancelaAChamada(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"response":"ok"}`))
	}))
	defer srv.Close()
	b, _, _ := brain.Pick(provider.Inventory{Ollama: []string{"qwen2.5"}}, "", brain.WithOllamaURL(srv.URL), brain.WithHTTP(&http.Client{}))
	got, err := b.Chat(context.Background(), []interview.Message{{Role: "user", Content: "x"}})
	if err != nil || got != "ok" {
		t.Fatalf("um *http.Client sem Timeout não pode virar prazo zero: %q %v", got, err)
	}
	if b.Timeout() != brain.DefaultTimeout {
		t.Errorf("o prazo do cérebro vem da configuração, não do cliente HTTP: %v", b.Timeout())
	}
}

func TestPick_ListaDeModelosLongaVemResumida(t *testing.T) {
	var ms []string
	for i := 1; i <= 120; i++ {
		ms = append(ms, fmt.Sprintf("modelo-%03d", i))
	}
	inv := provider.Inventory{RouterURL: "http://127.0.0.1:20128", RouterModels: ms}
	_, _, err := brain.Pick(inv, "inexistente")
	if !errors.Is(err, brain.ErrModeloDesconhecido) {
		t.Fatal(err)
	}
	msg := err.Error()
	if !strings.Contains(msg, "modelo-001") || strings.Contains(msg, "modelo-120") {
		t.Errorf("deve mostrar só os primeiros modelos: %.300s", msg)
	}
	if !strings.Contains(msg, "e mais 105") || !strings.Contains(msg, "legoloop doctor") {
		t.Errorf("deve dizer quantos faltam e onde ver todos: %.400s", msg)
	}
	if len(msg) > 600 {
		t.Errorf("mensagem longa demais (%d caracteres) para ler de uma vez", len(msg))
	}
	if vs := ste.Lint(msg); len(vs) != 0 {
		t.Errorf("a mensagem deve seguir o perfil STE-PT: %v", vs)
	}
	// com poucos modelos, lista todos e não menciona resumo
	_, _, err = brain.Pick(provider.Inventory{Ollama: []string{"a", "b"}}, "x")
	if strings.Contains(err.Error(), "e mais") {
		t.Errorf("lista curta não leva resumo: %v", err)
	}
}

func TestPick_SuporteAEffortEInferencia(t *testing.T) {
	inv := provider.Inventory{RouterURL: "http://127.0.0.1:20128", RouterModels: []string{"code-fast", "code-pro"}}
	b, _, err := brain.Pick(inv, "code-pro", brain.WithEffort(brain.EffortHigh))
	if err != nil {
		t.Fatalf("erro ao escolher modelo com effort: %v", err)
	}
	if b.Effort() != brain.EffortHigh {
		t.Errorf("Effort() = %q, quer %q", b.Effort(), brain.EffortHigh)
	}
}

func TestPick_SelecaoInteligentePorEsforco(t *testing.T) {
	inv := provider.Inventory{RouterURL: "http://127.0.0.1:20128", RouterModels: []string{"code-fast", "code-pro"}}
	// Com EffortLow, prefere code-fast
	bLow, _, err := brain.Pick(inv, "", brain.WithEffort(brain.EffortLow))
	if err != nil {
		t.Fatalf("erro com effort low: %v", err)
	}
	if bLow.Model() != "code-fast" {
		t.Errorf("EffortLow deveria escolher code-fast, mas escolheu: %s", bLow.Model())
	}

	// Com EffortHigh, prefere code-pro
	bHigh, _, err := brain.Pick(inv, "", brain.WithEffort(brain.EffortHigh))
	if err != nil {
		t.Fatalf("erro com effort high: %v", err)
	}
	if bHigh.Model() != "code-pro" {
		t.Errorf("EffortHigh deveria escolher code-pro, mas escolheu: %s", bHigh.Model())
	}
}

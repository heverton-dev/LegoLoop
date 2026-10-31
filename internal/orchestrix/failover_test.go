package orchestrix

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/provider"
)

// O motivo devolvido pelo claude no JSON aparece no erro (achado no e2e real).
func TestClaudeErrorReasonIsSurfaced(t *testing.T) {
	o := &Orchestrix{Inventory: provider.Inventory{CLIs: []provider.CLI{{Name: "claude", Path: os.Args[0]}}},
		SkillsDir: skillsDir(t), ExtraEnv: []string{"LEGOLOOP_HELPER=claude-api-error"}}
	_, err := o.Execute(context.Background(), skillNode(dag.RoutingConfig{Mode: dag.ModeManual, Model: "haiku"}, nil), t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "Credit balance is too low") {
		t.Fatalf("motivo do claude não reportado: %v", err)
	}
}

// Modo sugerido: se o modelo principal falhar em execução, o fallback do inventário assume.
func TestSuggestedFallsBackAtRuntime(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"response":"# SPEC via fallback","prompt_eval_count":3,"eval_count":4}`))
	}))
	defer srv.Close()
	o := &Orchestrix{
		Inventory: provider.Inventory{CLIs: []provider.CLI{{Name: "claude", Path: os.Args[0]}}, Ollama: []string{"llama"}},
		HTTP:      srv.Client(), OllamaURL: srv.URL, SkillsDir: skillsDir(t),
		ExtraEnv: []string{"LEGOLOOP_HELPER=claude-api-error"},
	}
	work := t.TempDir()
	res, err := o.Execute(context.Background(), skillNode(dag.RoutingConfig{Mode: dag.ModeSuggested, Model: "haiku", Fallback: "llama"}, nil), work)
	if err != nil {
		t.Fatalf("fallback não assumiu: %v", err)
	}
	if readArtifact(t, work) != "# SPEC via fallback" || res.Tokens.Output != 4 {
		t.Fatalf("artefato/tokens do fallback incorretos: %+v", res)
	}
	// Rastreabilidade: quem respondeu e por que o principal foi descartado.
	if res.Model != "ollama/llama" || len(res.Notes) != 1 || !strings.Contains(res.Notes[0], "Credit balance is too low") {
		t.Fatalf("failover sem rastreabilidade: model=%q notes=%v", res.Model, res.Notes)
	}
}

// Modo manual é travado pelo usuário: nunca troca de modelo, mesmo com fallback declarado.
func TestManualNeverFallsBack(t *testing.T) {
	o := &Orchestrix{
		Inventory: provider.Inventory{CLIs: []provider.CLI{{Name: "claude", Path: os.Args[0]}}, Ollama: []string{"llama"}},
		SkillsDir: skillsDir(t), ExtraEnv: []string{"LEGOLOOP_HELPER=claude-api-error"},
	}
	_, err := o.Execute(context.Background(), skillNode(dag.RoutingConfig{Mode: dag.ModeManual, Model: "haiku", Fallback: "llama"}, nil), t.TempDir())
	if err == nil {
		t.Fatal("modo manual não pode cair para outro modelo")
	}
}

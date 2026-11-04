package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/fabricar"
	"legoloop/internal/interview"
)

// sysLLM guarda o prompt do sistema da primeira chamada da entrevista.
type sysLLM struct {
	system string
	inner  *chatLLM
}

func (s *sysLLM) Chat(ctx context.Context, msgs []interview.Message) (string, error) {
	if s.system == "" && !strings.Contains(msgs[0].Content, "Arquiteto de Especificações") {
		s.system = msgs[0].Content
	}
	return s.inner.Chat(ctx, msgs)
}

// TestContextoCLILeAgentsEObjetivo: o fabricar interativo entrega ao cérebro o AGENTS.md do projeto e o objetivo.
func TestContextoCLILeAgentsEObjetivo(t *testing.T) {
	r := newInterRig(t, "\n0\n", lowReply)
	spy := &sysLLM{inner: r.llm}
	r.d.Fabricar.LLM = spy
	if err := os.MkdirAll(r.doc, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(r.doc, "AGENTS.md"), []byte("Go, stdlib only"), 0o644); err != nil {
		t.Fatal(err)
	}
	r.fabricar(t, "--goal", "refinar a etapa de BUILD do mega-pipeline")
	for _, want := range []string{"Go, stdlib only", "refinar a etapa de BUILD do mega-pipeline"} {
		if !strings.Contains(spy.system, want) {
			t.Errorf("o prompt do sistema não contém %q:\n%s", want, spy.system)
		}
	}
}

func TestContextoCLIObjetivoCaiNoNome(t *testing.T) {
	for _, c := range []struct{ goal, name, want string }{
		{"  refinar o build ", "meu-fluxo", "refinar o build"},
		{"", "meu-fluxo", "meu-fluxo"},
		{"   ", " meu-fluxo ", "meu-fluxo"},
	} {
		if got := goalOf(fabricar.Request{Goal: c.goal, RawName: c.name}); got != c.want {
			t.Errorf("goalOf(%q, %q) = %q, queria %q", c.goal, c.name, got, c.want)
		}
	}
}

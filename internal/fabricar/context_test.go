package fabricar

import (
	"context"
	"strings"
	"testing"

	"legoloop/internal/interview"
)

type recChat struct{ first []interview.Message }

func (r *recChat) Chat(_ context.Context, msgs []interview.Message) (string, error) {
	if r.first == nil {
		r.first = append([]interview.Message(nil), msgs...)
	}
	return "HIPÓTESE: h\nCONFIANÇA: 40%\nPERGUNTA: Qual etapa?\nPALPITE: build", nil
}

// O Clarifier repassa o contexto do projeto à entrevista: a primeira chamada ao cérebro já o carrega.
func TestContextoClarifierRepassaAoCerebro(t *testing.T) {
	llm := &recChat{}
	c := NewClarifier(llm)
	c.Context = interview.Context{Goal: "refinar o BUILD", Docs: []interview.Snippet{{Name: "AGENTS.md", Text: "Go, stdlib only"}}}
	ask := func(string, string) (string, error) { return "", ErrAborted }
	_, _ = c.Clarify(context.Background(), Request{RawName: "x", Interactive: true}, ask)
	var all strings.Builder
	for _, m := range llm.first {
		all.WriteString(m.Content)
	}
	for _, want := range []string{"Go, stdlib only", "refinar o BUILD"} {
		if !strings.Contains(all.String(), want) {
			t.Errorf("a primeira chamada não carrega %q", want)
		}
	}
}

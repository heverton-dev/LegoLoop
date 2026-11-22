package plangen_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/interview"
	"github.com/heverton-dev/LegoLoop/internal/plangen"
	"github.com/heverton-dev/LegoLoop/internal/upstream"
)

// gravarIntencaoAprovada grava a INTENT que a SPEC de teste referencia, com o cabeçalho do Stop-Gate 01.
func gravarIntencaoAprovada(t *testing.T, base string) string {
	t.Helper()
	p := filepath.Join(base, "docs", "intents", "2026-11-02_INTENT_SERVICO.md")
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("---\nid: INTENT-SERVICO\nstatus: APPROVED\n---\n# INTENT\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

// specMinima é uma SPEC que só traz a Referência da Intenção exigida pela trava.
const specMinima = "# SPEC Servico\n\n> **Referência da Intenção:** [docs/intents/2026-11-02_INTENT_SERVICO.md](../intents/2026-11-02_INTENT_SERVICO.md)\n"

type llmContador struct{ chamadas int }

func (l *llmContador) Chat(context.Context, []interview.Message) (string, error) {
	l.chamadas++
	return threeTickets, nil
}

// TestGenerate_RecusaSpecSemIntencaoAprovada: o PLAN só nasce de uma SPEC cuja INTENT foi aprovada.
// A recusa vem antes do LLM e não grava PLAN nem ROUTE.
func TestGenerate_RecusaSpecSemIntencaoAprovada(t *testing.T) {
	base, spec := writeSpec(t)
	intencao := gravarIntencaoAprovada(t, base)
	llm := &llmContador{}
	gen := plangen.NewGenerator(llm)

	if err := os.WriteFile(intencao, []byte("---\nid: INTENT-SERVICO\nstatus: UNDER_REVIEW\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := gen.Generate(context.Background(), base, spec, plangen.SupervisionEpic); !errors.Is(err, upstream.ErrIntencaoNaoAprovada) {
		t.Fatalf("INTENT em revisão: err = %v, quer ErrIntencaoNaoAprovada", err)
	}
	if err := os.Remove(intencao); err != nil {
		t.Fatal(err)
	}
	if _, err := gen.Generate(context.Background(), base, spec, plangen.SupervisionEpic); !errors.Is(err, upstream.ErrIntencaoAusente) {
		t.Fatalf("INTENT ausente: err = %v, quer ErrIntencaoAusente", err)
	}
	if llm.chamadas != 0 {
		t.Fatalf("o LLM foi chamado %d vez(es) antes da recusa", llm.chamadas)
	}
	if _, err := os.Stat(filepath.Join(base, "docs", "plans")); err == nil {
		t.Fatal("a recusa não pode gravar o PLAN")
	}
}

package specgen_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/graph"
	"github.com/heverton-dev/LegoLoop/internal/interview"
	"github.com/heverton-dev/LegoLoop/internal/specgen"
	"github.com/heverton-dev/LegoLoop/internal/upstream"
)

type llmContador struct{ chamadas int }

func (l *llmContador) Chat(_ context.Context, _ []interview.Message) (string, error) {
	l.chamadas++
	return "", nil
}

// TestGenerateRecusaIntencaoNaoAprovada: a SPEC só nasce de uma INTENT aprovada no Stop-Gate 01.
// A recusa vem antes de qualquer chamada ao LLM e não grava nada.
func TestGenerateRecusaIntencaoNaoAprovada(t *testing.T) {
	base := t.TempDir()
	intencao := filepath.Join(base, "docs", "intents", "2026-11-21_INTENT_X.md")
	if err := os.MkdirAll(filepath.Dir(intencao), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(intencao, []byte("---\nid: INTENT-X\nstatus: UNDER_REVIEW\n---\n# INTENT\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	idx, err := graph.NewIndex(graph.DefaultCatalog())
	if err != nil {
		t.Fatal(err)
	}
	llm := &llmContador{}
	_, err = specgen.NewGenerator(llm, idx).Generate(context.Background(), base, intencao)
	if !errors.Is(err, upstream.ErrIntencaoNaoAprovada) {
		t.Fatalf("err = %v, quer ErrIntencaoNaoAprovada", err)
	}
	if llm.chamadas != 0 {
		t.Fatalf("o LLM foi chamado %d vez(es) antes da recusa", llm.chamadas)
	}
	if _, err := os.Stat(filepath.Join(base, "docs", "specs")); err == nil {
		t.Fatal("a recusa não pode gravar a SPEC")
	}
}

package upstream_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/upstream"
)

func escrever(t *testing.T, base, rel, texto string) string {
	t.Helper()
	p := filepath.Join(base, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(texto), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

const specComIntencao = "# SPEC\n\n> **Referência da Intenção:** [docs/intents/2026-11-21_INTENT_X.md](../intents/2026-11-21_INTENT_X.md)  \n"

func TestVerificarSpecRecusaSpecSemReferenciaDeIntencao(t *testing.T) {
	base := t.TempDir()
	spec := escrever(t, base, "docs/specs/2026-11-21_SPEC_X.md", "# SPEC sem referência\n")
	if err := upstream.VerificarSpec(base, spec); !errors.Is(err, upstream.ErrSemIntencao) {
		t.Fatalf("err = %v, quer ErrSemIntencao", err)
	}
}

func TestVerificarSpecRecusaIntencaoAusente(t *testing.T) {
	base := t.TempDir()
	spec := escrever(t, base, "docs/specs/2026-11-21_SPEC_X.md", specComIntencao)
	if err := upstream.VerificarSpec(base, spec); !errors.Is(err, upstream.ErrIntencaoAusente) {
		t.Fatalf("err = %v, quer ErrIntencaoAusente", err)
	}
}

func TestVerificarSpecRecusaIntencaoNaoAprovada(t *testing.T) {
	base := t.TempDir()
	spec := escrever(t, base, "docs/specs/2026-11-21_SPEC_X.md", specComIntencao)
	escrever(t, base, "docs/intents/2026-11-21_INTENT_X.md", "---\nid: INTENT-x\nstatus: UNDER_REVIEW\n---\n# INTENT\n")
	if err := upstream.VerificarSpec(base, spec); !errors.Is(err, upstream.ErrIntencaoNaoAprovada) {
		t.Fatalf("err = %v, quer ErrIntencaoNaoAprovada", err)
	}
}

func TestVerificarSpecAceitaIntencaoAprovadaComQuebraDeLinhaWindows(t *testing.T) {
	base := t.TempDir()
	spec := escrever(t, base, "docs/specs/2026-11-21_SPEC_X.md", specComIntencao)
	escrever(t, base, "docs/intents/2026-11-21_INTENT_X.md", "---\r\nid: INTENT-x\r\nstatus: approved\r\n---\r\n# INTENT\r\n")
	if err := upstream.VerificarSpec(base, spec); err != nil {
		t.Fatalf("INTENT aprovada deveria passar: %v", err)
	}
}

func TestIntencaoAprovadaLeSoOCabecalho(t *testing.T) {
	base := t.TempDir()
	p := escrever(t, base, "i.md", "---\nid: I\nstatus: DRAFT\n---\n\nstatus: APPROVED no corpo não vale\n")
	if err := upstream.IntencaoAprovada(p); !errors.Is(err, upstream.ErrIntencaoNaoAprovada) {
		t.Fatalf("o corpo não pode aprovar a INTENT: %v", err)
	}
}

func TestRecusaReconheceAsTresRecusasEIgnoraOutrosErros(t *testing.T) {
	for _, err := range []error{upstream.ErrSemIntencao, upstream.ErrIntencaoAusente, upstream.ErrIntencaoNaoAprovada} {
		if !upstream.Recusa(fmt.Errorf("envolvido: %w", err)) {
			t.Errorf("Recusa(%v) deveria ser verdadeiro", err)
		}
	}
	if upstream.Recusa(errors.New("falha de rede")) || upstream.Recusa(nil) {
		t.Error("erro comum e nil não são recusa da trava")
	}
}

func TestVerificarSpecPropagaErroDeLeituraDaSpec(t *testing.T) {
	base := t.TempDir()
	err := upstream.VerificarSpec(base, filepath.Join(base, "inexistente.md"))
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("err = %v, quer os.ErrNotExist", err)
	}
	if upstream.Recusa(err) {
		t.Fatalf("erro de leitura não pode ser tratado como recusa: %v", err)
	}
}

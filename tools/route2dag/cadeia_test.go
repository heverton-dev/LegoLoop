package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/plangen"
	"github.com/heverton-dev/LegoLoop/internal/upstream"
)

// cadeiaDeTeste grava INTENT, SPEC e PLAN no disco e devolve a ROUTE que assina os três.
func cadeiaDeTeste(t *testing.T, statusIntencao string, comIntencao bool) (base string, raw []byte) {
	t.Helper()
	base = t.TempDir()
	arq := map[string]string{
		"docs/intents/i.md": "---\nid: I\nstatus: " + statusIntencao + "\n---\n# INTENT\n",
		"docs/specs/s.md":   "# SPEC\n",
		"docs/plans/p.md":   "# PLAN\n",
	}
	for rel, txt := range arq {
		p := filepath.Join(base, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(txt), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	sig := map[string]*plangen.FileSignature{
		"spec": plangen.FileSignatureFor(base, "docs/specs/s.md"),
		"plan": plangen.FileSignatureFor(base, "docs/plans/p.md"),
	}
	if comIntencao {
		sig["intent"] = plangen.FileSignatureFor(base, "docs/intents/i.md")
	}
	raw, err := json.Marshal(map[string]any{"upstream_signatures": sig})
	if err != nil {
		t.Fatal(err)
	}
	return base, raw
}

func TestCadeiaRecusaRouteSemAssinaturaDeIntencao(t *testing.T) {
	base, raw := cadeiaDeTeste(t, "APPROVED", false)
	if err := VerificarCadeia(raw, base); !errors.Is(err, upstream.ErrSemIntencao) {
		t.Fatalf("err = %v, quer ErrSemIntencao", err)
	}
}

func TestCadeiaRecusaIntencaoNaoAprovada(t *testing.T) {
	base, raw := cadeiaDeTeste(t, "UNDER_REVIEW", true)
	if err := VerificarCadeia(raw, base); !errors.Is(err, upstream.ErrIntencaoNaoAprovada) {
		t.Fatalf("err = %v, quer ErrIntencaoNaoAprovada", err)
	}
}

func TestCadeiaRecusaArtefatoAdulterado(t *testing.T) {
	base, raw := cadeiaDeTeste(t, "APPROVED", true)
	if err := os.WriteFile(filepath.Join(base, "docs", "plans", "p.md"), []byte("adulterado"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := VerificarCadeia(raw, base); err == nil {
		t.Fatal("PLAN adulterado deveria quebrar a cadeia")
	}
}

func TestCadeiaAceitaIntencaoAprovadaEHashesQueConferem(t *testing.T) {
	base, raw := cadeiaDeTeste(t, "APPROVED", true)
	if err := VerificarCadeia(raw, base); err != nil {
		t.Fatalf("cadeia íntegra deveria passar: %v", err)
	}
}

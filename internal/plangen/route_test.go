package plangen_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/plangen"
)

func TestHashNormalized_EquivalenceCRLF(t *testing.T) {
	crlf := []byte("linha 1\r\nlinha 2\r\nlinha 3\r\n")
	lf := []byte("linha 1\nlinha 2\nlinha 3\n")

	hashCRLF := plangen.HashNormalized(crlf)
	hashLF := plangen.HashNormalized(lf)

	if hashCRLF != hashLF {
		t.Fatalf("hashes devem ser identicos para CRLF e LF, obteve %s vs %s", hashCRLF, hashLF)
	}
	if len(hashCRLF) != 64 {
		t.Fatalf("tamanho do hash SHA-256 deve ser 64 hex, obteve %d", len(hashCRLF))
	}
}

func TestVerifyRouteIntegrity_DetectsTampering(t *testing.T) {
	base := t.TempDir()
	specRel := "docs/specs/2026-11-07_SPEC_TEST.md"
	specFull := filepath.Join(base, filepath.FromSlash(specRel))
	if err := os.MkdirAll(filepath.Dir(specFull), 0o755); err != nil {
		t.Fatal(err)
	}
	originalContent := "conteudo original da spec\n"
	if err := os.WriteFile(specFull, []byte(originalContent), 0o644); err != nil {
		t.Fatal(err)
	}

	sig := plangen.UpstreamSignatures{
		Spec: plangen.FileSignature{
			Path:   specRel,
			SHA256: plangen.HashNormalized([]byte(originalContent)),
		},
		Plan: plangen.FileSignature{
			Path:   "docs/plans/2026-11-07_PLAN_TEST.md",
			SHA256: strings.Repeat("0", 64),
		},
	}

	// 1. Integridade valida
	if err := plangen.VerifySignatures(base, sig); err != nil {
		t.Fatalf("esperava sucesso na integridade valida, obteve: %v", err)
	}

	// 2. Adulteracao da spec (mutacao de 1 caractere)
	tamperedContent := originalContent + "adulteracao!"
	if err := os.WriteFile(specFull, []byte(tamperedContent), 0o644); err != nil {
		t.Fatal(err)
	}

	err := plangen.VerifySignatures(base, sig)
	if err == nil {
		t.Fatal("esperava erro de integridade violada apos adulteracao, obteve nil")
	}
	if !strings.Contains(err.Error(), "integridade") && !strings.Contains(err.Error(), "mismatch") {
		t.Fatalf("erro inesperado: %v", err)
	}
}

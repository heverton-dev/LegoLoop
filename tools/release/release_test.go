package main

import (
	"path/filepath"
	"testing"
)

func raizRepoTeste(t *testing.T) string {
	t.Helper()
	raiz, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("erro ao resolver raiz do repositorio: %v", err)
	}
	return raiz
}

func TestBuildReproduzivel(t *testing.T) {
	raiz := raizRepoTeste(t)
	dir1 := t.TempDir()
	dir2 := t.TempDir()
	bin1 := filepath.Join(dir1, NomeBinario("linux"))
	bin2 := filepath.Join(dir2, NomeBinario("linux"))

	hash1, err := CompilarBinario(raiz, "linux", "amd64", "v0.1.0-test", bin1)
	if err != nil {
		t.Fatalf("primeira compilacao falhou: %v", err)
	}
	hash2, err := CompilarBinario(raiz, "linux", "amd64", "v0.1.0-test", bin2)
	if err != nil {
		t.Fatalf("segunda compilacao falhou: %v", err)
	}
	if len(hash1) != 64 {
		t.Fatalf("hash SHA-256 invalido: %q", hash1)
	}
	if hash1 != hash2 {
		t.Fatalf("build nao reproduzivel: hash1=%s hash2=%s", hash1, hash2)
	}
}

func TestNomeBinarioPorSO(t *testing.T) {
	if got := NomeBinario("windows"); got != "legoloop.exe" {
		t.Fatalf("esperado legoloop.exe, obtido %q", got)
	}
	if got := NomeBinario("linux"); got != "legoloop" {
		t.Fatalf("esperado legoloop, obtido %q", got)
	}
	if got := NomeBinario("darwin"); got != "legoloop" {
		t.Fatalf("esperado legoloop, obtido %q", got)
	}
}

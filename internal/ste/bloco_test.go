package ste_test

import (
	"bytes"
	"crypto/sha256"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/ste"
)

func hashArquivo(t *testing.T, caminho string) [32]byte {
	t.Helper()
	dados, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatalf("falha ao ler %s: %v", caminho, err)
	}
	return sha256.Sum256(dados)
}

func TestAplicarBloco(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	arq := filepath.Join(dir, "GEMINI.md")
	textoUsuario := "# Minhas Regras\n\nManter foco total.\n\n### Estilo\nFrases curtas.\n"
	if err := os.WriteFile(arq, []byte(textoUsuario), 0644); err != nil {
		t.Fatalf("falha ao criar arquivo inicial: %v", err)
	}

	conteudo := "# STE-PT\nRegra ativa."
	if err := ste.AplicarBloco(arq, "stept", conteudo); err != nil {
		t.Fatalf("primeira chamada falhou: %v", err)
	}

	aposPrimeira, err := os.ReadFile(arq)
	if err != nil {
		t.Fatalf("falha ao ler resultado: %v", err)
	}
	esperadoPrefixo := strings.TrimRight(textoUsuario, "\n")
	if !bytes.HasPrefix(aposPrimeira, []byte(esperadoPrefixo+"\n\n")) {
		t.Fatalf("texto do usuario alterado:\ngot:\n%q\nwant prefix:\n%q", string(aposPrimeira), esperadoPrefixo+"\n\n")
	}
	if !strings.Contains(string(aposPrimeira), conteudo) {
		t.Fatalf("conteudo gerenciado ausente na saida:\n%s", string(aposPrimeira))
	}

	hash1 := hashArquivo(t, arq)
	if err := ste.AplicarBloco(arq, "stept", conteudo); err != nil {
		t.Fatalf("segunda chamada falhou: %v", err)
	}
	hash2 := hashArquivo(t, arq)
	if hash1 != hash2 {
		t.Fatalf("segunda chamada alterou o arquivo: hash1=%x hash2=%x", hash1, hash2)
	}

	novoConteudo := "# STE-PT\nRegra atualizada."
	if err := ste.AplicarBloco(arq, "stept", novoConteudo); err != nil {
		t.Fatalf("terceira chamada falhou: %v", err)
	}
	aposTerceira, err := os.ReadFile(arq)
	if err != nil {
		t.Fatalf("falha ao ler apos atualizacao: %v", err)
	}
	if !bytes.HasPrefix(aposTerceira, []byte(esperadoPrefixo+"\n\n")) {
		t.Fatalf("texto do usuario alterado apos atualizacao")
	}
	if strings.Contains(string(aposTerceira), "Regra ativa.") || !strings.Contains(string(aposTerceira), "Regra atualizada.") {
		t.Fatalf("substituicao do trecho gerenciado falhou:\n%s", string(aposTerceira))
	}
}

func TestAplicarBlocoNovoArquivoEValidacao(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	novo := filepath.Join(dir, "sub", "CLAUDE.md")
	if err := ste.AplicarBloco(novo, "forja", "regra da forja"); err != nil {
		t.Fatalf("falha ao criar arquivo novo: %v", err)
	}
	dados, err := os.ReadFile(novo)
	if err != nil {
		t.Fatalf("falha ao ler novo arquivo: %v", err)
	}
	if !strings.Contains(string(dados), "<!-- forja:inicio") || !strings.Contains(string(dados), "<!-- forja:fim -->") {
		t.Fatalf("marcadores esperados ausentes:\n%s", string(dados))
	}

	if err := ste.AplicarBloco("", "stept", "x"); err == nil {
		t.Fatalf("esperava erro para caminho vazio")
	}
	if err := ste.AplicarBloco(novo, "", "x"); err == nil {
		t.Fatalf("esperava erro para nome vazio")
	}
}

func TestBlocoIgualAoGolden(t *testing.T) {
	t.Parallel()

	entrada, err := os.ReadFile(filepath.Join("testdata", "gemini_entrada.md"))
	if err != nil {
		t.Fatalf("falha ao ler entrada: %v", err)
	}
	regra, err := os.ReadFile(filepath.Join("testdata", "stept_regra.md"))
	if err != nil {
		t.Fatalf("falha ao ler regra: %v", err)
	}
	golden, err := os.ReadFile(filepath.Join("testdata", "gemini_golden.md"))
	if err != nil {
		t.Fatalf("falha ao ler golden: %v", err)
	}

	alvo := filepath.Join(t.TempDir(), "GEMINI.md")
	if err := os.WriteFile(alvo, entrada, 0644); err != nil {
		t.Fatalf("falha ao gravar entrada temporaria: %v", err)
	}

	if err := ste.AplicarBloco(alvo, "stept", string(regra)); err != nil {
		t.Fatalf("AplicarBloco falhou: %v", err)
	}
	obtido, err := os.ReadFile(alvo)
	if err != nil {
		t.Fatalf("falha ao ler saida: %v", err)
	}
	if !bytes.Equal(obtido, golden) {
		t.Fatalf("saida difere do golden:\n--- OBTIDO ---\n%q\n--- GOLDEN ---\n%q", string(obtido), string(golden))
	}

	hash1 := hashArquivo(t, alvo)
	if err := ste.AplicarBloco(alvo, "stept", string(regra)); err != nil {
		t.Fatalf("segunda chamada falhou: %v", err)
	}
	if hash2 := hashArquivo(t, alvo); hash1 != hash2 {
		t.Fatalf("segunda chamada sobre o golden alterou o hash")
	}
}

package forja

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSSOT(t *testing.T) {
	dir := t.TempDir()
	caminhoAgents := filepath.Join(dir, "AGENTS.md")
	arquivos := []string{
		filepath.Join(dir, "GEMINI.md"),
		filepath.Join(dir, "CLAUDE.md"),
	}
	original := "# Regras do usuario\n\nManter esta linha intacta.\n"

	for _, arq := range arquivos {
		if err := os.WriteFile(arq, []byte(original), 0o644); err != nil {
			t.Fatalf("falha ao preparar %s: %v", arq, err)
		}
		alterado, err := AplicarSSOT(arq, caminhoAgents)
		if err != nil {
			t.Fatalf("AplicarSSOT(%s) falhou: %v", arq, err)
		}
		if !alterado {
			t.Fatalf("esperava alterado=true na primeira execucao de %s", arq)
		}
		conteudo, err := os.ReadFile(arq)
		if err != nil {
			t.Fatalf("falha ao ler %s: %v", arq, err)
		}
		texto := string(conteudo)
		if !strings.Contains(texto, "Manter esta linha intacta.") {
			t.Errorf("conteudo original perdido em %s: %s", arq, texto)
		}
		esperadoAgents := filepath.ToSlash(caminhoAgents)
		if !strings.Contains(texto, esperadoAgents) || !strings.Contains(texto, "AGENTS.md") {
			t.Errorf("referencia ao AGENTS.md ausente em %s: %s", arq, texto)
		}
		if !strings.Contains(texto, "<!-- legoloop-ssot:inicio") || !strings.Contains(texto, "<!-- legoloop-ssot:fim -->") {
			t.Errorf("marcadores gerenciados ausentes em %s: %s", arq, texto)
		}

		alterado2, err := AplicarSSOT(arq, caminhoAgents)
		if err != nil {
			t.Fatalf("segunda chamada falhou em %s: %v", arq, err)
		}
		if alterado2 {
			t.Errorf("esperava idempotencia (alterado=false) em %s", arq)
		}
		depois, err := os.ReadFile(arq)
		if err != nil {
			t.Fatalf("falha ao reler %s: %v", arq, err)
		}
		if string(depois) != texto {
			t.Errorf("arquivo mudou na segunda execucao em %s", arq)
		}
	}
}

func TestSSOTNovoArquivoEEntradasInvalidas(t *testing.T) {
	dir := t.TempDir()
	novo := filepath.Join(dir, "sub", "CLAUDE.md")
	alterado, err := AplicarSSOT(novo, "AGENTS.md")
	if err != nil || !alterado {
		t.Fatalf("esperava criacao com sucesso, obteve alterado=%v err=%v", alterado, err)
	}
	if _, err := AplicarSSOT("", "AGENTS.md"); err == nil {
		t.Error("esperava erro para arquivo vazio")
	}
	if _, err := AplicarSSOT(novo, "   "); err == nil {
		t.Error("esperava erro para caminhoAgents vazio")
	}
}

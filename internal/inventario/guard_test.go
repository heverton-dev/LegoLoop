package inventario

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRecusaDelecaoTemplate prova a proteção de lei de templates e vocab.
func TestRecusaDelecaoTemplate(t *testing.T) {
	casos := []struct {
		alvo     string
		recusado bool
	}{
		{"templates/TEMPLATE_X.md", true},
		{"templates", true},
		{"vocab/dicionario.json", true},
		{"docs/templates/TEMPLATE_X.md", true},
		{`docs\vocab\dicionario.json`, true},
		{"reports/relatorio.md", false},
		{"archive/antigo.md", false},
		{"intents/2026-11-08_INTENT_X.md", false},
	}
	for _, c := range casos {
		err := RecusaDelecao(c.alvo)
		if c.recusado {
			if err == nil {
				t.Errorf("RecusaDelecao(%q) aceita deleção em diretório de lei", c.alvo)
				continue
			}
			if code := ExitCode(err); code == 0 {
				t.Errorf("ExitCode(%q) = 0; recusa exige exit != 0", c.alvo)
			}
		} else if err != nil {
			t.Errorf("RecusaDelecao(%q) recusou indevidamente: %v", c.alvo, err)
		}
	}
	if ExitCode(nil) != 0 {
		t.Error("ExitCode(nil) deve ser 0")
	}
}

// TestRegistrarTentativa prova o registro da tentativa no relatório do guard.
func TestRegistrarTentativa(t *testing.T) {
	raiz := t.TempDir()
	if err := RegistrarTentativa(raiz, "templates/TEMPLATE_X.md", ErrDelecaoRecusada); err != nil {
		t.Fatalf("RegistrarTentativa: %v", err)
	}
	registro := filepath.Join(raiz, "docs", "reports", "2026-11-16_GUARD_LOG.md")
	conteudo, err := os.ReadFile(registro)
	if err != nil {
		t.Fatalf("registro ausente: %v", err)
	}
	texto := string(conteudo)
	for _, trecho := range []string{"templates/TEMPLATE_X.md", "RECUSADA", "Lei 6/15"} {
		if !strings.Contains(texto, trecho) {
			t.Errorf("registro sem o trecho %q", trecho)
		}
	}
	if err := RegistrarTentativa(raiz, "reports/relatorio.md", nil); err != nil {
		t.Fatalf("RegistrarTentativa aceita: %v", err)
	}
	conteudo, _ = os.ReadFile(registro)
	if !strings.Contains(string(conteudo), "PERMITIDA") {
		t.Error("registro sem status PERMITIDA")
	}
}

package instalar

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/forja"
)

func TestInstalarEntregaRelatorio(t *testing.T) {
	home := t.TempDir()
	res, err := Executar(context.Background(), Opcao{
		Home:      home,
		Perfil:    "padrao",
		Aplicar:   true,
		Data:      "2026-11-22",
		SO:        "linux",
		Harnesses: []string{"claude"},
		Ambiente:  []string{"API_SECRET=segredo-oculto-123"},
	})
	if err != nil {
		t.Fatalf("Executar(--apply) falhou: %v", err)
	}
	if !res.Aplicado {
		t.Fatal("esperava Aplicado=true")
	}
	if res.Relatorio.Veredito != forja.VereditoOK {
		t.Fatalf("veredito=%q, esperava %q", res.Relatorio.Veredito, forja.VereditoOK)
	}
	verificarRelatoriosGravados(t, home, "2026-11-22", res)
}

func verificarRelatoriosGravados(t *testing.T, home, data string, res Resultado) {
	t.Helper()
	esperadoJSON := filepath.Join(home, ".legoloop", "forja", data, "RELATORIO.json")
	esperadoMD := filepath.Join(home, ".legoloop", "forja", data, "RELATORIO.md")
	if res.CaminhoJSON != esperadoJSON || res.CaminhoMD != esperadoMD {
		t.Fatalf("caminhos inesperados: json=%q md=%q", res.CaminhoJSON, res.CaminhoMD)
	}
	dadosJSON, err := os.ReadFile(esperadoJSON)
	if err != nil {
		t.Fatalf("RELATORIO.json ausente: %v", err)
	}
	var parsed forja.Relatorio
	if err := json.Unmarshal(dadosJSON, &parsed); err != nil {
		t.Fatalf("RELATORIO.json inválido: %v", err)
	}
	if parsed.Veredito != forja.VereditoOK {
		t.Fatalf("veredito no JSON=%q, esperava OK", parsed.Veredito)
	}
	dadosMD, err := os.ReadFile(esperadoMD)
	if err != nil {
		t.Fatalf("RELATORIO.md ausente: %v", err)
	}
	if !strings.Contains(string(dadosMD), "# Relatório da Forja") {
		t.Fatalf("RELATORIO.md sem título esperado: %s", string(dadosMD))
	}
}

func TestInstalarDryRunNaoGravaNada(t *testing.T) {
	home := t.TempDir()
	binFonte := criarBinarioFalso(t, t.TempDir(), "codebase-memory-mcp")
	res, err := Executar(context.Background(), Opcao{
		Home:         home,
		Perfil:       "padrao",
		Aplicar:      false,
		Data:         "2026-11-22",
		SO:           "linux",
		Harnesses:    []string{"claude"},
		BinarioGrafo: binFonte,
	})
	if err != nil {
		t.Fatalf("Executar(--dry-run) falhou: %v", err)
	}
	if res.Aplicado {
		t.Fatal("esperava Aplicado=false no dry-run")
	}
	entradas, err := os.ReadDir(home)
	if err != nil {
		t.Fatalf("falha ao ler home: %v", err)
	}
	if len(entradas) != 0 {
		t.Fatalf("dry-run gravou arquivos no home: %v", entradas)
	}
}

func TestInstalarPerfilInvalido(t *testing.T) {
	_, err := Executar(context.Background(), Opcao{
		Home:   t.TempDir(),
		Perfil: "invalido",
	})
	if err == nil {
		t.Fatal("esperava erro para perfil inválido")
	}
}

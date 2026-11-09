package ste_test

import (
	"strings"
	"testing"

	"legoloop/internal/ste"
)

func TestTextoManual(t *testing.T) {
	t.Parallel()

	manual := ste.TextoManual()
	if !strings.Contains(manual, "MANUAL DE USO DO ALTERNADOR UNIVERSAL STE-PT") {
		t.Errorf("manual nao contem titulo principal:\n%s", manual)
	}
	if !strings.Contains(manual, "┌") || !strings.Contains(manual, "└") {
		t.Errorf("manual deve conter moldura formatada estilo gw:\n%s", manual)
	}
	if !strings.Contains(manual, "hibrido") || !strings.Contains(manual, "caveman") || !strings.Contains(manual, "puro") {
		t.Errorf("manual deve descrever modos disponiveis incluindo puro:\n%s", manual)
	}

	// Alinhamento milimétrico de todas as linhas da moldura
	linhas := strings.Split(strings.TrimRight(manual, "\n"), "\n")
	larguraEsperada := 0
	for i, l := range linhas {
		tamRune := len([]rune(l))
		if i == 0 {
			larguraEsperada = tamRune
		}
		if tamRune != larguraEsperada {
			t.Errorf("linha %d fora de alinhamento (largura %d, esperada %d):\n%q", i+1, tamRune, larguraEsperada, l)
		}
	}
}

func TestRenderTabelaBox(t *testing.T) {
	t.Parallel()

	statusList := []ste.HarnessStatus{
		{
			ID:        ste.HarnessGemini,
			Nome:      "Google Gemini / Antigravity",
			Instalado: true,
			Ativo:     true,
			Perfil:    ste.ConfigPerfil{Modo: ste.ModoHibrido, Pausa: ste.PausaDestrutivo},
			Origem:    "padrao",
			Caminho:   "~/.gemini/config/rules/stept-rules.md",
		},
	}

	tabela := ste.RenderTabelaBox(statusList, "~/.stept/config.json")
	if !strings.Contains(tabela, "┌") || !strings.Contains(tabela, "└") {
		t.Errorf("tabela deve conter moldura estilo gw:\n%s", tabela)
	}
	if !strings.Contains(tabela, "HARNESS") || !strings.Contains(tabela, "Gemini") {
		t.Errorf("tabela deve conter cabecalhos e dados:\n%s", tabela)
	}
}

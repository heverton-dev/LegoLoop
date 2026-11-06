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
	if !strings.Contains(manual, "hibrido") || !strings.Contains(manual, "caveman") {
		t.Errorf("manual deve descrever modos disponiveis:\n%s", manual)
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

package ste_test

import (
	"strings"
	"testing"

	"legoloop/internal/ste"
)

func TestPromptRegra(t *testing.T) {
	t.Parallel()

	t.Run("ModoOff retorna vazio", func(t *testing.T) {
		prompt := ste.PromptRegra(ste.ConfigPerfil{Modo: ste.ModoOff, Pausa: ste.PausaNunca})
		if prompt != "" {
			t.Errorf("esperava string vazia para ModoOff, obteve %d bytes", len(prompt))
		}
	})

	t.Run("ModoCaveman inclui regras de thinking e output telegrafico", func(t *testing.T) {
		prompt := ste.PromptRegra(ste.ConfigPerfil{Modo: ste.ModoCaveman, Pausa: ste.PausaNunca})
		if !strings.Contains(prompt, "CAVEMAN THINKING") && !strings.Contains(prompt, "Caveman Thinking") {
			t.Errorf("esperava mencao a Caveman Thinking no prompt: %s", prompt)
		}
		if !strings.Contains(prompt, "telegráfico") && !strings.Contains(prompt, "telegrafico") {
			t.Errorf("esperava regra telegrafica no prompt caveman: %s", prompt)
		}
	})

	t.Run("ModoHibrido com pausa destrutivo contem triade completa", func(t *testing.T) {
		perfil := ste.ConfigPerfil{
			Modo:     ste.ModoHibrido,
			Pausa:    ste.PausaDestrutivo,
			Thinking: "caveman",
		}
		prompt := ste.PromptRegra(perfil)

		blocosObrigatorios := []string{
			"INPUT IMPERATIVO",
			"CAVEMAN THINKING",
			"TRÍADE STE-PT",
			"[S] Situação",
			"[T] Teoria",
			"[T] Tese",
			"[E] Execução",
			"[E] Evidência",
			"[1] seguir",
			"ações destrutivas",
		}

		for _, b := range blocosObrigatorios {
			if !strings.Contains(strings.ToUpper(prompt), strings.ToUpper(b)) {
				t.Errorf("prompt hibrido nao contem bloco obrigatorio %q", b)
			}
		}
	})

	t.Run("ModoStept com pausa sempre exige cartao para todas as acoes", func(t *testing.T) {
		perfil := ste.ConfigPerfil{
			Modo:  ste.ModoStept,
			Pausa: ste.PausaSempre,
		}
		prompt := ste.PromptRegra(perfil)

		if !strings.Contains(prompt, "TRÍADE STE-PT") {
			t.Errorf("prompt stept deve conter especificacao da triade")
		}
		if !strings.Contains(strings.ToLower(prompt), "todas as ações") && !strings.Contains(strings.ToLower(prompt), "todas as acoes") {
			t.Errorf("prompt stept com pausa sempre deve mencionar todas as acoes: %s", prompt)
		}
		if !strings.Contains(prompt, "1 a 3 linhas") {
			t.Errorf("prompt stept DEVE exigir teto de 1 a 3 linhas")
		}
	})
}

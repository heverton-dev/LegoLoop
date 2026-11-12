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

		secoesObrigatorias := []string{
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

		for _, b := range secoesObrigatorias {
			if !strings.Contains(strings.ToUpper(prompt), strings.ToUpper(b)) {
				t.Errorf("prompt hibrido nao contem secao obrigatoria %q", b)
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

	t.Run("ModoHibrido contem as 7 normas do ASD-STE100", func(t *testing.T) {
		perfil := ste.ConfigPerfil{
			Modo:     ste.ModoHibrido,
			Pausa:    ste.PausaDestrutivo,
			Thinking: "caveman",
		}
		prompt := ste.PromptRegra(perfil)

		regrasASD := []string{
			"ASD-STE100",
			"COMPRIMENTO DE FRASE",
			"UMA AÇÃO POR FRASE",
			"VOZ ATIVA E SUJEITO EXPLÍCITO",
			"ZERO DRIFT DE SINÔNIMOS",
			"ZERO HEDGING",
			"PRESERVAÇÃO DE ARTIGOS",
			"JARGÃO TÉCNICO PERMITIDO",
		}

		for _, r := range regrasASD {
			if !strings.Contains(strings.ToUpper(prompt), strings.ToUpper(r)) {
				t.Errorf("prompt deve conter a norma ASD-STE100 %q", r)
			}
		}
	})

	t.Run("ModoPuro contem ASD-STE100 e nao contem Triade", func(t *testing.T) {
		perfil := ste.ConfigPerfil{
			Modo:  ste.ModoPuro,
			Pausa: ste.PausaDestrutivo,
		}
		prompt := ste.PromptRegra(perfil)

		if !strings.Contains(prompt, "CONTRATO DE RESPOSTA PURO") {
			t.Errorf("prompt puro deve conter CONTRATO DE RESPOSTA PURO")
		}
		if !strings.Contains(prompt, "ASD-STE100") {
			t.Errorf("prompt puro deve conter normas ASD-STE100")
		}
		if strings.Contains(prompt, "TRÍADE STE-PT") {
			t.Errorf("prompt puro NÃO deve conter TRÍADE STE-PT")
		}
	})
}

func TestPromptPuroSemContradicao(t *testing.T) {
	t.Parallel()
	prompt := ste.PromptRegra(ste.ConfigPerfil{Modo: ste.ModoPuro, Pausa: ste.PausaDestrutivo})
	for _, proibido := range []string{"[E] Execução", "PÓS-AÇÃO", "PRÉ-AÇÃO", "Pare em [E]"} {
		if strings.Contains(prompt, proibido) {
			t.Errorf("prompt puro cita %q, que o modo puro proibe", proibido)
		}
	}
}

func TestModoProvaValido(t *testing.T) {
	t.Parallel()
	if err := (ste.ConfigPerfil{Modo: ste.ModoProva, Pausa: ste.PausaNunca}).Validate(); err != nil {
		t.Errorf("modo prova deve ser valido: %v", err)
	}
}

func TestPromptProvaProsa(t *testing.T) {
	t.Parallel()
	prompt := ste.PromptRegra(ste.ConfigPerfil{Modo: ste.ModoProva, Pausa: ste.PausaDestrutivo})
	for _, b := range []string{"veredito", "Recomendação:", "Não verificado", "(prova:", "10 linhas", "100 palavras", "detalhe", "Citar documento não é prova", "Exemplo:", "destrutiva", "[1] seguir"} {
		if !strings.Contains(prompt, b) {
			t.Errorf("prompt prova sem %q", b)
		}
	}
	for _, proibido := range []string{"TRÍADE", "Fatos comprovados", "Stop-Gate Humano"} {
		if strings.Contains(prompt, proibido) {
			t.Errorf("prompt prova nao deve conter %q", proibido)
		}
	}
	if len(prompt) > 1500 {
		t.Errorf("prompt prova deve ter no maximo 1500 bytes, tem %d", len(prompt))
	}
	if !strings.Contains(ste.PromptRegra(ste.ConfigPerfil{Modo: ste.ModoProva, Pausa: ste.PausaSempre}), "toda ação com efeito") {
		t.Errorf("pausa sempre deve alterar a regra de pausa")
	}
	if ste.PerfilPadraoRetorno().Modo != ste.ModoProva {
		t.Errorf("modo prova deve ser o padrao")
	}
}

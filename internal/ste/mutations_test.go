package ste_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/ste"
)

// TestMutationsCoverage é o alvo BDD do plano cobrindo as 4 mutações críticas M1 a M4.
func TestMutationsCoverage(t *testing.T) {
	TestRegra8_MutacoesDeliberadas(t)
}

// TestRegra8_MutacoesDeliberadas comprova matematicamente que qualquer desvio nos 4 pontos críticos faz o teste falhar.
func TestRegra8_MutacoesDeliberadas(t *testing.T) {
	t.Parallel()

	t.Run("M1: Seletor de Heranca nao pode ignorar override", func(t *testing.T) {
		cfg := ste.NovoConfig()
		cfg.Padrao = ste.ConfigPerfil{Modo: ste.ModoHibrido, Pausa: ste.PausaDestrutivo}
		_ = cfg.DefinirHarness(ste.HarnessClaude, ste.ConfigPerfil{Modo: ste.ModoStept, Pausa: ste.PausaSempre})

		// Simulação de mutação M1: forçar retorno do padrão
		mutacaoRetornouPadrao := func() bool {
			perfil, origem := cfg.ResolverPerfil(ste.HarnessClaude)
			return origem == "padrao" || perfil.Modo == ste.ModoHibrido
		}
		if mutacaoRetornouPadrao() {
			t.Errorf("M1 detectado: mutacao permitiu que override do claude fosse ignorado")
		}
	})

	t.Run("M2: Validacao nao pode aceitar modo invalido", func(t *testing.T) {
		mutacaoPerfilInvalido := ste.ConfigPerfil{Modo: "modo_fantasma", Pausa: ste.PausaNunca}
		err := mutacaoPerfilInvalido.Validate()
		if err == nil || !errors.Is(err, ste.ErrInvalido) {
			t.Errorf("M2 detectado: mutacao permitiu modo_fantasma sem retornar ErrInvalido")
		}
	})

	t.Run("M3: Prompt Hibrido nao pode omitir Caveman Thinking ou Cartao ASD", func(t *testing.T) {
		perfil := ste.ConfigPerfil{Modo: ste.ModoHibrido, Pausa: ste.PausaDestrutivo, Thinking: "caveman"}
		prompt := ste.PromptRegra(perfil)

		// Se mutar removendo Caveman Thinking ou Cartão
		if !strings.Contains(prompt, "CAVEMAN THINKING") {
			t.Errorf("M3 detectado: bloco CAVEMAN THINKING ausente")
		}
		if !strings.Contains(prompt, "CARTÃO DE AÇÃO STE-PT") {
			t.Errorf("M3 detectado: bloco CARTÃO DE AÇÃO STE-PT ausente")
		}
	})

	t.Run("M4: Modo Off nao pode deixar stept-rules.md no disco", func(t *testing.T) {
		tmpDir := t.TempDir()
		rulesDir := filepath.Join(tmpDir, ".gemini", "config", "rules")
		_ = os.MkdirAll(rulesDir, 0755)
		steptFile := filepath.Join(rulesDir, "stept-rules.md")
		_ = os.WriteFile(steptFile, []byte("regra antiga"), 0644)

		cfg := ste.NovoConfig()
		cfg.Padrao.Modo = ste.ModoOff
		_, err := ste.Sync(cfg, tmpDir, false)
		if err != nil {
			t.Fatalf("sync off falhou: %v", err)
		}

		if _, err := os.Stat(steptFile); !os.IsNotExist(err) {
			t.Errorf("M4 detectado: modo off nao deletou o arquivo stept-rules.md")
		}
	})
}

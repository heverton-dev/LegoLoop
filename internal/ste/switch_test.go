package ste_test

import (
	"errors"
	"path/filepath"
	"testing"

	"legoloop/internal/ste"
)

func TestConfigPerfil_Validate(t *testing.T) {
	t.Parallel()

	casosValidos := []ste.ConfigPerfil{
		{Modo: ste.ModoOff, Pausa: ste.PausaNunca},
		{Modo: ste.ModoStept, Pausa: ste.PausaDestrutivo},
		{Modo: ste.ModoCaveman, Pausa: ste.PausaSempre},
		{Modo: ste.ModoHibrido, Pausa: ste.PausaDestrutivo, Thinking: "caveman"},
		{Modo: ste.ModoPuro, Pausa: ste.PausaDestrutivo},
	}

	for _, p := range casosValidos {
		if err := p.Validate(); err != nil {
			t.Errorf("esperava valido para %+v, obteve %v", p, err)
		}
	}

	casosInvalidos := []ste.ConfigPerfil{
		{Modo: "invalido", Pausa: ste.PausaNunca},
		{Modo: ste.ModoHibrido, Pausa: "aleatoria"},
		{Modo: "", Pausa: ste.PausaDestrutivo},
		{Modo: ste.ModoStept, Pausa: ""},
	}

	for _, p := range casosInvalidos {
		err := p.Validate()
		if err == nil {
			t.Errorf("esperava erro para %+v, mas retornou nil", p)
		}
		if !errors.Is(err, ste.ErrInvalido) {
			t.Errorf("esperava ste.ErrInvalido para %+v, obteve %v", p, err)
		}
	}
}

func TestConfig_ResolveEHeranca(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.json")

	cfg := ste.NovoConfig()
	cfg.Padrao = ste.ConfigPerfil{
		Modo:     ste.ModoHibrido,
		Pausa:    ste.PausaDestrutivo,
		Thinking: "caveman",
	}

	overrideClaude := ste.ConfigPerfil{
		Modo:  ste.ModoStept,
		Pausa: ste.PausaSempre,
	}
	if err := cfg.DefinirHarness(ste.HarnessClaude, overrideClaude); err != nil {
		t.Fatalf("falha ao definir harness claude: %v", err)
	}

	// Salva e recarrega para testar persistência hermética
	if err := ste.SalvarConfig(cfgPath, cfg); err != nil {
		t.Fatalf("falha ao salvar config: %v", err)
	}

	carregado, err := ste.CarregarConfig(cfgPath)
	if err != nil {
		t.Fatalf("falha ao carregar config: %v", err)
	}

	// Resolver Gemini (sem override -> herda padrao)
	perfilGemini, origemGemini := carregado.ResolverPerfil(ste.HarnessGemini)
	if origemGemini != "padrao" {
		t.Errorf("esperava origem 'padrao' para gemini, obteve %q", origemGemini)
	}
	if perfilGemini.Modo != ste.ModoHibrido || perfilGemini.Pausa != ste.PausaDestrutivo {
		t.Errorf("perfil gemini incorreto: %+v", perfilGemini)
	}

	// Resolver Claude (com override)
	perfilClaude, origemClaude := carregado.ResolverPerfil(ste.HarnessClaude)
	if origemClaude != "override" {
		t.Errorf("esperava origem 'override' para claude, obteve %q", origemClaude)
	}
	if perfilClaude.Modo != ste.ModoStept || perfilClaude.Pausa != ste.PausaSempre {
		t.Errorf("perfil claude incorreto: %+v", perfilClaude)
	}

	// Testar fallback quando arquivo não existe
	naoExistePath := filepath.Join(tmpDir, "inexistente.json")
	padraoFallback, err := ste.CarregarConfig(naoExistePath)
	if err != nil {
		t.Fatalf("esperava retorno de config padrão quando arquivo não existe, obteve erro: %v", err)
	}
	if padraoFallback.Padrao.Modo != ste.ModoProva {
		t.Errorf("esperava modo padrão prova no fallback, obteve %q", padraoFallback.Padrao.Modo)
	}
}

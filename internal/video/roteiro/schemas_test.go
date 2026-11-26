package roteiro

import (
	"errors"
	"testing"
)

func projetoValidoTeste() Projeto {
	return Projeto{
		SchemaVersion:      1,
		Slug:               "pix-em-segundos",
		Canal:              "bilhoes-em-jogo",
		Nicho:              "APIs open source",
		PerfisFonte:        []string{"tecnologia", "geral"},
		EstruturaNarrativa: "explicador",
		Formatos:           []string{"16:9", "9:16"},
		DuracaoAlvoS:       60,
		Voz:                "pt-BR-ThalitaMultilingualNeural",
		EstiloVisual:       "vetor_2d_diorama",
		Animacao:           "codigo",
		PoliticaGates: PoliticaGates{
			Modo:         "hibrido",
			GatesHumanos: []string{"etapa1_pauta"},
		},
	}
}

func marcaValidaTeste() Marca {
	m := Marca{
		SchemaVersion: 1,
		Nome:          "Bilhoes em Jogo",
		Formatos:      []string{"16:9", "9:16", "1:1"},
	}
	m.Voz.Motor = "edge-tts"
	m.Voz.Voz = "pt-BR-ThalitaMultilingualNeural"
	m.Voz.Velocidade = "+0%"
	return m
}

func TestValidarConfiguracao(t *testing.T) {
	cat := CatalogoPadrao()

	t.Run("configuracao valida passa", func(t *testing.T) {
		if err := ValidarConfiguracao(projetoValidoTeste(), marcaValidaTeste(), cat); err != nil {
			t.Fatalf("esperava nil, obteve erro: %v", err)
		}
	})

	t.Run("nicho vazio bloqueia com ErrSchemaInvalido", func(t *testing.T) {
		p := projetoValidoTeste()
		p.Nicho = "   "
		err := ValidarConfiguracao(p, marcaValidaTeste(), cat)
		if !errors.Is(err, ErrSchemaInvalido) {
			t.Fatalf("esperava ErrSchemaInvalido, obteve: %v", err)
		}
	})

	t.Run("perfis_fonte vazio ou desconhecido bloqueia", func(t *testing.T) {
		p := projetoValidoTeste()
		p.PerfisFonte = nil
		if err := ValidarConfiguracao(p, marcaValidaTeste(), cat); !errors.Is(err, ErrSchemaInvalido) {
			t.Fatalf("esperava ErrSchemaInvalido para perfis vazios, obteve: %v", err)
		}
		p.PerfisFonte = []string{"perfil_inexistente"}
		if err := ValidarConfiguracao(p, marcaValidaTeste(), cat); !errors.Is(err, ErrSchemaInvalido) {
			t.Fatalf("esperava ErrSchemaInvalido para perfil desconhecido, obteve: %v", err)
		}
	})

	t.Run("politica_gates invalida ou hibrida sem gates_humanos bloqueia", func(t *testing.T) {
		p := projetoValidoTeste()
		p.PoliticaGates.Modo = "automatico"
		if err := ValidarConfiguracao(p, marcaValidaTeste(), cat); !errors.Is(err, ErrSchemaInvalido) {
			t.Fatalf("esperava ErrSchemaInvalido para modo invalido, obteve: %v", err)
		}
		p.PoliticaGates.Modo = "hibrido"
		p.PoliticaGates.GatesHumanos = nil
		if err := ValidarConfiguracao(p, marcaValidaTeste(), cat); !errors.Is(err, ErrSchemaInvalido) {
			t.Fatalf("esperava ErrSchemaInvalido para hibrido sem gates_humanos, obteve: %v", err)
		}
	})

	t.Run("inconsistencia de formato ou voz entre projeto e marca bloqueia", func(t *testing.T) {
		p := projetoValidoTeste()
		m := marcaValidaTeste()
		m.Formatos = []string{"16:9"}
		if err := ValidarConfiguracao(p, m, cat); !errors.Is(err, ErrConsistenciaConfig) {
			t.Fatalf("esperava ErrConsistenciaConfig para formato incompativel, obteve: %v", err)
		}
		m = marcaValidaTeste()
		m.Voz.Voz = "pt-BR-AntonioNeural"
		if err := ValidarConfiguracao(p, m, cat); !errors.Is(err, ErrConsistenciaConfig) {
			t.Fatalf("esperava ErrConsistenciaConfig para voz incompativel, obteve: %v", err)
		}
	})
}

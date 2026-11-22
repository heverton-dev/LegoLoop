package forja

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestContrato(t *testing.T) {
	t.Run("validate rejeita home vazio", func(t *testing.T) {
		op := Opcao{Home: "   ", Perfil: "padrao"}
		if err := op.Validate(); !errors.Is(err, ErrSemHome) {
			t.Fatalf("esperava ErrSemHome, obteve %v", err)
		}
	})

	t.Run("validate rejeita perfil desconhecido", func(t *testing.T) {
		op := Opcao{Home: t.TempDir(), Perfil: "xyz"}
		if err := op.Validate(); !errors.Is(err, ErrPerfilDesconhecido) {
			t.Fatalf("esperava ErrPerfilDesconhecido, obteve %v", err)
		}
	})

	t.Run("validate aceita perfis canonicos", func(t *testing.T) {
		home := t.TempDir()
		for _, p := range []string{"minimo", "padrao", "gateway"} {
			op := Opcao{Home: home, Perfil: p}
			if err := op.Validate(); err != nil {
				t.Fatalf("perfil %q deveria ser válido, erro: %v", p, err)
			}
		}
	})

	t.Run("serializacao JSON preserva contrato da SPEC", func(t *testing.T) {
		op := Opcao{
			Home:      "/tmp/home",
			Perfil:    "padrao",
			Sem:       []string{"opencode"},
			Aplicar:   true,
			Desfazer:  "2026-11-22",
			Data:      "2026-11-22",
			SO:        "linux",
			Harnesses: []string{"claude"},
		}
		b, err := json.Marshal(op)
		if err != nil {
			t.Fatalf("falha ao serializar Opcao: %v", err)
		}
		var m map[string]any
		if err := json.Unmarshal(b, &m); err != nil {
			t.Fatalf("falha ao desserializar Opcao: %v", err)
		}
		for _, campo := range []string{"home", "perfil", "sem", "aplicar", "desfazer"} {
			if _, ok := m[campo]; !ok {
				t.Errorf("campo JSON obrigatório ausente em Opcao: %s", campo)
			}
		}
		for _, interno := range []string{"Data", "SO", "Harnesses", "Geradores"} {
			if _, ok := m[interno]; ok {
				t.Errorf("campo interno não deveria ser serializado em JSON: %s", interno)
			}
		}
	})

	t.Run("relatorio e etapa serializam campos do contrato", func(t *testing.T) {
		rel := Relatorio{
			Versao:   VersaoContrato,
			Data:     "2026-11-22",
			SO:       "windows",
			Perfil:   "padrao",
			Veredito: VereditoOK,
			Etapas: []EtapaResultado{
				{
					Nome:      "forjar:claude",
					Harness:   "claude",
					Estado:    EstadoOK,
					Alterados: []string{"a.cfg"},
					Backup:    "backups/forja_2026-11-22/claude",
					Mensagem:  "aplicado",
				},
			},
			TokensAntes:  map[string]int{"claude": 1200},
			TokensDepois: map[string]int{"claude": 300},
		}
		b, err := json.Marshal(rel)
		if err != nil {
			t.Fatalf("falha ao serializar Relatorio: %v", err)
		}
		var out Relatorio
		if err := json.Unmarshal(b, &out); err != nil {
			t.Fatalf("falha ao desserializar Relatorio: %v", err)
		}
		if out.Veredito != VereditoOK || len(out.Etapas) != 1 || out.Etapas[0].Harness != "claude" {
			t.Fatalf("round-trip de Relatorio divergente: %+v", out)
		}
	})
}

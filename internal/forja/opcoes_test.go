package forja

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestSemEPerfil(t *testing.T) {
	t.Run("validar perfil rejeita desconhecido", func(t *testing.T) {
		if err := ValidarPerfil("xyz"); !errors.Is(err, ErrPerfilDesconhecido) {
			t.Fatalf("esperava ErrPerfilDesconhecido, obteve %v", err)
		}
		for _, p := range PerfisValidos() {
			if err := ValidarPerfil(p); err != nil {
				t.Fatalf("perfil %q deveria ser válido: %v", p, err)
			}
		}
	})

	t.Run("filtrar harnesses separa ativos e ignorados sem case sensitivity", func(t *testing.T) {
		ativos, ignorados := FiltrarHarnesses(
			[]string{"claude", "opencode", "gemini"},
			[]string{" OpenCode "},
		)
		if len(ativos) != 2 || ativos[0] != "claude" || ativos[1] != "gemini" {
			t.Fatalf("ativos inesperados: %v", ativos)
		}
		if len(ignorados) != 1 || ignorados[0] != "opencode" {
			t.Fatalf("ignorados inesperados: %v", ignorados)
		}
	})

	t.Run("sem opencode marca ignorada e nao toca nos arquivos dele", func(t *testing.T) {
		home := t.TempDir()
		pastaOpen := filepath.Join(home, ".opencode")
		if err := os.MkdirAll(pastaOpen, 0o755); err != nil {
			t.Fatal(err)
		}
		arqOpen := filepath.Join(pastaOpen, "config.json")
		original := []byte(`{"usuario":"intacto"}`)
		if err := os.WriteFile(arqOpen, original, 0o644); err != nil {
			t.Fatal(err)
		}

		op := Opcao{
			Home:      home,
			Perfil:    "padrao",
			Sem:       []string{"opencode"},
			Aplicar:   true,
			Data:      "2026-11-22",
			SO:        "linux",
			Harnesses: []string{"claude", "opencode"},
		}
		rel, err := Executar(context.Background(), op)
		if err != nil {
			t.Fatalf("Executar falhou: %v", err)
		}
		encontrouIgnorada := false
		for _, et := range rel.Etapas {
			if et.Harness == "opencode" {
				encontrouIgnorada = true
				if et.Estado != EstadoIgnorada {
					t.Fatalf("esperava estado %q para opencode, obteve %q", EstadoIgnorada, et.Estado)
				}
				if len(et.Alterados) != 0 {
					t.Fatalf("opencode ignorado não deve ter alterados: %v", et.Alterados)
				}
			}
		}
		if !encontrouIgnorada {
			t.Fatalf("etapa do opencode não encontrada no relatório: %+v", rel.Etapas)
		}
		depois, err := os.ReadFile(arqOpen)
		if err != nil {
			t.Fatal(err)
		}
		if string(depois) != string(original) {
			t.Fatalf("arquivo do opencode foi alterado: %q", string(depois))
		}
	})

	t.Run("executar com perfil invalido devolve ErrPerfilDesconhecido", func(t *testing.T) {
		_, err := Executar(context.Background(), Opcao{
			Home:      t.TempDir(),
			Perfil:    "xyz",
			Harnesses: []string{"claude"},
		})
		if !errors.Is(err, ErrPerfilDesconhecido) {
			t.Fatalf("esperava ErrPerfilDesconhecido, obteve %v", err)
		}
	})
}

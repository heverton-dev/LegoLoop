package taxonomy

import (
	"errors"
	"testing"
)

func TestValidarPapel(t *testing.T) {
	t.Run("papel homologado retorna nil e pertence a guilda construcao", func(t *testing.T) {
		if err := ValidarPapel("engenheiro-frontend"); err != nil {
			t.Fatalf("esperado nil, obtido %v", err)
		}
		papel, err := ObterPapel("engenheiro-frontend")
		if err != nil {
			t.Fatalf("esperado nil ao obter papel, obtido %v", err)
		}
		if papel.Guilda != GuildaConstrucao {
			t.Fatalf("esperado guilda %q, obtido %q", GuildaConstrucao, papel.Guilda)
		}
	})

	t.Run("papel desconhecido retorna ErrPapelInvalido", func(t *testing.T) {
		err := ValidarPapel("ninja-react")
		if !errors.Is(err, ErrPapelInvalido) {
			t.Fatalf("esperado ErrPapelInvalido, obtido %v", err)
		}
	})

	t.Run("papel vazio ou com espacos falha com ErrPapelInvalido", func(t *testing.T) {
		if err := ValidarPapel("   "); !errors.Is(err, ErrPapelInvalido) {
			t.Fatalf("esperado ErrPapelInvalido para vazio, obtido %v", err)
		}
		if err := ValidarPapel(" engenheiro-tdd "); err != nil {
			t.Fatalf("esperado nil com trim de espacos, obtido %v", err)
		}
	})
}

func TestListarGuildasECatalogo(t *testing.T) {
	guildas := ListarGuildas()
	if len(guildas) != 7 {
		t.Fatalf("esperado 7 guildas, obtido %d", len(guildas))
	}
	for _, g := range guildas {
		if err := ValidarGuilda(g); err != nil {
			t.Fatalf("guilda %q deveria ser valida: %v", g, err)
		}
	}
	if err := ValidarGuilda(Guilda("inexistente")); !errors.Is(err, ErrGuildaInvalida) {
		t.Fatalf("esperado ErrGuildaInvalida, obtido %v", err)
	}
	if len(CatalogoCanonico) != 20 {
		t.Fatalf("esperado 20 papeis no catalogo da spec, obtido %d", len(CatalogoCanonico))
	}
	for slug, p := range CatalogoCanonico {
		if p.Slug != slug {
			t.Fatalf("chave %q diverge de Slug %q", slug, p.Slug)
		}
		if err := ValidarGuilda(p.Guilda); err != nil {
			t.Fatalf("papel %q possui guilda invalida %q", slug, p.Guilda)
		}
	}
}

func TestSincronizacaoDicionarioJSON(t *testing.T) {
	dict, err := ObterDicionario()
	if err != nil {
		t.Fatalf("falha ao carregar dicionario: %v", err)
	}
	if len(dict.Guildas) != 7 {
		t.Fatalf("esperado 7 guildas no dicionario.json, obtido %d", len(dict.Guildas))
	}
	if len(dict.Papeis) != 20 {
		t.Fatalf("esperado 20 papeis no dicionario.json, obtido %d", len(dict.Papeis))
	}
	for _, p := range dict.Papeis {
		if err := ValidarPapel(p.Slug); err != nil {
			t.Fatalf("papel %q no dicionario.json nao e valido no codigo: %v", p.Slug, err)
		}
		if err := ValidarGuilda(Guilda(p.Guilda)); err != nil {
			t.Fatalf("guilda %q no dicionario.json nao e valida no codigo: %v", p.Guilda, err)
		}
	}
}

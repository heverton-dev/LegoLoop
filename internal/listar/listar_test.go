package listar_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"legoloop/internal/listar"
	"legoloop/internal/processo"
)

func etapa(id string, depende ...string) processo.Etapa {
	return processo.Etapa{
		ID: id, Titulo: "Etapa " + id, Tipo: "construir", Peca: "tdd-workflow", Depende: depende,
		Produz: processo.Artefato{Caminho: "out/" + id + ".md", Tipo: "markdown"},
		Gate:   "file_exists('out/" + id + ".md')",
	}
}

func proc(nome, versao, natureza string, es ...processo.Etapa) processo.Processo {
	return processo.Processo{
		Nome: nome, Versao: versao, Natureza: natureza, Objetivo: "Objetivo de " + nome + ".",
		Etapas: es, CriadoEm: time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC),
	}
}

func gravar(t *testing.T, vault string, p processo.Processo) {
	t.Helper()
	data, err := p.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	escrever(t, processo.CaminhoNoCofre(vault, p.Nome), string(data))
}

func escrever(t *testing.T, path, data string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

// cofreDeTres monta o cofre do AC-E03.1: legoloop-a vinculado a 2 harnesses, legoloop-b que compõe a, e
// legoloop-velho gravado pela v1 (só o .loop.json). Uma esteira sem prefixo legoloop- não é processo.
func cofreDeTres(t *testing.T) (vault, home string) {
	t.Helper()
	root := t.TempDir()
	vault, home = filepath.Join(root, "vault"), filepath.Join(root, "home")
	gravar(t, vault, proc("legoloop-a", "1.2.0", "dados", etapa("E01")))
	sub := processo.Etapa{ID: "E02", Titulo: "Sub", Tipo: "construir", Compoe: "legoloop-a@1.2.0", Depende: []string{"E01"},
		Produz: processo.Artefato{Caminho: "out/E02.md", Tipo: "markdown"}, Gate: "file_exists('out/E02.md')"}
	gravar(t, vault, proc("legoloop-b", "1.0.0", "software", etapa("E01"), sub))
	escrever(t, filepath.Join(vault, "loops", "legoloop-velho.loop.json"), `{"id":"legoloop-velho","name":"legoloop-velho","version":"0.3.0"}`)
	escrever(t, filepath.Join(vault, "loops", "legoloop-a.loop.json"), `{"id":"legoloop-a","version":"1.2.0"}`)
	escrever(t, filepath.Join(vault, "loops", "nova-esteira.loop.json"), `{"id":"nova-esteira","version":"2.0.0"}`)
	if err := os.MkdirAll(filepath.Join(vault, "loops", "legoloop-pasta.loop.json"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, h := range []string{"claude", "codex", "gemini"} {
		if err := os.MkdirAll(filepath.Join(home, "."+h, "skills"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	for _, h := range []string{"claude", "codex"} {
		escrever(t, filepath.Join(home, "."+h, "skills", "legoloop-a", "SKILL.md"), "# a\n")
	}
	return vault, home
}

func TestListar(t *testing.T) {
	t.Run("três processos", func(t *testing.T) {
		vault, home := cofreDeTres(t)
		res, err := listar.Run(context.Background(), listar.Deps{Store: processo.Store{Roots: processo.Roots{Vault: vault}}, Home: home})
		if err != nil {
			t.Fatal(err)
		}
		quer := []listar.Item{
			{Nome: "legoloop-a", Versao: "1.2.0", Natureza: "dados", Harnesses: []string{"claude", "codex"}, UsadoPor: []string{"legoloop-b"}, Origem: "v2"},
			{Nome: "legoloop-b", Versao: "1.0.0", Natureza: "software", Harnesses: []string{}, UsadoPor: []string{}, Origem: "v2"},
			{Nome: "legoloop-velho", Versao: "0.3.0", Natureza: listar.NaturezaDesconhecida, Harnesses: []string{}, UsadoPor: []string{}, Origem: "v1"},
		}
		if len(res.Processos) != len(quer) {
			t.Fatalf("processos: %+v", res.Processos)
		}
		for i, q := range quer {
			g := res.Processos[i]
			if g.Nome != q.Nome || g.Versao != q.Versao || g.Natureza != q.Natureza || g.Origem != q.Origem ||
				!slices.Equal(g.Harnesses, q.Harnesses) || !slices.Equal(g.UsadoPor, q.UsadoPor) || g.Erro != "" {
				t.Errorf("processo %d:\n got %+v\nquer %+v", i, g, q)
			}
		}
		if res.Vault != vault {
			t.Errorf("vault %q", res.Vault)
		}
	})
	t.Run("cofre vazio", func(t *testing.T) {
		vault := t.TempDir()
		res, err := listar.Run(context.Background(), listar.Deps{Store: processo.Store{Roots: processo.Roots{Vault: vault}}, Home: t.TempDir()})
		if err != nil || len(res.Processos) != 0 || res.Processos == nil {
			t.Fatalf("cofre vazio: %+v, %v", res, err)
		}
	})
	t.Run("cofre inexistente", func(t *testing.T) {
		vault := filepath.Join(t.TempDir(), "nao-existe")
		_, err := listar.Run(context.Background(), listar.Deps{Store: processo.Store{Roots: processo.Roots{Vault: vault}}})
		if !errors.Is(err, listar.ErrSemCofre) || !strings.Contains(err.Error(), vault) {
			t.Fatalf("esperava ErrSemCofre com o caminho: %v", err)
		}
	})
	t.Run("processo inválido não derruba a lista", func(t *testing.T) {
		vault, home := cofreDeTres(t)
		escrever(t, processo.CaminhoNoCofre(vault, "legoloop-quebrado"), "{ não é json")
		res, err := listar.Run(context.Background(), listar.Deps{Store: processo.Store{Roots: processo.Roots{Vault: vault}}, Home: home})
		if err != nil {
			t.Fatal(err)
		}
		i := slices.IndexFunc(res.Processos, func(it listar.Item) bool { return it.Nome == "legoloop-quebrado" })
		if len(res.Processos) != 4 || i < 0 || res.Processos[i].Erro == "" || res.Processos[i].Origem != "v2" {
			t.Fatalf("o inválido aparece com erro e os outros seguem: %+v", res.Processos)
		}
	})
	t.Run("contexto cancelado", func(t *testing.T) {
		vault, _ := cofreDeTres(t)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if _, err := listar.Run(ctx, listar.Deps{Store: processo.Store{Roots: processo.Roots{Vault: vault}}}); !errors.Is(err, context.Canceled) {
			t.Fatalf("esperava context.Canceled: %v", err)
		}
	})
}

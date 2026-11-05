package listar_test

import (
	"context"
	"encoding/json"
	"path/filepath"
	"slices"
	"testing"

	"legoloop/internal/listar"
	"legoloop/internal/matriz"
	"legoloop/internal/processo"
)

// comLacuna é um processo com a etapa coletar em lacuna (casos dados) e a etapa construir com peça.
func comLacuna(casos ...string) processo.Processo {
	col := etapa("E01")
	col.Tipo, col.Peca, col.Casos = "coletar", "", casos
	return proc("legoloop-coleta", "1.0.0", "dados", col, etapa("E02", "E01"))
}

func entrada(ref, nome, hash string, p *matriz.PerfilDeUso) matriz.Entrada {
	return matriz.Entrada{Ref: ref, Nome: nome, Tipo: "skill", Hash: hash, Perfil: p}
}

func gravarMatriz(t *testing.T, vault string, es ...matriz.Entrada) {
	t.Helper()
	b, err := json.Marshal(matriz.Registro{Versao: matriz.VersaoMatriz, Entradas: es})
	if err != nil {
		t.Fatal(err)
	}
	escrever(t, filepath.Join(vault, "registry_matrix.json"), string(b))
}

// matrizDeColeta: a e b fazem coletar sem declarar coleta-de-dados; c já declara; d faz outra etapa; e não tem perfil.
func matrizDeColeta(t *testing.T, vault string) {
	gravarMatriz(t, vault,
		entrada("~/skills/peca-b/SKILL.md", "peca-b", "hb", &matriz.PerfilDeUso{Etapas: []string{"coletar"}}),
		entrada("~/skills/peca-a/SKILL.md", "peca-a", "ha", &matriz.PerfilDeUso{Etapas: []string{"coletar"}, Naturezas: []string{"dados"}, CasosDeUso: []string{"pesquisa"}}),
		entrada("D:/copia/peca-a/SKILL.md", "peca-a", "ha", &matriz.PerfilDeUso{Etapas: []string{"coletar"}, Naturezas: []string{"dados"}, CasosDeUso: []string{"pesquisa"}}),
		entrada("~/skills/peca-c/SKILL.md", "peca-c", "hc", &matriz.PerfilDeUso{Etapas: []string{"coletar"}, CasosDeUso: []string{"coleta-de-dados"}}),
		entrada("~/skills/peca-d/SKILL.md", "peca-d", "hd", &matriz.PerfilDeUso{Etapas: []string{"construir"}}),
		entrada("~/skills/peca-e/SKILL.md", "peca-e", "he", nil),
	)
}

func listarUm(t *testing.T, vault string) listar.Item {
	t.Helper()
	res, err := listar.Run(context.Background(), listar.Deps{Store: processo.Store{Roots: processo.Roots{Vault: vault}}})
	if err != nil {
		t.Fatal(err)
	}
	i := slices.IndexFunc(res.Processos, func(it listar.Item) bool { return it.Nome == "legoloop-coleta" })
	if i < 0 {
		t.Fatalf("processo com lacuna não listado: %+v", res.Processos)
	}
	return res.Processos[i]
}

func TestSugestaoPerfilDefinir(t *testing.T) {
	t.Run("duas peças candidatas", func(t *testing.T) {
		vault := t.TempDir()
		gravar(t, vault, comLacuna("coleta-de-dados"))
		matrizDeColeta(t, vault)
		it := listarUm(t, vault)
		if !slices.Equal(it.Lacunas, []string{"coletar"}) || len(it.Sugestoes) != 1 {
			t.Fatalf("lacunas %v, sugestões %+v", it.Lacunas, it.Sugestoes)
		}
		s := it.Sugestoes[0]
		quer := []string{
			"legoloop matriz perfil definir peca-a --caso coleta-de-dados --caso pesquisa --etapa coletar --natureza dados",
			"legoloop matriz perfil definir peca-b --caso coleta-de-dados --etapa coletar",
		}
		if s.Etapa != "coletar" || s.Caso != "coleta-de-dados" || !slices.Equal(s.Pecas, []string{"peca-a", "peca-b"}) || !slices.Equal(s.Comandos, quer) {
			t.Errorf("sugestão:\n%+v\nquer comandos %q", s, quer)
		}
	})
	t.Run("sem peça candidata", func(t *testing.T) {
		vault := t.TempDir()
		gravar(t, vault, comLacuna("coleta-de-dados"))
		gravarMatriz(t, vault, entrada("~/skills/peca-d/SKILL.md", "peca-d", "hd", &matriz.PerfilDeUso{Etapas: []string{"construir"}}))
		it := listarUm(t, vault)
		if !slices.Equal(it.Lacunas, []string{"coletar"}) || len(it.Sugestoes) != 0 {
			t.Fatalf("só a lacuna, sem sugestão: %v %+v", it.Lacunas, it.Sugestoes)
		}
	})
	t.Run("sem matriz", func(t *testing.T) {
		vault := t.TempDir()
		gravar(t, vault, comLacuna("coleta-de-dados"))
		it := listarUm(t, vault)
		if !slices.Equal(it.Lacunas, []string{"coletar"}) || len(it.Sugestoes) != 0 {
			t.Fatalf("sem matriz: lacuna sem sugestão: %v %+v", it.Lacunas, it.Sugestoes)
		}
	})
	t.Run("etapa sem casos usa o vocabulário", func(t *testing.T) {
		vault := t.TempDir()
		gravar(t, vault, comLacuna())
		matrizDeColeta(t, vault)
		it := listarUm(t, vault)
		var casos []string
		for _, s := range it.Sugestoes {
			casos = append(casos, s.Caso)
		}
		if !slices.Equal(casos, []string{"pesquisa", "coleta-de-dados"}) {
			t.Fatalf("casos do vocabulário para coletar: %v", casos)
		}
		if !slices.Equal(it.Sugestoes[0].Pecas, []string{"peca-b", "peca-c"}) {
			t.Errorf("pesquisa: peca-a já declara; sobram b e c: %v", it.Sugestoes[0].Pecas)
		}
	})
	t.Run("nome ambíguo usa o caminho", func(t *testing.T) {
		vault := t.TempDir()
		gravar(t, vault, comLacuna("coleta-de-dados"))
		gravarMatriz(t, vault,
			entrada("~/skills/x/SKILL.md", "dup", "h1", &matriz.PerfilDeUso{Etapas: []string{"coletar"}}),
			entrada("D:/meus docs/x/SKILL.md", "dup", "h2", &matriz.PerfilDeUso{Etapas: []string{"coletar"}}),
		)
		s := listarUm(t, vault).Sugestoes
		quer := []string{
			`legoloop matriz perfil definir "D:/meus docs/x/SKILL.md" --caso coleta-de-dados --etapa coletar`,
			"legoloop matriz perfil definir ~/skills/x/SKILL.md --caso coleta-de-dados --etapa coletar",
		}
		if len(s) != 1 || !slices.Equal(s[0].Comandos, quer) {
			t.Fatalf("nomes iguais com conteúdos diferentes vão pelo caminho (com aspas se preciso): %+v", s)
		}
	})
	t.Run("lacuna repetida e etapa composta", func(t *testing.T) {
		vault := t.TempDir()
		p := comLacuna("coleta-de-dados")
		outra := p.Etapas[0]
		outra.ID, outra.Depende = "E03", []string{"E02"}
		outra.Produz.Caminho, outra.Gate = "out/E03.md", "file_exists('out/E03.md')"
		sub := etapa("E04", "E03")
		sub.Peca, sub.Compoe = "", "legoloop-outro@1.0.0"
		p.Etapas = append(p.Etapas, outra, sub)
		gravar(t, vault, p)
		matrizDeColeta(t, vault)
		it := listarUm(t, vault)
		if !slices.Equal(it.Lacunas, []string{"coletar"}) || len(it.Sugestoes) != 1 {
			t.Fatalf("a etapa repetida conta uma vez e a composta não é lacuna: %v %+v", it.Lacunas, it.Sugestoes)
		}
	})
	t.Run("vocabulário do cofre", func(t *testing.T) {
		vault := t.TempDir()
		gravar(t, vault, comLacuna())
		matrizDeColeta(t, vault)
		escrever(t, filepath.Join(vault, "vocab.json"), `{"versao": 2, "etapas": [{"nome": "coletar"}], "naturezas": [],
			"casos_de_uso": [{"nome": "so-do-cofre", "etapas": ["coletar"]}]}`)
		it := listarUm(t, vault)
		if len(it.Sugestoes) != 1 || it.Sugestoes[0].Caso != "so-do-cofre" {
			t.Fatalf("os casos vêm do vocabulário do cofre: %+v", it.Sugestoes)
		}
	})
	t.Run("processo sem lacuna", func(t *testing.T) {
		vault := t.TempDir()
		p := comLacuna("coleta-de-dados")
		p.Etapas[0].Peca = "peca-c"
		gravar(t, vault, p)
		matrizDeColeta(t, vault)
		if it := listarUm(t, vault); len(it.Lacunas) != 0 || len(it.Sugestoes) != 0 {
			t.Fatalf("sem lacuna, sem sugestão: %+v", it)
		}
	})
}

package ativos

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/video/roteiro"
)

func TestResolverElencoPadraoCasos(t *testing.T) {
	proj := roteiro.Projeto{Slug: "x", EstiloVisual: "hibrido"}

	explicito := []FichaPersonagem{{ID: "a"}}
	if got := resolverElencoPadrao(explicito, proj, nil); len(got) != 1 || got[0].ID != "a" {
		t.Fatalf("elenco explicito deveria vencer: %v", got)
	}

	fichas := []FichaAtivo{{ID: "heroi", Nome: "Heroi", Categoria: "personagem"}, {ID: "fundo1", Categoria: "cenario"}}
	got := resolverElencoPadrao(nil, proj, fichas)
	if len(got) != 1 || got[0].ID != "heroi" {
		t.Fatalf("deveria converter so a ficha de personagem: %v", got)
	}

	got2 := resolverElencoPadrao(nil, proj, nil)
	if len(got2) != 1 || got2[0].ID != "narrador" {
		t.Fatalf("sem elenco e sem fichas deveria usar o narrador padrao: %v", got2)
	}
}

func TestConverterFichaEmPersonagemComFolhaExplicita(t *testing.T) {
	f := FichaAtivo{ID: "x", Nome: "X", DescricaoImutavel: "d", FolhaMestrePNG: "ja/pronta.png"}
	p := converterFichaEmPersonagem(f)
	if p.Fichas["turnaround"] != "ja/pronta.png" {
		t.Fatalf("deveria reaproveitar FolhaMestrePNG: %+v", p)
	}
}

func TestLocalizarInventarioVariasOrigens(t *testing.T) {
	if got := localizarInventario(".", "explicito/dir"); got != "explicito/dir" {
		t.Fatalf("explicito deveria vencer: %q", got)
	}
	t.Setenv("LEGOLOOP_INVENTARIO_DIR", "do/ambiente")
	if got := localizarInventario(".", ""); got != "do/ambiente" {
		t.Fatalf("variavel de ambiente deveria vencer: %q", got)
	}
	t.Setenv("LEGOLOOP_INVENTARIO_DIR", "")
	raiz := t.TempDir()
	if err := os.MkdirAll(filepath.Join(raiz, "inventario", "video", "estilos"), 0o755); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(raiz, "projetos", "demo")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	got := localizarInventario(sub, "")
	if filepath.Clean(got) != filepath.Clean(filepath.Join(raiz, "inventario", "video")) {
		t.Fatalf("subindo pastas deveria achar %q, veio %q", filepath.Join(raiz, "inventario", "video"), got)
	}
}

func TestGravarJSONSeguroErroDeCaminho(t *testing.T) {
	if err := gravarJSONSeguro(t.TempDir(), "../fora", map[string]any{"a": 1}); err == nil {
		t.Fatal("esperava erro: caminho fora da pasta do projeto")
	}
}

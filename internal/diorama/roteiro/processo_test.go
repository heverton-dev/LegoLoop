package roteiro

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestProcessoLimiteDuasCorrecoesEGravacaoPosGate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<rss><channel><item><title>Protocolo Aberto 2026</title><link>https://github.com/golang/go/releases</link></item></channel></rss>`))
	}))
	defer srv.Close()

	cat := CatalogoPadrao()
	cat.Conectores = []ConectorFonte{{ID: "tech", Nome: "Tech", URL: srv.URL, Tipo: "rss", Perfis: []string{"tecnologia"}}}

	t.Run("para apos 2 correcoes falhas mesmo em modo mecanico e nao grava acervo", func(t *testing.T) {
		dir := t.TempDir()
		acervoPath := filepath.Join(dir, "acervo", "temas.json")
		p := projetoValidoTeste()
		p.PoliticaGates = PoliticaGates{Modo: "mecanico"}

		chamadas := 0
		pipe := Processo{
			DepsPreVoo: DepsPreVoo{
				ChecarLegoloop: func(context.Context) error { return nil },
				ChecarCofreRea: func() bool { return true },
				ChecarOpcional: func(string) bool { return true },
				HTTPClient:     srv.Client(),
			},
			GerarRoteiro: func(ctx context.Context, pauta PautaCandidata, anterior *RelatorioValidacao) (Script, string, error) {
				chamadas++
				return Script{
					SchemaVersion: 1,
					Title:         pauta.Titulo,
					Sources:       []Fonte{{ID: "FC", URL: "https://blog.invalido.org/x"}},
					Lines:         []LinhaRoteiro{{Text: "Sem prova.", Natureza: "fato_consolidado", Fontes: []string{"FC"}}},
				}, "# plano", nil
			},
		}
		_, err := pipe.Executar(context.Background(), OpcoesExecucao{
			ProjetoDir: dir, AcervoPath: acervoPath, Projeto: p, Marca: marcaValidaTeste(), Catalogo: cat,
		})
		if !errors.Is(err, ErrGateRoteiroReprovado) {
			t.Fatalf("esperava ErrGateRoteiroReprovado, obteve: %v", err)
		}
		if chamadas != 3 {
			t.Fatalf("esperava 1 geracao + 2 correcoes = 3 chamadas, obteve %d", chamadas)
		}
		if _, statErr := os.Stat(acervoPath); !os.IsNotExist(statErr) {
			t.Fatalf("nao deveria gravar acervo/temas.json apos falha")
		}
	})

	t.Run("grava script.json, plano.md e acervo/temas.json somente apos aprovacao", func(t *testing.T) {
		dir := t.TempDir()
		acervoPath := filepath.Join(dir, "acervo", "temas.json")
		p := projetoValidoTeste()
		p.PoliticaGates = PoliticaGates{Modo: "mecanico"}

		tentativa := 0
		pipe := Processo{
			DepsPreVoo: DepsPreVoo{
				ChecarLegoloop: func(context.Context) error { return nil },
				ChecarCofreRea: func() bool { return true },
				ChecarOpcional: func(string) bool { return true },
				HTTPClient:     srv.Client(),
			},
			GerarRoteiro: func(ctx context.Context, pauta PautaCandidata, anterior *RelatorioValidacao) (Script, string, error) {
				tentativa++
				if tentativa == 1 {
					return Script{
						SchemaVersion: 1, Title: pauta.Titulo,
						Sources: []Fonte{{ID: "FC", URL: "https://blog.x.com/1"}},
						Lines:   []LinhaRoteiro{{Text: "Falha inicial", Natureza: "fato_consolidado", Fontes: []string{"FC"}}},
					}, "# plano", nil
				}
				return Script{
					SchemaVersion: 1, Title: pauta.Titulo,
					Sources: []Fonte{{ID: "F1", URL: "https://github.com/golang/go/releases", Entidade: "Go", Data: "2026-11-24"}},
					Lines:   []LinhaRoteiro{{Text: "A release oficial entrou no ar.", Natureza: "fato_consolidado", Fontes: []string{"F1"}}},
				}, "# Plano aprovado", nil
			},
		}
		res, err := pipe.Executar(context.Background(), OpcoesExecucao{
			ProjetoDir: dir, AcervoPath: acervoPath, Projeto: p, Marca: marcaValidaTeste(), Catalogo: cat,
		})
		if err != nil {
			t.Fatalf("esperava sucesso na 1a correcao, obteve erro: %v", err)
		}
		if res.CorrecoesUsadas != 1 {
			t.Fatalf("esperava 1 correcao usada, obteve %d", res.CorrecoesUsadas)
		}
		ac, err := CarregarAcervoTemas(acervoPath)
		if err != nil || len(ac.Temas) != 1 {
			t.Fatalf("esperava 1 tema gravado em acervo/temas.json, obteve err=%v temas=%+v", err, ac.Temas)
		}
	})
}

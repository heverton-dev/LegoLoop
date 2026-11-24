package roteiro

import "testing"

func TestValidarRoteiro(t *testing.T) {
	cat := CatalogoPadrao()
	proj := projetoValidoTeste()

	t.Run("aprova roteiro universal cumprindo todas as 4 naturezas", func(t *testing.T) {
		sc := Script{
			SchemaVersion: 1,
			Title:         "Como funcionam APIs open source",
			Idioma:        "pt-BR",
			Estrutura:     "explicador",
			TipoPauta:     "perene",
			Sources: []Fonte{
				{ID: "F1", URL: "https://github.com/golang/go/releases", Entidade: "Go", Data: "2026-01-10"},
				{ID: "F2", URL: "https://www.reuters.com/tech/1", Entidade: "Reuters", Data: "2026-01-10"},
				{ID: "F3", URL: "https://www.bbc.com/news/2", Entidade: "BBC", Data: "2026-01-10"},
				{ID: "F4", URL: "https://arxiv.org/abs/2601.00001", Entidade: "arXiv", Data: "2026-01-10"},
			},
			Lines: []LinhaRoteiro{
				{Text: "Como um protocolo aberto conecta milhoes de servidores?", Natureza: "narrativa_transicao"},
				{Text: "A versao oficial lancou suporte nativo.", Natureza: "fato_consolidado", Fontes: []string{"F1"}},
				{Text: "Duas agencias confirmaram a adocao global.", Natureza: "fato_consolidado", Fontes: []string{"F2", "F3"}},
				{Text: "Um estudo preliminar com 200 repositorios indica ganho.", Natureza: "evidencia_preliminar", Fontes: []string{"F4"}, Preprint: true, Amostra: "200 repositorios"},
				{Text: "Segundo os mantenedores, o projeto reduz custos.", Natureza: "alegacao_terceiro", Fontes: []string{"F2"}, Atribuicao: "mantenedores"},
			},
		}
		rel := ValidarRoteiro(sc, proj, cat)
		if !rel.Aprovado {
			t.Fatalf("esperava aprovado, obteve violacoes: %+v", rel.Violacoes)
		}
	})

	t.Run("reprova fonte nivel C, fato com apenas 1 fonte B, preprint como fato e alegacao sem atribuicao", func(t *testing.T) {
		sc := Script{
			SchemaVersion: 1,
			Title:         "Roteiro invalido",
			Sources: []Fonte{
				{ID: "FC", URL: "https://blog-aleatorio.exemplo.com/a", Entidade: "Blog", Data: "2026-01-01"},
				{ID: "FB1", URL: "https://www.reuters.com/a", Entidade: "Reuters", Data: "2026-01-01"},
				{ID: "FP", URL: "https://arxiv.org/abs/2601.99999", Entidade: "arXiv", Data: "2026-01-01"},
			},
			Lines: []LinhaRoteiro{
				{Text: "Fato com blog C.", Natureza: "fato_consolidado", Fontes: []string{"FC"}},
				{Text: "Fato com apenas 1 B.", Natureza: "fato_consolidado", Fontes: []string{"FB1"}},
				{Text: "Preprint disfarcado de fato.", Natureza: "fato_consolidado", Fontes: []string{"FP"}},
				{Text: "Alegacao sem atribuicao.", Natureza: "alegacao_terceiro", Fontes: []string{"FB1"}},
			},
		}
		rel := ValidarRoteiro(sc, proj, cat)
		if rel.Aprovado || len(rel.Violacoes) < 4 {
			t.Fatalf("esperava >=4 violacoes, obteve aprovado=%v violacoes=%+v", rel.Aprovado, rel.Violacoes)
		}
	})
}

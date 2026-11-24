package roteiro

import "testing"

func TestClassificarFonte(t *testing.T) {
	cat := CatalogoPadrao()
	projTech := projetoValidoTeste()
	projSaude := projetoValidoTeste()
	projSaude.PerfisFonte = []string{"saude"}

	casos := []struct {
		nome     string
		url      string
		proj     Projeto
		nivel    string
		preprint bool
	}{
		{"gov.br oficial eh nivel A", "https://www.bcb.gov.br/estabilidadefinanceira/pix", projTech, "A", false},
		{"pubmed eh nivel A", "https://pubmed.ncbi.nlm.nih.gov/123456/", projSaude, "A", false},
		{"github no nicho tecnologia eh nivel A", "https://github.com/golang/go/releases/tag/go1.24", projTech, "A", false},
		{"github fora de tecnologia eh fail-closed C", "https://github.com/alguem/repo", projSaude, "C", false},
		{"reuters eh nivel B", "https://www.reuters.com/technology/ai-news", projTech, "B", false},
		{"arxiv marca preprint", "https://arxiv.org/abs/2601.00001", projTech, "A", true},
		{"blog desconhecido eh fail-closed C", "https://blog-aleatorio.exemplo.com/post", projTech, "C", false},
		{"wikipedia eh fail-closed C", "https://pt.wikipedia.org/wiki/Pix", projTech, "C", false},
	}

	for _, tc := range casos {
		t.Run(tc.nome, func(t *testing.T) {
			got := ClassificarFonte(tc.url, tc.proj, cat)
			if got.Nivel != tc.nivel || got.Preprint != tc.preprint {
				t.Fatalf("esperava (%s, preprint=%v), obteve (%s, preprint=%v)", tc.nivel, tc.preprint, got.Nivel, got.Preprint)
			}
		})
	}

	t.Run("fontes_extras promove dominio explicitamente aprovado", func(t *testing.T) {
		p := projetoValidoTeste()
		p.FontesExtras = []string{"docs.especialista.org"}
		got := ClassificarFonte("https://docs.especialista.org/spec", p, cat)
		if got.Nivel != "B" {
			t.Fatalf("esperava nivel B via fontes_extras, obteve %s", got.Nivel)
		}
	})
}

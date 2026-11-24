package roteiro

// CatalogoPadrao retorna o catalogo declarativo gratuito de fontes por nicho e dominio.
func CatalogoPadrao() CatalogoFontes {
	return CatalogoFontes{
		SchemaVersion: 1,
		Perfis:        []string{"tecnologia", "financas", "saude", "ciencia", "historia", "geral"},
		Conectores:    conectoresPadrao(),
		DominiosNivelA: []string{
			"bcb.gov.br", "cvm.gov.br", "sec.gov", "ibge.gov.br",
			"pubmed.ncbi.nlm.nih.gov", "doi.org", "who.int",
			"justice.gov", "clinicaltrials.gov", "europepmc.org",
		},
		SufixosNivelA:    []string{".gov.br", ".gov"},
		DominiosSoftware: []string{"github.com"},
		DominiosNivelB: []string{
			"reuters.com", "apnews.com", "bbc.com", "bloomberg.com",
			"cnbc.com", "folha.uol.com.br", "estadao.com.br",
			"g1.globo.com", "valor.globo.com",
		},
		Preprints: []string{"arxiv.org", "medrxiv.org", "biorxiv.org"},
	}
}

func conectoresPadrao() []ConectorFonte {
	return []ConectorFonte{
		{ID: "hacker_news", Nome: "Hacker News API", URL: "https://hacker-news.firebaseio.com/v0/topstories.json", Tipo: "hn_json", Perfis: []string{"tecnologia"}},
		{ID: "arxiv_cs", Nome: "arXiv CS RSS", URL: "https://rss.arxiv.org/rss/cs", Tipo: "rss", Perfis: []string{"tecnologia", "ciencia"}},
		{ID: "google_noticias", Nome: "Google Noticias RSS", URL: "https://news.google.com/rss/search?q={{QUERY}}&hl=pt-BR&gl=BR&ceid=BR:pt-419", Tipo: "rss", Perfis: []string{"geral", "financas", "tecnologia", "saude", "historia"}},
		{ID: "pubmed_rss", Nome: "PubMed E-utilities", URL: "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=pubmed&retmode=json&term={{QUERY}}", Tipo: "json", Perfis: []string{"saude", "ciencia"}, NivelPadrao: "A"},
		{ID: "bcb_noticias", Nome: "Banco Central RSS", URL: "https://www.bcb.gov.br/api/servico/sitebcb/noticias", Tipo: "json", Perfis: []string{"financas"}, NivelPadrao: "A"},
	}
}

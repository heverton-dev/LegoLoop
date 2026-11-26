package roteiro

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestColetarPautas(t *testing.T) {
	var hitsSaude atomic.Int32
	srvTech := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte(`<rss><channel>
			<item><title>API Open Source Go 1.24</title><link>https://github.com/golang/go/releases</link><pubDate>Mon, 24 Nov 2026</pubDate></item>
			<item><title>Tema Ja Feito Antes</title><link>https://www.reuters.com/tech/old</link><pubDate>Mon, 24 Nov 2026</pubDate></item>
		</channel></rss>`))
	}))
	defer srvTech.Close()

	srvSaude := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hitsSaude.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer srvSaude.Close()

	cat := CatalogoPadrao()
	cat.Conectores = []ConectorFonte{
		{ID: "tech_rss", Nome: "Tech RSS", URL: srvTech.URL, Tipo: "rss", Perfis: []string{"tecnologia"}},
		{ID: "saude_rss", Nome: "Saude RSS", URL: srvSaude.URL, Tipo: "rss", Perfis: []string{"saude"}},
	}
	proj := projetoValidoTeste()
	proj.PerfisFonte = []string{"tecnologia"}
	acervo := AcervoTemas{Temas: []RegistroTema{{Slug: "tema-ja-feito-antes", Titulo: "Tema Ja Feito Antes"}}}

	pautas, err := ColetarPautas(context.Background(), proj, cat, acervo, srvTech.Client())
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if hitsSaude.Load() != 0 {
		t.Fatalf("nao deveria consultar conector de saude para perfil tecnologia")
	}
	if len(pautas) != 1 || pautas[0].Titulo != "API Open Source Go 1.24" {
		t.Fatalf("esperava 1 pauta filtrada sem o tema repetido, obteve: %+v", pautas)
	}
}

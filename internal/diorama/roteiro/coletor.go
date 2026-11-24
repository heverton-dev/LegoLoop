package roteiro

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"strings"
	"time"
)

const limiteCorpoHTTP int64 = 10 << 20 // 10 MB

type PautaCandidata struct {
	Titulo    string `json:"titulo"`
	Slug      string `json:"slug"`
	URL       string `json:"url"`
	Conector  string `json:"conector"`
	Data      string `json:"data,omitempty"`
	Nivel     string `json:"nivel"`
	Preprint  bool   `json:"preprint,omitempty"`
	TipoPauta string `json:"tipo_pauta"`
	Score     int    `json:"score"`
}

type rssDoc struct {
	Channel struct {
		Items []rssItem `xml:"item"`
	} `xml:"channel"`
}

type rssItem struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`
}

// ColetarPautas consulta conectores gratuitos compativeis com perfis_fonte e descarta temas repetidos.
func ColetarPautas(ctx context.Context, p Projeto, cat CatalogoFontes, acervo AcervoTemas, client *http.Client) ([]PautaCandidata, error) {
	if client == nil {
		client = &http.Client{Timeout: 8 * time.Second}
	}
	ativos := cat.ConectoresParaProjeto(p)
	vistos := map[string]bool{}
	var saida []PautaCandidata
	for _, conn := range ativos {
		itens := buscarConector(ctx, conn, p, cat, client)
		saida = filtrarNovasPautas(saida, itens, acervo, vistos)
	}
	return saida, nil
}

func filtrarNovasPautas(acc, itens []PautaCandidata, acervo AcervoTemas, vistos map[string]bool) []PautaCandidata {
	for _, it := range itens {
		slug := NormalizarSlug(it.Titulo)
		if slug == "" || vistos[slug] || acervo.Contem(it.Titulo) {
			continue
		}
		it.Slug = slug
		vistos[slug] = true
		acc = append(acc, it)
	}
	return acc
}

func buscarConector(ctx context.Context, conn ConectorFonte, p Projeto, cat CatalogoFontes, client *http.Client) []PautaCandidata {
	alvo := montarURLConector(conn.URL, p.Nicho)
	corpo, err := baixarLimitado(ctx, client, alvo)
	if err != nil {
		return nil
	}
	if conn.Tipo == "rss" {
		return parsearRSS(corpo, conn, p, cat)
	}
	return parsearJSONBasico(corpo, conn, p, cat)
}

func baixarLimitado(ctx context.Context, client *http.Client, alvo string) ([]byte, error) {
	cctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(cctx, http.MethodGet, alvo, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "LegoLoop-Diorama-Roteiro/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, ErrPreVooBloqueado
	}
	return io.ReadAll(io.LimitReader(resp.Body, limiteCorpoHTTP))
}

func parsearRSS(corpo []byte, conn ConectorFonte, p Projeto, cat CatalogoFontes) []PautaCandidata {
	var doc rssDoc
	if err := xml.Unmarshal(corpo, &doc); err != nil {
		return nil
	}
	var out []PautaCandidata
	for _, it := range doc.Channel.Items {
		cls := ClassificarFonte(it.Link, p, cat)
		out = append(out, PautaCandidata{
			Titulo: strings.TrimSpace(it.Title), URL: strings.TrimSpace(it.Link),
			Conector: conn.ID, Data: strings.TrimSpace(it.PubDate),
			Nivel: cls.Nivel, Preprint: cls.Preprint, TipoPauta: "quente", Score: pontuarNivel(cls.Nivel),
		})
	}
	return out
}

func parsearJSONBasico(corpo []byte, conn ConectorFonte, p Projeto, cat CatalogoFontes) []PautaCandidata {
	var lista []struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	}
	if err := json.Unmarshal(corpo, &lista); err != nil {
		return nil
	}
	var out []PautaCandidata
	for _, it := range lista {
		cls := ClassificarFonte(it.URL, p, cat)
		out = append(out, PautaCandidata{
			Titulo: strings.TrimSpace(it.Title), URL: strings.TrimSpace(it.URL),
			Conector: conn.ID, Nivel: cls.Nivel, Preprint: cls.Preprint,
			TipoPauta: "perene", Score: pontuarNivel(cls.Nivel),
		})
	}
	return out
}

func pontuarNivel(nivel string) int {
	switch nivel {
	case "A":
		return 10
	case "B":
		return 7
	default:
		return 2
	}
}

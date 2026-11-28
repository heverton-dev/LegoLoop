package distribuicao

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// MinimoMetricas e o menor numero de medidas para gerar recomendacoes.
const MinimoMetricas = 3

// Metrica e a medida de um video ou corte publicado.
type Metrica struct {
	ItemID          string  `json:"item_id"`
	Tema            string  `json:"tema"`
	Formato         string  `json:"formato"`
	DuracaoS        float64 `json:"duracao_s"`
	Retencao        float64 `json:"retencao"`
	Alcance         int     `json:"alcance"`
	Hora            int     `json:"hora_utc"`
	NovosSeguidores int     `json:"novos_seguidores"`
}

func (m Metrica) validar() error {
	if m.ItemID == "" || m.Tema == "" || m.Formato == "" || m.DuracaoS <= 0 || m.Retencao < 0 || m.Retencao > 1 || m.Hora < 0 || m.Hora > 23 || m.Alcance < 0 || m.NovosSeguidores < 0 {
		return fmt.Errorf("%w: metrica %q fora do contrato", ErrPedidoInvalido, m.ItemID)
	}
	return nil
}

// media de um campo agrupado por chave; devolve a chave de maior media (empate: ordem alfabetica).
func melhor(ms []Metrica, chave func(Metrica) string, valor func(Metrica) float64) (string, float64) {
	soma, n := map[string]float64{}, map[string]int{}
	for _, m := range ms {
		k := chave(m)
		soma[k] += valor(m)
		n[k]++
	}
	chaves := make([]string, 0, len(soma))
	for k := range soma {
		chaves = append(chaves, k)
	}
	sort.Strings(chaves)
	melhorK, melhorV := "", -1.0
	for _, k := range chaves {
		if v := soma[k] / float64(n[k]); v > melhorV {
			melhorK, melhorV = k, v
		}
	}
	return melhorK, melhorV
}

func faixaDuracao(s float64) string {
	switch {
	case s < 30:
		return "ate 30 s"
	case s < 60:
		return "30 a 60 s"
	case s < 180:
		return "1 a 3 min"
	}
	return "mais de 3 min"
}

// GerarAprendizado transforma as medidas em recomendacoes que voltam para a descoberta de pautas.
func GerarAprendizado(ms []Metrica, agora time.Time) (string, error) {
	if len(ms) < MinimoMetricas {
		return "", fmt.Errorf("%w: %d medidas; minimo %d", ErrPedidoInvalido, len(ms), MinimoMetricas)
	}
	for _, m := range ms {
		if err := m.validar(); err != nil {
			return "", err
		}
	}
	ret := func(m Metrica) float64 { return m.Retencao }
	tema, rt := melhor(ms, func(m Metrica) string { return m.Tema }, ret)
	formato, rf := melhor(ms, func(m Metrica) string { return m.Formato }, ret)
	faixa, rd := melhor(ms, func(m Metrica) string { return faixaDuracao(m.DuracaoS) }, ret)
	hora, ah := melhor(ms, func(m Metrica) string { return fmt.Sprintf("%02d", m.Hora) }, func(m Metrica) float64 { return float64(m.Alcance) })
	item, ns := melhor(ms, func(m Metrica) string { return m.ItemID }, func(m Metrica) float64 { return float64(m.NovosSeguidores) })
	var b strings.Builder
	fmt.Fprintf(&b, "# Aprendizado semanal\n\nGerado em %s com %d medidas.\n\n## Recomendacoes\n\n", agora.UTC().Format("2006-01-02"), len(ms))
	fmt.Fprintf(&b, "- Tema com maior retencao media: %q (%.0f%%). Priorizar pautas parecidas.\n", tema, rt*100)
	fmt.Fprintf(&b, "- Formato com maior retencao media: %s (%.0f%%).\n", formato, rf*100)
	fmt.Fprintf(&b, "- Faixa de duracao com maior retencao media: %s (%.0f%%).\n", faixa, rd*100)
	fmt.Fprintf(&b, "- Horario (UTC) com maior alcance medio: %sh (%.0f).\n", hora, ah)
	fmt.Fprintf(&b, "- Item que mais trouxe seguidores: %s (%.0f). Repetir o gancho dele.\n", item, ns)
	return b.String(), nil
}

// GravarAprendizado escreve acervo/aprendizado.md, que a descoberta de pautas le antes de sugerir temas.
func GravarAprendizado(projDir string, ms []Metrica, agora time.Time) (string, error) {
	texto, err := GerarAprendizado(ms, agora)
	if err != nil {
		return "", err
	}
	caminho := filepath.Join(projDir, "acervo", "aprendizado.md")
	return caminho, gravarAtomico(caminho, []byte(texto))
}

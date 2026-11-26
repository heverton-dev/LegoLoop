package render

import (
	"fmt"
	"math"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/diorama/roteiro"
)

// LarguraSRT e o maximo de caracteres por linha de legenda no arquivo .srt.
const LarguraSRT = 42

// unidade e uma palavra normalizada com o tempo em que foi falada.
type unidade struct {
	txt      string
	ini, fim float64
}

// unidadesWhisper quebra as palavras ouvidas em unidades normalizadas.
// Sinais sem letra nem digito (como "%", falado "por cento") alongam a unidade anterior.
func unidadesWhisper(ps []PalavraTempo) []unidade {
	var us []unidade
	pendente, temPendente := 0.0, false
	for _, p := range ps {
		ws := palavrasDe(p.Texto)
		if len(ws) == 0 {
			if n := len(us); n > 0 {
				us[n-1].fim = math.Max(us[n-1].fim, p.FimS)
			} else if !temPendente {
				pendente, temPendente = p.InicioS, true
			}
			continue
		}
		primeira := len(us) == 0
		passo := (p.FimS - p.InicioS) / float64(len(ws))
		for k, w := range ws {
			us = append(us, unidade{txt: w, ini: p.InicioS + passo*float64(k), fim: p.InicioS + passo*float64(k+1)})
		}
		us[len(us)-1].fim = p.FimS
		if primeira && temPendente {
			us[0].ini = pendente
		}
	}
	return us
}

// parear alinha as unidades do roteiro (a) as ouvidas (b) pela menor distancia de edicao.
// Devolve, para cada unidade do roteiro, o indice em b ou -1 se nao foi ouvida.
func parear(a []string, b []unidade) []int {
	n, m := len(a), len(b)
	d := make([][]int, n+1)
	for i := range d {
		d[i] = make([]int, m+1)
		d[i][0] = i
	}
	for j := 0; j <= m; j++ {
		d[0][j] = j
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			d[i][j] = min(d[i-1][j-1]+custoTroca(a[i-1], b[j-1].txt), d[i-1][j]+1, d[i][j-1]+1)
		}
	}
	par := make([]int, n)
	for i := range par {
		par[i] = -1
	}
	for i, j := n, m; i > 0 && j > 0; {
		switch {
		case d[i][j] == d[i-1][j-1]+custoTroca(a[i-1], b[j-1].txt):
			par[i-1] = j - 1
			i, j = i-1, j-1
		case d[i][j] == d[i-1][j]+1:
			i--
		default:
			j--
		}
	}
	return par
}

func custoTroca(a, b string) int {
	if a == b {
		return 0
	}
	return 1
}

// temposDasUnidades da a cada unidade o tempo da unidade ouvida; as nao ouvidas dividem a lacuna entre vizinhas.
func temposDasUnidades(par []int, b []unidade, ini, fim float64) [][2]float64 {
	t := make([][2]float64, len(par))
	for i, k := range par {
		if k >= 0 {
			t[i] = [2]float64{b[k].ini, b[k].fim}
		}
	}
	limiteIni, limiteFim := ini, fim
	if len(b) > 0 {
		limiteIni, limiteFim = b[0].ini, b[len(b)-1].fim
	}
	for i := 0; i < len(par); {
		if par[i] >= 0 {
			i++
			continue
		}
		j := i
		for j < len(par) && par[j] < 0 {
			j++
		}
		preencherLacuna(t, i, j, limiteIni, limiteFim)
		i = j
	}
	return t
}

func preencherLacuna(t [][2]float64, i, j int, limiteIni, limiteFim float64) {
	antes, depois := limiteIni, limiteFim
	if i > 0 {
		antes = t[i-1][1]
	}
	if j < len(t) {
		depois = t[j][0]
	}
	vao, m := math.Max(depois-antes, 0), float64(j-i)
	for k := 0; k < j-i; k++ {
		t[i+k] = [2]float64{antes + vao*float64(k)/m, antes + vao*float64(k+1)/m}
	}
}

// montarTokens agrupa as unidades de volta nas palavras do roteiro; sinais soltos grudam na vizinha.
func montarTokens(texto string, t [][2]float64) []PalavraTempo {
	var saida []PalavraTempo
	u, prefixo := 0, ""
	for _, tok := range strings.Fields(texto) {
		n := len(palavrasDe(tok))
		if n == 0 {
			if k := len(saida); k > 0 {
				saida[k-1].Texto += " " + tok
			} else {
				prefixo = strings.TrimSpace(prefixo + " " + tok)
			}
			continue
		}
		p := PalavraTempo{Texto: tok, InicioS: t[u][0], FimS: t[u+n-1][1]}
		if prefixo != "" {
			p.Texto, prefixo = prefixo+" "+tok, ""
		}
		saida = append(saida, p)
		u += n
	}
	return saida
}

// limitarJanela garante ordem, ausencia de sobreposicao e permanencia dentro da janela.
func limitarJanela(ps []PalavraTempo, ini, fim float64) []PalavraTempo {
	cursor := ini
	for i, p := range ps {
		pi := math.Min(math.Max(p.InicioS, cursor), fim)
		pf := math.Min(math.Max(p.FimS, pi), fim)
		ps[i].InicioS, ps[i].FimS = pi, pf
		cursor = pf
	}
	return ps
}

func alinharLinha(texto string, l LinhaTempo) []PalavraTempo {
	a := palavrasDe(texto)
	b := unidadesWhisper(l.Palavras)
	tokens := montarTokens(texto, temposDasUnidades(parear(a, b), b, l.InicioS, l.FimS))
	if len(tokens) == 0 {
		return l.Palavras
	}
	return limitarJanela(tokens, l.InicioS, l.FimS)
}

// AlinharLegendas troca as palavras ouvidas pelas do roteiro, mantendo o tempo do Whisper.
// A legenda mostra o texto do roteiro, que e a fonte de verdade; o Whisper so da o tempo.
func AlinharLegendas(script roteiro.Script, tm Timings) Timings {
	linhas := append([]LinhaTempo(nil), tm.Linhas...)
	for i := range linhas {
		if i >= len(script.Lines) {
			break
		}
		linhas[i].Palavras = alinharLinha(script.Lines[i].Text, tm.Linhas[i])
	}
	tm.Linhas = linhas
	return tm
}

// AplicarLegendas grava, na camada de legenda de cada batida, as palavras que comecam dentro dela.
// Palavra que cruza o fim da batida e recortada nele. A entrada nao muda.
func AplicarLegendas(tl Timeline, tmL Timings) Timeline {
	saida := Timeline{SchemaVersion: tl.SchemaVersion, DuracaoS: tl.DuracaoS, Batidas: make([]Batida, len(tl.Batidas))}
	for i, b := range tl.Batidas {
		b.Camadas = append([]CamadaItem(nil), b.Camadas...)
		for k := range b.Camadas {
			if b.Camadas[k].Camada == "legenda" {
				b.Camadas[k].Palavras = palavrasDaBatida(b, tmL)
			}
		}
		saida.Batidas[i] = b
	}
	return saida
}

func palavrasDaBatida(b Batida, tmL Timings) []PalavraTempo {
	if b.IndiceLinha < 0 || b.IndiceLinha >= len(tmL.Linhas) {
		return nil
	}
	l := tmL.Linhas[b.IndiceLinha]
	ultima := b.FimS == l.FimS
	var ps []PalavraTempo
	for _, p := range l.Palavras {
		if p.InicioS >= b.InicioS && (p.InicioS < b.FimS || (ultima && p.InicioS <= b.FimS)) {
			p.FimS = math.Min(p.FimS, b.FimS)
			ps = append(ps, p)
		}
	}
	return ps
}

func camadaLegenda(b Batida) CamadaItem {
	for _, c := range b.Camadas {
		if c.Camada == "legenda" {
			return c
		}
	}
	return CamadaItem{}
}

func textosDasPalavras(ps []PalavraTempo) []string {
	t := make([]string, len(ps))
	for i, p := range ps {
		t[i] = p.Texto
	}
	return t
}

// quebrarLinhas junta palavras em linhas de ate largura caracteres, sem partir palavras.
// Quando cabem em duas linhas, equilibra o tamanho delas para nao deixar uma palavra solta.
func quebrarLinhas(palavras []string, largura int) []string {
	linhas := quebrarGuloso(palavras, largura)
	if len(linhas) == 2 {
		if b := quebrarBalanceado(palavras, largura); b != nil {
			return b
		}
	}
	return linhas
}

func quebrarGuloso(palavras []string, largura int) []string {
	var linhas []string
	atual := ""
	for _, w := range palavras {
		switch {
		case atual == "":
			atual = w
		case len([]rune(atual))+1+len([]rune(w)) <= largura:
			atual += " " + w
		default:
			linhas = append(linhas, atual)
			atual = w
		}
	}
	return append(linhas, atual)
}

// quebrarBalanceado escolhe o corte que minimiza a maior das duas linhas; nil se nenhuma cabe.
func quebrarBalanceado(p []string, largura int) []string {
	var melhor []string
	menor := math.MaxInt
	for k := 1; k < len(p); k++ {
		a, b := strings.Join(p[:k], " "), strings.Join(p[k:], " ")
		la, lb := len([]rune(a)), len([]rune(b))
		if la <= largura && lb <= largura && max(la, lb) < menor {
			menor, melhor = max(la, lb), []string{a, b}
		}
	}
	return melhor
}

// formatarTempoSRT escreve HH:MM:SS,mmm a partir do total de milissegundos arredondado.
func formatarTempoSRT(s float64) string {
	ms := int64(math.Round(math.Max(s, 0) * 1000))
	return fmt.Sprintf("%02d:%02d:%02d,%03d", ms/3600000, (ms/60000)%60, (ms/1000)%60, ms%1000)
}

// GerarSRT cria um bloco por batida com legenda, em ordem e sem sobreposicao.
func GerarSRT(tl Timeline) (string, error) {
	var sb strings.Builder
	n, fimAnterior := 0, 0.0
	for _, b := range tl.Batidas {
		ps := camadaLegenda(b).Palavras
		if len(ps) == 0 {
			continue
		}
		ini, fim := ps[0].InicioS, ps[len(ps)-1].FimS
		if ini < fimAnterior {
			return "", invalido("legenda do bloco %d comeca antes do fim do anterior", n+1)
		}
		n, fimAnterior = n+1, fim
		texto := strings.Join(quebrarLinhas(textosDasPalavras(ps), LarguraSRT), "\n")
		fmt.Fprintf(&sb, "%d\n%s --> %s\n%s\n\n", n, formatarTempoSRT(ini), formatarTempoSRT(fim), texto)
	}
	if n == 0 {
		return "", invalido("timeline sem legendas")
	}
	return sb.String(), nil
}

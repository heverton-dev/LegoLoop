package render

import (
	"fmt"
	"math"
	"sort"
)

// Metas do gate check_timeline (Mapa Macro, Etapa 7): 12 a 21 cortes por minuto, mediana de 2,5 a 4 s.
const (
	AlvoBatidaS     = 3.2
	MinCortesPorMin = 12.0
	MaxCortesPorMin = 21.0
	MedianaMinS     = 2.5
	MedianaMaxS     = 4.0
	derivaMaxRel    = 0.35
)

// RitmoMetricas resume o ritmo medido da timeline.
type RitmoMetricas struct {
	Batidas      int
	CortesPorMin float64
	MedianaS     float64
}

// numeroDeBatidas escolhe quantas batidas cabem numa fala, mirando AlvoBatidaS por batida.
func numeroDeBatidas(duracaoS float64) int {
	return max(1, arredondar(duracaoS/AlvoBatidaS))
}

// cameraDaBatida mantem a camera do storyboard nas batidas pares e alterna nas impares.
func cameraDaBatida(base string, k int) string {
	if k%2 == 0 {
		return base
	}
	if base == "zoom" {
		return "pan"
	}
	return "zoom"
}

// SegmentarBatidas divide cada fala longa em batidas, cortando nas fronteiras de palavra.
// Devolve uma timeline nova; a entrada nao muda.
func SegmentarBatidas(tl Timeline, tm Timings) Timeline {
	saida := Timeline{SchemaVersion: tl.SchemaVersion, DuracaoS: tl.DuracaoS}
	for _, b := range tl.Batidas {
		saida.Batidas = append(saida.Batidas, dividirBatida(b, palavrasDaLinha(tm, b.IndiceLinha))...)
	}
	return saida
}

func palavrasDaLinha(tm Timings, idx int) []PalavraTempo {
	if idx < 0 || idx >= len(tm.Linhas) {
		return nil
	}
	return tm.Linhas[idx].Palavras
}

func dividirBatida(b Batida, palavras []PalavraTempo) []Batida {
	n := numeroDeBatidas(b.FimS - b.InicioS)
	if n == 1 {
		b.Camadas = append([]CamadaItem(nil), b.Camadas...)
		b.SFX = append([]EventoSFX(nil), b.SFX...)
		return []Batida{b}
	}
	limites := append(append([]float64{b.InicioS}, pontosDeCorte(b.InicioS, b.FimS, palavras, n)...), b.FimS)
	saida := make([]Batida, n)
	for k := range saida {
		nb := b
		nb.InicioS, nb.FimS = limites[k], limites[k+1]
		nb.Camera = cameraDaBatida(b.Camera, k)
		nb.Camadas = camadasDaBatida(b.Camadas, k)
		nb.SFX = sfxDaBatida(b.SFX, nb.InicioS, nb.FimS, k == n-1)
		saida[k] = nb
	}
	return saida
}

// camadasDaBatida copia as camadas e da a cada legenda um id proprio por batida.
func camadasDaBatida(camadas []CamadaItem, k int) []CamadaItem {
	copia := append([]CamadaItem(nil), camadas...)
	for i := range copia {
		if copia[i].Camada == "legenda" {
			copia[i].ID = fmt.Sprintf("%s:b%d", copia[i].ID, k)
		}
	}
	return copia
}

// sfxDaBatida entrega cada evento a batida que contem seu tempo; o fim so vale na ultima.
func sfxDaBatida(eventos []EventoSFX, inicio, fim float64, ultima bool) []EventoSFX {
	var saida []EventoSFX
	for _, e := range eventos {
		if e.TempoS >= inicio && (e.TempoS < fim || (ultima && e.TempoS == fim)) {
			saida = append(saida, e)
		}
	}
	return saida
}

// pontosDeCorte busca, perto de cada ponto ideal, a fronteira de palavra mais proxima.
func pontosDeCorte(inicio, fim float64, palavras []PalavraTempo, n int) []float64 {
	passo := (fim - inicio) / float64(n)
	candidatas := fronteirasDePalavras(palavras)
	cortes := make([]float64, 0, n-1)
	anterior := inicio
	for k := 1; k < n; k++ {
		corte := escolherCorte(inicio+passo*float64(k), anterior, passo*derivaMaxRel, candidatas)
		cortes = append(cortes, corte)
		anterior = corte
	}
	return cortes
}

func fronteirasDePalavras(p []PalavraTempo) []float64 {
	var f []float64
	for i := 0; i+1 < len(p); i++ {
		f = append(f, (p[i].FimS+p[i+1].InicioS)/2)
	}
	return f
}

// escolherCorte limita o desvio a derivaMaxRel do passo e exige corte depois do anterior.
func escolherCorte(ideal, anterior, deriva float64, candidatas []float64) float64 {
	melhor, menor := ideal, math.Inf(1)
	for _, c := range candidatas {
		d := math.Abs(c - ideal)
		if c > anterior && d <= deriva && d < menor {
			melhor, menor = c, d
		}
	}
	return melhor
}

func mediana(v []float64) float64 {
	if len(v) == 0 {
		return 0
	}
	o := append([]float64(nil), v...)
	sort.Float64s(o)
	if len(o)%2 == 1 {
		return o[len(o)/2]
	}
	return (o[len(o)/2-1] + o[len(o)/2]) / 2
}

// ChecarRitmo e o gate check_timeline: confere cortes por minuto e mediana das batidas.
func ChecarRitmo(tl Timeline) (RitmoMetricas, error) {
	if len(tl.Batidas) == 0 || tl.DuracaoS <= 0 {
		return RitmoMetricas{}, fmt.Errorf("%w: timeline vazia", ErrTimelineRitmo)
	}
	durs := make([]float64, len(tl.Batidas))
	for i, b := range tl.Batidas {
		durs[i] = b.FimS - b.InicioS
	}
	m := RitmoMetricas{Batidas: len(durs), CortesPorMin: float64(len(durs)) / (tl.DuracaoS / 60), MedianaS: mediana(durs)}
	if m.CortesPorMin < MinCortesPorMin || m.CortesPorMin > MaxCortesPorMin {
		return m, fmt.Errorf("%w: cortes por minuto %.1f fora de %.0f a %.0f", ErrTimelineRitmo, m.CortesPorMin, MinCortesPorMin, MaxCortesPorMin)
	}
	if m.MedianaS < MedianaMinS || m.MedianaS > MedianaMaxS {
		return m, fmt.Errorf("%w: mediana de batida %.2f s fora de %.1f a %.1f s", ErrTimelineRitmo, m.MedianaS, MedianaMinS, MedianaMaxS)
	}
	return m, nil
}

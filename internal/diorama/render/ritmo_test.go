package render

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/diorama/ativos"
)

// linhaContigua cria n palavras coladas, cada uma com a duracao passo.
func linhaContigua(i int, inicio float64, n int, passo float64) LinhaTempo {
	l := LinhaTempo{IndiceLinha: i, InicioS: inicio, FimS: inicio + float64(n)*passo}
	for k := 0; k < n; k++ {
		l.Palavras = append(l.Palavras, PalavraTempo{Texto: fmt.Sprintf("w%d", k), InicioS: inicio + float64(k)*passo, FimS: inicio + float64(k+1)*passo})
	}
	return l
}

// timingsLongo tem 2 linhas de 9,75 s, como as falas longas do demo.
func timingsLongo() Timings {
	return Timings{SchemaVersion: 1, Provedor: "p", Voz: "v", Linhas: []LinhaTempo{
		linhaContigua(0, 0, 12, 0.8125), linhaContigua(1, 9.75, 12, 0.8125),
	}}
}

func batidaUnica(linha int, inicio, fim float64, camera string) Batida {
	return Batida{IndiceLinha: linha, InicioS: inicio, FimS: fim, CenarioID: "cen_01", Camera: camera, Camadas: []CamadaItem{
		{ID: "cenario:cen_01", Camada: "cenario", Profundidade: 0, Arquivo: "flow/cen_01.png", Caixa: ativos.BoundingBox{W: 10, H: 10}},
		{ID: "p", Camada: "sprite", Profundidade: 1, SpriteID: "narrador_pose_1", Caixa: ativos.BoundingBox{W: 5, H: 5}},
		{ID: fmt.Sprintf("legenda:linha_%02d", linha), Camada: "legenda", Profundidade: 2, Caixa: ativos.BoundingBox{W: 8, H: 2}},
	}}
}

func timelineDe(batidas ...Batida) Timeline {
	return Timeline{SchemaVersion: 1, DuracaoS: batidas[len(batidas)-1].FimS, Batidas: batidas}
}

func timelineDeDuracoes(duracoes ...float64) Timeline {
	var bs []Batida
	t := 0.0
	for i, d := range duracoes {
		bs = append(bs, batidaUnica(i, t, t+d, "pan"))
		t += d
	}
	return timelineDe(bs...)
}

func repetir(v float64, n int) []float64 {
	r := make([]float64, n)
	for i := range r {
		r[i] = v
	}
	return r
}

func TestRitmoNumeroDeBatidas(t *testing.T) {
	casos := map[float64]int{2.0: 1, 3.0: 1, 5.7: 2, 7.3: 2, 8.4: 3, 8.8: 3, 10.2: 3, 12.5: 4, 0.5: 1}
	for d, esperado := range casos {
		if got := numeroDeBatidas(d); got != esperado {
			t.Errorf("duracao %.1f: obteve %d batidas, quer %d", d, got, esperado)
		}
	}
}

func TestRitmoSegmentaSemBuracosEMesmoCortePalavra(t *testing.T) {
	tm := timingsLongo()
	entrada := timelineDe(batidaUnica(0, 0, 9.75, "pan"))
	saida := SegmentarBatidas(entrada, tm)
	if len(saida.Batidas) != 3 {
		t.Fatalf("9,75 s deve virar 3 batidas: %d", len(saida.Batidas))
	}
	if saida.Batidas[0].InicioS != 0 || saida.Batidas[2].FimS != 9.75 || saida.DuracaoS != 9.75 {
		t.Fatalf("deve cobrir a linha inteira: %+v", saida.Batidas)
	}
	for i := 0; i < 2; i++ {
		if saida.Batidas[i].FimS != saida.Batidas[i+1].InicioS {
			t.Fatalf("buraco ou sobreposicao entre batidas %d e %d", i, i+1)
		}
	}
	if math.Abs(saida.Batidas[0].FimS-3.25) > 1e-9 || math.Abs(saida.Batidas[1].FimS-6.5) > 1e-9 {
		t.Fatalf("cortes devem cair nas fronteiras de palavra mais proximas de 3,25 e 6,5: %v %v", saida.Batidas[0].FimS, saida.Batidas[1].FimS)
	}
	if err := saida.Validar(); err != nil {
		t.Fatalf("timeline segmentada deve validar: %v", err)
	}
}

func TestRitmoEncaixaNaFronteiraMaisProxima(t *testing.T) {
	limites := []float64{0, 1.0, 3.0, 3.7, 6.9, 7.6, 9.75}
	l := LinhaTempo{IndiceLinha: 0, InicioS: 0, FimS: 9.75}
	for i := 0; i+1 < len(limites); i++ {
		l.Palavras = append(l.Palavras, PalavraTempo{Texto: fmt.Sprintf("w%d", i), InicioS: limites[i], FimS: limites[i+1]})
	}
	saida := SegmentarBatidas(timelineDe(batidaUnica(0, 0, 9.75, "pan")), Timings{Linhas: []LinhaTempo{l}})
	if saida.Batidas[0].FimS != 3.0 || saida.Batidas[1].FimS != 6.9 {
		t.Fatalf("esperava cortes em 3.0 e 6.9, obteve %v e %v", saida.Batidas[0].FimS, saida.Batidas[1].FimS)
	}
}

func TestRitmoSemPalavrasCortaNoPontoIdeal(t *testing.T) {
	saida := SegmentarBatidas(timelineDe(batidaUnica(0, 0, 9.75, "pan")), Timings{Linhas: []LinhaTempo{{IndiceLinha: 0, InicioS: 0, FimS: 9.75}}})
	if math.Abs(saida.Batidas[0].FimS-3.25) > 1e-9 || math.Abs(saida.Batidas[1].FimS-6.5) > 1e-9 {
		t.Fatalf("sem palavras corta em 3,25 e 6,5: %+v", saida.Batidas)
	}
	saida = SegmentarBatidas(timelineDe(batidaUnica(5, 0, 9.75, "pan")), Timings{})
	if len(saida.Batidas) != 3 {
		t.Fatalf("linha sem timing deve segmentar pelo tempo da batida: %d", len(saida.Batidas))
	}
}

func TestRitmoCameraAlternaEntreBatidas(t *testing.T) {
	casos := map[string][]string{
		"pan": {"pan", "zoom", "pan"}, "zoom": {"zoom", "pan", "zoom"},
		"estatica": {"estatica", "zoom", "estatica"}, "shake": {"shake", "zoom", "shake"},
	}
	for base, esperado := range casos {
		saida := SegmentarBatidas(timelineDe(batidaUnica(0, 0, 9.75, base)), timingsLongo())
		for i, b := range saida.Batidas {
			if b.Camera != esperado[i] || cameraDaBatida(base, i) != esperado[i] {
				t.Errorf("base %s batida %d: camera %s, quer %s", base, i, b.Camera, esperado[i])
			}
		}
	}
}

func TestRitmoLegendaUnicaPorBatidaECopiaIndependente(t *testing.T) {
	entrada := timelineDe(batidaUnica(0, 0, 9.75, "pan"), batidaUnica(1, 9.75, 12.75, "pan"))
	saida := SegmentarBatidas(entrada, timingsLongo())
	vistos := map[string]bool{}
	for _, b := range saida.Batidas {
		for _, c := range b.Camadas {
			if c.Camada == "legenda" {
				if vistos[c.ID] {
					t.Fatalf("id de legenda repetido: %s", c.ID)
				}
				vistos[c.ID] = true
			}
		}
	}
	if !vistos["legenda:linha_00:b0"] || !vistos["legenda:linha_00:b2"] || !vistos["legenda:linha_01"] {
		t.Fatalf("ids esperados com sufixo por batida so quando a linha foi dividida: %v", vistos)
	}
	saida.Batidas[0].Camadas[0].ID = "alterado"
	if saida.Batidas[1].Camadas[0].ID == "alterado" || entrada.Batidas[0].Camadas[0].ID == "alterado" {
		t.Fatal("camadas devem ser copias independentes")
	}
}

func TestRitmoSegmentacaoDeterministicaENaoAlteraEntrada(t *testing.T) {
	entrada := timelineDe(batidaUnica(0, 0, 9.75, "pan"))
	a := SegmentarBatidas(entrada, timingsLongo())
	b := SegmentarBatidas(entrada, timingsLongo())
	if fmt.Sprintf("%+v", a) != fmt.Sprintf("%+v", b) {
		t.Fatal("duas segmentacoes devem ser identicas")
	}
	if len(entrada.Batidas) != 1 || entrada.Batidas[0].FimS != 9.75 || entrada.Batidas[0].Camadas[2].ID != "legenda:linha_00" {
		t.Fatal("a entrada nao pode ser alterada")
	}
}

func TestRitmoLinhaCurtaNaoDivide(t *testing.T) {
	saida := SegmentarBatidas(timelineDe(batidaUnica(0, 0, 3.0, "zoom")), timingsLongo())
	if len(saida.Batidas) != 1 || saida.Batidas[0].Camera != "zoom" || saida.Batidas[0].Camadas[2].ID != "legenda:linha_00" {
		t.Fatalf("linha de 3 s fica em 1 batida sem sufixo: %+v", saida.Batidas)
	}
}

func TestRitmoChecarAprova(t *testing.T) {
	m, err := ChecarRitmo(timelineDeDuracoes(repetir(52.9/17, 17)...))
	if err != nil {
		t.Fatalf("17 batidas em 52,9 s deve passar: %v", err)
	}
	if m.Batidas != 17 || math.Abs(m.CortesPorMin-19.28) > 0.01 || math.Abs(m.MedianaS-52.9/17) > 1e-9 {
		t.Fatalf("metricas incorretas: %+v", m)
	}
}

func TestRitmoChecarReprova(t *testing.T) {
	casos := map[string]struct {
		duracoes []float64
		trecho   string
	}{
		"poucos cortes": {append(repetir(8.8, 6), 0), "cortes por minuto"},
		"cortes demais": {repetir(1.5, 30), "cortes por minuto"},
		"mediana curta": {append(repetir(2, 6), repetir(8, 3)...), "mediana"},
		"mediana longa": {append(repetir(4.5, 5), 1, 1), "mediana"},
	}
	for nome, c := range casos {
		t.Run(nome, func(t *testing.T) {
			dur := c.duracoes
			if dur[len(dur)-1] == 0 {
				dur = dur[:len(dur)-1]
			}
			_, err := ChecarRitmo(timelineDeDuracoes(dur...))
			if !errors.Is(err, ErrTimelineRitmo) || !strings.Contains(err.Error(), c.trecho) {
				t.Fatalf("esperava ErrTimelineRitmo citando %q: %v", c.trecho, err)
			}
		})
	}
	if _, err := ChecarRitmo(Timeline{}); !errors.Is(err, ErrTimelineRitmo) {
		t.Fatalf("timeline vazia deve reprovar: %v", err)
	}
}

func TestRitmoMediana(t *testing.T) {
	if m := mediana([]float64{5, 2, 4, 3}); m != 3.5 {
		t.Fatalf("mediana par: %v", m)
	}
	if m := mediana([]float64{9, 1, 5}); m != 5 {
		t.Fatalf("mediana impar: %v", m)
	}
	if m := mediana(nil); m != 0 {
		t.Fatalf("mediana vazia: %v", m)
	}
}

func TestRitmoSFXVaiParaABatidaQueContemOTempo(t *testing.T) {
	b := batidaUnica(0, 0, 9.75, "pan")
	b.SFX = []EventoSFX{
		{Tipo: "troca_cena", TempoS: 0, Arquivo: "a.wav"},
		{Tipo: "movimento_objeto", TempoS: 3.25, Arquivo: "b.wav"},
		{Tipo: "destaque", TempoS: 9.75, Arquivo: "c.wav"},
	}
	saida := SegmentarBatidas(timelineDe(b), timingsLongo())
	quantos := []int{len(saida.Batidas[0].SFX), len(saida.Batidas[1].SFX), len(saida.Batidas[2].SFX)}
	if quantos[0] != 1 || quantos[1] != 1 || quantos[2] != 1 {
		t.Fatalf("cada evento vai para uma unica batida (o do fim exato, para a ultima): %v", quantos)
	}
	if saida.Batidas[1].SFX[0].Arquivo != "b.wav" || saida.Batidas[2].SFX[0].Arquivo != "c.wav" {
		t.Fatalf("evento no corte pertence a batida que comeca ali: %+v", saida.Batidas)
	}
	if err := saida.Validar(); err != nil {
		t.Fatalf("sfx dentro da janela de cada batida deve validar: %v", err)
	}
}

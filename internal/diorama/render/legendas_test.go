package render

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/diorama/roteiro"
)

func pt(txt string, ini, fim float64) PalavraTempo {
	return PalavraTempo{Texto: txt, InicioS: ini, FimS: fim}
}

func scriptDe(linhas ...string) roteiro.Script {
	var s roteiro.Script
	for _, l := range linhas {
		s.Lines = append(s.Lines, roteiro.LinhaRoteiro{Text: l})
	}
	return s
}

func tempoDe(linhas ...LinhaTempo) Timings {
	return Timings{SchemaVersion: 1, Provedor: "p", Voz: "v", Linhas: linhas}
}

func linhaCom(i int, ini, fim float64, ps ...PalavraTempo) LinhaTempo {
	return LinhaTempo{IndiceLinha: i, InicioS: ini, FimS: fim, Palavras: ps}
}

func textosDe(ps []PalavraTempo) []string {
	t := make([]string, len(ps))
	for i, p := range ps {
		t[i] = p.Texto
	}
	return t
}

func seq(textos []string, passo float64) []PalavraTempo {
	ps := make([]PalavraTempo, len(textos))
	for i, tx := range textos {
		ps[i] = pt(tx, float64(i)*passo, float64(i+1)*passo)
	}
	return ps
}

func TestLegendasUsaOTextoDoRoteiroComTempoDoWhisper(t *testing.T) {
	w := seq([]string{"A", "conta", "pública", "atingir", "um", "ponto."}, 0.5)
	saida := AlinharLegendas(scriptDe("A conta pública atingiu um ponto."), tempoDe(linhaCom(0, 0, 3, w...)))
	got := saida.Linhas[0].Palavras
	if !reflect.DeepEqual(textosDe(got), []string{"A", "conta", "pública", "atingiu", "um", "ponto."}) {
		t.Fatalf("a legenda deve mostrar o texto do roteiro e nao o ouvido: %v", textosDe(got))
	}
	for i, p := range got {
		if p.InicioS != float64(i)*0.5 || p.FimS != float64(i+1)*0.5 {
			t.Fatalf("tempo deve ser o do Whisper: %+v", p)
		}
	}
}

func TestLegendasNumerosEPontuacaoViramUmaPalavra(t *testing.T) {
	w := []PalavraTempo{pt("chegou", 0, .5), pt("a", .5, .7), pt("82", .7, 1.0), pt(",9", 1.0, 1.3), pt("%", 1.3, 1.8), pt("do", 1.8, 2.0), pt("PIB", 2.0, 2.5)}
	got := AlinharLegendas(scriptDe("chegou a 82,9% do PIB"), tempoDe(linhaCom(0, 0, 2.5, w...))).Linhas[0].Palavras
	if !reflect.DeepEqual(textosDe(got), []string{"chegou", "a", "82,9%", "do", "PIB"}) {
		t.Fatalf("palavras: %v", textosDe(got))
	}
	if got[2].InicioS != .7 || got[2].FimS != 1.8 || got[3].InicioS != 1.8 {
		t.Fatalf("82,9%% deve cobrir 82, ,9 e o %% falado (0,7 a 1,8): %+v", got[2:4])
	}
}

func TestLegendasTokenSemLetrasGrudaNoVizinho(t *testing.T) {
	w := seq([]string{"um", "dois"}, 1)
	got := AlinharLegendas(scriptDe("um — dois"), tempoDe(linhaCom(0, 0, 2, w...))).Linhas[0].Palavras
	if !reflect.DeepEqual(textosDe(got), []string{"um —", "dois"}) {
		t.Fatalf("travessao vai para a palavra anterior: %v", textosDe(got))
	}
	w = seq([]string{"começa", "agora"}, 1)
	got = AlinharLegendas(scriptDe("— começa agora"), tempoDe(linhaCom(0, 0, 2, w...))).Linhas[0].Palavras
	if !reflect.DeepEqual(textosDe(got), []string{"— começa", "agora"}) {
		t.Fatalf("travessao inicial vai para a proxima palavra: %v", textosDe(got))
	}
}

func TestLegendasInterpolaPalavraNaoOuvida(t *testing.T) {
	w := []PalavraTempo{pt("um", 0, 1), pt("tres", 2, 3)}
	got := AlinharLegendas(scriptDe("um dois tres"), tempoDe(linhaCom(0, 0, 3, w...))).Linhas[0].Palavras
	if len(got) != 3 || got[1].Texto != "dois" || got[1].InicioS != 1 || got[1].FimS != 2 {
		t.Fatalf("palavra nao ouvida ocupa a lacuna entre as vizinhas: %+v", got)
	}
}

func TestLegendasSemWhisperDistribuiNaJanela(t *testing.T) {
	got := AlinharLegendas(scriptDe("a bb cc"), tempoDe(linhaCom(0, 10, 12))).Linhas[0].Palavras
	if len(got) != 3 || got[0].InicioS != 10 || got[2].FimS != 12 {
		t.Fatalf("3 palavras devem cobrir a janela 10 a 12: %+v", got)
	}
	for i := 0; i < 2; i++ {
		if got[i].FimS != got[i+1].InicioS {
			t.Fatalf("buraco entre palavras: %+v", got)
		}
	}
}

func TestLegendasPreservaContagemOrdemEJanela(t *testing.T) {
	casos := []struct {
		texto   string
		whisper []string
	}{
		{"um dois tres quatro", []string{"a", "b", "c", "d", "e", "f", "g", "h"}},
		{"um dois tres quatro cinco seis sete oito", []string{"x", "y", "z"}},
		{"um dois tres", []string{"um", "dois", "tres"}},
		{"única", nil},
	}
	for i, c := range casos {
		l := linhaCom(0, 2, 8, seq(c.whisper, 0.5)...)
		for k := range l.Palavras {
			l.Palavras[k].InicioS += 2
			l.Palavras[k].FimS += 2
		}
		got := AlinharLegendas(scriptDe(c.texto), tempoDe(l)).Linhas[0].Palavras
		if len(got) != len(strings.Fields(c.texto)) {
			t.Fatalf("caso %d: contagem %d != %d", i, len(got), len(strings.Fields(c.texto)))
		}
		cursor := 2.0
		for _, p := range got {
			if p.InicioS < cursor || p.FimS < p.InicioS || p.FimS > 8 {
				t.Fatalf("caso %d: palavra fora da janela ou sobreposta: %+v", i, p)
			}
			cursor = p.FimS
		}
	}
}

func TestLegendasNaoAlteraEntradaEIgnoraLinhasSobrando(t *testing.T) {
	tm := tempoDe(linhaCom(0, 0, 1, pt("x", 0, 1)), linhaCom(1, 1, 2, pt("y", 1, 2)))
	saida := AlinharLegendas(scriptDe("alfa"), tm)
	if tm.Linhas[0].Palavras[0].Texto != "x" {
		t.Fatal("a entrada nao pode ser alterada")
	}
	if saida.Linhas[0].Palavras[0].Texto != "alfa" || saida.Linhas[1].Palavras[0].Texto != "y" {
		t.Fatalf("linha sem roteiro fica como esta: %+v", saida.Linhas)
	}
	if got := AlinharLegendas(scriptDe("a", "b", "c"), tm); len(got.Linhas) != 2 {
		t.Fatalf("roteiro com linhas a mais nao cria linhas: %d", len(got.Linhas))
	}
}

func timelineComLegendas(t *testing.T) (Timeline, Timings) {
	t.Helper()
	textos := make([]string, 12)
	for i := range textos {
		textos[i] = fmt.Sprintf("w%d", i)
	}
	tm := tempoDe(linhaCom(0, 0, 9.75, seq(textos, 0.8125)...))
	tmL := AlinharLegendas(scriptDe(strings.Join(textos, " ")), tm)
	tl := SegmentarBatidas(timelineDe(batidaUnica(0, 0, 9.75, "pan")), tmL)
	return AplicarLegendas(tl, tmL), tmL
}

func legendaDe(b Batida) CamadaItem {
	for _, c := range b.Camadas {
		if c.Camada == "legenda" {
			return c
		}
	}
	return CamadaItem{}
}

func TestLegendasKaraokeExatoEUmaVezPorPalavra(t *testing.T) {
	tl, tmL := timelineComLegendas(t)
	if err := tl.Validar(); err != nil {
		t.Fatalf("timeline com legendas deve validar: %v", err)
	}
	var todas []PalavraTempo
	for _, b := range tl.Batidas {
		ps := legendaDe(b).Palavras
		if len(ps) != 4 {
			t.Fatalf("12 palavras em 3 batidas: 4 por batida, obteve %d", len(ps))
		}
		todas = append(todas, ps...)
	}
	if !reflect.DeepEqual(todas, tmL.Linhas[0].Palavras) {
		t.Fatalf("cada palavra realca exatamente no tempo alinhado, uma vez so:\n%+v\n%+v", todas, tmL.Linhas[0].Palavras)
	}
}

func TestLegendasAplicarRecortaPalavraQueCruzaACortada(t *testing.T) {
	tmL := tempoDe(linhaCom(0, 0, 6, pt("t0", 0, 1), pt("t1", 2.5, 3.5), pt("t2", 3.5, 5)))
	tl := timelineDe(batidaUnica(0, 0, 3, "pan"), batidaUnica(0, 3, 6, "pan"))
	saida := AplicarLegendas(tl, tmL)
	a, b := legendaDe(saida.Batidas[0]).Palavras, legendaDe(saida.Batidas[1]).Palavras
	if len(a) != 2 || a[1].Texto != "t1" || a[1].FimS != 3 || len(b) != 1 || b[0].Texto != "t2" {
		t.Fatalf("palavra que comeca antes do corte fica na primeira batida, recortada: %+v | %+v", a, b)
	}
	if legendaDe(tl.Batidas[0]).Palavras != nil {
		t.Fatal("a entrada nao pode ser alterada")
	}
}

func TestSchemasLegendaPalavras(t *testing.T) {
	muta := map[string]func(*CamadaItem){
		"antes da batida":     func(c *CamadaItem) { c.Palavras = []PalavraTempo{pt("a", -1, 1)} },
		"depois da batida":    func(c *CamadaItem) { c.Palavras = []PalavraTempo{pt("a", 1, 9)} },
		"fora de ordem":       func(c *CamadaItem) { c.Palavras = []PalavraTempo{pt("a", 2, 3), pt("b", 1, 2)} },
		"sem texto":           func(c *CamadaItem) { c.Palavras = []PalavraTempo{pt(" ", 1, 2)} },
		"fim antes do inicio": func(c *CamadaItem) { c.Palavras = []PalavraTempo{pt("a", 2, 1)} },
	}
	for nome, f := range muta {
		t.Run(nome, func(t *testing.T) {
			b := batidaUnica(0, 0, 3, "pan")
			f(&b.Camadas[2])
			if err := timelineDe(b).Validar(); !errors.Is(err, ErrSchemaInvalido) {
				t.Fatalf("esperava ErrSchemaInvalido, obteve %v", err)
			}
		})
	}
}

func blocosSRT(t *testing.T, srt string) [][]string {
	t.Helper()
	var blocos [][]string
	for _, b := range strings.Split(strings.TrimSpace(srt), "\n\n") {
		blocos = append(blocos, strings.Split(b, "\n"))
	}
	return blocos
}

func TestLegendasSRT(t *testing.T) {
	tl, _ := timelineComLegendas(t)
	srt, err := GerarSRT(tl)
	if err != nil {
		t.Fatal(err)
	}
	blocos := blocosSRT(t, srt)
	esperado := [][]string{
		{"1", "00:00:00,000 --> 00:00:03,250", "w0 w1 w2 w3"},
		{"2", "00:00:03,250 --> 00:00:06,500", "w4 w5 w6 w7"},
		{"3", "00:00:06,500 --> 00:00:09,750", "w8 w9 w10 w11"},
	}
	if !reflect.DeepEqual(blocos, esperado) {
		t.Fatalf("srt inesperado:\n%s", srt)
	}
	if !strings.HasSuffix(srt, "\n") {
		t.Fatal("srt deve terminar com quebra de linha")
	}
}

func TestLegendasSRTQuebraLinhasLongas(t *testing.T) {
	var ps []PalavraTempo
	for i := 0; i < 8; i++ {
		ps = append(ps, pt(fmt.Sprintf("palavra%d", i), float64(i)*0.3, float64(i+1)*0.3))
	}
	b := batidaUnica(0, 0, 2.4, "pan")
	b.Camadas[2].Palavras = ps
	srt, err := GerarSRT(timelineDe(b))
	if err != nil {
		t.Fatal(err)
	}
	texto := blocosSRT(t, srt)[0][2:]
	if len(texto) < 2 {
		t.Fatalf("71 caracteres devem quebrar em mais de uma linha: %v", texto)
	}
	for _, l := range texto {
		if len([]rune(l)) > LarguraSRT {
			t.Fatalf("linha maior que %d: %q", LarguraSRT, l)
		}
	}
	if strings.Join(texto, " ") != "palavra0 palavra1 palavra2 palavra3 palavra4 palavra5 palavra6 palavra7" {
		t.Fatalf("quebra nao pode perder nem reordenar palavras: %v", texto)
	}
}

func TestLegendasSRTErros(t *testing.T) {
	if _, err := GerarSRT(timelineDe(batidaUnica(0, 0, 3, "pan"))); !errors.Is(err, ErrSchemaInvalido) {
		t.Fatalf("timeline sem nenhuma legenda deve falhar: %v", err)
	}
	a, b := batidaUnica(0, 0, 3, "pan"), batidaUnica(0, 3, 6, "pan")
	a.Camadas[2].Palavras = []PalavraTempo{pt("a", 0, 3)}
	b.Camadas[2].Palavras = []PalavraTempo{pt("b", 2, 6)}
	if _, err := GerarSRT(timelineDe(a, b)); !errors.Is(err, ErrSchemaInvalido) {
		t.Fatalf("blocos sobrepostos devem falhar: %v", err)
	}
}

func TestLegendasTempoSRT(t *testing.T) {
	casos := map[float64]string{0: "00:00:00,000", 3.25: "00:00:03,250", 3661.5: "01:01:01,500", 59.9996: "00:01:00,000", -1: "00:00:00,000"}
	for s, esperado := range casos {
		if got := formatarTempoSRT(s); got != esperado {
			t.Errorf("%v s: obteve %s, quer %s", s, got, esperado)
		}
	}
}

func TestLegendasQuebraBalanceada(t *testing.T) {
	palavras := strings.Fields("a Dívida Bruta do Governo Geral chegou a 82,9%")
	got := quebrarLinhas(palavras, LarguraSRT)
	if len(got) != 2 || len([]rune(got[1])) < 10 {
		t.Fatalf("duas linhas equilibradas, sem palavra solta na segunda: %q", got)
	}
	for _, l := range got {
		if len([]rune(l)) > LarguraSRT {
			t.Fatalf("linha maior que %d: %q", LarguraSRT, l)
		}
	}
	if strings.Join(got, " ") != strings.Join(palavras, " ") {
		t.Fatalf("quebra nao pode perder nem reordenar palavras: %q", got)
	}
	if curta := quebrarLinhas([]string{"oi", "mundo"}, LarguraSRT); len(curta) != 1 || curta[0] != "oi mundo" {
		t.Fatalf("texto curto fica em uma linha: %q", curta)
	}
}

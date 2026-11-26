package render

import (
	"encoding/json"
	"math"
	"strings"
	"testing"
)

func nomesVisemas(seq []visemaPeso) string {
	n := make([]string, len(seq))
	for i, v := range seq {
		n[i] = v.nome
	}
	return strings.Join(n, " ")
}

func TestVisemasSequencia(t *testing.T) {
	casos := map[string]string{
		"casa":    "CONS A CONS A",
		"bola":    "MBP O CONS A",
		"atingiu": "A CONS I CONS I U",
		"milhão":  "MBP I CONS A O",
		"queijo":  "CONS E I CONS O",
		"hoje":    "O CONS E",
		"PIB":     "MBP I MBP",
		"ação":    "A CONS A O",
		"efeito":  "E FV E I CONS O",
		"82":      "CONS",
		"%":       "CONS",
		",9":      "CONS",
		"h":       "CONS",
		"R$":      "CONS",
		"passo":   "MBP A CONS O",
	}
	for palavra, esperado := range casos {
		if got := nomesVisemas(sequenciaVisemas(palavra)); got != esperado {
			t.Errorf("%q: obteve %q, quer %q", palavra, got, esperado)
		}
	}
}

func TestVisemasPesoDaVogalMaior(t *testing.T) {
	seq := sequenciaVisemas("ba")
	if seq[0].peso != 1 || seq[1].peso != 2 {
		t.Fatalf("vogal deve pesar 2 e consoante 1: %+v", seq)
	}
	seq = sequenciaVisemas("assa")
	if nomesVisemas(seq) != "A CONS A" || seq[1].peso != 2 {
		t.Fatalf("consoantes iguais adjacentes devem fundir somando o peso: %+v", seq)
	}
}

func timingsParaVisemas() Timings {
	return Timings{SchemaVersion: 1, Provedor: "p", Voz: "v", Linhas: []LinhaTempo{
		{IndiceLinha: 0, InicioS: 0, FimS: 2,
			Palavras: []PalavraTempo{{Texto: "ba", InicioS: 0, FimS: 0.3}, {Texto: "conta", InicioS: 0.35, FimS: 1.2}},
			Visemas:  []VisemaTempo{{Viseme: "antigo", InicioS: 0, FimS: 0.1}}},
		{IndiceLinha: 1, InicioS: 2, FimS: 4,
			Palavras: []PalavraTempo{{Texto: "publica", InicioS: 2, FimS: 3}}},
	}}
}

func TestVisemasGerarCobrePalavrasEmOrdem(t *testing.T) {
	entrada := timingsParaVisemas()
	saida := GerarVisemas(entrada)
	if err := saida.Validar(); err != nil {
		t.Fatalf("timings com visemas deve validar: %v", err)
	}
	for _, l := range saida.Linhas {
		for _, p := range l.Palavras {
			dentro := 0
			for _, v := range l.Visemas {
				if v.InicioS >= p.InicioS && v.FimS <= p.FimS && v.Viseme != VisemaRepouso {
					dentro++
				}
			}
			if dentro == 0 {
				t.Fatalf("palavra %q sem visemas dentro de [%v,%v]", p.Texto, p.InicioS, p.FimS)
			}
		}
		for i := 1; i < len(l.Visemas); i++ {
			if l.Visemas[i].InicioS < l.Visemas[i-1].FimS {
				t.Fatalf("visemas sobrepostos na linha %d: %+v %+v", l.IndiceLinha, l.Visemas[i-1], l.Visemas[i])
			}
		}
	}
	if entrada.Linhas[0].Visemas[0].Viseme != "antigo" {
		t.Fatal("a entrada nao pode ser alterada")
	}
}

func TestVisemasDuracaoProporcionalEExataNosLimites(t *testing.T) {
	l := GerarVisemas(timingsParaVisemas()).Linhas[0]
	ba := l.Visemas[:2]
	if ba[0].Viseme != "MBP" || ba[1].Viseme != "A" {
		t.Fatalf("visemas de ba: %+v", ba)
	}
	if ba[0].InicioS != 0 || ba[1].FimS != 0.3 || ba[0].FimS != ba[1].InicioS {
		t.Fatalf("limites devem coincidir exatamente com a palavra e entre si: %+v", ba)
	}
	if math.Abs((ba[0].FimS-ba[0].InicioS)-0.1) > 1e-9 || math.Abs((ba[1].FimS-ba[1].InicioS)-0.2) > 1e-9 {
		t.Fatalf("vogal deve durar o dobro da consoante: %+v", ba)
	}
}

func TestVisemasRepousoNasLacunas(t *testing.T) {
	l := GerarVisemas(timingsParaVisemas()).Linhas[0]
	repousos := 0
	for _, v := range l.Visemas {
		if v.Viseme == VisemaRepouso {
			repousos++
		}
	}
	if repousos != 1 {
		t.Fatalf("so a cauda de 0,8 s gera repouso (o intervalo de 0,05 s nao): obteve %d: %+v", repousos, l.Visemas)
	}
}

func TestVisemasLacunaPequenaNaoGeraRepouso(t *testing.T) {
	tm := Timings{SchemaVersion: 1, Provedor: "p", Voz: "v", Linhas: []LinhaTempo{{IndiceLinha: 0, InicioS: 0, FimS: 1,
		Palavras: []PalavraTempo{{Texto: "a", InicioS: 0, FimS: 0.5}, {Texto: "b", InicioS: 0.55, FimS: 1}}}}}
	for _, v := range GerarVisemas(tm).Linhas[0].Visemas {
		if v.Viseme == VisemaRepouso {
			t.Fatalf("lacuna de 0,05 s nao deve gerar repouso: %+v", v)
		}
	}
}

func TestVisemasPalavraSemDuracaoEPalavraSemLetras(t *testing.T) {
	tm := Timings{SchemaVersion: 1, Provedor: "p", Voz: "v", Linhas: []LinhaTempo{{IndiceLinha: 0, InicioS: 0, FimS: 1,
		Palavras: []PalavraTempo{{Texto: "oi", InicioS: 0.2, FimS: 0.2}, {Texto: "82", InicioS: 0.2, FimS: 1}}}}}
	saida := GerarVisemas(tm)
	if err := saida.Validar(); err != nil {
		t.Fatalf("palavra sem duracao deve gerar visemas validos: %v", err)
	}
	if len(saida.Linhas[0].Visemas) < 3 {
		t.Fatalf("cada palavra recebe ao menos um viseme: %+v", saida.Linhas[0].Visemas)
	}
}

func TestVisemasReproduzivelByteAByte(t *testing.T) {
	a, _ := json.Marshal(GerarVisemas(timingsParaVisemas()))
	b, _ := json.Marshal(GerarVisemas(timingsParaVisemas()))
	if string(a) != string(b) {
		t.Fatal("duas execucoes devem gerar JSON identico")
	}
}

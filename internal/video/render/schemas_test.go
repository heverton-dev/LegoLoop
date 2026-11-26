package render

import (
	"errors"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/video/ativos"
)

func timingsValido() Timings {
	return Timings{SchemaVersion: 1, Provedor: "edge-tts", Voz: "pt-BR-ThalitaMultilingualNeural", Linhas: []LinhaTempo{
		{IndiceLinha: 0, InicioS: 0, FimS: 2, WER: 0.02,
			Palavras: []PalavraTempo{{Texto: "a", InicioS: 0, FimS: 0.5}, {Texto: "conta", InicioS: 0.5, FimS: 1.2}},
			Visemas:  []VisemaTempo{{Viseme: "A", InicioS: 0, FimS: 0.5}}},
		{IndiceLinha: 1, InicioS: 2, FimS: 4, WER: 0,
			Palavras: []PalavraTempo{{Texto: "publica", InicioS: 2, FimS: 3}}},
	}}
}

func timelineValida() Timeline {
	return Timeline{SchemaVersion: 1, DuracaoS: 6, Batidas: []Batida{
		{IndiceLinha: 0, InicioS: 0, FimS: 3, CenarioID: "cen_01", Camera: "zoom",
			Camadas: []CamadaItem{{ID: "fundo", Camada: "cenario", Profundidade: 0}, {ID: "nar", Camada: "sprite", Profundidade: 2, SpriteID: "narrador_pose_1", Caixa: ativos.BoundingBox{X: 1, Y: 1, W: 10, H: 10}}},
			SFX:     []EventoSFX{{Tipo: "troca_cena", TempoS: 0, Arquivo: "sfx/whoosh.wav"}}},
		{IndiceLinha: 0, InicioS: 3, FimS: 6, CenarioID: "cen_01", Camera: "pan"},
	}}
}

func TestSchemasConfigRender(t *testing.T) {
	padrao := ConfigPadrao()
	if len(padrao.Resolucoes) != 1 || padrao.Resolucoes[0] != "1080p" || padrao.ReprovasWERMax != 3 {
		t.Fatalf("padrao inesperado: %+v", padrao)
	}
	if err := padrao.Validar(); err != nil {
		t.Fatalf("padrao deve validar: %v", err)
	}
	multi := ConfigRender{SchemaVersion: 1, Resolucoes: []string{"720p", "1080p", "4k"}, ReprovasWERMax: 2}
	if err := multi.Validar(); err != nil {
		t.Fatalf("720p, 1080p e 4k devem validar: %v", err)
	}
	invalidos := []ConfigRender{
		{SchemaVersion: 1, Resolucoes: []string{"480p"}, ReprovasWERMax: 3},
		{SchemaVersion: 1, Resolucoes: []string{"1080p", "1080p"}, ReprovasWERMax: 3},
		{SchemaVersion: 1, Resolucoes: []string{"1080p"}, ReprovasWERMax: -1},
		{SchemaVersion: 0, Resolucoes: []string{"1080p"}, ReprovasWERMax: 3},
	}
	for i, c := range invalidos {
		if err := c.Validar(); !errors.Is(err, ErrSchemaInvalido) {
			t.Errorf("caso %d: esperava ErrSchemaInvalido, obteve %v", i, err)
		}
	}
}

func TestSchemasConfigNormalizar(t *testing.T) {
	c := ConfigRender{SchemaVersion: 1}.Normalizar()
	if len(c.Resolucoes) != 1 || c.Resolucoes[0] != "1080p" || c.ReprovasWERMax != 3 {
		t.Fatalf("normalizacao sem padroes: %+v", c)
	}
	c = ConfigRender{SchemaVersion: 1, Resolucoes: []string{"720p"}, ReprovasWERMax: 1}.Normalizar()
	if c.Resolucoes[0] != "720p" || c.ReprovasWERMax != 1 {
		t.Fatalf("normalizacao nao deve sobrescrever valores: %+v", c)
	}
}

func TestSchemasTimings(t *testing.T) {
	if err := timingsValido().Validar(); err != nil {
		t.Fatalf("timings valido reprovado: %v", err)
	}
	muta := map[string]func(*Timings){
		"sem provedor":          func(x *Timings) { x.Provedor = "" },
		"sem linhas":            func(x *Timings) { x.Linhas = nil },
		"indice fora de ordem":  func(x *Timings) { x.Linhas[1].IndiceLinha = 5 },
		"inicio >= fim":         func(x *Timings) { x.Linhas[0].FimS = 0 },
		"wer acima de 1":        func(x *Timings) { x.Linhas[0].WER = 1.5 },
		"palavra fora da linha": func(x *Timings) { x.Linhas[0].Palavras[1].FimS = 9 },
		"palavras sobrepostas":  func(x *Timings) { x.Linhas[0].Palavras[1].InicioS = 0.3 },
		"palavra vazia":         func(x *Timings) { x.Linhas[0].Palavras[0].Texto = "" },
		"linha sem palavras":    func(x *Timings) { x.Linhas[1].Palavras = nil },
		"viseme fora da linha":  func(x *Timings) { x.Linhas[0].Visemas[0].FimS = 7 },
		"linhas sobrepostas":    func(x *Timings) { x.Linhas[1].InicioS = 1 },
		"wer global acima de 1": func(x *Timings) { x.WERGlobal = 2 },
		"wer global negativo":   func(x *Timings) { x.WERGlobal = -0.1 },
	}
	for nome, f := range muta {
		t.Run(nome, func(t *testing.T) {
			x := timingsValido()
			f(&x)
			if err := x.Validar(); !errors.Is(err, ErrSchemaInvalido) {
				t.Fatalf("esperava ErrSchemaInvalido, obteve %v", err)
			}
		})
	}
}

func TestSchemasTimeline(t *testing.T) {
	if err := timelineValida().Validar(); err != nil {
		t.Fatalf("timeline valida reprovada: %v", err)
	}
	muta := map[string]func(*Timeline){
		"sem batidas":            func(x *Timeline) { x.Batidas = nil },
		"duracao zero":           func(x *Timeline) { x.DuracaoS = 0 },
		"batida inicio>=fim":     func(x *Timeline) { x.Batidas[0].FimS = 0 },
		"batidas sobrepostas":    func(x *Timeline) { x.Batidas[1].InicioS = 2 },
		"batida alem da duracao": func(x *Timeline) { x.Batidas[1].FimS = 99 },
		"camera invalida":        func(x *Timeline) { x.Batidas[0].Camera = "giro" },
		"sem cenario":            func(x *Timeline) { x.Batidas[0].CenarioID = "" },
		"camada invalida":        func(x *Timeline) { x.Batidas[0].Camadas[0].Camada = "xyz" },
		"profundidade igual":     func(x *Timeline) { x.Batidas[0].Camadas[1].Profundidade = 0 },
		"sfx tipo invalido":      func(x *Timeline) { x.Batidas[0].SFX[0].Tipo = "explosao" },
		"sfx fora da batida":     func(x *Timeline) { x.Batidas[0].SFX[0].TempoS = 4 },
		"sfx sem arquivo":        func(x *Timeline) { x.Batidas[0].SFX[0].Arquivo = "" },
	}
	for nome, f := range muta {
		t.Run(nome, func(t *testing.T) {
			x := timelineValida()
			f(&x)
			if err := x.Validar(); !errors.Is(err, ErrSchemaInvalido) {
				t.Fatalf("esperava ErrSchemaInvalido, obteve %v", err)
			}
		})
	}
}

func TestSchemasRelatorio(t *testing.T) {
	ok := RelatorioRender{Encoder: "libx264", ProvedorVoz: "edge-tts", LUFS: -14, TruePeakDBTP: -1.6}
	if err := ok.Validar(); err != nil {
		t.Fatalf("relatorio valido reprovado: %v", err)
	}
	for i, r := range []RelatorioRender{{ProvedorVoz: "x"}, {Encoder: "x"}} {
		if err := r.Validar(); !errors.Is(err, ErrSchemaInvalido) {
			t.Errorf("caso %d: esperava ErrSchemaInvalido, obteve %v", i, err)
		}
	}
}

func TestSchemasErrosDistintos(t *testing.T) {
	lista := []error{ErrPreVooBloqueado, ErrProvedorDesconhecido, ErrCreditosNegados, ErrWERAcimaLimite,
		ErrTimelineRitmo, ErrLoudnessReprovado, ErrRenderInvalido, ErrCorteVertical, ErrSchemaInvalido}
	vistos := map[string]bool{}
	for _, e := range lista {
		if e == nil || vistos[e.Error()] {
			t.Fatalf("erro nulo ou duplicado: %v", e)
		}
		vistos[e.Error()] = true
	}
}

package motor

import (
	"encoding/json"
	"math"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/video/ativos"
	"github.com/heverton-dev/LegoLoop/internal/video/inventario"
	"github.com/heverton-dev/LegoLoop/internal/video/render"
)

func TestPonteEmiteLayersCartaoCameraESubcortes(t *testing.T) {
	dir := t.TempDir()
	e := Entrada{
		Cenas: ativos.Storyboard{
			SchemaVersion: 2,
			AlvoCorteS:    3.2,
			Musica:        "dark",
			Cenas: []ativos.Cena{
				{
					IndiceLinha:       0,
					CenarioID:         "escritorio",
					CenarioID2:        "servidor",
					ParametrosCenario: map[string]any{"calendar": "NOV 2026"},
					TextoTela:         "90 HORAS",
					CorTexto:          "#FF5A5A",
					TextoTela2:        "META DOBRADA",
					CorTexto2:         "#FFD24A",
					Camera:            "in",
					FX:                []string{"dust", "bokeh"},
					Camadas:           []map[string]any{{"kind": "stopwatch", "x": 1620.0, "y": 800.0, "s": 0.55, "in": "pop", "hud": true}},
					ElementosMoveis: []ativos.ElementoMovel{
						{ID: "p0", Tipo: "personagem", PersonagemID: "dev", X: 600, Y: 980, Escala: 1.1, Entrada: "left", Pose: "point", Humor: "worry", Fala: true},
					},
				},
				{
					IndiceLinha: 1,
					Cartao:      true,
					Texto:       "Por que a promessa virou sobrecarga?",
				},
			},
		},
		Elenco: ElencoPonte{Personagens: []PersonagemPonte{
			{ID: "dev", Cast: map[string]any{"skin": "#E6BE8E", "suit": "#1F2F5A", "glasses": true}},
		}},
		Timings: render.Timings{SchemaVersion: 1, Linhas: []render.LinhaTempo{
			{IndiceLinha: 0, InicioS: 0.0, FimS: 6.4},
			{IndiceLinha: 1, InicioS: 6.4, FimS: 9.6},
		}},
		Fundos: []inventario.Item{
			{ID: "escritorio", Tipo: "cenario", Parametros: map[string]any{"kind": "office", "time": "night"}},
			{ID: "servidor", Tipo: "cenario", Parametros: map[string]any{"kind": "server", "time": "night"}},
		},
		Estilo:      inventario.Estilo{ID: "molde-de-pilula", Luz: &inventario.Luz{Rim: true, Cor: "#FFB347"}},
		DuracaoVozS: 9.6,
		FPS:         30,
	}
	tl, err := Montar(e, dir)
	if err != nil {
		t.Fatalf("Montar diorama rico: %v", err)
	}
	if len(tl.Shots) != 3 {
		t.Fatalf("esperava 3 cortes (2 subcortes da linha 0 + 1 cartao da linha 1), obteve %d", len(tl.Shots))
	}
	s0, s1, s2 := tl.Shots[0], tl.Shots[1], tl.Shots[2]
	if s0.Set["calendar"] != "NOV 2026" || s1.Set["kind"] != "server" {
		t.Fatalf("parametros de set ou cenario do subcorte B incorretos: s0=%+v s1=%+v", s0.Set, s1.Set)
	}
	if len(s0.Layers) < 2 || len(s0.Cam) != 2 || len(s0.FX) != 2 {
		t.Fatalf("corte A sem layers/cam/fx: layers=%+v cam=%+v fx=%+v", s0.Layers, s0.Cam, s0.FX)
	}
	if s0.Chars[0].X != 600 || s0.Chars[0].Enter != "left" {
		t.Fatalf("boneco nao respeitou X/Enter: %+v", s0.Chars[0])
	}
	if s2.Type != "card" || s2.Text != "Por que a promessa virou sobrecarga?" || s2.CPS <= 0 {
		t.Fatalf("cartao preto digitado incorreto: %+v", s2)
	}
	soma := s0.Duration + s1.Duration + s2.Duration
	if math.Abs(soma-e.DuracaoVozS) > 1.0/30 {
		t.Fatalf("soma %.3f diverge da voz %.3f", soma, e.DuracaoVozS)
	}
	raw, _ := json.Marshal(tl)
	if !strings.Contains(string(raw), `"90 HORAS"`) || !strings.Contains(string(raw), `"META DOBRADA"`) {
		t.Fatalf("textos foco ausentes no JSON: %s", raw)
	}
}

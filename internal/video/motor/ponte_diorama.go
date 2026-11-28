package motor

import (
	"fmt"
	"maps"
	"math"
	"unicode/utf8"

	"github.com/heverton-dev/LegoLoop/internal/video/ativos"
)

func montarCortesDaCena(e Entrada, dir string, c ativos.Cena, cast map[string]map[string]any, quadros int) ([]Corte, error) {
	durTotal := float64(quadros) / float64(e.FPS)
	if c.Cartao {
		return []Corte{montarCartao(c, durTotal)}, nil
	}
	n := calcularSubcortes(durTotal, e.Cenas.AlvoCorteS)
	if n == 1 {
		ct, err := montarSubcorteCena(e, dir, c, cast, 0, 1, durTotal)
		if err != nil {
			return nil, err
		}
		return []Corte{ct}, nil
	}
	return dividirCenaEmSubcortes(e, dir, c, cast, quadros, n)
}

func montarCartao(c ativos.Cena, dur float64) Corte {
	txt := c.TextoCartao
	if txt == "" {
		txt = c.Texto
	}
	runas := float64(utf8.RuneCountInString(txt))
	cps := math.Round((runas/math.Max(0.8, (dur-0.3)*0.55))*10) / 10
	return Corte{
		ID:       fmt.Sprintf("l%03d", c.IndiceLinha),
		Type:     "card",
		Lines:    [2]int{c.IndiceLinha, c.IndiceLinha},
		Text:     txt,
		CPS:      cps,
		Duration: dur,
		Chars:    []Boneco{},
		Layers:   []map[string]any{},
	}
}

func calcularSubcortes(dur, alvo float64) int {
	if alvo <= 0 {
		return 1
	}
	n := int(math.Round(dur / alvo))
	if n < 1 {
		return 1
	}
	return min(n, 4)
}

func dividirCenaEmSubcortes(e Entrada, dir string, c ativos.Cena, cast map[string]map[string]any, quadros, n int) ([]Corte, error) {
	saida := make([]Corte, 0, n)
	acum := 0
	for v := 0; v < n; v++ {
		fimQ := quadros * (v + 1) / n
		durV := float64(fimQ-acum) / float64(e.FPS)
		acum = fimQ
		ct, err := montarSubcorteCena(e, dir, c, cast, v, n, durV)
		if err != nil {
			return nil, err
		}
		saida = append(saida, ct)
	}
	return saida, nil
}

func montarSubcorteCena(e Entrada, dir string, c ativos.Cena, cast map[string]map[string]any, v, total int, dur float64) (Corte, error) {
	cenID := escolherCenarioSubcorte(c, v)
	bg, setBase, err := resolverFundo(e.Fundos, dir, cenID)
	if err != nil {
		return Corte{}, err
	}
	set := mesclarSet(setBase, c.ParametrosCenario, v)
	bonecos, err := montarBonecosSubcorte(c.ElementosMoveis, cast, v)
	if err != nil {
		return Corte{}, err
	}
	id := fmt.Sprintf("l%03d", c.IndiceLinha)
	if total > 1 {
		id = fmt.Sprintf("l%03d%c", c.IndiceLinha, 'a'+rune(v))
	}
	return construirStructSubcorte(id, e, c, bg, set, bonecos, v, dur), nil
}

func construirStructSubcorte(id string, e Entrada, c ativos.Cena, bg string, set map[string]any, bonecos []Boneco, v int, dur float64) Corte {
	camNome := alternarCameraSubcorte(c.Camera, v)
	return Corte{
		ID:       id,
		Type:     "scene",
		Lines:    [2]int{c.IndiceLinha, c.IndiceLinha},
		BG:       bg,
		Set:      set,
		Duration: dur,
		Chars:    bonecos,
		Layers:   montarLayersSubcorte(c, v),
		Cam:      quadrosCamera(camNome, dur),
		FX:       append([]string(nil), c.FX...),
		Shake:    copiarMapas(c.Shake),
		Grade:    c.Grade,
		Tin:      c.TransicaoEntrada,
		Tout:     c.TransicaoSaida,
		Light:    luzDoEstilo(e.Estilo),
	}
}

func escolherCenarioSubcorte(c ativos.Cena, v int) string {
	if v > 0 && c.CenarioID2 != "" {
		return c.CenarioID2
	}
	return c.CenarioID
}

func mesclarSet(base, extra map[string]any, v int) map[string]any {
	if len(base) == 0 && len(extra) == 0 {
		return nil
	}
	out := map[string]any{}
	maps.Copy(out, base)
	if v == 0 {
		maps.Copy(out, extra)
	}
	return out
}

func montarBonecosSubcorte(elementos []ativos.ElementoMovel, cast map[string]map[string]any, v int) ([]Boneco, error) {
	bs, err := montarBonecos(elementos, cast)
	if err != nil {
		return nil, err
	}
	if v > 0 {
		for i := range bs {
			bs[i].Enter = ""
		}
	}
	return bs, nil
}

func montarLayersSubcorte(c ativos.Cena, v int) []map[string]any {
	var ls []map[string]any
	if v == 0 {
		if c.TextoTela != "" {
			ls = append(ls, camadaTextoFoco(c.TextoTela, c.CorTexto, c.TextoTelaX, c.TextoTelaY))
		}
		ls = append(ls, copiarMapas(c.Camadas)...)
	} else if v == 1 {
		if c.TextoTela2 != "" {
			ls = append(ls, camadaTextoFoco(c.TextoTela2, c.CorTexto2, 960, 220))
		}
		ls = append(ls, copiarMapas(c.Camadas2)...)
	}
	return ls
}

func camadaTextoFoco(texto, cor string, x, y float64) map[string]any {
	if cor == "" {
		cor = "#FFD24A"
	}
	if x <= 0 {
		x = 960
	}
	if y <= 0 {
		y = 210
	}
	n := utf8.RuneCountInString(texto)
	tam := 104.0
	switch {
	case n > 26:
		tam = 62
	case n > 18:
		tam = 76
	case n > 12:
		tam = 90
	}
	return map[string]any{
		"kind": "text", "text": texto, "x": x, "y": y, "size": tam,
		"font": "brush", "color": cor, "stroke": "#18161B", "sw": 10.0,
		"start": 0.35, "in": "pop", "hud": true,
	}
}

func alternarCameraSubcorte(base string, v int) string {
	if base == "" {
		return ""
	}
	ciclo := []string{"in", "panR", "out", "low", "panL"}
	if v == 0 {
		return base
	}
	for i, c := range ciclo {
		if c == base {
			return ciclo[(i+v)%len(ciclo)]
		}
	}
	return ciclo[v%len(ciclo)]
}

func quadrosCamera(kind string, dur float64) []QuadroCamera {
	presets := map[string][2][3]float64{
		"in":   {{960, 540, 1.0}, {960, 560, 1.18}},
		"out":  {{960, 560, 1.2}, {960, 540, 1.0}},
		"panL": {{1060, 540, 1.08}, {860, 550, 1.1}},
		"panR": {{860, 540, 1.08}, {1060, 550, 1.1}},
		"low":  {{960, 540, 1.0}, {960, 520, 1.06}},
	}
	ab, ok := presets[kind]
	if !ok {
		return nil
	}
	fim := math.Round(dur*100) / 100
	return []QuadroCamera{
		{T: 0, X: ab[0][0], Y: ab[0][1], Z: ab[0][2]},
		{T: fim, X: ab[1][0], Y: ab[1][1], Z: ab[1][2]},
	}
}

func copiarMapas(orig []map[string]any) []map[string]any {
	if len(orig) == 0 {
		return nil
	}
	out := make([]map[string]any, len(orig))
	for i, m := range orig {
		cp := map[string]any{}
		maps.Copy(cp, m)
		out[i] = cp
	}
	return out
}

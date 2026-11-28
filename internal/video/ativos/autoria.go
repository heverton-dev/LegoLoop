package ativos

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/video/roteiro"
)

// CompilarCenasAutorais aplica as receitas do AUTORIA.md a cada linha do roteiro:
// textos foco, camadas graficas (doc, stamp, stopwatch, bars, poly), cartoes em perguntas,
// camera ciclada, fx e maximo de 3 cenas seguidas com personagem.
func CompilarCenasAutorais(linhas []roteiro.LinhaRoteiro, cenarios []Cenario, spriteID string, elencoIDs []string) []Cena {
	cenas := make([]Cena, 0, len(linhas))
	seguidasComHeroi := 0
	cams := []string{"in", "panR", "out", "low", "panL"}
	for i, l := range linhas {
		cenA := cenarios[i%len(cenarios)].ID
		cenB := cenarios[(i+1)%len(cenarios)].ID
		cartao := ehPerguntaDeVirada(i, l.Text)
		incluirHeroi := !cartao && seguidasComHeroi < 3
		if incluirHeroi {
			seguidasComHeroi++
		} else {
			seguidasComHeroi = 0
		}
		heroiID := escolherHeroiCena(i, elencoIDs)
		c := montarCenaAutoral(i, l, cenA, cenB, cams[i%len(cams)], spriteID, heroiID, incluirHeroi, cartao)
		cenas = append(cenas, c)
	}
	return cenas
}

func ehPerguntaDeVirada(idx int, texto string) bool {
	t := strings.TrimSpace(texto)
	return idx > 0 && strings.HasSuffix(t, "?") && len([]rune(t)) <= 78
}

func escolherHeroiCena(idx int, ids []string) string {
	if len(ids) == 0 {
		return ""
	}
	return ids[idx%len(ids)]
}

func montarCenaAutoral(idx int, l roteiro.LinhaRoteiro, cenA, cenB, cam, spriteID, heroiID string, comHeroi, cartao bool) Cena {
	txt1, cor1 := extrairTextoFoco(l.Text, l.Natureza)
	ladoHeroi := 580.0
	ladoTexto := 1280.0
	if idx%2 == 1 {
		ladoHeroi, ladoTexto = 1320.0, 640.0
	}
	if !comHeroi {
		ladoTexto = 960.0
	}
	return Cena{
		IndiceLinha:     idx,
		Texto:           l.Text,
		CenarioID:       cenA,
		CenarioID2:      cenB,
		MetaforaVisual:  derivarMetaforaVisual(l),
		Cartao:          cartao,
		TextoTela:       txt1,
		CorTexto:        cor1,
		TextoTelaX:      ladoTexto,
		TextoTelaY:      210,
		Camera:          cam,
		FX:              []string{"dust", "bokeh"},
		Camadas:         camadasSemanticasDaLinha(l, ladoTexto),
		ElementosMoveis: elementosDaCenaAutoral(idx, spriteID, heroiID, ladoHeroi, l.Natureza, comHeroi),
	}
}

func extrairTextoFoco(texto, natureza string) (string, string) {
	cor := corPorSentimento(texto, natureza)
	reNum := regexp.MustCompile(`(?i)((?:US\$|R\$)\s?\d[\d.,]*\s?(?:bi|mi|milh[oõ]es|bilh[oõ]es)?|\d[\d.,]*\s?(?:%|por cento|horas|dias|anos|vezes|x))`)
	if m := strings.TrimSpace(reNum.FindString(texto)); m != "" {
		return strings.ToUpper(m), cor
	}
	palavras := palavrasRelevantes(texto)
	if len(palavras) >= 2 {
		return strings.ToUpper(strings.Join(palavras[:min(3, len(palavras))], " ")), cor
	}
	return "", cor
}

func corPorSentimento(texto, natureza string) string {
	tl := strings.ToLower(texto)
	for _, r := range []string{"risco", "falha", "press", "esgotamento", "sobrecarga", "perigo", "crise", "queda", "90 horas", "estresse"} {
		if strings.Contains(tl, r) {
			return "#FF5A5A"
		}
	}
	if natureza == "fato_consolidado" {
		return "#3DDC84"
	}
	return "#FFD24A"
}

func palavrasRelevantes(texto string) []string {
	ignorar := map[string]bool{
		"para": true, "como": true, "mais": true, "pelo": true, "pela": true, "isso": true,
		"essa": true, "esse": true, "quando": true, "entanto": true, "apenas": true, "muito": true,
	}
	var out []string
	for _, p := range strings.Fields(texto) {
		limpa := strings.Trim(p, ".,;:!?\"'()-")
		if len([]rune(limpa)) >= 5 && !ignorar[strings.ToLower(limpa)] {
			out = append(out, limpa)
		}
	}
	return out
}

func camadasSemanticasDaLinha(l roteiro.LinhaRoteiro, x float64) []map[string]any {
	var cs []map[string]any
	tl := strings.ToLower(l.Text)
	if strings.ContainsAny(tl, "0123456789") && (strings.Contains(tl, "hora") || strings.Contains(tl, "dia") || strings.Contains(tl, "rápido") || strings.Contains(tl, "tempo")) {
		cs = append(cs, map[string]any{"kind": "stopwatch", "x": 1640.0, "y": 800.0, "s": 0.55, "dur": 3.5, "to": 270, "in": "pop", "start": 0.3, "hud": true})
	}
	if l.Natureza == "fato_consolidado" || l.Natureza == "evidencia_preliminar" {
		hdr := "DADO VERIFICADO"
		if l.Natureza == "evidencia_preliminar" {
			hdr = "ESTUDO PRELIMINAR"
		}
		cs = append(cs, map[string]any{"kind": "doc", "x": x, "y": 560.0, "w": 500.0, "header": hdr, "rows": []map[string]any{{"text": "Fonte auditada", "hl": true}, {"gray": true, "w": 0.7}}, "in": "pop", "start": 0.35, "hud": true})
	}
	if l.Natureza == "alegacao_terceiro" {
		cs = append(cs, map[string]any{"kind": "stamp", "text": "RELATO / ALEGAÇÃO", "x": x, "y": 340.0, "size": 58.0, "color": "#FFD24A", "start": 0.8, "in": "stamp", "rot": -6.0, "hud": true})
	}
	if len(cs) == 0 {
		cs = append(cs, map[string]any{"kind": "bars", "x": x, "y": 880.0, "n": 18, "w": 520.0, "h": 180.0, "appear": 1.1, "hud": true})
	}
	return cs
}

func elementosDaCenaAutoral(idx int, spriteID, heroiID string, x float64, natureza string, comHeroi bool) []ElementoMovel {
	pose, humor := poseHumorPorNatureza(idx, natureza)
	ent := "left"
	if x > 960 {
		ent = "right"
	}
	els := []ElementoMovel{
		{ID: fmt.Sprintf("l%d_objeto", idx), Tipo: "objeto", Acao: "destaque_foco"},
		{ID: fmt.Sprintf("l%d_camera", idx), Tipo: "camera", Acao: "travelling_suave"},
		{ID: fmt.Sprintf("l%d_efeito", idx), Tipo: "efeito", Acao: "particulas_profundidade"},
	}
	if comHeroi {
		p := ElementoMovel{
			ID: fmt.Sprintf("l%d_personagem", idx), Tipo: "personagem", SpriteID: spriteID,
			Acao: "gesto_enfatiza", PersonagemID: heroiID, Pose: pose, Humor: humor, Fala: true,
			X: x, Y: 1000, Escala: 1.05, Entrada: ent,
		}
		return append([]ElementoMovel{p}, els...)
	}
	els = append(els, ElementoMovel{ID: fmt.Sprintf("l%d_grafico", idx), Tipo: "objeto", Acao: "grafico_tela"})
	return els
}

func poseHumorPorNatureza(idx int, natureza string) (string, string) {
	poses := []string{"present", "point", "gesture", "think"}
	switch natureza {
	case "alegacao_terceiro":
		return "gesture", "worry"
	case "evidencia_preliminar":
		return "think", "neutral"
	case "fato_consolidado":
		return "point", "neutral"
	default:
		return poses[idx%len(poses)], "neutral"
	}
}

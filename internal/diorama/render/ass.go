package render

import (
	"fmt"
	"math"
	"strings"
)

const cabecalhoASS = `[Script Info]
ScriptType: v4.00+
PlayResX: 1920
PlayResY: 1080
WrapStyle: 2
ScaledBorderAndShadow: yes

[V4+ Styles]
Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding
Style: Karaoke,Arial,64,&H0000E6FF,&H00FFFFFF,&H00000000,&H80000000,-1,0,0,0,100,100,0,0,1,4,1,2,96,96,140,1

[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
`

// limparASS troca os caracteres de controle do ASS; sem variavel global.
func limparASS(texto string) string {
	return strings.NewReplacer("{", "(", "}", ")", "\\", "").Replace(texto)
}

func centesimos(s float64) int { return int(math.Round(math.Max(s, 0) * 100)) }

func tempoASS(s float64) string {
	cs := centesimos(s)
	return fmt.Sprintf("%d:%02d:%02d.%02d", cs/360000, (cs/6000)%60, (cs/100)%60, cs%100)
}

// textoKaraoke marca cada palavra com {\k<centesimos>}. A pausa entre palavras vira {\k} vazio,
// e a quebra de linha segue a mesma regra do SRT (42 colunas, duas linhas equilibradas).
func textoKaraoke(ps []PalavraTempo) string {
	linhas := quebrarLinhas(textosDasPalavras(ps), LarguraSRT)
	var sb strings.Builder
	i, fimAnterior := 0, centesimos(ps[0].InicioS)
	for l, linha := range linhas {
		n := len(strings.Fields(linha))
		for j := 0; j < n; j, i = j+1, i+1 {
			p := ps[i]
			ini, fim := centesimos(p.InicioS), centesimos(p.FimS)
			if ini > fimAnterior {
				fmt.Fprintf(&sb, "{\\k%d}", ini-fimAnterior)
			}
			fmt.Fprintf(&sb, "{\\k%d}%s", fim-ini, limparASS(p.Texto))
			fimAnterior = fim
			if j < n-1 {
				sb.WriteString(" ")
			}
		}
		if l < len(linhas)-1 {
			sb.WriteString("\\N")
		}
	}
	return sb.String()
}

// GerarASS cria a legenda karaoke (uma fala por batida com palavras) no quadro de 1920x1080.
func GerarASS(tl Timeline) string {
	var sb strings.Builder
	sb.WriteString(cabecalhoASS)
	for _, b := range tl.Batidas {
		ps := camadaLegenda(b).Palavras
		if len(ps) == 0 {
			continue
		}
		fmt.Fprintf(&sb, "Dialogue: 0,%s,%s,Karaoke,,0,0,0,,%s\n", tempoASS(ps[0].InicioS), tempoASS(ps[len(ps)-1].FimS), textoKaraoke(ps))
	}
	return sb.String()
}

package render

import (
	"strings"
	"unicode"
)

// Visemas simples para sincronia labial dos sprites. Os nomes sao o contrato com a timeline.
const (
	VisemaRepouso = "X"
	VisemaA       = "A"
	VisemaE       = "E"
	VisemaI       = "I"
	VisemaO       = "O"
	VisemaU       = "U"
	VisemaMBP     = "MBP"
	VisemaFV      = "FV"
	VisemaCons    = "CONS"
)

// LacunaRepousoS e a pausa minima entre palavras que fecha a boca do sprite.
const LacunaRepousoS = 0.08

type visemaPeso struct {
	nome string
	peso int
}

// visemaDeLetra devolve o visema da letra; letra sem visema devolve ok=false.
func visemaDeLetra(r rune) (string, bool) {
	switch r {
	case 'a', 'á', 'à', 'â', 'ã':
		return VisemaA, true
	case 'e', 'é', 'ê':
		return VisemaE, true
	case 'i', 'í', 'y':
		return VisemaI, true
	case 'o', 'ó', 'ô', 'õ':
		return VisemaO, true
	case 'u', 'ú', 'ü':
		return VisemaU, true
	case 'm', 'b', 'p':
		return VisemaMBP, true
	case 'f', 'v':
		return VisemaFV, true
	}
	return "", false
}

func ehVogalViseme(nome string) bool {
	switch nome {
	case VisemaA, VisemaE, VisemaI, VisemaO, VisemaU:
		return true
	}
	return false
}

// sequenciaVisemas converte a grafia PT-BR em visemas com peso (vogal 2, consoante 1).
// Digrafos lh, nh e ch e os grupos qu/gu antes de e/i viram uma consoante; h isolado e mudo.
func sequenciaVisemas(palavra string) []visemaPeso {
	runas := []rune(strings.ToLower(palavra))
	var seq []visemaPeso
	for i := 0; i < len(runas); i++ {
		r := runas[i]
		if !unicode.IsLetter(r) {
			continue
		}
		nome, salto := classificar(runas, i)
		i += salto
		if nome != "" {
			seq = anexarViseme(seq, nome)
		}
	}
	if len(seq) == 0 {
		return []visemaPeso{{nome: VisemaCons, peso: 1}}
	}
	return seq
}

// classificar devolve o viseme da letra na posicao i e quantas letras extras ela consome.
func classificar(runas []rune, i int) (string, int) {
	r := runas[i]
	var prox rune
	if i+1 < len(runas) {
		prox = runas[i+1]
	}
	switch {
	case (r == 'l' || r == 'n' || r == 'c') && prox == 'h':
		return VisemaCons, 1
	case (r == 'q' || r == 'g') && prox == 'u' && i+2 < len(runas) && strings.ContainsRune("eéêií", runas[i+2]):
		return VisemaCons, 1
	case r == 'h':
		return "", 0
	}
	if nome, ok := visemaDeLetra(r); ok {
		return nome, 0
	}
	return VisemaCons, 0
}

func anexarViseme(seq []visemaPeso, nome string) []visemaPeso {
	peso := 1
	if ehVogalViseme(nome) {
		peso = 2
	}
	if n := len(seq); n > 0 && seq[n-1].nome == nome {
		seq[n-1].peso += peso
		return seq
	}
	return append(seq, visemaPeso{nome: nome, peso: peso})
}

// GerarVisemas devolve uma copia dos timings com visemas por linha; a entrada nao muda.
func GerarVisemas(t Timings) Timings {
	linhas := make([]LinhaTempo, len(t.Linhas))
	for i, l := range t.Linhas {
		l.Visemas = visemasDaLinha(l)
		linhas[i] = l
	}
	t.Linhas = linhas
	return t
}

func visemasDaLinha(l LinhaTempo) []VisemaTempo {
	var vs []VisemaTempo
	cursor := l.InicioS
	for _, p := range l.Palavras {
		vs = anexarRepouso(vs, cursor, p.InicioS)
		vs = append(vs, visemasDaPalavra(p)...)
		cursor = p.FimS
	}
	return anexarRepouso(vs, cursor, l.FimS)
}

func anexarRepouso(vs []VisemaTempo, de, ate float64) []VisemaTempo {
	if ate-de >= LacunaRepousoS {
		vs = append(vs, VisemaTempo{Viseme: VisemaRepouso, InicioS: de, FimS: ate})
	}
	return vs
}

// visemasDaPalavra divide a janela da palavra por peso; os extremos coincidem exatamente com ela.
func visemasDaPalavra(p PalavraTempo) []VisemaTempo {
	seq := sequenciaVisemas(p.Texto)
	total := 0
	for _, v := range seq {
		total += v.peso
	}
	dur := p.FimS - p.InicioS
	saida := make([]VisemaTempo, len(seq))
	acumulado, inicio := 0, p.InicioS
	for i, v := range seq {
		acumulado += v.peso
		fim := p.InicioS + dur*float64(acumulado)/float64(total)
		if i == len(seq)-1 {
			fim = p.FimS
		}
		saida[i] = VisemaTempo{Viseme: v.nome, InicioS: inicio, FimS: fim}
		inicio = fim
	}
	return saida
}

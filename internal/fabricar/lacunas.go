package fabricar

import (
	"slices"
	"strings"

	"legoloop/internal/matriz"
	"legoloop/internal/ste"
)

// maxDesconhecidasCitadas é quantas peças sem perfil os avisos citam pelo nome; o resto vira um resumo.
const maxDesconhecidasCitadas = 5

func init() {
	ste.Default.Register("fabricar.lacuna_etapa", "Lacuna: a etapa %s não tem peça apta.")
	ste.Default.Register("fabricar.caso_sem_etapa", "O caso de uso %s não contou: ele pertence a %s, e nenhuma dessas etapas foi pedida.")
	ste.Default.Register("fabricar.aptidao_desconhecida", "Aptidão desconhecida: %s. A peça não foi escolhida. Rode legoloop matriz gerar.")
	ste.Default.Register("fabricar.aptidao_desconhecida_mais", "Aptidão desconhecida em mais %d peças. Veja todas com legoloop matriz ver.")
}

// Decisao é o que a Fusão decidiu sobre uma peça ou uma etapa. Peca vazia com Etapa preenchida é uma lacuna:
// a etapa foi pedida e nenhuma peça é apta a ela.
type Decisao struct {
	Etapa      string  `json:"etapa,omitempty"`
	Peca       string  `json:"peca"`
	Entrou     bool    `json:"entrou"`
	Aptidao    float64 `json:"aptidao"`
	Score      float64 `json:"score"`
	Motivo     string  `json:"motivo"`
	Origem     string  `json:"origem,omitempty"`     // origem do perfil da peça escolhida: declarado, confirmado ou inferido
	Explicacao string  `json:"explicacao,omitempty"` // escrita pelo cérebro; nunca decide
}

const (
	motivoDesconhecida = "aptidão desconhecida: a peça não tem perfil"
	motivoLacuna       = "lacuna: nenhuma peça apta a esta etapa"
	motivoSemCaso      = "lacuna: nenhum caso de uso pedido pertence a esta etapa, então não há como medir a aptidão"
)

// addLacuna registra uma etapa pedida sem nenhuma peça apta, com o motivo, e a destaca nos avisos.
func (m *Merged) addLacuna(etapa, motivo string) {
	m.Lacunas = append(m.Lacunas, etapa)
	m.Decisoes = append(m.Decisoes, Decisao{Etapa: etapa, Motivo: motivo})
	m.Warnings = append(m.Warnings, ste.Msg("fabricar.lacuna_etapa", etapa))
}

// noteUnknown registra as skills sem perfil: nunca são escolhidas, e a aptidão delas é desconhecida.
// As de maior Score vêm primeiro, porque são as que mais valem catalogar.
func (m *Merged) noteUnknown(skills []matriz.Entrada) {
	var sem []matriz.Entrada
	for _, e := range skills {
		if e.Perfil == nil {
			sem = append(sem, e)
		}
	}
	slices.SortFunc(sem, func(a, b matriz.Entrada) int {
		if a.Rating.TotalScore != b.Rating.TotalScore {
			return cmpScoreDesc(a.Rating.TotalScore, b.Rating.TotalScore)
		}
		return strings.Compare(a.Ref, b.Ref)
	})
	for i, e := range sem {
		m.Unknown = append(m.Unknown, pieceOf(e))
		m.Decisoes = append(m.Decisoes, Decisao{Peca: e.Nome, Score: e.Rating.TotalScore, Motivo: motivoDesconhecida})
		if i < maxDesconhecidasCitadas {
			m.Warnings = append(m.Warnings, ste.Msg("fabricar.aptidao_desconhecida", e.Nome))
		}
	}
	if extra := len(sem) - maxDesconhecidasCitadas; extra > 0 {
		m.Warnings = append(m.Warnings, ste.Msg("fabricar.aptidao_desconhecida_mais", extra))
	}
}

func cmpScoreDesc(a, b float64) int {
	if a > b {
		return -1
	}
	return 1
}

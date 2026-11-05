package fabricar

import (
	"fmt"
	"slices"
	"strings"

	"legoloop/internal/matriz"
	"legoloop/internal/ste"
)

// MergeByFit é a Camada 3 pela aptidão ao caso de uso: para cada etapa pedida escolhe a peça mais apta da matriz
// (matriz.Rank), e o Score só desempata. Só skills entram (MCPs, hooks, scripts e loops nunca são herdados);
// cópias idênticas (mesmo hash) valem uma só. A mesma peça em duas etapas vira um passo só. Sem nenhuma peça
// apta valem os passos padrão. Peça excluída por uso (quando_nao_usar ou natureza) é avisada com o motivo.
func MergeByFit(n matriz.Needs, es []matriz.Entrada, opts MergeOptions) Merged {
	max := opts.MaxSteps
	if max <= 0 {
		max = defaultMaxSteps
	}
	skills := dedupeCopies(onlySkills(es))
	grupos := [][]string{nil} // sem etapas pedidas: um ranking geral
	if len(n.Etapas) > 0 {
		grupos = nil
		for _, e := range n.Etapas {
			grupos = append(grupos, []string{e})
		}
	}
	var m Merged
	avisadas := map[string]bool{}
	for _, etapas := range grupos {
		casos := n.Casos
		if len(etapas) == 1 && len(n.Casos) > 0 { // o caso conta só na etapa a que pertence
			if casos = opts.Vocab.CasosDaEtapa(n.Casos, etapas[0]); len(casos) == 0 {
				m.addLacuna(etapas[0], motivoSemCaso)
				continue
			}
		}
		ranked := matriz.Rank(matriz.Needs{Natureza: n.Natureza, Etapas: etapas, Casos: casos}, skills, 0)
		m.noteExclusions(ranked, avisadas)
		if !m.addBest(ranked, etapas, max) {
			break
		}
	}
	m.noteCasosSemEtapa(n, opts.Vocab)
	m.noteUnknown(skills)
	if len(m.Steps) == 0 {
		m.Steps = DefaultSteps()
		m.Warnings = append(m.Warnings, "nenhuma peça apta na matriz; usando os passos padrão")
	}
	return m
}

// noteExclusions registra, uma vez por peça, as exclusões por uso.
func (m *Merged) noteExclusions(ranked []matriz.Ranked, avisadas map[string]bool) {
	for _, r := range ranked {
		if !r.ExcluidaPorUso() || avisadas[r.Entrada.Ref] {
			continue
		}
		avisadas[r.Entrada.Ref] = true
		m.Excluded = append(m.Excluded, pieceOf(r.Entrada))
		m.Decisoes = append(m.Decisoes, Decisao{Peca: r.Entrada.Nome, Score: r.Entrada.Rating.TotalScore, Motivo: r.Motivo})
		m.Warnings = append(m.Warnings, fmt.Sprintf("peça %q excluída: %s", r.Entrada.Nome, r.Motivo))
	}
}

// addBest acrescenta a melhor peça apta à etapa (todas, no ranking geral, que não tem etapa) respeitando o
// limite de passos. Etapa sem peça apta vira lacuna. Devolve false quando o limite foi atingido.
func (m *Merged) addBest(ranked []matriz.Ranked, etapas []string, max int) bool {
	geral, etapa := len(etapas) == 0, ""
	if !geral {
		etapa = etapas[0]
	}
	achou := false
	for _, r := range ranked {
		if r.Excluida {
			break
		}
		achou = true
		if slices.ContainsFunc(m.Steps, func(s Step) bool { return s.Piece == r.Entrada.Nome }) {
			if !geral {
				m.decide(etapa, r)
				return true
			}
			continue
		}
		if len(m.Steps) >= max {
			m.Warnings = append(m.Warnings, fmt.Sprintf("limite de %d passos atingido; as demais peças foram ignoradas", max))
			return false
		}
		m.Steps = append(m.Steps, Step{Piece: r.Entrada.Nome})
		m.Used = append(m.Used, pieceOf(r.Entrada))
		m.decide(etapa, r)
		if !geral {
			return true
		}
	}
	if !achou && !geral {
		m.addLacuna(etapa, motivoSemPeca(ranked))
	}
	return true
}

// motivoSemPeca explica a lacuna: se alguma peça fazia a etapa mas não cobria o caso pedido, o motivo diz qual
// caso ninguém cobre; senão, nenhuma peça é apta à etapa.
func motivoSemPeca(ranked []matriz.Ranked) string {
	for _, r := range ranked {
		if strings.HasPrefix(r.Motivo, "caso de uso:") {
			return "lacuna: nenhuma peça desta etapa cobre o caso pedido (" + strings.TrimPrefix(r.Motivo, "caso de uso: a peça não cobre ") + ")"
		}
	}
	return motivoLacuna
}

// noteCasosSemEtapa avisa o caso pedido que não pertence a nenhuma etapa pedida: ele não contou em nenhuma.
func (m *Merged) noteCasosSemEtapa(n matriz.Needs, v matriz.Vocab) {
	if len(n.Etapas) == 0 {
		return
	}
	for _, c := range n.Casos {
		dele := v.EtapasDoCaso(c)
		if slices.ContainsFunc(n.Etapas, func(e string) bool { return slices.Contains(dele, e) }) {
			continue
		}
		onde := strings.Join(dele, ", ")
		if onde == "" {
			onde = "nenhuma etapa do vocabulário"
		}
		m.Warnings = append(m.Warnings, ste.Msg("fabricar.caso_sem_etapa", c, onde))
	}
}

// decide registra que a peça entrou para a etapa.
func (m *Merged) decide(etapa string, r matriz.Ranked) {
	m.Decisoes = append(m.Decisoes, Decisao{Etapa: etapa, Peca: r.Entrada.Nome, Entrou: true, Aptidao: r.Aptidao,
		Score: r.Entrada.Rating.TotalScore, Motivo: r.Motivo})
}

func pieceOf(e matriz.Entrada) DiscoveredPiece {
	return DiscoveredPiece{Name: e.Nome, Kind: KindSkill, Path: e.Ref, SHA256: e.Hash, Score: e.Rating.TotalScore, Tier: e.Rating.Tier,
		Untrusted: naoConfiavel(e.Ref)}
}

func onlySkills(es []matriz.Entrada) []matriz.Entrada {
	var out []matriz.Entrada
	for _, e := range es {
		if e.Tipo == string(KindSkill) {
			out = append(out, e)
		}
	}
	return out
}

// dedupeCopies junta as cópias idênticas (mesmo hash): fica a de menor Ref. Entrada sem hash nunca é cópia.
func dedupeCopies(es []matriz.Entrada) []matriz.Entrada {
	best := map[string]string{}
	for _, e := range es {
		if cur, ok := best[e.Hash]; e.Hash != "" && (!ok || e.Ref < cur) {
			best[e.Hash] = e.Ref
		}
	}
	var out []matriz.Entrada
	for _, e := range es {
		if e.Hash == "" || best[e.Hash] == e.Ref {
			out = append(out, e)
		}
	}
	return out
}

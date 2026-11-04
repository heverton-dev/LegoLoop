package fabricar

import (
	"fmt"
	"slices"

	"legoloop/internal/matriz"
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
		ranked := matriz.Rank(matriz.Needs{Natureza: n.Natureza, Etapas: etapas, Casos: n.Casos}, skills, 0)
		m.noteExclusions(ranked, avisadas)
		if !m.addBest(ranked, etapas == nil, max) {
			break
		}
	}
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
		m.Warnings = append(m.Warnings, fmt.Sprintf("peça %q excluída: %s", r.Entrada.Nome, r.Motivo))
	}
}

// addBest acrescenta a melhor peça ainda não usada (todas, se for o ranking geral) respeitando o limite.
// Devolve false quando o limite de passos foi atingido.
func (m *Merged) addBest(ranked []matriz.Ranked, geral bool, max int) bool {
	for _, r := range ranked {
		if r.Excluida {
			break
		}
		if slices.ContainsFunc(m.Steps, func(s Step) bool { return s.Piece == r.Entrada.Nome }) {
			if geral {
				continue
			}
			return true
		}
		if len(m.Steps) >= max {
			m.Warnings = append(m.Warnings, fmt.Sprintf("limite de %d passos atingido; as demais peças foram ignoradas", max))
			return false
		}
		m.Steps = append(m.Steps, Step{Piece: r.Entrada.Nome})
		m.Used = append(m.Used, pieceOf(r.Entrada))
		if !geral {
			return true
		}
	}
	return true
}

func pieceOf(e matriz.Entrada) DiscoveredPiece {
	return DiscoveredPiece{Name: e.Nome, Kind: KindSkill, Path: e.Ref, SHA256: e.Hash, Score: e.Rating.TotalScore, Tier: e.Rating.Tier}
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

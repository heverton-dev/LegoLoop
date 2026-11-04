package fabricar

import (
	"fmt"
	"slices"

	"legoloop/internal/reputation"
)

// MergeOptions controla a Camada 3.
type MergeOptions struct {
	TrustRemovable bool // true = aceita peças Untrusted
	MaxSteps       int  // padrão 5
}

const defaultMaxSteps = 5

// Merged é o resultado da fusão: os passos do processo e o que ficou de fora.
type Merged struct {
	Steps    []Step
	Used     []DiscoveredPiece
	Excluded []DiscoveredPiece
	Warnings []string
	// Só preenchidos por MergeByFit:
	Decisoes []Decisao         // o que a Fusão decidiu sobre cada etapa e peça
	Lacunas  []string          // etapas pedidas sem nenhuma peça apta
	Unknown  []DiscoveredPiece // skills sem perfil: nunca escolhidas, aptidão desconhecida
	// Preenchido por Explain: ordens que o cérebro propôs e que o ranking calculado não seguiu.
	Divergencias []string
}

// DefaultSteps são os passos usados quando nenhuma peça curada é utilizável: especificar e planejar.
func DefaultSteps() []Step { return []Step{{Piece: "aidd-spec"}, {Piece: "aidd-tickets"}} }

// ArchitectMerge é a Camada 3. Seleciona por Score decrescente e Path crescente e só aproveita skills:
// MCPs, hooks, scripts e loops varridos nunca são herdados (default-deny). Fica de fora o que é de origem
// não confiável (salvo TrustRemovable) ou de tier tóxico; de nomes repetidos fica só a melhor.
func ArchitectMerge(pieces []DiscoveredPiece, opts MergeOptions) Merged {
	max := opts.MaxSteps
	if max <= 0 {
		max = defaultMaxSteps
	}
	sorted := slices.Clone(pieces)
	sortPieces(sorted)
	var m Merged
	seen := map[string]bool{}
	for _, p := range sorted {
		if p.Kind != KindSkill || p.Name == "" || seen[p.Name] {
			continue
		}
		if reason := excludeReason(p, opts); reason != "" {
			m.Excluded = append(m.Excluded, p)
			m.Warnings = append(m.Warnings, fmt.Sprintf("peça %q excluída: %s (%s)", p.Name, reason, p.Drive))
			continue
		}
		if len(m.Steps) >= max {
			m.Warnings = append(m.Warnings, fmt.Sprintf("limite de %d passos atingido; as demais peças foram ignoradas", max))
			break
		}
		seen[p.Name] = true
		m.Steps = append(m.Steps, Step{Piece: p.Name})
		m.Used = append(m.Used, p)
	}
	if len(m.Steps) == 0 {
		m.Steps = DefaultSteps()
		m.Warnings = append(m.Warnings, "nenhuma peça curada utilizável; usando os passos padrão")
	}
	return m
}

func excludeReason(p DiscoveredPiece, opts MergeOptions) string {
	switch {
	case p.Untrusted && !opts.TrustRemovable:
		return "origem não confiável (mídia removível ou rede)"
	case p.Tier == reputation.TierToxica:
		return "tier tóxico"
	}
	return ""
}

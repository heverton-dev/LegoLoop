package fabricar

import "legoloop/internal/matriz"

// Step é um passo escolhido pela Fusão: a skill que executa a etapa.
type Step struct {
	Piece string // nome da skill
}

// MergeOptions controla a Camada 3.
type MergeOptions struct {
	MaxSteps int          // padrão 5
	Vocab    matriz.Vocab // diz a etapa de cada caso de uso; vazio = o vocabulário padrão
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

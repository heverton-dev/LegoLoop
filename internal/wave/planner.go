package wave

import (
	"legoloop/internal/board"
)

// PlanWaves organiza uma lista de cartões em lotes ordenados de execução concorrente segura (waves).
// Cartões que tocam nos mesmos arquivos nunca são colocados na mesma onda.
func PlanWaves(cards []board.Card) []board.Wave {
	if len(cards) == 0 {
		return nil
	}

	type waveSlot struct {
		cards        []board.Card
		occupiedSets [][]string
	}

	var slots []*waveSlot

	for _, card := range cards {
		cardFiles := ExtractAffectedFiles(card)
		placed := false

		// Tenta alocar no primeiro slot compatível (sem overlap com nenhum card já no slot)
		for _, slot := range slots {
			conflict := false
			for _, existingFiles := range slot.occupiedSets {
				if HasOverlap(cardFiles, existingFiles) {
					conflict = true
					break
				}
			}

			if !conflict {
				slot.cards = append(slot.cards, card)
				slot.occupiedSets = append(slot.occupiedSets, cardFiles)
				placed = true
				break
			}
		}

		// Se conflitar com todos os slots anteriores, cria um novo slot (nova onda)
		if !placed {
			slots = append(slots, &waveSlot{
				cards:        []board.Card{card},
				occupiedSets: [][]string{cardFiles},
			})
		}
	}

	result := make([]board.Wave, len(slots))
	for i, slot := range slots {
		result[i] = board.Wave{
			Index: i,
			Cards: slot.cards,
		}
	}

	return result
}

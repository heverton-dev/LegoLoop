package wave

import (
	"legoloop/internal/board"
)

type waveSlot struct {
	cards        []board.Card
	occupiedSets [][]string
}

// accepts informa se o cartão não tem overlap com nenhum cartão já no slot.
func (w *waveSlot) accepts(files []string) bool {
	for _, existing := range w.occupiedSets {
		if HasOverlap(files, existing) {
			return false
		}
	}
	return true
}

// PlanWaves organiza uma lista de cartões em lotes ordenados de execução concorrente segura (waves).
// Cartões que tocam nos mesmos arquivos nunca são colocados na mesma onda.
func PlanWaves(cards []board.Card) []board.Wave {
	if len(cards) == 0 {
		return nil
	}
	var slots []*waveSlot
	for _, card := range cards {
		cardFiles := ExtractAffectedFiles(card)
		// Primeiro slot compatível; se conflitar com todos, abre uma nova onda
		placed := false
		for _, slot := range slots {
			if slot.accepts(cardFiles) {
				slot.cards = append(slot.cards, card)
				slot.occupiedSets = append(slot.occupiedSets, cardFiles)
				placed = true
				break
			}
		}
		if !placed {
			slots = append(slots, &waveSlot{cards: []board.Card{card}, occupiedSets: [][]string{cardFiles}})
		}
	}
	result := make([]board.Wave, len(slots))
	for i, slot := range slots {
		result[i] = board.Wave{Index: i, Cards: slot.cards}
	}
	return result
}

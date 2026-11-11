package wave

import (
	"legoloop/internal/board"
)

type waveSlot struct {
	tickets      []board.Ticket
	occupiedSets [][]string
}

// accepts informa se o ticket não tem overlap com nenhum ticket já no slot.
func (w *waveSlot) accepts(files []string) bool {
	for _, existing := range w.occupiedSets {
		if HasOverlap(files, existing) {
			return false
		}
	}
	return true
}

// PlanWaves organiza uma lista de tickets em lotes ordenados de execução concorrente segura (waves).
// Tickets que tocam nos mesmos arquivos nunca são colocados na mesma onda.
func PlanWaves(tickets []board.Ticket) []board.Wave {
	if len(tickets) == 0 {
		return nil
	}
	var slots []*waveSlot
	for _, ticket := range tickets {
		ticketFiles := ExtractAffectedFiles(ticket)
		placed := false
		for _, slot := range slots {
			if slot.accepts(ticketFiles) {
				slot.tickets = append(slot.tickets, ticket)
				slot.occupiedSets = append(slot.occupiedSets, ticketFiles)
				placed = true
				break
			}
		}
		if !placed {
			slots = append(slots, &waveSlot{tickets: []board.Ticket{ticket}, occupiedSets: [][]string{ticketFiles}})
		}
	}
	result := make([]board.Wave, len(slots))
	for i, slot := range slots {
		result[i] = board.Wave{Index: i, Cards: slot.tickets}
	}
	return result
}

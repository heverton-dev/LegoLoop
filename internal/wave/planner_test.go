package wave_test

import (
	"testing"

	"legoloop/internal/board"
	"legoloop/internal/wave"
)

func TestPlanWaves_DisjointTicketsRunInSameWave(t *testing.T) {
	tickets := []board.Ticket{
		{
			Number: 1,
			Title:  "Feature A",
			Body:   "Touches internal/pkg_a/a.go",
		},
		{
			Number: 2,
			Title:  "Feature B",
			Body:   "Touches internal/pkg_b/b.go",
		},
		{
			Number: 3,
			Title:  "Feature C",
			Body:   "Touches internal/pkg_c/c.go",
		},
	}

	waves := wave.PlanWaves(tickets)

	if len(waves) != 1 {
		t.Fatalf("esperava 1 onda para tickets totalmente disjuntos, obteve %d", len(waves))
	}
	if len(waves[0].Cards) != 3 {
		t.Fatalf("esperava 3 tickets na primeira onda, obteve %d", len(waves[0].Cards))
	}
}

func TestPlanWaves_ConflictingTicketsPartitionedAcrossWaves(t *testing.T) {
	tickets := []board.Ticket{
		{
			Number: 1,
			Title:  "Ticket 1 modifies main",
			Body:   "Touches cmd/legoloop/main.go",
		},
		{
			Number: 2,
			Title:  "Ticket 2 also modifies main",
			Body:   "Touches cmd/legoloop/main.go and internal/cli/cli.go",
		},
		{
			Number: 3,
			Title:  "Ticket 3 touches cli",
			Body:   "Touches internal/cli/cli.go",
		},
		{
			Number: 4,
			Title:  "Ticket 4 touches unrelated docs",
			Body:   "Touches docs/README.md",
		},
	}

	waves := wave.PlanWaves(tickets)

	for wIdx, w := range waves {
		for i := 0; i < len(w.Cards); i++ {
			filesI := wave.ExtractAffectedFiles(w.Cards[i])
			for j := i + 1; j < len(w.Cards); j++ {
				filesJ := wave.ExtractAffectedFiles(w.Cards[j])
				if wave.HasOverlap(filesI, filesJ) {
					t.Fatalf("conflito detectado na onda %d entre Ticket #%d e Ticket #%d:\nArquivos I: %v\nArquivos J: %v",
						wIdx, w.Cards[i].Number, w.Cards[j].Number, filesI, filesJ)
				}
			}
		}
	}

	if len(waves) < 2 {
		t.Fatalf("esperava pelo menos 2 ondas para resolver conflito em main.go, obteve %d", len(waves))
	}
}

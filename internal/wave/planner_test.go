package wave_test

import (
	"testing"

	"legoloop/internal/board"
	"legoloop/internal/wave"
)

func TestPlanWaves_DisjointCardsRunInSameWave(t *testing.T) {
	cards := []board.Card{
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

	waves := wave.PlanWaves(cards)

	if len(waves) != 1 {
		t.Fatalf("esperava 1 onda para cartões totalmente disjuntos, obteve %d", len(waves))
	}
	if len(waves[0].Cards) != 3 {
		t.Fatalf("esperava 3 cartões na primeira onda, obteve %d", len(waves[0].Cards))
	}
}

func TestPlanWaves_ConflictingCardsPartitionedAcrossWaves(t *testing.T) {
	cards := []board.Card{
		{
			Number: 1,
			Title:  "Card 1 modifies main",
			Body:   "Touches cmd/legoloop/main.go",
		},
		{
			Number: 2,
			Title:  "Card 2 also modifies main",
			Body:   "Touches cmd/legoloop/main.go and internal/cli/cli.go",
		},
		{
			Number: 3,
			Title:  "Card 3 touches cli",
			Body:   "Touches internal/cli/cli.go",
		},
		{
			Number: 4,
			Title:  "Card 4 touches unrelated docs",
			Body:   "Touches docs/README.md",
		},
	}

	waves := wave.PlanWaves(cards)

	// Nenhum card na mesma onda pode ter overlap
	for wIdx, w := range waves {
		for i := 0; i < len(w.Cards); i++ {
			filesI := wave.ExtractAffectedFiles(w.Cards[i])
			for j := i + 1; j < len(w.Cards); j++ {
				filesJ := wave.ExtractAffectedFiles(w.Cards[j])
				if wave.HasOverlap(filesI, filesJ) {
					t.Fatalf("conflito detectado na onda %d entre Card #%d e Card #%d:\nArquivos I: %v\nArquivos J: %v",
						wIdx, w.Cards[i].Number, w.Cards[j].Number, filesI, filesJ)
				}
			}
		}
	}

	// Card 1 e Card 2 conflitam em main.go, logo deve haver pelo menos 2 ondas
	if len(waves) < 2 {
		t.Fatalf("esperava pelo menos 2 ondas para resolver conflito em main.go, obteve %d", len(waves))
	}
}

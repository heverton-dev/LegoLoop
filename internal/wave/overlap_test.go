package wave_test

import (
	"reflect"
	"sort"
	"testing"

	"legoloop/internal/board"
	"legoloop/internal/wave"
)

func TestExtractAffectedFiles(t *testing.T) {
	card := board.Card{
		Number: 42,
		Title:  "Refactor scanner and update dag types",
		Body: `This ticket addresses performance issues in the scanner.
Files affected:
- internal/assets/scanner.go
- ` + "`internal/dag/types.go`" + `
- cmd/legoloop/main.go
Also touches docs/specs/SPEC.md.`,
	}

	files := wave.ExtractAffectedFiles(card)
	sort.Strings(files)

	expected := []string{
		"cmd/legoloop/main.go",
		"docs/specs/SPEC.md",
		"internal/assets/scanner.go",
		"internal/dag/types.go",
	}

	if !reflect.DeepEqual(files, expected) {
		t.Fatalf("arquivos extraídos incorretos.\nEsperado: %v\nObtido:   %v", expected, files)
	}
}

func TestExtractAffectedFiles_FallbackToTitle(t *testing.T) {
	card := board.Card{
		Number: 99,
		Title:  "Fix bug in internal/board/guard.go",
		Body:   "Simple fix with no explicit file list.",
	}

	files := wave.ExtractAffectedFiles(card)
	expected := []string{"internal/board/guard.go"}

	if !reflect.DeepEqual(files, expected) {
		t.Fatalf("arquivos extraídos do título incorretos.\nEsperado: %v\nObtido:   %v", expected, files)
	}
}

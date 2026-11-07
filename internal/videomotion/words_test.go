package videomotion_test

import (
	"strings"
	"testing"

	"legoloop/internal/videomotion"
)

func TestParseAndValidateWords(t *testing.T) {
	validJSON := `[
		{"index": 0, "word": "Bem-vindo", "start_ms": 100, "end_ms": 450, "confidence": 0.99},
		{"index": 1, "word": "ao", "start_ms": 460, "end_ms": 550, "confidence": 0.98},
		{"index": 2, "word": "LegoLoop", "start_ms": 560, "end_ms": 1100, "confidence": 0.95}
	]`

	timeline, err := videomotion.ParseAndValidateWords(strings.NewReader(validJSON))
	if err != nil {
		t.Fatalf("ParseAndValidateWords falhou inesperadamente: %v", err)
	}

	if len(timeline.Words) != 3 {
		t.Fatalf("esperava 3 palavras, obteve %d", len(timeline.Words))
	}

	if timeline.TotalDurationMs != 1100 {
		t.Errorf("duração total esperada 1100ms, obteve %d", timeline.TotalDurationMs)
	}

	if timeline.Words[2].Word != "LegoLoop" {
		t.Errorf("palavra 2 esperada 'LegoLoop', obteve '%s'", timeline.Words[2].Word)
	}
}

func TestParseAndValidateWords_InvertedTimestamps(t *testing.T) {
	invalidJSON := `[
		{"index": 0, "word": "Erro", "start_ms": 500, "end_ms": 200, "confidence": 0.99}
	]`

	_, err := videomotion.ParseAndValidateWords(strings.NewReader(invalidJSON))
	if err == nil {
		t.Fatal("esperava erro para start_ms > end_ms, obteve nil")
	}
}

func TestParseAndValidateWords_EmptyStream(t *testing.T) {
	_, err := videomotion.ParseAndValidateWords(strings.NewReader("[]"))
	if err == nil {
		t.Fatal("esperava erro para lista vazia de palavras, obteve nil")
	}
}

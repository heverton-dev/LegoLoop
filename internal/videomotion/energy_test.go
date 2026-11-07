package videomotion_test

import (
	"testing"

	"legoloop/internal/videomotion"
)

func TestDetectPunchWords(t *testing.T) {
	words := []videomotion.WordItem{
		{Index: 0, Word: "Este", StartMs: 100, EndMs: 300},
		{Index: 1, Word: "é", StartMs: 310, EndMs: 400},
		{Index: 2, Word: "o", StartMs: 410, EndMs: 500},
		{Index: 3, Word: "FUTURO", StartMs: 510, EndMs: 900}, // palavra falada com ênfase alta
		{Index: 4, Word: "do", StartMs: 910, EndMs: 1000},
		{Index: 5, Word: "código", StartMs: 1010, EndMs: 1400},
	}

	frames := make([]videomotion.AudioEnergyFrame, 15)
	for i := range frames {
		tMs := int64(i * 100)
		e := 0.2 // energia média normal
		if tMs >= 500 && tMs <= 900 {
			e = 0.85 // pico evidente de energia vocal na palavra FUTURO
		}
		frames[i] = videomotion.AudioEnergyFrame{TimestampMs: tMs, Energy: e}
	}

	cfg := videomotion.EnergyDetectorConfig{
		MinRelativeBoost: 1.4, // 40% acima da média da fala
	}

	punched, err := videomotion.DetectPunchWords(words, frames, cfg)
	if err != nil {
		t.Fatalf("DetectPunchWords falhou: %v", err)
	}

	if len(punched) != len(words) {
		t.Fatalf("tamanho esperado %d, obteve %d", len(words), len(punched))
	}

	// Palavra 3 ("FUTURO") deve ser punch word
	if !punched[3].IsPunchWord {
		t.Errorf("esperava que a palavra 3 'FUTURO' fosse punch word, mas IsPunchWord é false")
	}

	// Palavra 0 ("Este") deve ser normal
	if punched[0].IsPunchWord {
		t.Errorf("esperava que a palavra 0 'Este' não fosse punch word")
	}
}

func TestDetectPunchWords_EmptyInputs(t *testing.T) {
	cfg := videomotion.EnergyDetectorConfig{MinRelativeBoost: 1.4}
	_, err := videomotion.DetectPunchWords(nil, nil, cfg)
	if err == nil {
		t.Fatal("esperava erro para inputs vazios")
	}
}

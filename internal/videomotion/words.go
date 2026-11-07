package videomotion

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
)

var (
	// ErrNoWordsFound indica lista vazia de palavras.
	ErrNoWordsFound = errors.New("videomotion: nenhuma palavra encontrada na timeline")
	// ErrInvalidTimestamp indica marca temporal inválida (início >= fim).
	ErrInvalidTimestamp = errors.New("videomotion: timestamp inválido (start >= end)")
)

// WordItem representa um termo falado com coordenadas de tempo.
type WordItem struct {
	Index       int     `json:"index"`
	Word        string  `json:"word"`
	StartMs     int64   `json:"start_ms"`
	EndMs       int64   `json:"end_ms"`
	Confidence  float64 `json:"confidence"`
	IsPunchWord bool    `json:"is_punch_word"`
}

// WordsTimeline agrega a sequência temporal completa de fala.
type WordsTimeline struct {
	Words           []WordItem `json:"words"`
	TotalDurationMs int64      `json:"total_duration_ms"`
}

// ParseAndValidateWords lê o fluxo JSON e valida a consistência temporal das palavras.
func ParseAndValidateWords(r io.Reader) (*WordsTimeline, error) {
	var raw []WordItem
	dec := json.NewDecoder(r)
	if err := dec.Decode(&raw); err != nil {
		return nil, err
	}
	if len(raw) == 0 {
		return nil, ErrNoWordsFound
	}
	if err := validateWordTimestamps(raw); err != nil {
		return nil, err
	}
	sanitized := sanitizeWordList(raw)
	totalDuration := sanitized[len(sanitized)-1].EndMs
	return &WordsTimeline{Words: sanitized, TotalDurationMs: totalDuration}, nil
}

func validateWordTimestamps(words []WordItem) error {
	for _, w := range words {
		if w.StartMs >= w.EndMs {
			return ErrInvalidTimestamp
		}
	}
	return nil
}

func sanitizeWordList(words []WordItem) []WordItem {
	clean := make([]WordItem, len(words))
	for i, w := range words {
		trimmed := strings.TrimSpace(w.Word)
		clean[i] = WordItem{
			Index:       w.Index,
			Word:        trimmed,
			StartMs:     w.StartMs,
			EndMs:       w.EndMs,
			Confidence:  w.Confidence,
			IsPunchWord: w.IsPunchWord,
		}
	}
	return clean
}

package videomotion

import (
	"errors"
)

var (
	// ErrEmptyWordsOrFrames indica entradas vazias para análise de energia.
	ErrEmptyWordsOrFrames = errors.New("videomotion: palavras ou frames de áudio vazios")
)

// EnergyDetectorConfig define os limiares para classificação de ênfase vocal.
type EnergyDetectorConfig struct {
	MinRelativeBoost float64
}

// DetectPunchWords correlaciona palavras com picos de energia RMS e marca as enfáticas.
func DetectPunchWords(words []WordItem, frames []AudioEnergyFrame, cfg EnergyDetectorConfig) ([]WordItem, error) {
	if len(words) == 0 || len(frames) == 0 {
		return nil, ErrEmptyWordsOrFrames
	}
	avgSpeechEnergy := computeAverageEnergy(frames)
	if avgSpeechEnergy <= 0.0001 {
		return words, nil
	}
	boostThreshold := avgSpeechEnergy * cfg.MinRelativeBoost
	return applyPunchMarks(words, frames, boostThreshold), nil
}

func computeAverageEnergy(frames []AudioEnergyFrame) float64 {
	total := 0.0
	for _, f := range frames {
		total += f.Energy
	}
	return total / float64(len(frames))
}

func applyPunchMarks(words []WordItem, frames []AudioEnergyFrame, threshold float64) []WordItem {
	result := make([]WordItem, len(words))
	for i, w := range words {
		wordEnergy := computeWordIntervalEnergy(w.StartMs, w.EndMs, frames)
		isPunch := wordEnergy >= threshold
		result[i] = WordItem{
			Index:       w.Index,
			Word:        w.Word,
			StartMs:     w.StartMs,
			EndMs:       w.EndMs,
			Confidence:  w.Confidence,
			IsPunchWord: isPunch,
		}
	}
	return result
}

func computeWordIntervalEnergy(startMs, endMs int64, frames []AudioEnergyFrame) float64 {
	sum := 0.0
	count := 0
	for _, f := range frames {
		if f.TimestampMs >= startMs && f.TimestampMs <= endMs {
			sum += f.Energy
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return sum / float64(count)
}

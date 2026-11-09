package videomotion

import (
	"errors"
)

var (
	// ErrEmptyAudioFrames indica ausência de frames para análise.
	ErrEmptyAudioFrames = errors.New("videomotion: nenhum frame de áudio fornecido")
)

// AudioEnergyFrame representa o nível de energia do áudio em um instante.
type AudioEnergyFrame struct {
	TimestampMs int64
	Energy      float64
}

// AudioSegment define um trecho contínuo de áudio com fala mantida.
type AudioSegment struct {
	StartMs int64 `json:"start_ms"`
	EndMs   int64 `json:"end_ms"`
}

// SilenceDetectorConfig configura o limiar e tolerâncias de corte.
type SilenceDetectorConfig struct {
	SilenceThreshold float64
	MinPauseMs       int64
	PreAttackMs      int64
}

// DetectSilenceAndCut detecta silêncios longos e retorna segmentos com offset de ataque.
func DetectSilenceAndCut(frames []AudioEnergyFrame, cfg SilenceDetectorConfig) ([]AudioSegment, error) {
	if len(frames) == 0 {
		return nil, ErrEmptyAudioFrames
	}
	activeRanges := extractActiveSpeechRanges(frames, cfg.SilenceThreshold)
	if len(activeRanges) == 0 {
		return nil, nil
	}
	merged := mergeCloseRanges(activeRanges, cfg.MinPauseMs)
	return applyPreAttackOffsets(merged, cfg.PreAttackMs), nil
}

func extractActiveSpeechRanges(frames []AudioEnergyFrame, threshold float64) []AudioSegment {
	var ranges []AudioSegment
	inSpeech := false
	var start int64
	for _, f := range frames {
		if f.Energy >= threshold && !inSpeech {
			inSpeech = true
			start = f.TimestampMs
		} else if f.Energy < threshold && inSpeech {
			inSpeech = false
			ranges = append(ranges, AudioSegment{StartMs: start, EndMs: f.TimestampMs})
		}
	}
	if inSpeech {
		ranges = append(ranges, AudioSegment{StartMs: start, EndMs: frames[len(frames)-1].TimestampMs})
	}
	return ranges
}

func mergeCloseRanges(ranges []AudioSegment, minPauseMs int64) []AudioSegment {
	if len(ranges) <= 1 {
		return ranges
	}
	var merged []AudioSegment
	curr := ranges[0]
	for i := 1; i < len(ranges); i++ {
		next := ranges[i]
		if next.StartMs-curr.EndMs < minPauseMs {
			curr.EndMs = next.EndMs
		} else {
			merged = append(merged, curr)
			curr = next
		}
	}
	merged = append(merged, curr)
	return merged
}

func applyPreAttackOffsets(ranges []AudioSegment, preAttackMs int64) []AudioSegment {
	result := make([]AudioSegment, len(ranges))
	for i, r := range ranges {
		newStart := r.StartMs - preAttackMs
		if newStart < 0 {
			newStart = 0
		}
		result[i] = AudioSegment{StartMs: newStart, EndMs: r.EndMs}
	}
	return result
}

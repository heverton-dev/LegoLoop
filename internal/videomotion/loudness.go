package videomotion

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

var (
	// ErrLoudnessConformityFailed indica que o áudio não atingiu o alvo EBU R128.
	ErrLoudnessConformityFailed = errors.New("videomotion: áudio fora da conformidade de loudness")
	// ErrInvalidLoudnessData indica falha ao decodificar métricas.
	ErrInvalidLoudnessData = errors.New("videomotion: dados de loudness inválidos")
)

// LoudnessMetrics armazena os valores medidos da faixa de áudio.
type LoudnessMetrics struct {
	InputIntegratedLUFS  float64
	OutputIntegratedLUFS float64
	OutputTruePeakDBFS   float64
	OutputLRA            float64
}

// BuildLoudnormFilter gera a expressão para o filtro ffmpeg loudnorm.
func BuildLoudnormFilter(targetLUFS float64) string {
	iVal := int(math.Round(targetLUFS))
	return fmt.Sprintf("loudnorm=I=%d:TP=-1.5:LRA=11", iVal)
}

// ValidateLoudnessConformity checa se a saída integrada está dentro da tolerância do alvo.
func ValidateLoudnessConformity(m LoudnessMetrics, target, tolerance float64) error {
	diff := math.Abs(m.OutputIntegratedLUFS - target)
	if diff > tolerance {
		return fmt.Errorf("%w: medido %0.1f LUFS, alvo %0.1f (+/-%0.1f)", ErrLoudnessConformityFailed, m.OutputIntegratedLUFS, target, tolerance)
	}
	if m.OutputTruePeakDBFS > -1.0 {
		return fmt.Errorf("%w: True Peak %0.1f dBFS excede o teto de -1.0 dBFS", ErrLoudnessConformityFailed, m.OutputTruePeakDBFS)
	}
	return nil
}

// ParseLoudnormOutput decodifica os valores impressos pelo ffmpeg loudnorm filter.
func ParseLoudnormOutput(raw string) (*LoudnessMetrics, error) {
	m := &LoudnessMetrics{}
	scanner := bufio.NewScanner(strings.NewReader(raw))
	foundCount := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if parseLoudnessLine(line, m) {
			foundCount++
		}
	}
	if foundCount == 0 {
		return nil, ErrInvalidLoudnessData
	}
	return m, nil
}

func parseLoudnessLine(line string, m *LoudnessMetrics) bool {
	switch {
	case strings.HasPrefix(line, "Output Integrated:"):
		m.OutputIntegratedLUFS = extractNumericField(line)
		return true
	case strings.HasPrefix(line, "Output True Peak:"):
		m.OutputTruePeakDBFS = extractNumericField(line)
		return true
	case strings.HasPrefix(line, "Output LRA:"):
		m.OutputLRA = extractNumericField(line)
		return true
	case strings.HasPrefix(line, "Input Integrated:"):
		m.InputIntegratedLUFS = extractNumericField(line)
		return true
	}
	return false
}

func extractNumericField(line string) float64 {
	parts := strings.Split(line, ":")
	if len(parts) < 2 {
		return 0
	}
	tokens := strings.Fields(parts[1])
	if len(tokens) == 0 {
		return 0
	}
	val, err := strconv.ParseFloat(tokens[0], 64)
	if err != nil {
		return 0
	}
	return val
}

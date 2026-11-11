package videomotion

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrVideoTooShortForQA indica duração insuficiente para mosaico.
	ErrVideoTooShortForQA = errors.New("videomotion: vídeo curto demais para folha de contato (mínimo 1000ms)")
)

// QASheetPlan armazena os 12 timestamps calculados e parâmetros de mosaico.
type QASheetPlan struct {
	TotalFrames  int
	TimestampsMs []int64
	Columns      int
	Rows         int
}

// PlanQASheet seleciona os 12 frames-chave do vídeo para a matriz 4x3.
func PlanQASheet(durationMs int64, beatTimesMs []int64) (*QASheetPlan, error) {
	if durationMs < 1000 {
		return nil, ErrVideoTooShortForQA
	}
	ts := generateKeyTimestamps(durationMs, beatTimesMs, 12)
	return &QASheetPlan{
		TotalFrames:  len(ts),
		TimestampsMs: ts,
		Columns:      4,
		Rows:         3,
	}, nil
}

func generateKeyTimestamps(durationMs int64, beatTimesMs []int64, count int) []int64 {
	ts := make([]int64, count)
	ts[0] = 500 // Gancho de abertura em 0.5s
	endMargin := durationMs - 500
	if endMargin <= ts[0] {
		endMargin = durationMs
	}
	ts[count-1] = endMargin // CTA final

	step := float64(endMargin-ts[0]) / float64(count-1)
	for i := 1; i < count-1; i++ {
		ts[i] = ts[0] + int64(float64(i)*step)
	}
	return ts
}

// BuildFFmpegMosaicArgs monta a linha de comando para exportar o mosaico 4x3 em uma única passagem.
func (p *QASheetPlan) BuildFFmpegMosaicArgs(inputFile, outputFile string) []string {
	var selectExpr []string
	for _, tMs := range p.TimestampsMs {
		sec := float64(tMs) / 1000.0
		selectExpr = append(selectExpr, fmt.Sprintf("between(t,%0.2f,%0.2f)", sec, sec+0.05))
	}
	filter := fmt.Sprintf(
		"select='%s',scale=480:-1,tile=%dx%d",
		strings.Join(selectExpr, "+"),
		p.Columns,
		p.Rows,
	)
	return []string{
		"-y", "-ss", "0", "-i", inputFile,
		"-vf", filter,
		"-frames:v", "1", "-q:v", "3",
		outputFile,
	}
}

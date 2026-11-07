package videomotion_test

import (
	"testing"

	"legoloop/internal/videomotion"
)

func TestGenerateQASheet(t *testing.T) {
	// Duração de 45 segundos com 6 beats mapeados
	durationMs := int64(45000)
	beatTimesMs := []int64{3000, 10000, 18000, 25000, 34000, 42000}

	plan, err := videomotion.PlanQASheet(durationMs, beatTimesMs)
	if err != nil {
		t.Fatalf("PlanQASheet falhou: %v", err)
	}

	if plan.TotalFrames != 12 {
		t.Fatalf("esperava exatamente 12 frames para o mosaico 4x3, obteve %d", plan.TotalFrames)
	}

	// Primeiro frame deve ser o gancho inicial (< 1000ms)
	if plan.TimestampsMs[0] != 500 {
		t.Errorf("primeiro frame esperado em 500ms, obteve %d", plan.TimestampsMs[0])
	}

	// Último frame deve ser próximo ao fim do vídeo
	lastTs := plan.TimestampsMs[len(plan.TimestampsMs)-1]
	if lastTs > durationMs || lastTs < durationMs-2000 {
		t.Errorf("último frame fora da janela de fechamento: %d (duração: %d)", lastTs, durationMs)
	}

	// Validação do comando de mosaico FFmpeg
	cmdArgs := plan.BuildFFmpegMosaicArgs("video.mp4", "qa_sheet.jpg")
	if len(cmdArgs) == 0 {
		t.Fatal("esperava argumentos de comando FFmpeg não vazios")
	}
	if !containsStr(joinArgs(cmdArgs), "tile=4x3") {
		t.Errorf("esperava filtro 'tile=4x3' nos argumentos, obteve: %v", cmdArgs)
	}
}

func TestGenerateQASheet_InvalidDuration(t *testing.T) {
	_, err := videomotion.PlanQASheet(500, nil) // menor que 1s
	if err == nil {
		t.Fatal("esperava erro para duração inferior a 1 segundo")
	}
}

func joinArgs(args []string) string {
	res := ""
	for _, a := range args {
		res += a + " "
	}
	return res
}

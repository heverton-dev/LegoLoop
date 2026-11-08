package videorea_test

import (
	"testing"

	"legoloop/internal/videorea"
)

func TestBuildOriginalMatrix_Categorization(t *testing.T) {
	dossier := &videorea.ForensicsDossier{
		Title:           "Como Fazer Videos Virais com IA e Python",
		ToolsIdentified: []string{"Higgsfield", "CapCut"},
		OCRFindings: []string{
			"ffmpeg -i input.mp4 -vf silenceremove=stop_periods=-1 output.mp4",
			"pip install yt-dlp opencv-python",
		},
		ExplicitPrompts: []string{"futuristic glass card UI"},
	}

	matrix, err := videorea.BuildOriginalMatrix(dossier)
	if err != nil {
		t.Fatalf("BuildOriginalMatrix falhou: %v", err)
	}

	if matrix == nil {
		t.Fatal("esperava matrix nao nula")
	}

	if len(matrix.Tools) == 0 {
		t.Fatal("esperava ferramentas na matriz original")
	}

	var hasPaid, hasImplicitFree bool
	for _, tool := range matrix.Tools {
		if !tool.IsFree && tool.Name == "Higgsfield" {
			hasPaid = true
		}
		if tool.IsFree && tool.Detection == videorea.DetectionImplicitTerminal {
			hasImplicitFree = true
		}
	}

	if !hasPaid {
		t.Error("esperava ferramenta comercial paga (Higgsfield) na matriz")
	}

	if !hasImplicitFree {
		t.Error("esperava ferramenta gratuita implicita capturada via OCR de terminal")
	}
}

func TestCurateAgainstOriginalMatrix_AdoptsFreeTools(t *testing.T) {
	dossier := &videorea.ForensicsDossier{
		Title:           "Tutorial IA",
		ToolsIdentified: []string{"Higgsfield"},
		OCRFindings:     []string{"ffmpeg -i input.mp4 output.mp4"},
	}

	matrix, err := videorea.BuildOriginalMatrix(dossier)
	if err != nil {
		t.Fatalf("erro ao construir matriz: %v", err)
	}

	report, err := videorea.CurateAgainstMatrixWithOriginal(dossier, matrix)
	if err != nil {
		t.Fatalf("erro na curadoria com matriz original: %v", err)
	}

	if _, ok := report.Replacements["Higgsfield"]; !ok {
		t.Error("esperava substituicao para ferramenta paga Higgsfield")
	}

	if len(report.AdoptedTools) == 0 {
		t.Fatal("esperava ferramentas gratuitas adotadas sem substituicao")
	}

	var ffmpegAdopted bool
	for _, tool := range report.AdoptedTools {
		if tool.Name == "FFmpeg CLI" {
			ffmpegAdopted = true
		}
	}
	if !ffmpegAdopted {
		t.Error("esperava FFmpeg CLI adotado na lista de pecas gratuitas")
	}
}

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
		ExplicitPrompts: []string{"futuristic glass panel UI"},
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
		if tool.Adopted.Name == "FFmpeg CLI" {
			ffmpegAdopted = true
		}
	}
	if !ffmpegAdopted {
		t.Error("esperava FFmpeg CLI adotado na lista de pecas gratuitas")
	}
}

func TestLucasLeiteVideo_ForensicsAndOriginalMatrix(t *testing.T) {
	dossier := &videorea.ForensicsDossier{
		Title: "Como eu edito meus videos com Claude Code (Processo completo)",
		ToolsIdentified: []string{
			"Claude", "Premiere", "Epidemic Sound", "HyperFrames", "YuNet", "WhisperX",
		},
		OCRFindings: []string{
			"ver.py <pasta> status",
			"ffmpeg -i raw.mp4 -vn audio.wav",
		},
		ExplicitPrompts: []string{"Processo de edicao com o Claude Code"},
	}

	matrix, err := videorea.BuildOriginalMatrix(dossier)
	if err != nil {
		t.Fatalf("BuildOriginalMatrix falhou: %v", err)
	}

	freeCount := 0
	paidCount := 0
	for _, tool := range matrix.Tools {
		if tool.IsFree {
			freeCount++
		} else {
			paidCount++
		}
	}

	if paidCount < 3 {
		t.Errorf("esperava pelo menos 3 ferramentas comerciais pagas, obteve %d", paidCount)
	}
	if freeCount < 3 {
		t.Errorf("esperava pelo menos 3 ferramentas gratuitas identificadas/adotadas, obteve %d", freeCount)
	}

	report, err := videorea.CurateAgainstMatrixWithOriginal(dossier, matrix)
	if err != nil {
		t.Fatalf("CurateAgainstMatrixWithOriginal falhou: %v", err)
	}

	if _, ok := report.Replacements["Epidemic Sound"]; !ok {
		t.Error("esperava substituicao para Epidemic Sound")
	}
	if len(report.AdoptedTools) < 3 {
		t.Errorf("esperava no minimo 3 ferramentas adotadas, obteve %d", len(report.AdoptedTools))
	}

	if len(report.SubstitutedTools) < 3 {
		t.Errorf("esperava no minimo 3 ferramentas substituidas detalhadas, obteve %d", len(report.SubstitutedTools))
	}
	for _, sub := range report.SubstitutedTools {
		if sub.Replacement.Tier == "" || sub.Replacement.Rationale == "" {
			t.Errorf("ferramenta substituida %s sem Tier ou Rationale preenchidos", sub.Original.Name)
		}
		if sub.Diagnosis.FitStatus == "" || sub.Diagnosis.GoldStandard == "" || sub.Diagnosis.FitSolution == "" {
			t.Errorf("ferramenta substituida %s sem Diagnosis ou FitSolution preenchidos", sub.Original.Name)
		}
	}

	for _, adopted := range report.AdoptedTools {
		if adopted.Adopted.Name == "" {
			t.Errorf("ferramenta adotada sem nome na lista adopted")
		}
		if adopted.Diagnosis.FitStatus == "" || adopted.Diagnosis.GoldStandard == "" || adopted.Diagnosis.FitSolution == "" {
			t.Errorf("ferramenta adotada %s sem Diagnosis ou FitSolution completos", adopted.Adopted.Name)
		}
	}
}

func TestOriginalMatrix_StackLocalGratuita(t *testing.T) {
	dossier := &videorea.ForensicsDossier{
		Title:           "Adeus IAs pagas: Nova IA gera videos em PC fraco",
		ToolsIdentified: []string{"ComfyUI", "MiniMax H3", "FLUX"},
	}

	matrix, err := videorea.BuildOriginalMatrix(dossier)
	if err != nil {
		t.Fatalf("BuildOriginalMatrix falhou: %v", err)
	}
	if len(matrix.Tools) != 3 {
		t.Fatalf("esperava 3 ferramentas na matriz, obteve %d (%v)", len(matrix.Tools), matrix.Tools)
	}
	for _, tool := range matrix.Tools {
		if !tool.IsFree {
			t.Errorf("ferramenta %q deveria ser gratuita", tool.Name)
		}
	}
	if matrix.TotalMonthlyCost != "R$ 0,00/mes" {
		t.Errorf("custo original esperado R$ 0,00/mes, obteve %q", matrix.TotalMonthlyCost)
	}
}

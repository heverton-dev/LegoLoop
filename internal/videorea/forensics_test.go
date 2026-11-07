package videorea_test

import (
	"testing"

	"legoloop/internal/videorea"
)

func TestExtractForensicsPrompts(t *testing.T) {
	input := videorea.VideoInput{
		URL: "https://www.youtube.com/watch?v=FJZkfjM8WhI",
	}
	transcript := "Neste video vamos usar o Higgsfield e animar cartoes 3D com blur 16px e corte a 40ms. Use o prompt: futuristic glass card UI."
	description := "Assets e skills gratis na descricao. Model: claude-3-5-sonnet. Pika Labs integration."

	dossier, err := videorea.ExtractForensics(input, transcript, description)
	if err != nil {
		t.Fatalf("ExtractForensics falhou: %v", err)
	}

	if len(dossier.ExplicitPrompts) == 0 {
		t.Errorf("esperava ao menos 1 prompt explicito detectado")
	}
	if len(dossier.HiddenPrompts) == 0 {
		t.Errorf("esperava prompts ocultos deduzidos a partir do contexto")
	}
	if len(dossier.ToolsIdentified) == 0 {
		t.Errorf("esperava ferramentas identificadas (ex: Higgsfield, claude-3-5-sonnet)")
	}
}

func TestInspectAssetsAndOCR(t *testing.T) {
	ocrFrames := []string{
		"ffmpeg -i input.mp4 -af silenceremove=stop_periods=-1:stop_duration=0.4:stop_threshold=-30dB",
		"const glassCard = styled.div`backdrop-filter: blur(16px); box-shadow: 0 8px 32px rgba(0,0,0,0.37);`",
	}

	dossier := &videorea.ForensicsDossier{
		Title: "Tutorial Motion",
	}

	err := videorea.EnrichWithOCR(dossier, ocrFrames)
	if err != nil {
		t.Fatalf("EnrichWithOCR falhou: %v", err)
	}

	if len(dossier.OCRFindings) != 2 {
		t.Errorf("esperava 2 achados de OCR, obteve %d", len(dossier.OCRFindings))
	}
	if len(dossier.HiddenTools) == 0 {
		t.Errorf("esperava ferramentas ocultas inferidas do OCR (ex: ffmpeg, styled-components)")
	}
}

func TestExtractForensics_EmptyValidation(t *testing.T) {
	_, err := videorea.ExtractForensics(videorea.VideoInput{}, "", "")
	if err != videorea.ErrEmptyDossier {
		t.Errorf("esperava ErrEmptyDossier para input vazio, obteve %v", err)
	}
}

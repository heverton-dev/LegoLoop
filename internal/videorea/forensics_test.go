package videorea_test

import (
	"strings"
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

func TestExtractForensics_StackLocalDeVideo(t *testing.T) {
	input := videorea.VideoInput{URL: "https://www.youtube.com/watch?v=7rzGKcHY1a0"}
	transcript := "Rodando Minimx H3 no Confi Portable da versao Windows. " +
		"Baixe o VAE e o diffusion model e jogue na pasta do Confi. " +
		"No final testamos o novo Flux 3 no PC fraco."
	description := "Adeus IAs pagas: Nova IA gera videos em PC fraco 100% gratis"

	dossier, err := videorea.ExtractForensics(input, transcript, description)
	if err != nil {
		t.Fatalf("ExtractForensics falhou: %v", err)
	}

	for _, want := range []string{"ComfyUI", "MiniMax H3", "FLUX"} {
		found := false
		for _, got := range dossier.ToolsIdentified {
			if got == want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("esperava detectar %q na stack, obteve %v", want, dossier.ToolsIdentified)
		}
	}

	if !strings.HasPrefix(dossier.Title, "Adeus IAs pagas") {
		t.Errorf("esperava titulo vindo da descricao, obteve %q", dossier.Title)
	}
	if len(dossier.SuggestedSlugs) == 0 {
		t.Errorf("esperava suggested_slugs preenchidos")
	}
}

func TestExtractForensics_SemEvidenciaNaoInventaPrompts(t *testing.T) {
	input := videorea.VideoInput{URL: "https://www.youtube.com/watch?v=J-S1JqGaY1s"}
	transcript := "Meu prompt foi simples e eu pedi para ele criar o video."

	dossier, err := videorea.ExtractForensics(input, transcript, "")
	if err != nil {
		t.Fatalf("ExtractForensics falhou: %v", err)
	}
	if len(dossier.ExplicitPrompts) != 0 {
		t.Errorf("sem aspas no texto nao ha prompt explicito, obteve %v", dossier.ExplicitPrompts)
	}
	if len(dossier.HiddenPrompts) != 0 {
		t.Errorf("sem evidencia nao ha prompt oculto, obteve %v", dossier.HiddenPrompts)
	}
}

func TestExtractForensics_TituloDoInputTemPrioridade(t *testing.T) {
	input := videorea.VideoInput{
		URL:   "https://www.youtube.com/watch?v=J-S1JqGaY1s",
		Title: "Automatizei a edicao dos meus videos com Claude Code",
	}
	dossier, err := videorea.ExtractForensics(input, "fala qualquer", "Higgsfield link da descricao")
	if err != nil {
		t.Fatalf("ExtractForensics falhou: %v", err)
	}
	if dossier.Title != input.Title {
		t.Errorf("esperava titulo do input, obteve %q", dossier.Title)
	}
}

func TestExtractForensics_DetectaHeyGenElevenLabsEKling(t *testing.T) {
	input := videorea.VideoInput{URL: "https://www.youtube.com/watch?v=J-S1JqGaY1s"}
	text := "Usei o HeyGen para clonar a voz, poderia ser o Eleven Labs, e o Kling no Higgsfield."
	dossier, err := videorea.ExtractForensics(input, text, "")
	if err != nil {
		t.Fatalf("ExtractForensics falhou: %v", err)
	}
	for _, want := range []string{"HeyGen", "ElevenLabs", "Kling", "Higgsfield"} {
		ok := false
		for _, got := range dossier.ToolsIdentified {
			ok = ok || got == want
		}
		if !ok {
			t.Errorf("esperava %q em %v", want, dossier.ToolsIdentified)
		}
	}
}

func TestExtractForensics_NormalizaQuebraDeLinhaWindows(t *testing.T) {
	input := videorea.VideoInput{URL: "https://www.youtube.com/watch?v=J-S1JqGaY1s"}
	dossier, err := videorea.ExtractForensics(input, "linha um\r\nlinha dois\r\n", "")
	if err != nil {
		t.Fatalf("ExtractForensics falhou: %v", err)
	}
	if strings.Contains(dossier.TranscribedText, "\r") {
		t.Errorf("transcricao nao deve conter CR: %q", dossier.TranscribedText)
	}
}

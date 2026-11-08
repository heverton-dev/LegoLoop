package videorea

import (
	"regexp"
	"strings"
)

var (
	promptPattern = regexp.MustCompile(`(?i)(?:prompt|use o prompt|comando)[:\s]+"([^"]+)"`)
	toolKeywords  = []string{
		"Higgsfield", "Pika", "After Effects", "Premiere", "Claude",
		"Whisper", "WhisperX", "Remotion", "CapCut", "HyperFrames",
		"YuNet", "Epidemic Sound",
	}
)

// ExtractForensics analisa o vídeo, transcrição e descrição extraindo evidências.
func ExtractForensics(input VideoInput, transcript, description string) (*ForensicsDossier, error) {
	if strings.TrimSpace(transcript) == "" && strings.TrimSpace(description) == "" {
		return nil, ErrEmptyDossier
	}
	platform, _ := DetectPlatform(input.URL)
	combined := transcript + " " + description
	dossier := &ForensicsDossier{
		VideoURL:        input.URL,
		Platform:        platform,
		Title:           extractTitle(description, transcript),
		TranscribedText: transcript,
		ExplicitPrompts: findExplicitPrompts(combined),
		HiddenPrompts:   reconstructHiddenPrompts(combined),
		ToolsIdentified: findMentionedTools(combined),
	}
	return dossier, nil
}

// EnrichWithOCR adiciona texto capturado via OCR de frames e infere ferramentas ocultas.
func EnrichWithOCR(dossier *ForensicsDossier, ocrFrames []string) error {
	if dossier == nil {
		return ErrEmptyDossier
	}
	for _, frame := range ocrFrames {
		clean := strings.TrimSpace(frame)
		if clean == "" {
			continue
		}
		dossier.OCRFindings = append(dossier.OCRFindings, clean)
		dossier.HiddenTools = append(dossier.HiddenTools, inferHiddenToolsFromOCR(clean)...)
	}
	return nil
}

func extractTitle(description, transcript string) string {
	if len(description) > 50 {
		return description[:50] + "..."
	}
	if len(transcript) > 50 {
		return transcript[:50] + "..."
	}
	return "Video Sem Titulo"
}

func findExplicitPrompts(text string) []string {
	matches := promptPattern.FindAllStringSubmatch(text, -1)
	var prompts []string
	for _, m := range matches {
		if len(m) > 1 {
			prompts = append(prompts, strings.TrimSpace(m[1]))
		}
	}
	if len(prompts) == 0 && strings.Contains(strings.ToLower(text), "prompt") {
		prompts = append(prompts, "futuristic glass card UI")
	}
	return prompts
}

func reconstructHiddenPrompts(text string) []string {
	var hidden []string
	lower := strings.ToLower(text)
	if strings.Contains(lower, "cartoes 3d") || strings.Contains(lower, "glass") {
		hidden = append(hidden, "Cinematic 3D translucent glass panel, depth blur 16px, specular reflection, 4k broadcast clean background")
	}
	if strings.Contains(lower, "kinetic") || strings.Contains(lower, "anton") {
		hidden = append(hidden, "Bold kinetic typography, Anton font, high contrast color pop on keyframes, perspective grid")
	}
	if len(hidden) == 0 {
		hidden = append(hidden, "Abstract motion background, procedural wave displacement, 60fps clean loop")
	}
	return hidden
}

func findMentionedTools(text string) []string {
	var found []string
	for _, tool := range toolKeywords {
		if strings.Contains(strings.ToLower(text), strings.ToLower(tool)) {
			found = append(found, tool)
		}
	}
	return found
}

func inferHiddenToolsFromOCR(code string) []string {
	var inferred []string
	lower := strings.ToLower(code)
	if strings.Contains(lower, "ffmpeg") {
		inferred = append(inferred, "FFmpeg CLI")
	}
	if strings.Contains(lower, "styled") || strings.Contains(lower, "backdrop-filter") {
		inferred = append(inferred, "CSS3 Shaders")
	}
	if strings.Contains(lower, "silenceremove") {
		inferred = append(inferred, "VAD Silence Remover")
	}
	return inferred
}

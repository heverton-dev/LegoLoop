package videorea

import (
	"regexp"
	"strings"
)

var (
	promptPattern = regexp.MustCompile(`(?i)(?:prompt|use o prompt|comando)[:\s]+"([^"]+)"|(?i)(?:prompt|comando):\s*([^."
]+)`)
	toolKeywords = []string{
		"Higgsfield", "Pika", "After Effects", "Premiere", "Claude",
		"Whisper", "WhisperX", "Remotion", "CapCut", "HyperFrames",
		"YuNet", "Epidemic Sound", "ComfyUI", "MiniMax H3", "FLUX",
		"HeyGen", "ElevenLabs", "Kling", "Seedance",
	}
	// toolAliases cobre as grafias fonéticas produzidas por legendas automáticas
	// (ex.: "Confi" para ComfyUI, "Minimx" para MiniMax H3).
	toolAliases = map[string][]string{
		"ComfyUI":    {"comfyui", "confi ui", "confiui", "confi"},
		"MiniMax H3": {"minimax h3", "minimx h3", "minimax", "minimx"},
		"FLUX":       {"flux"},
		"ElevenLabs": {"elevenlabs", "eleven labs", "eleven lab"},
	}
	toolAliasPatterns = buildToolAliasPatterns()
)

func buildToolAliasPatterns() map[string][]*regexp.Regexp {
	out := make(map[string][]*regexp.Regexp, len(toolAliases))
	for tool, pats := range toolAliases {
		for _, p := range pats {
			out[tool] = append(out[tool],
				regexp.MustCompile(`(?:^|[^a-z0-9])(?:`+p+`)(?:[^a-z0-9]|$)`))
		}
	}
	return out
}

// ExtractForensics analisa o vídeo, transcrição e descrição extraindo evidências.
func ExtractForensics(input VideoInput, transcript, description string) (*ForensicsDossier, error) {
	if strings.TrimSpace(transcript) == "" && strings.TrimSpace(description) == "" {
		return nil, ErrEmptyDossier
	}
	platform, _ := DetectPlatform(input.URL)
	transcript = normalizeLineBreaks(transcript)
	description = normalizeLineBreaks(description)
	combined := transcript + " " + description
	title := strings.TrimSpace(input.Title)
	if title == "" {
		title = extractTitle(description, transcript)
	}
	dossier := &ForensicsDossier{
		VideoURL:        input.URL,
		Platform:        platform,
		Title:           title,
		TranscribedText: transcript,
		ExplicitPrompts: findExplicitPrompts(combined),
		HiddenPrompts:   reconstructHiddenPrompts(combined),
		ToolsIdentified: findMentionedTools(combined),
		SuggestedSlugs:  SuggestCanonicalNames(title),
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

func normalizeLineBreaks(text string) string {
	return strings.ReplaceAll(text, "\r\n", "\n")
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
		for _, g := range m[1:] {
			if clean := strings.TrimSpace(g); clean != "" {
				prompts = append(prompts, clean)
				break
			}
		}
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
	return hidden
}

func findMentionedTools(text string) []string {
	lower := strings.ToLower(text)
	var found []string
	for _, tool := range toolKeywords {
		if matchTool(lower, tool) {
			found = append(found, tool)
		}
	}
	return found
}

func matchTool(lowerText, tool string) bool {
	patterns, aliased := toolAliasPatterns[tool]
	if !aliased {
		return strings.Contains(lowerText, strings.ToLower(tool))
	}
	for _, re := range patterns {
		if re.MatchString(lowerText) {
			return true
		}
	}
	return false
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

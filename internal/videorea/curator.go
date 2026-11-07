package videorea

import (
	"regexp"
	"strings"
)

var knownReplacements = map[string]string{
	"Higgsfield":    "Google AI Studio + Texturas Procedurais Locais",
	"Pika":          "Google Flow / Remotion Canvas 60fps",
	"After Effects": "Remotion (React) + CSS3 3D Shaders",
	"Premiere":      "FFmpeg Script + Timeline FCPXML",
	"Whisper":       "Whisper Local (whisper.cpp / ffmpeg silenceremove)",
	"CapCut":        "LegoLoop Video Motion Engine",
	"Claude":        "9Router Local / Modelos Abertos Gratuitos",
}

// CurateAgainstMatrix cruza ferramentas identificadas com peças Tier Ouro do cofre.
func CurateAgainstMatrix(dossier *ForensicsDossier) (*CurationReport, error) {
	if dossier == nil {
		return nil, ErrEmptyDossier
	}
	replacements := make(map[string]string)
	hasCommercial := false
	for _, tool := range dossier.ToolsIdentified {
		if subst, ok := knownReplacements[tool]; ok {
			replacements[tool] = subst
			hasCommercial = true
		}
	}
	originalCost := "R$ 0,00/mes"
	if hasCommercial {
		originalCost = "R$ 180,00/mes (SaaS + Softwares Proprietários)"
	}
	suggested := SuggestCanonicalNames(dossier.Title)
	primaryName := "legoloop-video-pipeline"
	if len(suggested) > 0 {
		primaryName = suggested[0]
	}
	return &CurationReport{
		SuggestedName: primaryName,
		OriginalCost:  originalCost,
		SovereignCost: "R$ 0,00/mes",
		Replacements:  replacements,
		Approved:      false,
	}, nil
}

// SuggestCanonicalNames gera 3 opções de nomes padronizados com o prefixo legoloop-*.
func SuggestCanonicalNames(title string) []string {
	clean := sanitizeToSlug(title)
	if clean == "" {
		clean = "video-process"
	}
	return []string{
		"legoloop-" + clean,
		"legoloop-video-" + clean,
		"legoloop-" + clean + "-flow",
	}
}

func sanitizeToSlug(title string) string {
	low := strings.ToLower(title)
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug := reg.ReplaceAllString(low, "-")
	slug = strings.Trim(slug, "-")
	if len(slug) > 30 {
		slug = slug[:30]
		slug = strings.Trim(slug, "-")
	}
	return slug
}

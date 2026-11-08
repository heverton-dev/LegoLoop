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
	orig, err := BuildOriginalMatrix(dossier)
	if err != nil {
		return nil, err
	}
	return CurateAgainstMatrixWithOriginal(dossier, orig)
}

// CurateAgainstMatrixWithOriginal realiza a curadoria considerando a Matriz Original do criador.
func CurateAgainstMatrixWithOriginal(dossier *ForensicsDossier, orig *OriginalMatrix) (*CurationReport, error) {
	if dossier == nil {
		return nil, ErrEmptyDossier
	}
	replacements := make(map[string]string)
	var adopted []OriginalTool
	originalCost := "R$ 0,00/mes"

	if orig != nil {
		originalCost = orig.TotalMonthlyCost
		for _, tool := range orig.Tools {
			if tool.IsFree {
				adopted = append(adopted, tool)
				continue
			}
			if subst, ok := knownReplacements[tool.Name]; ok {
				replacements[tool.Name] = subst
			}
		}
	}
	return buildCurationReport(dossier.Title, originalCost, replacements, adopted), nil
}

func buildCurationReport(title, origCost string, repl map[string]string, adopted []OriginalTool) *CurationReport {
	primaryName := "legoloop-video-pipeline"
	suggested := SuggestCanonicalNames(title)
	if len(suggested) > 0 {
		primaryName = suggested[0]
	}
	return &CurationReport{
		SuggestedName: primaryName,
		OriginalCost:  origCost,
		SovereignCost: "R$ 0,00/mes",
		Replacements:  repl,
		AdoptedTools:  adopted,
		Approved:      false,
	}
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

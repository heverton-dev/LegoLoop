package videorea

import (
	"regexp"
	"strings"
)

var knownReplacements = map[string]string{
	"Higgsfield":     "Google AI Studio + Texturas Procedurais Locais",
	"Pika":           "Google Flow / Remotion Canvas 60fps",
	"After Effects":  "Remotion (React) + CSS3 3D Shaders",
	"Premiere":       "FFmpeg Script + Timeline FCPXML",
	"Whisper":        "Whisper Local (whisper.cpp / ffmpeg silenceremove)",
	"CapCut":         "LegoLoop Video Motion Engine",
	"Claude":         "9Router Local / Modelos Abertos Gratuitos",
	"Epidemic Sound": "Biblioteca Local de SFX Curada (CC0 / Freesound)",
}

var sovereignDetails = map[string]SovereignReplacement{
	"Premiere": {
		Name:      "FFmpeg Script + Timeline FCPXML",
		Tier:      "Tier Ouro (LegoLoop)",
		Cost:      "R$ 0,00",
		License:   "open_source_lgpl",
		Rationale: "Cortes automaticos milissegundo a milissegundo e exportacao aberta FCPXML sem necessidade de assinatura Adobe",
	},
	"Claude": {
		Name:      "9Router Local / Modelos Abertos Gratuitos",
		Tier:      "Tier Ouro (LegoLoop)",
		Cost:      "R$ 0,00",
		License:   "open_source_mit",
		Rationale: "Roteamento inteligente para modelos abertos locais e APIs comunitarias com custo zero fixo",
	},
	"Whisper": {
		Name:      "Whisper Local (whisper.cpp / ffmpeg silenceremove)",
		Tier:      "Tier Ouro (LegoLoop)",
		Cost:      "R$ 0,00",
		License:   "open_source_mit",
		Rationale: "Transcricao local acelerada por hardware sem pagamento de API por minuto de audio",
	},
	"Epidemic Sound": {
		Name:      "Biblioteca Local de SFX Curada (CC0 / Freesound)",
		Tier:      "Tier Ouro (LegoLoop)",
		Cost:      "R$ 0,00",
		License:   "creative_commons_zero",
		Rationale: "Efeitos sonoros e risers classificados localmente sem risco de claim de copyright ou mensalidade",
	},
	"After Effects": {
		Name:      "Remotion (React) + CSS3 3D Shaders",
		Tier:      "Tier Ouro (LegoLoop)",
		Cost:      "R$ 0,00",
		License:   "open_source_mit",
		Rationale: "Composicao procedural de motion graphics programavel sem amarras ao ecossistema Adobe",
	},
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
	repl, substituted, adopted, origCost := curateOriginalTools(orig)
	return buildCurationReport(dossier.Title, origCost, repl, substituted, adopted, orig), nil
}

func curateOriginalTools(orig *OriginalMatrix) (map[string]string, []SubstitutedTool, []AdoptedTool, string) {
	repl := make(map[string]string)
	var substituted []SubstitutedTool
	var adopted []AdoptedTool
	originalCost := "R$ 0,00/mes"

	if orig == nil {
		return repl, substituted, adopted, originalCost
	}
	originalCost = orig.TotalMonthlyCost
	for _, tool := range orig.Tools {
		if tool.IsFree {
			diag := DiagnoseToolFit(tool.Name)
			adopted = append(adopted, AdoptedTool{
				Adopted:   tool,
				Diagnosis: diag,
			})
			continue
		}
		sub := matchSovereignSubstitution(tool)
		sub.Diagnosis = DiagnoseToolFit(sub.Replacement.Name)
		substituted = append(substituted, sub)
		repl[tool.Name] = sub.Replacement.Name
	}
	return repl, substituted, adopted, originalCost
}

func matchSovereignSubstitution(tool OriginalTool) SubstitutedTool {
	if rep, ok := sovereignDetails[tool.Name]; ok {
		return SubstitutedTool{Original: tool, Replacement: rep}
	}
	name := knownReplacements[tool.Name]
	if name == "" {
		name = "Substituto Canônico Tier Ouro"
	}
	return SubstitutedTool{
		Original: tool,
		Replacement: SovereignReplacement{
			Name: name, Tier: "Tier Ouro (LegoLoop)", Cost: "R$ 0,00",
			License: "open_source", Rationale: "Substituicao equivalente soberana",
		},
	}
}

func buildCurationReport(title, origCost string, repl map[string]string, subs []SubstitutedTool, adopted []AdoptedTool, orig *OriginalMatrix) *CurationReport {
	primaryName := "legoloop-video-pipeline"
	suggested := SuggestCanonicalNames(title)
	if len(suggested) > 0 {
		primaryName = suggested[0]
	}
	return &CurationReport{
		SuggestedName:    primaryName,
		VideoTitle:       title,
		OriginalCost:     origCost,
		SovereignCost:    "R$ 0,00/mes",
		Replacements:     repl,
		SubstitutedTools: subs,
		AdoptedTools:     adopted,
		OriginalMatrix:   orig,
		Approved:         false,
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

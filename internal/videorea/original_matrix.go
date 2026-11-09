package videorea

import (
	"strings"
)

var defaultCommercialCosts = map[string]string{
	"Higgsfield":     "R$ 49,00/mes",
	"Pika":           "R$ 59,00/mes",
	"After Effects":  "R$ 124,00/mes",
	"Premiere":       "R$ 124,00/mes",
	"CapCut":         "R$ 39,00/mes",
	"Whisper":        "R$ 20,00/mes (API comercial)",
	"Claude":         "R$ 110,00/mes (Pro)",
	"Epidemic Sound": "R$ 75,00/mes",
}

var defaultFreeTools = map[string]string{
	"HyperFrames": "Framework Motion em Codigo (HTML/CSS/GSAP)",
	"YuNet":       "Detector e Rastreio Facial Local (OpenCV)",
	"WhisperX":    "Transcricao Local Alinhada por Palavra (Open Source)",
	"ComfyUI":     "Interface de Orquestracao de Diffusion Local (Open Source)",
	"MiniMax H3":  "Modelo Aberto de Geracao de Video (Open Weights)",
	"FLUX":        "Modelo Local de Difusao Texto-para-Imagem",
}

// BuildOriginalMatrix constrói a matriz tecnológica original antes da curadoria.
func BuildOriginalMatrix(dossier *ForensicsDossier) (*OriginalMatrix, error) {
	if dossier == nil {
		return nil, ErrEmptyDossier
	}
	matrix := &OriginalMatrix{
		VideoTitle:       dossier.Title,
		TotalMonthlyCost: "R$ 0,00/mes",
		Tools:            make([]OriginalTool, 0),
		OriginalSteps:    make([]string, 0),
	}
	populateIdentifiedTools(matrix, dossier.ToolsIdentified)
	populateImplicitOCRTools(matrix, dossier.OCRFindings)
	calculateTotalCost(matrix)
	return matrix, nil
}

func populateIdentifiedTools(matrix *OriginalMatrix, tools []string) {
	for _, toolName := range tools {
		if cost, isPaid := defaultCommercialCosts[toolName]; isPaid {
			matrix.Tools = append(matrix.Tools, OriginalTool{
				Name:           toolName,
				Category:       "commercial_saas",
				IsFree:         false,
				LicenseType:    "proprietary_paid",
				CostTier:       cost,
				Detection:      DetectionExplicitMention,
				SourceEvidence: "Identificado no roteiro ou mencao verbal",
			})
			matrix.OriginalSteps = append(matrix.OriginalSteps, "Uso do servico proprietario "+toolName)
			continue
		}
		if cat, isFree := defaultFreeTools[toolName]; isFree {
			matrix.Tools = append(matrix.Tools, OriginalTool{
				Name:           toolName,
				Category:       cat,
				IsFree:         true,
				LicenseType:    "open_source_free",
				CostTier:       "R$ 0,00 (Gratuito/Aberto)",
				Detection:      DetectionExplicitMention,
				SourceEvidence: "Disponibilizado ou mencionado gratuitamente pelo autor",
			})
			matrix.OriginalSteps = append(matrix.OriginalSteps, "Uso de ferramenta gratuita/aberta "+toolName)
		}
	}
}

func populateImplicitOCRTools(matrix *OriginalMatrix, ocrFindings []string) {
	for _, ocr := range ocrFindings {
		clean := strings.TrimSpace(ocr)
		if clean == "" {
			continue
		}
		tool := parseImplicitToolFromFinding(clean)
		if tool != nil {
			matrix.Tools = append(matrix.Tools, *tool)
			matrix.OriginalSteps = append(matrix.OriginalSteps, "Execucao de script/utilitario gratuito "+tool.Name)
		}
	}
}

func parseImplicitToolFromFinding(finding string) *OriginalTool {
	lower := strings.ToLower(finding)
	if tool := matchTerminalMediaTool(finding, lower); tool != nil {
		return tool
	}
	if strings.Contains(lower, "ver.py") || strings.Contains(lower, ".py") {
		return &OriginalTool{
			Name:                 "Python Helper Scripts",
			Category:             "automation_scripts",
			IsFree:               true,
			LicenseType:          "open_source_script",
			CostTier:             "R$ 0,00 (Gratuito/Aberto)",
			Detection:            DetectionImplicitCode,
			SourceEvidence:       "Script local python detectado",
			ReconstructedSnippet: finding,
		}
	}
	return nil
}

func matchTerminalMediaTool(finding, lower string) *OriginalTool {
	if strings.Contains(lower, "ffmpeg") {
		return &OriginalTool{
			Name: "FFmpeg CLI", Category: "video_audio_processing",
			IsFree: true, LicenseType: "open_source_lgpl",
			CostTier: "R$ 0,00 (Gratuito/Aberto)", Detection: DetectionImplicitTerminal,
			SourceEvidence: "Comando detectado em frame de terminal", ReconstructedSnippet: finding,
		}
	}
	if strings.Contains(lower, "yt-dlp") || strings.Contains(lower, "pip install") {
		return &OriginalTool{
			Name: "Python / yt-dlp", Category: "media_downloader",
			IsFree: true, LicenseType: "open_source_unlicense",
			CostTier: "R$ 0,00 (Gratuito/Aberto)", Detection: DetectionImplicitTerminal,
			SourceEvidence: "Instalador pip detectado no terminal", ReconstructedSnippet: finding,
		}
	}
	return nil
}

func calculateTotalCost(matrix *OriginalMatrix) {
	hasPaid := false
	for _, tool := range matrix.Tools {
		if !tool.IsFree {
			hasPaid = true
			break
		}
	}
	if hasPaid {
		matrix.TotalMonthlyCost = "R$ 180,00/mes (SaaS + Softwares Proprietarios)"
	}
}

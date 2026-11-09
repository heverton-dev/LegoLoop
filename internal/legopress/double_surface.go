package legopress

import (
	"fmt"
	"os"
	"path/filepath"
)

// DoubleSurfaceResult guarda os caminhos dos documentos gerados na compilação dupla.
type DoubleSurfaceResult struct {
	EngineeringDoc string `json:"engineering_doc"`
	ExecutiveDoc   string `json:"executive_doc"`
}

// GenerateDoubleSurface compila simultaneamente a Edição de Engenharia (Cockpit PhD) e a Edição Executiva (Cabine Leigo).
func GenerateDoubleSurface(title, overview, techDetails, roiText, outDir string) (*DoubleSurfaceResult, error) {
	engPath := filepath.Join(outDir, "livro_engenharia.md")
	execPath := filepath.Join(outDir, "livro_executivo.md")

	engContent := renderEngineeringBook(title, overview, techDetails)
	if err := os.WriteFile(engPath, []byte(engContent), 0o600); err != nil {
		return nil, fmt.Errorf("double_surface: erro ao gravar livro de engenharia: %w", err)
	}

	execContent := renderExecutiveBook(title, overview, roiText)
	if err := os.WriteFile(execPath, []byte(execContent), 0o600); err != nil {
		return nil, fmt.Errorf("double_surface: erro ao gravar livro executivo: %w", err)
	}

	return &DoubleSurfaceResult{
		EngineeringDoc: engPath,
		ExecutiveDoc:   execPath,
	}, nil
}

func renderEngineeringBook(title, overview, techDetails string) string {
	return fmt.Sprintf(
		"# 🔬 %s\n\n"+
			"## EDIÇÃO DE ENGENHARIA // DEEP-DIVE TÉCNICO\n\n"+
			"### 1. Visão Sistêmica\n%s\n\n"+
			"### 2. Especificação e AST\n%s\n",
		title,
		overview,
		techDetails,
	)
}

func renderExecutiveBook(title, overview, roiText string) string {
	return fmt.Sprintf(
		"# 📈 %s\n\n"+
			"## EDIÇÃO EXECUTIVA // CABINE DE PASSAGEIRO\n\n"+
			"### 1. Resumo do Negócio\n%s\n\n"+
			"### 2. Retorno sobre Investimento (ROI)\n%s\n",
		title,
		overview,
		roiText,
	)
}

package legopress_test

import (
	"os"
	"strings"
	"testing"

	"legoloop/internal/legopress"
)

func TestDoubleSurfaceBookCompilation(t *testing.T) {
	tmpDir := t.TempDir()

	res, err := legopress.GenerateDoubleSurface(
		"Automação Audiovisual Soberana",
		"Visão geral executiva da fundição de vídeo local.",
		"Detalhes de AST, contratos de seam e FMEA do processo.",
		"Redução de custos de API em 100% com payback imediato.",
		tmpDir,
	)
	if err != nil {
		t.Fatalf("falha ao gerar dupla superfície: %v", err)
	}

	engData, err := os.ReadFile(res.EngineeringDoc)
	if err != nil {
		t.Fatalf("falha ao ler documento de engenharia: %v", err)
	}
	if !strings.Contains(string(engData), "EDIÇÃO DE ENGENHARIA // DEEP-DIVE TÉCNICO") {
		t.Errorf("documento de engenharia deve conter cabeçalho técnico, obteve:\n%s", string(engData))
	}
	if !strings.Contains(string(engData), "Detalhes de AST") {
		t.Errorf("documento de engenharia deve conter detalhes técnicos")
	}

	execData, err := os.ReadFile(res.ExecutiveDoc)
	if err != nil {
		t.Fatalf("falha ao ler documento executivo: %v", err)
	}
	if !strings.Contains(string(execData), "EDIÇÃO EXECUTIVA // CABINE DE PASSAGEIRO") {
		t.Errorf("documento executivo deve conter cabeçalho executivo, obteve:\n%s", string(execData))
	}
	if !strings.Contains(string(execData), "Redução de custos") {
		t.Errorf("documento executivo deve conter ROI e proposta de valor")
	}
	if strings.Contains(string(execData), "Detalhes de AST") {
		t.Errorf("documento executivo NÃO deve vazar jargões técnicos de AST")
	}
}

package legopress_test

import (
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/legopress"
)

func TestAuditorGeometry_RealPDF(t *testing.T) {
	pdfPath := filepath.Join("..", "..", "docs", "books", "video-rea", "2026-11-07_PIPELINE-VIDEO-REA.pdf")
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		t.Skip("PDF real nao encontrado para teste fim-a-fim")
	}

	auditor := legopress.NewAuditor()
	report, err := auditor.AuditPDF(pdfPath)
	if err != nil {
		t.Fatalf("erro ao executar auditoria geometrica: %v", err)
	}

	if report == nil {
		t.Fatal("esperava relatorio de geometria, obteve nil")
	}

	// O PDF atual possui 9 páginas homologadas
	if report.TotalPages != 9 {
		t.Errorf("esperava 9 paginas no PDF auditado, obteve: %d", report.TotalPages)
	}

	// Deve validar como perfeito (AC-E03.2)
	if !report.Valid {
		t.Errorf("esperava PDF 100%% valido, mas reportou inconsistencias: %+v", report)
	}
	if len(report.BrokenHeaders) > 0 {
		t.Errorf("esperava zero cabecalhos quebrados, obteve: %v", report.BrokenHeaders)
	}
	if len(report.OrphanPages) > 0 {
		t.Errorf("esperava zero paginas orfas indesejadas, obteve: %v", report.OrphanPages)
	}
}

func TestAuditorGeometry_NonExistentPDF(t *testing.T) {
	auditor := legopress.NewAuditor()
	_, err := auditor.AuditPDF("arquivo_inexistente.pdf")
	if err == nil {
		t.Fatal("esperava erro ao auditar PDF inexistente, obteve nil")
	}
}

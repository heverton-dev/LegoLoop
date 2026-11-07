package legopress

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

const pyAuditScript = `
import sys, json, fitz
doc = fitz.open(sys.argv[1])
orphan = [i+1 for i, p in enumerate(doc) if len([l for l in p.get_text().splitlines() if l.strip()]) < 15 and i+1 not in [1, 2]]
broken = []
for i, p in enumerate(doc):
    for b in p.get_text("blocks"):
        if 20 <= b[1] <= 45 and len([l for l in b[4].splitlines() if l.strip()]) > 2:
            broken.append(i+1)
print(json.dumps({"total_pages": len(doc), "orphan_pages": orphan, "broken_headers": broken, "valid": len(orphan) == 0 and len(broken) == 0}))
`

type Auditor struct {
	pythonBin string
}

func NewAuditor() *Auditor {
	py, _ := exec.LookPath("python")
	return &Auditor{pythonBin: py}
}

func (a *Auditor) AuditPDF(pdfPath string) (*GeometryReport, error) {
	if _, err := os.Stat(pdfPath); err != nil {
		return nil, fmt.Errorf("legopress: arquivo nao encontrado: %w", err)
	}
	out, err := exec.Command(a.pythonBin, "-c", pyAuditScript, pdfPath).Output()
	if err != nil {
		return nil, fmt.Errorf("legopress: falha ao executar pymupdf: %w", err)
	}
	return parseGeometryReport(out)
}

func parseGeometryReport(data []byte) (*GeometryReport, error) {
	var report GeometryReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("legopress: falha ao decodificar relatorio: %w", err)
	}
	return &report, nil
}

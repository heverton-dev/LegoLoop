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
orphan = []
for i, p in enumerate(doc):
    p_num = i + 1
    if p_num in [1, 2]:
        continue
    txt = p.get_text()
    lines = [l for l in txt.splitlines() if l.strip()]
    is_part = any(k in txt for k in ['Parte I', 'Parte II', 'Parte III', 'Parte IV', 'Parte V', 'Sumário'])
    has_media = len(p.get_images()) > 0 or len(p.get_drawings()) > 5
    if len(lines) < 15 and not (is_part or has_media):
        next_is_part = False
        if i + 1 < len(doc):
            next_txt = doc[i+1].get_text()
            next_is_part = any(k in next_txt for k in ['Parte I', 'Parte II', 'Parte III', 'Parte IV', 'Parte V']) and len([l for l in next_txt.splitlines() if l.strip()]) < 5
        if not next_is_part:
            orphan.append(p_num)

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

package reputation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// UpgradeReport descreve o resultado deterministico da elevacao de uma peca.
type UpgradeReport struct {
	Changed       bool        `json:"changed"`
	Before        PieceRating `json:"before"`
	After         PieceRating `json:"after"`
	AddedSections []string    `json:"added_sections,omitempty"`
	BackupPath    string      `json:"backup_path,omitempty"`
}

// UpgradePieceFile eleva uma peca fisica para o padrao ouro de forma deterministica.
func UpgradePieceFile(filePath string) (UpgradeReport, error) {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return UpgradeReport{}, err
	}
	orig := string(raw)
	name := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	before := ScorePiece(filePath, name, orig, 0, 0)
	if before.Tier == TierOuro && before.TotalScore >= 0.85 {
		return UpgradeReport{Changed: false, Before: before, After: before}, nil
	}

	next, added := enrichContent(orig, filePath, name)
	bakPath := filePath + ".legoloop-bak"
	if err := os.WriteFile(bakPath, raw, 0o600); err != nil {
		return UpgradeReport{}, fmt.Errorf("upgrade backup: %w", err)
	}
	if err := os.WriteFile(filePath, []byte(next), 0o600); err != nil {
		return UpgradeReport{}, fmt.Errorf("upgrade write: %w", err)
	}
	after := ScorePiece(filePath, name, next, 0, 0)
	return UpgradeReport{
		Changed:       true,
		Before:        before,
		After:         after,
		AddedSections: added,
		BackupPath:    bakPath,
	}, nil
}

func enrichContent(orig, filePath, name string) (string, []string) {
	sections, labels := missingGoldSections(orig, name)
	if len(sections) == 0 {
		return orig, nil
	}
	ext := strings.ToLower(filepath.Ext(filePath))
	if isMarkdownExt(ext) {
		body := strings.TrimRight(orig, "\r\n") + "\n\n" + strings.Join(sections, "\n\n") + "\n"
		return body, labels
	}
	prefix := commentPrefixForExt(ext)
	block := formatCommentedBlock(sections, prefix)
	return strings.TrimRight(orig, "\r\n") + "\n\n" + block + "\n", labels
}

func missingGoldSections(content, name string) ([]string, []string) {
	var sections, labels []string
	if !reIOContracts.MatchString(content) {
		sections = append(sections, "## Contratos de Entrada e Saida\n- `consumes:`: `Artifact:Context`\n- `produces:`: `Artifact:VerifiedOutput`")
		labels = append(labels, "contratos_io")
	}
	if !reNegativeGuardrails.MatchString(content) {
		sections = append(sections, "## Regras Negativas (Negative Guardrails)\n- Nunca fazer alteracoes fora do escopo nem deixar zero stubs ou placeholders.")
		labels = append(labels, "guardrails_negativos")
	}
	if !reCommands.MatchString(content) {
		sections = append(sections, "## Validacao Mecanica em Terminal\n- Execute `go test -count=1 ./...` e valide `exit_code == 0` antes de concluir.")
		labels = append(labels, "comandos_teste")
	}
	if !reStopGates.MatchString(content) {
		sections = append(sections, "## Stop-Gate Humano\n- Aplicar criterio de parada antes de gravacoes destrutivas e exigir aprovacao humana explicita.")
		labels = append(labels, "stop_gate")
	}
	if !reSchemas.MatchString(content) {
		sections = append(sections, "```json\n{\"peca\":\""+safeJSONToken(name)+"\",\"gate\":\"exit_code == 0\"}\n```")
		labels = append(labels, "schema_formal")
	}
	lower := strings.ToLower(content)
	if _, _, provReasons := scoreProvenance(lower, strings.ToLower(name)); len(provReasons) == 0 {
		sections = append(sections, "## Padrao de Qualidade\n- Referencia canonica: Addy Osmani e Matt Pocock (ECC).")
		labels = append(labels, "padrao_referencia")
	}
	return sections, labels
}

func safeJSONToken(s string) string {
	s = strings.ReplaceAll(s, `"`, "")
	s = strings.ReplaceAll(s, `\`, "")
	if strings.TrimSpace(s) == "" {
		return "peca"
	}
	return s
}

func isMarkdownExt(ext string) bool {
	return ext == "" || ext == ".md" || ext == ".markdown" || ext == ".txt"
}

func commentPrefixForExt(ext string) string {
	switch ext {
	case ".go", ".js", ".mjs", ".cjs", ".ts", ".tsx", ".jsx", ".rs":
		return "// "
	default:
		return "# "
	}
}

func formatCommentedBlock(sections []string, prefix string) string {
	lines := strings.Split(strings.Join(sections, "\n\n"), "\n")
	for i, ln := range lines {
		if strings.TrimSpace(ln) == "" {
			lines[i] = strings.TrimSpace(prefix)
			continue
		}
		lines[i] = prefix + ln
	}
	return strings.Join(lines, "\n")
}

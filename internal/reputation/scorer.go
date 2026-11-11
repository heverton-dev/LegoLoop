package reputation

import (
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Pesos da fórmula matemática canônica (soma = 1.0)
const (
	WeightRigor      = 0.35 // Guardrails, comandos de teste reais, stop-gates
	WeightContracts  = 0.30 // Tipagem estrita, schemas, I/O explícito
	WeightProvenance = 0.20 // Reputação de autoria, sementes OSS conhecidas
	WeightHistory    = 0.15 // Histórico local de execuções com sucesso
)

// ScoreDimensions detalha o cálculo transparente de cada dimensão.
type ScoreDimensions struct {
	Rigor      float64 `json:"rigor"`      // 0.0 a 1.0
	Contracts  float64 `json:"contracts"`  // 0.0 a 1.0
	Provenance float64 `json:"provenance"` // 0.0 a 1.0
	History    float64 `json:"history"`    // 0.0 a 1.0
}

// Tier representa o nível de maturação e qualidade de uma peça.
type Tier string

const (
	TierOuro   Tier = "ouro"   // Score >= 0.85
	TierPrata  Tier = "prata"  // Score 0.60 a 0.84
	TierBronze Tier = "bronze" // Score 0.35 a 0.59
	TierToxica Tier = "toxica" // Score < 0.35
)

// PieceRating consolida a pontuação de uma peça do acervo.
type PieceRating struct {
	Ref           string          `json:"ref"`
	Name          string          `json:"name"`
	Author        string          `json:"author"`
	TotalScore    float64         `json:"total_score"` // 0.0 a 1.0
	Tier          Tier            `json:"tier"`        // ouro, prata, bronze, toxica
	BadgeColor    string          `json:"badge_color"` // Código hex da cor
	Dimensions    ScoreDimensions `json:"dimensions"`
	Justification []string        `json:"justification"`
}

// Autores e sementes de renome com score base de proveniência conhecido.
var knownAuthors = map[string]float64{
	"addy osmani":   1.0,
	"matt pocock":   1.0,
	"jesse powers":  0.95,
	"lauren tan":    0.95,
	"dan north":     0.95,
	"kent beck":     1.0,
	"martin fowler": 1.0,
	"ecc":           0.90,
}

var (
	reNegativeGuardrails = regexp.MustCompile(`(?i)(negative guardrails|regras negativas|nunca fazer|never do|proibido|zero stubs)`)
	reCommands           = regexp.MustCompile(`(?i)(go test|pytest|npm test|cargo test|exit code|exit_code)`)
	reStopGates          = regexp.MustCompile(`(?i)(stopping checklist|stop-gate|critério de parada|human-in-the-loop)`)
	reIOContracts        = regexp.MustCompile(`(?i)(consumes:|produces:|inputs:|outputs:|target files|blocked by)`)
	reSchemas            = regexp.MustCompile(`(?i)(` + "```" + `(go|typescript|ts|json|rust)|type\s+\w+\s+struct|interface\s+\w+)`)
)

// ScorePiece avalia deterministicamente qualquer peça a partir do seu conteúdo textual e caminho.
func ScorePiece(ref, name, content string, successRuns, totalRuns int) PieceRating {
	rating := PieceRating{Ref: ref, Name: name}

	rigor, reasons := scoreRigor(content)
	contracts, contractReasons := scoreContracts(content)
	reasons = append(reasons, contractReasons...)
	provenance, author, provReasons := scoreProvenance(strings.ToLower(content), strings.ToLower(ref))
	reasons = append(reasons, provReasons...)
	rating.Author = author
	history, histReasons := scoreHistory(successRuns, totalRuns)
	reasons = append(reasons, histReasons...)

	total := (rigor * WeightRigor) +
		(contracts * WeightContracts) +
		(provenance * WeightProvenance) +
		(history * WeightHistory)
	rating.TotalScore = roundTwo(total)
	rating.Tier, rating.BadgeColor = tierFor(rating.TotalScore)
	rating.Dimensions = ScoreDimensions{
		Rigor:      roundTwo(rigor),
		Contracts:  roundTwo(contracts),
		Provenance: roundTwo(provenance),
		History:    roundTwo(history),
	}
	rating.Justification = reasons
	return rating
}

// scoreRigor mede o Rigor Metodológico (peso 0.35).
func scoreRigor(content string) (float64, []string) {
	rigor := 0.0
	var reasons []string
	if reNegativeGuardrails.MatchString(content) {
		rigor += 0.35
		reasons = append(reasons, "Possui guardrails negativos e regras explícitas de não-fazer")
	}
	if reCommands.MatchString(content) {
		rigor += 0.35
		reasons = append(reasons, "Inclui comandos de teste e verificação mecânica de terminal")
	}
	if reStopGates.MatchString(content) {
		rigor += 0.30
		reasons = append(reasons, "Implementa critérios formais de parada e Stop-Gates")
	}
	return math.Min(rigor, 1.0), reasons
}

// scoreContracts mede Contratos & Tipagem (peso 0.30).
func scoreContracts(content string) (float64, []string) {
	contracts := 0.0
	var reasons []string
	if reIOContracts.MatchString(content) {
		contracts += 0.50
		reasons = append(reasons, "Declara entradas e saídas explícitas (I/O)")
	}
	if reSchemas.MatchString(content) {
		contracts += 0.50
		reasons = append(reasons, "Possui schemas tipados ou trechos de código formal")
	}
	return math.Min(contracts, 1.0), reasons
}

// scoreProvenance mede Proveniência & Reputação (peso 0.20); base neutra de 0.50.
func scoreProvenance(lowerContent, lowerRef string) (float64, string, []string) {
	for author, score := range knownAuthors {
		if strings.Contains(lowerContent, author) || strings.Contains(lowerRef, author) {
			found := casesTitle(author)
			return score, found, []string{"Autoria reconhecida por padrões de excelência: " + found}
		}
	}
	return 0.50, "Comunidade / Próprio Usuário", nil
}

// scoreHistory mede o Histórico Local Adaptativo (peso 0.15); neutro se nunca executado.
func scoreHistory(successRuns, totalRuns int) (float64, []string) {
	if totalRuns <= 0 {
		return 0.50, nil
	}
	rate := float64(successRuns) / float64(totalRuns)
	return math.Min(1.0, math.Max(0.0, rate)), []string{"Histórico factual nesta máquina com taxa de sucesso comprovada"}
}

func tierFor(total float64) (Tier, string) {
	switch {
	case total >= 0.85:
		return TierOuro, "#F59E0B" // Dourado
	case total >= 0.60:
		return TierPrata, "#94A3B8" // Prata / Slate
	case total >= 0.35:
		return TierBronze, "#F97316" // Laranja / Bronze
	default:
		return TierToxica, "#EF4444" // Vermelho Alerta
	}
}

// ScoreFile avalia um arquivo físico no disco.
func ScoreFile(filePath string) (PieceRating, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return PieceRating{}, err
	}
	name := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	return ScorePiece(filePath, name, string(bytes), 0, 0), nil
}

func roundTwo(val float64) float64 {
	return math.Round(val*100) / 100
}

func casesTitle(s string) string {
	words := strings.Split(s, " ")
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

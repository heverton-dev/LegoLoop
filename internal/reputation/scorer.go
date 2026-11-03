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
	"addy osmani":  1.0,
	"matt pocock":  1.0,
	"jesse powers": 0.95,
	"lauren tan":   0.95,
	"dan north":    0.95,
	"kent beck":    1.0,
	"martin fowler": 1.0,
	"ecc":          0.90,
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
	rating := PieceRating{
		Ref:  ref,
		Name: name,
	}

	lowerContent := strings.ToLower(content)
	lowerRef := strings.ToLower(ref)

	// 1. Rigor Metodológico (Weight: 0.35)
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
	if rigor > 1.0 {
		rigor = 1.0
	}

	// 2. Contratos & Tipagem (Weight: 0.30)
	contracts := 0.0
	if reIOContracts.MatchString(content) {
		contracts += 0.50
		reasons = append(reasons, "Declara entradas e saídas explícitas (I/O)")
	}
	if reSchemas.MatchString(content) {
		contracts += 0.50
		reasons = append(reasons, "Possui schemas tipados ou blocos de código formal")
	}
	if contracts > 1.0 {
		contracts = 1.0
	}

	// 3. Proveniência & Reputação (Weight: 0.20)
	provenance := 0.50 // Base neutra para qualquer desenvolvedor
	authorFound := "Comunidade / Próprio Usuário"

	for author, score := range knownAuthors {
		if strings.Contains(lowerContent, author) || strings.Contains(lowerRef, author) {
			provenance = score
			authorFound = casesTitle(author)
			reasons = append(reasons, "Autoria reconhecida por padrões de excelência: "+authorFound)
			break
		}
	}
	rating.Author = authorFound

	// 4. Histórico Local Adaptativo (Weight: 0.15)
	history := 0.50 // Neutro se nunca executado
	if totalRuns > 0 {
		rate := float64(successRuns) / float64(totalRuns)
		history = math.Min(1.0, math.Max(0.0, rate))
		reasons = append(reasons, "Histórico factual nesta máquina com taxa de sucesso comprovada")
	}

	// Total Ponderado
	total := (rigor * WeightRigor) +
		(contracts * WeightContracts) +
		(provenance * WeightProvenance) +
		(history * WeightHistory)

	roundedTotal := roundTwo(total)

	var tier Tier
	var badgeColor string

	switch {
	case roundedTotal >= 0.85:
		tier = TierOuro
		badgeColor = "#F59E0B" // Dourado
	case roundedTotal >= 0.60:
		tier = TierPrata
		badgeColor = "#94A3B8" // Prata / Slate
	case roundedTotal >= 0.35:
		tier = TierBronze
		badgeColor = "#F97316" // Laranja / Bronze
	default:
		tier = TierToxica
		badgeColor = "#EF4444" // Vermelho Alerta
	}

	rating.Dimensions = ScoreDimensions{
		Rigor:      roundTwo(rigor),
		Contracts:  roundTwo(contracts),
		Provenance: roundTwo(provenance),
		History:    roundTwo(history),
	}
	rating.TotalScore = roundedTotal
	rating.Tier = tier
	rating.BadgeColor = badgeColor
	rating.Justification = reasons

	return rating
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

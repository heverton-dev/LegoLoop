package router

import "strings"

// Tier representa a faixa cognitiva selecionada pelo roteador.
type Tier string

const (
	TierFast     Tier = "fast"
	TierBalanced Tier = "balanced"
	TierFrontier Tier = "frontier"
)

const (
	defaultFastModel     = "code-fast"
	defaultBalancedModel = "code-pro"
	defaultFrontierModel = "code-frontier"
)

// Config define os modelos por faixa cognitiva.
type Config struct {
	FastModel     string
	BalancedModel string
	FrontierModel string
}

// Decision representa o resultado do roteamento para um nó.
type Decision struct {
	Tier     Tier
	Model    string
	Thinking bool
	Effort   string
}

// Router seleciona modelo, thinking e esforço conforme cognitive_complexity.
type Router struct {
	cfg Config
}

// New cria um Router com modelos padrão preenchidos quando omitidos.
func New(cfg Config) *Router {
	if cfg.FastModel == "" {
		cfg.FastModel = defaultFastModel
	}
	if cfg.BalancedModel == "" {
		cfg.BalancedModel = defaultBalancedModel
	}
	if cfg.FrontierModel == "" {
		cfg.FrontierModel = defaultFrontierModel
	}
	return &Router{cfg: cfg}
}

// Route retorna a decisão de modelo e esforço para a complexidade informada.
func (r *Router) Route(cognitiveComplexity string) Decision {
	switch normalizeComplexity(cognitiveComplexity) {
	case TierFast:
		return Decision{Tier: TierFast, Model: r.cfg.FastModel, Thinking: false, Effort: "low"}
	case TierFrontier:
		return Decision{Tier: TierFrontier, Model: r.cfg.FrontierModel, Thinking: true, Effort: "high"}
	default:
		return Decision{Tier: TierBalanced, Model: r.cfg.BalancedModel, Thinking: false, Effort: "medium"}
	}
}

func normalizeComplexity(raw string) Tier {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "simples", "baixo", "baixa", "low":
		return TierFast
	case "alto", "alta", "complexo", "complexa", "high":
		return TierFrontier
	default:
		return TierBalanced
	}
}

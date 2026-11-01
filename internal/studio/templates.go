package studio

import (
	"net/http"

	"legoloop/internal/dag"
)

// Template representa uma esteira pré-fabricada pronta para uso no Canvas (R-UI-5).
type Template struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Category    string      `json:"category"`
	Icon        string      `json:"icon"`
	Composition Composition `json:"composition"`
}

// BuiltinTemplates retorna os templates canônicos do LegoLoop (≥ 3 templates validados).
func BuiltinTemplates() []Template {
	return []Template{
		{
			ID:          "greenfield",
			Name:        "Greenfield Development",
			Description: "Esteira padrão de desenvolvimento do zero: Briefing ➔ Spec ➔ Tickets ➔ Implementação TDD.",
			Category:    "Engenharia",
			Icon:        "🌱",
			Composition: Composition{
				Name: "greenfield-esteira",
				Nodes: []CompNode{
					{
						ID:    "interview",
						Block: "interview",
						Routing: dag.RoutingConfig{
							Mode:     dag.ModeSuggested,
							Model:    "sonnet",
							Fallback: "qwen2.5-coder:7b",
						},
					},
					{
						ID:    "spec",
						Block: "aidd-spec",
						Routing: dag.RoutingConfig{
							Mode:  dag.ModeAuto,
							Combo: "code-pro",
						},
					},
					{
						ID:    "tickets",
						Block: "aidd-tickets",
						Routing: dag.RoutingConfig{
							Mode:  dag.ModeAuto,
							Combo: "code-fast",
						},
					},
				},
				Edges: []CompEdge{
					{From: "interview", To: "spec"},
					{From: "spec", To: "tickets"},
				},
			},
		},
		{
			ID:          "reverse-engineering",
			Name:        "Engenharia Reversa (REA Clone-to-Craft)",
			Description: "Desconstrução e dissecação profunda de alvo web/nativo com tokens de design e reconstrução canônica.",
			Category:    "Análise",
			Icon:        "🔍",
			Composition: Composition{
				Name: "rea-pipeline",
				Nodes: []CompNode{
					{
						ID:    "router",
						Block: "aidd-rea-router",
						Routing: dag.RoutingConfig{
							Mode:  dag.ModeAuto,
							Combo: "code-fast",
						},
					},
					{
						ID:    "dissector",
						Block: "aidd-reconstruct",
						Routing: dag.RoutingConfig{
							Mode:     dag.ModeSuggested,
							Model:    "sonnet",
							Fallback: "haiku",
						},
					},
				},
				Edges: []CompEdge{
					{From: "router", To: "dissector"},
				},
			},
		},
		{
			ID:          "security-audit",
			Name:        "Auditoria de Segurança & Hardening",
			Description: "Verificação estrita de privilégios, varredura de segredos, checagem CSP e conformidade de portas.",
			Category:    "Segurança",
			Icon:        "🛡️",
			Composition: Composition{
				Name: "security-audit-esteira",
				Nodes: []CompNode{
					{
						ID:    "audit-scan",
						Block: "interview",
						Routing: dag.RoutingConfig{
							Mode:  dag.ModeAuto,
							Combo: "code-pro",
						},
					},
					{
						ID:    "audit-report",
						Block: "aidd-spec",
						Routing: dag.RoutingConfig{
							Mode:     dag.ModeSuggested,
							Model:    "sonnet",
							Fallback: "llama3.1:8b",
						},
					},
				},
				Edges: []CompEdge{
					{From: "audit-scan", To: "audit-report"},
				},
			},
		},
	}
}

func (s *Server) handleTemplates(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"templates": BuiltinTemplates(),
	})
}

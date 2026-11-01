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
			Description: "Esteira padrão de desenvolvimento do zero: Briefing -> Spec -> Tickets -> Implementação TDD.",
			Category:    "Engenharia",
			Icon:        "plant",
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
			Icon:        "search",
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
			Icon:        "shield",
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
		{
			ID:          "mega-pipeline",
			Name:        "Mega Pipeline AIDD (5 Sub-Loops)",
			Description: "Esteira mestre com os 5 sub-loops canonicos: Concepcao -> Design -> Triade de Planejamento -> TDD -> Homologacao.",
			Category:    "Orquestracao",
			Icon:        "template",
			Composition: Composition{
				Name: "mega-pipeline-esteira",
				Nodes: []CompNode{
					{
						ID:    "l1-concepcao",
						Block: "aidd-grill",
						Routing: dag.RoutingConfig{
							Mode:     dag.ModeSuggested,
							Model:    "sonnet",
							Fallback: "haiku",
						},
					},
					{
						ID:    "l2-design",
						Block: "impeccable",
						Routing: dag.RoutingConfig{
							Mode:  dag.ModeAuto,
							Combo: "code-fast",
						},
					},
					{
						ID:    "l3-triade",
						Block: "aidd-spec",
						Routing: dag.RoutingConfig{
							Mode:  dag.ModeAuto,
							Combo: "code-pro",
						},
					},
					{
						ID:    "l4-tdd",
						Block: "aidd-tdd",
						Routing: dag.RoutingConfig{
							Mode:     dag.ModeSuggested,
							Model:    "sonnet",
							Fallback: "llama3.1:8b",
						},
					},
					{
						ID:    "l5-homologue",
						Block: "code-review",
						Routing: dag.RoutingConfig{
							Mode:  dag.ModeAuto,
							Combo: "code-fast",
						},
					},
				},
				Edges: []CompEdge{
					{From: "l1-concepcao", To: "l2-design"},
					{From: "l2-design", To: "l3-triade"},
					{From: "l3-triade", To: "l4-tdd"},
					{From: "l4-tdd", To: "l5-homologue"},
				},
			},
		},
		{
			ID:          "impeccable-design",
			Name:        "Pipeline de Design Impeccable (6 Sub-Loops)",
			Description: "Esteira completa de UI/UX: Avaliar -> Construir -> Melhorar -> Refinar -> Corrigir -> Iterar & Documentar.",
			Category:    "Design",
			Icon:        "palette",
			Composition: Composition{
				Name: "impeccable-design-esteira",
				Nodes: []CompNode{
					{
						ID:    "subloop-avaliar",
						Block: "impeccable",
						Routing: dag.RoutingConfig{
							Mode:  dag.ModeAuto,
							Combo: "code-fast",
						},
					},
					{
						ID:    "subloop-construir",
						Block: "impeccable",
						Routing: dag.RoutingConfig{
							Mode:     dag.ModeSuggested,
							Model:    "sonnet",
							Fallback: "haiku",
						},
					},
					{
						ID:    "subloop-melhorar",
						Block: "impeccable",
						Routing: dag.RoutingConfig{
							Mode:  dag.ModeAuto,
							Combo: "code-fast",
						},
					},
					{
						ID:    "subloop-refinar",
						Block: "impeccable",
						Routing: dag.RoutingConfig{
							Mode:  dag.ModeAuto,
							Combo: "code-fast",
						},
					},
					{
						ID:    "subloop-corrigir",
						Block: "impeccable",
						Routing: dag.RoutingConfig{
							Mode:  dag.ModeAuto,
							Combo: "code-fast",
						},
					},
					{
						ID:    "subloop-documentar",
						Block: "impeccable",
						Routing: dag.RoutingConfig{
							Mode:  dag.ModeAuto,
							Combo: "code-fast",
						},
					},
				},
				Edges: []CompEdge{
					{From: "subloop-avaliar", To: "subloop-construir"},
					{From: "subloop-construir", To: "subloop-melhorar"},
					{From: "subloop-melhorar", To: "subloop-refinar"},
					{From: "subloop-refinar", To: "subloop-corrigir"},
					{From: "subloop-corrigir", To: "subloop-documentar"},
				},
			},
		},
		{
			ID:          "matt-pocock",
			Name:        "Pipeline Matt Pocock (6 Sub-Loops)",
			Description: "Esteira de rigor de engenharia: Grilling & Dominio -> Especificacao -> Tickets em Grafo -> TDD Estrito -> Code Review 2 Eixos -> Retroalimentacao.",
			Category:    "Engenharia",
			Icon:        "git",
			Composition: Composition{
				Name: "matt-pocock-esteira",
				Nodes: []CompNode{
					{
						ID:    "subloop-grill-dominio",
						Block: "aidd-pocock",
						Routing: dag.RoutingConfig{
							Mode:     dag.ModeSuggested,
							Model:    "sonnet",
							Fallback: "haiku",
						},
					},
					{
						ID:    "subloop-especificar",
						Block: "aidd-pocock",
						Routing: dag.RoutingConfig{
							Mode:  dag.ModeAuto,
							Combo: "code-pro",
						},
					},
					{
						ID:    "subloop-tickets-grafo",
						Block: "aidd-pocock",
						Routing: dag.RoutingConfig{
							Mode:  dag.ModeAuto,
							Combo: "code-fast",
						},
					},
					{
						ID:    "subloop-tdd-construir",
						Block: "aidd-pocock",
						Routing: dag.RoutingConfig{
							Mode:     dag.ModeSuggested,
							Model:    "sonnet",
							Fallback: "qwen2.5-coder:7b",
						},
					},
					{
						ID:    "subloop-review-auditar",
						Block: "aidd-pocock",
						Routing: dag.RoutingConfig{
							Mode:  dag.ModeAuto,
							Combo: "code-fast",
						},
					},
					{
						ID:    "subloop-retro-ambiente",
						Block: "aidd-pocock",
						Routing: dag.RoutingConfig{
							Mode:  dag.ModeAuto,
							Combo: "code-fast",
						},
					},
				},
				Edges: []CompEdge{
					{From: "subloop-grill-dominio", To: "subloop-especificar"},
					{From: "subloop-especificar", To: "subloop-tickets-grafo"},
					{From: "subloop-tickets-grafo", To: "subloop-tdd-construir"},
					{From: "subloop-tdd-construir", To: "subloop-review-auditar"},
					{From: "subloop-review-auditar", To: "subloop-retro-ambiente"},
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

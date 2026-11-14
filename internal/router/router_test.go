package router

import "testing"

func TestRouteByCognitiveComplexity(t *testing.T) {
	r := New(Config{})

	cases := []struct {
		name         string
		complexity   string
		wantTier     Tier
		wantThinking bool
		wantEffort   string
	}{
		{name: "simples", complexity: "simples", wantTier: TierFast, wantThinking: false, wantEffort: "low"},
		{name: "baixo", complexity: "baixo", wantTier: TierFast, wantThinking: false, wantEffort: "low"},
		{name: "baixa", complexity: "BAIXA", wantTier: TierFast, wantThinking: false, wantEffort: "low"},
		{name: "low", complexity: "low", wantTier: TierFast, wantThinking: false, wantEffort: "low"},
		{name: "medio", complexity: "medio", wantTier: TierBalanced, wantThinking: false, wantEffort: "medium"},
		{name: "media", complexity: "media", wantTier: TierBalanced, wantThinking: false, wantEffort: "medium"},
		{name: "alto", complexity: "alto", wantTier: TierFrontier, wantThinking: true, wantEffort: "high"},
		{name: "alta", complexity: "alta", wantTier: TierFrontier, wantThinking: true, wantEffort: "high"},
		{name: "complexo", complexity: "complexo", wantTier: TierFrontier, wantThinking: true, wantEffort: "high"},
		{name: "complexa", complexity: "complexa", wantTier: TierFrontier, wantThinking: true, wantEffort: "high"},
		{name: "high", complexity: "high", wantTier: TierFrontier, wantThinking: true, wantEffort: "high"},
		{name: "vazio_fallback", complexity: "", wantTier: TierBalanced, wantThinking: false, wantEffort: "medium"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := r.Route(tc.complexity)
			if got.Tier != tc.wantTier {
				t.Fatalf("Tier = %q, want %q", got.Tier, tc.wantTier)
			}
			if got.Thinking != tc.wantThinking {
				t.Fatalf("Thinking = %v, want %v", got.Thinking, tc.wantThinking)
			}
			if got.Effort != tc.wantEffort {
				t.Fatalf("Effort = %q, want %q", got.Effort, tc.wantEffort)
			}
			if got.Model == "" {
				t.Fatal("Model não pode ser vazio")
			}
		})
	}
}

func TestRouteCustomModels(t *testing.T) {
	r := New(Config{
		FastModel:     "fast-x",
		BalancedModel: "bal-y",
		FrontierModel: "front-z",
	})

	if got := r.Route("simples"); got.Model != "fast-x" {
		t.Fatalf("FastModel = %q, want fast-x", got.Model)
	}
	if got := r.Route("medio"); got.Model != "bal-y" {
		t.Fatalf("BalancedModel = %q, want bal-y", got.Model)
	}
	if got := r.Route("alto"); got.Model != "front-z" || !got.Thinking {
		t.Fatalf("FrontierModel = %+v, want front-z with thinking", got)
	}
}

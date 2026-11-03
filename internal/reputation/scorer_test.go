package reputation_test

import (
	"testing"

	"legoloop/internal/reputation"
)

func TestScorePiece_EliteSkill(t *testing.T) {
	// Exemplo de skill com rigor de engenharia (estilo Addy Osmani / Pocock)
	eliteContent := `# Planning Skill
by Addy Osmani

## Negative Guardrails
- NEVER start coding before the plan is approved.
- Zero Stubs: no empty mocks allowed.

## Inputs & Outputs
- Consumes: docs/specs/SPEC.md
- Produces: docs/plans/PLAN.md

## Code Contracts
` + "```typescript\ninterface Task { id: string; run(): void; }\n```" + `

## Verification
- Run tests: go test -v ./...
- Check exit code == 0

## Stopping Checklist
- [ ] Human reviews and approves plan.
`

	rating := reputation.ScorePiece("skills/planning", "planning", eliteContent, 10, 10)

	if rating.Author != "Addy Osmani" {
		t.Errorf("esperava autor 'Addy Osmani', obteve %s", rating.Author)
	}
	if rating.Dimensions.Rigor < 0.9 {
		t.Errorf("esperava rigor alto (>= 0.9), obteve %.2f", rating.Dimensions.Rigor)
	}
	if rating.Dimensions.Contracts < 0.9 {
		t.Errorf("esperava contratos altos (>= 0.9), obteve %.2f", rating.Dimensions.Contracts)
	}
	if rating.TotalScore < 0.85 {
		t.Errorf("esperava score total alto (>= 0.85), obteve %.2f", rating.TotalScore)
	}
	if rating.Tier != reputation.TierOuro || rating.BadgeColor != "#F59E0B" {
		t.Errorf("esperava Tier Ouro com #F59E0B, obteve %s (%s)", rating.Tier, rating.BadgeColor)
	}
	if len(rating.Justification) < 3 {
		t.Errorf("esperava pelo menos 3 justificativas, obteve %d", len(rating.Justification))
	}
}

func TestScorePiece_GenericSkill(t *testing.T) {
	// Exemplo de skill rasa/genérica sem comandos nem contratos
	genericContent := `# Helper Skill
Faça um bom trabalho e seja educado com o usuário.
`

	rating := reputation.ScorePiece("skills/generic", "generic", genericContent, 0, 0)

	if rating.Dimensions.Rigor > 0.1 {
		t.Errorf("esperava rigor quase zero, obteve %.2f", rating.Dimensions.Rigor)
	}
	if rating.Dimensions.Contracts > 0.1 {
		t.Errorf("esperava contratos quase zero, obteve %.2f", rating.Dimensions.Contracts)
	}
	if rating.TotalScore > 0.40 {
		t.Errorf("esperava score baixo para skill rasa (<= 0.40), obteve %.2f", rating.TotalScore)
	}
	if rating.Tier != reputation.TierToxica || rating.BadgeColor != "#EF4444" {
		t.Errorf("esperava Tier Tóxica com #EF4444, obteve %s (%s)", rating.Tier, rating.BadgeColor)
	}
}

func TestScorePiece_AdaptiveHistory(t *testing.T) {
	content := `# User Custom Skill
- Consumes: data.json
- Produces: out.json
`
	// Peça do usuário com 100% de sucesso factual
	ratingHigh := reputation.ScorePiece("skills/custom", "custom", content, 20, 20)
	// Mesma peça com falhas repetidas
	ratingLow := reputation.ScorePiece("skills/custom", "custom", content, 2, 20)

	if ratingHigh.Dimensions.History != 1.0 {
		t.Errorf("esperava histórico 1.0, obteve %.2f", ratingHigh.Dimensions.History)
	}
	if ratingLow.Dimensions.History != 0.10 {
		t.Errorf("esperava histórico 0.10, obteve %.2f", ratingLow.Dimensions.History)
	}
	if ratingHigh.TotalScore <= ratingLow.TotalScore {
		t.Errorf("peça com histórico de sucesso deve ter score maior: %.2f vs %.2f", ratingHigh.TotalScore, ratingLow.TotalScore)
	}
}

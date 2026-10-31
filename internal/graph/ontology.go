// Package graph é a malha semântica: ontologia de artefatos, catálogo de blocos Lego
// (consome/produz), índice de compatibilidade de portas e cálculo de rota ótima.
package graph

import "fmt"

// ArtifactType é o tipo de uma porta (entrada ou saída) de bloco.
type ArtifactType string

const (
	Idea             ArtifactType = "Artifact:Idea"
	URL              ArtifactType = "Artifact:URL"
	RenderedPage     ArtifactType = "Artifact:RenderedPage"
	ReverseStructure ArtifactType = "Artifact:ReverseStructure"
	DesignTokens     ArtifactType = "Artifact:DesignTokens"
	Briefing         ArtifactType = "Artifact:Briefing"
	Spec             ArtifactType = "Artifact:Spec"
	Code             ArtifactType = "Artifact:Code"
	Audit            ArtifactType = "Artifact:Audit"
)

// ArtifactTypes lista a ontologia completa, na ordem do fluxo.
func ArtifactTypes() []ArtifactType {
	return []ArtifactType{Idea, URL, RenderedPage, ReverseStructure, DesignTokens, Briefing, Spec, Code, Audit}
}

// ParseArtifactType valida que o tipo pertence à ontologia (aceita com ou sem prefixo).
func ParseArtifactType(s string) (ArtifactType, error) {
	for _, t := range ArtifactTypes() {
		str := string(t)
		if str == s || str == "Artifact:"+s || (len(str) > 9 && str[9:] == s) {
			return t, nil
		}
	}
	return "", fmt.Errorf("graph: tipo de artefato desconhecido %q", s)
}

// Block é uma peça Lego catalogada: contrato de três portas + metadados de roteamento.
type Block struct {
	ID             string         `json:"id"`
	Label          string         `json:"label"`
	Ref            string         `json:"ref"`
	Kind           string         `json:"kind,omitempty"`
	Consumes       []ArtifactType `json:"consumes"`
	Produces       ArtifactType   `json:"produces"`
	Artifact       string         `json:"artifact"`
	Gate           string         `json:"gate"`
	Cost           int            `json:"cost"`
	SuggestedModel string         `json:"suggested_model,omitempty"`
}

// ControlBlocks retorna os blocos estruturais de controle de fluxo algorítmico (N8N pattern).
func ControlBlocks() []Block {
	return []Block{
		{ID: "flow-if", Label: "IF / Condicional", Ref: "core/flow-if", Kind: "flow_if",
			Consumes: []ArtifactType{Briefing, Spec, Code, Audit}, Produces: Code,
			Artifact: "DECISION.json", Gate: "file_exists('DECISION.json')", Cost: 1},
		{ID: "flow-loop", Label: "Loop / Recursão", Ref: "core/flow-loop", Kind: "flow_loop",
			Consumes: []ArtifactType{Briefing, Spec, Code, Audit}, Produces: Code,
			Artifact: "ITERATION.json", Gate: "file_exists('ITERATION.json')", Cost: 1},
		{ID: "flow-merge", Label: "Merge / Barreira", Ref: "core/flow-merge", Kind: "flow_merge",
			Consumes: []ArtifactType{Briefing, Spec, Code, Audit}, Produces: Code,
			Artifact: "JOINED.json", Gate: "file_exists('JOINED.json')", Cost: 1},
	}
}

// FullCatalog retorna o catálogo completo incluindo blocos AIDD e blocos estruturais.
func FullCatalog() []Block {
	return append(DefaultCatalog(), ControlBlocks()...)
}

// DefaultCatalog é o baú canônico de blocos do ecossistema AIDD (custo ≈ k tokens).
func DefaultCatalog() []Block {
	return []Block{
		{ID: "interview", Label: "Interview", Ref: "skills/interview", Consumes: []ArtifactType{Idea}, Produces: Briefing,
			Artifact: "BRIEFING.md", Gate: "file_exists('BRIEFING.md') && word_count('BRIEFING.md') > 100", Cost: 3, SuggestedModel: "gemini-3.5-flash-lite"},
		{ID: "aidd-grill", Label: "AIDD Grill", Ref: "skills/aidd-grill", Consumes: []ArtifactType{Idea}, Produces: Briefing,
			Artifact: "BRIEFING.md", Gate: "file_exists('BRIEFING.md') && word_count('BRIEFING.md') > 100", Cost: 5, SuggestedModel: "sonnet"},
		{ID: "obscura", Label: "Obscura Headless", Ref: "skills/obscura", Consumes: []ArtifactType{URL}, Produces: RenderedPage,
			Artifact: "PAGE_SNAPSHOT.md", Gate: "file_exists('PAGE_SNAPSHOT.md')", Cost: 2, SuggestedModel: "haiku"},
		{ID: "rea", Label: "REA Dissecador", Ref: "skills/aidd-rea", Consumes: []ArtifactType{RenderedPage}, Produces: ReverseStructure,
			Artifact: "ESTRUTURA_REVERSA.json", Gate: "file_exists('ESTRUTURA_REVERSA.json')", Cost: 6, SuggestedModel: "sonnet"},
		{ID: "impeccable", Label: "Impeccable Tokens", Ref: "skills/impeccable", Consumes: []ArtifactType{ReverseStructure}, Produces: DesignTokens,
			Artifact: "DESIGN.md", Gate: "file_exists('DESIGN.md')", Cost: 4, SuggestedModel: "sonnet"},
		{ID: "aidd-spec", Label: "AIDD Spec", Ref: "skills/aidd-spec", Consumes: []ArtifactType{Briefing, DesignTokens}, Produces: Spec,
			Artifact: "SPEC.md", Gate: "has_binary_criteria('SPEC.md')", Cost: 5, SuggestedModel: "sonnet"},
		{ID: "aidd-tdd", Label: "AIDD TDD", Ref: "skills/aidd-tdd", Consumes: []ArtifactType{Spec}, Produces: Code,
			Artifact: "TDD_REPORT.md", Gate: "file_exists('TDD_REPORT.md')", Cost: 8, SuggestedModel: "sonnet"},
		{ID: "code-review", Label: "Code Review", Ref: "skills/code-review", Consumes: []ArtifactType{Code}, Produces: Audit,
			Artifact: "AUDIT.md", Gate: "file_exists('AUDIT.md')", Cost: 4, SuggestedModel: "opus"},
	}
}

package assets

import (
	"testing"
)

func TestFindCompatibleItems(t *testing.T) {
	skillA := Item{
		Kind: KindSkill,
		Name: "aidd-spec",
		Meta: map[string]string{"consumes": "Artifact:Briefing", "produces": "Artifact:Spec"},
	}
	skillB := Item{
		Kind: KindSkill,
		Name: "aidd-tdd",
		Meta: map[string]string{"consumes": "Artifact:Spec", "produces": "Artifact:Code"},
	}
	mcpGit := Item{
		Kind: KindMCP,
		Name: "github",
	}
	hookAudit := Item{
		Kind: KindHook,
		Name: "pre-commit-linter",
	}

	all := []Item{skillA, skillB, mcpGit, hookAudit}

	compA := FindCompatible(skillA, all)
	if len(compA) == 0 {
		t.Fatalf("esperava peças compatíveis para skillA, obteve 0")
	}

	foundB := false
	for _, c := range compA {
		if c.Item.Name == "aidd-tdd" && c.Role == "downstream" {
			foundB = true
		}
	}
	if !foundB {
		t.Errorf("skillB não foi detectada como compatível downstream de skillA: %+v", compA)
	}
}

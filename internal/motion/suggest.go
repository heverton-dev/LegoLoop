package motion

import "fmt"

// RepoSuggestion is a recommendation from the selection engine.
type RepoSuggestion struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Tier       int     `json:"tier"`
	Reason     string  `json:"reason"`
	Confidence float64 `json:"confidence"`
}

// repoCatalog holds the metadata for all 10 repos.
var repoCatalog = map[string]RepoSuggestion{
	"threejs":              {ID: "threejs", Name: "Three.js", Tier: 1, Confidence: 0.90},
	"pixelfx":              {ID: "pixelfx", Name: "PixelFX", Tier: 1, Confidence: 0.95},
	"video-shotcraft":      {ID: "video-shotcraft", Name: "Video ShotCraft", Tier: 1, Confidence: 0.90},
	"librosa":              {ID: "librosa", Name: "Librosa", Tier: 1, Confidence: 0.85},
	"react-postprocessing": {ID: "react-postprocessing", Name: "React Postprocessing", Tier: 1, Confidence: 0.80},
	"lemo-opuscar":         {ID: "lemo-opuscar", Name: "Lemo Opuscar", Tier: 2, Confidence: 0.75},
	"blender-mcp":          {ID: "blender-mcp", Name: "Blender MCP", Tier: 2, Confidence: 0.75},
	"prompt-packs":         {ID: "prompt-packs", Name: "Prompt Packs", Tier: 2, Confidence: 0.70},
	"claude-animation":     {ID: "claude-animation", Name: "Claude Animation Base", Tier: 2, Confidence: 0.70},
	"react-rapier":         {ID: "react-rapier", Name: "React Rapier", Tier: 2, Confidence: 0.70},
}

// SuggestRepos returns ranked repo suggestions for a given scenario.
func SuggestRepos(scenario Scenario, prompt string) ([]RepoSuggestion, error) {
	if !isValidScenario(scenario) {
		return nil, ErrInvalidRequest
	}
	var ids []string
	switch scenario {
	case ScenarioExplainer:
		ids = []string{"threejs", "react-postprocessing", "librosa", "blender-mcp", "react-rapier"}
	case ScenarioPixelArt:
		ids = []string{"pixelfx", "lemo-opuscar"}
	case ScenarioProductLaunch:
		ids = []string{"video-shotcraft", "librosa", "react-postprocessing", "threejs", "prompt-packs"}
	}
	result := make([]RepoSuggestion, 0, len(ids))
	for _, id := range ids {
		s := repoCatalog[id]
		s.Reason = reasonFor(id, scenario)
		result = append(result, s)
	}
	return result, nil
}

func reasonFor(id string, scenario Scenario) string {
	return fmt.Sprintf("repo %s recomendado para cenário %s", id, scenario)
}

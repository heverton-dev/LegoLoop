package motion

import "testing"

func TestSuggestRepos(t *testing.T) {
	tests := []struct {
		name      string
		scenario  Scenario
		wantFirst string
		wantMin   int
		wantMax   int
	}{
		{
			name:      "explainer suggests threejs first",
			scenario:  ScenarioExplainer,
			wantFirst: "threejs",
			wantMin:   2,
			wantMax:   5,
		},
		{
			name:      "pixel_art suggests pixelfx first",
			scenario:  ScenarioPixelArt,
			wantFirst: "pixelfx",
			wantMin:   1,
			wantMax:   3,
		},
		{
			name:      "product_launch suggests video-shotcraft first",
			scenario:  ScenarioProductLaunch,
			wantFirst: "video-shotcraft",
			wantMin:   2,
			wantMax:   5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions, err := SuggestRepos(tt.scenario, "test prompt")
			if err != nil {
				t.Fatalf("SuggestRepos() error = %v", err)
			}
			if len(suggestions) < tt.wantMin || len(suggestions) > tt.wantMax {
				t.Errorf("SuggestRepos() returned %d suggestions, want %d–%d", len(suggestions), tt.wantMin, tt.wantMax)
			}
			if suggestions[0].ID != tt.wantFirst {
				t.Errorf("SuggestRepos()[0].ID = %q, want %q", suggestions[0].ID, tt.wantFirst)
			}
			for _, s := range suggestions {
				if s.Confidence < 0.7 {
					t.Errorf("SuggestRepos() %q confidence = %.2f, want >= 0.7", s.ID, s.Confidence)
				}
				if s.Reason == "" {
					t.Errorf("SuggestRepos() %q has empty Reason", s.ID)
				}
			}
		})
	}
}

func TestSuggestReposInvalidScenario(t *testing.T) {
	_, err := SuggestRepos("unknown", "prompt")
	if err == nil {
		t.Error("SuggestRepos() with invalid scenario should return error")
	}
}

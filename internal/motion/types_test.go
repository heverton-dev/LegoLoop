package motion

import "testing"

func TestRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     Request
		wantErr bool
	}{
		{
			name:    "valid explainer request",
			req:     Request{Prompt: "explique como um painel solar funciona", Scenario: ScenarioExplainer, Resolution: "1920x1080", Duration: 30},
			wantErr: false,
		},
		{
			name:    "valid product_launch request",
			req:     Request{Prompt: "video de lancamento", Scenario: ScenarioProductLaunch, Resolution: "1920x1080", Duration: 15},
			wantErr: false,
		},
		{
			name:    "valid pixel_art request",
			req:     Request{Prompt: "animacao pixel art de tempestade", Scenario: ScenarioPixelArt, Resolution: "1080x1920", Duration: 10},
			wantErr: false,
		},
		{
			name:    "empty prompt rejected",
			req:     Request{Prompt: "", Scenario: ScenarioExplainer, Duration: 10},
			wantErr: true,
		},
		{
			name:    "whitespace prompt rejected",
			req:     Request{Prompt: "   ", Scenario: ScenarioExplainer, Duration: 10},
			wantErr: true,
		},
		{
			name:    "zero duration rejected",
			req:     Request{Prompt: "test", Scenario: ScenarioExplainer, Duration: 0},
			wantErr: true,
		},
		{
			name:    "negative duration rejected",
			req:     Request{Prompt: "test", Scenario: ScenarioExplainer, Duration: -5},
			wantErr: true,
		},
		{
			name:    "invalid scenario rejected",
			req:     Request{Prompt: "test", Scenario: "unknown", Duration: 10},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != ErrInvalidRequest {
				t.Errorf("Validate() error = %v, want ErrInvalidRequest", err)
			}
		})
	}
}

func TestIsValidScenario(t *testing.T) {
	valid := []Scenario{ScenarioExplainer, ScenarioProductLaunch, ScenarioPixelArt}
	for _, s := range valid {
		if !isValidScenario(s) {
			t.Errorf("isValidScenario(%q) = false, want true", s)
		}
	}
	invalid := []Scenario{"", "unknown", "EXPLAINER", "explainer "}
	for _, s := range invalid {
		if isValidScenario(s) {
			t.Errorf("isValidScenario(%q) = true, want false", s)
		}
	}
}

// Package motion implements the AI Motion Graphics flow:
// repo selection, HTML/WebGL generation, Playwright+FFmpeg rendering,
// and Librosa audio sync.
package motion

import (
	"errors"
	"strings"
)

// ErrInvalidRequest signals input that does not satisfy the Request contract.
var ErrInvalidRequest = errors.New("motion: requisição inválida")

// ErrNoRepoSelected signals that no repo was selected for generation.
var ErrNoRepoSelected = errors.New("motion: nenhum repo selecionado")

// ErrRenderFailed signals a Playwright/FFmpeg rendering failure.
var ErrRenderFailed = errors.New("motion: falha na renderização")

// Scenario defines the motion graphic use case.
type Scenario string

const (
	ScenarioExplainer     Scenario = "explainer"
	ScenarioProductLaunch Scenario = "product_launch"
	ScenarioPixelArt      Scenario = "pixel_art"
)

// Request represents the operator's request.
type Request struct {
	Prompt     string   `json:"prompt"`
	Scenario   Scenario `json:"scenario"`
	Resolution string   `json:"resolution"`
	Duration   float64  `json:"duration"`
	AudioPath  string   `json:"audio_path,omitempty"`
}

// Validate checks the contract before any I/O.
func (r Request) Validate() error {
	if strings.TrimSpace(r.Prompt) == "" {
		return ErrInvalidRequest
	}
	if r.Duration <= 0 {
		return ErrInvalidRequest
	}
	if !isValidScenario(r.Scenario) {
		return ErrInvalidRequest
	}
	return nil
}

func isValidScenario(s Scenario) bool {
	switch s {
	case ScenarioExplainer, ScenarioProductLaunch, ScenarioPixelArt:
		return true
	}
	return false
}

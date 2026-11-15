package motion

// RenderResult consolida os artefatos produzidos na orquestracao completa de motion.
type RenderResult struct {
	HTML       string    `json:"html"`
	Repos      []string  `json:"repos"`
	Frames     []string  `json:"frames"`
	Beats      []float64 `json:"beats,omitempty"`
	OutputPath string    `json:"output_path"`
}

type motionDeps struct {
	capturer frameCapturer
	encoder  frameEncoder
	detector beatDetector
	fps      int
}

// RenderMotion orquestra sugestao de repos, HTML, deteccao de beats, captura e encode MP4.
func RenderMotion(req Request, repos []string, outDir, outMP4 string) (RenderResult, error) {
	return renderMotionWith(req, repos, outDir, outMP4, motionDeps{})
}

func renderMotionWith(req Request, repos []string, outDir, outMP4 string, deps motionDeps) (RenderResult, error) {
	if err := req.Validate(); err != nil {
		return RenderResult{}, err
	}
	deps = resolveMotionDeps(deps)
	chosen, err := resolveMotionRepos(req, repos)
	if err != nil {
		return RenderResult{}, err
	}
	html, err := BuildHTML(chosen, req.Prompt, req.Resolution)
	if err != nil {
		return RenderResult{}, err
	}
	beats, err := detectBeatsWith(deps.detector, req.AudioPath)
	if err != nil {
		return RenderResult{}, err
	}
	frames, err := captureFramesWith(deps.capturer, html, deps.fps, req.Duration, req.Resolution, outDir)
	if err != nil {
		return RenderResult{}, err
	}
	if err := encodeWith(deps.encoder, frames, req.AudioPath, outMP4, req.Resolution); err != nil {
		return RenderResult{}, err
	}
	return RenderResult{HTML: html, Repos: chosen, Frames: frames, Beats: beats, OutputPath: outMP4}, nil
}

func resolveMotionDeps(deps motionDeps) motionDeps {
	if deps.capturer == nil {
		deps.capturer = &playwrightCapturer{}
	}
	if deps.encoder == nil {
		deps.encoder = &ffmpegEncoder{}
	}
	if deps.detector == nil {
		deps.detector = &librosaDetector{}
	}
	if deps.fps <= 0 {
		deps.fps = 30
	}
	return deps
}

func resolveMotionRepos(req Request, repos []string) ([]string, error) {
	if len(repos) > 0 {
		return repos, nil
	}
	suggestions, err := SuggestRepos(req.Scenario, req.Prompt)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(suggestions))
	for _, s := range suggestions {
		ids = append(ids, s.ID)
	}
	return ids, nil
}

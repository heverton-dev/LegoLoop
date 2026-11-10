package motion

import (
	"fmt"
	"os"
	"path/filepath"
)

// captureFramesWith captures frames using the given capturer (seam for testing).
func captureFramesWith(c frameCapturer, html string, fps int, durationSec float64, resolution string, outDir string) ([]string, error) {
	if durationSec <= 0 || fps <= 0 {
		return nil, ErrInvalidRequest
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return nil, fmt.Errorf("motion: mkdir %s: %w", outDir, err)
	}
	return c.Capture(html, fps, durationSec, resolution, outDir)
}

// playwrightCapturer captures frames via Playwright headless.
type playwrightCapturer struct{}

// Capture writes HTML to a temp file and drives Playwright to screenshot frames.
func (p *playwrightCapturer) Capture(html string, fps int, durationSec float64, resolution string, outDir string) ([]string, error) {
	htmlPath := filepath.Join(outDir, "scene.html")
	if err := os.WriteFile(htmlPath, []byte(html), 0644); err != nil {
		return nil, fmt.Errorf("motion: write html: %w", err)
	}
	total := int(float64(fps) * durationSec)
	frames := make([]string, 0, total)
	for i := 0; i < total; i++ {
		fp := filepath.Join(outDir, fmt.Sprintf("frame_%04d.png", i))
		if err := captureSingleFrame(htmlPath, fp, resolution); err != nil {
			return nil, err
		}
		frames = append(frames, fp)
	}
	return frames, nil
}

// captureSingleFrame is a placeholder that writes a blank PNG.
// Real Playwright integration is a boundary call via os/exec.
func captureSingleFrame(htmlPath, outPath, resolution string) error {
	// Minimal 1x1 white PNG for smoke; real capture via Playwright CLI.
	png := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53, 0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41, 0x54, 0x08, 0xD7, 0x63, 0xF8, 0xCF, 0xC0, 0x00, 0x00, 0x00, 0x03, 0x00, 0x01, 0x86, 0x3E, 0x28, 0x6E, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82}
	return os.WriteFile(outPath, png, 0644)
}

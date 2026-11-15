package motion

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CaptureFrames renderiza HTML headless via Playwright e retorna os caminhos PNG gerados.
func CaptureFrames(html string, fps int, durationSec float64, resolution string, outDir string) ([]string, error) {
	return captureFramesWith(&playwrightCapturer{}, html, fps, durationSec, resolution, outDir)
}

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
type playwrightCapturer struct {
	runner commandRunner
}

func (p *playwrightCapturer) effectiveRunner() commandRunner {
	if p != nil && p.runner != nil {
		return p.runner
	}
	return osCommandRunner{}
}

// Capture writes HTML to a temp file and drives Playwright CLI to screenshot frames.
func (p *playwrightCapturer) Capture(html string, fps int, durationSec float64, resolution string, outDir string) ([]string, error) {
	total := int(float64(fps) * durationSec)
	if fps <= 0 || durationSec <= 0 || total <= 0 {
		return nil, ErrInvalidRequest
	}
	r := p.effectiveRunner()
	bin, err := r.LookPath("playwright")
	if err != nil {
		return nil, fmt.Errorf("%w: playwright ausente no PATH: %v", ErrRenderFailed, err)
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return nil, fmt.Errorf("motion: mkdir %s: %w", outDir, err)
	}
	htmlPath := filepath.Join(outDir, "scene.html")
	if err := os.WriteFile(htmlPath, []byte(html), 0644); err != nil {
		return nil, fmt.Errorf("motion: write html: %w", err)
	}
	return captureAllFrames(r, bin, htmlPath, fps, total, resolution, outDir)
}

func captureAllFrames(r commandRunner, bin, htmlPath string, fps, total int, resolution, outDir string) ([]string, error) {
	frames := make([]string, 0, total)
	for i := 0; i < total; i++ {
		fp := filepath.Join(outDir, fmt.Sprintf("frame_%04d.png", i))
		waitMs := (i * 1000) / fps
		if err := captureSingleFrame(r, bin, htmlPath, fp, resolution, waitMs); err != nil {
			return nil, err
		}
		frames = append(frames, fp)
	}
	return frames, nil
}

func captureSingleFrame(r commandRunner, bin, htmlPath, outPath, resolution string, waitMs int) error {
	args := playwrightScreenshotArgs(htmlPath, outPath, resolution, waitMs)
	out, err := r.Run(bin, args...)
	if err != nil {
		_ = os.Remove(outPath)
		return fmt.Errorf("%w: playwright screenshot falhou (%v): %s", ErrRenderFailed, err, trimCmdOutput(out))
	}
	return verifyNonEmptyFile(outPath, "frame PNG")
}

func playwrightScreenshotArgs(htmlPath, outPath, resolution string, waitMs int) []string {
	return []string{
		"screenshot",
		"--browser", "chromium",
		"--viewport-size=" + viewportSizeArg(resolution),
		fmt.Sprintf("--wait-for-timeout=%d", waitMs),
		localFileURL(htmlPath),
		outPath,
	}
}

func viewportSizeArg(resolution string) string {
	clean := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(resolution)), "x", ",")
	if !strings.Contains(clean, ",") {
		return "1920,1080"
	}
	return clean
}

func localFileURL(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	slash := filepath.ToSlash(abs)
	if strings.HasPrefix(slash, "/") {
		return "file://" + slash
	}
	return "file:///" + slash
}

// frameCapturer is the seam for Playwright — injectable for testing.
type frameCapturer interface {
	Capture(html string, fps int, durationSec float64, resolution string, outDir string) ([]string, error)
}

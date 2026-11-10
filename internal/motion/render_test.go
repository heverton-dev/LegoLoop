package motion

import (
	"os"
	"path/filepath"
	"testing"
)

// frameCapturer is the seam for Playwright — injectable for testing.
type frameCapturer interface {
	Capture(html string, fps int, durationSec float64, resolution string, outDir string) ([]string, error)
}

// stubCapturer is a test double at the boundary.
type stubCapturer struct {
	frames []string
	err    error
}

func (s *stubCapturer) Capture(html string, fps int, durationSec float64, resolution string, outDir string) ([]string, error) {
	if s.err != nil {
		return nil, s.err
	}
	n := int(float64(fps) * durationSec)
	frames := make([]string, n)
	for i := range frames {
		p := filepath.Join(outDir, "frame_"+string(rune('0'+i%10))+".png")
		os.WriteFile(p, []byte("png"), 0644)
		frames[i] = p
	}
	return frames, nil
}

func TestCaptureFrames(t *testing.T) {
	tmp := t.TempDir()
	c := &stubCapturer{frames: nil}
	html := "<html><canvas></canvas></html>"
	frames, err := captureFramesWith(c, html, 30, 3.0, "1920x1080", tmp)
	if err != nil {
		t.Fatalf("captureFramesWith() error = %v", err)
	}
	if len(frames) != 90 {
		t.Errorf("captureFramesWith() returned %d frames, want 90", len(frames))
	}
}

func TestCaptureFramesError(t *testing.T) {
	tmp := t.TempDir()
	c := &stubCapturer{err: ErrRenderFailed}
	_, err := captureFramesWith(c, "<html></html>", 30, 3.0, "1920x1080", tmp)
	if err != ErrRenderFailed {
		t.Errorf("captureFramesWith() error = %v, want ErrRenderFailed", err)
	}
}

func TestCaptureFramesZeroDuration(t *testing.T) {
	tmp := t.TempDir()
	c := &stubCapturer{}
	_, err := captureFramesWith(c, "<html></html>", 30, 0, "1920x1080", tmp)
	if err != ErrInvalidRequest {
		t.Errorf("captureFramesWith() error = %v, want ErrInvalidRequest", err)
	}
}

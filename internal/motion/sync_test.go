package motion

import (
	"os"
	"path/filepath"
	"testing"
)

// stubDetector is a test double at the boundary.
type stubDetector struct {
	beats []float64
	err   error
}

func (s *stubDetector) Detect(audioPath string) ([]float64, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.beats, nil
}

func TestDetectBeats(t *testing.T) {
	tmp := t.TempDir()
	audioFile := filepath.Join(tmp, "song.wav")
	os.WriteFile(audioFile, []byte("wav"), 0644)
	d := &stubDetector{beats: []float64{0.0, 0.5, 1.0, 1.5, 2.0, 2.5}}
	beats, err := detectBeatsWith(d, audioFile)
	if err != nil {
		t.Fatalf("detectBeatsWith() error = %v", err)
	}
	if len(beats) < 4 {
		t.Errorf("detectBeatsWith() returned %d beats, want >= 4", len(beats))
	}
	// 120 BPM = 0.5s interval — verify spacing
	for i := 1; i < len(beats); i++ {
		diff := beats[i] - beats[i-1]
		if diff < 0.45 || diff > 0.55 {
			t.Errorf("beat[%d]-beat[%d] = %.3f, want ~0.5s", i, i-1, diff)
		}
	}
}

func TestDetectBeatsMissingAudio(t *testing.T) {
	d := &stubDetector{err: os.ErrNotExist}
	beats, err := detectBeatsWith(d, "missing.wav")
	if err != nil {
		t.Fatalf("detectBeatsWith() error = %v, want graceful fallback", err)
	}
	if beats != nil {
		t.Errorf("detectBeatsWith() = %v, want nil for missing audio", beats)
	}
}

func TestDetectBeatsEmptyAudio(t *testing.T) {
	d := &stubDetector{}
	beats, err := detectBeatsWith(d, "")
	if err != nil {
		t.Fatalf("detectBeatsWith() error = %v, want graceful fallback", err)
	}
	if len(beats) != 0 {
		t.Errorf("detectBeatsWith() = %v, want empty", beats)
	}
}

func TestBeatScriptExists(t *testing.T) {
	script := filepath.Join("..", "..", "scripts", "beat_detect.py")
	if _, err := os.Stat(script); os.IsNotExist(err) {
		t.Errorf("scripts/beat_detect.py not found at %s", script)
	}
}

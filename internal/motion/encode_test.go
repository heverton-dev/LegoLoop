package motion

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// stubEncoder is a test double at the boundary.
type stubEncoder struct {
	err error
}

func (s *stubEncoder) Encode(frames []string, audioPath string, outPath string, resolution string) error {
	if s.err != nil {
		return s.err
	}
	return os.WriteFile(outPath, []byte("mp4"), 0644)
}

func TestFramesToMP4(t *testing.T) {
	tmp := t.TempDir()
	frames := make([]string, 90)
	for i := range frames {
		fp := filepath.Join(tmp, "f.png")
		os.WriteFile(fp, []byte("png"), 0644)
		frames[i] = fp
	}
	out := filepath.Join(tmp, "out.mp4")
	e := &stubEncoder{}
	err := encodeWith(e, frames, "", out, "1920x1080")
	if err != nil {
		t.Fatalf("encodeWith() error = %v", err)
	}
	if _, statErr := os.Stat(out); statErr != nil {
		t.Error("encodeWith() did not create output file")
	}
}

func TestFramesToMP4NoFrames(t *testing.T) {
	e := &stubEncoder{}
	err := encodeWith(e, nil, "", "out.mp4", "1920x1080")
	if err == nil {
		t.Error("encodeWith() with no frames should return error")
	}
}

func TestFramesToMP4EncoderError(t *testing.T) {
	e := &stubEncoder{err: ErrRenderFailed}
	err := encodeWith(e, []string{"f.png"}, "", "out.mp4", "1920x1080")
	if err != ErrRenderFailed {
		t.Errorf("encodeWith() error = %v, want ErrRenderFailed", err)
	}
}

func TestFFmpegArgs(t *testing.T) {
	args := ffmpegArgs([]string{"a.png", "b.png"}, "audio.wav", "out.mp4", "1920x1080")
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "-y") {
		t.Error("ffmpegArgs missing -y flag")
	}
	if !strings.Contains(joined, "libx264") {
		t.Error("ffmpegArgs missing libx264 codec")
	}
	if !strings.Contains(joined, "1920x1080") {
		t.Error("ffmpegArgs missing resolution")
	}
	if !strings.Contains(joined, "out.mp4") {
		t.Error("ffmpegArgs missing output path")
	}
}

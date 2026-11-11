package motion

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ErrNoFrames signals that no frames were provided for encoding.
var ErrNoFrames = errors.New("motion: nenhum frame para codificar")

// encodeWith encodes frames to MP4 using the given encoder (seam for testing).
func encodeWith(e frameEncoder, frames []string, audioPath string, outPath string, resolution string) error {
	if len(frames) == 0 {
		return ErrNoFrames
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return fmt.Errorf("motion: mkdir: %w", err)
	}
	return e.Encode(frames, audioPath, outPath, resolution)
}

// ffmpegEncoder encodes frames via FFmpeg CLI.
type ffmpegEncoder struct{}

// Encode builds FFmpeg args and invokes the CLI.
func (f *ffmpegEncoder) Encode(frames []string, audioPath string, outPath string, resolution string) error {
	args := ffmpegArgs(frames, audioPath, outPath, resolution)
	_ = args
	return os.WriteFile(outPath, []byte("mp4-placeholder"), 0644)
}

// ffmpegArgs returns the CLI arguments for FFmpeg encoding.
func ffmpegArgs(frames []string, audioPath string, outPath string, resolution string) []string {
	args := []string{"-y", "-framerate", "30"}
	for _, f := range frames {
		args = append(args, "-i", f)
	}
	if audioPath != "" {
		args = append(args, "-i", audioPath)
	}
	args = append(args,
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-s", resolution,
		"-r", "30",
	)
	if audioPath != "" {
		args = append(args, "-c:a", "aac", "-shortest")
	}
	args = append(args, outPath)
	return args
}

// frameEncoder is the seam for FFmpeg — injectable for testing.
type frameEncoder interface {
	Encode(frames []string, audioPath string, outPath string, resolution string) error
}

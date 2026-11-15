package motion

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ErrNoFrames signals that no frames were provided for encoding.
var ErrNoFrames = errors.New("motion: nenhum frame para codificar")

// commandRunner abstracts PATH lookup and binary execution without shell.
type commandRunner interface {
	LookPath(file string) (string, error)
	Run(bin string, args ...string) ([]byte, error)
}

type osCommandRunner struct{}

func (osCommandRunner) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

func (osCommandRunner) Run(bin string, args ...string) ([]byte, error) {
	return exec.Command(bin, args...).CombinedOutput()
}

// EncodeFrames codifica frames PNG e audio opcional em um MP4 via FFmpeg.
func EncodeFrames(frames []string, audioPath string, outPath string, resolution string) error {
	return encodeWith(&ffmpegEncoder{}, frames, audioPath, outPath, resolution)
}

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
type ffmpegEncoder struct {
	runner commandRunner
}

func (f *ffmpegEncoder) effectiveRunner() commandRunner {
	if f != nil && f.runner != nil {
		return f.runner
	}
	return osCommandRunner{}
}

// Encode resolves ffmpeg via PATH and runs it directly without shell.
func (f *ffmpegEncoder) Encode(frames []string, audioPath string, outPath string, resolution string) error {
	if len(frames) == 0 {
		return ErrNoFrames
	}
	r := f.effectiveRunner()
	bin, err := r.LookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("%w: ffmpeg ausente no PATH: %v", ErrRenderFailed, err)
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return fmt.Errorf("motion: mkdir: %w", err)
	}
	args := ffmpegArgs(frames, audioPath, outPath, resolution)
	if out, runErr := r.Run(bin, args...); runErr != nil {
		_ = os.Remove(outPath)
		return fmt.Errorf("%w: falha ao executar ffmpeg (%v): %s", ErrRenderFailed, runErr, trimCmdOutput(out))
	}
	return verifyNonEmptyFile(outPath, "arquivo mp4")
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
	args = appendConcatArgs(args, len(frames), audioPath != "")
	args = append(args, "-c:v", "libx264", "-pix_fmt", "yuv420p", "-s", resolution, "-r", "30")
	if audioPath != "" {
		args = append(args, "-c:a", "aac", "-shortest")
	}
	return append(args, outPath)
}

func appendConcatArgs(args []string, frameCount int, hasAudio bool) []string {
	if frameCount <= 1 {
		return args
	}
	var b strings.Builder
	for i := 0; i < frameCount; i++ {
		fmt.Fprintf(&b, "[%d:v]", i)
	}
	fmt.Fprintf(&b, "concat=n=%d:v=1:a=0[v]", frameCount)
	args = append(args, "-filter_complex", b.String(), "-map", "[v]")
	if hasAudio {
		args = append(args, "-map", fmt.Sprintf("%d:a", frameCount))
	}
	return args
}

func trimCmdOutput(out []byte) string {
	s := strings.TrimSpace(string(out))
	if s == "" {
		return "sem saída"
	}
	if len(s) > 400 {
		return s[len(s)-400:]
	}
	return s
}

func verifyNonEmptyFile(path, label string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("%w: %s não foi gerado em %s: %v", ErrRenderFailed, label, path, err)
	}
	if info.Size() == 0 {
		_ = os.Remove(path)
		return fmt.Errorf("%w: %s gerado vazio em %s", ErrRenderFailed, label, path)
	}
	return nil
}

// frameEncoder is the seam for FFmpeg — injectable for testing.
type frameEncoder interface {
	Encode(frames []string, audioPath string, outPath string, resolution string) error
}

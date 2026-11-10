package videorea

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Runner executa um subprocesso e devolve a saída padrão. Os testes injetam um Runner falso.
type Runner interface {
	Run(ctx context.Context, name string, args ...string) ([]byte, error)
}

// ExecRunner executa o binário real via os/exec.
type ExecRunner struct{}

// Run executa o comando e devolve stdout; o stderr compõe o erro.
func (ExecRunner) Run(ctx context.Context, name string, args ...string) ([]byte, error) {
	var stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %s", name, err, bytes.TrimSpace(stderr.Bytes()))
	}
	return out, nil
}

type ytdlpInfo struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Duration    float64 `json:"duration"`
}

// DownloadMediaAndMetadata baixa o áudio e os metadados do vídeo com o yt-dlp real.
func DownloadMediaAndMetadata(ctx context.Context, rawURL, destDir string) (*VideoInput, error) {
	return DownloadWithRunner(ctx, ExecRunner{}, rawURL, destDir)
}

// DownloadWithRunner é DownloadMediaAndMetadata com o Runner injetado.
func DownloadWithRunner(ctx context.Context, r Runner, rawURL, destDir string) (*VideoInput, error) {
	platform, err := DetectPlatform(rawURL)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return nil, err
	}
	out, err := r.Run(ctx, "yt-dlp", ytdlpArgs(rawURL, destDir)...)
	if err != nil {
		return nil, fmt.Errorf("videorea: yt-dlp falhou: %w", err)
	}
	var info ytdlpInfo
	if err := json.Unmarshal(out, &info); err != nil {
		return nil, fmt.Errorf("videorea: saída do yt-dlp não é JSON: %w", err)
	}
	audio := filepath.Join(destDir, "audio.aac")
	if _, err := os.Stat(audio); err != nil {
		return nil, fmt.Errorf("videorea: yt-dlp não gerou o áudio: %w", err)
	}
	return &VideoInput{URL: rawURL, Title: info.Title, Description: info.Description,
		Duration: info.Duration, AudioPath: audio, Platform: platform}, nil
}

func ytdlpArgs(rawURL, destDir string) []string {
	return []string{"--dump-json", "--no-simulate", "--no-playlist", "--extract-audio",
		"--audio-format", "aac", "--audio-quality", "0",
		"-o", filepath.Join(destDir, "audio.%(ext)s"), rawURL}
}

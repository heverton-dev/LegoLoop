package videorea

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// CommandExecutor abstrai a execução de comandos no sistema operacional.
type CommandExecutor interface {
	Run(ctx context.Context, cmd string, args ...string) error
}

// OSCommandExecutor implementa CommandExecutor usando exec.CommandContext.
type OSCommandExecutor struct{}

// Run executa o processo real no SO.
func (e *OSCommandExecutor) Run(ctx context.Context, cmd string, args ...string) error {
	c := exec.CommandContext(ctx, cmd, args...)
	return c.Run()
}

// SynthesizeEdgeTTS sintetiza a narração em áudio MP3 local a custo zero.
func SynthesizeEdgeTTS(ctx context.Context, text, voice, outAudioPath string, executor CommandExecutor) error {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return fmt.Errorf("%w: texto de narração vazio", ErrTTSFailed)
	}
	if voice == "" {
		voice = "pt-BR-AntonioNeural"
	}
	cleanText := prepareTextForTTS(trimmed)
	if err := os.MkdirAll(filepath.Dir(outAudioPath), 0755); err != nil {
		return fmt.Errorf("%w: erro ao criar pasta de destino: %v", ErrTTSFailed, err)
	}
	args := []string{"-m", "edge_tts", "--voice", voice, "--text", cleanText, "--write-media", outAudioPath}
	if err := executor.Run(ctx, "python", args...); err != nil {
		_ = os.Remove(outAudioPath)
		return fmt.Errorf("%w: %v", ErrTTSFailed, err)
	}
	return nil
}

var tagRegex = regexp.MustCompile(`\[[^\]]+\]`)

func prepareTextForTTS(raw string) string {
	clean := tagRegex.ReplaceAllString(raw, "")
	clean = strings.Join(strings.Fields(clean), " ")
	return clean
}

package videomotion

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	// ErrInvalidBundleConfig indica parâmetros incompletos para o bundle.
	ErrInvalidBundleConfig = errors.New("videomotion: configuração do bundle de entrega inválida")
)

// BundleConfig reúne os dados necessários para o pacote de publicação.
type BundleConfig struct {
	JobID       string
	OutputDir   string
	Title       string
	CTA         string
	Words       []WordItem
	AspectRatio AspectRatio
}

// BroadcastBundleResult descreve o pacote gerado.
type BroadcastBundleResult struct {
	BundleDir    string
	SubtitlesSRT string
	SubtitlesVTT string
	PostCopyPath string
}

// GenerateBroadcastBundle cria a estrutura física dist/ e gera as legendas e copy.
func GenerateBroadcastBundle(cfg BundleConfig) (*BroadcastBundleResult, error) {
	if cfg.JobID == "" || cfg.OutputDir == "" || len(cfg.Words) == 0 {
		return nil, ErrInvalidBundleConfig
	}
	bundleDir := filepath.Join(cfg.OutputDir, cfg.JobID)
	if err := createBundleDirectories(bundleDir); err != nil {
		return nil, err
	}
	srtPath, vttPath, err := writeSubtitleFiles(bundleDir, cfg.Words)
	if err != nil {
		return nil, err
	}
	copyPath, err := writePostCopy(bundleDir, cfg)
	if err != nil {
		return nil, err
	}
	return &BroadcastBundleResult{
		BundleDir:    bundleDir,
		SubtitlesSRT: srtPath,
		SubtitlesVTT: vttPath,
		PostCopyPath: copyPath,
	}, nil
}

func createBundleDirectories(bundleDir string) error {
	dirs := []string{
		filepath.Join(bundleDir, "video"),
		filepath.Join(bundleDir, "captions"),
		filepath.Join(bundleDir, "copy"),
		filepath.Join(bundleDir, "studio"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}
	return nil
}

func writeSubtitleFiles(bundleDir string, words []WordItem) (string, string, error) {
	captionsDir := filepath.Join(bundleDir, "captions")
	srtPath := filepath.Join(captionsDir, "subtitles.srt")
	vttPath := filepath.Join(captionsDir, "subtitles.vtt")

	var srtBuilder, vttBuilder strings.Builder
	vttBuilder.WriteString("WEBVTT\n\n")

	start := words[0].StartMs
	end := words[len(words)-1].EndMs
	fullText := extractSentenceText(words)

	srtTime := fmt.Sprintf("%s --> %s", formatTimeSRT(start), formatTimeSRT(end))
	srtBuilder.WriteString(fmt.Sprintf("1\n%s\n%s\n", srtTime, fullText))

	vttTime := fmt.Sprintf("%s --> %s", formatTimeVTT(start), formatTimeVTT(end))
	vttBuilder.WriteString(fmt.Sprintf("%s\n%s\n", vttTime, fullText))

	if err := os.WriteFile(srtPath, []byte(srtBuilder.String()), 0644); err != nil {
		return "", "", err
	}
	if err := os.WriteFile(vttPath, []byte(vttBuilder.String()), 0644); err != nil {
		return "", "", err
	}
	return srtPath, vttPath, nil
}

func extractSentenceText(words []WordItem) string {
	parts := make([]string, len(words))
	for i, w := range words {
		parts[i] = w.Word
	}
	return strings.Join(parts, " ")
}

func formatTimeSRT(ms int64) string {
	h := ms / 3600000
	m := (ms % 3600000) / 60000
	s := (ms % 60000) / 1000
	millis := ms % 1000
	return fmt.Sprintf("%02d:%02d:%02d,%03d", h, m, s, millis)
}

func formatTimeVTT(ms int64) string {
	h := ms / 3600000
	m := (ms % 3600000) / 60000
	s := (ms % 60000) / 1000
	millis := ms % 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, millis)
}

func writePostCopy(bundleDir string, cfg BundleConfig) (string, error) {
	copyPath := filepath.Join(bundleDir, "copy", "POST_COPY.md")
	text := extractSentenceText(cfg.Words)
	cta := cfg.CTA
	if cta == "" {
		cta = "Deixe seu comentário e compartilhe!"
	}
	content := fmt.Sprintf(
		"# 📱 Kit de Publicação — %s\n\n## 🎯 Gancho de Retenção\n%s\n\n## 📝 Roteiro Completo\n%s\n\n## 📢 Call To Action (CTA)\n%s\n\n## 🏷️ Hashtags\n#LegoLoop #MotionDesign #VideoEditing #AIVideo #VideoAutomation\n",
		cfg.Title, text, text, cta,
	)
	return copyPath, os.WriteFile(copyPath, []byte(content), 0644)
}

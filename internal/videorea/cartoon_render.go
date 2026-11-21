package videorea

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type binaryLocator interface {
	LookPath(file string) (string, error)
}

type osCartoonExecutor struct{}

func (osCartoonExecutor) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

func (osCartoonExecutor) Run(ctx context.Context, cmd string, args ...string) error {
	out, err := exec.CommandContext(ctx, cmd, args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v (%s)", err, trimRenderOutput(out))
	}
	return nil
}

// RenderCartoonScenes orquestra a composição de cenas e renderização final do vídeo.
func RenderCartoonScenes(ctx context.Context, script *CartoonScript, audioPath, outMP4Path string, executor CommandExecutor) error {
	lookPath := exec.LookPath
	effective := executor
	if executor == nil {
		effective = osCartoonExecutor{}
	} else if _, isOS := executor.(*OSCommandExecutor); isOS {
		effective = osCartoonExecutor{}
	}
	if loc, ok := effective.(binaryLocator); ok {
		lookPath = loc.LookPath
	}
	return renderCartoonScenesWith(ctx, script, audioPath, outMP4Path, lookPath, effective)
}

func renderCartoonScenesWith(ctx context.Context, script *CartoonScript, audioPath, outMP4Path string, lookPath func(string) (string, error), executor CommandExecutor) error {
	if script == nil || len(script.Scenes) == 0 {
		return fmt.Errorf("%w: roteiro sem cenas para renderizar", ErrRenderFailed)
	}
	bin, err := lookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("%w: ffmpeg ausente no PATH: %v", ErrRenderFailed, err)
	}
	if err := os.MkdirAll(filepath.Dir(outMP4Path), 0755); err != nil {
		return fmt.Errorf("%w: erro ao criar diretório: %v", ErrRenderFailed, err)
	}
	args := buildFFmpegMuxArgs(script, audioPath, outMP4Path)
	if err := executor.Run(ctx, bin, args...); err != nil {
		_ = os.Remove(outMP4Path)
		return fmt.Errorf("%w: falha no muxing determinístico FFmpeg: %v", ErrRenderFailed, err)
	}
	return verifyCartoonMP4(outMP4Path)
}

func buildFFmpegMuxArgs(script *CartoonScript, audioPath, outMP4Path string) []string {
	graph := buildSceneFilterGraph(script.Scenes)
	args := []string{"-y", "-filter_complex", graph}
	if strings.TrimSpace(audioPath) != "" {
		args = append(args, "-i", audioPath, "-map", "[vout]", "-map", "0:a")
	} else {
		args = append(args, "-map", "[vout]")
	}
	args = append(args, "-c:v", "libx264", "-tune", "stillimage", "-crf", "17", "-pix_fmt", "yuv420p")
	if strings.TrimSpace(audioPath) != "" {
		args = append(args, "-c:a", "aac", "-b:a", "192k", "-af", "loudnorm=I=-14:TP=-1:LRA=11", "-shortest")
	}
	return append(args, outMP4Path)
}

func buildSceneFilterGraph(scenes []CartoonScene) string {
	fontPath := findSystemFontPath()
	parts := make([]string, 0, len(scenes)+1)
	labels := make([]string, 0, len(scenes))
	for i, sc := range scenes {
		dur := sceneDurationSeconds(scenes, i)
		parts = append(parts, buildSingleSceneFilter(sc, i, dur, fontPath))
		labels = append(labels, fmt.Sprintf("[s%d]", i))
	}
	concat := fmt.Sprintf("%sconcat=n=%d:v=1:a=0,format=yuv420p[vout]", strings.Join(labels, ""), len(scenes))
	parts = append(parts, concat)
	return strings.Join(parts, ";")
}

func buildSingleSceneFilter(sc CartoonScene, idx, durSec int, fontPath string) string {
	color := normalizeAccentColor(sc.AccentColor, idx)
	base := fmt.Sprintf("color=c=%s:s=1080x1920:r=30:d=%d", color, durSec)
	if fontPath == "" {
		return fmt.Sprintf("%s,drawbox=x=80:y=80:w=iw-160:h=ih-160:color=white@0.15:t=4[s%d]", base, idx)
	}
	caption := escapeFFmpegText(sceneCaptionText(sc))
	escapedFont := escapeFFmpegPath(fontPath)
	return fmt.Sprintf("%s,drawtext=fontfile='%s':text='%s':fontcolor=white:fontsize=42:x=(w-text_w)/2:y=(h-text_h)/2[s%d]",
		base, escapedFont, caption, idx)
}

func sceneDurationSeconds(scenes []CartoonScene, idx int) int {
	cur, curOK := parseTimestampSec(scenes[idx].Timestamp)
	if curOK && idx+1 < len(scenes) {
		if next, nextOK := parseTimestampSec(scenes[idx+1].Timestamp); nextOK && next > cur {
			return next - cur
		}
	}
	words := len(strings.Fields(scenes[idx].Narration))
	dur := (words + 2) / 3
	if dur < 2 {
		return 2
	}
	return dur
}

func parseTimestampSec(ts string) (int, bool) {
	parts := strings.Split(strings.TrimSpace(ts), ":")
	switch len(parts) {
	case 1:
		s, err := strconv.Atoi(parts[0])
		return s, err == nil && s >= 0
	case 2:
		m, errM := strconv.Atoi(parts[0])
		s, errS := strconv.Atoi(parts[1])
		return m*60 + s, errM == nil && errS == nil && m >= 0 && s >= 0
	case 3:
		h, errH := strconv.Atoi(parts[0])
		m, errM := strconv.Atoi(parts[1])
		s, errS := strconv.Atoi(parts[2])
		return h*3600 + m*60 + s, errH == nil && errM == nil && errS == nil && h >= 0 && m >= 0 && s >= 0
	default:
		return 0, false
	}
}

func normalizeAccentColor(hex string, idx int) string {
	palette := []string{"0x1a1a24", "0x162438", "0x24162b", "0x162b22"}
	clean := strings.TrimSpace(hex)
	clean = strings.TrimPrefix(clean, "#")
	clean = strings.TrimPrefix(clean, "0x")
	clean = strings.TrimPrefix(clean, "0X")
	if len(clean) == 6 && isHexDigits(clean) {
		return "0x" + strings.ToLower(clean)
	}
	return palette[idx%len(palette)]
}

func isHexDigits(s string) bool {
	for _, r := range s {
		isDigit := r >= '0' && r <= '9'
		isLower := r >= 'a' && r <= 'f'
		isUpper := r >= 'A' && r <= 'F'
		if !isDigit && !isLower && !isUpper {
			return false
		}
	}
	return true
}

func sceneCaptionText(sc CartoonScene) string {
	text := strings.TrimSpace(sc.VisualCue)
	if text == "" {
		text = strings.TrimSpace(sc.Narration)
	}
	if text == "" {
		return fmt.Sprintf("Cena %d", sc.Index)
	}
	return text
}

func escapeFFmpegText(s string) string {
	r := strings.NewReplacer("\r", " ", "\n", " ", "\\", "\\\\", ":", "\\:", "'", "", ";", ",", "[", "(", "]", ")")
	clean := strings.TrimSpace(r.Replace(s))
	runes := []rune(clean)
	if len(runes) > 48 {
		return string(runes[:48])
	}
	return clean
}

func escapeFFmpegPath(p string) string {
	slash := filepath.ToSlash(p)
	return strings.ReplaceAll(slash, ":", "\\:")
}

func findSystemFontPath() string {
	return findSystemFontPathForOS(runtime.GOOS, isRegularFile)
}

func isRegularFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func findSystemFontPathForOS(goos string, existsFn func(string) bool) string {
	for _, c := range fontCandidatesForOS(goos) {
		if existsFn != nil && existsFn(c) {
			return c
		}
	}
	return ""
}

func fontCandidatesForOS(goos string) []string {
	switch goos {
	case "windows":
		return []string{
			`C:\Windows\Fonts\arial.ttf`,
			`C:\Windows\Fonts\segoeui.ttf`,
		}
	case "darwin":
		return []string{
			"/System/Library/Fonts/Helvetica.ttc",
			"/Library/Fonts/Arial.ttf",
		}
	default:
		return []string{
			"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
			"/usr/share/fonts/TTF/DejaVuSans.ttf",
		}
	}
}

func pythonForOS(goos string) string {
	if goos == "windows" {
		return "python"
	}
	return "python3"
}

func shellForOS(goos string) []string {
	if goos == "windows" {
		return []string{"powershell", "-NoProfile"}
	}
	return []string{"sh", "-c"}
}

// GenerateCartoonFixture3s produz um arquivo MP4 de 3 segundos real via FFmpeg sem fallback placeholder.
func GenerateCartoonFixture3s(outMP4Path string) error {
	return generateCartoonFixture3sWith(outMP4Path, exec.LookPath, runOSCombined)
}

func generateCartoonFixture3sWith(outMP4Path string, lookPath func(string) (string, error), runFn func(bin string, args ...string) ([]byte, error)) error {
	bin, err := lookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("%w: ffmpeg ausente no PATH: %v", ErrRenderFailed, err)
	}
	if err := os.MkdirAll(filepath.Dir(outMP4Path), 0755); err != nil {
		return fmt.Errorf("%w: erro ao criar diretório: %v", ErrRenderFailed, err)
	}
	args := []string{
		"-y", "-f", "lavfi", "-i", "color=c=0x16a34a:s=1080x1920:r=30:d=3",
		"-f", "lavfi", "-i", "anullsrc=r=44100:cl=stereo:d=3",
		"-c:v", "libx264", "-c:a", "aac", "-pix_fmt", "yuv420p", outMP4Path,
	}
	if out, runErr := runFn(bin, args...); runErr != nil {
		_ = os.Remove(outMP4Path)
		return fmt.Errorf("%w: falha ao gerar fixture MP4 com ffmpeg (%v): %s", ErrRenderFailed, runErr, trimRenderOutput(out))
	}
	return verifyCartoonMP4(outMP4Path)
}

func runOSCombined(bin string, args ...string) ([]byte, error) {
	return exec.Command(bin, args...).CombinedOutput()
}

func verifyCartoonMP4(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("%w: arquivo MP4 não foi gerado em %s: %v", ErrRenderFailed, path, err)
	}
	if info.Size() <= 32 {
		_ = os.Remove(path)
		return fmt.Errorf("%w: arquivo MP4 inválido (%d bytes) em %s", ErrRenderFailed, info.Size(), path)
	}
	return nil
}

func trimRenderOutput(out []byte) string {
	s := strings.TrimSpace(string(out))
	if s == "" {
		return "sem saída"
	}
	if len(s) > 400 {
		return s[len(s)-400:]
	}
	return s
}

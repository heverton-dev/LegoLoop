package motion

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderMotion_OrchestrationAndErrors(t *testing.T) {
	t.Run("rejeita request invalida", func(t *testing.T) {
		_, err := RenderMotion(Request{}, nil, t.TempDir(), filepath.Join(t.TempDir(), "out.mp4"))
		if !errors.Is(err, ErrInvalidRequest) {
			t.Fatalf("esperava ErrInvalidRequest, obteve %v", err)
		}
	})

	t.Run("orquestra processo completo com audio e beats via seams", func(t *testing.T) {
		tmp := t.TempDir()
		audio := filepath.Join(tmp, "beat.wav")
		if err := os.WriteFile(audio, []byte("wav"), 0644); err != nil {
			t.Fatal(err)
		}
		outDir := filepath.Join(tmp, "frames")
		outMP4 := filepath.Join(tmp, "video.mp4")

		req := Request{
			Prompt:     "grafico animado",
			Scenario:   ScenarioExplainer,
			Resolution: "1280x720",
			Duration:   1.0,
			AudioPath:  audio,
		}
		res, err := renderMotionWith(req, nil, outDir, outMP4, motionDeps{
			capturer: &stubCapturer{},
			encoder:  &stubEncoder{},
			detector: &stubDetector{beats: []float64{0.0, 0.5, 1.0}},
			fps:      2,
		})
		if err != nil {
			t.Fatalf("renderMotionWith falhou: %v", err)
		}
		if !strings.Contains(res.HTML, "<canvas") {
			t.Errorf("HTML sem canvas: %s", res.HTML)
		}
		if len(res.Repos) == 0 {
			t.Errorf("esperava repos sugeridos automaticamente")
		}
		if len(res.Frames) != 2 {
			t.Errorf("frames = %d, want 2", len(res.Frames))
		}
		if len(res.Beats) != 3 {
			t.Errorf("beats = %v, want 3 itens", res.Beats)
		}
		if _, err := os.Stat(outMP4); err != nil {
			t.Errorf("MP4 nao gerado em %s: %v", outMP4, err)
		}
	})

	t.Run("propaga erro do detector de beats sem fallback silencioso", func(t *testing.T) {
		tmp := t.TempDir()
		audio := filepath.Join(tmp, "beat.wav")
		if err := os.WriteFile(audio, []byte("wav"), 0644); err != nil {
			t.Fatal(err)
		}
		req := Request{
			Prompt:     "lancamento",
			Scenario:   ScenarioProductLaunch,
			Resolution: "1280x720",
			Duration:   1.0,
			AudioPath:  audio,
		}
		_, err := renderMotionWith(req, []string{"threejs"}, tmp, filepath.Join(tmp, "out.mp4"), motionDeps{
			capturer: &stubCapturer{},
			encoder:  &stubEncoder{},
			detector: &stubDetector{err: ErrBeatDetectFailed},
			fps:      1,
		})
		if !errors.Is(err, ErrBeatDetectFailed) {
			t.Fatalf("esperava ErrBeatDetectFailed, obteve %v", err)
		}
	})

	t.Run("funcoes exportadas de captura encode e detect falham explicitamente sem binarios", func(t *testing.T) {
		t.Setenv("PATH", t.TempDir())
		if _, err := CaptureFrames("<html></html>", 1, 1.0, "1280x720", t.TempDir()); !errors.Is(err, ErrRenderFailed) {
			t.Errorf("CaptureFrames sem playwright: got %v", err)
		}
		frame := filepath.Join(t.TempDir(), "f.png")
		_ = os.WriteFile(frame, []byte("png"), 0644)
		if err := EncodeFrames([]string{frame}, "", filepath.Join(t.TempDir(), "out.mp4"), "1280x720"); !errors.Is(err, ErrRenderFailed) {
			t.Errorf("EncodeFrames sem ffmpeg: got %v", err)
		}
		audio := filepath.Join(t.TempDir(), "a.wav")
		_ = os.WriteFile(audio, []byte("wav"), 0644)
		if _, err := DetectBeats(audio); !errors.Is(err, ErrBeatDetectFailed) {
			t.Errorf("DetectBeats sem python: got %v", err)
		}
	})
}

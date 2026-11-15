package motion

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type spyCommandRunner struct {
	lookPathFn func(file string) (string, error)
	runFn      func(bin string, args ...string) ([]byte, error)
	lookedUp   []string
	calls      []spyCmdCall
}

type spyCmdCall struct {
	bin  string
	args []string
}

func (s *spyCommandRunner) LookPath(file string) (string, error) {
	s.lookedUp = append(s.lookedUp, file)
	if s.lookPathFn != nil {
		return s.lookPathFn(file)
	}
	return "/usr/bin/" + file, nil
}

func (s *spyCommandRunner) Run(bin string, args ...string) ([]byte, error) {
	cp := append([]string(nil), args...)
	s.calls = append(s.calls, spyCmdCall{bin: bin, args: cp})
	if s.runFn != nil {
		return s.runFn(bin, args...)
	}
	return nil, nil
}

func TestEncoderERendererSemStubPlaceholder(t *testing.T) {
	t.Run("encoder falha explicitamente sem ffmpeg no PATH e nao grava placeholder", func(t *testing.T) {
		t.Setenv("PATH", t.TempDir())
		tmp := t.TempDir()
		frame := filepath.Join(tmp, "f0.png")
		if err := os.WriteFile(frame, []byte("png-data"), 0644); err != nil {
			t.Fatal(err)
		}
		out := filepath.Join(tmp, "out.mp4")

		enc := &ffmpegEncoder{}
		err := enc.Encode([]string{frame}, "", out, "1920x1080")
		if err == nil {
			t.Fatalf("esperava erro quando ffmpeg esta ausente no PATH, obteve nil")
		}
		if !errors.Is(err, ErrRenderFailed) {
			t.Errorf("erro = %v, esperava ErrRenderFailed", err)
		}
		if !strings.Contains(strings.ToLower(err.Error()), "ffmpeg") {
			t.Errorf("mensagem de erro deveria mencionar ffmpeg: %v", err)
		}
		if _, statErr := os.Stat(out); !os.IsNotExist(statErr) {
			t.Errorf("arquivo placeholder nao deveria existir em %s", out)
		}
	})

	t.Run("encoder propaga falha real do comando ffmpeg e limpa saida", func(t *testing.T) {
		tmp := t.TempDir()
		out := filepath.Join(tmp, "out.mp4")
		_ = os.WriteFile(out, []byte("parcial"), 0644)

		spy := &spyCommandRunner{
			runFn: func(bin string, args ...string) ([]byte, error) {
				return []byte("codec libx264 error"), errors.New("exit status 1")
			},
		}
		enc := &ffmpegEncoder{runner: spy}
		err := enc.Encode([]string{"a.png", "b.png"}, "audio.wav", out, "1280x720")
		if !errors.Is(err, ErrRenderFailed) {
			t.Fatalf("erro = %v, esperava ErrRenderFailed", err)
		}
		if !strings.Contains(err.Error(), "codec libx264 error") {
			t.Errorf("erro deveria incluir saida detalhada do ffmpeg: %v", err)
		}
		if _, statErr := os.Stat(out); !os.IsNotExist(statErr) {
			t.Errorf("arquivo corrompido deveria ser removido apos falha")
		}
	})

	t.Run("encoder valida argumentos reais sem shell e exige arquivo nao vazio", func(t *testing.T) {
		tmp := t.TempDir()
		out := filepath.Join(tmp, "out.mp4")
		spy := &spyCommandRunner{
			runFn: func(bin string, args ...string) ([]byte, error) {
				return nil, os.WriteFile(out, make([]byte, 512), 0644)
			},
		}
		enc := &ffmpegEncoder{runner: spy}
		if err := enc.Encode([]string{"f0.png", "f1.png"}, "voice.wav", out, "1920x1080"); err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if len(spy.lookedUp) != 1 || spy.lookedUp[0] != "ffmpeg" {
			t.Fatalf("LookPath = %v, want [ffmpeg]", spy.lookedUp)
		}
		if len(spy.calls) != 1 {
			t.Fatalf("esperava 1 chamada ao runner, obteve %d", len(spy.calls))
		}
		joined := strings.Join(spy.calls[0].args, " ")
		for _, want := range []string{"-y", "f0.png", "f1.png", "voice.wav", "concat=n=2", "libx264", "1920x1080", out} {
			if !strings.Contains(joined, want) {
				t.Errorf("args ffmpeg ausente %q em: %s", want, joined)
			}
		}
	})

	t.Run("encoder falha se ffmpeg retornar exit 0 sem gerar arquivo mp4", func(t *testing.T) {
		tmp := t.TempDir()
		out := filepath.Join(tmp, "missing.mp4")
		spy := &spyCommandRunner{}
		enc := &ffmpegEncoder{runner: spy}
		err := enc.Encode([]string{"f0.png"}, "", out, "1920x1080")
		if !errors.Is(err, ErrRenderFailed) {
			t.Fatalf("esperava ErrRenderFailed para arquivo nao criado, obteve %v", err)
		}
	})

	t.Run("renderer falha explicitamente sem playwright no PATH e nao grava PNG 1x1", func(t *testing.T) {
		t.Setenv("PATH", t.TempDir())
		outDir := t.TempDir()
		cap := &playwrightCapturer{}
		frames, err := cap.Capture("<html><body>teste</body></html>", 2, 1.0, "1280x720", outDir)
		if err == nil {
			t.Fatalf("esperava erro quando playwright esta ausente no PATH, obteve nil (frames=%v)", frames)
		}
		if !errors.Is(err, ErrRenderFailed) {
			t.Errorf("erro = %v, esperava ErrRenderFailed", err)
		}
		if !strings.Contains(strings.ToLower(err.Error()), "playwright") {
			t.Errorf("erro deveria mencionar playwright: %v", err)
		}
		fp0 := filepath.Join(outDir, "frame_0000.png")
		if _, statErr := os.Stat(fp0); !os.IsNotExist(statErr) {
			t.Errorf("PNG placeholder nao deveria ser gravado em %s", fp0)
		}
	})

	t.Run("renderer propaga erro do playwright e passa argumentos reais sem shell", func(t *testing.T) {
		outDir := t.TempDir()
		spy := &spyCommandRunner{
			runFn: func(bin string, args ...string) ([]byte, error) {
				return []byte("Executable doesn't exist"), errors.New("exit status 1")
			},
		}
		cap := &playwrightCapturer{runner: spy}
		_, err := cap.Capture("<html></html>", 2, 1.0, "1280x720", outDir)
		if !errors.Is(err, ErrRenderFailed) {
			t.Fatalf("erro = %v, esperava ErrRenderFailed", err)
		}
		if !strings.Contains(err.Error(), "Executable doesn't exist") {
			t.Errorf("erro deveria incluir saida do playwright: %v", err)
		}
		if len(spy.calls) != 1 {
			t.Fatalf("deveria interromper no primeiro frame com erro, chamadas=%d", len(spy.calls))
		}
		joined := strings.Join(spy.calls[0].args, " ")
		for _, want := range []string{"screenshot", "--viewport-size=1280,720", "file://", "frame_0000.png"} {
			if !strings.Contains(joined, want) {
				t.Errorf("arg playwright ausente %q em: %s", want, joined)
			}
		}
	})

	t.Run("renderer captura todos os frames via seam e falha se PNG nao existir", func(t *testing.T) {
		outDir := t.TempDir()
		spy := &spyCommandRunner{
			runFn: func(bin string, args ...string) ([]byte, error) {
				outPNG := args[len(args)-1]
				return nil, os.WriteFile(outPNG, []byte("real-png-bytes-123456"), 0644)
			},
		}
		cap := &playwrightCapturer{runner: spy}
		frames, err := cap.Capture("<html></html>", 2, 1.5, "1920x1080", outDir)
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if len(frames) != 3 || len(spy.calls) != 3 {
			t.Fatalf("frames=%d calls=%d, want 3", len(frames), len(spy.calls))
		}

		spyMissing := &spyCommandRunner{}
		capMissing := &playwrightCapturer{runner: spyMissing}
		if _, err := capMissing.Capture("<html></html>", 1, 1.0, "1920x1080", t.TempDir()); !errors.Is(err, ErrRenderFailed) {
			t.Errorf("esperava ErrRenderFailed quando PNG nao e escrito, obteve %v", err)
		}
	})

	t.Run("encoder real com ffmpeg do sistema gera MP4 valido", func(t *testing.T) {
		if _, err := exec.LookPath("ffmpeg"); err != nil {
			t.Skip("ffmpeg nao disponivel neste host")
		}
		tmp := t.TempDir()
		p1 := filepath.Join(tmp, "f0.png")
		p2 := filepath.Join(tmp, "f1.png")
		for idx, spec := range []struct {
			path  string
			color string
		}{{p1, "red"}, {p2, "blue"}} {
			cmd := exec.Command("ffmpeg", "-y", "-f", "lavfi", "-i", "color=c="+spec.color+":s=64x64:d=0.1", "-frames:v", "1", spec.path)
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("falha gerando frame %d: %v (%s)", idx, err, string(out))
			}
		}
		out := filepath.Join(tmp, "real.mp4")
		enc := &ffmpegEncoder{}
		if err := enc.Encode([]string{p1, p2}, "", out, "64x64"); err != nil {
			t.Fatalf("Encode real falhou: %v", err)
		}
		data, err := os.ReadFile(out)
		if err != nil {
			t.Fatal(err)
		}
		if len(data) < 512 || bytes.Equal(data, []byte("mp4-placeholder")) || !bytes.Contains(data[:32], []byte("ftyp")) {
			t.Fatalf("arquivo MP4 real invalido: len=%d", len(data))
		}
	})
}

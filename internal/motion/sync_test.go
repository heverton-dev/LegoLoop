package motion

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
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

func TestLibrosaDetector_RealExecAndErrors(t *testing.T) {
	t.Run("falha explicitamente quando python3 e python estao ausentes no PATH", func(t *testing.T) {
		t.Setenv("PATH", t.TempDir())
		audio := filepath.Join(t.TempDir(), "audio.wav")
		_ = os.WriteFile(audio, []byte("wav"), 0644)

		d := &librosaDetector{}
		beats, err := d.Detect(audio)
		if err == nil {
			t.Fatalf("esperava erro explicito sem python no PATH, obteve beats=%v", beats)
		}
		if !errors.Is(err, ErrBeatDetectFailed) {
			t.Errorf("erro = %v, want ErrBeatDetectFailed", err)
		}
		if !strings.Contains(strings.ToLower(err.Error()), "python") {
			t.Errorf("mensagem deveria mencionar python: %v", err)
		}
	})

	t.Run("falha explicitamente quando script beat_detect.py nao existe", func(t *testing.T) {
		audio := filepath.Join(t.TempDir(), "audio.wav")
		_ = os.WriteFile(audio, []byte("wav"), 0644)
		spy := &spyCommandRunner{}
		d := &librosaDetector{runner: spy, scriptPath: filepath.Join(t.TempDir(), "missing.py")}

		_, err := d.Detect(audio)
		if !errors.Is(err, ErrBeatDetectFailed) {
			t.Fatalf("esperava ErrBeatDetectFailed para script ausente, obteve %v", err)
		}
		if len(spy.calls) != 0 {
			t.Errorf("nao deveria executar python quando script nao existe, calls=%d", len(spy.calls))
		}
	})

	t.Run("tenta python3 e faz fallback para python passando script e audio sem shell", func(t *testing.T) {
		tmp := t.TempDir()
		audio := filepath.Join(tmp, "audio.wav")
		script := filepath.Join(tmp, "beat_detect.py")
		_ = os.WriteFile(audio, []byte("wav"), 0644)
		_ = os.WriteFile(script, []byte("print('[0.5]')"), 0644)

		spy := &spyCommandRunner{
			lookPathFn: func(file string) (string, error) {
				if file == "python3" {
					return "", errors.New("not found")
				}
				return "/usr/bin/python", nil
			},
			runFn: func(bin string, args ...string) ([]byte, error) {
				return []byte("[0.0, 0.5, 1.0, 1.5]\n"), nil
			},
		}
		d := &librosaDetector{runner: spy, scriptPath: script}
		beats, err := d.Detect(audio)
		if err != nil {
			t.Fatalf("Detect erro inesperado: %v", err)
		}
		if len(beats) != 4 || beats[1] != 0.5 {
			t.Fatalf("beats = %v, want [0 0.5 1 1.5]", beats)
		}
		if len(spy.lookedUp) != 2 || spy.lookedUp[0] != "python3" || spy.lookedUp[1] != "python" {
			t.Fatalf("LookPath = %v, want [python3 python]", spy.lookedUp)
		}
		if len(spy.calls) != 1 || spy.calls[0].bin != "/usr/bin/python" || len(spy.calls[0].args) != 2 || spy.calls[0].args[0] != script || spy.calls[0].args[1] != audio {
			t.Fatalf("chamada invalida: %+v", spy.calls)
		}
	})

	t.Run("propaga erro quando python falha ou retorna JSON invalido sem fallback silencioso", func(t *testing.T) {
		tmp := t.TempDir()
		audio := filepath.Join(tmp, "audio.wav")
		script := filepath.Join(tmp, "beat_detect.py")
		_ = os.WriteFile(audio, []byte("wav"), 0644)
		_ = os.WriteFile(script, []byte(""), 0644)

		dErr := &librosaDetector{
			runner: &spyCommandRunner{
				runFn: func(bin string, args ...string) ([]byte, error) {
					return []byte("Traceback: No module named librosa"), errors.New("exit status 1")
				},
			},
			scriptPath: script,
		}
		if _, err := detectBeatsWith(dErr, audio); !errors.Is(err, ErrBeatDetectFailed) || !strings.Contains(err.Error(), "No module named librosa") {
			t.Fatalf("detectBeatsWith deveria propagar falha do detector: %v", err)
		}

		dBadJSON := &librosaDetector{
			runner: &spyCommandRunner{
				runFn: func(bin string, args ...string) ([]byte, error) {
					return []byte("{bad json"), nil
				},
			},
			scriptPath: script,
		}
		if _, err := detectBeatsWith(dBadJSON, audio); !errors.Is(err, ErrBeatDetectFailed) {
			t.Fatalf("detectBeatsWith deveria falhar para JSON invalido: %v", err)
		}
	})
}

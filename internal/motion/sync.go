package motion

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ErrBeatDetectFailed sinaliza falha explicita na deteccao de beats de audio.
var ErrBeatDetectFailed = errors.New("motion: falha na deteccao de beats")

// DetectBeats analisa beats do arquivo de audio usando scripts/beat_detect.py.
func DetectBeats(audioPath string) ([]float64, error) {
	return detectBeatsWith(&librosaDetector{}, audioPath)
}

// detectBeatsWith detects beats using the given detector (seam for testing).
// Returns nil when audio is empty or missing; propagates detector failures.
func detectBeatsWith(d beatDetector, audioPath string) ([]float64, error) {
	if audioPath == "" {
		return nil, nil
	}
	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		return nil, nil
	}
	return d.Detect(audioPath)
}

type osStdoutRunner struct{}

func (osStdoutRunner) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

func (osStdoutRunner) Run(bin string, args ...string) ([]byte, error) {
	out, err := exec.Command(bin, args...).Output()
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) && len(ee.Stderr) > 0 {
			return ee.Stderr, err
		}
		return out, err
	}
	return out, nil
}

// librosaDetector invokes scripts/beat_detect.py via Python.
type librosaDetector struct {
	runner     commandRunner
	scriptPath string
}

func (l *librosaDetector) effectiveRunner() commandRunner {
	if l != nil && l.runner != nil {
		return l.runner
	}
	return osStdoutRunner{}
}

// Detect calls Python scripts/beat_detect.py without shell and parses JSON beats.
func (l *librosaDetector) Detect(audioPath string) ([]float64, error) {
	script, err := resolveBeatScript(l.scriptPath)
	if err != nil {
		return nil, err
	}
	r := l.effectiveRunner()
	bin, err := resolvePythonBin(r)
	if err != nil {
		return nil, err
	}
	out, err := r.Run(bin, script, audioPath)
	if err != nil {
		return nil, formatBeatExecError(err, out)
	}
	var beats []float64
	if err := json.Unmarshal(bytes.TrimSpace(out), &beats); err != nil {
		return nil, fmt.Errorf("%w: JSON invalido de beat_detect: %v", ErrBeatDetectFailed, err)
	}
	return beats, nil
}

func resolveBeatScript(custom string) (string, error) {
	candidates := []string{filepath.Join("scripts", "beat_detect.py"), filepath.Join("..", "..", "scripts", "beat_detect.py")}
	if custom != "" {
		candidates = []string{custom}
	}
	for _, path := range candidates {
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return path, nil
		}
	}
	return "", fmt.Errorf("%w: script beat_detect.py ausente (%v)", ErrBeatDetectFailed, candidates)
}

func resolvePythonBin(r commandRunner) (string, error) {
	var lastErr error
	for _, cand := range []string{"python3", "python"} {
		bin, err := r.LookPath(cand)
		if err == nil {
			return bin, nil
		}
		lastErr = err
	}
	return "", fmt.Errorf("%w: python3/python ausente no PATH: %v", ErrBeatDetectFailed, lastErr)
}

func formatBeatExecError(err error, out []byte) error {
	msg := trimCmdOutput(out)
	if msg != "" {
		return fmt.Errorf("%w: %v: %s", ErrBeatDetectFailed, err, msg)
	}
	return fmt.Errorf("%w: %v", ErrBeatDetectFailed, err)
}

// beatDetector is the seam for Librosa — injectable for testing.
type beatDetector interface {
	Detect(audioPath string) ([]float64, error)
}

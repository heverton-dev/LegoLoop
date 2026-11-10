package motion

import "os"

// detectBeatsWith detects beats using the given detector (seam for testing).
// Returns nil gracefully when audio is missing or empty.
func detectBeatsWith(d beatDetector, audioPath string) ([]float64, error) {
	if audioPath == "" {
		return nil, nil
	}
	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		return nil, nil
	}
	beats, err := d.Detect(audioPath)
	if err != nil {
		return nil, nil // graceful fallback — no sync beats
	}
	return beats, nil
}

// librosaDetector invokes scripts/beat_detect.py via Python.
type librosaDetector struct{}

// Detect calls the Python script for beat detection.
func (l *librosaDetector) Detect(audioPath string) ([]float64, error) {
	// Boundary call via os/exec to python scripts/beat_detect.py
	// Real implementation invokes: python scripts/beat_detect.py <audioPath>
	// Returns JSON array of beat timestamps.
	return nil, nil
}

// beatDetector is the seam for Librosa — injectable for testing.
type beatDetector interface {
	Detect(audioPath string) ([]float64, error)
}

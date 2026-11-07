package legopress

import (
	"strings"
)

type SyncEngine struct{}

func NewSyncEngine() *SyncEngine {
	return &SyncEngine{}
}

func (s *SyncEngine) DetectStaleParts(diffs []string, m Manifest) []Part {
	hasCodeDiffs := false
	hasProcessDiffs := false

	for _, d := range diffs {
		norm := strings.ToLower(d)
		if strings.HasPrefix(norm, "internal/") || strings.HasPrefix(norm, "cmd/") {
			hasCodeDiffs = true
		}
		if strings.Contains(norm, "intent") || strings.Contains(norm, "spec") || strings.Contains(norm, "route") {
			hasProcessDiffs = true
		}
	}

	var stale []Part
	for _, p := range m.Parts {
		fileLower := strings.ToLower(p.File)
		titleLower := strings.ToLower(p.Title)

		if hasCodeDiffs && (strings.Contains(fileLower, "micro") || strings.Contains(titleLower, "micro") || strings.Contains(fileLower, "motor")) {
			stale = append(stale, p)
		} else if hasProcessDiffs && (strings.Contains(fileLower, "meso") || strings.Contains(titleLower, "meso")) {
			stale = append(stale, p)
		}
	}
	return stale
}

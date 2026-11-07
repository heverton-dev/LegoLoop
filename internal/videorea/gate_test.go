package videorea_test

import (
	"testing"

	"legoloop/internal/videorea"
)

func TestStopGateCuration(t *testing.T) {
	report := &videorea.CurationReport{
		SuggestedName: "legoloop-video-talking-head",
		OriginalCost:  "R$ 180,00/mes",
		SovereignCost: "R$ 0,00/mes",
		Replacements: map[string]string{
			"Higgsfield": "Google AI Studio",
		},
	}

	// Caso de recusa
	err := videorea.EvaluateCurationGate(report, false)
	if err != videorea.ErrGateRejected {
		t.Errorf("esperava ErrGateRejected quando aprovacao e false, obteve: %v", err)
	}

	// Caso de aprovação
	err = videorea.EvaluateCurationGate(report, true)
	if err != nil {
		t.Fatalf("EvaluateCurationGate com true nao deveria falhar: %v", err)
	}
	if !report.Approved {
		t.Errorf("report.Approved deveria ser true apos o gate")
	}
}

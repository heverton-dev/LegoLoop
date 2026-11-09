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

func TestEvaluateStopGateAllTypes(t *testing.T) {
	gates := []videorea.StopGateType{
		videorea.GateCuradoria,
		videorea.GateIntent,
		videorea.GateSpec,
		videorea.GateRoute,
		videorea.GateDAG,
	}

	for _, g := range gates {
		// Rejeição
		d, err := videorea.EvaluateStopGate(g, false, "ajustar threshold")
		if err != videorea.ErrGateRejected {
			t.Errorf("[%s] esperava ErrGateRejected, obteve: %v", g, err)
		}
		if d.Approved {
			t.Errorf("[%s] Approved deveria ser false", g)
		}
		if d.Notes != "ajustar threshold" {
			t.Errorf("[%s] Notes incorreto: %s", g, d.Notes)
		}

		// Aprovação
		d, err = videorea.EvaluateStopGate(g, true, "aprovado pelo operador")
		if err != nil {
			t.Errorf("[%s] erro inesperado na aprovacao: %v", g, err)
		}
		if !d.Approved {
			t.Errorf("[%s] Approved deveria ser true", g)
		}
		if d.Gate != g {
			t.Errorf("[%s] Gate retornado incorreto: %s", g, d.Gate)
		}
	}
}

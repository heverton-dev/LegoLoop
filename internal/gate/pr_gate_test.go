package gate_test

import (
	"testing"

	"legoloop/internal/gate"
)

func TestEvaluatePRMergeGate_AllPassed(t *testing.T) {
	checks := gate.PRMergeChecks{
		BuildClean:    true,
		TestsPassed:   true,
		ZeroWarnings:  true,
		SecurityClean: true,
	}

	result := gate.EvaluatePRMergeGate(checks)
	if !result.Approved {
		t.Fatalf("esperava gate aprovado com todos checks true, motivo: %s", result.Reason)
	}
}

func TestEvaluatePRMergeGate_FailsOnSecurityOrTests(t *testing.T) {
	checks := gate.PRMergeChecks{
		BuildClean:    true,
		TestsPassed:   false,
		ZeroWarnings:  true,
		SecurityClean: true,
	}

	result := gate.EvaluatePRMergeGate(checks)
	if result.Approved {
		t.Fatalf("esperava reprovação do gate devido a falha nos testes")
	}

	checks2 := gate.PRMergeChecks{
		BuildClean:    true,
		TestsPassed:   true,
		ZeroWarnings:  true,
		SecurityClean: false, // Vulnerabilidade encontrada em pr-security-review
	}

	result2 := gate.EvaluatePRMergeGate(checks2)
	if result2.Approved {
		t.Fatalf("esperava reprovação do gate devido a falha de segurança")
	}
}

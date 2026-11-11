package plangen_test

import (
	"testing"

	"legoloop/internal/plangen"
)

func TestConcurrencyStrategy_Validation(t *testing.T) {
	// 1. Configuração válida padrão
	valid := plangen.ConcurrencyStrategy{
		Mode:                 "hybrid",
		MaxParallelSubagents: 4,
		HarnessCapability:    "subagent_cli_isolated",
		RateLimitThrottleMS:  500,
	}
	if err := plangen.ValidateConcurrency(valid); err != nil {
		t.Fatalf("esperava sucesso para configuração válida, obteve: %v", err)
	}

	// 2. Modo inválido
	invalidMode := valid
	invalidMode.Mode = "invalid_mode_xyz"
	if err := plangen.ValidateConcurrency(invalidMode); err == nil {
		t.Fatal("esperava erro para modo de concorrência inválido, obteve nil")
	}

	// 3. Limite estourado (> 8)
	overflow := valid
	overflow.MaxParallelSubagents = 10
	if err := plangen.ValidateConcurrency(overflow); err == nil {
		t.Fatal("esperava erro para subagentes > 8, obteve nil")
	}

	// 4. Limite zerado ou negativo
	underflow := valid
	underflow.MaxParallelSubagents = 0
	if err := plangen.ValidateConcurrency(underflow); err == nil {
		t.Fatal("esperava erro para subagentes <= 0, obteve nil")
	}
}

func TestConcurrencyStrategy_GracefulDegradation(t *testing.T) {
	degraded := plangen.ConcurrencyStrategy{
		Mode:                 "hybrid",
		MaxParallelSubagents: 4,
		HarnessCapability:    "none",
	}
	if eff := plangen.EffectiveParallelism(degraded); eff != 1 {
		t.Fatalf("esperava paralelismo efetivo 1 em caso de harness_capability 'none', obteve %d", eff)
	}

	serial := plangen.ConcurrencyStrategy{
		Mode:                 "serial",
		MaxParallelSubagents: 4,
		HarnessCapability:    "subagent_cli_isolated",
	}
	if eff := plangen.EffectiveParallelism(serial); eff != 1 {
		t.Fatalf("esperava paralelismo efetivo 1 para mode 'serial', obteve %d", eff)
	}

	capable := plangen.ConcurrencyStrategy{
		Mode:                 "hybrid",
		MaxParallelSubagents: 3,
		HarnessCapability:    "subagent_cli_isolated",
	}
	if eff := plangen.EffectiveParallelism(capable); eff != 3 {
		t.Fatalf("esperava paralelismo efetivo 3, obteve %d", eff)
	}
}

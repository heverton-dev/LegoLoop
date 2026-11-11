package plangen

import (
	"errors"
	"fmt"
	"slices"
)

var (
	ErrConcurrencyMode  = errors.New("plangen: modo de concorrência inválido")
	ErrConcurrencyLimit = errors.New("plangen: teto de subagentes fora do limite permitido [1, 8]")
)

var (
	allowedModes        = []string{"serial", "parallel_clean", "hybrid"}
	allowedCapabilities = []string{"none", "subagent_cli_isolated", "subagent_native"}
)

// ValidateConcurrency confere se a estratégia de concorrência é válida e segura.
func ValidateConcurrency(c ConcurrencyStrategy) error {
	if !slices.Contains(allowedModes, c.Mode) {
		return fmt.Errorf("%w: %q (esperado: %v)", ErrConcurrencyMode, c.Mode, allowedModes)
	}
	if c.MaxParallelSubagents < 1 || c.MaxParallelSubagents > 8 {
		return fmt.Errorf("%w: %d", ErrConcurrencyLimit, c.MaxParallelSubagents)
	}
	if c.HarnessCapability != "" && !slices.Contains(allowedCapabilities, c.HarnessCapability) {
		return fmt.Errorf("plangen: harness_capability inválido: %q", c.HarnessCapability)
	}
	return nil
}

// EffectiveParallelism calcula o número real de workers paralelos com degradação graciosa.
func EffectiveParallelism(c ConcurrencyStrategy) int {
	if c.HarnessCapability == "none" || c.Mode == "serial" {
		return 1
	}
	if c.MaxParallelSubagents < 1 {
		return 1
	}
	if c.MaxParallelSubagents > 8 {
		return 8
	}
	return c.MaxParallelSubagents
}

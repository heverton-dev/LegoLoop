package runner

import (
	"fmt"
	"strings"
)

// ValidateRedBlueBoundary aplica a regra criptográfica Anti-Reward Hacking (Motor I).
// Agentes com papel Blue-Team são proibidos de declarar ou gravar em arquivos de teste (*_test.go).
func ValidateRedBlueBoundary(role string, contextIn, contextOut []string) error {
	if role != "blue_team" {
		return nil
	}

	for _, out := range contextOut {
		if isTestFile(out) {
			return fmt.Errorf("reward_hacking: agente blue_team proibido de modificar arquivo de teste %q", out)
		}
	}

	return nil
}

func isTestFile(path string) bool {
	norm := strings.ToLower(path)
	return strings.HasSuffix(norm, "_test.go") || strings.HasSuffix(norm, ".test.js") || strings.HasSuffix(norm, ".test.ts")
}

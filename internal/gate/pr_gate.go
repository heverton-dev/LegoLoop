package gate

import (
	"strings"
)

// PRMergeChecks representa os pré-requisitos determinísticos para autorização de merge automático.
type PRMergeChecks struct {
	BuildClean    bool `json:"build_clean"`
	TestsPassed   bool `json:"tests_passed"`
	ZeroWarnings  bool `json:"zero_warnings"`
	SecurityClean bool `json:"security_clean"`
}

// PRMergeVerdict é a decisão binária do portão de fusão.
type PRMergeVerdict struct {
	Approved bool     `json:"approved"`
	Reason   string   `json:"reason"`
	Failures []string `json:"failures,omitempty"`
}

// EvaluatePRMergeGate avalia se uma Pull Request atende a todos os critérios inegociáveis para auto-merge.
func EvaluatePRMergeGate(checks PRMergeChecks) PRMergeVerdict {
	var failures []string

	if !checks.BuildClean {
		failures = append(failures, "compilação com falhas ou warnings bloqueantes")
	}
	if !checks.TestsPassed {
		failures = append(failures, "suíte de testes unitários ou de integração falhou")
	}
	if !checks.ZeroWarnings {
		failures = append(failures, "lint ou analisador estático reportou advertências")
	}
	if !checks.SecurityClean {
		failures = append(failures, "revisão de segurança (pr-security-review) detectou potenciais vulnerabilidades")
	}

	if len(failures) > 0 {
		return PRMergeVerdict{
			Approved: false,
			Reason:   strings.Join(failures, "; "),
			Failures: failures,
		}
	}

	return PRMergeVerdict{
		Approved: true,
		Reason:   "todos os requisitos de build, testes, lint e segurança foram cumpridos",
	}
}

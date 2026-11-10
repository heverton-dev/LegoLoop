package orchestrix

import (
	"strings"
	"testing"
)

// O claude filho precisa de CLAUDE_CONFIG_DIR para achar o login do perfil; segredos continuam fora.
func TestSterileEnvPassesClaudeConfigDirOnly(t *testing.T) {
	vars := map[string]string{
		"PATH":                  `C:\bin`,
		"CLAUDE_CONFIG_DIR":     `C:\perfil\claude`,
		"ANTHROPIC_API_KEY":     "segredo",
		"ORCA_AGENT_HOOK_TOKEN": "segredo",
	}
	env := sterileEnv(func(k string) string { return vars[k] })
	joined := "\n" + strings.Join(env, "\n") + "\n"
	if !strings.Contains(joined, "\nCLAUDE_CONFIG_DIR="+`C:\perfil\claude`+"\n") {
		t.Fatalf("CLAUDE_CONFIG_DIR ausente: %v", env)
	}
	for _, bad := range []string{"ANTHROPIC_API_KEY", "ORCA_AGENT_HOOK_TOKEN"} {
		if strings.Contains(joined, "\n"+bad+"=") {
			t.Fatalf("%s vazou para o filho", bad)
		}
	}
}

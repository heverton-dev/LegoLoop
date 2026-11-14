package runner

import (
	"fmt"
	"strings"

	"legoloop/internal/dag"
)

// ValidateToolAllowlist verifica se a ferramenta solicitada consta na allowlist do nó.
func ValidateToolAllowlist(n dag.Node, requestedTool string) error {
	allowed := extractAllowedTools(n.Tooling)
	if len(allowed) == 0 {
		return nil
	}
	target := strings.TrimSpace(requestedTool)
	for _, a := range allowed {
		if a == target {
			return nil
		}
	}
	return fmt.Errorf("allowlist: ferramenta %q não permitida no nó %s", target, n.ID)
}

func extractAllowedTools(tooling map[string]any) []string {
	if tooling == nil {
		return nil
	}
	for _, key := range []string{"allowlist", "tools"} {
		if raw, ok := tooling[key]; ok {
			return toStringSlice(raw)
		}
	}
	return nil
}

func toStringSlice(val any) []string {
	switch v := val.(type) {
	case []string:
		return v
	case []any:
		out := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok && strings.TrimSpace(s) != "" {
				out = append(out, strings.TrimSpace(s))
			}
		}
		return out
	}
	return nil
}

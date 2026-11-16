package cli

import (
	"sort"

	"legoloop/internal/mcp"
)

// serverKeys lista as chaves onde cada formato de harness guarda seus MCPs.
func serverKeys(path string) []string {
	if isMimoOrOpencode(path) {
		return []string{"mcp"}
	}
	return []string{"mcpServers", "servers"}
}

// lazyExtract move os MCPs stdio eager de obj para specs e os remove de obj.
// Remotos (URL), desabilitados e entradas sem command ficam como estao.
func lazyExtract(path string, obj map[string]any) map[string]mcp.ServerSpec {
	specs := map[string]mcp.ServerSpec{}
	for _, key := range serverKeys(path) {
		section, ok := obj[key].(map[string]any)
		if !ok {
			continue
		}
		for name, raw := range section {
			if name == "legoloop" || entryDisabled(raw) {
				continue
			}
			if spec, ok := toServerSpec(raw); ok {
				specs[name] = spec
				delete(section, name)
			}
		}
	}
	return specs
}

func entryDisabled(raw any) bool {
	m, ok := raw.(map[string]any)
	if !ok {
		return false
	}
	return m["disabled"] == true || m["enabled"] == false
}

// toServerSpec converte uma entrada de harness (formato padrao ou OpenCode/Mimo) em ServerSpec.
func toServerSpec(raw any) (mcp.ServerSpec, bool) {
	m, ok := raw.(map[string]any)
	if !ok {
		return mcp.ServerSpec{}, false
	}
	switch cmd := m["command"].(type) {
	case string:
		return mcp.ServerSpec{Command: cmd, Args: toStrings(m["args"]), Env: toEnv(m["env"])}, cmd != ""
	case []any:
		parts := toStrings(cmd)
		if len(parts) == 0 {
			return mcp.ServerSpec{}, false
		}
		return mcp.ServerSpec{Command: parts[0], Args: parts[1:], Env: toEnv(m["environment"])}, true
	}
	return mcp.ServerSpec{}, false
}

func toStrings(v any) []string {
	items, _ := v.([]any)
	out := make([]string, 0, len(items))
	for _, it := range items {
		if s, ok := it.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func toEnv(v any) map[string]string {
	m, _ := v.(map[string]any)
	if len(m) == 0 {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, val := range m {
		if s, ok := val.(string); ok {
			out[k] = s
		}
	}
	return out
}

// persistSpecs acrescenta os servidores ao registro canonico (add-only).
func persistSpecs(home string, specs map[string]mcp.ServerSpec) error {
	if len(specs) == 0 {
		return nil
	}
	path := mcp.RegistryPath(home)
	current, err := mcp.LoadRegistry(path)
	if err != nil {
		return err
	}
	return mcp.SaveRegistry(path, mcp.MergeRegistry(current, specs))
}

func specNames(specs map[string]mcp.ServerSpec) []string {
	names := make([]string, 0, len(specs))
	for n := range specs {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

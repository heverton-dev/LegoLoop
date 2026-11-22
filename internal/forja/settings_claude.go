package forja

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func ferramentasPesadasBase() []string {
	return []string{
		"Task",
		"Agent",
		"PowerShell",
		"Workflow",
		"ScheduleWakeup",
		"ReportFindings",
		"ShareOnboardingGuide",
		"ListAgents",
	}
}

// PlanejarSettingsClaude mescla o JSON de settings.json preservando chaves desconhecidas.
func PlanejarSettingsClaude(atual []byte, perfil string) ([]byte, error) {
	denyAlvo, err := ferramentasNegadasPorPerfil(perfil)
	if err != nil {
		return nil, err
	}
	doc, err := decodificarObjetoJSON(atual)
	if err != nil {
		return nil, err
	}
	doc["syncClaudeAiSkills"] = false
	doc["skillListingBudgetFraction"] = 0.001
	aplicarEnvPadrao(doc)
	aplicarPermissionsDeny(doc, denyAlvo)
	saida, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(saida, '\n'), nil
}

func ferramentasNegadasPorPerfil(perfil string) ([]string, error) {
	base := ferramentasPesadasBase()
	switch perfil {
	case "padrao":
		return base, nil
	case "minimo":
		return append(base, "Grep", "Glob"), nil
	case "gateway":
		return append(base, "Bash", "Read", "Edit", "Write", "Grep", "Glob", "Skill"), nil
	default:
		return nil, fmt.Errorf("forja: perfil desconhecido: %s", perfil)
	}
}

func decodificarObjetoJSON(dados []byte) (map[string]any, error) {
	if len(bytes.TrimSpace(dados)) == 0 {
		return map[string]any{}, nil
	}
	var doc map[string]any
	if err := json.Unmarshal(dados, &doc); err != nil {
		return nil, err
	}
	if doc == nil {
		doc = map[string]any{}
	}
	return doc, nil
}

func aplicarEnvPadrao(doc map[string]any) {
	env, ok := doc["env"].(map[string]any)
	if !ok || env == nil {
		env = map[string]any{}
	}
	env["ENABLE_CLAUDEAI_MCP_SERVERS"] = "false"
	doc["env"] = env
}

func aplicarPermissionsDeny(doc map[string]any, extras []string) {
	perms, ok := doc["permissions"].(map[string]any)
	if !ok || perms == nil {
		perms = map[string]any{}
	}
	perms["deny"] = mesclarListaUnica(perms["deny"], extras)
	doc["permissions"] = perms
}

func mesclarListaUnica(atual any, extras []string) []string {
	vistos := map[string]bool{}
	var res []string
	if bruto, ok := atual.([]any); ok {
		for _, item := range bruto {
			if texto, ok := item.(string); ok && !vistos[texto] {
				vistos[texto] = true
				res = append(res, texto)
			}
		}
	}
	for _, nome := range extras {
		if !vistos[nome] {
			vistos[nome] = true
			res = append(res, nome)
		}
	}
	return res
}

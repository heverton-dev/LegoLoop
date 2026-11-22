package forja

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestSettingsClaude(t *testing.T) {
	entrada := []byte(`{
  "x_usuario": {"ativo": true, "nota": "preservar"},
  "env": {"MINHA_VAR": "ok"},
  "permissions": {
    "allow": ["Bash(git status)"],
    "deny": ["CustomTool"]
  }
}`)
	saida, err := PlanejarSettingsClaude(entrada, "padrao")
	if err != nil {
		t.Fatalf("PlanejarSettingsClaude falhou: %v", err)
	}

	var doc map[string]any
	if err := json.Unmarshal(saida, &doc); err != nil {
		t.Fatalf("JSON invalido na saida: %v", err)
	}
	xUsr, ok := doc["x_usuario"].(map[string]any)
	if !ok || xUsr["nota"] != "preservar" {
		t.Fatalf("chave x_usuario nao foi preservada: %#v", doc["x_usuario"])
	}
	env, ok := doc["env"].(map[string]any)
	if !ok || env["MINHA_VAR"] != "ok" || env["ENABLE_CLAUDEAI_MCP_SERVERS"] != "false" {
		t.Fatalf("env incorreto: %#v", doc["env"])
	}

	perms, ok := doc["permissions"].(map[string]any)
	if !ok {
		t.Fatalf("permissions ausente: %#v", doc)
	}
	allow, _ := perms["allow"].([]any)
	if len(allow) != 1 || allow[0] != "Bash(git status)" {
		t.Fatalf("permissions.allow nao foi preservado: %#v", perms["allow"])
	}

	denySet := map[string]bool{}
	for _, item := range perms["deny"].([]any) {
		denySet[item.(string)] = true
	}
	esperadas := []string{
		"CustomTool",
		"Task", "Agent", "PowerShell", "Workflow",
		"ScheduleWakeup", "ReportFindings", "ShareOnboardingGuide", "ListAgents",
	}
	for _, nome := range esperadas {
		if !denySet[nome] {
			t.Errorf("ferramenta %q ausente em permissions.deny: %v", nome, perms["deny"])
		}
	}

	saida2, err := PlanejarSettingsClaude(saida, "padrao")
	if err != nil {
		t.Fatalf("segunda chamada falhou: %v", err)
	}
	if !bytes.Equal(saida, saida2) {
		t.Errorf("PlanejarSettingsClaude nao e idempotente:\n1: %s\n2: %s", saida, saida2)
	}
}

func TestSettingsClaudePerfisEErros(t *testing.T) {
	base := []byte(`{"x_usuario": 42}`)
	for _, perfil := range []string{"minimo", "gateway"} {
		out, err := PlanejarSettingsClaude(base, perfil)
		if err != nil {
			t.Fatalf("perfil %s falhou: %v", perfil, err)
		}
		var doc map[string]any
		if err := json.Unmarshal(out, &doc); err != nil {
			t.Fatalf("JSON invalido para %s: %v", perfil, err)
		}
		if doc["x_usuario"] != float64(42) {
			t.Errorf("x_usuario perdido no perfil %s", perfil)
		}
		if doc["syncClaudeAiSkills"] != false {
			t.Errorf("syncClaudeAiSkills esperado false no perfil %s", perfil)
		}
	}
	if _, err := PlanejarSettingsClaude(nil, "padrao"); err != nil {
		t.Errorf("entrada vazia deveria ser aceita: %v", err)
	}
	if _, err := PlanejarSettingsClaude(base, "xyz"); err == nil {
		t.Error("perfil desconhecido deveria retornar erro")
	}
	if _, err := PlanejarSettingsClaude([]byte("{invalido"), "padrao"); err == nil {
		t.Error("JSON invalido deveria retornar erro")
	}
}

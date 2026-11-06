package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSyncHarnesses(t *testing.T) {
	tmpHome := t.TempDir()

	// 1. Cria fixture de ~/.claude.json
	claudePath := filepath.Join(tmpHome, ".claude.json")
	claudeData := `{"theme":"dark","mcpServers":{"old_tool":{"command":"node"}}}`
	if err := os.WriteFile(claudePath, []byte(claudeData), 0644); err != nil {
		t.Fatalf("erro ao criar fixture claude: %v", err)
	}

	// 2. Cria fixture de ~/.config/mimocode/mimocode.jsonc
	mimoDir := filepath.Join(tmpHome, ".config", "mimocode")
	if err := os.MkdirAll(mimoDir, 0755); err != nil {
		t.Fatalf("erro ao criar dir mimocode: %v", err)
	}
	mimoPath := filepath.Join(mimoDir, "mimocode.jsonc")
	mimoData := `{"provider":{"custom":{}}}`
	if err := os.WriteFile(mimoPath, []byte(mimoData), 0644); err != nil {
		t.Fatalf("erro ao criar fixture mimocode: %v", err)
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	deps := Deps{
		Stdin:  strings.NewReader(""),
		Stdout: stdout,
		Stderr: stderr,
		Home:   tmpHome,
	}

	ctx := context.Background()

	// 3. Testa dry-run: nao deve modificar arquivos fisicos
	exitDry := Run(ctx, []string{"sync-harnesses", "--dry-run"}, deps)
	if exitDry != 0 {
		t.Fatalf("dry-run esperava exit 0, obteve %d: %s", exitDry, stderr.String())
	}
	if _, err := os.Stat(claudePath + ".bak_legoloop"); !os.IsNotExist(err) {
		t.Errorf("dry-run nao deveria criar arquivo de backup")
	}

	// 4. Testa apply: deve criar backup e aplicar configuracao canônica
	stdout.Reset()
	stderr.Reset()
	exitApply := Run(ctx, []string{"sync-harnesses", "--apply"}, deps)
	if exitApply != 0 {
		t.Fatalf("apply esperava exit 0, obteve %d: %s", exitApply, stderr.String())
	}

	// Verifica se gerou backup
	if _, err := os.Stat(claudePath + ".bak_legoloop"); os.IsNotExist(err) {
		t.Errorf("apply deveria ter criado claude.json.bak_legoloop")
	}

	// Verifica se a config do claude agora contem o servidor legoloop mcp-router
	updatedClaude, err := os.ReadFile(claudePath)
	if err != nil {
		t.Fatalf("erro ao ler claude atualizado: %v", err)
	}
	var claudeObj map[string]any
	if err := json.Unmarshal(updatedClaude, &claudeObj); err != nil {
		t.Fatalf("erro ao decodificar claude atualizado: %v", err)
	}
	if claudeObj["theme"] != "dark" {
		t.Errorf("chave existente 'theme' foi corrompida: %v", claudeObj["theme"])
	}
	mcpServers, ok := claudeObj["mcpServers"].(map[string]any)
	if !ok || mcpServers["legoloop"] == nil {
		t.Errorf("esperava mcpServers.legoloop no claude.json, obteve: %+v", claudeObj["mcpServers"])
	}
}

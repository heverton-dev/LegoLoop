package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func parseSyncHarnessesFlags(args []string) (dryRun, apply bool, err error) {
	fs := flag.NewFlagSet("sync-harnesses", flag.ContinueOnError)
	fs.BoolVar(&dryRun, "dry-run", false, "Simula alterações sem gravar")
	fs.BoolVar(&apply, "apply", false, "Aplica as configurações e cria backups")
	err = fs.Parse(args)
	return dryRun, apply, err
}

func targetConfigs(home string) []string {
	return []string{
		filepath.Join(home, ".claude.json"),
		filepath.Join(home, ".config", "mimocode", "mimocode.jsonc"),
		filepath.Join(home, ".cursor", "mcp.json"),
		filepath.Join(home, ".codeium", "windsurf", "mcp_config.json"),
	}
}

func cmdSyncHarnesses(_ context.Context, args []string, d Deps) int {
	dryRun, apply, err := parseSyncHarnessesFlags(args)
	if err != nil || (!dryRun && !apply) {
		fmt.Fprintln(d.Stderr, "Uso: legoloop sync-harnesses (--dry-run|--apply)")
		return 2
	}
	home := resolveHome(d, "")
	for _, target := range targetConfigs(home) {
		if _, err := os.Stat(target); os.IsNotExist(err) {
			continue
		}
		if err := syncFile(target, dryRun, d.Stdout); err != nil {
			fmt.Fprintf(d.Stderr, "Erro ao sincronizar %s: %v\n", target, err)
			return 1
		}
	}
	return 0
}

func backupFile(path string) error {
	bak := path + ".bak_legoloop"
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return os.WriteFile(bak, data, 0644)
}

func injectServerConfig(raw []byte) ([]byte, error) {
	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		return nil, err
	}
	var servers map[string]any
	if existing, ok := obj["mcpServers"].(map[string]any); ok && existing != nil {
		servers = existing
	} else {
		servers = make(map[string]any)
	}
	servers["legoloop"] = map[string]any{
		"command": "legoloop",
		"args":    []string{"mcp-router"},
	}
	obj["mcpServers"] = servers
	return json.MarshalIndent(obj, "", "  ")
}

func syncFile(path string, dryRun bool, out io.Writer) error {
	if dryRun {
		fmt.Fprintf(out, "[dry-run] Sincronizaria %s com legoloop mcp-router\n", path)
		return nil
	}
	if err := backupFile(path); err != nil {
		return err
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	updated, err := injectServerConfig(raw)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, updated, 0644); err != nil {
		return err
	}
	fmt.Fprintf(out, "Sincronizado %s com sucesso (backup criado)\n", path)
	return nil
}

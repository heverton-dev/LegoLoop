package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"legoloop/internal/mcp"
)

type syncTargetResult struct {
	Path   string   `json:"path"`
	Status string   `json:"status"`
	Backup string   `json:"backup,omitempty"`
	Moved  []string `json:"moved_to_registry,omitempty"`
}

type syncHarnessesReport struct {
	DryRun  bool               `json:"dry_run"`
	Targets []syncTargetResult `json:"targets"`
}

func parseSyncHarnessesFlags(args []string, errOut io.Writer) (dryRun, apply, jsonOut bool, extra []string, err error) {
	fs := flag.NewFlagSet("sync-harnesses", flag.ContinueOnError)
	if errOut != nil {
		fs.SetOutput(errOut)
	}
	fs.BoolVar(&dryRun, "dry-run", false, "Simula alterações sem gravar")
	fs.BoolVar(&apply, "apply", false, "Aplica as configurações e cria backups")
	fs.BoolVar(&jsonOut, "json", false, "Emite saída estruturada em JSON")
	cfg := fs.String("config", "", "Arquivos de config extras (separados por vírgula), ex.: perfil Claude do Orca")
	err = fs.Parse(args)
	if *cfg != "" {
		extra = strings.Split(*cfg, ",")
	}
	return dryRun, apply, jsonOut, extra, err
}

func targetConfigs(home string) []string {
	return []string{
		filepath.Join(home, ".claude.json"),
		filepath.Join(home, ".config", "mimocode", "mimocode.jsonc"),
		filepath.Join(home, ".config", "opencode", "opencode.json"),
		filepath.Join(home, ".cursor", "mcp.json"),
		filepath.Join(home, ".codeium", "windsurf", "mcp_config.json"),
		filepath.Join(home, ".gemini", "config", "mcp_config.json"),
		filepath.Join(home, ".gemini", "settings.json"),
		filepath.Join(home, ".gemini", "antigravity", "mcp_config.json"),
		filepath.Join(home, ".gemini", "antigravity-ide", "mcp_config.json"),
		filepath.Join(home, ".copilot", "mcp-config.json"),
		filepath.Join(home, "AppData", "Roaming", "Code", "User", "mcp.json"),
		filepath.Join(home, "AppData", "Roaming", "Trae", "User", "mcp.json"),
		filepath.Join(home, ".codex", "config.toml"),
		filepath.Join(home, ".pi", "agent", "mcp.json"),
		filepath.Join(home, ".omp", "agent", "mcp.json"),
		filepath.Join(home, "AppData", "Local", "hermes", "config.yaml"),
	}
}

func cmdSyncHarnesses(_ context.Context, args []string, d Deps) int {
	dryRun, apply, jsonOut, extra, err := parseSyncHarnessesFlags(args, d.Stderr)
	if err != nil || (!dryRun && !apply) {
		fmt.Fprintln(d.Stderr, "Uso: legoloop sync-harnesses (--dry-run|--apply) [--json]")
		return 2
	}
	home := resolveHome(d, "")
	report := syncHarnessesReport{DryRun: dryRun, Targets: []syncTargetResult{}}
	for _, target := range append(targetConfigs(home), extra...) {
		if _, err := os.Stat(target); os.IsNotExist(err) {
			continue
		}
		res, err := syncFileResult(home, target, dryRun)
		if err != nil {
			fmt.Fprintf(d.Stderr, "Erro ao sincronizar %s: %v\n", target, err)
			return 1
		}
		report.Targets = append(report.Targets, res)
		if !jsonOut {
			writeSyncLine(d.Stdout, res)
		}
	}
	if jsonOut {
		return printJSON(d, report)
	}
	return 0
}

func backupFile(path string) error {
	bak := path + ".bak_legoloop"
	if _, err := os.Stat(bak); err == nil {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return os.WriteFile(bak, data, 0644)
}

func injectMimoConfig(obj map[string]any) {
	delete(obj, "mcpServers")
	var mcpMap map[string]any
	if existing, ok := obj["mcp"].(map[string]any); ok && existing != nil {
		mcpMap = existing
	} else {
		mcpMap = make(map[string]any)
	}
	mcpMap["legoloop"] = map[string]any{
		"type":    "local",
		"command": []string{"legoloop", "mcp-router"},
		"enabled": true,
	}
	obj["mcp"] = mcpMap
}

func injectStandardConfig(obj map[string]any) {
	key, entry := "mcpServers", map[string]any{"command": "legoloop", "args": []string{"mcp-router"}}
	if _, hasStd := obj["mcpServers"].(map[string]any); !hasStd {
		if _, hasVS := obj["servers"].(map[string]any); hasVS {
			key, entry["type"] = "servers", "stdio"
		}
	}
	servers, _ := obj[key].(map[string]any)
	if servers == nil {
		servers = make(map[string]any)
	}
	servers["legoloop"] = entry
	obj[key] = servers
}

func isMimoOrOpencode(path string) bool {
	lower := strings.ToLower(path)
	return strings.Contains(lower, "mimocode") || strings.Contains(lower, "opencode")
}

func injectServerConfig(path string, raw []byte) ([]byte, map[string]mcp.ServerSpec, error) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".toml":
		out, specs := injectTOML(raw)
		return out, specs, nil
	case ".yaml", ".yml":
		out, specs := injectYAML(raw)
		return out, specs, nil
	}
	var obj map[string]any
	if err := json.Unmarshal(stripJSONC(raw), &obj); err != nil {
		return nil, nil, err
	}
	specs := lazyExtract(path, obj)
	if isMimoOrOpencode(path) {
		injectMimoConfig(obj)
	} else {
		injectStandardConfig(obj)
	}
	encoded, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return nil, nil, err
	}
	if header := extractJSONCHeader(raw); header != "" {
		return []byte(header + string(encoded) + "\n"), specs, nil
	}
	return encoded, specs, nil
}

func extractJSONCHeader(raw []byte) string {
	var lines []string
	for _, line := range strings.Split(string(raw), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") || strings.HasPrefix(trimmed, "*") {
			lines = append(lines, strings.TrimRight(line, "\r"))
			continue
		}
		if trimmed == "" {
			continue
		}
		break
	}
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n") + "\n"
}

func stripJSONC(raw []byte) []byte {
	noComments := stripJSONCComments(raw)
	return stripTrailingCommas(noComments)
}

func stripJSONCComments(src []byte) []byte {
	out := make([]byte, 0, len(src))
	inStr, esc := false, false
	for i := 0; i < len(src); i++ {
		c := src[i]
		if inStr {
			out = append(out, c)
			inStr, esc = stepJSONString(c, esc)
			continue
		}
		if c == '"' {
			inStr = true
			out = append(out, c)
			continue
		}
		if next, skipped := skipComment(src, i); skipped {
			i = next
			continue
		}
		out = append(out, c)
	}
	return out
}

func stepJSONString(c byte, esc bool) (bool, bool) {
	if esc {
		return true, false
	}
	if c == '\\' {
		return true, true
	}
	return c != '"', false
}

func skipComment(src []byte, i int) (int, bool) {
	if i+1 >= len(src) || src[i] != '/' {
		return i, false
	}
	if src[i+1] == '/' {
		for i < len(src) && src[i] != '\n' {
			i++
		}
		return i - 1, true
	}
	if src[i+1] == '*' {
		i += 2
		for i+1 < len(src) && !(src[i] == '*' && src[i+1] == '/') {
			i++
		}
		return i + 1, true
	}
	return i, false
}

func stripTrailingCommas(src []byte) []byte {
	out := make([]byte, 0, len(src))
	inStr, esc := false, false
	for i := 0; i < len(src); i++ {
		c := src[i]
		if inStr {
			out = append(out, c)
			inStr, esc = stepJSONString(c, esc)
			continue
		}
		if c == '"' {
			inStr = true
			out = append(out, c)
			continue
		}
		if c == ',' && isClosingNext(src, i+1) {
			continue
		}
		out = append(out, c)
	}
	return out
}

func isClosingNext(src []byte, idx int) bool {
	for j := idx; j < len(src); j++ {
		b := src[j]
		if b == ' ' || b == '\t' || b == '\r' || b == '\n' {
			continue
		}
		return b == '}' || b == ']'
	}
	return false
}

func writeSyncLine(out io.Writer, res syncTargetResult) {
	if res.Status == "dry-run" {
		fmt.Fprintf(out, "[dry-run] Sincronizaria %s com legoloop mcp-router\n", res.Path)
		return
	}
	fmt.Fprintf(out, "Sincronizado %s com sucesso (backup criado)\n", res.Path)
}

func syncFileResult(home, path string, dryRun bool) (syncTargetResult, error) {
	if dryRun {
		return syncTargetResult{Path: path, Status: "dry-run"}, nil
	}
	if err := backupFile(path); err != nil {
		return syncTargetResult{}, err
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return syncTargetResult{}, err
	}
	updated, specs, err := injectServerConfig(path, raw)
	if err != nil {
		return syncTargetResult{}, err
	}
	if err := persistSpecs(home, specs); err != nil {
		return syncTargetResult{}, err
	}
	if err := os.WriteFile(path, updated, 0644); err != nil {
		return syncTargetResult{}, err
	}
	return syncTargetResult{Path: path, Status: "synced", Backup: path + ".bak_legoloop", Moved: specNames(specs)}, nil
}

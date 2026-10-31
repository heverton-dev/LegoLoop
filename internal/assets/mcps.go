package assets

import (
	"context"
	"encoding/json"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// maxConfigBytes limita a leitura de arquivos de configuração (o ~/.claude.json pode ser grande).
const maxConfigBytes = 16 << 20

type mcpSource struct {
	path, origin, scope, key string
	nested                   bool // lê também projects.<caminho>.mcpServers (formato do ~/.claude.json)
}

func (s *Scanner) mcpSources(ctx context.Context) []mcpSource {
	h := s.cfg.Home
	src := []mcpSource{
		{path: filepath.Join(h, ".claude.json"), origin: "claude", scope: "global", key: "mcpServers", nested: true},
		{path: filepath.Join(h, ".claude", "settings.json"), origin: "claude", scope: "global", key: "mcpServers"},
		{path: filepath.Join(h, ".cursor", "mcp.json"), origin: "cursor", scope: "global", key: "mcpServers"},
		{path: filepath.Join(h, ".gemini", "settings.json"), origin: "gemini", scope: "global", key: "mcpServers"},
		{path: filepath.Join(h, ".codeium", "windsurf", "mcp_config.json"), origin: "windsurf", scope: "global", key: "mcpServers"},
		{path: filepath.Join(h, ".config", "opencode", "opencode.json"), origin: "opencode", scope: "global", key: "mcp"},
		{path: filepath.Join(h, ".config", "Claude", "claude_desktop_config.json"), origin: "claude-desktop", scope: "global", key: "mcpServers"},
	}
	if appdata := s.cfg.Getenv("APPDATA"); appdata != "" {
		src = append(src, mcpSource{path: filepath.Join(appdata, "Claude", "claude_desktop_config.json"),
			origin: "claude-desktop", scope: "global", key: "mcpServers"})
	}
	for _, p := range s.allProjects(ctx) {
		src = append(src,
			mcpSource{path: filepath.Join(p, ".mcp.json"), origin: "claude", scope: "project", key: "mcpServers"},
			mcpSource{path: filepath.Join(p, ".cursor", "mcp.json"), origin: "cursor", scope: "project", key: "mcpServers"},
			mcpSource{path: filepath.Join(p, "opencode.json"), origin: "opencode", scope: "project", key: "mcp"},
		)
	}
	return src
}

func (s *Scanner) scanMCPs(ctx context.Context) []Item {
	var out []Item
	for _, src := range s.mcpSources(ctx) {
		if ctx.Err() != nil {
			return out
		}
		out = append(out, s.load(KindMCP, src.path, func(p string) []Item { return parseMCPs(p, src) })...)
	}
	return out
}

// parseMCPs extrai só nome, transporte, comando-base e host. env, args, headers e query de
// URL nunca são copiados.
func parseMCPs(path string, src mcpSource) []Item {
	obj := readJSONObject(path)
	if obj == nil {
		return nil
	}
	seen := map[string]bool{}
	var out []Item
	add := func(servers any, scope string) {
		m, _ := servers.(map[string]any)
		for name, raw := range m {
			entry, ok := raw.(map[string]any)
			key := scope + "\x00" + name
			if !ok || seen[key] {
				continue
			}
			seen[key] = true
			out = append(out, Item{Kind: KindMCP, Name: name, Path: path, Scope: scope, Origin: src.origin, Meta: mcpMeta(entry)})
		}
	}
	add(obj[src.key], src.scope)
	if projects, ok := obj["projects"].(map[string]any); ok && src.nested {
		for _, p := range projects {
			if pm, ok := p.(map[string]any); ok {
				add(pm[src.key], "project")
			}
		}
	}
	return out
}

func mcpMeta(entry map[string]any) map[string]string {
	meta := map[string]string{}
	cmd := commandOf(entry["command"])
	rawURL, _ := entry["url"].(string)
	typ, _ := entry["type"].(string)
	switch {
	case rawURL != "":
		meta["transport"] = remoteTransport(typ, rawURL)
		if u, err := url.Parse(rawURL); err == nil && u.Hostname() != "" {
			meta["host"] = u.Hostname()
		}
	case cmd != "":
		meta["transport"] = "stdio"
	}
	if cmd != "" {
		meta["command"] = baseCommand(cmd)
	}
	return meta
}

func remoteTransport(typ, rawURL string) string {
	if typ == "sse" || strings.Contains(strings.ToLower(rawURL), "/sse") {
		return "sse"
	}
	return "http"
}

// commandOf aceita "command" como texto ou como array (opencode); só o primeiro elemento importa.
func commandOf(v any) string {
	switch c := v.(type) {
	case string:
		return c
	case []any:
		if len(c) > 0 {
			s, _ := c[0].(string)
			return s
		}
	}
	return ""
}

// baseCommand devolve só o nome do executável, sem diretório nem argumentos.
func baseCommand(cmd string) string {
	cmd = strings.Trim(strings.TrimSpace(cmd), `"'`)
	if i := strings.LastIndexAny(cmd, `/\`); i >= 0 {
		cmd = cmd[i+1:]
	}
	if f := strings.Fields(cmd); len(f) > 0 {
		return strings.Trim(f[0], `"'`)
	}
	return ""
}

// readJSONObject lê um objeto JSON de até maxConfigBytes; qualquer problema vira nil.
func readJSONObject(path string) map[string]any {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, maxConfigBytes+1))
	if err != nil || len(data) > maxConfigBytes {
		return nil
	}
	var obj map[string]any
	if json.Unmarshal(data, &obj) != nil {
		return nil
	}
	return obj
}

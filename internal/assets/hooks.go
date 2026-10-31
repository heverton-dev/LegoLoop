package assets

import (
	"context"
	"path/filepath"
	"regexp"
	"unicode/utf8"
)

// maxHookDesc limita, em runas, o comando guardado no inventário (proteção de memória). A exibição
// é mais curta e quem expõe o texto mascara caminhos ANTES de truncar, para nunca mostrar um
// caminho cortado ao meio.
const maxHookDesc = 1000

var (
	secretAssignRe = regexp.MustCompile(`(?i)\b([A-Za-z0-9_]*(?:token|key|secret|password|passwd|pwd|authorization))(\s*[=:]\s*)\S+`)
	bearerRe       = regexp.MustCompile(`(?i)\b(bearer)\s+\S+`)
)

// redact mascara valores que parecem segredo (token=..., API_KEY=..., Authorization: Bearer ...).
func redact(cmd string) string {
	// Bearer primeiro: senão "Authorization: Bearer x" perderia só a palavra "Bearer" e vazaria o token.
	cmd = bearerRe.ReplaceAllString(cmd, "${1} ***")
	return secretAssignRe.ReplaceAllString(cmd, "${1}${2}***")
}

// ContainsSecret diz se o comando traz um valor que parece credencial (as mesmas regras do redact).
func ContainsSecret(cmd string) bool { return redact(cmd) != cmd }

func truncate(s string, n int) string {
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	return string([]rune(s)[:n]) + "…"
}

type hookSource struct{ path, scope string }

func (s *Scanner) hookSources(ctx context.Context) []hookSource {
	src := []hookSource{{filepath.Join(s.cfg.Home, ".claude", "settings.json"), "global"}}
	for _, p := range s.allProjects(ctx) {
		src = append(src,
			hookSource{filepath.Join(p, ".claude", "settings.json"), "project"},
			hookSource{filepath.Join(p, ".claude", "settings.local.json"), "local"},
		)
	}
	return append(src, s.scaffoldHookSources()...)
}

func (s *Scanner) scanHooks(ctx context.Context) []Item {
	var out []Item
	for _, src := range s.hookSources(ctx) {
		if ctx.Err() != nil {
			return out
		}
		out = append(out, s.load(KindHook, src.path, func(p string) []Item { return parseHooks(p, src.scope) })...)
	}
	return out
}

// parseHooks lê settings.json: hooks.<Evento> = [{matcher, hooks: [{type, command}]}].
// Entradas malformadas são ignoradas; o comando sai redigido e truncado.
func parseHooks(path, scope string) []Item {
	obj := readJSONObject(path)
	events, _ := obj["hooks"].(map[string]any)
	var out []Item
	for event, raw := range events {
		groups, _ := raw.([]any)
		for _, g := range groups {
			group, _ := g.(map[string]any)
			matcher, _ := group["matcher"].(string)
			list, _ := group["hooks"].([]any)
			for _, h := range list {
				hook, _ := h.(map[string]any)
				typ, _ := hook["type"].(string)
				cmd, _ := hook["command"].(string)
				if typ != "command" || cmd == "" {
					continue
				}
				out = append(out, hookItem(path, scope, event, matcher, typ, cmd))
			}
		}
	}
	return out
}

func hookItem(path, scope, event, matcher, typ, cmd string) Item {
	name := event
	meta := map[string]string{"event": event, "type": typ}
	if matcher != "" {
		name = event + ":" + matcher
		meta["matcher"] = matcher
	}
	raw := cmd
	if len(raw) > maxCommandLen {
		raw = "" // comando longo demais para virar argv com fidelidade; nunca guardar truncado
	}
	return Item{Kind: KindHook, Name: name, Description: truncate(redact(cmd), maxHookDesc),
		Path: path, Scope: scope, Origin: "claude", Meta: meta, Command: raw}
}

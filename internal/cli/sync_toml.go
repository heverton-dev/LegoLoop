package cli

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/mcp"
)

const legoloopTOML = "[mcp_servers.legoloop]\ncommand = \"legoloop\"\nargs = [\"mcp-router\"]\n"

var (
	tomlHeader = regexp.MustCompile(`^\s*\[([^\[\]]+)\]\s*$`)
	tomlStrLit = `"(?:[^"\\]|\\.)*"|'[^']*'`
	tomlStrRe  = regexp.MustCompile(tomlStrLit)
	tomlEnvRe  = regexp.MustCompile(`(?m)^\s*([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(` + tomlStrLit + `)`)
	tomlArgsRe = regexp.MustCompile(`(?ms)^\s*args\s*=\s*\[(.*?)\]`)
	tomlOffRe  = regexp.MustCompile(`(?m)^\s*enabled\s*=\s*false`)
)

type textSeg struct {
	header string
	lines  []string
}

// splitTOML agrupa as linhas por tabela; unir os grupos devolve o texto original.
func splitTOML(raw string) []textSeg {
	segs := []textSeg{{}}
	for _, line := range strings.Split(raw, "\n") {
		if m := tomlHeader.FindStringSubmatch(strings.TrimRight(line, "\r")); m != nil {
			segs = append(segs, textSeg{header: strings.TrimSpace(m[1])})
		}
		last := &segs[len(segs)-1]
		last.lines = append(last.lines, line)
	}
	return segs
}

func tomlServerOf(header string) (name string, isEnv, ok bool) {
	rest, found := strings.CutPrefix(header, "mcp_servers.")
	if !found {
		return "", false, false
	}
	if base, isE := strings.CutSuffix(rest, ".env"); isE {
		return strings.Trim(base, `"`), true, true
	}
	return strings.Trim(rest, `"`), false, true
}

func tomlUnquote(s string) string {
	if strings.HasPrefix(s, "'") {
		return strings.Trim(s, "'")
	}
	if u, err := strconv.Unquote(s); err == nil {
		return u
	}
	return strings.Trim(s, `"`)
}

func tomlValue(txt, key string) (string, bool) {
	re := regexp.MustCompile(`(?m)^\s*` + key + `\s*=\s*(` + tomlStrLit + `)`)
	m := re.FindStringSubmatch(txt)
	if m == nil {
		return "", false
	}
	return tomlUnquote(m[1]), true
}

func tomlSpec(base, env string) (mcp.ServerSpec, bool) {
	cmd, ok := tomlValue(base, "command")
	if _, remote := tomlValue(base, "url"); !ok || remote || tomlOffRe.MatchString(base) {
		return mcp.ServerSpec{}, false
	}
	spec := mcp.ServerSpec{Command: cmd}
	spec.Cwd, _ = tomlValue(base, "cwd")
	if m := tomlArgsRe.FindStringSubmatch(base); m != nil {
		for _, it := range tomlStrRe.FindAllString(m[1], -1) {
			spec.Args = append(spec.Args, tomlUnquote(it))
		}
	}
	for _, m := range tomlEnvRe.FindAllStringSubmatch(env, -1) {
		if spec.Env == nil {
			spec.Env = map[string]string{}
		}
		spec.Env[m[1]] = tomlUnquote(m[2])
	}
	return spec, true
}

func tomlSpecs(segs []textSeg) map[string]mcp.ServerSpec {
	base, env := map[string]string{}, map[string]string{}
	for _, s := range segs {
		name, isEnv, ok := tomlServerOf(s.header)
		if !ok {
			continue
		}
		if isEnv {
			env[name] = strings.Join(s.lines, "\n")
		} else {
			base[name] = strings.Join(s.lines, "\n")
		}
	}
	specs := map[string]mcp.ServerSpec{}
	for name, txt := range base {
		if spec, ok := tomlSpec(txt, env[name]); ok && name != "legoloop" {
			specs[name] = spec
		}
	}
	return specs
}

// injectTOML move MCPs stdio para o registro e garante a tabela legoloop, preservando o resto.
func injectTOML(raw []byte) ([]byte, map[string]mcp.ServerSpec) {
	segs := splitTOML(string(raw))
	specs := tomlSpecs(segs)
	hasLego := false
	var out []string
	for _, s := range segs {
		name, _, ok := tomlServerOf(s.header)
		hasLego = hasLego || (ok && name == "legoloop")
		if _, moved := specs[name]; ok && moved {
			continue
		}
		out = append(out, s.lines...)
	}
	text := strings.Join(out, "\n")
	if !hasLego {
		text = strings.TrimRight(text, "\r\n") + "\n\n" + legoloopTOML
	}
	return []byte(text), specs
}

package cli

import (
	"regexp"
	"strings"

	"legoloop/internal/mcp"
)

var legoloopYAML = []string{"  legoloop:", "    command: legoloop", "    args:", "      - mcp-router"}

var (
	yamlKeyRe  = regexp.MustCompile(`^ {4}([A-Za-z_]+):\s*(.*)$`)
	yamlItemRe = regexp.MustCompile(`^ {6}-\s*(.*)$`)
	yamlEnvRe  = regexp.MustCompile(`^ {6}([A-Za-z_][A-Za-z0-9_]*):\s*(.*)$`)
	yamlKidRe  = regexp.MustCompile(`^ {2}([^\s#:][^:]*):\s*$`)
)

type yamlChild struct {
	name  string
	lines []string
}

func yamlUnquote(s string) string {
	s = strings.TrimSpace(strings.TrimRight(s, "\r"))
	if len(s) >= 2 && (s[0] == '"' || s[0] == '\'') && s[len(s)-1] == s[0] {
		return s[1 : len(s)-1]
	}
	return s
}

func yamlBlockBounds(lines []string) (start, end int) {
	start = -1
	for i, l := range lines {
		if strings.TrimRight(l, "\r ") == "mcp_servers:" {
			start = i
			break
		}
	}
	if start < 0 {
		return -1, -1
	}
	end = len(lines)
	for i := start + 1; i < len(lines); i++ {
		if t := strings.TrimRight(lines[i], "\r "); t != "" && t[0] != ' ' && t[0] != '\t' {
			return start, i
		}
	}
	return start, end
}

func yamlChildren(block []string) []yamlChild {
	var kids []yamlChild
	for _, l := range block {
		if m := yamlKidRe.FindStringSubmatch(strings.TrimRight(l, "\r")); m != nil {
			kids = append(kids, yamlChild{name: strings.TrimSpace(m[1])})
		}
		if len(kids) > 0 {
			kids[len(kids)-1].lines = append(kids[len(kids)-1].lines, l)
		}
	}
	return kids
}

// yamlSpec le command/args/env/cwd de um filho; remoto, desabilitado ou args inline nao sao extraidos.
func yamlSpec(lines []string) (mcp.ServerSpec, bool) {
	var spec mcp.ServerSpec
	section, remote, off, inline := "", false, false, false
	for _, raw := range lines[1:] {
		l := strings.TrimRight(raw, "\r")
		if m := yamlKeyRe.FindStringSubmatch(l); m != nil {
			section = m[1]
			v := yamlUnquote(m[2])
			spec, remote, off, inline = yamlApplyKey(spec, m[1], v, remote, off, inline)
			continue
		}
		spec = yamlApplyItem(spec, section, l)
	}
	return spec, spec.Command != "" && !remote && !off && !inline
}

func yamlApplyKey(spec mcp.ServerSpec, key, v string, remote, off, inline bool) (mcp.ServerSpec, bool, bool, bool) {
	switch key {
	case "command":
		spec.Command = v
	case "cwd":
		spec.Cwd = v
	case "url":
		remote = true
	case "enabled":
		off = off || v == "false"
	case "args":
		inline = inline || v != ""
	}
	return spec, remote, off, inline
}

func yamlApplyItem(spec mcp.ServerSpec, section, l string) mcp.ServerSpec {
	if m := yamlItemRe.FindStringSubmatch(l); m != nil && section == "args" {
		spec.Args = append(spec.Args, yamlUnquote(m[1]))
	}
	if m := yamlEnvRe.FindStringSubmatch(l); m != nil && section == "env" {
		if spec.Env == nil {
			spec.Env = map[string]string{}
		}
		spec.Env[m[1]] = yamlUnquote(m[2])
	}
	return spec
}

// injectYAML move MCPs stdio do trecho mcp_servers para o registro e garante o legoloop.
func injectYAML(raw []byte) ([]byte, map[string]mcp.ServerSpec) {
	lines := strings.Split(string(raw), "\n")
	start, end := yamlBlockBounds(lines)
	if start < 0 {
		return raw, nil
	}
	block := lines[start+1 : end]
	trail := 0
	for trail < len(block) && strings.TrimSpace(block[len(block)-1-trail]) == "" {
		trail++
	}
	kept, specs := yamlKeep(yamlChildren(block[:len(block)-trail]))
	out := append([]string{}, lines[:start+1]...)
	out = append(out, kept...)
	out = append(out, block[len(block)-trail:]...)
	return []byte(strings.Join(append(out, lines[end:]...), "\n")), specs
}

func yamlKeep(kids []yamlChild) ([]string, map[string]mcp.ServerSpec) {
	specs := map[string]mcp.ServerSpec{}
	var kept []string
	hasLego := false
	for _, c := range kids {
		hasLego = hasLego || c.name == "legoloop"
		if spec, ok := yamlSpec(c.lines); ok && c.name != "legoloop" {
			specs[c.name] = spec
			continue
		}
		kept = append(kept, c.lines...)
	}
	if !hasLego {
		kept = append(kept, legoloopYAML...)
	}
	return kept, specs
}

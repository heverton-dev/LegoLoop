package assets

import (
	"context"
	"os"
	"path/filepath"

	"legoloop/internal/gate"
)

type harnessDef struct {
	name string
	bins []string
	dirs []string // relativos ao HOME
}

// harnessDefs: ambientes de agente conhecidos. Um harness conta se o binário está no PATH
// OU a pasta de configuração existe.
func harnessDefs() []harnessDef {
	return []harnessDef{
		{"claude", []string{"claude"}, []string{".claude"}},
		{"codex", []string{"codex"}, []string{".codex"}},
		{"gemini", []string{"gemini"}, []string{".gemini"}},
		{"opencode", []string{"opencode"}, []string{".config/opencode"}},
		{"cursor", []string{"cursor-agent", "cursor"}, []string{".cursor"}},
		{"agy", []string{"agy"}, nil},
		{"omp", []string{"omp"}, []string{".omp"}},
		{"mimo", []string{"mimo"}, []string{".mimocode"}},
		{"kiro", []string{"kiro-cli"}, []string{".kiro"}},
	}
}

// knownBinaries: ferramentas comuns de desenvolvimento que as peças podem exigir.
func knownBinaries() []string {
	return []string{
		"git", "go", "node", "npm", "pnpm", "bun", "python", "python3", "pip", "uv", "cargo", "rustc",
		"docker", "ollama", "rg", "jq", "gh", "curl", "ffmpeg", "pandoc", "make", "dotnet", "java",
	}
}

func yesNo(b bool) string {
	if b {
		return "sim"
	}
	return "não"
}

func (s *Scanner) scanHarnesses(ctx context.Context) []Item {
	var out []Item
	for _, h := range harnessDefs() {
		if ctx.Err() != nil {
			return out
		}
		binPath := ""
		for _, b := range h.bins {
			if p := s.findBinary(b); p != "" {
				binPath = p
				break
			}
		}
		cfgDir := ""
		for _, d := range h.dirs {
			if p := filepath.Join(s.cfg.Home, filepath.FromSlash(d)); isDir(p) {
				cfgDir = p
				break
			}
		}
		if binPath == "" && cfgDir == "" {
			continue
		}
		path := binPath
		if path == "" {
			path = cfgDir
		}
		out = append(out, Item{Kind: KindHarness, Name: h.name, Path: path, Scope: "system", Origin: "path",
			Meta: map[string]string{"cli": yesNo(binPath != ""), "config": yesNo(cfgDir != "")}})
	}
	return out
}

func (s *Scanner) scanBinaries(ctx context.Context) []Item {
	var out []Item
	for _, name := range knownBinaries() {
		if ctx.Err() != nil {
			return out
		}
		if p := s.findBinary(name); p != "" {
			out = append(out, Item{Kind: KindBinary, Name: name, Path: p, Scope: "system", Origin: "path"})
		}
	}
	return out
}

func (s *Scanner) findBinary(name string) string {
	if p, err := s.cfg.LookPath(name); err == nil && p != "" {
		return p
	}
	if s.cfg.Home == "" {
		return ""
	}
	dirs := []string{
		filepath.Join(s.cfg.Home, ".cargo", "bin"),
		filepath.Join(s.cfg.Home, "go", "bin"),
		filepath.Join(s.cfg.Home, ".local", "bin"),
		filepath.Join(s.cfg.Home, ".dotnet", "tools"),
	}
	for _, d := range dirs {
		for _, target := range []string{name, name + ".exe"} {
			p := filepath.Join(d, target)
			if st, err := os.Stat(p); err == nil && !st.IsDir() {
				return p
			}
		}
	}
	return ""
}

func (s *Scanner) scanGates(context.Context) []Item {
	var out []Item
	for _, d := range gate.Catalog() {
		out = append(out, Item{Kind: KindGate, Name: d.Name, Description: d.Description, Scope: "builtin",
			Origin: "legoloop", Meta: map[string]string{"signature": d.Signature}})
	}
	return out
}

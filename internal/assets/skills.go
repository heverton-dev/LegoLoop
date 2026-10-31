package assets

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

const (
	maxSkillDepth = 5        // nível máximo de diretórios abaixo da raiz
	headBytes     = 64 << 10 // trecho inicial usado para ler o frontmatter
	maxNameRunes  = 128
	maxDescRunes  = 300
)

func defaultGlobalSkillRoots() []root {
	return []root{
		{".claude/skills", "", "claude"}, {".aidd/canonical/skills", "", "aidd"}, {".aidd/skills", "", "aidd"},
		{".agents/skills", "", "agents"}, {".codex/skills", "", "codex"}, {".gemini/skills", "", "gemini"},
		{".gemini/antigravity-cli/skills", "", "gemini"}, {".gemini/antigravity/skills", "", "gemini"},
		{".gemini/antigravity-ide/skills", "", "gemini"}, {".gemini/config/skills", "", "gemini"},
		{".cline/skills", "", "cline"}, {".factory/skills", "", "factory"},
		{".config/opencode/skills", "", "opencode"}, {".config/opencode/skill", "", "opencode"},
	}
}

func (s *Scanner) skillRoots() []root {
	roots := s.joinRoots(
		defaultGlobalSkillRoots(),
		[]root{{".claude/skills", "", "claude"}, {".agents/skills", "", "agents"}, {"skills", "", "project"}},
	)
	return append(roots, s.scaffoldSkillRoots()...)
}

func (s *Scanner) scanSkills(ctx context.Context) []Item {
	var out []Item
	for _, r := range s.skillRoots() {
		for _, file := range findSkillFiles(ctx, r.dir, 0) {
			for _, it := range s.load(KindSkill, file, parseSkill) {
				it.Scope, it.Origin = r.scope, r.origin
				out = append(out, it)
			}
		}
	}
	return out
}

// findSkillFiles devolve cada SKILL.md abaixo de dir, descendo recursivamente até maxSkillDepth.
func findSkillFiles(ctx context.Context, dir string, depth int) []string {
	if ctx.Err() != nil || depth >= maxSkillDepth {
		return nil
	}
	var out []string
	if depth > 0 {
		if f := filepath.Join(dir, "SKILL.md"); isRegular(f) {
			out = append(out, f)
		}
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return out
	}
	for _, e := range entries {
		sub := filepath.Join(dir, e.Name())
		switch {
		case skipDirName(e.Name()):
		case e.IsDir():
			out = append(out, findSkillFiles(ctx, sub, depth+1)...)
		case e.Type()&os.ModeSymlink != 0 && isDir(sub):
			if f := filepath.Join(sub, "SKILL.md"); isRegular(f) {
				out = append(out, f)
			}
		}
	}
	return out
}

func isRegular(path string) bool {
	st, err := os.Stat(path)
	return err == nil && st.Mode().IsRegular()
}

func isDir(path string) bool {
	st, err := os.Stat(path)
	return err == nil && st.IsDir()
}

// parseSkill lê o arquivo inteiro uma vez: o SHA-256 cobre todo o conteúdo e o frontmatter
// vem só do trecho inicial.
func parseSkill(path string) []Item {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	h := sha256.New()
	head := &headBuffer{limit: headBytes}
	if _, err := io.Copy(io.MultiWriter(h, head), f); err != nil {
		return nil
	}
	name, desc := parseFrontmatter(head.buf.Bytes())
	if name == "" {
		name = filepath.Base(filepath.Dir(path))
	}
	return []Item{{
		Kind: KindSkill, Name: name, Description: desc, Path: path,
		Meta: map[string]string{"sha256": hex.EncodeToString(h.Sum(nil))},
	}}
}

// headBuffer guarda só os primeiros limit bytes escritos.
type headBuffer struct {
	buf   bytes.Buffer
	limit int
}

func (b *headBuffer) Write(p []byte) (int, error) {
	if room := b.limit - b.buf.Len(); room > 0 {
		b.buf.Write(p[:min(room, len(p))])
	}
	return len(p), nil
}

// parseFrontmatter extrai name e description do bloco YAML inicial (--- ... ---).
func parseFrontmatter(data []byte) (name, desc string) {
	text := strings.ReplaceAll(strings.TrimPrefix(string(data), "\xef\xbb\xbf"), "\r\n", "\n")
	lines := strings.Split(text, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return "", ""
	}
	for i := 1; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "---" {
			break
		}
		key, val, ok := strings.Cut(line, ":")
		if !ok || strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			continue
		}
		switch strings.TrimSpace(key) {
		case "name":
			name = limitRunes(unquote(val), maxNameRunes)
		case "description":
			desc, i = descriptionValue(lines, i, val)
		}
	}
	return name, desc
}

// descriptionValue lê a descrição simples ou em bloco (> | >- |-), devolvendo a última linha consumida.
func descriptionValue(lines []string, i int, val string) (string, int) {
	v := strings.TrimSpace(val)
	if v != ">" && v != "|" && v != ">-" && v != "|-" {
		return limitRunes(unquote(v), maxDescRunes), i
	}
	var parts []string
	for i+1 < len(lines) && (strings.HasPrefix(lines[i+1], " ") || strings.HasPrefix(lines[i+1], "\t") || lines[i+1] == "") {
		i++
		if t := strings.TrimSpace(lines[i]); t != "" {
			parts = append(parts, t)
		}
	}
	return limitRunes(strings.Join(parts, " "), maxDescRunes), i
}

func unquote(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && (s[0] == '"' && s[len(s)-1] == '"' || s[0] == '\'' && s[len(s)-1] == '\'') {
		return s[1 : len(s)-1]
	}
	return s
}

func limitRunes(s string, n int) string {
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	return string([]rune(s)[:n])
}

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
	name, desc, ont := parseFrontmatter(head.buf.Bytes())
	if name == "" {
		name = filepath.Base(filepath.Dir(path))
	}
	inferOntologyFallback(name, desc, ont)
	meta := map[string]string{"sha256": hex.EncodeToString(h.Sum(nil))}
	for k, v := range ont {
		if v != "" {
			meta[k] = v
		}
	}
	return []Item{{
		Kind: KindSkill, Name: name, Description: desc, Path: path,
		Meta: meta,
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

// parseFrontmatter extrai name, description e campos de contrato ontológico (--- ... ---).
func parseFrontmatter(data []byte) (name, desc string, ont map[string]string) {
	ont = make(map[string]string)
	text := strings.ReplaceAll(strings.TrimPrefix(string(data), "\xef\xbb\xbf"), "\r\n", "\n")
	lines := strings.Split(text, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return "", "", ont
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
		k := strings.TrimSpace(key)
		switch k {
		case "name":
			name = limitRunes(unquote(val), maxNameRunes)
		case "description":
			desc, i = descriptionValue(lines, i, val)
		case "consumes", "produces", "artifact", "gate", "suggested_model":
			v := strings.TrimSpace(val)
			v = strings.Trim(v, "[]")
			ont[k] = unquote(v)
		}
	}
	return name, desc, ont
}

type ontologyRule struct {
	keys                                 []string
	consumes, produces, artifact, gateEx string
}

// ontologyRules mapeia palavras do nome da skill ao contrato; a primeira regra que casa vence.
var ontologyRules = []ontologyRule{
	{[]string{"mega-pipeline", "megapipeline"}, "Artifact:Idea", "Artifact:Audit", "MEGA_REPORT.md", "file_exists('MEGA_REPORT.md')"},
	{[]string{"interview", "grill", "briefing"}, "Artifact:Idea", "Artifact:Briefing", "BRIEFING.md", "file_exists('BRIEFING.md')"},
	{[]string{"spec"}, "Artifact:Briefing", "Artifact:Spec", "SPEC.md", "file_exists('SPEC.md')"},
	{[]string{"ticket", "plan"}, "Artifact:Spec", "Artifact:Spec", "PLAN.md", "file_exists('PLAN.md')"},
	{[]string{"tdd", "implementacao"}, "Artifact:Spec", "Artifact:Code", "CODE.patch", "exit_code == 0"},
	{[]string{"homolog", "audit"}, "Artifact:Code", "Artifact:Audit", "REPORT.md", "file_exists('REPORT.md')"},
	{[]string{"rea", "obscura"}, "Artifact:URL", "Artifact:ReverseStructure", "ESTRUTURA_REVERSA.json", "file_exists('ESTRUTURA_REVERSA.json')"},
	{[]string{"impeccable"}, "Artifact:ReverseStructure", "Artifact:DesignTokens", "DESIGN.md", "file_exists('DESIGN.md')"},
}

func (r ontologyRule) matches(lower string) bool {
	for _, k := range r.keys {
		if strings.Contains(lower, k) {
			return true
		}
	}
	return false
}

// inferOntologyFallback preenche heurística ontológica para skills canônicas sem contrato explícito.
func inferOntologyFallback(name, desc string, ont map[string]string) {
	lower := strings.ToLower(name)
	for _, r := range ontologyRules {
		if !r.matches(lower) {
			continue
		}
		setDefault(ont, "consumes", r.consumes)
		setDefault(ont, "produces", r.produces)
		setDefault(ont, "artifact", r.artifact)
		setDefault(ont, "gate", r.gateEx)
		return
	}
}

func setDefault(m map[string]string, k, v string) {
	if m[k] == "" {
		m[k] = v
	}
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

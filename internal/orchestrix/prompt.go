package orchestrix

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"legoloop/internal/dag"
)

// maxInputBytes limita cada artefato de entrada injetado no prompt (contexto 2k–8k tokens).
const maxInputBytes = 256 << 10

// coreInviolableRules define o núcleo mínimo inviolável injetado em cada nó.
const coreInviolableRules = "- STE-PT: frases curtas (<20 palavras), voz ativa, PT-BR, fatos com prova.\n- Zero Stubs: proibido TODO, pass, mock simulado ou implementação incompleta.\n- Path Jail: proibido escape de diretório (.., symlinks externos ou caminhos fora da raiz)."

var fileInputRe = regexp.MustCompile(`^file\('([^']+)'\)$`)

// buildPrompt monta o prompt único e estéril: núcleo mínimo + skill + tooling + context_in + entradas.
func (o *Orchestrix) buildPrompt(n dag.Node, work *os.Root) (string, error) {
	skill, err := o.readSkill(n.Ref)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	fmt.Fprintf(&b, "# Orchestrix isolada — nó %s\n\n", n.ID)
	fmt.Fprintf(&b, "Você é um executor de tarefa única, sem histórico. Responda SOMENTE com o conteúdo final do artefato `%s`, sem comentários extras.\n\n", n.Output.Artifact)
	fmt.Fprintf(&b, "Condição de saída (gate): `%s`\n\n## Núcleo Mínimo Inviolável\n\n%s\n\n## Skill\n\n%s\n", n.Output.Gate, coreInviolableRules, strings.TrimSpace(skill))
	o.appendRuleTags(&b, n.RuleTags)
	o.appendTooling(&b, n.Tooling)
	o.appendContextIn(&b, n.ContextIn, work)
	if err := o.appendInputs(&b, n.Inputs, work); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (o *Orchestrix) appendInputs(b *strings.Builder, inputs map[string]string, work *os.Root) error {
	keys := make([]string, 0, len(inputs))
	for k := range inputs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v, err := o.resolveInput(inputs[k], work)
		if err != nil {
			return fmt.Errorf("entrada %q: %w", k, err)
		}
		fmt.Fprintf(b, "\n## Entrada: %s\n\n%s\n", k, v)
	}
	return nil
}

func (o *Orchestrix) appendContextIn(b *strings.Builder, contextIn []string, work *os.Root) {
	if work == nil || len(contextIn) == 0 {
		return
	}
	for _, path := range contextIn {
		clean := filepath.FromSlash(strings.TrimSpace(path))
		if clean == "" || !filepath.IsLocal(clean) {
			continue
		}
		data, err := readLimited(work, clean)
		if err != nil {
			continue
		}
		sanitized := SanitizeContextIn(strings.TrimSpace(string(data)))
		fmt.Fprintf(b, "\n## Contexto de Entrada: %s\n\n<<<CONTEXT_IN_DATA\n%s\nCONTEXT_IN_DATA>>>\n", path, sanitized)
	}
}

var injectionRules = []struct {
	re   *regexp.Regexp
	repl string
}{
	{regexp.MustCompile(`(?i)ignore\s+(all\s+)?(previous|prior|above)\s+instructions?`), "[SANITIZED:instruction-override]"},
	{regexp.MustCompile(`(?i)you\s+are\s+now\s+`), "[SANITIZED:role-override]"},
	{regexp.MustCompile(`(?im)^\s*(system|assistant|developer)\s*:`), "[SANITIZED:role-tag]"},
	{regexp.MustCompile(`(?i)<\|im_start\|>|<\|im_end\|>|<\|system\|>|<system>|</system>`), "[SANITIZED:control-token]"},
}

// SanitizeContextIn neutraliza padrões conhecidos de prompt injection em dados externos.
func SanitizeContextIn(raw string) string {
	out := raw
	for _, r := range injectionRules {
		out = r.re.ReplaceAllString(out, r.repl)
	}
	return out
}

func (o *Orchestrix) appendTooling(b *strings.Builder, tooling map[string]any) {
	for _, key := range []string{"skills", "rules", "skill", "rule"} {
		for _, ref := range extractStringRefs(tooling[key]) {
			content, err := o.readSkill(ref)
			if err != nil {
				continue
			}
			fmt.Fprintf(b, "\n## Regra/Skill (%s)\n\n%s\n", ref, strings.TrimSpace(content))
		}
	}
}

func extractStringRefs(val any) []string {
	switch v := val.(type) {
	case string:
		if strings.TrimSpace(v) != "" {
			return []string{strings.TrimSpace(v)}
		}
	case []string:
		return v
	case []any:
		out := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok && strings.TrimSpace(s) != "" {
				out = append(out, strings.TrimSpace(s))
			}
		}
		return out
	}
	return nil
}

func (o *Orchestrix) readSkill(ref string) (string, error) {
	if !filepath.IsLocal(filepath.FromSlash(ref)) {
		return "", fmt.Errorf("ref %q fora do cofre de skills", ref)
	}
	root, err := os.OpenRoot(o.SkillsDir)
	if err != nil {
		return "", fmt.Errorf("cofre de skills: %w", err)
	}
	defer root.Close()
	data, err := readLimited(root, filepath.Join(filepath.FromSlash(ref), "SKILL.md"))
	if err != nil {
		return "", fmt.Errorf("skill %q: %w", ref, err)
	}
	return string(data), nil
}

// resolveInput: "var.X" ⇒ variável; "file('P')" ⇒ artefato do workdir; senão literal.
func (o *Orchestrix) resolveInput(spec string, work *os.Root) (string, error) {
	if name, ok := strings.CutPrefix(spec, "var."); ok {
		v, ok := o.Vars[name]
		if !ok {
			return "", fmt.Errorf("variável %q não definida", name)
		}
		return v, nil
	}
	if m := fileInputRe.FindStringSubmatch(spec); m != nil {
		data, err := readLimited(work, filepath.FromSlash(m[1]))
		return string(data), err
	}
	return spec, nil
}

// readLimited abre uma vez, exige arquivo regular e lê no máximo o limite.
func readLimited(root *os.Root, name string) ([]byte, error) {
	f, err := root.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if info, err := f.Stat(); err != nil || !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%q não é arquivo regular", name)
	}
	data, err := io.ReadAll(io.LimitReader(f, maxInputBytes+1))
	if err == nil && len(data) > maxInputBytes {
		err = fmt.Errorf("%q excede %d bytes", name, maxInputBytes)
	}
	return data, err
}

var domainRulesByTag = map[string]string{
	"go":        "Go: execute `gofmt`, `go vet` e mantenha funcoes <= 35 linhas.",
	"ui":        "UI: respeite tokens de design, contraste WCAG e responsividade.",
	"seguranca": "Seguranca: valide entradas, aplique Path Jail e proiba segredos em logs.",
}

func (o *Orchestrix) appendRuleTags(b *strings.Builder, tags []string) {
	if len(tags) == 0 {
		return
	}
	var regras []string
	for _, tag := range tags {
		if r, ok := domainRulesByTag[strings.ToLower(strings.TrimSpace(tag))]; ok {
			regras = append(regras, "- "+r)
		}
	}
	if len(regras) == 0 {
		return
	}
	fmt.Fprintf(b, "\n## Regras de Dominio (Lazy)\n\n%s\n", strings.Join(regras, "\n"))
}

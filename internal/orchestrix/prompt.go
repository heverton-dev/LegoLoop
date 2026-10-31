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

var fileInputRe = regexp.MustCompile(`^file\('([^']+)'\)$`)

// buildPrompt monta o prompt único e estéril: skill + entradas resolvidas + contrato de saída.
func (o *Orchestrix) buildPrompt(n dag.Node, work *os.Root) (string, error) {
	skill, err := o.readSkill(n.Ref)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	fmt.Fprintf(&b, "# Orchestrix isolada — nó %s\n\n", n.ID)
	fmt.Fprintf(&b, "Você é um executor de tarefa única, sem histórico. Responda SOMENTE com o conteúdo final do artefato `%s`, sem comentários extras.\n\n", n.Output.Artifact)
	fmt.Fprintf(&b, "Condição de saída (gate): `%s`\n\n## Skill\n\n%s\n", n.Output.Gate, strings.TrimSpace(skill))
	keys := make([]string, 0, len(n.Inputs))
	for k := range n.Inputs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v, err := o.resolveInput(n.Inputs[k], work)
		if err != nil {
			return "", fmt.Errorf("entrada %q: %w", k, err)
		}
		fmt.Fprintf(&b, "\n## Entrada: %s\n\n%s\n", k, v)
	}
	return b.String(), nil
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

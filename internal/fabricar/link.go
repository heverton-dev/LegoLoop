package fabricar

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// harnessMatches diz se a pasta de skills pertence ao harness (ex.: "claude" casa com ~/.claude/skills).
func harnessMatches(harness, dir string) bool {
	return strings.Contains(strings.ToLower(filepath.ToSlash(dir))+"/", "/."+strings.ToLower(strings.TrimSpace(harness))+"/")
}

// linkHarnesses liga a skill às pastas dos harnesses pedidos. Só age com opt-in (LinkHarnesses) e fora do dry-run:
// o cofre não é lido pelos harnesses, então sem o link o /<slug> não aparece no Chat.
func linkHarnesses(req Request, dirs []string, skillPath, name string) ([]string, []string) {
	if len(req.LinkHarnesses) == 0 || req.DryRun {
		return nil, nil
	}
	var links, warns []string
	for _, h := range req.LinkHarnesses {
		matched := false
		for _, dir := range dirs {
			if !harnessMatches(h, dir) {
				continue
			}
			matched = true
			link, err := linkSkill(filepath.Dir(skillPath), filepath.Join(dir, name))
			if err != nil {
				warns = append(warns, fmt.Sprintf("harness %q: %v", h, err))
				continue
			}
			links = append(links, link)
		}
		if !matched {
			warns = append(warns, fmt.Sprintf("harness %q sem pasta de skills configurada", h))
		}
	}
	return links, warns
}

// linkSkill cria dst como symlink para o diretório da skill; se o sistema não permitir, copia o SKILL.md.
func linkSkill(srcDir, dst string) (string, error) {
	if _, err := os.Lstat(dst); err == nil {
		return "", fmt.Errorf("já existe: %s", dst)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return "", err
	}
	if err := os.Symlink(srcDir, dst); err == nil {
		return dst, nil
	}
	data, err := os.ReadFile(filepath.Join(srcDir, "SKILL.md"))
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return "", err
	}
	return dst, os.WriteFile(filepath.Join(dst, "SKILL.md"), data, 0o644)
}

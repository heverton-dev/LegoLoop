package assets

import (
	"os"
	"path/filepath"
)

// scaffoldDirs lista as pastas diretas de scaffolds do ecossistema AIDD.
func (s *Scanner) scaffoldDirs() []string {
	if s.cfg.Home == "" {
		return nil
	}
	base := filepath.Join(s.cfg.Home, ".aidd", "scaffolds")
	entries, err := os.ReadDir(base)
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() && !IsIgnoredDir(e.Name()) {
			out = append(out, filepath.Join(base, e.Name()))
		}
	}
	return out
}

// scaffoldSkillRoots devolve raízes de skills de scaffolds com scope "scaffold".
func (s *Scanner) scaffoldSkillRoots() []root {
	var roots []root
	for _, d := range s.scaffoldDirs() {
		sub := filepath.Join(d, "skills")
		if st, err := os.Stat(sub); err == nil && st.IsDir() {
			roots = append(roots, root{dir: sub, scope: "scaffold", origin: "aidd"})
		}
	}
	return roots
}

// scaffoldHookSources devolve arquivos settings.json de scaffolds com scope "scaffold".
func (s *Scanner) scaffoldHookSources() []hookSource {
	var sources []hookSource
	for _, d := range s.scaffoldDirs() {
		for _, rel := range []string{
			filepath.Join(".claude", "settings.json"),
			"settings.json",
		} {
			p := filepath.Join(d, rel)
			if st, err := os.Stat(p); err == nil && !st.IsDir() {
				sources = append(sources, hookSource{path: p, scope: "scaffold"})
			}
		}
	}
	return sources
}

// scaffoldLoopRoots devolve raízes de loops de scaffolds com scope "scaffold".
func (s *Scanner) scaffoldLoopRoots() []root {
	var roots []root
	for _, d := range s.scaffoldDirs() {
		sub := filepath.Join(d, "loops")
		if st, err := os.Stat(sub); err == nil && st.IsDir() {
			roots = append(roots, root{dir: sub, scope: "scaffold", origin: "aidd"})
		}
	}
	return roots
}

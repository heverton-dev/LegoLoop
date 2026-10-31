package assets

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var ignoredDirNames = map[string]bool{
	"node_modules": true,
	".git":         true,
	"vendor":       true,
	".venv":        true,
	"venv":         true,
	"__pycache__":  true,
	".cache":       true,
	"dist":         true,
	"build":        true,
	".next":        true,
	".nuget":       true,
	"target":       true,
}

var projectMarkers = []string{
	".git",
	"package.json",
	"go.mod",
	"Cargo.toml",
	"pyproject.toml",
	".claude",
	".gemini",
	"pom.xml",
	"requirements.txt",
}

// IsIgnoredDir diz se um diretório deve ser podado da varredura.
func IsIgnoredDir(name string) bool {
	return ignoredDirNames[strings.ToLower(name)]
}

// isProjectDir verifica se o diretório contém algum marcador canônico de projeto.
func isProjectDir(dir string) bool {
	for _, m := range projectMarkers {
		p := filepath.Join(dir, m)
		if _, err := os.Stat(p); err == nil {
			return true
		}
	}
	return false
}

// DiscoverProjects descobre concorrentemente diretórios de projeto a partir das raízes.
func DiscoverProjects(ctx context.Context, roots []string, maxDepth int) []string {
	var mu sync.Mutex
	found := make(map[string]bool)
	var wg sync.WaitGroup

	for _, r := range roots {
		if st, err := os.Stat(r); err != nil || !st.IsDir() {
			continue
		}
		wg.Add(1)
		go func(rootPath string) {
			defer wg.Done()
			walkProjects(ctx, rootPath, 0, maxDepth, func(p string) {
				mu.Lock()
				found[p] = true
				mu.Unlock()
			})
		}(r)
	}
	wg.Wait()

	out := make([]string, 0, len(found))
	for p := range found {
		out = append(out, p)
	}
	return out
}

func walkProjects(ctx context.Context, cur string, depth, maxDepth int, emit func(string)) {
	if depth > maxDepth || ctx.Err() != nil {
		return
	}
	entries, err := os.ReadDir(cur)
	if err != nil {
		return
	}
	if depth > 0 && isProjectDir(cur) {
		emit(cur)
	}
	for _, e := range entries {
		if !e.IsDir() || IsIgnoredDir(e.Name()) {
			continue
		}
		walkProjects(ctx, filepath.Join(cur, e.Name()), depth+1, maxDepth, emit)
	}
}

// DefaultDiscoveryRoots devolve as pastas comuns de trabalho de desenvolvimento no SO.
func DefaultDiscoveryRoots(home string) []string {
	if home == "" {
		return nil
	}
	candidates := []string{
		"Desktop", "Projects", "Workspace", "repos", "src", "code", "dev", "Documents",
	}
	var out []string
	for _, c := range candidates {
		p := filepath.Join(home, c)
		if st, err := os.Stat(p); err == nil && st.IsDir() {
			out = append(out, p)
		}
	}
	return out
}

// allProjects reúne projetos declarados e descobertos sem duplicação.
func (s *Scanner) allProjects(ctx context.Context) []string {
	seen := make(map[string]bool)
	var list []string
	add := func(p string) {
		clean := filepath.Clean(p)
		if !seen[clean] {
			seen[clean] = true
			list = append(list, clean)
		}
	}
	for _, p := range s.cfg.Projects {
		add(p)
	}
	if s.cfg.AutoDiscover {
		roots := s.cfg.DiscoveryRoots
		if len(roots) == 0 {
			roots = DefaultDiscoveryRoots(s.cfg.Home)
		}
		for _, p := range DiscoverProjects(ctx, roots, 4) {
			add(p)
		}
	}
	return list
}

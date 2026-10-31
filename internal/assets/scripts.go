package assets

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

// maxScriptDepth é o nível máximo de subpastas abaixo de uma raiz de scripts.
const maxScriptDepth = 3

var scriptExts = map[string]bool{
	".ps1": true, ".sh": true, ".bash": true, ".py": true, ".js": true,
	".mjs": true, ".cjs": true, ".bat": true, ".cmd": true,
}

func (s *Scanner) scriptRoots(ctx context.Context) []root {
	var roots []root
	seen := make(map[string]bool)
	add := func(r root) {
		clean := filepath.Clean(r.dir)
		if !seen[clean] {
			seen[clean] = true
			roots = append(roots, r)
		}
	}
	for _, r := range s.joinRoots(
		[]root{{".aidd/scripts", "", "aidd"}, {".claude/scripts", "", "claude"}},
		[]root{{"scripts", "", "project"}, {"tools", "", "project"}, {".claude/scripts", "", "claude"}},
	) {
		add(r)
	}
	for _, p := range s.allProjects(ctx) {
		for _, sub := range []string{"scripts", "tools", "bin"} {
			dir := filepath.Join(p, sub)
			if st, err := os.Stat(dir); err == nil && st.IsDir() {
				add(root{dir: dir, scope: "project", origin: "project"})
			}
		}
	}
	for _, d := range s.scaffoldDirs() {
		sub := filepath.Join(d, "scripts")
		if st, err := os.Stat(sub); err == nil && st.IsDir() {
			add(root{dir: sub, scope: "scaffold", origin: "aidd"})
		}
	}
	return roots
}

func (s *Scanner) scanScripts(ctx context.Context) []Item {
	var out []Item
	for _, r := range s.scriptRoots(ctx) {
		walkScripts(ctx, r, r.dir, 0, &out)
	}
	return out
}

func walkScripts(ctx context.Context, r root, dir string, depth int, out *[]Item) {
	if ctx.Err() != nil {
		return
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		path := filepath.Join(dir, e.Name())
		switch {
		case e.IsDir():
			if !skipDirName(e.Name()) && depth < maxScriptDepth {
				walkScripts(ctx, r, path, depth+1, out)
			}
		case e.Type().IsRegular() && scriptExts[strings.ToLower(filepath.Ext(e.Name()))]:
			*out = append(*out, Item{Kind: KindScript, Name: e.Name(), Path: path, Scope: r.scope, Origin: r.origin})
		}
	}
}

package assets

import (
	"context"
	"encoding/json"
	"legoloop/internal/caminhos"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"legoloop/internal/dag"
)

func (s *Scanner) loopRoots(ctx context.Context) []root {
	var roots []root
	if s.cfg.LoopsDir != "" {
		roots = append(roots, root{dir: s.cfg.LoopsDir, scope: "global", origin: "acervo"})
	}
	if s.cfg.Home != "" {
		roots = append(roots,
			root{dir: filepath.Join(s.cfg.Home, ".aidd", "canonical", caminhos.DirLoops), scope: "global", origin: "acervo"},
			root{dir: filepath.Join(s.cfg.Home, ".aidd", caminhos.DirLoops), scope: "global", origin: "acervo"},
		)
	}
	roots = append(roots, s.scaffoldLoopRoots()...)
	for _, p := range s.allProjects(ctx) {
		roots = append(roots, root{dir: p, scope: "project", origin: "project"})
		sub := caminhos.DirLoopsEm(p)
		if st, err := os.Stat(sub); err == nil && st.IsDir() {
			roots = append(roots, root{dir: sub, scope: "project", origin: "project"})
		}
	}
	return roots
}

func (s *Scanner) scanLoops(ctx context.Context) []Item {
	var items []Item
	seen := map[string]bool{}
	for _, r := range s.loopRoots(ctx) {
		entries, err := os.ReadDir(r.dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if ctx.Err() != nil {
				return items
			}
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".loop.json") {
				continue
			}
			p := filepath.Join(r.dir, e.Name())
			if seen[p] {
				continue
			}
			seen[p] = true
			items = append(items, s.load(KindLoop, p, func(path string) []Item {
				return parseLoopItem(path, r)
			})...)
		}
	}
	return items
}

func parseLoopItem(path string, r root) []Item {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var lp struct {
		ID          string       `json:"id"`
		Name        string       `json:"name"`
		Description string       `json:"description"`
		Trigger     *dag.Trigger `json:"trigger"`
		Nodes       []dag.Node   `json:"nodes"`
	}
	if err := json.Unmarshal(data, &lp); err != nil {
		return nil
	}
	id := lp.ID
	if id == "" {
		id = strings.TrimSuffix(filepath.Base(path), ".loop.json")
	}
	name := lp.Name
	if name == "" {
		name = id
	}
	meta := buildLoopMeta(lp.Nodes, lp.Trigger)
	return []Item{{
		Kind:        KindLoop,
		Name:        name,
		Description: lp.Description,
		Path:        path,
		Scope:       r.scope,
		Origin:      r.origin,
		Meta:        meta,
	}}
}

func buildLoopMeta(nodes []dag.Node, trg *dag.Trigger) map[string]string {
	meta := map[string]string{
		"nodes": strconv.Itoa(len(nodes)),
	}
	if len(nodes) > 0 {
		meta["artifact"] = nodes[len(nodes)-1].Output.Artifact
	}
	if trg != nil && trg.Type != "" {
		meta["trigger"] = string(trg.Type)
	} else {
		meta["trigger"] = "manual"
	}
	return meta
}

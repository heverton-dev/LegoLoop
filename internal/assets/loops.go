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

const sufixoAntigo = caminhos.SufixoLoop

var kindProcesso = Kinds()[len(Kinds())-1]

func isProcessoDoc(name string) bool {
	return strings.HasSuffix(name, sufixoAntigo) || strings.HasSuffix(name, caminhos.SufixoStudio)
}

func trimProcessoDoc(name string) string {
	for _, suf := range []string{sufixoAntigo, caminhos.SufixoStudio} {
		name = strings.TrimSuffix(name, suf)
	}
	return name
}

func (s *Scanner) processoRoots(ctx context.Context) []root {
	var roots []root
	if dir := s.cfg.LoopsDir; dir != "" {
		roots = append(roots, root{dir: dir, scope: "global", origin: "acervo"})
	}
	if s.cfg.Home != "" {
		for _, base := range []string{filepath.Join(s.cfg.Home, ".aidd", "canonical"), filepath.Join(s.cfg.Home, ".aidd")} {
			for _, dir := range caminhos.DirsLeitura(base) {
				roots = append(roots, root{dir: dir, scope: "global", origin: "acervo"})
			}
		}
	}
	roots = append(roots, s.scaffoldProcessoRoots()...)
	for _, p := range s.allProjects(ctx) {
		roots = append(roots, root{dir: p, scope: "project", origin: "project"})
		for _, sub := range caminhos.DirsLeitura(p) {
			roots = append(roots, root{dir: sub, scope: "project", origin: "project"})
		}
	}
	return roots
}

func (s *Scanner) scanProcessos(ctx context.Context) []Item {
	var items []Item
	seen := map[string]bool{}
	for _, r := range s.processoRoots(ctx) {
		entries, err := os.ReadDir(r.dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if ctx.Err() != nil {
				return items
			}
			if e.IsDir() || !isProcessoDoc(e.Name()) {
				continue
			}
			p := filepath.Join(r.dir, e.Name())
			if seen[p] {
				continue
			}
			seen[p] = true
			items = append(items, s.load(kindProcesso, p, func(path string) []Item {
				return parseProcessoItem(path, r)
			})...)
		}
	}
	return items
}

func parseProcessoItem(path string, r root) []Item {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var doc struct {
		ID          string       `json:"id"`
		Name        string       `json:"name"`
		Description string       `json:"description"`
		Trigger     *dag.Trigger `json:"trigger"`
		Nodes       []dag.Node   `json:"nodes"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil
	}
	id := doc.ID
	if id == "" {
		id = trimProcessoDoc(filepath.Base(path))
	}
	name := doc.Name
	if name == "" {
		name = id
	}
	meta := buildProcessoMeta(doc.Nodes, doc.Trigger)
	return []Item{{
		Kind:        kindProcesso,
		Name:        name,
		Description: doc.Description,
		Path:        path,
		Scope:       r.scope,
		Origin:      r.origin,
		Meta:        meta,
	}}
}

func buildProcessoMeta(nodes []dag.Node, trg *dag.Trigger) map[string]string {
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

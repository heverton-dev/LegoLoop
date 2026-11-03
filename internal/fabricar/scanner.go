package fabricar

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"legoloop/internal/assets"
)

// DiskScanner varre raízes (unidades, volumes) reaproveitando o scanner de internal/assets, que já poda
// node_modules e .git, não segue symlinks e limita a profundidade. Não existe um segundo walker.
type DiskScanner struct {
	// ScanRoot é a borda de leitura de uma raiz; o padrão usa internal/assets. Injetável em teste.
	ScanRoot func(ctx context.Context, root string) ([]assets.Item, error)
	// Untrusted diz se a origem é removível ou de rede; o padrão consulta o sistema.
	Untrusted func(root string) bool
}

// NewDiskScanner devolve o scanner com as bordas reais.
func NewDiskScanner() *DiskScanner {
	return &DiskScanner{ScanRoot: assetsScanRoot, Untrusted: isUntrustedRoot}
}

// ScanStorageLocations varre as raízes com o scanner padrão.
func ScanStorageLocations(ctx context.Context, roots []string, opts ScanOptions) (ScanResult, error) {
	return NewDiskScanner().Scan(ctx, roots, opts)
}

// Scan varre cada raiz com prazo próprio. Raiz lenta ou inacessível vira SkippedRoot; limites viram Warnings.
func (s *DiskScanner) Scan(ctx context.Context, roots []string, opts ScanOptions) (ScanResult, error) {
	opts = opts.withDefaults()
	var res ScanResult
	for _, root := range roots {
		if err := ctx.Err(); err != nil {
			return res, err
		}
		items, reason := s.scanOne(ctx, root, opts)
		if reason != "" {
			res.Skipped = append(res.Skipped, SkippedRoot{Root: root, Reason: reason})
			continue
		}
		pieces, warns := s.convert(root, items, opts)
		res.Pieces = append(res.Pieces, pieces...)
		res.Warnings = append(res.Warnings, warns...)
	}
	sortPieces(res.Pieces)
	return res, nil
}

// scanOne roda ScanRoot em goroutine própria para que uma raiz travada não passe do prazo.
func (s *DiskScanner) scanOne(ctx context.Context, root string, opts ScanOptions) ([]assets.Item, string) {
	rctx, cancel := context.WithTimeout(ctx, opts.PerRootTimeout)
	defer cancel()
	type out struct {
		items []assets.Item
		err   error
	}
	ch := make(chan out, 1)
	go func() {
		items, err := s.ScanRoot(rctx, root)
		ch <- out{items, err}
	}()
	select {
	case o := <-ch:
		if o.err != nil {
			return nil, reasonOf(o.err)
		}
		return o.items, ""
	case <-rctx.Done():
		return nil, "timeout"
	}
}

func reasonOf(err error) string {
	if errors.Is(err, context.DeadlineExceeded) {
		return "timeout"
	}
	return "inacessivel"
}

// convert aplica profundidade, MaxFiles e MaxPieces, sempre com aviso em vez de erro.
func (s *DiskScanner) convert(root string, items []assets.Item, opts ScanOptions) ([]DiscoveredPiece, []string) {
	var warns []string
	items = dedupe(items)
	slices.SortFunc(items, func(a, b assets.Item) int { return strings.Compare(a.Path+a.Name, b.Path+b.Name) })
	if kept := withinDepth(root, items, opts.MaxDepth); len(kept) < len(items) {
		warns = append(warns, fmt.Sprintf("raiz %s: %d item(ns) além da profundidade %d foram ignorados", root, len(items)-len(kept), opts.MaxDepth))
		items = kept
	}
	if len(items) > opts.MaxFiles {
		warns = append(warns, fmt.Sprintf("raiz %s: limite de %d itens inspecionados atingido", root, opts.MaxFiles))
		items = items[:opts.MaxFiles]
	}
	untrusted := s.Untrusted != nil && s.Untrusted(root)
	var pieces []DiscoveredPiece
	for _, it := range items {
		if p, ok := newPiece(it, root, untrusted); ok {
			pieces = append(pieces, p)
		}
	}
	sortPieces(pieces)
	if len(pieces) > opts.MaxPieces {
		warns = append(warns, fmt.Sprintf("raiz %s: limite de %d peças atingido", root, opts.MaxPieces))
		pieces = pieces[:opts.MaxPieces]
	}
	for i := range pieces {
		pieces[i].SHA256 = fileSHA(pieces[i].Path)
	}
	return pieces, warns
}

func dedupe(items []assets.Item) []assets.Item {
	seen := map[string]bool{}
	var out []assets.Item
	for _, it := range items {
		key := string(it.Kind) + "\x00" + it.Name + "\x00" + it.Path
		if !seen[key] {
			seen[key] = true
			out = append(out, it)
		}
	}
	return out
}

// depthOf conta os diretórios entre o perfil (raiz ou Users/<x>, home/<x>) e o arquivo.
func depthOf(root, path string) int {
	rel, err := filepath.Rel(root, path)
	if err != nil || strings.HasPrefix(rel, "..") {
		return 0
	}
	parts := strings.Split(filepath.ToSlash(rel), "/")
	if len(parts) > 2 && (parts[0] == "Users" || parts[0] == "home") {
		parts = parts[2:]
	}
	return len(parts) - 1
}

func withinDepth(root string, items []assets.Item, max int) []assets.Item {
	var kept []assets.Item
	for _, it := range items {
		if depthOf(root, it.Path) <= max {
			kept = append(kept, it)
		}
	}
	return kept
}

// sortPieces ordena por Score decrescente e depois por Path crescente (ordem estável).
func sortPieces(p []DiscoveredPiece) {
	slices.SortStableFunc(p, func(a, b DiscoveredPiece) int {
		switch {
		case a.Score > b.Score:
			return -1
		case a.Score < b.Score:
			return 1
		}
		return strings.Compare(a.Path, b.Path)
	})
}

func noLookPath(string) (string, error) { return "", os.ErrNotExist }
func noEnv(string) string               { return "" }

// assetsScanRoot inventaria uma raiz de forma dirigida: a própria raiz e os perfis em Users/* e home/*.
// Variáveis de ambiente e PATH ficam desligados, para que o inventário do usuário atual não vaze para a varredura.
func assetsScanRoot(ctx context.Context, root string) ([]assets.Item, error) {
	if _, err := os.Stat(root); err != nil {
		return nil, err
	}
	var items []assets.Item
	for _, home := range candidateHomes(root) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		sc := assets.NewScanner(assets.Config{Home: home, Projects: []string{home}, LookPath: noLookPath, Getenv: noEnv})
		items = append(items, sc.Scan(ctx).Items...)
	}
	return items, ctx.Err()
}

const maxProfiles = 50

// candidateHomes lista os perfis de usuário de uma raiz: ela mesma, Users/* e home/* (sem symlinks).
func candidateHomes(root string) []string {
	homes := []string{root}
	for _, parent := range []string{"Users", "home"} {
		entries, err := os.ReadDir(filepath.Join(root, parent))
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() && len(homes) < maxProfiles {
				homes = append(homes, filepath.Join(root, parent, e.Name()))
			}
		}
	}
	return homes
}

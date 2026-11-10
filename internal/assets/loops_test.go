package assets_test

import (
	"context"
	"legoloop/internal/caminhos"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/assets"
)

func TestScanLoops(t *testing.T) {
	dir := t.TempDir()
	loopsDir := caminhos.DirLoopsEm(dir)
	if err := os.MkdirAll(loopsDir, 0700); err != nil {
		t.Fatal(err)
	}

	content := `{
		"id": "pesquisa",
		"name": "Pesquisa de Mercado",
		"description": "Loop de pesquisa",
		"trigger": {"type": "manual"},
		"nodes": [
			{"id": "n1", "output": {"artifact": "BRIEFING.md"}}
		]
	}`
	if err := os.WriteFile(filepath.Join(loopsDir, "pesquisa.loop.json"), []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	scanner := assets.NewScanner(assets.Config{Home: t.TempDir(), LoopsDir: loopsDir})
	snap := scanner.Scan(context.Background())

	loops := snap.Of(assets.KindLoop)
	if len(loops) != 1 {
		t.Fatalf("esperava 1 loop, obteve %d", len(loops))
	}
	l := loops[0]
	if l.Name != "Pesquisa de Mercado" || l.Description != "Loop de pesquisa" {
		t.Fatalf("dados do loop divergentes: %+v", l)
	}
	if l.Meta["artifact"] != "BRIEFING.md" || l.Meta["nodes"] != "1" {
		t.Fatalf("meta divergente: %+v", l.Meta)
	}
}

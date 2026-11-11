package assets_test

import (
	"context"
	"legoloop/internal/caminhos"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/assets"
)

var kindProcesso = assets.Kinds()[len(assets.Kinds())-1]

func TestScanProcessos(t *testing.T) {
	dir := t.TempDir()
	processosDir := filepath.Join(dir, caminhos.DirProcessos)
	if err := os.MkdirAll(processosDir, 0700); err != nil {
		t.Fatal(err)
	}

	content := `{
		"id": "pesquisa",
		"name": "Pesquisa de Mercado",
		"description": "Processo de pesquisa",
		"trigger": {"type": "manual"},
		"nodes": [
			{"id": "n1", "output": {"artifact": "BRIEFING.md"}}
		]
	}`
	if err := os.WriteFile(filepath.Join(processosDir, "pesquisa"+caminhos.SufixoStudio), []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	scanner := assets.NewScanner(assets.Config{Home: t.TempDir(), LoopsDir: processosDir})
	snap := scanner.Scan(context.Background())

	processos := snap.Of(kindProcesso)
	if len(processos) != 1 {
		t.Fatalf("esperava 1 Processo, obteve %d", len(processos))
	}
	p := processos[0]
	if p.Name != "Pesquisa de Mercado" || p.Description != "Processo de pesquisa" {
		t.Fatalf("dados do Processo divergentes: %+v", p)
	}
	if p.Meta["artifact"] != "BRIEFING.md" || p.Meta["nodes"] != "1" {
		t.Fatalf("meta divergente: %+v", p.Meta)
	}
}

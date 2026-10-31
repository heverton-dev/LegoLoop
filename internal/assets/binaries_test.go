package assets_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/assets"
)

func TestScanBinariesInUserToolDirs(t *testing.T) {
	home := t.TempDir()
	cargoBin := filepath.Join(home, ".cargo", "bin")
	if err := os.MkdirAll(cargoBin, 0755); err != nil {
		t.Fatal(err)
	}

	// Cria um binário simulado (ex: "rg" ou "cargo") apenas dentro de ~/.cargo/bin
	fakeBin := filepath.Join(cargoBin, "rg")
	if err := os.WriteFile(fakeBin, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatal(err)
	}

	// LookPath que falha (simulando que ~/.cargo/bin não está no PATH)
	cfg := assets.Config{
		Home: home,
		LookPath: func(string) (string, error) {
			return "", os.ErrNotExist
		},
	}
	scanner := assets.NewScanner(cfg)
	snap := scanner.Scan(context.Background())

	var found bool
	for _, it := range snap.Of(assets.KindBinary) {
		if it.Name == "rg" {
			found = true
			if !filepath.IsAbs(it.Path) {
				t.Errorf("caminho do binário deveria ser absoluto: %s", it.Path)
			}
		}
	}
	if !found {
		t.Fatal("binário 'rg' presente em ~/.cargo/bin deveria ter sido descoberto mesmo fora do PATH")
	}
}

package assets

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// testCfg monta uma Config isolada: HOME e projeto em diretórios temporários e PATH falso.
func testCfg(t *testing.T, onPath ...string) (Config, string, string) {
	t.Helper()
	home, project := t.TempDir(), t.TempDir()
	found := map[string]string{}
	for _, n := range onPath {
		found[n] = filepath.Join(home, "bin", n)
	}
	return Config{
		Home:     home,
		Projects: []string{project},
		LookPath: func(n string) (string, error) {
			if p, ok := found[n]; ok {
				return p, nil
			}
			return "", errors.New("não encontrado")
		},
		Getenv: func(k string) string {
			if k == "APPDATA" {
				return filepath.Join(home, "AppData", "Roaming")
			}
			return ""
		},
	}, home, project
}

func scanKind(t *testing.T, cfg Config, k Kind) []Item {
	t.Helper()
	return NewScanner(cfg).Scan(t.Context()).Of(k)
}

func find(items []Item, name string) (Item, bool) {
	for _, it := range items {
		if it.Name == name {
			return it, true
		}
	}
	return Item{}, false
}

func mustFind(t *testing.T, items []Item, name string) Item {
	t.Helper()
	it, ok := find(items, name)
	if !ok {
		names := []string{}
		for _, i := range items {
			names = append(names, i.Name)
		}
		t.Fatalf("item %q não encontrado em %v", name, names)
	}
	return it
}

func mustNotLeak(t *testing.T, items []Item, segredos ...string) {
	t.Helper()
	b, err := json.Marshal(items)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range segredos {
		if strings.Contains(string(b), s) {
			t.Errorf("segredo %q vazou para o inventário: %s", s, b)
		}
	}
}

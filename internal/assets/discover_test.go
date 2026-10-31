package assets_test

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"legoloop/internal/assets"
)

func TestDiscoverProjects(t *testing.T) {
	root := t.TempDir()

	// Cria estrutura com:
	// - projeto1 (com .git)
	// - projeto2 (com package.json)
	// - sub/projeto3 (com go.mod aninhado a profundidade 2)
	// - node_modules/falso_projeto (deve ser podado / ignorado)
	// - .git/objects (deve ser podado)
	// - pasta_vazia (sem marcadores)

	p1 := filepath.Join(root, "proj1")
	p2 := filepath.Join(root, "proj2")
	p3 := filepath.Join(root, "workspace", "proj3")
	bad1 := filepath.Join(root, "node_modules", "lib1")
	bad2 := filepath.Join(p1, "node_modules", "lib2")
	empty := filepath.Join(root, "pasta_sem_projeto")

	for _, d := range []string{
		filepath.Join(p1, ".git"),
		p2,
		p3,
		bad1,
		bad2,
		empty,
	} {
		if err := os.MkdirAll(d, 0755); err != nil {
			t.Fatal(err)
		}
	}

	_ = os.WriteFile(filepath.Join(p2, "package.json"), []byte("{}"), 0644)
	_ = os.WriteFile(filepath.Join(p3, "go.mod"), []byte("module proj3"), 0644)
	_ = os.WriteFile(filepath.Join(bad1, "package.json"), []byte("{}"), 0644)
	_ = os.WriteFile(filepath.Join(bad2, "package.json"), []byte("{}"), 0644)

	projects := assets.DiscoverProjects(context.Background(), []string{root}, 4)

	var names []string
	for _, p := range projects {
		names = append(names, filepath.Base(p))
	}
	sort.Strings(names)

	want := []string{"proj1", "proj2", "proj3"}
	if len(names) != len(want) {
		t.Fatalf("projetos descobertos = %v, quer %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Errorf("projeto %d: got %s, want %s", i, names[i], want[i])
		}
	}
}

func TestDiscoverProjectsRespectsMaxDepth(t *testing.T) {
	root := t.TempDir()
	deep := filepath.Join(root, "a", "b", "c", "d", "e", "f", "deep_proj")
	if err := os.MkdirAll(deep, 0755); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(deep, "go.mod"), []byte("module deep"), 0644)

	// Com maxDepth = 3, não deve achar o deep_proj que está a profundidade 7
	projects := assets.DiscoverProjects(context.Background(), []string{root}, 3)
	if len(projects) != 0 {
		t.Fatalf("esperava 0 projetos além do limite de profundidade, achou: %v", projects)
	}
}

func TestDiscoverProjectsPrunesIgnoredDirs(t *testing.T) {
	for _, ign := range []string{"node_modules", ".venv", "vendor", "__pycache__", ".cache"} {
		if !assets.IsIgnoredDir(ign) {
			t.Errorf("esperava que %s fosse ignorado", ign)
		}
	}
	if assets.IsIgnoredDir("meu_projeto") {
		t.Error("meu_projeto não deve ser ignorado")
	}
}

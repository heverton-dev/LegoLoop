package roteamento

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"legoloop/internal/processo"
)

func TestInventarioPorHarness(t *testing.T) {
	home := t.TempDir()
	for _, n := range []string{".claude", ".opencode"} {
		if err := os.MkdirAll(filepath.Join(home, n), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	binario := func(n string) string { return "/bin/" + n }
	modelos := map[string][]string{"claude": {"haiku", "sonnet"}, "opencode": {"modelo-x"}}

	inv := LerInventario(processo.OSFS(), home, binario, modelos)
	want := []HarnessInv{
		{Nome: "claude", Binario: "/bin/claude", Modelos: []string{"haiku", "sonnet"}},
		{Nome: "opencode", Binario: "/bin/opencode", Modelos: []string{"modelo-x"}},
	}
	if !reflect.DeepEqual(inv.Harnesses, want) {
		t.Fatalf("harnesses = %+v, quer %+v", inv.Harnesses, want)
	}
	if inv.Vazio() || inv.SO == "" {
		t.Fatalf("inventário preenchido inválido: %+v", inv)
	}

	vazio := LerInventario(processo.OSFS(), t.TempDir(), binario, modelos)
	if !vazio.Vazio() {
		t.Fatalf("inventário deveria ser vazio: %+v", vazio)
	}
}

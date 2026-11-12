package roteamento

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"legoloop/internal/processo"
	"legoloop/internal/provider"
)

func TestInventarioUnificado(t *testing.T) {
	home := t.TempDir()
	for _, n := range []string{".claude", ".opencode", ".mimo", ".omp", ".agy"} {
		if err := os.MkdirAll(filepath.Join(home, n), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	prov := provider.Inventory{
		RouterModels: []string{"combo-a"},
		CLIs: []provider.CLI{
			{Name: "claude", Path: "/bin/claude"},
			{Name: "opencode", Path: "/bin/opencode"},
			{Name: "agy", Path: "/bin/agy"},
		},
	}
	store := map[string]Homologacao{
		"claude|oficial": {Estado: EstadoHomologado, PesoMB: 300},
		"agy|oficial":    {Estado: EstadoReprovado, Motivo: "tos"},
	}
	itens, err := LerInventarioUnificado(processo.OSFS(), home, prov, store)
	if err != nil {
		t.Fatal(err)
	}
	porPar := map[string]ItemInv{}
	for _, it := range itens {
		if it.Par.Harness == "9router" {
			t.Fatalf("9router não pode ser harness: %+v", it)
		}
		porPar[it.Par.Chave()] = it
	}
	if len(itens) != 10 {
		t.Fatalf("itens = %d, quer 10", len(itens))
	}
	c := porPar["claude|oficial"]
	if c.Estado != EstadoHomologado || c.PesoMB != 300 || c.Binario != "/bin/claude" ||
		!reflect.DeepEqual(c.Modelos, []string{"haiku", "sonnet", "opus"}) {
		t.Fatalf("claude oficial = %+v", c)
	}
	if r := porPar["claude|9router"]; !reflect.DeepEqual(r.Modelos, []string{"combo-a"}) || r.Estado != EstadoPendente {
		t.Fatalf("claude 9router = %+v", r)
	}
	if a := porPar["agy|oficial"]; a.Estado != EstadoReprovado || a.Motivo != "tos" {
		t.Fatalf("agy oficial = %+v", a)
	}
	if m := porPar["mimo|oficial"]; m.Estado != EstadoNaoInstalado {
		t.Fatalf("mimo oficial = %+v", m)
	}
	if o := porPar["opencode|oficial"]; o.Estado != EstadoPendente || !reflect.DeepEqual(o.Modelos, []string{"opencode:default"}) {
		t.Fatalf("opencode oficial = %+v", o)
	}
	if _, err := LerInventarioUnificado(processo.OSFS(), t.TempDir(), provider.Inventory{}, store); err != ErrInventarioVazio {
		t.Fatalf("err = %v, quer ErrInventarioVazio", err)
	}
}

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

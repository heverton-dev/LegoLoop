package espinha_test

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/cli"
	"github.com/heverton-dev/LegoLoop/internal/espinha"
)

func TestTodoVerboEstaNoRegistro(t *testing.T) {
	caminho := filepath.Join("..", "..", "espinha", "registro.json")
	reg, err := espinha.CarregarRegistro(caminho)
	if err != nil {
		t.Fatalf("CarregarRegistro(%s) falhou: %v", caminho, err)
	}
	verbosCLI := cli.Verbos()
	if len(verbosCLI) == 0 {
		t.Fatal("cli.Verbos() retornou lista vazia")
	}
	if !slices.IsSorted(verbosCLI) {
		t.Fatalf("cli.Verbos() nao esta ordenado: %v", verbosCLI)
	}
	noRegistro := mapVerbosRegistro(reg.Verbos)
	for _, verbo := range verbosCLI {
		if !noRegistro[verbo] {
			t.Errorf("verbo %q de cli.Verbos() ausente em espinha/registro.json", verbo)
		}
	}
	verificarVerbosForaDoHelp(t, verbosCLI, noRegistro)
}

func mapVerbosRegistro(verbos []espinha.Verbo) map[string]bool {
	m := make(map[string]bool, len(verbos))
	for _, v := range verbos {
		m[v.Nome] = true
	}
	return m
}

func verificarVerbosForaDoHelp(t *testing.T, verbosCLI []string, noRegistro map[string]bool) {
	t.Helper()
	for _, obrigatorio := range []string{"commit", "legoloop-cartoon-shorts"} {
		if !slices.Contains(verbosCLI, obrigatorio) {
			t.Errorf("verbo %q ausente em cli.Verbos()", obrigatorio)
		}
		if !noRegistro[obrigatorio] {
			t.Errorf("verbo %q ausente em espinha/registro.json", obrigatorio)
		}
	}
}

package processo_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
)

var fatiasCanonicas = []string{
	"legoloop/internal/fabricar",
	"legoloop/internal/listar",
	"legoloop/internal/refinar",
	"legoloop/internal/remover",
	"legoloop/internal/compor",
}

func acharRaizModulo(t *testing.T) string {
	t.Helper()
	_, arquivo, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("não foi possível obter o caminho do arquivo de teste via runtime.Caller")
	}
	dir := filepath.Dir(arquivo)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		pai := filepath.Dir(dir)
		if pai == dir {
			t.Fatal("raiz do módulo (go.mod) não encontrada subindo a partir de " + arquivo)
		}
		dir = pai
	}
}

func TestFatiasIndependentes(t *testing.T) {
	raiz := acharRaizModulo(t)

	for _, fatia := range fatiasCanonicas {
		fatia := fatia
		t.Run(fatia, func(t *testing.T) {
			cmd := exec.Command("go", "list", "-deps", "-f", "{{.ImportPath}}", fatia)
			cmd.Dir = raiz
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				t.Fatalf("go list falhou para o pacote %s: %v\nstderr: %s", fatia, err, stderr.String())
			}

			linhas := strings.Split(strings.TrimSpace(stdout.String()), "\n")
			if len(linhas) == 0 || (len(linhas) == 1 && linhas[0] == "") {
				t.Fatalf("go list não retornou dependências para o pacote %s", fatia)
			}

			deps := make(map[string]bool, len(linhas))
			for _, l := range linhas {
				deps[strings.TrimSpace(l)] = true
			}

			for _, outra := range fatiasCanonicas {
				if outra == fatia {
					continue
				}
				if deps[outra] {
					t.Errorf("a fatia %s depende indevidamente da fatia %s; fatias devem ser independentes", fatia, outra)
				}
			}
		})
	}
}

func TestFatiasCanonicasDefinidas(t *testing.T) {
	esperadas := []string{
		"legoloop/internal/fabricar",
		"legoloop/internal/listar",
		"legoloop/internal/refinar",
		"legoloop/internal/remover",
		"legoloop/internal/compor",
	}
	if len(fatiasCanonicas) != 5 {
		t.Fatalf("esperava 5 fatias canônicas, tem %d", len(fatiasCanonicas))
	}
	for _, esp := range esperadas {
		if !slices.Contains(fatiasCanonicas, esp) {
			t.Errorf("fatia canônica ausente: %s", esp)
		}
	}
}

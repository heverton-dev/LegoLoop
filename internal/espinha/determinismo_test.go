package espinha

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPacotesDeterministicosSemReferenciasProibidas(t *testing.T) {
	proibidos := []string{
		"LLM" + "Client",
		"Resolvedor" + "LLM",
		"time." + "Now",
	}
	dirs := []string{
		filepath.Join("..", "espinha"),
		filepath.Join("..", "descoberta"),
	}
	for _, dir := range dirs {
		verificarDirSemProibidos(t, dir, proibidos)
	}
}

func verificarDirSemProibidos(t *testing.T, dir string, proibidos []string) {
	t.Helper()
	entradas, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ler %s: %v", dir, err)
	}
	for _, e := range entradas {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		caminho := filepath.Join(dir, e.Name())
		dados, err := os.ReadFile(caminho)
		if err != nil {
			t.Fatalf("ler %s: %v", caminho, err)
		}
		texto := string(dados)
		for _, termo := range proibidos {
			if strings.Contains(texto, termo) {
				t.Errorf("%s referencia simbolo proibido %q (AC-NF.1)", caminho, termo)
			}
		}
	}
}

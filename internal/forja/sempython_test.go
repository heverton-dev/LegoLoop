package forja

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCaminhoSemPython(t *testing.T) {
	raiz := localizarRaizModulo(t)
	pacotes := []string{
		filepath.Join(raiz, "internal", "forja"),
		filepath.Join(raiz, "internal", "ste"),
		filepath.Join(raiz, "internal", "instalar"),
	}
	for _, dir := range pacotes {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}
		verificarArvoreSemPython(t, dir)
	}
	verificarComandosInstalacaoSemPython(t, filepath.Join(raiz, "internal", "cli"))
}

func TestDetectarPythonEmProducao(t *testing.T) {
	dir := t.TempDir()
	arqProd := filepath.Join(dir, "prod.go")
	arqTeste := filepath.Join(dir, "prod_test.go")
	if err := os.WriteFile(arqProd, []byte("package x\nconst c = \"Python3\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(arqTeste, []byte("package x\nconst c = \"python\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	achados, err := buscarPythonEmDir(dir)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(achados) != 1 || achados[0] != arqProd {
		t.Fatalf("esperava apenas %q, obteve %v", arqProd, achados)
	}
}

func localizarRaizModulo(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("falha ao obter diretório atual: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		pai := filepath.Dir(dir)
		if pai == dir {
			t.Fatal("go.mod não encontrado acima de internal/forja")
		}
		dir = pai
	}
}

func verificarArvoreSemPython(t *testing.T, dir string) {
	t.Helper()
	achados, err := buscarPythonEmDir(dir)
	if err != nil {
		t.Fatalf("falha ao inspecionar %s: %v", dir, err)
	}
	for _, arq := range achados {
		t.Errorf("referência proibida a python em arquivo de produção: %s", arq)
	}
}

func verificarComandosInstalacaoSemPython(t *testing.T, dirCLI string) {
	t.Helper()
	prefixos := []string{"instalar", "forjar", "atualizar"}
	entradas, err := os.ReadDir(dirCLI)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		t.Fatalf("falha ao listar %s: %v", dirCLI, err)
	}
	for _, e := range entradas {
		if e.IsDir() || !ehGoProducao(e.Name()) {
			continue
		}
		for _, pref := range prefixos {
			if strings.HasPrefix(e.Name(), pref) {
				arq := filepath.Join(dirCLI, e.Name())
				contem, err := arquivoContemPython(arq)
				if err != nil {
					t.Fatalf("falha ao ler %s: %v", arq, err)
				}
				if contem {
					t.Errorf("referência proibida a python no comando %s", arq)
				}
			}
		}
	}
}

func buscarPythonEmDir(dir string) ([]string, error) {
	var achados []string
	err := filepath.WalkDir(dir, func(caminho string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !ehGoProducao(d.Name()) {
			return nil
		}
		contem, err := arquivoContemPython(caminho)
		if err != nil {
			return err
		}
		if contem {
			achados = append(achados, caminho)
		}
		return nil
	})
	return achados, err
}

func ehGoProducao(nome string) bool {
	return strings.HasSuffix(nome, ".go") && !strings.HasSuffix(nome, "_test.go")
}

func arquivoContemPython(caminho string) (bool, error) {
	dados, err := os.ReadFile(caminho)
	if err != nil {
		return false, err
	}
	return strings.Contains(strings.ToLower(string(dados)), "python"), nil
}

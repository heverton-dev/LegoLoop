package dag

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDAGsDoRepositorioValidam(t *testing.T) {
	raiz := filepath.Join("..", "..")
	arquivos := coletarDAGsDoRepositorio(t, raiz)
	if len(arquivos) == 0 {
		t.Fatal("nenhum arquivo DAG encontrado em processos/ ou docs/dags/")
	}
	for _, arq := range arquivos {
		rel, err := filepath.Rel(raiz, arq)
		if err != nil {
			rel = arq
		}
		t.Run(filepath.ToSlash(rel), func(t *testing.T) {
			validarDAGDoRepositorio(t, arq)
		})
	}
}

func coletarDAGsDoRepositorio(t *testing.T, raiz string) []string {
	t.Helper()
	var lista []string
	coletarPorFiltro(t, filepath.Join(raiz, "processos"), func(nome string) bool {
		return strings.HasSuffix(nome, ".dag.json")
	}, &lista)
	coletarPorFiltro(t, filepath.Join(raiz, "docs", "dags"), func(nome string) bool {
		return strings.HasSuffix(nome, ".json")
	}, &lista)
	return lista
}

func coletarPorFiltro(t *testing.T, dir string, aceita func(string) bool, destino *[]string) {
	t.Helper()
	err := filepath.WalkDir(dir, func(caminho string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && aceita(d.Name()) {
			*destino = append(*destino, caminho)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("falha ao percorrer %s: %v", dir, err)
	}
}

func validarDAGDoRepositorio(t *testing.T, caminho string) {
	t.Helper()
	dados, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatalf("erro ao ler %s: %v", caminho, err)
	}
	d, err := Parse(dados)
	if err != nil {
		t.Fatalf("Parse falhou em %s: %v", caminho, err)
	}
	if err := d.ValidarExecutores(); err != nil {
		t.Errorf("ValidarExecutores falhou em %s: %v", caminho, err)
	}
	for _, n := range d.Nodes {
		if n.Type != TypeStopGate {
			continue
		}
		exec, _ := n.Isolation["executor"].(string)
		if exec != ExecutorLocalProcess {
			t.Errorf("nó stop_gate %q em %s tem isolation.executor=%q, esperado %q", n.ID, caminho, exec, ExecutorLocalProcess)
		}
	}
}

package compor_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/caminhos"
	"github.com/heverton-dev/LegoLoop/internal/triagem"
)

// canonizar move o DAG e o Studio de um processo para o layout do fabricar: uma pasta por processo no cofre.
func canonizar(t *testing.T, c cofre, nome string) (dag, studio string) {
	t.Helper()
	base := pastaDoCofre(c.vault)
	pasta := filepath.Join(base, nome)
	if err := os.MkdirAll(pasta, 0o755); err != nil {
		t.Fatal(err)
	}
	dag, studio = caminhos.CaminhoDAG(pasta), filepath.Join(pasta, nome+caminhos.SufixoStudio)
	for de, para := range map[string]string{filepath.Join(base, nome+".dag.json"): dag, filepath.Join(base, nome+caminhos.SufixoStudio): studio} {
		if err := os.Rename(de, para); err != nil {
			t.Fatal(err)
		}
	}
	return dag, studio
}

func TestCompor_ComLayoutCanonicoAtualizaAPastaEGeraATriagem(t *testing.T) {
	c := semear(t, nomeA, nomeB)
	dag, studio := canonizar(t, c, nomeA)
	if _, err := ligar(c, nomeA, nomeB+"@1.0.0"); err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{dag, studio} {
		b, err := os.ReadFile(f)
		if err != nil || !strings.Contains(string(b), "e04") {
			t.Errorf("%s deve ganhar a etapa e04: %v", filepath.Base(f), err)
		}
	}
	base := pastaDoCofre(c.vault)
	for _, solto := range []string{nomeA + ".dag.json", nomeA + caminhos.SufixoStudio} {
		if _, err := os.Stat(filepath.Join(base, solto)); err == nil {
			t.Errorf("não deve criar %s solto ao lado da pasta", solto)
		}
	}
	ids := []string{"e01", "e02", "e03", "e04"}
	if err := triagem.Conferir(filepath.Dir(dag), ids); err != nil {
		t.Errorf("com 4 nós a triagem deve existir e cobrir todos: %v", err)
	}
}

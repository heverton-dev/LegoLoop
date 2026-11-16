package refinar_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/caminhos"
	"legoloop/internal/paths"
	"legoloop/internal/refinar"
)

// canonizar move o DAG e o documento do Studio para o layout que o fabricar grava: pastaCofre/<nome>/.
func canonizar(t *testing.T, c cofre) (dag, studio string) {
	t.Helper()
	pastaCofre := filepath.Join(c.vault, paths.DirAntigo)
	pasta := filepath.Join(pastaCofre, nome)
	if err := os.MkdirAll(pasta, 0o755); err != nil {
		t.Fatal(err)
	}
	dag, studio = filepath.Join(pasta, caminhos.ArquivoDAG), filepath.Join(pasta, nome+caminhos.SufixoStudio)
	for de, para := range map[string]string{filepath.Join(pastaCofre, nome+".dag.json"): dag, filepath.Join(pastaCofre, nome+paths.SufixoAntigo): studio} {
		if err := os.Rename(de, para); err != nil {
			t.Fatal(err)
		}
	}
	return dag, studio
}

func TestRefinar_ComLayoutCanonicoAtualizaOMesmoDAG(t *testing.T) {
	c := semear(t, processoDeTeste(), inicio, true)
	dag, studio := canonizar(t, c)
	m := refinar.Mudanca{Etapa: "e04", Peca: "legoloop-nova", Casos: []string{"a"}, Adicionar: true}
	if _, err := refinar.Run(context.Background(), refinar.Request{Nome: nome, Mudanca: m}, depsDe(c, nil, seguir, inicio)); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(ler(t, dag), "e04") || !strings.Contains(ler(t, studio), "e04") {
		t.Error("o DAG e o Studio da pasta do processo devem ganhar a etapa e04")
	}
	pastaCofre := filepath.Join(c.vault, paths.DirAntigo)
	for _, solto := range []string{nome + ".dag.json", nome + paths.SufixoAntigo} {
		if _, err := os.Stat(filepath.Join(pastaCofre, solto)); err == nil {
			t.Errorf("não deve criar %s solto ao lado da pasta do processo", solto)
		}
	}
}

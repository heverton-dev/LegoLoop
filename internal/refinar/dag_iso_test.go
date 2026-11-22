package refinar_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/processo"
	"github.com/heverton-dev/LegoLoop/internal/refinar"
)

func TestRefinar_AcompanhaODagEspelhoQueAProjetoJaTem(t *testing.T) {
	c := semear(t, processoDeTeste(), inicio, true)
	pr, err := processo.Projetar(processoDeTeste(), inicio)
	if err != nil {
		t.Fatal(err)
	}
	iso := filepath.Join(c.docs, "docs", "dags", pr.DAGNome)
	if err := os.MkdirAll(filepath.Dir(iso), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(iso, pr.DAG, 0o644); err != nil {
		t.Fatal(err)
	}
	m := refinar.Mudanca{Etapa: "e04", Peca: "legoloop-nova", Casos: []string{"a"}, Adicionar: true}
	if _, err := refinar.Run(context.Background(), refinar.Request{Nome: nome, Mudanca: m}, depsDe(c, nil, seguir, inicio.AddDate(0, 0, 3))); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(ler(t, iso), "e04") {
		t.Error("o DAG espelho deve ganhar a etapa e04")
	}
	skill := ler(t, filepath.Join(c.vault, "skills", nome, "SKILL.md"))
	if !strings.Contains(skill, "docs/dags/"+pr.DAGNome) {
		t.Errorf("a skill deve apontar o espelho %s", pr.DAGNome)
	}
}

func TestRefinar_SemDagEspelhoNaoCriaUm(t *testing.T) {
	c := semear(t, processoDeTeste(), inicio, true)
	m := refinar.Mudanca{Etapa: "e04", Peca: "legoloop-nova", Casos: []string{"a"}, Adicionar: true}
	if _, err := refinar.Run(context.Background(), refinar.Request{Nome: nome, Mudanca: m}, depsDe(c, nil, seguir, inicio)); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(c.docs, "docs", "dags")); err == nil {
		t.Error("sem espelho no projeto, o refinar não cria docs/dags")
	}
}

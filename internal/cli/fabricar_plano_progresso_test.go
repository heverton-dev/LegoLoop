package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestFabricarPlanoApontaProgresso cobre o achado A19: o DAG gerado aponta o progresso.json do PLAN,
// e o arquivo apontado existe. Sem o ponteiro o swarm run nao atualizava o progresso vivo.
func TestFabricarPlanoApontaProgresso(t *testing.T) {
	raiz, spec, ticketsFile := prepararAmbienteFabricarPlano(t, ticketsDoFabricarPlano)
	args := []string{"plano", "--raiz", raiz, "--spec", spec, "--tickets", ticketsFile, "--data", "2026-11-21"}
	if code, saida := rodarFabricar(args...); code != 0 {
		t.Fatalf("fabricar plano falhou (%d): %s", code, saida)
	}
	raw, err := os.ReadFile(filepath.Join(raiz, "docs", "dags", "2026-11-21_DAG_TESTE.json"))
	if err != nil {
		t.Fatal(err)
	}
	var d struct {
		Progresso string `json:"progresso"`
	}
	if err := json.Unmarshal(raw, &d); err != nil {
		t.Fatal(err)
	}
	want := "docs/progress/2026-11-21_plan-teste/progresso.json"
	if d.Progresso != want {
		t.Fatalf("progresso = %q; quer %q", d.Progresso, want)
	}
	if _, err := os.Stat(filepath.Join(raiz, filepath.FromSlash(d.Progresso))); err != nil {
		t.Fatalf("progresso apontado nao existe: %v", err)
	}
}

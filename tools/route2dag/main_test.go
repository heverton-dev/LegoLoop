package main

import (
	"encoding/json"
	"errors"
	"os"
	"slices"
	"strings"
	"testing"

	"legoloop/internal/dag"
)

// TestProjecaoValidaNoParser (AC-E01.6): a projeção da ROUTE desta demanda passa no parser,
// mantém todos os nós, tem artefato único por nó e o command de cada tarefa é o primeiro
// argv de verify.
func TestProjecaoValidaNoParser(t *testing.T) {
	raw, err := os.ReadFile(rotaLei15)
	if err != nil {
		t.Fatal(err)
	}
	out, nome, err := Projetar(raw)
	if err != nil {
		t.Fatalf("projeção falhou: %v", err)
	}
	if nome != "2026-11-10_DAG_NORMALIZACAO_TAXONOMICA_LEI15.json" {
		t.Errorf("nome = %q", nome)
	}
	d, err := dag.Parse(out)
	if err != nil {
		t.Fatalf("DAG projetado não valida: %v", err)
	}
	var r route
	if err := json.Unmarshal(raw, &r); err != nil {
		t.Fatal(err)
	}
	if len(d.Nodes) != len(r.Nodes) {
		t.Fatalf("nós = %d, esperado %d", len(d.Nodes), len(r.Nodes))
	}
	conferirNos(t, d, r)
}

func conferirNos(t *testing.T, d *dag.DAG, r route) {
	t.Helper()
	vistos := map[string]bool{}
	for i, n := range d.Nodes {
		rn := r.Nodes[i]
		if n.ID != rn.ID {
			t.Errorf("nó %d: id = %q, esperado %q", i, n.ID, rn.ID)
		}
		if strings.HasPrefix(n.ID, "TICK-") && !strings.Contains(n.ID, ".") {
			t.Errorf("%s: id de ticket sem ponto", n.ID)
		}
		if vistos[n.Output.Artifact] {
			t.Errorf("%s: artefato repetido %q", n.ID, n.Output.Artifact)
		}
		vistos[n.Output.Artifact] = true
		var want []string
		if rn.Type != "stop_gate" && len(rn.Output.Verify) > 0 {
			want = rn.Output.Verify[0]
		}
		if !slices.Equal(n.Command, want) {
			t.Errorf("%s: command = %v, esperado %v", n.ID, n.Command, want)
		}
	}
}

// Mutação: dois nós com o mesmo artefato (ids iguais sem caixa) fazem a projeção falhar.
func TestProjecaoRecusaArtefatoRepetido(t *testing.T) {
	rota := `{"name":"X","version":"1","date":"2026-01-01","nodes":[
	 {"id":"T1","type":"unit_task","nature":"micro","routing":{"mode":"manual","model":"sonnet"},"output":{"gate":"exit_code == 0"}},
	 {"id":"t1","type":"unit_task","nature":"micro","routing":{"mode":"manual","model":"sonnet"},"output":{"gate":"exit_code == 0"}}]}`
	_, _, err := Projetar([]byte(rota))
	if !errors.Is(err, ErrArtefatoRepetido) {
		t.Fatalf("err = %v, esperado ErrArtefatoRepetido", err)
	}
}

// Mutação: ROUTE com ciclo gera DAG que o parser recusa, e a projeção falha.
func TestProjecaoRecusaCiclo(t *testing.T) {
	rota := `{"name":"X","version":"1","date":"2026-01-01","nodes":[
	 {"id":"A","type":"unit_task","nature":"micro","depends_on":["B"],"routing":{"mode":"manual","model":"sonnet"},"output":{"gate":"exit_code == 0"}},
	 {"id":"B","type":"unit_task","nature":"micro","depends_on":["A"],"routing":{"mode":"manual","model":"sonnet"},"output":{"gate":"exit_code == 0"}}]}`
	if _, _, err := Projetar([]byte(rota)); err == nil {
		t.Fatal("ROUTE com ciclo deveria falhar na projeção")
	}
}

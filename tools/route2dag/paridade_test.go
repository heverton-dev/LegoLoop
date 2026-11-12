package main

import (
	"encoding/json"
	"os"
	"testing"

	"legoloop/internal/dag"
)

func carregarRota(t *testing.T) route {
	t.Helper()
	raw, err := os.ReadFile(rotaLei15)
	if err != nil {
		t.Fatal(err)
	}
	var r route
	if err := json.Unmarshal(raw, &r); err != nil {
		t.Fatal(err)
	}
	return r
}

func indicesUnidade(r route) []int {
	var idx []int
	for i, n := range r.Nodes {
		if n.Type != "stop_gate" {
			idx = append(idx, i)
		}
	}
	return idx
}

func TestParidadeViaEHarness(t *testing.T) {
	r := carregarRota(t)
	idx := indicesUnidade(r)
	if len(idx) < 2 {
		t.Fatal("ROUTE sem dois nós executáveis")
	}
	r.Nodes[idx[0]].Routing.Harness = "direto"
	r.Nodes[idx[0]].Routing.Via = "9router"
	r.Nodes[idx[1]].Routing.Harness = "claude"
	r.Nodes[idx[1]].Routing.Via = "oficial"
	raw, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	out, _, err := Projetar(raw)
	if err != nil {
		t.Fatal(err)
	}
	d, err := dag.Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range idx {
		n := r.Nodes[i]
		got := d.Node(n.ID).Routing
		if got.Harness != n.Routing.Harness || got.Via != n.Routing.Via || got.Model != n.Routing.Model {
			t.Errorf("%s: DAG %q/%q/%q diverge da ROUTE %q/%q/%q", n.ID,
				got.Harness, got.Via, got.Model, n.Routing.Harness, n.Routing.Via, n.Routing.Model)
		}
		if n.Routing.Effort != nil && string(got.Effort) != *n.Routing.Effort {
			t.Errorf("%s: effort DAG %q diverge da ROUTE %q", n.ID, got.Effort, *n.Routing.Effort)
		}
	}
}

func TestParidadeDiretoComViaOficialFalha(t *testing.T) {
	r := carregarRota(t)
	i := indicesUnidade(r)[0]
	r.Nodes[i].Routing.Harness = "direto"
	r.Nodes[i].Routing.Via = "oficial"
	raw, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := Projetar(raw); err == nil {
		t.Fatal("nó direto com via oficial deveria falhar")
	}
}

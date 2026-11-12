package main

import (
	"encoding/json"
	"os"
	"testing"

	"legoloop/internal/dag"
)

func TestParidadeRouteDag(t *testing.T) {
	raw, err := os.ReadFile(rotaLei15)
	if err != nil {
		t.Fatal(err)
	}
	var r route
	if err := json.Unmarshal(raw, &r); err != nil {
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
	comHarness := 0
	for _, n := range r.Nodes {
		if n.Type == "stop_gate" {
			continue
		}
		got := d.Node(n.ID).Routing
		if got.Harness != n.Routing.Harness || got.Model != n.Routing.Model {
			t.Errorf("%s: DAG %q/%q diverge da ROUTE %q/%q", n.ID, got.Harness, got.Model, n.Routing.Harness, n.Routing.Model)
		}
		if n.Routing.Effort != nil && string(got.Effort) != *n.Routing.Effort {
			t.Errorf("%s: effort DAG %q diverge da ROUTE %q", n.ID, got.Effort, *n.Routing.Effort)
		}
		if n.Routing.Harness != "" {
			comHarness++
		}
	}
	if comHarness == 0 {
		t.Fatal("ROUTE sem nenhum harness resolvido")
	}
}

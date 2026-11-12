package main

import "testing"

func TestStopGateLevaSpecRef(t *testing.T) {
	ref := "AC-I3-05.2"
	nd, err := projetarNo(routeNode{ID: "SG-X", Type: "stop_gate", SpecRef: &ref})
	if err != nil {
		t.Fatal(err)
	}
	if nd["spec_ref"] != ref {
		t.Errorf("spec_ref do Stop-Gate = %v, esperava %s", nd["spec_ref"], ref)
	}
	semRef, err := projetarNo(routeNode{ID: "SG-Y", Type: "stop_gate"})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := semRef["spec_ref"]; ok {
		t.Error("Stop-Gate sem spec_ref não deve gravar o campo")
	}
}

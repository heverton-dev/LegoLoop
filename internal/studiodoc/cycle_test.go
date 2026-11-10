package studiodoc_test

import (
	"fmt"
	"strings"
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/studiodoc"
)

func makeSubloopNode(id, loopRef string) dag.Node {
	return dag.Node{
		ID:   id,
		Type: dag.TypeSubloop,
		Ref:  loopRef,
		Routing: dag.RoutingConfig{
			Mode:  dag.ModeManual,
			Model: "subloop",
		},
		Output: dag.HandoffContract{
			Artifact: fmt.Sprintf("%s_OUT.md", id),
			Gate:     "true",
		},
	}
}

func TestCheckSubloopsNoSubloops(t *testing.T) {
	dir := t.TempDir()
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	lp := sampleValidLoop("standalone", "Standalone")
	if err := store.Save(lp); err != nil {
		t.Fatal(err)
	}

	if err := store.ValidateHierarchy(lp); err != nil {
		t.Fatalf("esperava sucesso para loop sem sub-loops, obteve: %v", err)
	}
}

func TestCheckSubloopsLinearValid(t *testing.T) {
	dir := t.TempDir()
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	// C é folha
	c := sampleValidLoop("loop-c", "Loop C")
	if err := store.Save(c); err != nil {
		t.Fatal(err)
	}

	// B chama C
	b := &studiodoc.Loop{
		ID:      "loop-b",
		Name:    "Loop B",
		Version: "2.0.0",
		Trigger: &dag.Trigger{Type: dag.TriggerManual},
		Nodes:   []dag.Node{makeSubloopNode("nb", "loop-c")},
	}
	if err := store.Save(b); err != nil {
		t.Fatal(err)
	}

	// A chama B
	a := &studiodoc.Loop{
		ID:      "loop-a",
		Name:    "Loop A",
		Version: "2.0.0",
		Trigger: &dag.Trigger{Type: dag.TriggerManual},
		Nodes:   []dag.Node{makeSubloopNode("na", "loop-b")},
	}
	if err := store.Save(a); err != nil {
		t.Fatal(err)
	}

	if err := store.ValidateHierarchy(a); err != nil {
		t.Fatalf("esperava sucesso para cadeia A -> B -> C, obteve: %v", err)
	}
}

func TestCheckSubloopsDirectSelfCycle(t *testing.T) {
	dir := t.TempDir()
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	a := &studiodoc.Loop{
		ID:      "loop-a",
		Name:    "Loop A",
		Version: "2.0.0",
		Trigger: &dag.Trigger{Type: dag.TriggerManual},
		Nodes:   []dag.Node{makeSubloopNode("na", "loop-a")},
	}
	if err := store.Save(a); err != nil {
		t.Fatal(err)
	}

	err = store.ValidateHierarchy(a)
	if err == nil {
		t.Fatal("esperava erro de ciclo direto A -> A")
	}
	if !strings.Contains(err.Error(), "ciclo") {
		t.Fatalf("mensagem de erro deveria conter 'ciclo', obteve: %v", err)
	}
}

func TestCheckSubloopsMutualCycle(t *testing.T) {
	dir := t.TempDir()
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	// A chama B
	a := &studiodoc.Loop{
		ID:      "loop-a",
		Name:    "Loop A",
		Version: "2.0.0",
		Trigger: &dag.Trigger{Type: dag.TriggerManual},
		Nodes:   []dag.Node{makeSubloopNode("na", "loop-b")},
	}
	// B chama A
	b := &studiodoc.Loop{
		ID:      "loop-b",
		Name:    "Loop B",
		Version: "2.0.0",
		Trigger: &dag.Trigger{Type: dag.TriggerManual},
		Nodes:   []dag.Node{makeSubloopNode("nb", "loop-a")},
	}

	if err := store.Save(a); err != nil {
		t.Fatal(err)
	}
	if err := store.Save(b); err != nil {
		t.Fatal(err)
	}

	err = store.ValidateHierarchy(a)
	if err == nil {
		t.Fatal("esperava erro de ciclo mútuo A -> B -> A")
	}
	if !strings.Contains(err.Error(), "ciclo") {
		t.Fatalf("mensagem esperada com 'ciclo', obteve: %v", err)
	}
}

func TestCheckSubloopsMaxDepthExceeded(t *testing.T) {
	dir := t.TempDir()
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	// Cria cadeia de profundidade 9: l1 -> l2 -> ... -> l9 -> leaf
	// Max allowed depth is 8
	leaf := sampleValidLoop("loop-leaf", "Leaf")
	if err := store.Save(leaf); err != nil {
		t.Fatal(err)
	}

	prev := "loop-leaf"
	for i := 9; i >= 1; i-- {
		id := fmt.Sprintf("loop-%d", i)
		l := &studiodoc.Loop{
			ID:      id,
			Name:    fmt.Sprintf("Loop %d", i),
			Version: "2.0.0",
			Trigger: &dag.Trigger{Type: dag.TriggerManual},
			Nodes:   []dag.Node{makeSubloopNode("n", prev)},
		}
		if err := store.Save(l); err != nil {
			t.Fatal(err)
		}
		prev = id
	}

	root, err := store.Get("loop-1")
	if err != nil {
		t.Fatal(err)
	}

	err = store.ValidateHierarchy(root)
	if err == nil {
		t.Fatal("esperava erro de profundidade máxima excedida")
	}
	if !strings.Contains(err.Error(), "profundidade") {
		t.Fatalf("mensagem esperada com 'profundidade', obteve: %v", err)
	}
}

func TestCheckSubloopsMissingRef(t *testing.T) {
	dir := t.TempDir()
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	a := &studiodoc.Loop{
		ID:      "loop-a",
		Name:    "Loop A",
		Version: "2.0.0",
		Trigger: &dag.Trigger{Type: dag.TriggerManual},
		Nodes:   []dag.Node{makeSubloopNode("na", "fantasma")},
	}
	if err := store.Save(a); err != nil {
		t.Fatal(err)
	}

	err = store.ValidateHierarchy(a)
	if err == nil {
		t.Fatal("esperava erro para subloop inexistente")
	}
	if !strings.Contains(err.Error(), "não encontrado") {
		t.Fatalf("mensagem esperada com 'não encontrado', obteve: %v", err)
	}
}

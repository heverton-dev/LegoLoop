package studiodoc_test

import (
	"fmt"
	"strings"
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/studiodoc"
)

func makeSubprocessoNode(id, ref string) dag.Node {
	return dag.Node{
		ID:   id,
		Type: dag.TypeSubprocesso,
		Ref:  ref,
		Routing: dag.RoutingConfig{
			Mode:  dag.ModeManual,
			Model: "subprocesso",
		},
		Output: dag.HandoffContract{
			Artifact: fmt.Sprintf("%s_OUT.md", id),
			Gate:     "true",
		},
	}
}

func TestCheckSubprocessosNoSubprocessos(t *testing.T) {
	dir := t.TempDir()
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	lp := sampleValidProcesso("standalone", "Standalone")
	if err := store.Save(lp); err != nil {
		t.Fatal(err)
	}

	if err := store.ValidateHierarchy(lp); err != nil {
		t.Fatalf("esperava sucesso para Processo sem subprocessos, obteve: %v", err)
	}
}

func TestCheckSubprocessosLinearValid(t *testing.T) {
	dir := t.TempDir()
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	c := sampleValidProcesso("proc-c", "Processo C")
	if err := store.Save(c); err != nil {
		t.Fatal(err)
	}

	b := &studiodoc.Processo{
		ID:      "proc-b",
		Name:    "Processo B",
		Version: "2.0.0",
		Trigger: &dag.Trigger{Type: dag.TriggerManual},
		Nodes:   []dag.Node{makeSubprocessoNode("nb", "proc-c")},
	}
	if err := store.Save(b); err != nil {
		t.Fatal(err)
	}

	a := &studiodoc.Processo{
		ID:      "proc-a",
		Name:    "Processo A",
		Version: "2.0.0",
		Trigger: &dag.Trigger{Type: dag.TriggerManual},
		Nodes:   []dag.Node{makeSubprocessoNode("na", "proc-b")},
	}
	if err := store.Save(a); err != nil {
		t.Fatal(err)
	}

	if err := store.ValidateHierarchy(a); err != nil {
		t.Fatalf("esperava sucesso para cadeia A -> B -> C, obteve: %v", err)
	}
}

func TestCheckSubprocessosDirectSelfCycle(t *testing.T) {
	dir := t.TempDir()
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	a := &studiodoc.Processo{
		ID:      "proc-a",
		Name:    "Processo A",
		Version: "2.0.0",
		Trigger: &dag.Trigger{Type: dag.TriggerManual},
		Nodes:   []dag.Node{makeSubprocessoNode("na", "proc-a")},
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

func TestCheckSubprocessosMutualCycle(t *testing.T) {
	dir := t.TempDir()
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	a := &studiodoc.Processo{
		ID:      "proc-a",
		Name:    "Processo A",
		Version: "2.0.0",
		Trigger: &dag.Trigger{Type: dag.TriggerManual},
		Nodes:   []dag.Node{makeSubprocessoNode("na", "proc-b")},
	}
	b := &studiodoc.Processo{
		ID:      "proc-b",
		Name:    "Processo B",
		Version: "2.0.0",
		Trigger: &dag.Trigger{Type: dag.TriggerManual},
		Nodes:   []dag.Node{makeSubprocessoNode("nb", "proc-a")},
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

func TestCheckSubprocessosMaxDepthExceeded(t *testing.T) {
	dir := t.TempDir()
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	leaf := sampleValidProcesso("proc-leaf", "Leaf")
	if err := store.Save(leaf); err != nil {
		t.Fatal(err)
	}

	prev := "proc-leaf"
	for i := 9; i >= 1; i-- {
		id := fmt.Sprintf("proc-%d", i)
		l := &studiodoc.Processo{
			ID:      id,
			Name:    fmt.Sprintf("Processo %d", i),
			Version: "2.0.0",
			Trigger: &dag.Trigger{Type: dag.TriggerManual},
			Nodes:   []dag.Node{makeSubprocessoNode("n", prev)},
		}
		if err := store.Save(l); err != nil {
			t.Fatal(err)
		}
		prev = id
	}

	root, err := store.Get("proc-1")
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

func TestCheckSubprocessosMissingRef(t *testing.T) {
	dir := t.TempDir()
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	a := &studiodoc.Processo{
		ID:      "proc-a",
		Name:    "Processo A",
		Version: "2.0.0",
		Trigger: &dag.Trigger{Type: dag.TriggerManual},
		Nodes:   []dag.Node{makeSubprocessoNode("na", "fantasma")},
	}
	if err := store.Save(a); err != nil {
		t.Fatal(err)
	}

	err = store.ValidateHierarchy(a)
	if err == nil {
		t.Fatal("esperava erro para subprocesso inexistente")
	}
	if !strings.Contains(err.Error(), "não encontrado") {
		t.Fatalf("mensagem esperada com 'não encontrado', obteve: %v", err)
	}
}

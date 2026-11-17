package runner

import (
	"context"
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/triagem"
)

func TestPreExecNodeGravaEtapaNaMemoria(t *testing.T) {
	dir := t.TempDir()
	d := &dag.DAG{Nodes: []dag.Node{
		{ID: "n1"},
		{ID: "n2", DependsOn: []string{"n1"}},
	}}
	r := &run{d: d, opts: Options{Workdir: dir, SessaoID: "sess-teste"}}
	n1 := d.Nodes[0]
	if err := r.preExecNode(context.Background(), &n1, t.TempDir()); err != nil {
		t.Fatalf("preExecNode: %v", err)
	}
	mem, err := triagem.LerMemoria(dir, "sess-teste")
	if err != nil {
		t.Fatalf("LerMemoria: %v", err)
	}
	if mem.EtapaAtual != "n1" {
		t.Errorf("etapa atual = %q, esperado \"n1\"", mem.EtapaAtual)
	}
	if mem.ProximoPasso != "n2" {
		t.Errorf("próximo passo = %q, esperado \"n2\"", mem.ProximoPasso)
	}
}

func TestProximosDeAchaSucessores(t *testing.T) {
	d := &dag.DAG{Nodes: []dag.Node{
		{ID: "a"},
		{ID: "b", DependsOn: []string{"a"}},
		{ID: "c", DependsOn: []string{"a", "b"}},
	}}
	got := proximosDe(d, "a")
	if len(got) != 2 || got[0] != "b" || got[1] != "c" {
		t.Errorf("proximosDe(a) = %v, esperado [b c]", got)
	}
	if fim := proximosDe(d, "c"); len(fim) != 0 {
		t.Errorf("proximosDe(c) = %v, esperado vazio", fim)
	}
}

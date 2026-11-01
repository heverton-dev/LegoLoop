package runner

import (
	"context"
	"strings"
	"testing"

	"legoloop/internal/dag"
)

func guardDAG(t *testing.T, extra string) *dag.DAG {
	t.Helper()
	src := `{"name":"g","version":"2.0.0","trigger":{"type":"manual"},"nodes":[{"id":"a","type":"script",
	"command":["x"],"routing":{"mode":"manual","model":"none"},
	"output":{"artifact":"a.md","gate":"exit_code == 0"}` + extra + `}]}`
	d, err := dag.Parse([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	return d
}

// Recursos do contrato v2 que o runner ainda não executa NUNCA são ignorados em silêncio:
// um hook ou gate de slot descartado faria o usuário crer que a proteção rodou.
func TestRunRefusesUnsupportedV2Features(t *testing.T) {
	cases := map[string]string{
		"slots": `,"slots":{"hooks":{"pre":[["go","version"]]}}`,
	}
	for name, extra := range cases {
		ex := &fakeExec{}
		_, err := Run(context.Background(), guardDAG(t, extra), Options{Workdir: t.TempDir(), Executor: ex})
		if err == nil || !strings.Contains(err.Error(), "a") || !strings.Contains(err.Error(), "ainda não") {
			t.Errorf("%s: esperado recusa explícita, veio %v", name, err)
		}
		if len(ex.order) != 0 {
			t.Errorf("%s: nada deve executar quando o DAG é recusado", name)
		}
	}
	sub, err := dag.Parse([]byte(`{"name":"g","version":"2.0.0","trigger":{"type":"manual"},"nodes":[{"id":"s",
	"type":"subloop","ref":"outro-loop","routing":{"mode":"manual","model":"none"},
	"output":{"artifact":"s.md","gate":"true"}}]}`))
	if err != nil {
		t.Fatal(err)
	}
	ex := &fakeExec{}
	if _, err := Run(context.Background(), sub, Options{Workdir: t.TempDir(), Executor: ex}); err != nil {
		t.Fatalf("subloop deve executar normalmente via executor: %v", err)
	}
	if len(ex.order) != 1 || ex.order[0] != "s" {
		t.Fatalf("subloop deveria ter sido despachado para o executor, ordem: %v", ex.order)
	}
	// inline explícito e v1 continuam executando
	if _, err := Run(context.Background(), guardDAG(t, `,"execution":{"mode":"inline"}`),
		Options{Workdir: t.TempDir(), Executor: &fakeExec{}}); err != nil {
		t.Fatalf("inline deveria executar: %v", err)
	}
}

// Slots vazios são escolha válida do usuário: o nó executa normalmente, sem recusa.
func TestRunEmptySlotsExecute(t *testing.T) {
	ex := &fakeExec{}
	empty := `,"slots":{"gates":[],"scripts":[],"hooks":{"pre":[],"post":[]}}`
	rep, err := Run(context.Background(), guardDAG(t, empty), Options{Workdir: t.TempDir(), Executor: ex})
	if err != nil || !rep.OK || len(ex.order) != 1 {
		t.Fatalf("slots vazios deveriam executar: err=%v ok=%v ran=%v", err, rep.OK, ex.order)
	}
}

package runner

import (
	"context"
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

// Recursos do contrato v2 agora são todos executados pelo runner: slots, isolated e subprocesso (TICK-1107).
func TestRunSupportsAllV2Features(t *testing.T) {
	cases := map[string]string{
		"slots":    `,"slots":{"gates":["exit_code == 0"]}`,
		"isolated": `,"execution":{"mode":"isolated"}`,
	}
	for name, extra := range cases {
		ex := &fakeExec{}
		rep, err := Run(context.Background(), guardDAG(t, extra), Options{Workdir: t.TempDir(), Executor: ex, AllowScripts: true})
		if err != nil || !rep.OK {
			t.Errorf("%s: esperado sucesso de execução v2, veio err=%v rep=%+v", name, err, rep)
		}
		if len(ex.order) != 1 {
			t.Errorf("%s: o nó deve executar", name)
		}
	}
	sub, err := dag.Parse([]byte(`{"name":"g","version":"2.0.0","trigger":{"type":"manual"},"nodes":[{"id":"s",
	"type":"subprocesso","ref":"outro-processo","routing":{"mode":"manual","model":"none"},
	"output":{"artifact":"s.md","gate":"true"}}]}`))
	if err != nil {
		t.Fatal(err)
	}
	ex := &fakeExec{}
	if _, err := Run(context.Background(), sub, Options{Workdir: t.TempDir(), Executor: ex}); err != nil {
		t.Fatalf("subprocesso deve executar normalmente via executor: %v", err)
	}
	if len(ex.order) != 1 || ex.order[0] != "s" {
		t.Fatalf("subprocesso deveria ter sido despachado para o executor, ordem: %v", ex.order)
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

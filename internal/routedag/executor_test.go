package routedag_test

import (
	"errors"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/dag"
	"github.com/heverton-dev/LegoLoop/internal/routedag"
)

func routeValidaMultiExecutores() string {
	return `{
  "name": "2026-11-23_ROUTE_TESTE_EXECUTORES",
  "version": "1.0.0",
  "date": "2026-11-23",
  "execution_strategy": {
    "supervision_level": "epic"
  },
  "nodes": [
    {
      "id": "SG-01",
      "nature": "meta",
      "type": "stop_gate",
      "status": "PENDING",
      "isolation": {"executor": "local_process"},
      "output": {"gate": "prompt('ok?')"}
    },
    {
      "id": "TICK-01",
      "nature": "micro",
      "type": "unit_no",
      "status": "PENDING",
      "depends_on": ["SG-01"],
      "isolation": {"executor": "native_swarm", "write_boundary": "strict"},
      "routing": {"mode": "suggested", "harness": "9router", "model": "code-pro"},
      "output": {
        "artifacts": ["out/t1.txt"],
        "verify": [["go", "test", "./..."]],
        "gate": "exit_code == 0"
      }
    },
    {
      "id": "TICK-02",
      "nature": "micro",
      "type": "unit_no",
      "status": "PENDING",
      "depends_on": ["TICK-01"],
      "isolation": {"executor": "orchestrix_hermetic", "write_boundary": "strict"},
      "routing": {"mode": "suggested", "harness": "9router", "model": "code-pro"},
      "output": {
        "artifacts": ["out/t2.txt"],
        "verify": [["go", "test", "./..."]],
        "gate": "exit_code == 0"
      }
    }
  ]
}`
}

func TestProjecaoPreservaExecutoresEDagParseVerde(t *testing.T) {
	dagBytes, _, err := routedag.Converter([]byte(routeValidaMultiExecutores()))
	if err != nil {
		t.Fatalf("Converter falhou: %v", err)
	}

	d, err := dag.Parse(dagBytes)
	if err != nil {
		t.Fatalf("dag.Parse falhou: %v\nDAG:\n%s", err, string(dagBytes))
	}

	nSG := d.Node("SG-01")
	if nSG == nil {
		t.Fatal("nó SG-01 não encontrado")
	}
	if exec, _ := nSG.Isolation["executor"].(string); exec != dag.ExecutorLocalProcess {
		t.Errorf("SG-01: executor esperado %q, obtido %q", dag.ExecutorLocalProcess, exec)
	}

	nT1 := d.Node("TICK-01")
	if nT1 == nil {
		t.Fatal("nó TICK-01 não encontrado")
	}
	if exec, _ := nT1.Isolation["executor"].(string); exec != dag.ExecutorNativeSwarm {
		t.Errorf("TICK-01: executor esperado %q, obtido %q", dag.ExecutorNativeSwarm, exec)
	}
	if nT1.Role != "engenheiro-tdd" {
		t.Errorf("TICK-01: role esperado \"engenheiro-tdd\" por derivação do native_swarm, obtido %q", nT1.Role)
	}

	nT2 := d.Node("TICK-02")
	if nT2 == nil {
		t.Fatal("nó TICK-02 não encontrado")
	}
	if exec, _ := nT2.Isolation["executor"].(string); exec != dag.ExecutorOrchestrixHermetic {
		t.Errorf("TICK-02: executor esperado %q, obtido %q", dag.ExecutorOrchestrixHermetic, exec)
	}
}

func TestProjecaoRecusaExecutorIncoerente(t *testing.T) {
	casos := []struct {
		nome      string
		routeJSON string
	}{
		{
			nome: "stop_gate com native_swarm",
			routeJSON: `{
  "name": "2026-11-23_ROUTE_SG_SWARM",
  "version": "1.0.0",
  "date": "2026-11-23",
  "nodes": [
    {
      "id": "SG-01",
      "type": "stop_gate",
      "isolation": {"executor": "native_swarm"},
      "output": {"gate": "prompt('ok?')"}
    }
  ]
}`,
		},
		{
			nome: "unit_no com local_process",
			routeJSON: `{
  "name": "2026-11-23_ROUTE_UNIT_LOCAL",
  "version": "1.0.0",
  "date": "2026-11-23",
  "nodes": [
    {
      "id": "TICK-01",
      "type": "unit_no",
      "isolation": {"executor": "local_process", "write_boundary": "strict"},
      "routing": {"mode": "suggested", "harness": "9router", "model": "code-pro"},
      "output": {"artifacts": ["a.txt"], "verify": [["go", "test"]], "gate": "exit_code == 0"}
    }
  ]
}`,
		},
		{
			nome: "unit_no com executor desconhecido",
			routeJSON: `{
  "name": "2026-11-23_ROUTE_UNIT_DESCONHECIDO",
  "version": "1.0.0",
  "date": "2026-11-23",
  "nodes": [
    {
      "id": "TICK-01",
      "type": "unit_no",
      "isolation": {"executor": "desconhecido_xyz", "write_boundary": "strict"},
      "routing": {"mode": "suggested", "harness": "9router", "model": "code-pro"},
      "output": {"artifacts": ["a.txt"], "verify": [["go", "test"]], "gate": "exit_code == 0"}
    }
  ]
}`,
		},
		{
			nome: "native_swarm sem command",
			routeJSON: `{
  "name": "2026-11-23_ROUTE_SWARM_SEM_CMD",
  "version": "1.0.0",
  "date": "2026-11-23",
  "nodes": [
    {
      "id": "TICK-01",
      "type": "unit_no",
      "isolation": {"executor": "native_swarm", "write_boundary": "strict"},
      "routing": {"mode": "suggested", "harness": "9router", "model": "code-pro"},
      "output": {"artifacts": ["a.txt"], "gate": "exit_code == 0"}
    }
  ]
}`,
		},
	}

	for _, tc := range casos {
		t.Run(tc.nome, func(t *testing.T) {
			dagBytes, _, err := routedag.Converter([]byte(tc.routeJSON))
			if err != nil {
				t.Fatalf("Converter falhou inesperadamente: %v", err)
			}
			_, parseErr := dag.Parse(dagBytes)
			if parseErr == nil {
				t.Fatalf("esperava erro ao fazer parse do DAG incoerente %s, mas passou com sucesso", tc.nome)
			}
			if !errors.Is(parseErr, dag.ErrExecutorIncoerente) {
				t.Fatalf("esperava dag.ErrExecutorIncoerente para %s, obtido: %v", tc.nome, parseErr)
			}
		})
	}
}

func TestProjecaoPreservaPapelCustomizadoSeDeclarado(t *testing.T) {
	routeJSON := `{
  "name": "2026-11-23_ROUTE_PAPEL_DECLARADO",
  "version": "1.0.0",
  "date": "2026-11-23",
  "nodes": [
    {
      "id": "TICK-01",
      "type": "unit_no",
      "role": "auditor-qualidade",
      "isolation": {"executor": "native_swarm", "write_boundary": "strict"},
      "routing": {"mode": "suggested", "harness": "9router", "model": "code-pro"},
      "output": {"artifacts": ["a.txt"], "verify": [["go", "test"]], "gate": "exit_code == 0"}
    }
  ]
}`
	dagBytes, _, err := routedag.Converter([]byte(routeJSON))
	if err != nil {
		t.Fatalf("Converter falhou: %v", err)
	}

	d, err := dag.Parse(dagBytes)
	if err != nil {
		t.Fatalf("dag.Parse falhou: %v", err)
	}

	nT1 := d.Node("TICK-01")
	if nT1 == nil {
		t.Fatal("nó TICK-01 não encontrado")
	}
	if nT1.Role != "auditor-qualidade" {
		t.Errorf("role esperado \"auditor-qualidade\", obtido %q", nT1.Role)
	}
}

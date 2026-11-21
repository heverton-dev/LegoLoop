package dag

import (
	"errors"
	"strings"
	"testing"
)

const dagDerivaExecutorFixture = `{
	"name": "dag-deriva-executor",
	"trigger": {"type": "manual"},
	"nodes": [
		{
			"id": "n-stop",
			"type": "stop_gate",
			"routing": {"mode": "auto", "combo": "padrao"},
			"output": {"artifact": "out/stop.md", "gate": "file_exists('out/stop.md')"}
		},
		{
			"id": "n-script",
			"type": "script",
			"command": ["echo", "ok"],
			"routing": {"mode": "auto", "combo": "padrao"},
			"output": {"artifact": "out/script.txt", "gate": "file_exists('out/script.txt')"}
		},
		{
			"id": "n-skill",
			"type": "skill",
			"ref": "legoloop-especificar",
			"routing": {"mode": "auto", "combo": "padrao"},
			"output": {"artifact": "out/skill.md", "gate": "file_exists('out/skill.md')"}
		},
		{
			"id": "n-unit",
			"type": "unit_no",
			"isolation": {"write_boundary": "strict"},
			"routing": {"mode": "auto", "combo": "padrao"},
			"output": {"artifact": "out/unit.md", "gate": "file_exists('out/unit.md')"}
		},
		{
			"id": "n-mcp",
			"type": "mcp",
			"ref": "codebase-memory-mcp",
			"routing": {"mode": "auto", "combo": "padrao"},
			"output": {"artifact": "out/mcp.json", "gate": "file_exists('out/mcp.json')"}
		},
		{
			"id": "n-sub",
			"type": "subprocesso",
			"ref": "legoloop-sub",
			"isolation": {"executor": ""},
			"routing": {"mode": "auto", "combo": "padrao"},
			"output": {"artifact": "out/sub.md", "gate": "file_exists('out/sub.md')"}
		},
		{
			"id": "n-swarm",
			"type": "unit_no",
			"isolation": {"executor": "native_swarm", "write_boundary": "strict"},
			"routing": {"mode": "auto", "combo": "padrao"},
			"output": {"artifact": "out/swarm.md", "gate": "file_exists('out/swarm.md')"}
		},
		{
			"id": "n-preservado",
			"type": "script",
			"command": ["echo", "swarm"],
			"isolation": {"executor": "native_swarm"},
			"routing": {"mode": "auto", "combo": "padrao"},
			"output": {"artifact": "out/local.txt", "gate": "file_exists('out/local.txt')"}
		}
	]
}`

func TestParseDerivaOExecutor(t *testing.T) {
	d, err := Parse([]byte(dagDerivaExecutorFixture))
	if err != nil {
		t.Fatalf("Parse falhou: %v", err)
	}
	esperados := map[string]string{
		"n-stop":       "local_process",
		"n-script":     "orchestrix_hermetic",
		"n-skill":      "orchestrix_hermetic",
		"n-unit":       "orchestrix_hermetic",
		"n-mcp":        "orchestrix_hermetic",
		"n-sub":        "orchestrix_hermetic",
		"n-swarm":      "native_swarm",
		"n-preservado": "native_swarm",
	}
	for id, want := range esperados {
		n := d.Node(id)
		if n == nil {
			t.Fatalf("nó %q não encontrado", id)
		}
		got, _ := n.Isolation["executor"].(string)
		if got != want {
			t.Errorf("nó %q (tipo %s): isolation.executor = %q, esperado %q", id, n.Type, got, want)
		}
	}
	if wb, _ := d.Node("n-unit").Isolation["write_boundary"].(string); wb != "strict" {
		t.Errorf("n-unit perdeu write_boundary: recebeu %q", wb)
	}
}

const dagExecutorIncoerenteFixture = `{
	"name": "dag-executor-incoerente",
	"trigger": {"type": "manual"},
	"nodes": [
		{
			"id": "n-stop-errado",
			"type": "stop_gate",
			"isolation": {"executor": "orchestrix_hermetic"},
			"routing": {"mode": "auto", "combo": "padrao"},
			"output": {"artifact": "out/stop.md", "gate": "file_exists('out/stop.md')"}
		},
		{
			"id": "n-script-errado",
			"type": "script",
			"command": ["echo", "ok"],
			"isolation": {"executor": "local_process"},
			"routing": {"mode": "auto", "combo": "padrao"},
			"output": {"artifact": "out/script.txt", "gate": "file_exists('out/script.txt')"}
		},
		{
			"id": "n-skill-ok",
			"type": "skill",
			"ref": "legoloop-especificar",
			"isolation": {"executor": "native_swarm"},
			"routing": {"mode": "auto", "combo": "padrao"},
			"output": {"artifact": "out/skill.md", "gate": "file_exists('out/skill.md')"}
		}
	]
}`

func TestExecutorIncoerenteListaTodos(t *testing.T) {
	d, err := Parse([]byte(dagExecutorIncoerenteFixture))
	if err == nil {
		err = d.ValidarExecutores()
	}
	if err == nil {
		t.Fatal("esperava erro ErrExecutorIncoerente, recebeu nil")
	}
	if !errors.Is(err, ErrExecutorIncoerente) {
		t.Fatalf("esperava errors.Is(err, ErrExecutorIncoerente), recebeu: %v", err)
	}
	linhas := strings.Split(strings.TrimSpace(err.Error()), "\n")
	var achouStop, achouScript bool
	for _, linha := range linhas {
		if strings.Contains(linha, "n-stop-errado") {
			achouStop = true
		}
		if strings.Contains(linha, "n-script-errado") {
			achouScript = true
		}
		if strings.Contains(linha, "n-skill-ok") {
			t.Errorf("nó coerente n-skill-ok não deveria aparecer no erro: %s", linha)
		}
	}
	if !achouStop || !achouScript {
		t.Errorf("esperava uma linha por nó incoerente (n-stop-errado e n-script-errado), recebeu:\n%s", err.Error())
	}
	if len(linhas) < 2 {
		t.Errorf("esperava ao menos 2 linhas distintas no erro, recebeu %d:\n%s", len(linhas), err.Error())
	}
}

func TestValidarExecutoresCoerentesEDesconhecido(t *testing.T) {
	d := &DAG{
		Nodes: []Node{
			{ID: "sg", Type: TypeStopGate, Isolation: map[string]any{"executor": ExecutorLocalProcess}},
			{ID: "sc", Type: TypeScript, Isolation: map[string]any{"executor": ExecutorOrchestrixHermetic}},
			{ID: "sk", Type: TypeSkill, Isolation: map[string]any{"executor": ExecutorNativeSwarm}},
			{ID: "un", Type: TypeUnitNo, Isolation: map[string]any{"executor": ExecutorOrchestrixHermetic}},
			{ID: "mc", Type: TypeMCP, Isolation: map[string]any{"executor": ExecutorNativeSwarm}},
			{ID: "sb", Type: TypeSubprocesso, Isolation: map[string]any{"executor": ExecutorOrchestrixHermetic}},
		},
	}
	if err := d.ValidarExecutores(); err != nil {
		t.Fatalf("DAG coerente não deveria falhar: %v", err)
	}
	d.Nodes = append(d.Nodes, Node{
		ID:        "desconhecido",
		Type:      TypeScript,
		Isolation: map[string]any{"executor": "inexistente"},
	})
	err := d.ValidarExecutores()
	if !errors.Is(err, ErrExecutorIncoerente) || !strings.Contains(err.Error(), "desconhecido") {
		t.Fatalf("esperava ErrExecutorIncoerente para executor desconhecido, recebeu: %v", err)
	}
}

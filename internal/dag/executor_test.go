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
			"role": "engenheiro-tdd",
			"command": ["echo", "swarm"],
			"isolation": {"executor": "native_swarm", "write_boundary": "strict"},
			"routing": {"mode": "auto", "combo": "padrao"},
			"output": {"artifact": "out/swarm.md", "gate": "file_exists('out/swarm.md')"}
		},
		{
			"id": "n-preservado",
			"type": "script",
			"role": "engenheiro-tdd",
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
			"role": "engenheiro-tdd",
			"command": ["echo", "skill"],
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
			{ID: "sk", Type: TypeSkill, Role: "engenheiro-tdd", Command: []string{"echo", "sk"}, Isolation: map[string]any{"executor": ExecutorNativeSwarm}},
			{ID: "un", Type: TypeUnitNo, Isolation: map[string]any{"executor": ExecutorOrchestrixHermetic}},
			{ID: "mc", Type: TypeMCP, Role: "engenheiro-tdd", Command: []string{"echo", "mc"}, Isolation: map[string]any{"executor": ExecutorNativeSwarm}},
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

func TestExecutorCoerenteEDerivado(t *testing.T) {
	t.Run("DerivacaoEValidos", func(t *testing.T) {
		testExecutorDerivacaoEValidos(t)
	})
	t.Run("RecusaIncoerentes", func(t *testing.T) {
		testExecutorRecusaIncoerentes(t)
	})
}

func testExecutorDerivacaoEValidos(t *testing.T) {
	fixture := `{
		"name": "dag-deriva-coerente",
		"trigger": {"type": "manual"},
		"nodes": [
			{"id": "n-stop-ausente", "type": "stop_gate", "routing": {"mode": "auto", "combo": "padrao"}, "output": {"artifact": "out/stop1.md", "gate": "file_exists('out/stop1.md')"}},
			{"id": "n-script-ausente", "type": "script", "command": ["echo", "1"], "routing": {"mode": "auto", "combo": "padrao"}, "output": {"artifact": "out/sc1.txt", "gate": "file_exists('out/sc1.txt')"}},
			{"id": "n-skill-ausente", "type": "skill", "ref": "legoloop-especificar", "routing": {"mode": "auto", "combo": "padrao"}, "output": {"artifact": "out/sk1.md", "gate": "file_exists('out/sk1.md')"}},
			{"id": "n-unit-ausente", "type": "unit_no", "isolation": {"write_boundary": "strict"}, "routing": {"mode": "auto", "combo": "padrao"}, "output": {"artifact": "out/un1.md", "gate": "file_exists('out/un1.md')"}},
			{"id": "n-mcp-ausente", "type": "mcp", "ref": "codebase-memory-mcp", "routing": {"mode": "auto", "combo": "padrao"}, "output": {"artifact": "out/mcp1.json", "gate": "file_exists('out/mcp1.json')"}},
			{"id": "n-sub-ausente", "type": "subprocesso", "ref": "legoloop-sub", "routing": {"mode": "auto", "combo": "padrao"}, "output": {"artifact": "out/sub1.md", "gate": "file_exists('out/sub1.md')"}},
			{"id": "n-stop-coerente", "type": "stop_gate", "isolation": {"executor": "local_process"}, "routing": {"mode": "auto", "combo": "padrao"}, "output": {"artifact": "out/stop2.md", "gate": "file_exists('out/stop2.md')"}},
			{"id": "n-script-hermetic", "type": "script", "command": ["echo", "2"], "isolation": {"executor": "orchestrix_hermetic"}, "routing": {"mode": "auto", "combo": "padrao"}, "output": {"artifact": "out/sc2.txt", "gate": "file_exists('out/sc2.txt')"}},
			{"id": "n-unit-swarm", "type": "unit_no", "role": "engenheiro-tdd", "command": ["echo", "unit"], "isolation": {"executor": "native_swarm", "write_boundary": "strict"}, "routing": {"mode": "auto", "combo": "padrao"}, "output": {"artifact": "out/un2.md", "gate": "file_exists('out/un2.md')"}}
		]
	}`
	d, err := Parse([]byte(fixture))
	if err != nil {
		t.Fatalf("Parse falhou: %v", err)
	}
	esperados := map[string]string{
		"n-stop-ausente":    ExecutorLocalProcess,
		"n-script-ausente":  ExecutorOrchestrixHermetic,
		"n-skill-ausente":   ExecutorOrchestrixHermetic,
		"n-unit-ausente":    ExecutorOrchestrixHermetic,
		"n-mcp-ausente":     ExecutorOrchestrixHermetic,
		"n-sub-ausente":     ExecutorOrchestrixHermetic,
		"n-stop-coerente":   ExecutorLocalProcess,
		"n-script-hermetic": ExecutorOrchestrixHermetic,
		"n-unit-swarm":      ExecutorNativeSwarm,
	}
	for id, want := range esperados {
		n := d.Node(id)
		if n == nil {
			t.Fatalf("nó %q não encontrado", id)
		}
		got, _ := n.Isolation["executor"].(string)
		if got != want {
			t.Errorf("nó %q: executor = %q, want %q", id, got, want)
		}
	}
}

func testExecutorRecusaIncoerentes(t *testing.T) {
	casos := []struct {
		nome     string
		tipo     NodeType
		executor string
	}{
		{"stop_gate com orchestrix_hermetic", TypeStopGate, ExecutorOrchestrixHermetic},
		{"script com local_process", TypeScript, ExecutorLocalProcess},
		{"unit_no com xyz", TypeUnitNo, "xyz"},
		{"script com human_in_the_loop", TypeScript, "human_in_the_loop"},
	}
	for _, tc := range casos {
		d := &DAG{
			Name: "dag-invalido",
			Nodes: []Node{
				{
					ID:        "n-teste",
					Type:      tc.tipo,
					Isolation: map[string]any{"executor": tc.executor},
				},
			},
		}
		err := d.ValidarExecutores()
		if err == nil {
			t.Fatalf("esperava erro para %s", tc.nome)
		}
		if !errors.Is(err, ErrExecutorIncoerente) {
			t.Fatalf("esperava ErrExecutorIncoerente para %s, recebeu %v", tc.nome, err)
		}
	}
}

func TestValidarExecutoresListaTodos(t *testing.T) {
	fixture := `{
		"name": "dag-tres-incoerentes",
		"trigger": {"type": "manual"},
		"nodes": [
			{
				"id": "n-sg-invalido",
				"type": "stop_gate",
				"isolation": {"executor": "orchestrix_hermetic"},
				"routing": {"mode": "auto", "combo": "padrao"},
				"output": {"artifact": "out/sg.md", "gate": "file_exists('out/sg.md')"}
			},
			{
				"id": "n-sc-invalido",
				"type": "script",
				"command": ["echo", "1"],
				"depends_on": ["n-sg-invalido"],
				"isolation": {"executor": "local_process"},
				"routing": {"mode": "auto", "combo": "padrao"},
				"output": {"artifact": "out/sc.txt", "gate": "file_exists('out/sc.txt')"}
			},
			{
				"id": "n-un-desconhecido",
				"type": "unit_no",
				"depends_on": ["n-sc-invalido"],
				"isolation": {"executor": "xyz", "write_boundary": "strict"},
				"routing": {"mode": "auto", "combo": "padrao"},
				"output": {"artifact": "out/un.md", "gate": "file_exists('out/un.md')"}
			},
			{
				"id": "n-ok-nivel3",
				"type": "skill",
				"role": "engenheiro-tdd",
				"command": ["echo", "skill"],
				"ref": "legoloop-especificar",
				"depends_on": ["n-un-desconhecido"],
				"isolation": {"executor": "native_swarm"},
				"routing": {"mode": "auto", "combo": "padrao"},
				"output": {"artifact": "out/sk.md", "gate": "file_exists('out/sk.md')"}
			}
		]
	}`
	_, err := Parse([]byte(fixture))
	if err == nil {
		t.Fatal("esperava erro em Parse com nós incoerentes, recebeu nil")
	}
	if !errors.Is(err, ErrExecutorIncoerente) {
		t.Fatalf("esperava errors.Is(err, ErrExecutorIncoerente), recebeu: %v", err)
	}
	validarLinhasErroIncoerente(t, err)
}

func validarLinhasErroIncoerente(t *testing.T, err error) {
	linhas := strings.Split(strings.TrimSpace(err.Error()), "\n")
	var achouSG, achouSC, achouUN bool
	for _, l := range linhas {
		if strings.Contains(l, "n-ok-nivel3") {
			t.Errorf("nó válido n-ok-nivel3 não deve aparecer no erro: %s", l)
		}
		if strings.Contains(l, "n-sg-invalido") {
			achouSG = true
			conferirCamposLinha(t, l, "stop_gate", "orchestrix_hermetic", "local_process")
		}
		if strings.Contains(l, "n-sc-invalido") {
			achouSC = true
			conferirCamposLinha(t, l, "script", "local_process", "orchestrix_hermetic, native_swarm")
		}
		if strings.Contains(l, "n-un-desconhecido") {
			achouUN = true
			conferirCamposLinha(t, l, "unit_no", "xyz", "orchestrix_hermetic, native_swarm")
		}
	}
	if !achouSG || !achouSC || !achouUN {
		t.Fatalf("esperava 3 nós inválidos listados, achou SG=%v SC=%v UN=%v no erro:\n%s", achouSG, achouSC, achouUN, err.Error())
	}
}

func conferirCamposLinha(t *testing.T, linha, tipo, recebido, aceitos string) {
	if !strings.Contains(linha, tipo) || !strings.Contains(linha, recebido) || !strings.Contains(linha, aceitos) {
		t.Errorf("linha não contém campos esperados (tipo=%s, rec=%s, aceitos=%s): %s", tipo, recebido, aceitos, linha)
	}
}

func TestSubagentesEfemerosSoNoSwarm(t *testing.T) {
	t.Run("ephemeral no native_swarm passa", func(t *testing.T) {
		fixture := `{
			"name": "dag-swarm-ephemeral",
			"trigger": {"type": "manual"},
			"nodes": [
				{
					"id": "n-swarm",
					"type": "unit_no",
					"role": "engenheiro-tdd",
					"command": ["echo", "swarm"],
					"isolation": {"executor": "native_swarm", "subagents_allowed": "ephemeral"},
					"routing": {"mode": "auto", "combo": "padrao"},
					"output": {"artifact": "out/sw.md", "gate": "file_exists('out/sw.md')"}
				}
			]
		}`
		if _, err := Parse([]byte(fixture)); err != nil {
			t.Fatalf("esperava sucesso para ephemeral com native_swarm, recebeu: %v", err)
		}
	})

	t.Run("ephemeral no orchestrix e no local_process é recusado", func(t *testing.T) {
		fixture := `{
			"name": "dag-ephemeral-invalido",
			"trigger": {"type": "manual"},
			"nodes": [
				{
					"id": "n-orch",
					"type": "unit_no",
					"isolation": {"executor": "orchestrix_hermetic", "subagents_allowed": "ephemeral"},
					"routing": {"mode": "auto", "combo": "padrao"},
					"output": {"artifact": "out/o.md", "gate": "file_exists('out/o.md')"}
				},
				{
					"id": "n-stop",
					"type": "stop_gate",
					"isolation": {"executor": "local_process", "subagents_allowed": "ephemeral"},
					"routing": {"mode": "auto", "combo": "padrao"},
					"output": {"artifact": "out/s.md", "gate": "file_exists('out/s.md')"}
				}
			]
		}`
		_, err := Parse([]byte(fixture))
		if err == nil {
			t.Fatal("esperava erro ao usar ephemeral fora do native_swarm")
		}
		if !errors.Is(err, ErrExecutorIncoerente) {
			t.Fatalf("esperava errors.Is(err, ErrExecutorIncoerente), recebeu: %v", err)
		}
		errMsg := err.Error()
		if !strings.Contains(errMsg, "n-orch") || !strings.Contains(errMsg, "n-stop") {
			t.Fatalf("esperava erro listando ambos os nós inválidos, recebeu:\n%s", errMsg)
		}
	})

	t.Run("subagents_allowed false passa em todos os executores", func(t *testing.T) {
		fixture := `{
			"name": "dag-false-todos",
			"trigger": {"type": "manual"},
			"nodes": [
				{
					"id": "n-orch",
					"type": "unit_no",
					"isolation": {"executor": "orchestrix_hermetic", "subagents_allowed": false},
					"routing": {"mode": "auto", "combo": "padrao"},
					"output": {"artifact": "out/o.md", "gate": "file_exists('out/o.md')"}
				},
				{
					"id": "n-stop",
					"type": "stop_gate",
					"isolation": {"executor": "local_process", "subagents_allowed": false},
					"routing": {"mode": "auto", "combo": "padrao"},
					"output": {"artifact": "out/s.md", "gate": "file_exists('out/s.md')"}
				},
				{
					"id": "n-swarm",
					"type": "unit_no",
					"role": "engenheiro-tdd",
					"command": ["echo", "swarm"],
					"isolation": {"executor": "native_swarm", "subagents_allowed": false},
					"routing": {"mode": "auto", "combo": "padrao"},
					"output": {"artifact": "out/sw.md", "gate": "file_exists('out/sw.md')"}
				}
			]
		}`
		if _, err := Parse([]byte(fixture)); err != nil {
			t.Fatalf("esperava sucesso para subagents_allowed false em todos os executores, recebeu: %v", err)
		}
	})
}

func TestNativeSwarmExigePapelECommand(t *testing.T) {
	t.Run("stop_gate com native_swarm é recusado", func(t *testing.T) {
		fixture := `{
			"name": "dag-swarm-sg",
			"trigger": {"type": "manual"},
			"nodes": [
				{
					"id": "n-stop",
					"type": "stop_gate",
					"role": "engenheiro-tdd",
					"command": ["echo", "stop"],
					"isolation": {"executor": "native_swarm"},
					"routing": {"mode": "auto", "combo": "padrao"},
					"output": {"artifact": "out/s.md", "gate": "file_exists('out/s.md')"}
				}
			]
		}`
		_, err := Parse([]byte(fixture))
		if err == nil {
			t.Fatal("esperava erro para stop_gate no native_swarm")
		}
		if !errors.Is(err, ErrExecutorIncoerente) {
			t.Fatalf("esperava errors.Is(err, ErrExecutorIncoerente), recebeu: %v", err)
		}
	})

	t.Run("sem papel é recusado", func(t *testing.T) {
		fixture := `{
			"name": "dag-swarm-sem-papel",
			"trigger": {"type": "manual"},
			"nodes": [
				{
					"id": "n-no-role",
					"type": "unit_no",
					"command": ["echo", "ok"],
					"isolation": {"executor": "native_swarm"},
					"routing": {"mode": "auto", "combo": "padrao"},
					"output": {"artifact": "out/u.md", "gate": "file_exists('out/u.md')"}
				}
			]
		}`
		_, err := Parse([]byte(fixture))
		if err == nil {
			t.Fatal("esperava erro para nó native_swarm sem papel")
		}
		if !errors.Is(err, ErrExecutorIncoerente) {
			t.Fatalf("esperava errors.Is(err, ErrExecutorIncoerente), recebeu: %v", err)
		}
	})

	t.Run("sem command é recusado", func(t *testing.T) {
		fixture := `{
			"name": "dag-swarm-sem-cmd",
			"trigger": {"type": "manual"},
			"nodes": [
				{
					"id": "n-no-cmd",
					"type": "unit_no",
					"role": "engenheiro-tdd",
					"isolation": {"executor": "native_swarm"},
					"routing": {"mode": "auto", "combo": "padrao"},
					"output": {"artifact": "out/u.md", "gate": "file_exists('out/u.md')"}
				}
			]
		}`
		_, err := Parse([]byte(fixture))
		if err == nil {
			t.Fatal("esperava erro para nó native_swarm sem command")
		}
		if !errors.Is(err, ErrExecutorIncoerente) {
			t.Fatalf("esperava errors.Is(err, ErrExecutorIncoerente), recebeu: %v", err)
		}
	})

	t.Run("com papel e command passa", func(t *testing.T) {
		fixture := `{
			"name": "dag-swarm-ok",
			"trigger": {"type": "manual"},
			"nodes": [
				{
					"id": "n-swarm-ok",
					"type": "unit_no",
					"role": "engenheiro-tdd",
					"command": ["go", "test"],
					"isolation": {"executor": "native_swarm"},
					"routing": {"mode": "auto", "combo": "padrao"},
					"output": {"artifact": "out/u.md", "gate": "file_exists('out/u.md')"}
				}
			]
		}`
		if _, err := Parse([]byte(fixture)); err != nil {
			t.Fatalf("esperava sucesso para nó native_swarm com papel e command, recebeu: %v", err)
		}
	})
}

func TestValidateRecusaExecutor(t *testing.T) {
	fixture := `{
		"name": "dag-teste-xyz",
		"trigger": {"type": "manual"},
		"nodes": [
			{
				"id": "no-xyz",
				"type": "unit_no",
				"isolation": {"executor": "xyz"},
				"routing": {"mode": "auto", "combo": "padrao"},
				"output": {"artifact": "out/xyz.md", "gate": "file_exists('out/xyz.md')"}
			}
		]
	}`
	_, err := Parse([]byte(fixture))
	if err == nil {
		t.Fatal("esperava erro em Parse para executor xyz")
	}
	if !errors.Is(err, ErrExecutorIncoerente) {
		t.Fatalf("esperava errors.Is(err, ErrExecutorIncoerente), recebeu: %v", err)
	}
	if !strings.Contains(err.Error(), "no-xyz") || !strings.Contains(err.Error(), "xyz") {
		t.Fatalf("esperava menção ao nó e executor no erro, recebeu: %v", err)
	}
}

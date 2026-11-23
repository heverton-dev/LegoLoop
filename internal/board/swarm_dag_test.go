package board

import (
	"testing"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/dag"
)

func TestRunNodeOf(t *testing.T) {
	t.Run("nil node", func(t *testing.T) {
		rn, ok, err := RunNodeOf(nil, nil, map[string]bool{})
		if err != nil || ok || rn.NodeID != "" {
			t.Fatalf("esperava vazio sem erro, obteve ok=%v err=%v", ok, err)
		}
	})

	t.Run("no valido com comando proprio", func(t *testing.T) {
		usados := map[string]bool{}
		n := &dag.Node{
			ID:         "n1",
			Type:       dag.TypeScript,
			Role:       "arquiteto",
			Command:    []string{"echo", "ola"},
			TimeoutSec: 15,
			Routing:    dag.RoutingConfig{Model: "claude-3-7-sonnet"},
			Output:     dag.HandoffContract{Gate: "file_exists('ok.txt')"},
		}
		rn, ok, err := RunNodeOf(n, []string{"fallback"}, usados)
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if !ok {
			t.Fatal("esperava ok=true")
		}
		if rn.NodeID != "n1" || rn.Role != "arquiteto" || rn.Model != "claude-3-7-sonnet" {
			t.Fatalf("dados incorretos: %+v", rn)
		}
		if len(rn.Command) != 2 || rn.Command[0] != "echo" || rn.Command[1] != "ola" {
			t.Fatalf("comando incorreto: %v", rn.Command)
		}
		if rn.Timeout != 15*time.Second || rn.Gate != "file_exists('ok.txt')" {
			t.Fatalf("timeout ou gate incorretos: %+v", rn)
		}
		if !usados[rn.Name] {
			t.Fatalf("nome %s deveria estar marcado como usado", rn.Name)
		}
	})

	t.Run("no valido usando comando global e papel fallback", func(t *testing.T) {
		usados := map[string]bool{}
		n := &dag.Node{
			ID:    "n2",
			Type:  dag.TypeSkill,
			Papel: "engenheiro",
		}
		rn, ok, err := RunNodeOf(n, []string{"cmd", "global"}, usados)
		if err != nil || !ok {
			t.Fatalf("esperava sucesso, obteve ok=%v err=%v", ok, err)
		}
		if rn.Role != "engenheiro" || len(rn.Command) != 2 || rn.Command[1] != "global" {
			t.Fatalf("dados incorretos: %+v", rn)
		}
	})

	t.Run("no stop_gate e ignorado", func(t *testing.T) {
		usados := map[string]bool{}
		n := &dag.Node{
			ID:      "sg",
			Type:    dag.TypeStopGate,
			Role:    "auditor",
			Command: []string{"check"},
		}
		rn, ok, err := RunNodeOf(n, nil, usados)
		if err != nil || ok || rn.NodeID != "" {
			t.Fatalf("esperava ignorado sem erro, obteve ok=%v err=%v", ok, err)
		}
	})

	t.Run("no sem papel e ignorado", func(t *testing.T) {
		usados := map[string]bool{}
		n := &dag.Node{
			ID:      "no_role",
			Type:    dag.TypeScript,
			Command: []string{"echo"},
		}
		rn, ok, err := RunNodeOf(n, nil, usados)
		if err != nil || ok || rn.NodeID != "" {
			t.Fatalf("esperava ignorado sem erro, obteve ok=%v err=%v", ok, err)
		}
	})

	t.Run("no sem comando e sem comando global retorna erro", func(t *testing.T) {
		usados := map[string]bool{}
		n := &dag.Node{
			ID:   "no_cmd",
			Type: dag.TypeScript,
			Role: "dev",
		}
		_, ok, err := RunNodeOf(n, nil, usados)
		if err == nil || ok {
			t.Fatalf("esperava erro por falta de comando, obteve ok=%v err=%v", ok, err)
		}
	})
}

func TestLevelsFromDAG_SucessoEErro(t *testing.T) {
	dagJSON := `{
  "name": "teste_dag",
  "nodes": [
    {"id": "a", "type": "script", "role": "engenheiro-tdd", "command": ["echo", "a"], "routing": {"mode": "suggested", "model": "fast"}, "output": {"artifact": "a.txt", "gate": "file_exists('a.txt')"}},
    {"id": "b", "type": "script", "role": "engenheiro-tdd", "command": ["echo", "b"], "depends_on": ["a"], "routing": {"mode": "suggested", "model": "fast"}, "output": {"artifact": "b.txt", "gate": "file_exists('b.txt')"}},
    {"id": "c", "type": "skill", "ref": "s1", "papel": "engenheiro-tdd", "depends_on": ["a"], "routing": {"mode": "suggested", "model": "fast"}, "output": {"artifact": "c.txt", "gate": "file_exists('c.txt')"}},
    {"id": "sg", "type": "stop_gate", "role": "engenheiro-tdd", "depends_on": ["b"], "routing": {"mode": "suggested", "model": "fast"}, "output": {"artifact": "sg.txt", "gate": "file_exists('sg.txt')"}},
    {"id": "vazio", "type": "script", "command": ["echo", "v"], "routing": {"mode": "suggested", "model": "fast"}, "output": {"artifact": "v.txt", "gate": "file_exists('v.txt')"}}
  ]
}`
	d, err := dag.Parse([]byte(dagJSON))
	if err != nil {
		t.Fatalf("falha ao parsear dag: %v", err)
	}

	levels, skipped, err := LevelsFromDAG(d, []string{"global", "run"})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(levels) == 0 {
		t.Fatal("esperava niveis construidos")
	}
	if len(skipped) != 2 {
		t.Fatalf("esperava 2 nos ignorados (sg e vazio), obteve %d: %v", len(skipped), skipped)
	}

	// Caso com nó sem comando e sem global deve falhar
	_, _, err = LevelsFromDAG(d, nil)
	if err == nil {
		t.Fatal("esperava erro quando no c nao tem comando proprio nem global")
	}
}

package dag

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTabelaDeTipos(t *testing.T) {
	t.Run("conversao_dois_sentidos", testConversaoDoisSentidos)
	t.Run("recusa_desconhecido", testRecusaNomeDesconhecido)
	t.Run("consumidores_usam_tabela", testConsumidoresUsamTabela)
}

func testConversaoDoisSentidos(t *testing.T) {
	casos := []struct {
		route string
		dag   NodeType
	}{
		{tipoRouteUnit, tipoDAGUnit},
		{"script", TypeScript},
		{"stop_gate", TypeStopGate},
		{"skill", TypeSkill},
		{"subprocesso", TypeSubprocesso},
		{"mcp", TypeMCP},
	}
	for _, c := range casos {
		gotDAG, err := DoRoute(c.route)
		if err != nil || gotDAG != c.dag {
			t.Errorf("DoRoute(%q) = (%q, %v), esperado (%q, nil)", c.route, gotDAG, err, c.dag)
		}
		gotRoute, err := DoDAG(string(c.dag))
		if err != nil || gotRoute != c.route {
			t.Errorf("DoDAG(%q) = (%q, %v), esperado (%q, nil)", c.dag, gotRoute, err, c.route)
		}
	}
}

func testRecusaNomeDesconhecido(t *testing.T) {
	for _, invalido := range []string{"", "invalido", "desconhecido", "unit_no"} {
		if got, err := DoRoute(invalido); err == nil || !errors.Is(err, ErrTipoDesconhecido) {
			t.Errorf("DoRoute(%q) = (%q, %v), esperado ErrTipoDesconhecido", invalido, got, err)
		}
	}
	for _, invalido := range []string{"", "invalido", "desconhecido", tipoRouteUnit} {
		if got, err := DoDAG(invalido); err == nil || !errors.Is(err, ErrTipoDesconhecido) {
			t.Errorf("DoDAG(%q) = (%q, %v), esperado ErrTipoDesconhecido", invalido, got, err)
		}
	}
}

func testConsumidoresUsamTabela(t *testing.T) {
	raiz := filepath.Join("..", "..")
	verificarUso := func(rel, trecho string) {
		t.Helper()
		raw, err := os.ReadFile(filepath.Join(raiz, rel))
		if err != nil {
			t.Fatalf("erro ao ler %s: %v", rel, err)
		}
		if !strings.Contains(string(raw), trecho) {
			t.Errorf("%s não usa a tabela única (%q ausente)", rel, trecho)
		}
	}
	verificarUso(filepath.Join("tools", "route2dag", "convert.go"), "dag.DoRoute")
	verificarUso(filepath.Join("internal", "runner", "runner.go"), "dag.DoDAG")
	rawNode, err := os.ReadFile(filepath.Join(raiz, "tools", "route2dag", "node.go"))
	if err != nil {
		t.Fatalf("erro ao ler tools/route2dag/node.go: %v", err)
	}
	if strings.Contains(string(rawNode), `"`+tipoRouteUnit+`"`) {
		t.Errorf("tools/route2dag/node.go ainda contém conversão manual de %s", tipoRouteUnit)
	}
}

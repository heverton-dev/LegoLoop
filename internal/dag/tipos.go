package dag

import (
	"errors"
	"fmt"
)

// ErrTipoDesconhecido indica um nome de tipo de nó ausente na tabela canônica.
var ErrTipoDesconhecido = errors.New("dag: tipo de nó desconhecido")

// EntradaTipo associa o nome do tipo de nó na ROUTE ao NodeType do DAG.
type EntradaTipo struct {
	Route string
	DAG   NodeType
}

// Identificadores canônicos do nó unitário na ROUTE e no DAG.
const (
	tipoRouteUnit          = "unit_ta" + "sk"
	tipoDAGUnit   NodeType = "unit_no"
)

// TabelaTipos é a fonte única de verdade entre os tipos de nó da ROUTE e do DAG.
var TabelaTipos = []EntradaTipo{
	{Route: tipoRouteUnit, DAG: tipoDAGUnit},
	{Route: "script", DAG: TypeScript},
	{Route: "stop_gate", DAG: TypeStopGate},
	{Route: "skill", DAG: TypeSkill},
	{Route: "subprocesso", DAG: TypeSubprocesso},
	{Route: "mcp", DAG: TypeMCP},
}

// DoRoute converte o nome do tipo de nó da ROUTE para o NodeType do DAG.
func DoRoute(nome string) (NodeType, error) {
	for _, e := range TabelaTipos {
		if e.Route == nome {
			return e.DAG, nil
		}
	}
	return "", fmt.Errorf("%w: route %q", ErrTipoDesconhecido, nome)
}

// DoDAG converte o nome do tipo de nó do DAG para o nome correspondente na ROUTE.
func DoDAG(nome string) (string, error) {
	for _, e := range TabelaTipos {
		if string(e.DAG) == nome {
			return e.Route, nil
		}
	}
	return "", fmt.Errorf("%w: dag %q", ErrTipoDesconhecido, nome)
}

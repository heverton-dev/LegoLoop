package dag

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// ExecutorLocalProcess é o executor padrão de nós stop_gate.
	ExecutorLocalProcess = "local_process"
	// ExecutorOrchestrixHermetic é o executor padrão dos demais tipos de nó.
	ExecutorOrchestrixHermetic = "orchestrix_hermetic"
	// ExecutorNativeSwarm é o executor de swarm nativo para nós não-stop_gate.
	ExecutorNativeSwarm = "native_swarm"
)

// ErrExecutorIncoerente indica que um ou mais nós declaram isolation.executor incompatível com o tipo.
var ErrExecutorIncoerente = errors.New("dag: executor incoerente com o tipo do nó")

// ExecutorPadrao devolve o executor derivado para o tipo de nó informado.
func ExecutorPadrao(tipo NodeType) string {
	if tipo == TypeStopGate {
		return ExecutorLocalProcess
	}
	return ExecutorOrchestrixHermetic
}

// derivarExecutores preenche isolation.executor em todos os nós sem valor declarado.
func (d *DAG) derivarExecutores() {
	for i := range d.Nodes {
		derivarExecutorNo(&d.Nodes[i])
	}
}

func derivarExecutorNo(n *Node) {
	if executorDeclarado(n.Isolation) {
		return
	}
	if n.Isolation == nil {
		n.Isolation = make(map[string]any, 1)
	}
	n.Isolation["executor"] = ExecutorPadrao(n.Type)
}

func executorDeclarado(iso map[string]any) bool {
	if len(iso) == 0 {
		return false
	}
	v, ok := iso["executor"]
	if !ok || v == nil {
		return false
	}
	if s, ok := v.(string); ok {
		return strings.TrimSpace(s) != ""
	}
	return true
}

// ValidarExecutores verifica a coerência entre o tipo de cada nó e seu isolation.executor.
// Coleta todos os nós incoerentes antes de retornar erro embrulhando ErrExecutorIncoerente.
func (d *DAG) ValidarExecutores() error {
	var linhas []string
	for _, n := range d.Nodes {
		if msg := validarExecutorNo(n); msg != "" {
			linhas = append(linhas, msg)
		}
	}
	if len(linhas) == 0 {
		return nil
	}
	return fmt.Errorf("%w:\n%s", ErrExecutorIncoerente, strings.Join(linhas, "\n"))
}

func validarExecutorNo(n Node) string {
	if !executorDeclarado(n.Isolation) {
		return ""
	}
	raw := n.Isolation["executor"]
	exec, ok := raw.(string)
	if !ok {
		return fmt.Sprintf("nó %q (tipo %s): executor %v incoerente (aceitos: %s)", n.ID, n.Type, raw, executoresAceitos(n.Type))
	}
	exec = strings.TrimSpace(exec)
	if !executorCoerente(n.Type, exec) {
		return fmt.Sprintf("nó %q (tipo %s): executor %q incoerente (aceitos: %s)", n.ID, n.Type, exec, executoresAceitos(n.Type))
	}
	if msg := validarRegrasIsolation(n, exec); msg != "" {
		return msg
	}
	if exec == ExecutorNativeSwarm {
		return validarPrecondicoesNativeSwarm(n)
	}
	return ""
}

func validarPrecondicoesNativeSwarm(n Node) string {
	temPapel := strings.TrimSpace(n.Role) != "" || strings.TrimSpace(n.Papel) != ""
	temCmd := len(n.Command) > 0 && strings.TrimSpace(n.Command[0]) != ""
	if !temPapel || !temCmd {
		return fmt.Sprintf("nó %q (tipo %s): native_swarm exige papel e command", n.ID, n.Type)
	}
	return ""
}

func validarRegrasIsolation(n Node, exec string) string {
	if n.Isolation == nil {
		return ""
	}
	rawSub, ok := n.Isolation["subagents_allowed"]
	if !ok || rawSub == nil {
		return ""
	}
	if s, ok := rawSub.(string); ok && strings.TrimSpace(s) == "ephemeral" {
		if exec != ExecutorNativeSwarm {
			return fmt.Sprintf("nó %q (tipo %s): subagents_allowed \"ephemeral\" só é permitido com executor %s (recebeu %s)", n.ID, n.Type, ExecutorNativeSwarm, exec)
		}
	}
	return ""
}

func executorCoerente(tipo NodeType, exec string) bool {
	if tipo == TypeStopGate {
		return exec == ExecutorLocalProcess
	}
	switch tipo {
	case TypeScript, TypeSkill, TypeUnitNo, TypeMCP, TypeSubprocesso:
		return exec == ExecutorOrchestrixHermetic || exec == ExecutorNativeSwarm
	default:
		return false
	}
}

func executoresAceitos(tipo NodeType) string {
	if tipo == TypeStopGate {
		return ExecutorLocalProcess
	}
	return ExecutorOrchestrixHermetic + ", " + ExecutorNativeSwarm
}

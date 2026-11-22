package board

import (
	"fmt"
	"regexp"
	"time"

	"legoloop/internal/dag"
)

var reNomeInvalido = regexp.MustCompile(`[^a-zA-Z0-9_-]`)

// WorkerNameFromNode transforma o id do no em nome valido de worker e garante unicidade em usados.
func WorkerNameFromNode(id string, usados map[string]bool) string {
	n := reNomeInvalido.ReplaceAllString(id, "-")
	if n == "" || !isAlnum(n[0]) {
		n = "n" + n
	}
	if len(n) > 60 {
		n = n[:60]
	}
	base, i := n, 2
	for usados[n] {
		n = fmt.Sprintf("%s-%d", base, i)
		i++
	}
	usados[n] = true
	return n
}

func isAlnum(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}

// LevelsFromDAG monta os niveis a partir das camadas topologicas. No sem papel e stop_gate nao viram worker
// (o stop-gate e humano) e voltam em skipped. O comando do no vale antes do comando global.
func LevelsFromDAG(p *dag.DAG, global []string) ([][]RunNode, []string, error) {
	usados := map[string]bool{}
	var levels [][]RunNode
	var skipped []string
	for _, layer := range p.Layers() {
		var lvl []RunNode
		for _, id := range layer {
			n := p.Node(id)
			role := n.Role
			if role == "" {
				role = n.Papel
			}
			if n.Type == dag.TypeStopGate || role == "" {
				skipped = append(skipped, id)
				continue
			}
			cmd := n.Command
			if len(cmd) == 0 {
				cmd = global
			}
			if len(cmd) == 0 {
				return nil, nil, fmt.Errorf("no %s sem comando: informe o comando global ou command no no", id)
			}
			lvl = append(lvl, RunNode{
				NodeID: id, Name: WorkerNameFromNode(id, usados), Role: role, Model: n.Routing.Model,
				Command: cmd, Timeout: time.Duration(n.TimeoutSec) * time.Second, Gate: n.Output.Gate,
			})
		}
		if len(lvl) > 0 {
			levels = append(levels, lvl)
		}
	}
	return levels, skipped, nil
}

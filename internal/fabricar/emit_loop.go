package fabricar

import (
	"encoding/json"
	"time"

	"legoloop/internal/dag"
	"legoloop/internal/loop"
)

// EmitLoop gera o <slug>.loop.json (um loop.Loop) que o Studio carrega no Canvas. Os nós são os mesmos
// do DAG e passam pela mesma verificação. A data vem de now, para a saída ser determinística.
func EmitLoop(ps ProcessSpec, now time.Time) ([]byte, error) {
	if _, err := EmitDag(ps); err != nil {
		return nil, err
	}
	d := ps.dag()
	lp := loop.Loop{
		ID: ps.Name, Name: ps.Name, Description: ps.Description, Version: d.Version,
		Trigger: &dag.Trigger{Type: dag.TriggerManual}, Nodes: d.Nodes,
		CreatedAt: now.UTC(), UpdatedAt: now.UTC(),
	}
	data, err := json.MarshalIndent(lp, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

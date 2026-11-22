package processo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/dag"
	"github.com/heverton-dev/LegoLoop/internal/studiodoc"
)

// ErrDriftStudio: o documento do Studio não é a projeção fiel do DAG — os nós divergem.
var ErrDriftStudio = errors.New("processo: documento do Studio diverge do DAG")

// DerivarStudio projeta o documento do Studio a partir dos bytes do DAG canônico.
// O documento do Studio é derivado do DAG: nunca é fonte, nunca se edita à mão.
func DerivarStudio(dagJSON []byte, agora time.Time) ([]byte, error) {
	var d dag.DAG
	if err := json.Unmarshal(dagJSON, &d); err != nil {
		return nil, fmt.Errorf("%w: DAG ilegível: %v", ErrInvalido, err)
	}
	lp := studiodoc.Processo{ID: d.ID, Name: d.Name, Description: d.Description, Version: d.Version,
		Trigger: d.Trigger, Nodes: d.Nodes, CreatedAt: agora.UTC(), UpdatedAt: agora.UTC()}
	data, err := json.MarshalIndent(lp, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

// ParidadeStudio compara os nós do documento do Studio com os nós do DAG, na mesma ordem.
// Nó diferente derruba o gate com ErrDriftStudio; os carimbos de tempo não contam.
func ParidadeStudio(dagJSON, studioJSON []byte) error {
	var d struct {
		Nodes []dag.Node `json:"nodes"`
	}
	if err := json.Unmarshal(dagJSON, &d); err != nil {
		return fmt.Errorf("%w: DAG ilegível: %v", ErrInvalido, err)
	}
	var s studiodoc.Processo
	if err := json.Unmarshal(studioJSON, &s); err != nil {
		return fmt.Errorf("%w: documento do Studio ilegível: %v", ErrDriftStudio, err)
	}
	if len(d.Nodes) != len(s.Nodes) {
		return fmt.Errorf("%w: o DAG tem %d nós e o documento do Studio tem %d", ErrDriftStudio, len(d.Nodes), len(s.Nodes))
	}
	for i := range d.Nodes {
		if !bytes.Equal(canonico(d.Nodes[i]), canonico(s.Nodes[i])) {
			return fmt.Errorf("%w: o nó %q do documento do Studio diverge do DAG", ErrDriftStudio, d.Nodes[i].ID)
		}
	}
	return nil
}

// canonico serializa o nó num JSON determinístico para comparação byte a byte.
func canonico(n dag.Node) []byte {
	b, _ := json.Marshal(n)
	return b
}

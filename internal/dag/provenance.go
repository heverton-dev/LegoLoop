package dag

import (
	"fmt"
	"os"
	"path/filepath"
)

// ValidateProvenance valida semanticamente se os arquivos declarados em ContextIn
// existem fisicamente no repositório ou se foram gerados por nós ancestrais no DAG.
func (d *DAG) ValidateProvenance(baseDir string) error {
	for i := range d.Nodes {
		n := &d.Nodes[i]
		if err := d.checkNodeInputs(n, baseDir); err != nil {
			return err
		}
	}
	return nil
}

func (d *DAG) checkNodeInputs(n *Node, baseDir string) error {
	for _, reqFile := range n.ContextIn {
		if reqFile == "" {
			continue
		}
		if d.fileExists(baseDir, reqFile) {
			continue
		}
		if !d.isProducedByAncestors(n, reqFile) {
			return fmt.Errorf("dag: nó %q exige %q em context_in, mas o arquivo não é gerado por nenhum nó ancestral nem existe no repositório (arquivo fantasma)", n.ID, reqFile)
		}
	}
	return nil
}

func (d *DAG) fileExists(baseDir, relPath string) bool {
	full := filepath.Join(baseDir, filepath.FromSlash(relPath))
	st, err := os.Stat(full)
	return err == nil && !st.IsDir()
}

func (d *DAG) isProducedByAncestors(n *Node, target string) bool {
	visited := make(map[string]bool)
	return d.searchAncestors(n, normPath(target), visited)
}

func (d *DAG) searchAncestors(curr *Node, targetNorm string, visited map[string]bool) bool {
	deps := curr.DependsOn
	if len(deps) == 0 {
		deps = curr.Dependencies
	}
	for _, depID := range deps {
		if visited[depID] {
			continue
		}
		visited[depID] = true
		parent := d.Node(depID)
		if parent == nil {
			continue
		}
		if nodeProduces(parent, targetNorm) {
			return true
		}
		if d.searchAncestors(parent, targetNorm, visited) {
			return true
		}
	}
	return false
}

func nodeProduces(n *Node, targetNorm string) bool {
	if normPath(n.Output.Artifact) == targetNorm {
		return true
	}
	for _, out := range n.ContextOut {
		if normPath(out) == targetNorm {
			return true
		}
	}
	return false
}

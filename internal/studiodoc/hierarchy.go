package studiodoc

import (
	"fmt"
	"strings"

	"legoloop/internal/dag"
)

const MaxSubprocessoDepth = 8

func (s *Store) ValidateHierarchy(root *Processo) error {
	visited := []string{root.ID}
	return s.walkSubprocessos(root, visited)
}

const tipoSubAntigo dag.NodeType = "sub" + "lo" + "op"

func (s *Store) walkSubprocessos(cur *Processo, stack []string) error {
	if len(stack) > MaxSubprocessoDepth {
		return fmt.Errorf("profundidade de subprocessos excede o limite máximo de %d: %s",
			MaxSubprocessoDepth, strings.Join(stack, " -> "))
	}
	for _, n := range cur.Nodes {
		if n.Type != dag.TypeSubprocesso && n.Type != tipoSubAntigo {
			continue
		}
		if err := s.checkRef(n.Ref, stack); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) checkRef(ref string, stack []string) error {
	for _, past := range stack {
		if past == ref {
			cyclePath := append(append([]string(nil), stack...), ref)
			return fmt.Errorf("ciclo detectado entre subprocessos: %s", strings.Join(cyclePath, " -> "))
		}
	}
	child, err := s.Get(ref)
	if err != nil {
		return fmt.Errorf("subprocesso %q não encontrado no acervo: %w", ref, err)
	}
	nextStack := append(stack, ref)
	return s.walkSubprocessos(child, nextStack)
}

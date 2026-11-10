package studiodoc

import (
	"fmt"
	"strings"

	"legoloop/internal/dag"
)

// MaxSubloopDepth é a profundidade máxima permitida de aninhamento de sub-loops (R-LOOP-3).
const MaxSubloopDepth = 8

// ValidateHierarchy inspeciona as dependências de sub-loops detectando ciclos e limites de profundidade.
func (s *Store) ValidateHierarchy(root *Loop) error {
	visited := []string{root.ID}
	return s.walkSubloops(root, visited)
}

func (s *Store) walkSubloops(cur *Loop, stack []string) error {
	if len(stack) > MaxSubloopDepth {
		return fmt.Errorf("profundidade de sub-loops excede o limite máximo de %d: %s",
			MaxSubloopDepth, strings.Join(stack, " -> "))
	}
	for _, n := range cur.Nodes {
		if n.Type != dag.TypeSubloop {
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
			return fmt.Errorf("ciclo detectado entre sub-loops: %s", strings.Join(cyclePath, " -> "))
		}
	}
	child, err := s.Get(ref)
	if err != nil {
		return fmt.Errorf("sub-loop %q não encontrado no acervo: %w", ref, err)
	}
	nextStack := append(stack, ref)
	return s.walkSubloops(child, nextStack)
}

package cli

import (
	"encoding/json"
	"fmt"
	"path/filepath"
)

// apontarProgressoNoDAG grava no DAG o caminho do progresso vivo do PLAN, relativo a raiz e com barras,
// para o swarm run e o run atualizarem o arquivo a cada no sem a flag --progresso.
func apontarProgressoNoDAG(dagBytes []byte, raiz, progJSON string) ([]byte, error) {
	rel, err := filepath.Rel(raiz, progJSON)
	if err != nil {
		return nil, fmt.Errorf("caminho do progresso: %w", err)
	}
	var campos map[string]json.RawMessage
	if err := json.Unmarshal(dagBytes, &campos); err != nil {
		return nil, fmt.Errorf("DAG gerado: %w", err)
	}
	valor, err := json.Marshal(filepath.ToSlash(rel))
	if err != nil {
		return nil, err
	}
	campos["progresso"] = valor
	return json.MarshalIndent(campos, "", "  ")
}

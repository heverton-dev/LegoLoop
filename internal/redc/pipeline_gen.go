package redc

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// GenerateMethodPipelines gera micro-pipelines executáveis para os métodos detectados.
func GenerateMethodPipelines(rootDir, canalSlug string, clusters []MethodCluster) ([]string, error) {
	var created []string
	baseDir := filepath.Join(rootDir, "docs", "channels", canalSlug, "pipelines")
	for _, c := range clusters {
		pDir := filepath.Join(baseDir, c.Slug)
		if err := os.MkdirAll(pDir, 0755); err != nil {
			return nil, fmt.Errorf("redc: erro ao criar pasta de pipeline %s: %w", c.Slug, err)
		}
		dagPath, err := writeMethodDAG(pDir, c)
		if err != nil {
			return nil, err
		}
		created = append(created, dagPath)
	}
	return created, nil
}

func writeMethodDAG(dir string, c MethodCluster) (string, error) {
	dag := buildMethodDAGPayload(c)
	b, err := json.MarshalIndent(dag, "", "  ")
	if err != nil {
		return "", fmt.Errorf("redc: erro ao serializar dag do metodo %s: %w", c.Slug, err)
	}
	dagPath := filepath.Join(dir, "pipeline.dag.json")
	if err := os.WriteFile(dagPath, b, 0644); err != nil {
		return "", fmt.Errorf("redc: erro ao gravar %s: %w", dagPath, err)
	}
	return dagPath, nil
}

func buildMethodDAGPayload(c MethodCluster) map[string]any {
	return map[string]any{
		"version": "1.0.0",
		"name":    fmt.Sprintf("pipeline-%s", c.Slug),
		"nodes": []map[string]any{
			{
				"id":     "N01_ENTRADA",
				"type":   "unit_task",
				"label":  fmt.Sprintf("Ingestão de Parâmetros para %s", c.Metodo),
				"inputs": []string{"contexto"},
			},
			{
				"id":         "N02_EXECUCAO",
				"type":       "unit_task",
				"label":      fmt.Sprintf("Execução da Metodologia %s", c.Metodo),
				"depends_on": []string{"N01_ENTRADA"},
				"outputs":    []string{"resultado.md"},
			},
		},
	}
}

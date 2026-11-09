package redc

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WriteREPC persiste determinísticamente os artefatos de DNA na pasta repc/.
func WriteREPC(rootDir string, dna ChannelDNA) ([]string, error) {
	repcDir := filepath.Join(rootDir, "docs", "channels", dna.CanalSlug, "repc")
	if err := os.MkdirAll(repcDir, 0755); err != nil {
		return nil, fmt.Errorf("redc: falha ao criar diretorio repc: %w", err)
	}

	var created []string
	if err := writeJSONArtifact(repcDir, "01_dossie_dna.json", dna, &created); err != nil {
		return nil, err
	}
	if err := writeJSONArtifact(repcDir, "02_design_tokens.json", dna.DesignTokens, &created); err != nil {
		return nil, err
	}
	if err := writeJSONArtifact(repcDir, "03_motion_tokens.json", dna.MotionTokens, &created); err != nil {
		return nil, err
	}
	if err := writeGanchosMD(repcDir, dna, &created); err != nil {
		return nil, err
	}
	return created, nil
}

func writeJSONArtifact(dir, filename string, data any, created *[]string) error {
	path := filepath.Join(dir, filename)
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("redc: falha ao serializar %s: %w", filename, err)
	}
	if err := os.WriteFile(path, b, 0644); err != nil {
		return fmt.Errorf("redc: falha ao escrever %s: %w", filename, err)
	}
	*created = append(*created, path)
	return nil
}

func writeGanchosMD(dir string, dna ChannelDNA, created *[]string) error {
	path := filepath.Join(dir, "04_matriz_ganchos.md")
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# 🎣 Matriz de Ganchos Virais — Canal: %s\n\n", dna.CanalSlug))
	sb.WriteString("> Fórmulas comprovadas de retenção nos primeiros 3 segundos.\n\n")
	for i, g := range dna.TopGanchos {
		sb.WriteString(fmt.Sprintf("%d. **Gancho:** \"%s\"\n", i+1, g))
	}
	if err := os.WriteFile(path, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("redc: falha ao escrever matriz de ganchos: %w", err)
	}
	*created = append(*created, path)
	return nil
}

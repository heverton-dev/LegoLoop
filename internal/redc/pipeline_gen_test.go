package redc_test

import (
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/redc"
)

func TestGenerateMethodPipelines_Success(t *testing.T) {
	tempDir := t.TempDir()
	clusters := []redc.MethodCluster{
		{
			Slug:        "engenharia-reversa-viral",
			Metodo:      "Engenharia Reversa Viral",
			ClusterTema: "crescimento",
			VideoIDs:    []string{"v1", "v2"},
			ScoreTotal:  2500,
		},
	}

	created, err := redc.GenerateMethodPipelines(tempDir, "canal-alvo", clusters)
	if err != nil {
		t.Fatalf("esperava sucesso ao gerar pipelines de metodos, obteve: %v", err)
	}

	if len(created) == 0 {
		t.Errorf("esperava arquivos de pipeline criados")
	}

	expectedDag := filepath.Join(tempDir, "docs", "channels", "canal-alvo", "pipelines", "engenharia-reversa-viral", "pipeline.dag.json")
	if _, err := os.Stat(expectedDag); os.IsNotExist(err) {
		t.Errorf("arquivo esperado pipeline.dag.json nao foi encontrado: %s", expectedDag)
	}
}

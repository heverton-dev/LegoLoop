package redc_test

import (
	"legoloop/internal/caminhos"
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
		t.Fatalf("esperava sucesso ao gerar Processos de metodos, obteve: %v", err)
	}

	if len(created) == 0 {
		t.Errorf("esperava arquivos de Processo criados")
	}

	expectedDag := filepath.Join(tempDir, "docs", "channels", "canal-alvo", "processos", "engenharia-reversa-viral", caminhos.ArquivoDAG)
	if _, err := os.Stat(expectedDag); os.IsNotExist(err) {
		t.Errorf("arquivo esperado %s nao foi encontrado: %s", caminhos.ArquivoDAG, expectedDag)
	}
}

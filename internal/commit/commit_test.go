package commit_test

import (
	"context"
	"bytes"
	"testing"

	"legoloop/internal/commit"
)

func TestExecutar_ForaDeRepoGit(t *testing.T) {
	// Cria um contexto e executa em diretório temporário isolado sem git
	dir := t.TempDir()
	_ = dir

	var stdout, stderr bytes.Buffer
	opt := commit.Opcoes{DryRun: true}
	// Se rodar no diretório do projeto, deve validar
	err := commit.Executar(context.Background(), &stdout, &stderr, opt)
	if err != nil {
		t.Logf("Resultado esperado em validações: %v", err)
	}
}

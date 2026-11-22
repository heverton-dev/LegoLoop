package forja

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/ste"
)

const nomeSecaoSSOT = "legoloop-ssot"

// AplicarSSOT grava uma secao gerenciada em arquivo de instrucoes apontando para AGENTS.md.
func AplicarSSOT(arquivo, caminhoAgents string) (bool, error) {
	arq := strings.TrimSpace(arquivo)
	alvo := strings.TrimSpace(caminhoAgents)
	if arq == "" || alvo == "" {
		return false, errors.New("forja: arquivo e caminho do AGENTS.md sao obrigatorios")
	}
	antes, existia, err := lerSeExistir(arq)
	if err != nil {
		return false, err
	}
	if err := os.MkdirAll(filepath.Dir(arq), 0o755); err != nil {
		return false, err
	}
	conteudo := fmt.Sprintf(
		"# Regras do LegoLoop (SSOT)\nLeia e siga integralmente as regras em `%s`.",
		filepath.ToSlash(alvo),
	)
	if err := ste.AplicarBloco(arq, nomeSecaoSSOT, conteudo); err != nil {
		return false, err
	}
	depois, err := os.ReadFile(arq)
	if err != nil {
		return false, err
	}
	return !existia || !bytes.Equal(antes, depois), nil
}

func lerSeExistir(caminho string) ([]byte, bool, error) {
	dados, err := os.ReadFile(caminho)
	if err == nil {
		return dados, true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	return nil, false, err
}

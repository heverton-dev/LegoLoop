package kraften

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// NomeReceipt é o nome fixo do artefato de prova na raiz do projeto.
const NomeReceipt = "DEPLOY_RECEIPT.json"

// GravarReceipt emite o DEPLOY_RECEIPT.json em dir e devolve o caminho.
func GravarReceipt(r DeployReceipt, dir string) (string, error) {
	corpo, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", err
	}
	caminho := filepath.Join(dir, NomeReceipt)
	if err := os.WriteFile(caminho, corpo, 0o644); err != nil {
		return "", err
	}
	return caminho, nil
}

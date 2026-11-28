package distribuicao

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RegistroEnvio e a linha de auditoria de um envio real ao Postiz (acervo/envios.json).
type RegistroEnvio struct {
	ItemID       string    `json:"item_id"`
	Canal        string    `json:"canal"`
	Visibilidade string    `json:"visibilidade"`
	Modo         string    `json:"modo"`
	EnviadoEm    time.Time `json:"enviado_em"`
	Resposta     string    `json:"resposta"`
}

// CaminhoEnvios devolve acervo/envios.json.
func CaminhoEnvios(projDir string) string {
	return filepath.Join(projDir, "acervo", "envios.json")
}

// RegistrarEnvio acrescenta um registro ao arquivo de auditoria; nunca grava a chave nem o corpo inteiro.
func RegistrarEnvio(projDir string, r RegistroEnvio) error {
	var lista []RegistroEnvio
	dados, err := os.ReadFile(CaminhoEnvios(projDir))
	switch {
	case errors.Is(err, os.ErrNotExist):
	case err != nil:
		return err
	default:
		if err := json.Unmarshal(dados, &lista); err != nil {
			return fmt.Errorf("%w: envios.json: %v", ErrAcervoInvalido, err)
		}
	}
	r.Resposta = trecho(r.Resposta, 120)
	saida, err := json.MarshalIndent(append(lista, r), "", "  ")
	if err != nil {
		return err
	}
	return gravarAtomico(CaminhoEnvios(projDir), append(saida, '\n'))
}

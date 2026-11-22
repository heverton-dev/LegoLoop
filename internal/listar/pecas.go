package listar

import (
	"os"
	"path/filepath"

	"github.com/heverton-dev/LegoLoop/internal/processo"
)

// ResumoContrato é o contrato da peça em tipos de artefato: o que entra, o que sai e os efeitos.
type ResumoContrato struct {
	Tipo      string   `json:"tipo"`
	Naturezas []string `json:"naturezas"`
	Entra     []string `json:"entra"`
	Sai       []string `json:"sai"`
	Efeitos   []string `json:"efeitos"`
}

// ItemPeca é uma peça solta do cofre com contrato declarado. Erro diz por que o contrato não valeu.
type ItemPeca struct {
	Nome     string         `json:"nome"`
	Contrato ResumoContrato `json:"contrato"`
	Erro     string         `json:"erro,omitempty"`
}

// pecasComContrato varre skills/<peca>/contrato.json do cofre e resume cada contrato. A peça sem
// contrato.json não é alvo; o contrato inválido aparece com Erro, como o processo quebrado.
func pecasComContrato(vault string) []ItemPeca {
	ents, err := os.ReadDir(filepath.Join(vault, "skills"))
	if err != nil {
		return nil
	}
	var out []ItemPeca
	for _, e := range ents {
		if !e.IsDir() {
			continue
		}
		caminho := filepath.Join(vault, "skills", e.Name(), "contrato.json")
		if _, err := os.Stat(caminho); err != nil {
			continue
		}
		out = append(out, itemPeca(e.Name(), caminho))
	}
	return out
}

// itemPeca resume o contrato da peça. Contrato inválido não derruba a lista: vira Erro no item.
func itemPeca(nome, caminho string) ItemPeca {
	c, err := processo.ContratoDe(caminho)
	if err != nil {
		return ItemPeca{Nome: nome, Erro: err.Error()}
	}
	return ItemPeca{Nome: nome, Contrato: resumirContrato(c)}
}

// resumirContrato recorta o contrato aos tipos de artefato de entrada, saída e efeitos.
func resumirContrato(c processo.Contrato) ResumoContrato {
	r := ResumoContrato{Tipo: c.Tipo, Naturezas: c.Naturezas, Efeitos: c.Efeitos,
		Entra: make([]string, 0, len(c.Entra)), Sai: make([]string, 0, len(c.Sai))}
	for _, e := range c.Entra {
		r.Entra = append(r.Entra, e.Tipo)
	}
	for _, s := range c.Sai {
		r.Sai = append(r.Sai, s.Tipo)
	}
	return r
}

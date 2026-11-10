package processo

import (
	"encoding/json"
	"fmt"
	"legoloop/internal/caminhos"
	"os"
	"path/filepath"
)

// sufixoProcesso é o fim do arquivo do processo no cofre: loops/<nome>.processo.json, ao lado do DAG e do loop.
const sufixoProcesso = ".processo.json"

// CaminhoNoCofre devolve onde o processo fica no cofre. É o arquivo que listar, refinar e compor leem.
func CaminhoNoCofre(vault, nome string) string {
	return filepath.Join(vault, caminhos.DirLoops, nome+sufixoProcesso)
}

// Marshal valida o processo e o serializa em JSON indentado, com quebra de linha no fim.
// Um processo inválido nunca chega ao disco.
func (p Processo) Marshal() ([]byte, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

// Ler lê um processo gravado e o valida. JSON quebrado ou contrato violado viram ErrInvalido.
func Ler(path string) (Processo, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Processo{}, err
	}
	var p Processo
	if err := json.Unmarshal(b, &p); err != nil {
		return Processo{}, invalido("%s: %v", filepath.Base(path), err)
	}
	if err := p.Validate(); err != nil {
		return Processo{}, fmt.Errorf("%s: %w", filepath.Base(path), err)
	}
	return p, nil
}

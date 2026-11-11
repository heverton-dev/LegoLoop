// Package caminhos centraliza os nomes de pasta e de arquivo do contrato de
// Processo em disco. Nenhum outro pacote grava esses literais: todos leem
// daqui. A leitura aceita os dois dialetos, com preferência para o novo.
package caminhos

import (
	"os"
	"path/filepath"
)

const (
	DirLoops      = "loops"
	ArquivoDAG    = "pipeline.dag.json"
	ChaveAPILoops = "loops"
	SufixoLoop    = ".loop.json"
)

const (
	// DirProcessos é a pasta dos Processos no dialeto novo.
	DirProcessos = "processos"
	// ArquivoProcessoDAG é o arquivo do contrato DAG no dialeto novo.
	ArquivoProcessoDAG = "processo.dag.json"
	// SufixoStudio é a extensão do documento do Studio no dialeto novo.
	SufixoStudio = ".studio.json"
)

func DirLoopsEm(base string) string {
	return filepath.Join(base, DirLoops)
}

// CaminhoDAG resolve o caminho do arquivo DAG dentro da pasta informada.
func CaminhoDAG(dir string) string {
	return filepath.Join(dir, ArquivoDAG)
}

// DirsLeitura lista as pastas de Processos que existem sob a base, na ordem de
// preferência: a do dialeto novo primeiro, depois a do antigo. Pasta ausente
// não entra na lista.
func DirsLeitura(base string) []string {
	var out []string
	for _, nome := range []string{DirProcessos, DirLoops} {
		dir := filepath.Join(base, nome)
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			out = append(out, dir)
		}
	}
	return out
}

// CaminhoDAGLeitura resolve o arquivo DAG dentro da pasta informada. Prefere o
// nome novo quando o arquivo existe; senão devolve o nome antigo, mesmo que
// ele também não exista, para o chamador reportar o erro de abertura.
func CaminhoDAGLeitura(dir string) string {
	novo := filepath.Join(dir, ArquivoProcessoDAG)
	if info, err := os.Stat(novo); err == nil && info.Mode().IsRegular() {
		return novo
	}
	return CaminhoDAG(dir)
}

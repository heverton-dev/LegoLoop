// Package paths é a fachada de resolução única dos nomes de pasta e de arquivo
// do contrato de Processo em disco. Os literais vivem só em internal/caminhos;
// este pacote os reexpõe sem repetir nenhum deles, para que exista um único
// ponto de verdade (AC-E01.1).
package paths

import "legoloop/internal/caminhos"

const (
	// DirLoops é a pasta dos Processos no dialeto antigo.
	DirLoops = caminhos.DirLoops
	// ArquivoDAG é o arquivo do contrato DAG no dialeto antigo.
	ArquivoDAG = caminhos.ArquivoDAG
	// DirProcessos é a pasta dos Processos no dialeto novo.
	DirProcessos = caminhos.DirProcessos
	// ArquivoProcessoDAG é o arquivo do contrato DAG no dialeto novo.
	ArquivoProcessoDAG = caminhos.ArquivoProcessoDAG
)

// PacoteCentral é o caminho, relativo a internal/, do único pacote autorizado
// a declarar os literais de nome.
const PacoteCentral = "caminhos"

// Literais devolve os valores que só o pacote central pode declarar.
func Literais() []string {
	return []string{DirLoops, ArquivoDAG}
}

// DirLoopsEm resolve a pasta dos Processos sob a base informada.
func DirLoopsEm(base string) string {
	return caminhos.DirLoopsEm(base)
}

// CaminhoDAG resolve o caminho do arquivo DAG dentro da pasta informada.
func CaminhoDAG(dir string) string {
	return caminhos.CaminhoDAG(dir)
}

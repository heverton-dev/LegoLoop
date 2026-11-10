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
	// SufixoLoop é a extensão do documento do Studio no dialeto antigo.
	SufixoLoop = caminhos.SufixoLoop
	// DirProcessos é a pasta dos Processos no dialeto novo.
	DirProcessos = caminhos.DirProcessos
	// ArquivoProcessoDAG é o arquivo do contrato DAG no dialeto novo.
	ArquivoProcessoDAG = caminhos.ArquivoProcessoDAG
	// SufixoStudio é a extensão do documento do Studio no dialeto novo.
	SufixoStudio = caminhos.SufixoStudio
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

// DirsLeitura lista as pastas de Processos existentes sob a base, com o
// dialeto novo antes do antigo.
func DirsLeitura(base string) []string {
	return caminhos.DirsLeitura(base)
}

// CaminhoDAGLeitura resolve o arquivo DAG para leitura, com preferência para o
// nome do dialeto novo quando ele existe.
func CaminhoDAGLeitura(dir string) string {
	return caminhos.CaminhoDAGLeitura(dir)
}

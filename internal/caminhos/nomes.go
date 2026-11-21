// Package caminhos centraliza os nomes de pasta e de arquivo do contrato de
// Processo em disco. Nenhum outro pacote grava esses literais: todos leem
// daqui. A leitura aceita os dois dialetos, com preferência para o novo.
package caminhos

import (
	"os"
	"path/filepath"
)

const (
	DirLoops              = "loops"
	ArquivoDAG            = "pipeline.dag.json"
	ChaveAPILoops         = "loops"
	ChaveMetaProcessoID   = "loop_id"
	SufixoLoop            = ".loop.json"
	TipoSubprocessoAntigo = "subloop"
	TemplatePecaAntigo    = "TEMPLATE_PIECE_LOOP.json"
)

const (
	// DirProcessos é a pasta dos Processos no dialeto novo.
	DirProcessos = "processos"
	// ArquivoProcessoDAG é o arquivo do contrato DAG no dialeto novo.
	ArquivoProcessoDAG = "processo.dag.json"
	// ArquivoToolsDoProcesso é o manifesto de tools que um Processo expõe via router.
	ArquivoToolsDoProcesso = "tools.json"
	// SufixoStudio é a extensão do documento do Studio no dialeto novo.
	SufixoStudio = ".studio.json"
	// DirAntigo é o apelido canônico da pasta dos Processos no dialeto antigo.
	DirAntigo = DirLoops
	// SufixoAntigo é o apelido canônico da extensão do Studio no dialeto antigo.
	SufixoAntigo = SufixoLoop
	// PacoteCentral é o caminho, relativo a internal/, do único pacote autorizado
	// a declarar os literais de nome.
	PacoteCentral = "caminhos"
)

// Literais devolve os valores que só o pacote central pode declarar.
func Literais() []string {
	return []string{DirAntigo, ArquivoDAG}
}

func DirLoopsEm(base string) string {
	return filepath.Join(base, DirLoops)
}

// DirAntigoEm resolve a pasta dos Processos no dialeto antigo sob a base informada.
func DirAntigoEm(base string) string {
	return DirLoopsEm(base)
}

// DirProcessosEm resolve a pasta dos Processos no dialeto de gravação sob a base.
func DirProcessosEm(base string) string {
	return filepath.Join(base, DirProcessos)
}

// CaminhoDAG resolve o caminho do arquivo DAG dentro da pasta informada.
func CaminhoDAG(dir string) string {
	return filepath.Join(dir, ArquivoDAG)
}

// CaminhoToolsDoProcesso resolve o manifesto de tools dentro da pasta do Processo.
func CaminhoToolsDoProcesso(dir string) string {
	return filepath.Join(dir, ArquivoToolsDoProcesso)
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

// DirAcervoLeitura devolve a pasta de Processos a abrir como acervo: a do
// dialeto novo quando existe, depois a do antigo, e a nova se nenhuma existe.
func DirAcervoLeitura(base string) string {
	if dirs := DirsLeitura(base); len(dirs) > 0 {
		return dirs[0]
	}
	return DirProcessosEm(base)
}

// DestinosDoProcesso diz onde gravar o DAG e o documento do Studio de um processo do cofre. O dialeto de
// gravação é processos/: se a pasta processos/<nome>/ existe (layout do fabricar), os dois ficam nela;
// senão, ficam soltos em processos/ (layout antigo de arquivos soltos).
func DestinosDoProcesso(vault, nome, sufixoSolto string) (dag, studio string) {
	proc := DirProcessosEm(vault)
	pasta := filepath.Join(proc, nome)
	if info, err := os.Stat(pasta); err == nil && info.IsDir() {
		return CaminhoDAG(pasta), filepath.Join(pasta, nome+SufixoStudio)
	}
	return filepath.Join(proc, nome+".dag.json"), filepath.Join(proc, nome+sufixoSolto)
}

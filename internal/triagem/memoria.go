package triagem

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// ArquivoMemoria é o nome do arquivo de memória de curto prazo por processo.
const ArquivoMemoria = "memoria_execucao.json"

// MemoriaExecucao registra o estado atual de execução de um processo.
type MemoriaExecucao struct {
	Processo     string    `json:"processo"`
	EtapaAtual   string    `json:"etapa_atual"`
	ProximoPasso string    `json:"proximo_passo"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

// LerMemoria lê a memória de curto prazo do diretório informado.
func LerMemoria(dir string) (MemoriaExecucao, error) {
	dados, err := os.ReadFile(filepath.Join(dir, ArquivoMemoria))
	if err != nil {
		return MemoriaExecucao{}, err
	}
	var mem MemoriaExecucao
	if err := json.Unmarshal(dados, &mem); err != nil {
		return MemoriaExecucao{}, err
	}
	return mem, nil
}

// SalvarMemoria persiste a memória de curto prazo no diretório informado.
func SalvarMemoria(dir string, mem MemoriaExecucao) error {
	if mem.AtualizadoEm.IsZero() {
		mem.AtualizadoEm = time.Now().UTC()
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	dados, err := json.MarshalIndent(mem, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, ArquivoMemoria), dados, 0o644)
}

// LimparMemoria remove o arquivo de memória de curto prazo do diretório informado.
func LimparMemoria(dir string) error {
	err := os.Remove(filepath.Join(dir, ArquivoMemoria))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

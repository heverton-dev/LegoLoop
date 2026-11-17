package triagem

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ArquivoMemoria é o nome do arquivo de memória de curto prazo por processo.
const ArquivoMemoria = "memoria_execucao.json"

// ArquivoLockMemoria é o nome do arquivo de lock da memória do processo.
const ArquivoLockMemoria = "memoria.lock"

// SessaoLegado identifica o estado gravado sem chave de sessão completa.
const SessaoLegado = "legado"

// ErrSessaoInvalida sinaliza sessão sem ID ou sem PID válido.
var ErrSessaoInvalida = errors.New("triagem: sessão inválida")

// ErrLockOcupado sinaliza lock retido por outra sessão além do tempo limite.
var ErrLockOcupado = errors.New("triagem: lock ocupado")

// timeoutLockMemoria limita a espera pelo lock de arquivo (padrão internal/board).
var timeoutLockMemoria = 10 * time.Second

// TokenBudgetTelemetry registra a atuação do gate mecânico de token_budget.
type TokenBudgetTelemetry struct {
	NodeID          string `json:"node_id"`
	EstimatedTokens int    `json:"estimated_tokens"`
	MaxPromptTokens int    `json:"max_prompt_tokens"`
	Action          string `json:"action"` // "pass", "truncate", "abort"
}

// MemoriaExecucao registra o estado atual de execução de um processo.
type MemoriaExecucao struct {
	Processo     string                `json:"processo"`
	SessaoID     string                `json:"sessao_id"` // UUID da sessão
	PID          int                   `json:"pid"`       // PID do detentor do lock
	EtapaAtual   string                `json:"etapa_atual"`
	ProximoPasso string                `json:"proximo_passo"`
	AtualizadoEm time.Time             `json:"atualizado_em"`
	TokenBudget  *TokenBudgetTelemetry `json:"token_budget,omitempty"`
}

// EstadoMemoria agrupa o estado por sessão persistido em ArquivoMemoria.
type EstadoMemoria struct {
	Sessoes map[string]MemoriaExecucao `json:"sessoes"`
}

// Validate verifica a chave de sessão: exige ID não vazio e PID positivo.
func (m MemoriaExecucao) Validate() error {
	if m.SessaoID == "" || m.PID <= 0 {
		return ErrSessaoInvalida
	}
	return nil
}

// lockMemoria retém o lock de arquivo da memória de um diretório.
type lockMemoria struct {
	arquivo *os.File
	caminho string
}

// liberar solta o lock de arquivo e remove o arquivo de lock.
func (l *lockMemoria) liberar() {
	if l == nil || l.arquivo == nil {
		return
	}
	_ = l.arquivo.Close()
	for i := 0; i < 25; i++ {
		err := os.Remove(l.caminho)
		if err == nil || errors.Is(err, os.ErrNotExist) {
			return
		}
		time.Sleep(2 * time.Millisecond)
	}
}

// lockOcupado informa se o erro de abertura marca lock retido por outrem;
// no Windows o handle aberto sem compartilhamento devolve ErrPermission.
func lockOcupado(err error) bool {
	return errors.Is(err, os.ErrExist) || errors.Is(err, os.ErrPermission)
}

// adquirirLockMemoria cria o lock com o PID da sessão e espera até o tempo
// limite, no padrão de lock de arquivo do internal/board.
func adquirirLockMemoria(dir string, pid int, sessao string) (*lockMemoria, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	caminho := filepath.Join(dir, ArquivoLockMemoria)
	prazo := time.Now().Add(timeoutLockMemoria)
	for {
		f, err := os.OpenFile(caminho, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err == nil {
			_, _ = fmt.Fprintf(f, "%d %s", pid, sessao)
			return &lockMemoria{arquivo: f, caminho: caminho}, nil
		}
		if !lockOcupado(err) {
			return nil, err
		}
		if time.Now().After(prazo) {
			return nil, fmt.Errorf("%w: %s", ErrLockOcupado, ArquivoLockMemoria)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

// chaveSessao devolve a chave de estado da memória; sem identidade completa
// usa a sessão legado e com identidade inválida rejeita sem gravar.
func chaveSessao(mem MemoriaExecucao) (string, error) {
	if mem.SessaoID == "" && mem.PID == 0 {
		return SessaoLegado, nil
	}
	if err := mem.Validate(); err != nil {
		return "", err
	}
	return mem.SessaoID, nil
}

// lerEstadoDisco carrega o envelope por sessão; sem arquivo devolve vazio.
func lerEstadoDisco(dir string) (EstadoMemoria, error) {
	vazio := EstadoMemoria{Sessoes: map[string]MemoriaExecucao{}}
	dados, err := os.ReadFile(filepath.Join(dir, ArquivoMemoria))
	if errors.Is(err, os.ErrNotExist) {
		return vazio, nil
	}
	if err != nil {
		return vazio, err
	}
	return desserializarEstado(dados)
}

// desserializarEstado interpreta o envelope por sessão ou o formato legado.
func desserializarEstado(dados []byte) (EstadoMemoria, error) {
	est := EstadoMemoria{Sessoes: map[string]MemoriaExecucao{}}
	var envelope EstadoMemoria
	if err := json.Unmarshal(dados, &envelope); err != nil {
		return est, err
	}
	if envelope.Sessoes != nil {
		est.Sessoes = envelope.Sessoes
		return est, nil
	}
	var mem MemoriaExecucao
	if err := json.Unmarshal(dados, &mem); err != nil {
		return est, err
	}
	est.Sessoes[SessaoLegado] = mem
	return est, nil
}

// escreverAtomico grava em arquivo temporário e renomeia para o destino.
func escreverAtomico(caminho string, dados []byte) error {
	tmp := caminho + ".tmp"
	if err := os.WriteFile(tmp, dados, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, caminho)
}

// persistirEstado grava o envelope completo por sessão com escrita atômica.
func persistirEstado(dir string, est EstadoMemoria) error {
	dados, err := json.MarshalIndent(est, "", "  ")
	if err != nil {
		return err
	}
	return escreverAtomico(filepath.Join(dir, ArquivoMemoria), dados)
}

// gravarEstado atualiza a sessão no envelope e grava com escrita atômica.
func gravarEstado(dir, chave string, mem MemoriaExecucao) error {
	est, err := lerEstadoDisco(dir)
	if err != nil {
		return err
	}
	est.Sessoes[chave] = mem
	return persistirEstado(dir, est)
}

// SalvarMemoria persiste a memória sob lock de arquivo com PID por sessão.
func SalvarMemoria(dir string, mem MemoriaExecucao) error {
	chave, err := chaveSessao(mem)
	if err != nil {
		return err
	}
	if mem.AtualizadoEm.IsZero() {
		mem.AtualizadoEm = time.Now().UTC()
	}
	lk, err := adquirirLockMemoria(dir, mem.PID, mem.SessaoID)
	if err != nil {
		return err
	}
	defer lk.liberar()
	return gravarEstado(dir, chave, mem)
}

// LerMemoria lê o estado da memória; com sessaoID devolve só aquela sessão.
func LerMemoria(dir string, sessaoID ...string) (MemoriaExecucao, error) {
	chave := SessaoLegado
	if len(sessaoID) > 0 {
		chave = sessaoID[0]
	}
	lk, err := adquirirLockMemoria(dir, os.Getpid(), chave)
	if err != nil {
		return MemoriaExecucao{}, err
	}
	defer lk.liberar()
	est, err := lerEstadoDisco(dir)
	if err != nil {
		return MemoriaExecucao{}, err
	}
	mem, ok := est.Sessoes[chave]
	if !ok {
		return MemoriaExecucao{}, fmt.Errorf("triagem: sessão %q ausente em %s: %w",
			chave, ArquivoMemoria, os.ErrNotExist)
	}
	return mem, nil
}

// LimparMemoria remove o arquivo de memória de curto prazo do diretório informado.
func LimparMemoria(dir string) error {
	err := os.Remove(filepath.Join(dir, ArquivoMemoria))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// ZerarSessao zera e remove a sessão dona do envelope da memória do Processo.
func ZerarSessao(dir, sessaoID string) error {
	if sessaoID == "" {
		return ErrSessaoInvalida
	}
	if err := limparLockOrfao(dir); err != nil {
		return err
	}
	lk, err := adquirirLockMemoria(dir, os.Getpid(), sessaoID)
	if err != nil {
		return err
	}
	defer lk.liberar()
	est, err := lerEstadoDisco(dir)
	if err != nil {
		return err
	}
	if _, ok := est.Sessoes[sessaoID]; !ok {
		return fmt.Errorf("triagem: sessão %q inexistente em %s: %w",
			sessaoID, ArquivoMemoria, os.ErrNotExist)
	}
	delete(est.Sessoes, sessaoID)
	return persistirEstado(dir, est)
}

// LimparOrfaos remove sessões com PID morto e o lock órfão do Processo.
func LimparOrfaos(dir string) (int, error) {
	if err := limparLockOrfao(dir); err != nil {
		return 0, err
	}
	lk, err := adquirirLockMemoria(dir, os.Getpid(), SessaoLegado)
	if err != nil {
		return 0, err
	}
	defer lk.liberar()
	est, err := lerEstadoDisco(dir)
	if err != nil {
		return 0, err
	}
	removidas := limparSessoesMortas(est)
	if removidas == 0 {
		return 0, nil
	}
	return removidas, persistirEstado(dir, est)
}

// limparSessoesMortas remove do envelope as sessões sem PID vivo.
func limparSessoesMortas(est EstadoMemoria) int {
	removidas := 0
	for chave, mem := range est.Sessoes {
		if !pidVivo(mem.PID) {
			delete(est.Sessoes, chave)
			removidas++
		}
	}
	return removidas
}

// limparLockOrfao remove o lock de arquivo retido por um PID morto.
func limparLockOrfao(dir string) error {
	caminho := filepath.Join(dir, ArquivoLockMemoria)
	dados, err := os.ReadFile(caminho)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if pidVivo(pidDoLock(dados)) {
		return nil
	}
	return os.Remove(caminho)
}

// pidDoLock extrai o PID gravado no lock de arquivo ("<pid> <sessão>").
func pidDoLock(dados []byte) int {
	campos := strings.Fields(string(dados))
	if len(campos) == 0 {
		return 0
	}
	pid, err := strconv.Atoi(campos[0])
	if err != nil {
		return 0
	}
	return pid
}

//go:build windows

package triagem

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

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
		// Lock órfão (PID morto) é tomado em vez de esperar o tempo limite (AC-E03.2).
		_ = limparLockOrfao(dir)
		if time.Now().After(prazo) {
			return nil, fmt.Errorf("%w: %s", ErrLockOcupado, ArquivoLockMemoria)
		}
		time.Sleep(5 * time.Millisecond)
	}
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

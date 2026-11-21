//go:build !windows

package triagem

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// lockMemoria retém o lock exclusivo (flock) da memória de um diretório.
type lockMemoria struct {
	arquivo *os.File
	caminho string
}

// liberar remove o arquivo de lock enquanto retém o flock e depois o solta.
func (l *lockMemoria) liberar() {
	if l == nil || l.arquivo == nil {
		return
	}
	defer l.arquivo.Close()
	defer syscall.Flock(int(l.arquivo.Fd()), syscall.LOCK_UN)
	if err := os.Remove(l.caminho); err != nil && !errors.Is(err, os.ErrNotExist) {
		return
	}
}

// mesmoArquivoLock confirma que o descritor aberto aponta para o caminho atual.
func mesmoArquivoLock(f *os.File, caminho string) bool {
	fiAberto, err1 := f.Stat()
	fiDisco, err2 := os.Stat(caminho)
	return err1 == nil && err2 == nil && os.SameFile(fiAberto, fiDisco)
}

// lockFlockOcupado informa se o erro de flock indica lock retido por outro dono.
func lockFlockOcupado(err error) bool {
	return errors.Is(err, syscall.EWOULDBLOCK) || errors.Is(err, syscall.EAGAIN)
}

// gravarIdentidadeLock registra PID e sessão no arquivo de lock já travado.
func gravarIdentidadeLock(f *os.File, pid int, sessao string) error {
	if err := f.Truncate(0); err != nil {
		return err
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return err
	}
	_, err := fmt.Fprintf(f, "%d %s", pid, sessao)
	return err
}

// tentarLockPosix tenta adquirir flock exclusivo não-bloqueante no caminho.
func tentarLockPosix(caminho string, pid int, sessao string) (*lockMemoria, bool, error) {
	f, err := os.OpenFile(caminho, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, false, err
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		defer f.Close()
		if lockFlockOcupado(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	if !mesmoArquivoLock(f, caminho) {
		defer f.Close()
		defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		return nil, false, nil
	}
	if err := gravarIdentidadeLock(f, pid, sessao); err != nil {
		defer f.Close()
		defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		return nil, false, err
	}
	return &lockMemoria{arquivo: f, caminho: caminho}, true, nil
}

// adquirirLockMemoria adquire o lock exclusivo com o PID da sessão até o prazo.
func adquirirLockMemoria(dir string, pid int, sessao string) (*lockMemoria, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	caminho := filepath.Join(dir, ArquivoLockMemoria)
	prazo := time.Now().Add(timeoutLockMemoria)
	for {
		lk, ok, err := tentarLockPosix(caminho, pid, sessao)
		if err != nil {
			return nil, err
		}
		if ok {
			return lk, nil
		}
		if time.Now().After(prazo) {
			return nil, fmt.Errorf("%w: %s", ErrLockOcupado, ArquivoLockMemoria)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

// limparLockOrfao remove o arquivo de lock órfão quando ninguém retém o flock
// e o PID registrado no arquivo não está mais vivo.
func limparLockOrfao(dir string) error {
	caminho := filepath.Join(dir, ArquivoLockMemoria)
	f, err := os.OpenFile(caminho, os.O_RDWR, 0)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		if lockFlockOcupado(err) {
			return nil
		}
		return err
	}
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	dados, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	if pidVivo(pidDoLock(dados)) {
		return nil
	}
	if err := os.Remove(caminho); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

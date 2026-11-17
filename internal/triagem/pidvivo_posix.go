//go:build !windows

package triagem

import (
	"os"
	"syscall"
)

// pidVivo informa se o PID ainda tem processo ativo; no POSIX o sinal 0
// devolve erro quando o processo não existe mais.
func pidVivo(pid int) bool {
	if pid <= 0 {
		return false
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	defer p.Release()
	return p.Signal(syscall.Signal(0)) == nil
}

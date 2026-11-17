//go:build windows

package triagem

import "os"

// pidVivo informa se o PID ainda tem processo ativo; no Windows a abertura
// do processo falha quando o PID não existe mais.
func pidVivo(pid int) bool {
	if pid <= 0 {
		return false
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	_ = p.Release()
	return true
}

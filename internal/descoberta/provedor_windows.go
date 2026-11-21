//go:build windows

package descoberta

import (
	"syscall"
	"unsafe"
)

const (
	driveRemovable = 2
	driveFixed     = 3
	driveRemote    = 4
)

type provedorWindows struct{}

// ProvedorPadrao devolve o provedor real somente leitura para Windows.
func ProvedorPadrao() Provedor {
	return provedorWindows{}
}

func (provedorWindows) SO() string { return "windows" }

func (provedorWindows) ListarUnidades() []Unidade {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getDrives := kernel32.NewProc("GetLogicalDrives")
	getType := kernel32.NewProc("GetDriveTypeW")
	mask, _, _ := getDrives.Call()
	var out []Unidade
	for i := 0; i < 26; i++ {
		if mask&(1<<uint(i)) == 0 {
			continue
		}
		letra := string(rune('A'+i)) + ":"
		out = append(out, Unidade{
			Caminho: letra,
			Tipo:    consultarTipoWindows(getType, letra+"\\"),
		})
	}
	return out
}

func consultarTipoWindows(proc *syscall.LazyProc, raiz string) TipoUnidade {
	ptr, err := syscall.UTF16PtrFromString(raiz)
	if err != nil {
		return TipoOutro
	}
	codigo, _, _ := proc.Call(uintptr(unsafe.Pointer(ptr)))
	switch codigo {
	case driveFixed:
		return TipoFixa
	case driveRemovable:
		return TipoRemovivel
	case driveRemote:
		return TipoRede
	default:
		return TipoOutro
	}
}

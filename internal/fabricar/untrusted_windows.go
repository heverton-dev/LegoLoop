//go:build windows

package fabricar

import (
	"syscall"
	"unsafe"
)

// driveFixed é DRIVE_FIXED de GetDriveTypeW: disco interno.
const driveFixed = 3

var getDriveType = syscall.NewLazyDLL("kernel32.dll").NewProc("GetDriveTypeW")

// isUntrustedRoot trata como não confiável tudo que não é disco fixo: removível (USB, cartão), rede,
// CD e RAM disk. Limite conhecido: alguns HDs externos USB se anunciam como fixos.
func isUntrustedRoot(root string) bool {
	p, err := syscall.UTF16PtrFromString(root)
	if err != nil {
		return true
	}
	t, _, _ := getDriveType.Call(uintptr(unsafe.Pointer(p)))
	return t != driveFixed
}

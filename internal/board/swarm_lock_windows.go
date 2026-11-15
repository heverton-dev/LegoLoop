//go:build windows

package board

import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"
)

var (
	lockKernel32     = syscall.NewLazyDLL("kernel32.dll")
	procLockFileEx   = lockKernel32.NewProc("LockFileEx")
	procUnlockFileEx = lockKernel32.NewProc("UnlockFileEx")
)

const (
	lockfileFailImmediately = 0x00000001
	lockfileExclusiveLock   = 0x00000002
)

func acquireFileLock(path string, timeout time.Duration) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	deadline := time.Now().Add(timeout)
	for {
		var ov syscall.Overlapped
		r1, _, callErr := procLockFileEx.Call(
			uintptr(f.Fd()),
			lockfileExclusiveLock|lockfileFailImmediately,
			0, 1, 0,
			uintptr(unsafe.Pointer(&ov)),
		)
		if r1 != 0 {
			return f, nil
		}
		if time.Now().After(deadline) {
			_ = f.Close()
			return nil, fmt.Errorf("board: timeout ao obter lock %s: %w", path, callErr)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func releaseFileLock(f *os.File) {
	if f == nil {
		return
	}
	var ov syscall.Overlapped
	_, _, _ = procUnlockFileEx.Call(uintptr(f.Fd()), 0, 1, 0, uintptr(unsafe.Pointer(&ov)))
	_ = f.Close()
}

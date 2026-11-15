//go:build windows

package board

import (
	"os/exec"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)

var (
	procKernel32           = syscall.NewLazyDLL("kernel32.dll")
	procOpenProcess        = procKernel32.NewProc("OpenProcess")
	procTerminateProcess   = procKernel32.NewProc("TerminateProcess")
	procGetExitCodeProcess = procKernel32.NewProc("GetExitCodeProcess")
	procWaitSingleObject   = procKernel32.NewProc("WaitForSingleObject")
	procCloseHandle        = procKernel32.NewProc("CloseHandle")
)

const (
	winProcessTerminate        = 0x0001
	winProcessSynchronize      = 0x00100000
	winProcessQueryLimitedInfo = 0x1000
	winStillActive             = 259
)

func configureWorkerCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
}

func killOSProcess(pid int) error {
	if pid <= 0 {
		return nil
	}
	_ = exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(pid)).Run()
	if isOSProcessAlive(pid) {
		terminateDirectWindows(pid)
	}
	waitProcessGone(pid, 3*time.Second)
	return nil
}

func terminateDirectWindows(pid int) {
	rights := uintptr(winProcessTerminate | winProcessSynchronize | winProcessQueryLimitedInfo)
	h, _, _ := procOpenProcess.Call(rights, 0, uintptr(pid))
	if h == 0 {
		return
	}
	defer func() { _, _, _ = procCloseHandle.Call(h) }()
	_, _, _ = procTerminateProcess.Call(h, 1)
	_, _, _ = procWaitSingleObject.Call(h, 3000)
}

func isOSProcessAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	h, _, _ := procOpenProcess.Call(uintptr(winProcessQueryLimitedInfo|winProcessSynchronize), 0, uintptr(pid))
	if h == 0 {
		return false
	}
	defer func() { _, _, _ = procCloseHandle.Call(h) }()
	var code uint32
	r1, _, _ := procGetExitCodeProcess.Call(h, uintptr(unsafe.Pointer(&code)))
	if r1 == 0 {
		return false
	}
	return code == winStillActive
}

func waitProcessGone(pid int, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if !isOSProcessAlive(pid) {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
}

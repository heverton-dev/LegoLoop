//go:build !windows

package main

import (
	"os/exec"
	"syscall"
	"time"
)

const esperaPosCancelamento = 3 * time.Second

func configurarCancelamentoArvore(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.WaitDelay = esperaPosCancelamento
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return nil
		}
		if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
			return cmd.Process.Kill()
		}
		return nil
	}
}

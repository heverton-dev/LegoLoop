//go:build windows

package main

import (
	"os/exec"
	"strconv"
	"time"
)

const esperaPosCancelamento = 3 * time.Second

func configurarCancelamentoArvore(cmd *exec.Cmd) {
	cmd.WaitDelay = esperaPosCancelamento
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return nil
		}
		pid := strconv.Itoa(cmd.Process.Pid)
		if err := exec.Command("taskkill", "/F", "/T", "/PID", pid).Run(); err != nil {
			return cmd.Process.Kill()
		}
		return nil
	}
}

package board

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// SetWorkerModel atualiza o modelo de um worker em runtime no quadro persistente.
func (sb *SwarmBoard) SetWorkerModel(name, model string) (Worker, error) {
	model = strings.TrimSpace(model)
	if model == "" {
		return Worker{}, fmt.Errorf("modelo vazio")
	}
	sb.mu.Lock()
	defer sb.mu.Unlock()
	st, err := sb.loadLocked()
	if err != nil {
		return Worker{}, err
	}
	for i := range st.Workers {
		if st.Workers[i].Name == name {
			st.Workers[i].Model = model
			st.Workers[i].LastActiveAt = time.Now().UTC()
			if err := sb.saveLocked(st); err != nil {
				return Worker{}, err
			}
			return st.Workers[i], nil
		}
	}
	return Worker{}, fmt.Errorf("worker %q não encontrado", name)
}

// SpawnWorker inicia o processo OS real de um worker e registra seu PID no quadro.
func (sb *SwarmBoard) SpawnWorker(ctx context.Context, name string, cmdArgs []string) (int, error) {
	if len(cmdArgs) == 0 {
		return 0, fmt.Errorf("comando vazio para spawn")
	}
	workDir, err := sb.workerWorkDir(name)
	if err != nil {
		return 0, err
	}
	cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
	cmd.Dir = workDir
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	pid := cmd.Process.Pid
	if err := sb.recordSpawn(name, pid); err != nil {
		return pid, err
	}
	go sb.waitWorkerExit(name, cmd)
	return pid, nil
}

func (sb *SwarmBoard) workerWorkDir(name string) (string, error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	st, err := sb.loadLocked()
	if err != nil {
		return "", err
	}
	for _, w := range st.Workers {
		if w.Name == name {
			if w.WorktreePath != "" {
				return w.WorktreePath, nil
			}
			return sb.dir, nil
		}
	}
	return "", fmt.Errorf("worker %q não encontrado", name)
}

func (sb *SwarmBoard) recordSpawn(name string, pid int) error {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	st, err := sb.loadLocked()
	if err != nil {
		return err
	}
	for i := range st.Workers {
		if st.Workers[i].Name == name {
			st.Workers[i].PID = pid
			st.Workers[i].Status = WorkerWorking
			st.Workers[i].LastActiveAt = time.Now().UTC()
			return sb.saveLocked(st)
		}
	}
	return fmt.Errorf("worker %q não encontrado", name)
}

func (sb *SwarmBoard) waitWorkerExit(name string, cmd *exec.Cmd) {
	_ = cmd.Wait()
	_ = sb.TouchWorker(name, WorkerStopped, time.Now().UTC())
}

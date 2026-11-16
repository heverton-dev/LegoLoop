package board

import (
	"context"
	"fmt"
	"os"
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
	var updated Worker
	err := sb.withLock(func() error {
		st, err := sb.loadLocked()
		if err != nil {
			return err
		}
		return updateWorkerModel(&st, name, model, sb, &updated)
	})
	return updated, err
}

func updateWorkerModel(st *swarmState, name, model string, sb *SwarmBoard, out *Worker) error {
	for i := range st.Workers {
		if st.Workers[i].Name != name {
			continue
		}
		if st.Workers[i].Status == WorkerKilled {
			return fmt.Errorf("worker %q ja encerrado", name)
		}
		st.Workers[i].Model = model
		st.Workers[i].LastActiveAt = time.Now().UTC()
		if err := sb.saveLocked(*st); err != nil {
			return err
		}
		*out = st.Workers[i]
		return nil
	}
	return fmt.Errorf("worker %q nao encontrado", name)
}

// SpawnWorker inicia o processo OS real de um worker e registra seu PID no quadro.
func (sb *SwarmBoard) SpawnWorker(ctx context.Context, name string, cmdArgs []string) (int, error) {
	if len(cmdArgs) == 0 {
		return 0, fmt.Errorf("comando vazio para spawn")
	}
	workDir, model, err := sb.prepareSpawnSpec(name)
	if err != nil {
		return 0, err
	}
	cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
	cmd.Dir = workDir
	cmd.Env = buildWorkerEnv(name, model, sb.dir)
	configureWorkerCmd(cmd)
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	pid := cmd.Process.Pid
	done := sb.trackActive(pid)
	if err := sb.recordSpawn(name, pid); err != nil {
		_ = killOSProcess(pid)
		_ = cmd.Wait()
		sb.untrackActive(pid, done)
		return 0, err
	}
	go sb.waitWorkerExit(name, pid, cmd, done)
	return pid, nil
}

func buildWorkerEnv(name, model, boardDir string) []string {
	env := append([]string{}, os.Environ()...)
	env = append(env, "LEGOLOOP_SWARM_WORKER="+name, "LEGOLOOP_SWARM_BOARD="+boardDir)
	if model != "" {
		env = append(env, "LEGOLOOP_SWARM_MODEL="+model)
	}
	return env
}

func (sb *SwarmBoard) prepareSpawnSpec(name string) (string, string, error) {
	var workDir, model string
	err := sb.withLock(func() error {
		st, err := sb.loadLocked()
		if err != nil {
			return err
		}
		for _, w := range st.Workers {
			if w.Name != name {
				continue
			}
			if w.Status == WorkerKilled {
				return fmt.Errorf("worker %q ja encerrado", name)
			}
			if w.PID > 0 && isOSProcessAlive(w.PID) {
				return fmt.Errorf("worker %q ja possui processo ativo (%d)", name, w.PID)
			}
			workDir = pickWorkerWorkDir(w, sb.projDir)
			model = w.Model
			return nil
		}
		return fmt.Errorf("worker %q nao encontrado", name)
	})
	return workDir, model, err
}

func pickWorkerWorkDir(w Worker, projDir string) string {
	if w.WorktreePath != "" {
		return w.WorktreePath
	}
	if w.RepoDir != "" {
		return w.RepoDir
	}
	return projDir
}

func (sb *SwarmBoard) recordSpawn(name string, pid int) error {
	return sb.withLock(func() error {
		st, err := sb.loadLocked()
		if err != nil {
			return err
		}
		for i := range st.Workers {
			if st.Workers[i].Name == name {
				st.Workers[i].PID = pid
				st.Workers[i].ExitCode = nil
				st.Workers[i].Status = WorkerWorking
				st.Workers[i].LastActiveAt = time.Now().UTC()
				return sb.saveLocked(st)
			}
		}
		return fmt.Errorf("worker %q nao encontrado", name)
	})
}

func (sb *SwarmBoard) trackActive(pid int) chan struct{} {
	ch := make(chan struct{})
	sb.actMu.Lock()
	sb.active[pid] = ch
	sb.actMu.Unlock()
	return ch
}

func (sb *SwarmBoard) untrackActive(pid int, done chan struct{}) {
	sb.actMu.Lock()
	delete(sb.active, pid)
	sb.actMu.Unlock()
	close(done)
}

func (sb *SwarmBoard) waitActiveExit(pid int, timeout time.Duration) {
	sb.actMu.Lock()
	ch, ok := sb.active[pid]
	sb.actMu.Unlock()
	if !ok {
		return
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-ch:
	case <-timer.C:
	}
}

func (sb *SwarmBoard) waitWorkerExit(name string, pid int, cmd *exec.Cmd, done chan struct{}) {
	defer sb.untrackActive(pid, done)
	_ = cmd.Wait()
	code := -1
	if cmd.ProcessState != nil {
		code = cmd.ProcessState.ExitCode()
	}
	_ = sb.finishWorkerExit(name, pid, code)
}

func (sb *SwarmBoard) finishWorkerExit(name string, pid, code int) error {
	return sb.withLock(func() error {
		st, err := sb.loadLocked()
		if err != nil {
			return err
		}
		for i := range st.Workers {
			w := &st.Workers[i]
			if w.Name != name {
				continue
			}
			switch {
			case w.Status == WorkerKilled:
				return nil
			case w.PID == pid:
				w.Status = WorkerStopped
				w.PID = 0
			case w.PID == 0 && w.Status == WorkerStopped && w.ExitCode == nil:
				// Status() ja marcou o PID morto como stopped; falta so o codigo de saida.
			default:
				return nil
			}
			c := code
			w.ExitCode = &c
			w.LastActiveAt = time.Now().UTC()
			return sb.saveLocked(st)
		}
		return nil
	})
}

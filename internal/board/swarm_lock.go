package board

import (
	"os"
	"path/filepath"
	"time"
)

const fileLockTimeout = 10 * time.Second

func (sb *SwarmBoard) withLock(fn func() error) error {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	if err := os.MkdirAll(sb.dir, 0o755); err != nil {
		return err
	}
	f, err := acquireFileLock(filepath.Join(sb.dir, "swarm.lock"), fileLockTimeout)
	if err != nil {
		return err
	}
	defer releaseFileLock(f)
	return fn()
}

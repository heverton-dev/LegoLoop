package board

import (
	"context"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestSetWorkerModelAndSpawnWorker(t *testing.T) {
	boardRoot := t.TempDir()
	sb, err := OpenSwarmBoard(boardRoot)
	if err != nil {
		t.Fatalf("OpenSwarmBoard: %v", err)
	}
	w, err := sb.AddWorker(AddWorkerRequest{Name: "w1", Role: "engenheiro-tdd", Model: "m1"})
	if err != nil {
		t.Fatalf("AddWorker: %v", err)
	}
	if w.Model != "m1" {
		t.Fatalf("expected model m1, got %q", w.Model)
	}

	updated, err := sb.SetWorkerModel("w1", "m2")
	if err != nil {
		t.Fatalf("SetWorkerModel: %v", err)
	}
	if updated.Model != "m2" {
		t.Fatalf("expected model m2, got %q", updated.Model)
	}

	cmdArgs := []string{"sh", "-c", "sleep 5"}
	if runtime.GOOS == "windows" {
		cmdArgs = []string{"cmd", "/c", "ping -n 6 127.0.0.1 >nul"}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pid, err := sb.SpawnWorker(ctx, "w1", cmdArgs)
	if err != nil {
		t.Fatalf("SpawnWorker: %v", err)
	}
	if pid <= 0 {
		t.Fatalf("expected positive pid, got %d", pid)
	}
	t.Cleanup(func() { _ = sb.KillWorker("w1", "") })

	if err := sb.KillWorker("w1", ""); err != nil {
		t.Fatalf("KillWorker: %v", err)
	}
	if isOSProcessAlive(pid) {
		t.Fatalf("processo %d ainda ativo apos KillWorker", pid)
	}
	if err := os.RemoveAll(boardRoot); err != nil {
		t.Fatalf("vazamento impediu remocao de boardRoot: %v", err)
	}
}

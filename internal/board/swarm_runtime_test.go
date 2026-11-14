package board

import (
	"context"
	"runtime"
	"testing"
	"time"
)

func TestSetWorkerModelAndSpawnWorker(t *testing.T) {
	sb, err := OpenSwarmBoard(t.TempDir())
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

	cmdArgs := []string{"sh", "-c", "sleep 0.1"}
	if runtime.GOOS == "windows" {
		cmdArgs = []string{"cmd", "/c", "ping -n 1 127.0.0.1 >nul"}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pid, err := sb.SpawnWorker(ctx, "w1", cmdArgs)
	if err != nil {
		t.Fatalf("SpawnWorker: %v", err)
	}
	if pid <= 0 {
		t.Fatalf("expected positive pid, got %d", pid)
	}
}

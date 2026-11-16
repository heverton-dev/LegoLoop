package mcp

import (
	"os"
	"testing"
	"time"

	"legoloop/internal/board"
)

func TestSwarmSpawnToolRealProcess(t *testing.T) {
	sb, err := board.OpenSwarmBoard(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := sb.AddWorker(board.AddWorkerRequest{Name: "w1", Role: "engenheiro-tdd"}); err != nil {
		t.Fatal(err)
	}
	if _, _, err := dispatchSwarmTool(sb, "swarm_spawn", map[string]any{"name": "w1"}); err == nil {
		t.Fatal("swarm_spawn sem command deveria falhar")
	}
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	res, handled, err := dispatchSwarmTool(sb, "swarm_spawn", map[string]any{
		"name": "w1", "command": []any{exe, "-test.run=^$"},
	})
	if err != nil || !handled {
		t.Fatalf("swarm_spawn falhou: handled=%v err=%v", handled, err)
	}
	if pid, _ := res.(map[string]any)["pid"].(int); pid <= 0 {
		t.Fatalf("pid invalido: %+v", res)
	}
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		st, _ := sb.Status()
		if len(st.Workers) == 1 && st.Workers[0].Status == board.WorkerStopped && st.Workers[0].PID == 0 {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("worker nao transicionou para stopped apos saida natural")
}

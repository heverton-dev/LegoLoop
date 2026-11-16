package mcp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/board"
)

func writeDAG(t *testing.T, withCommand bool) string {
	t.Helper()
	node := func(id, role string, deps ...string) map[string]any {
		return map[string]any{
			"id": id, "name": id, "type": "skill", "nature": "micro", "cognitive_complexity": "simples",
			"ref": "skills/" + id, "role": role, "depends_on": deps, "timeout_sec": 60,
			"routing": map[string]any{"mode": "suggested", "model": "sonnet", "effort": "medium", "fallback": "code-fast"},
			"output":  map[string]any{"artifact": "out/" + id + ".md", "gate": "file_exists('out/" + id + ".md')"},
		}
	}
	d := map[string]any{
		"$schema": "https://legoloop.dev/schemas/dag.v1.json", "id": "demo", "name": "demo", "version": "1.0.0",
		"trigger": map[string]any{"type": "manual"},
		"nodes":   []any{node("a", "engenheiro-backend"), node("b", "engenheiro-tdd"), node("c", "auditor-qualidade", "a", "b")},
	}
	b, _ := json.Marshal(d)
	p := filepath.Join(t.TempDir(), "pipeline.dag.json")
	if err := os.WriteFile(p, b, 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestSwarmRunToolPorNivelEEncerra(t *testing.T) {
	sb, err := board.OpenSwarmBoard(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	exe, _ := os.Executable()
	res, handled, err := dispatchSwarmTool(sb, "swarm_run", map[string]any{
		"dag": writeDAG(t, true), "workspace": "inherit", "timeout": "30s",
		"command": []any{exe, "-test.run=^$"},
	})
	if err != nil || !handled {
		t.Fatalf("swarm_run falhou: handled=%v err=%v", handled, err)
	}
	levels := res.(map[string]any)["levels"].([]board.RunLevelResult)
	if len(levels) != 2 || len(levels[0].Workers) != 2 || len(levels[1].Workers) != 1 {
		t.Fatalf("niveis inesperados: %+v", levels)
	}
	st, _ := sb.Status()
	if len(st.Workers) != 3 {
		t.Fatalf("esperava 3 workers: %+v", st.Workers)
	}
	for _, w := range st.Workers {
		if w.Status != board.WorkerKilled || w.PID != 0 {
			t.Errorf("worker %s deveria estar killed: %+v", w.Name, w)
		}
	}
}

func TestSwarmRunToolErros(t *testing.T) {
	sb, _ := board.OpenSwarmBoard(t.TempDir())
	if _, _, err := dispatchSwarmTool(sb, "swarm_run", map[string]any{}); err == nil {
		t.Error("sem dag deveria falhar")
	}
	if _, _, err := dispatchSwarmTool(sb, "swarm_run", map[string]any{"dag": writeDAG(t, false)}); err == nil {
		t.Error("no sem comando deveria falhar")
	}
}

// swarm_status somente leitura mostra o mesmo estado que Status() grava: PID morto vira stopped, sem gravar.
func TestSwarmStatusSomenteLeituraReconciliaPIDMorto(t *testing.T) {
	proj := t.TempDir()
	sb, err := board.OpenSwarmBoard(proj)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := sb.AddWorker(board.AddWorkerRequest{Name: "w1", Role: "engenheiro-tdd"}); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(proj, ".legoloop-swarm", "swarm.json")
	raw, _ := os.ReadFile(file)
	var st map[string]any
	_ = json.Unmarshal(raw, &st)
	w := st["workers"].([]any)[0].(map[string]any)
	w["status"], w["pid"] = "working", float64(1<<22+12345) // PID inexistente
	edited, _ := json.Marshal(st)
	if err := os.WriteFile(file, edited, 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := readSwarmStatusReadOnly(proj)
	if err != nil {
		t.Fatal(err)
	}
	if got.Workers[0].Status != board.WorkerStopped || got.Workers[0].PID != 0 {
		t.Errorf("deveria mostrar stopped com PID 0: %+v", got.Workers[0])
	}
	after, _ := os.ReadFile(file)
	if string(after) != string(edited) {
		t.Error("leitura somente leitura nao pode gravar no quadro")
	}
}

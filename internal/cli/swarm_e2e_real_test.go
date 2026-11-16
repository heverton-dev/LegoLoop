package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"legoloop/internal/board"
)

// buildLegoloopBin compila o binario real; o E2E nao usa Run() em processo.
func buildLegoloopBin(t *testing.T) string {
	t.Helper()
	name := "legoloop"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	bin := filepath.Join(t.TempDir(), name)
	cmd := exec.Command("go", "build", "-o", bin, "../../cmd/legoloop")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build legoloop: %v\n%s", err, out)
	}
	return bin
}

func e2eGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	steps := [][]string{
		{"init"},
		{"config", "user.email", "e2e@legoloop.local"},
		{"config", "user.name", "LegoLoop E2E"},
	}
	for _, a := range steps {
		c := exec.Command("git", a...)
		c.Dir = dir
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v (%s)", a, err, out)
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# e2e\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, a := range [][]string{{"add", "."}, {"commit", "-m", "init"}} {
		c := exec.Command("git", a...)
		c.Dir = dir
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v (%s)", a, err, out)
		}
	}
	return dir
}

func runBin(t *testing.T, bin string, args ...string) string {
	t.Helper()
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(bin, args...)
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("legoloop %v: %v\nstdout=%s\nstderr=%s", args, err, stdout.String(), stderr.String())
	}
	return stdout.String()
}

func statusOf(t *testing.T, bin, repo string) map[string]board.Worker {
	t.Helper()
	var st board.SwarmStatus
	if err := json.Unmarshal([]byte(runBin(t, bin, "swarm", "status", "-dir", repo)), &st); err != nil {
		t.Fatalf("status json: %v", err)
	}
	out := make(map[string]board.Worker, len(st.Workers))
	for _, w := range st.Workers {
		out[w.Name] = w
	}
	return out
}

func waitStatus(t *testing.T, bin, repo, name string, want board.WorkerStatus) board.Worker {
	t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		if w := statusOf(t, bin, repo)[name]; w.Status == want {
			return w
		}
		time.Sleep(150 * time.Millisecond)
	}
	t.Fatalf("worker %s nao chegou a %s: %+v", name, want, statusOf(t, bin, repo)[name])
	return board.Worker{}
}

// TestSwarmE2ERealParallelWorkersAndKill prova o enxame sem harness externo: CLI proprio,
// worktrees isolados (workspace share), quadro em disco, 2 workers paralelos e kill limpo.
func TestSwarmE2ERealParallelWorkersAndKill(t *testing.T) {
	if testing.Short() {
		t.Skip("E2E real compila o binario")
	}
	bin := buildLegoloopBin(t)
	repo := e2eGitRepo(t)

	for _, w := range []struct{ name, role string }{{"w1", "engenheiro-backend"}, {"w2", "engenheiro-tdd"}, {"w3", "auditor-qualidade"}} {
		runBin(t, bin, "swarm", "add", "-dir", repo, "-repo", repo, "-name", w.name, "-role", w.role, "-workspace", "share")
	}
	t.Cleanup(func() {
		for _, n := range []string{"w1", "w2", "w3"} {
			_ = exec.Command(bin, "swarm", "kill", "-dir", repo, "-repo", repo, "-name", n).Run()
		}
	})

	// w1 e w2 bloqueiam em "swarm wait" (alvo e quadro herdados do ambiente do spawn).
	for _, n := range []string{"w1", "w2"} {
		runBin(t, bin, "swarm", "spawn", "-dir", repo, "-name", n, "--", bin, "swarm", "wait", "-timeout", "60s")
	}
	// w3 posta no quadro usando apenas o remetente herdado do ambiente.
	runBin(t, bin, "swarm", "spawn", "-dir", repo, "-name", "w3", "--", bin, "swarm", "post", "-to", "@coordenador", "-body", "ola do w3")

	st := waitStatus(t, bin, repo, "w3", board.WorkerStopped)
	if st.PID != 0 {
		t.Fatalf("w3 parado deveria ter PID=0: %+v", st)
	}
	msgs := runBin(t, bin, "swarm", "read", "-dir", repo)
	var read []board.Message
	if err := json.Unmarshal([]byte(msgs), &read); err != nil || len(read) != 1 || read[0].From != "w3" {
		t.Fatalf("mensagem do worker nao veio com remetente herdado: %s err=%v", msgs, err)
	}

	// Paralelismo real: 2 PIDs vivos distintos, cada um em worktree proprio.
	cur := statusOf(t, bin, repo)
	w1, w2 := cur["w1"], cur["w2"]
	if w1.Status != board.WorkerWorking || w2.Status != board.WorkerWorking || w1.PID <= 0 || w2.PID <= 0 || w1.PID == w2.PID {
		t.Fatalf("esperava 2 workers vivos e distintos: w1=%+v w2=%+v", w1, w2)
	}
	if w1.WorktreePath == w2.WorktreePath {
		t.Fatalf("worktrees nao isolados: %s", w1.WorktreePath)
	}
	for _, w := range []board.Worker{w1, w2} {
		if _, err := os.Stat(filepath.Join(w.WorktreePath, "README.md")); err != nil {
			t.Fatalf("worktree de %s sem checkout: %v", w.Name, err)
		}
	}

	// Mensagem endereçada acorda so o w1; o w2 segue vivo.
	runBin(t, bin, "swarm", "post", "-dir", repo, "-to", "@w1", "-body", "go")
	waitStatus(t, bin, repo, "w1", board.WorkerStopped)
	if got := statusOf(t, bin, repo)["w2"]; got.Status != board.WorkerWorking || got.PID != w2.PID {
		t.Fatalf("w2 nao devia ser afetado: %+v", got)
	}

	// Kill limpo: PID morto, status killed, worktree removido do disco e do git.
	runBin(t, bin, "swarm", "kill", "-dir", repo, "-repo", repo, "-name", "w2")
	runBin(t, bin, "swarm", "kill", "-dir", repo, "-repo", repo, "-name", "w1")
	final := statusOf(t, bin, repo)
	for _, w := range []board.Worker{w1, w2} {
		got := final[w.Name]
		if got.Status != board.WorkerKilled || got.PID != 0 {
			t.Fatalf("%s deveria estar killed com PID=0: %+v", w.Name, got)
		}
		if _, err := os.Stat(w.WorktreePath); !os.IsNotExist(err) {
			t.Fatalf("worktree de %s sobrou: err=%v", w.Name, err)
		}
	}
	c := exec.Command("git", "worktree", "list", "--porcelain")
	c.Dir = repo
	out, err := c.Output()
	if err != nil || bytes.Count(out, []byte("worktree ")) != 2 { // principal + w3 (parado, nunca killed)
		t.Fatalf("git worktree list inesperado (err=%v): %s", err, out)
	}
}

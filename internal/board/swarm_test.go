package board

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func initTestGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@legoloop.local"},
		{"git", "config", "user.name", "LegoLoop Test"},
	}
	for _, args := range cmds {
		c := exec.Command(args[0], args[1:]...)
		c.Dir = dir
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("git setup %v: %v (%s)", args, err, out)
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"git", "add", "."}, {"git", "commit", "-m", "init"}} {
		c := exec.Command(args[0], args[1:]...)
		c.Dir = dir
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("git commit %v: %v (%s)", args, err, out)
		}
	}
	return dir
}

func TestSwarmBoardAddAndKillWorktree(t *testing.T) {
	repo := initTestGitRepo(t)
	boardDir := filepath.Join(t.TempDir(), "board")
	sb, err := OpenSwarmBoard(boardDir)
	if err != nil {
		t.Fatalf("OpenSwarmBoard: %v", err)
	}
	if _, err := sb.AddWorker(AddWorkerRequest{Role: "papel-invalido", RepoDir: repo}); err == nil {
		t.Fatal("esperava erro para papel invalido")
	}
	w, err := sb.AddWorker(AddWorkerRequest{
		Role:          "engenheiro-backend",
		WorkspaceMode: WorkspaceShare,
		Model:         "sonnet",
		RepoDir:       repo,
	})
	if err != nil {
		t.Fatalf("AddWorker share: %v", err)
	}
	if w.WorktreePath == "" {
		t.Fatal("WorktreePath vazio para workspace share")
	}
	if _, err := os.Stat(w.WorktreePath); err != nil {
		t.Fatalf("worktree nao existe em disco: %v", err)
	}
	if err := sb.KillWorker(w.Name, repo); err != nil {
		t.Fatalf("KillWorker: %v", err)
	}
	st, err := sb.Status()
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if st.Workers[0].Status != WorkerKilled {
		t.Fatalf("esperava killed, obteve %s", st.Workers[0].Status)
	}
}

func TestSwarmBoardPostReadAndWait(t *testing.T) {
	boardDir := filepath.Join(t.TempDir(), "board")
	sb, err := OpenSwarmBoard(boardDir)
	if err != nil {
		t.Fatalf("OpenSwarmBoard: %v", err)
	}
	m1, err := sb.PostMessage("coordenador", "@engenheiro-backend-1", "iniciar")
	if err != nil {
		t.Fatalf("PostMessage 1: %v", err)
	}
	m2, err := sb.PostMessage("engenheiro-backend-1", "@all", "pronto")
	if err != nil {
		t.Fatalf("PostMessage 2: %v", err)
	}
	sinceMsgs, err := sb.ReadMessages(m1.Seq, nil)
	if err != nil || len(sinceMsgs) != 1 || sinceMsgs[0].Seq != m2.Seq {
		t.Fatalf("ReadMessages since falhou: %+v err=%v", sinceMsgs, err)
	}
	seqMsgs, err := sb.ReadMessages(0, []int64{m1.Seq})
	if err != nil || len(seqMsgs) != 1 || seqMsgs[0].Seq != m1.Seq {
		t.Fatalf("ReadMessages seqs falhou: %+v err=%v", seqMsgs, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	go func() {
		time.Sleep(50 * time.Millisecond)
		_, _ = sb.PostMessage("coordenador", "@engenheiro-tdd-1", "validar")
	}()
	waited, err := sb.WaitMessages(ctx, "@engenheiro-tdd-1", m2.Seq, 20*time.Millisecond)
	if err != nil || len(waited) != 1 {
		t.Fatalf("WaitMessages falhou: %+v err=%v", waited, err)
	}
}

func TestSwarmBoardStallMonitor(t *testing.T) {
	boardDir := filepath.Join(t.TempDir(), "board")
	sb, err := OpenSwarmBoard(boardDir)
	if err != nil {
		t.Fatalf("OpenSwarmBoard: %v", err)
	}
	w, err := sb.AddWorker(AddWorkerRequest{Role: "engenheiro-tdd", WorkspaceMode: WorkspaceInherit})
	if err != nil {
		t.Fatalf("AddWorker: %v", err)
	}
	if err := sb.TouchWorker(w.Name, WorkerWorking, time.Now().Add(-2*time.Minute)); err != nil {
		t.Fatalf("TouchWorker: %v", err)
	}
	stalled, err := sb.CheckStall(30 * time.Second)
	if err != nil || len(stalled) != 1 || stalled[0].Name != w.Name {
		t.Fatalf("CheckStall falhou: %+v err=%v", stalled, err)
	}
}

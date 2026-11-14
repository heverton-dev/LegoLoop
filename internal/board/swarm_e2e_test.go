package board

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestSwarmE2EDAGMultiNoReal(t *testing.T) {
	repo := initTestGitRepo(t)
	boardDir := filepath.Join(t.TempDir(), "swarm-e2e")
	sb, err := OpenSwarmBoard(boardDir)
	if err != nil {
		t.Fatalf("OpenSwarmBoard: %v", err)
	}
	w1, err := sb.AddWorker(AddWorkerRequest{Role: "engenheiro-backend", WorkspaceMode: WorkspaceShare, Model: "sonnet", RepoDir: repo})
	if err != nil {
		t.Fatalf("AddWorker backend: %v", err)
	}
	w2, err := sb.AddWorker(AddWorkerRequest{Role: "engenheiro-tdd", WorkspaceMode: WorkspaceShare, Model: "sonnet", RepoDir: repo})
	if err != nil {
		t.Fatalf("AddWorker tdd: %v", err)
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		_ = os.WriteFile(filepath.Join(w1.WorktreePath, "backend.txt"), []byte("ok\n"), 0o644)
		_, _ = sb.PostMessage(w1.Name, "@coordenador", "nivel-0 backend concluido")
	}()
	go func() {
		defer wg.Done()
		_ = os.WriteFile(filepath.Join(w2.WorktreePath, "tdd.txt"), []byte("ok\n"), 0o644)
		_, _ = sb.PostMessage(w2.Name, "@coordenador", "nivel-0 tdd concluido")
	}()
	wg.Wait()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	msgs, err := sb.WaitMessages(ctx, "@coordenador", 0, 20*time.Millisecond)
	if err != nil || len(msgs) < 2 {
		t.Fatalf("WaitMessages nivel 0 falhou: len=%d err=%v", len(msgs), err)
	}
	w3, err := sb.AddWorker(AddWorkerRequest{Role: "auditor-qualidade", WorkspaceMode: WorkspaceInherit, Model: "opus"})
	if err != nil {
		t.Fatalf("AddWorker auditor: %v", err)
	}
	if _, err := sb.PostMessage(w3.Name, "@all", "nivel-1 convergencia aprovada"); err != nil {
		t.Fatalf("PostMessage auditor: %v", err)
	}
	for _, name := range []string{w1.Name, w2.Name, w3.Name} {
		if err := sb.KillWorker(name, repo); err != nil {
			t.Fatalf("KillWorker %s: %v", name, err)
		}
	}
	st, err := sb.Status()
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if len(st.Workers) != 3 || st.MessageCount != 3 {
		t.Fatalf("estado final inesperado: %+v", st)
	}
}

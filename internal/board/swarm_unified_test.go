package board

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestHelperSwarmUnifiedProcess(t *testing.T) {
	mode := os.Getenv("GO_SWARM_UNIFIED_HELPER")
	if mode == "" {
		return
	}
	switch mode {
	case "sleep":
		if modelOut := os.Getenv("GO_SWARM_MODEL_OUT"); modelOut != "" {
			_ = os.WriteFile(modelOut, []byte(os.Getenv("LEGOLOOP_SWARM_MODEL")), 0o644)
		}
		time.Sleep(30 * time.Second)
	case "post":
		dir := os.Getenv("GO_SWARM_HELPER_DIR")
		from := os.Getenv("GO_SWARM_HELPER_FROM")
		count, _ := strconv.Atoi(os.Getenv("GO_SWARM_HELPER_COUNT"))
		sb, err := OpenSwarmBoard(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "OpenSwarmBoard: %v\n", err)
			os.Exit(2)
		}
		for i := 0; i < count; i++ {
			if _, err := sb.PostMessage(from, "@all", fmt.Sprintf("msg-%d", i)); err != nil {
				fmt.Fprintf(os.Stderr, "PostMessage: %v\n", err)
				os.Exit(3)
			}
		}
	}
	os.Exit(0)
}

func TestUnifiedPathFileLockExactMentionAndPIDKill(t *testing.T) {
	t.Run("CanonicalPathAndLegacyMigration", testUnifiedCanonicalPathAndLegacy)
	t.Run("FileLockCrossProcessAndAtomicUniqueTmp", testUnifiedFileLockCrossProcess)
	t.Run("SanitizeWorkerNameAndWorkspaceModes", testUnifiedNameAndWorkspaceModes)
	t.Run("ExactMentionWordBoundary", testUnifiedExactMentionBoundary)
	t.Run("SpawnSetModelAndRealPIDKill", testUnifiedSpawnSetModelAndPIDKill)
}

func testUnifiedCanonicalPathAndLegacy(t *testing.T) {
	proj := t.TempDir()
	sb, err := OpenSwarmBoard(proj)
	if err != nil {
		t.Fatalf("OpenSwarmBoard(proj): %v", err)
	}
	canonical := filepath.Join(proj, ".legoloop-swarm", "swarm.json")
	if _, err := os.Stat(canonical); err != nil {
		t.Fatalf("arquivo canonico ausente em %s: %v", canonical, err)
	}
	if _, err := sb.PostMessage("coordenador", "@all", "canonico"); err != nil {
		t.Fatalf("PostMessage: %v", err)
	}

	// Abrir apontando diretamente para .legoloop-swarm nao deve aninhar .legoloop-swarm/.legoloop-swarm.
	sbDirect, err := OpenSwarmBoard(filepath.Join(proj, ".legoloop-swarm"))
	if err != nil {
		t.Fatalf("OpenSwarmBoard(.legoloop-swarm): %v", err)
	}
	nested := filepath.Join(proj, ".legoloop-swarm", ".legoloop-swarm", "swarm.json")
	if _, err := os.Stat(nested); err == nil {
		t.Fatalf("nao deveria criar caminho aninhado %s", nested)
	}
	msgs, err := sbDirect.ReadMessages(0, nil)
	if err != nil || len(msgs) != 1 || msgs[0].Body != "canonico" {
		t.Fatalf("leitura direta falhou: msgs=%+v err=%v", msgs, err)
	}

	// Retrocompatibilidade: .swarm/swarm.json e swarm_board.json.
	root := t.TempDir()
	legacySwarmProj := filepath.Join(root, "proc-swarm")
	if err := os.MkdirAll(filepath.Join(legacySwarmProj, ".swarm"), 0o755); err != nil {
		t.Fatal(err)
	}
	legacyPayload := `{"workers":[],"messages":[{"seq":1,"from":"coord","to":"@all","body":"legado-swarm","created_at":"2026-11-14T00:00:00Z"}],"next_seq":2}`
	if err := os.WriteFile(filepath.Join(legacySwarmProj, ".swarm", "swarm.json"), []byte(legacyPayload), 0o644); err != nil {
		t.Fatal(err)
	}

	legacyBoardProj := filepath.Join(root, "proc-board")
	if err := os.MkdirAll(legacyBoardProj, 0o755); err != nil {
		t.Fatal(err)
	}
	legacyBoardPayload := `{"workers":[],"messages":[{"seq":1,"from":"coord","to":"@all","body":"legado-board","created_at":"2026-11-14T00:00:00Z"}],"next_seq":2}`
	if err := os.WriteFile(filepath.Join(legacyBoardProj, "swarm_board.json"), []byte(legacyBoardPayload), 0o644); err != nil {
		t.Fatal(err)
	}

	boards := DiscoverBoards(root, "", "")
	if len(boards) != 2 {
		t.Fatalf("esperava 2 quadros descobertos, obteve %+v", boards)
	}
	sbLeg1, err := OpenSwarmBoard(boards["proc-swarm"])
	if err != nil {
		t.Fatalf("OpenSwarmBoard proc-swarm: %v", err)
	}
	if _, err := sbLeg1.PostMessage("coord", "@all", "novo-apos-legado"); err != nil {
		t.Fatalf("PostMessage legado1: %v", err)
	}
	migratedPath := filepath.Join(legacySwarmProj, ".legoloop-swarm", "swarm.json")
	if _, err := os.Stat(migratedPath); err != nil {
		t.Fatalf("escrita nao foi unificada em %s: %v", migratedPath, err)
	}
	m1, err := sbLeg1.ReadMessages(0, nil)
	if err != nil || len(m1) != 2 || m1[0].Body != "legado-swarm" {
		t.Fatalf("migracao perdeu mensagens antigas: %+v err=%v", m1, err)
	}

	sbLeg2, err := OpenSwarmBoard(boards["proc-board"])
	if err != nil {
		t.Fatalf("OpenSwarmBoard proc-board: %v", err)
	}
	m2, err := sbLeg2.ReadMessages(0, nil)
	if err != nil || len(m2) != 1 || m2[0].Body != "legado-board" {
		t.Fatalf("leitura retrocompativel swarm_board.json falhou: %+v err=%v", m2, err)
	}
}

func testUnifiedFileLockCrossProcess(t *testing.T) {
	dir := t.TempDir()
	if _, err := OpenSwarmBoard(dir); err != nil {
		t.Fatalf("OpenSwarmBoard: %v", err)
	}
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	const procs = 4
	const perProc = 8
	var wg sync.WaitGroup
	errs := make(chan error, procs)
	for p := 0; p < procs; p++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			cmd := exec.Command(exe, "-test.run=^TestHelperSwarmUnifiedProcess$")
			cmd.Env = append(os.Environ(),
				"GO_SWARM_UNIFIED_HELPER=post",
				"GO_SWARM_HELPER_DIR="+dir,
				fmt.Sprintf("GO_SWARM_HELPER_FROM=worker-%d", idx),
				fmt.Sprintf("GO_SWARM_HELPER_COUNT=%d", perProc),
			)
			if out, err := cmd.CombinedOutput(); err != nil {
				errs <- fmt.Errorf("proc %d: %w (%s)", idx, err, strings.TrimSpace(string(out)))
			}
		}(p)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatal(err)
	}
	sb, err := OpenSwarmBoard(dir)
	if err != nil {
		t.Fatal(err)
	}
	msgs, err := sb.ReadMessages(0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != procs*perProc {
		t.Fatalf("esperava %d mensagens sem perda cross-processo, obteve %d", procs*perProc, len(msgs))
	}
	for i, m := range msgs {
		if m.Seq != int64(i+1) {
			t.Fatalf("sequencia nao linearizavel em [%d]: obteve %d", i, m.Seq)
		}
	}
	entries, err := os.ReadDir(filepath.Join(dir, ".legoloop-swarm"))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.Contains(e.Name(), ".tmp") {
			t.Fatalf("arquivo temporario residual encontrado: %s", e.Name())
		}
	}
	t1 := atomicTempPath(filepath.Join(dir, ".legoloop-swarm", "swarm.json"))
	t2 := atomicTempPath(filepath.Join(dir, ".legoloop-swarm", "swarm.json"))
	if t1 == t2 || !strings.Contains(t1, fmt.Sprintf(".tmp.%d.", os.Getpid())) {
		t.Fatalf("atomicTempPath deve ser unico com PID+ns: t1=%q t2=%q", t1, t2)
	}
}

func testUnifiedNameAndWorkspaceModes(t *testing.T) {
	repo := initTestGitRepo(t)
	sb, err := OpenSwarmBoard(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	invalidNames := []string{"../evil", "..\\evil", "a/b", "a\\b", ".", "..", "-rf", "bad name", "<img/onerror>", "all"}
	for _, bad := range invalidNames {
		if _, err := sb.AddWorker(AddWorkerRequest{Name: bad, Role: "engenheiro-tdd", RepoDir: repo}); err == nil {
			t.Fatalf("esperava rejeicao para nome invalido %q", bad)
		}
	}
	if _, err := sb.AddWorker(AddWorkerRequest{Name: "w-invalid-ws", Role: "engenheiro-tdd", WorkspaceMode: WorkspaceMode("invalid")}); err == nil {
		t.Fatal("esperava rejeicao para workspace_mode invalido")
	}

	wInherit, err := sb.AddWorker(AddWorkerRequest{Name: "w-inherit", Role: "engenheiro-tdd", WorkspaceMode: WorkspaceInherit, RepoDir: repo})
	if err != nil {
		t.Fatalf("AddWorker inherit: %v", err)
	}
	if wInherit.WorktreePath != "" || wInherit.BranchName != "" {
		t.Fatalf("inherit nao deve criar worktree nem branch: %+v", wInherit)
	}

	wShare, err := sb.AddWorker(AddWorkerRequest{Name: "w-share", Role: "engenheiro-tdd", WorkspaceMode: WorkspaceShare, RepoDir: repo})
	if err != nil {
		t.Fatalf("AddWorker share: %v", err)
	}
	if wShare.WorktreePath == "" || wShare.BranchName != "" {
		t.Fatalf("share deve criar worktree detached sem branch dedicada: %+v", wShare)
	}
	headOut, err := exec.Command("git", "-C", wShare.WorktreePath, "rev-parse", "--abbrev-ref", "HEAD").CombinedOutput()
	if err != nil || strings.TrimSpace(string(headOut)) != "HEAD" {
		t.Fatalf("share deve estar em HEAD detached, obteve %q err=%v", strings.TrimSpace(string(headOut)), err)
	}

	wBranch, err := sb.AddWorker(AddWorkerRequest{Name: "w-branch", Role: "engenheiro-tdd", WorkspaceMode: WorkspaceBranch, RepoDir: repo})
	if err != nil {
		t.Fatalf("AddWorker branch: %v", err)
	}
	if wBranch.WorktreePath == "" || wBranch.BranchName == "" {
		t.Fatalf("branch deve criar worktree e branch dedicada: %+v", wBranch)
	}
	brOut, err := exec.Command("git", "-C", wBranch.WorktreePath, "rev-parse", "--abbrev-ref", "HEAD").CombinedOutput()
	if err != nil || strings.TrimSpace(string(brOut)) != wBranch.BranchName {
		t.Fatalf("branch worktree deveria apontar para %q, obteve %q err=%v", wBranch.BranchName, strings.TrimSpace(string(brOut)), err)
	}
	t.Cleanup(func() {
		_ = sb.KillWorker(wShare.Name, repo)
		_ = sb.KillWorker(wBranch.Name, repo)
	})
}

func testUnifiedExactMentionBoundary(t *testing.T) {
	sb, err := OpenSwarmBoard(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := sb.AddWorker(AddWorkerRequest{Name: "w1", Role: "engenheiro-tdd"}); err != nil {
		t.Fatal(err)
	}
	if _, err := sb.AddWorker(AddWorkerRequest{Name: "w10", Role: "engenheiro-tdd"}); err != nil {
		t.Fatal(err)
	}

	mFalse, err := sb.PostMessage("coordenador", "@w10", "oi @w10, ignore @w1-extra e email@w1 e @alligator")
	if err != nil {
		t.Fatal(err)
	}
	st, err := sb.Status()
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range st.Workers {
		if w.Name == "w1" && w.Status == WorkerWorking {
			t.Fatalf("w1 nao deveria ser ativado por @w10/@w1-extra/email@w1")
		}
		if w.Name == "w10" && w.Status != WorkerWorking {
			t.Fatalf("w10 deveria estar working")
		}
	}

	ctxFail, cancelFail := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancelFail()
	if got, err := sb.WaitMessages(ctxFail, "@w1", 0, 20*time.Millisecond); err == nil || len(got) > 0 {
		t.Fatalf("WaitMessages(@w1) nao devia casar com @w10/@w1-extra/email@w1: got=%+v", got)
	}
	ctxAllFail, cancelAllFail := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancelAllFail()
	if got, err := sb.WaitMessages(ctxAllFail, "@all", 0, 20*time.Millisecond); err == nil || len(got) > 0 {
		t.Fatalf("WaitMessages(@all) nao devia casar com @alligator: got=%+v", got)
	}

	if _, err := sb.PostMessage("coordenador", "", "agora sim @w1 @w10"); err != nil {
		t.Fatal(err)
	}
	ctxOK, cancelOK := context.WithTimeout(context.Background(), time.Second)
	defer cancelOK()
	got, err := sb.WaitMessages(ctxOK, "@w1", mFalse.Seq, 20*time.Millisecond)
	if err != nil || len(got) != 1 {
		t.Fatalf("WaitMessages(@w1) exato falhou: got=%+v err=%v", got, err)
	}
}

func testUnifiedSpawnSetModelAndPIDKill(t *testing.T) {
	repo := initTestGitRepo(t)
	boardRoot := filepath.Join(t.TempDir(), "swarm-proj")
	sb, err := OpenSwarmBoard(boardRoot)
	if err != nil {
		t.Fatal(err)
	}
	w, err := sb.AddWorker(AddWorkerRequest{Name: "w-real", Role: "engenheiro-backend", WorkspaceMode: WorkspaceBranch, Model: "m1", RepoDir: repo})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := sb.SetWorkerModel(w.Name, "m2"); err != nil {
		t.Fatalf("SetWorkerModel: %v", err)
	}

	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	modelFile := filepath.Join(t.TempDir(), "model.txt")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	t.Setenv("GO_SWARM_UNIFIED_HELPER", "sleep")
	t.Setenv("GO_SWARM_MODEL_OUT", modelFile)

	pid, err := sb.SpawnWorker(ctx, w.Name, []string{exe, "-test.run=^TestHelperSwarmUnifiedProcess$"})
	if err != nil {
		t.Fatalf("SpawnWorker: %v", err)
	}
	if pid <= 0 {
		t.Fatalf("PID invalido: %d", pid)
	}
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		if b, err := os.ReadFile(modelFile); err == nil && string(b) == "m2" {
			break
		}
		time.Sleep(25 * time.Millisecond)
	}
	if b, err := os.ReadFile(modelFile); err != nil || string(b) != "m2" {
		t.Fatalf("subprocesso nao recebeu modelo atualizado m2: %q err=%v", string(b), err)
	}

	// Kill cross-processo via nova abertura do mesmo board.
	sb2, err := OpenSwarmBoard(boardRoot)
	if err != nil {
		t.Fatal(err)
	}
	if err := sb2.KillWorker(w.Name, repo); err != nil {
		t.Fatalf("KillWorker: %v", err)
	}
	// O kill do Windows é assíncrono: espera ativa pela morte real do PID (janela de SO, sem data fixa).
	morto := time.Now().Add(10 * time.Second)
	for isOSProcessAlive(pid) && time.Now().Before(morto) {
		time.Sleep(25 * time.Millisecond)
	}
	if isOSProcessAlive(pid) {
		t.Fatalf("PID %d continuou vivo apos KillWorker", pid)
	}
	if _, err := os.Stat(w.WorktreePath); !os.IsNotExist(err) {
		t.Fatalf("worktree %s deveria ter sido removido apos KillWorker, err=%v", w.WorktreePath, err)
	}
	var finalKilled Worker
	prazoKilled := time.Now().Add(30 * time.Second)
	for time.Now().Before(prazoKilled) {
		st, err := sb2.Status()
		if err != nil {
			t.Fatal(err)
		}
		if len(st.Workers) == 1 && st.Workers[0].Status == WorkerKilled && st.Workers[0].PID == 0 {
			finalKilled = st.Workers[0]
			break
		}
		time.Sleep(25 * time.Millisecond)
	}
	if finalKilled.Status != WorkerKilled {
		t.Fatalf("esperava killed com PID=0, obteve %+v", finalKilled)
	}
	if _, err := sb2.SetWorkerModel(w.Name, "m3"); err == nil {
		t.Fatal("nao devia permitir SetWorkerModel em worker killed")
	}
	if err := os.RemoveAll(boardRoot); err != nil {
		t.Fatalf("vazamento impediu remocao de boardRoot no Windows: %v", err)
	}

	// Subprocesso com saida natural deve transicionar para stopped com PID=0.
	sbNat, err := OpenSwarmBoard(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := sbNat.AddWorker(AddWorkerRequest{Name: "w-exit", Role: "engenheiro-tdd"}); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GO_SWARM_UNIFIED_HELPER", "exit0")
	if _, err := sbNat.SpawnWorker(ctx, "w-exit", []string{exe, "-test.run=^TestHelperSwarmUnifiedProcess$"}); err != nil {
		t.Fatalf("SpawnWorker exit0: %v", err)
	}
	waitStopped := time.Now().Add(30 * time.Second)
	var finalWorker Worker
	for time.Now().Before(waitStopped) {
		stNat, err := sbNat.Status()
		if err == nil && len(stNat.Workers) == 1 && stNat.Workers[0].Status == WorkerStopped {
			finalWorker = stNat.Workers[0]
			break
		}
		time.Sleep(25 * time.Millisecond)
	}
	if finalWorker.Status != WorkerStopped || finalWorker.PID != 0 {
		t.Fatalf("esperava stopped e PID=0 apos saida natural, obteve %+v", finalWorker)
	}
}

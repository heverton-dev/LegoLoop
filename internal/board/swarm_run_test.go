package board

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"legoloop/internal/dag"
)

func TestHelperProcessRunLevels(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS_BOARD") != "1" {
		return
	}
	for i, arg := range os.Args {
		if arg != "--" {
			continue
		}
		rest := os.Args[i+1:]
		if len(rest) == 0 {
			os.Exit(0)
		}
		switch rest[0] {
		case "touch":
			if len(rest) > 1 {
				_ = os.WriteFile(rest[1], []byte("ok"), 0o600)
			}
			os.Exit(0)
		case "exit":
			code := 1
			if len(rest) > 1 {
				if n, err := strconv.Atoi(rest[1]); err == nil {
					code = n
				}
			}
			os.Exit(code)
		case "sleep":
			time.Sleep(5 * time.Second)
			os.Exit(0)
		}
		os.Exit(0)
	}
}

func helperCmd(args ...string) []string {
	cmd := []string{os.Args[0], "-test.run=^TestHelperProcessRunLevels$", "--"}
	return append(cmd, args...)
}

func TestWorkerNameFromNode(t *testing.T) {
	usados := map[string]bool{}
	if got := WorkerNameFromNode("no.1/a", usados); got != "no-1-a" {
		t.Fatalf("esperava no-1-a, obtive %q", got)
	}
	if got := WorkerNameFromNode("no.1/a", usados); got != "no-1-a-2" {
		t.Fatalf("esperava colisão no-1-a-2, obtive %q", got)
	}
	if got := WorkerNameFromNode("-invalido", usados); got != "n-invalido" {
		t.Fatalf("esperava prefixo n, obtive %q", got)
	}
	longo := strings.Repeat("a", 80)
	if got := WorkerNameFromNode(longo, usados); len(got) != 60 {
		t.Fatalf("esperava truncamento para 60, obtive %d", len(got))
	}
}

func TestLevelsFromDAG(t *testing.T) {
	jsonDAG := `{
  "name": "demo",
  "nodes": [
    {"id": "n1", "type": "script", "role": "engenheiro-tdd", "command": ["echo", "1"], "timeout_sec": 5, "routing": {"mode": "suggested", "model": "fast"}, "output": {"artifact": "a1.txt", "gate": "file_exists('a1.txt')"}},
    {"id": "n2", "type": "script", "papel": "engenheiro-tdd", "depends_on": ["n1"], "command": ["echo", "2"], "routing": {"mode": "suggested", "model": "fast"}, "output": {"artifact": "a2.txt", "gate": "file_exists('a2.txt')"}},
    {"id": "n3", "type": "skill", "ref": "s1", "role": "engenheiro-tdd", "depends_on": ["n2"], "routing": {"mode": "suggested", "model": "fast"}, "output": {"artifact": "a3.txt", "gate": "file_exists('a3.txt')"}},
    {"id": "sg", "type": "stop_gate", "role": "engenheiro-tdd", "routing": {"mode": "suggested", "model": "fast"}, "output": {"artifact": "sg.txt", "gate": "file_exists('sg.txt')"}},
    {"id": "sem_papel", "type": "script", "command": ["echo", "x"], "routing": {"mode": "suggested", "model": "fast"}, "output": {"artifact": "sp.txt", "gate": "file_exists('sp.txt')"}}
  ]
}`
	d, err := dag.Parse([]byte(jsonDAG))
	if err != nil {
		t.Fatalf("dag.Parse falhou: %v", err)
	}
	levels, skipped, err := LevelsFromDAG(d, []string{"echo", "global"})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(levels) != 3 {
		t.Fatalf("esperava 3 níveis, obtive %d", len(levels))
	}
	if len(skipped) != 2 {
		t.Fatalf("esperava 2 ignorados, obtive %v", skipped)
	}
	if _, _, err := LevelsFromDAG(d, nil); err == nil {
		t.Fatal("esperava erro quando nó não tem comando nem global")
	}
}

func TestRunLevelsSucessoEGates(t *testing.T) {
	t.Setenv("GO_WANT_HELPER_PROCESS_BOARD", "1")
	proj := t.TempDir()
	sb, err := OpenSwarmBoard(proj)
	if err != nil {
		t.Fatal(err)
	}
	okFile := filepath.Join(proj, "ok.txt")
	levels := [][]RunNode{
		{{Name: "w1", Role: "engenheiro-tdd", Command: helperCmd("touch", okFile), Gate: "file_exists('ok.txt')", Timeout: 15 * time.Second}},
		{{Name: "w2", Role: "engenheiro-tdd", Command: helperCmd("exit", "0")}},
	}
	res, err := sb.RunLevels(context.Background(), levels, RunOptions{
		Mode:           WorkspaceInherit,
		RepoDir:        proj,
		DefaultTimeout: 10 * time.Second,
		Poll:           30 * time.Millisecond,
		CheckGates:     true,
	})
	if err != nil {
		t.Fatalf("RunLevels falhou: %v", err)
	}
	if len(res) != 2 || res[0].Workers[0].GateOK == nil || !*res[0].Workers[0].GateOK {
		t.Fatalf("resultado inesperado: %+v", res)
	}
}

func TestRunLevelsFalhas(t *testing.T) {
	t.Setenv("GO_WANT_HELPER_PROCESS_BOARD", "1")

	t.Run("exit_code_nonzero", func(t *testing.T) {
		proj := t.TempDir()
		sb, _ := OpenSwarmBoard(proj)
		_, err := sb.RunLevels(context.Background(), [][]RunNode{
			{{Name: "wfail", Role: "engenheiro-tdd", Command: helperCmd("exit", "3")}},
		}, RunOptions{Mode: WorkspaceInherit, RepoDir: proj, Poll: 20 * time.Millisecond})
		if err == nil || !strings.Contains(err.Error(), "saiu com codigo 3") {
			t.Fatalf("esperava erro de codigo 3, obtive %v", err)
		}
	})

	t.Run("gate_reprovado", func(t *testing.T) {
		proj := t.TempDir()
		sb, _ := OpenSwarmBoard(proj)
		_, err := sb.RunLevels(context.Background(), [][]RunNode{
			{{Name: "wgate", Role: "engenheiro-tdd", Command: helperCmd("exit", "0"), Gate: "file_exists('nao_existe.txt')"}},
		}, RunOptions{Mode: WorkspaceInherit, RepoDir: proj, Poll: 20 * time.Millisecond, CheckGates: true})
		if err == nil || !strings.Contains(err.Error(), "gate reprovou") {
			t.Fatalf("esperava erro de gate reprovado, obtive %v", err)
		}
	})

	t.Run("gate_invalido", func(t *testing.T) {
		proj := t.TempDir()
		sb, _ := OpenSwarmBoard(proj)
		_, err := sb.RunLevels(context.Background(), [][]RunNode{
			{{Name: "wbad", Role: "engenheiro-tdd", Command: helperCmd("exit", "0"), Gate: "funcao_inexistente("}},
		}, RunOptions{Mode: WorkspaceInherit, RepoDir: proj, Poll: 20 * time.Millisecond, CheckGates: true})
		if err == nil || !strings.Contains(err.Error(), "invalido") {
			t.Fatalf("esperava erro de gate invalido, obtive %v", err)
		}
	})

	t.Run("timeout", func(t *testing.T) {
		proj := t.TempDir()
		sb, _ := OpenSwarmBoard(proj)
		_, err := sb.RunLevels(context.Background(), [][]RunNode{
			{{Name: "wslow", Role: "engenheiro-tdd", Command: helperCmd("sleep")}},
		}, RunOptions{Mode: WorkspaceInherit, RepoDir: proj, DefaultTimeout: 80 * time.Millisecond, Poll: 20 * time.Millisecond})
		if err == nil || !strings.Contains(err.Error(), "espera dos workers") {
			t.Fatalf("esperava erro de timeout, obtive %v", err)
		}
	})
}

func TestGateRoot(t *testing.T) {
	if got := gateRoot(Worker{WorktreePath: "/wt", RepoDir: "/repo"}, "/arg", "/proj"); got != "/wt" {
		t.Fatalf("esperava /wt, obtive %q", got)
	}
	if got := gateRoot(Worker{RepoDir: "/repo"}, "/arg", "/proj"); got != "/repo" {
		t.Fatalf("esperava /repo, obtive %q", got)
	}
	if got := gateRoot(Worker{}, "/arg", "/proj"); got != "/arg" {
		t.Fatalf("esperava /arg, obtive %q", got)
	}
	if got := gateRoot(Worker{}, "", "/proj"); got != "/proj" {
		t.Fatalf("esperava /proj, obtive %q", got)
	}
}

func TestBoardsPayload(t *testing.T) {
	root := t.TempDir()
	proj := filepath.Join(root, "proc-a")
	sb, err := OpenSwarmBoard(proj)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := sb.PostMessage("coord", "@all", "ola"); err != nil {
		t.Fatal(err)
	}
	raw, err := BoardsPayload(root, "", "")
	if err != nil {
		t.Fatalf("BoardsPayload falhou: %v", err)
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("JSON inválido: %v", err)
	}
	if _, ok := out["proc-a"]; !ok {
		t.Fatalf("esperava proc-a no payload, obtive %+v", out)
	}
	rawSingle, err := BoardsPayload("", filepath.Join(proj, ".legoloop-swarm"), "unico")
	if err != nil || !strings.Contains(string(rawSingle), "unico") {
		t.Fatalf("BoardsPayload singleDir falhou: %v (%s)", err, string(rawSingle))
	}
}

func TestClientMoveTicketEGuardReset(t *testing.T) {
	g := NewTokenBucketGuard(0)
	if err := g.Check(5, time.Now().Add(time.Minute).Unix()); err == nil {
		t.Fatal("esperava halted com remaining 5")
	}
	if !g.IsHalted() {
		t.Fatal("guard deveria estar halted")
	}

	statusCode := http.StatusOK
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer tok" {
			t.Errorf("header Authorization ausente: %q", r.Header.Get("Authorization"))
		}
		w.WriteHeader(statusCode)
	}))
	defer srv.Close()

	cli := NewClient(ClientConfig{BaseURL: srv.URL, Token: "tok", Guard: g})
	if err := cli.MoveTicket(context.Background(), "T1", StatusBuilding); err != ErrRateLimitExhausted {
		t.Fatalf("esperava ErrRateLimitExhausted, obtive %v", err)
	}

	g.Reset()
	if g.IsHalted() {
		t.Fatal("guard deveria ter sido resetado")
	}
	if err := cli.MoveTicket(context.Background(), "T1", StatusBuilding); err != nil {
		t.Fatalf("MoveTicket 200 falhou: %v", err)
	}
	statusCode = http.StatusNoContent
	if err := cli.MoveTicket(context.Background(), "T1", StatusDone); err != nil {
		t.Fatalf("MoveTicket 204 falhou: %v", err)
	}
	statusCode = http.StatusInternalServerError
	if err := cli.MoveTicket(context.Background(), "T1", StatusDone); err == nil {
		t.Fatal("esperava erro com status 500")
	}
}

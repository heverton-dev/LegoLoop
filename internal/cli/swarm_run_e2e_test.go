package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"legoloop/internal/board"
)

// TestHelperSwarmRunWorker e o corpo dos workers do E2E: um processo real, sem harness externo.
// Nivel 0 (a, b): so passa se o OUTRO worker tambem ja iniciou (barreira), o que prova paralelismo real.
// Nivel 1 (c): so passa se a e b ja terminaram, o que prova a ordem por nivel.
func TestHelperSwarmRunWorker(t *testing.T) {
	out := os.Getenv("GO_SWARM_RUN_OUT")
	if os.Getenv("GO_SWARM_RUN_HELPER") != "1" || out == "" {
		return
	}
	name := os.Getenv("LEGOLOOP_SWARM_WORKER")
	stamp := func(suffix string) {
		_ = os.WriteFile(filepath.Join(out, name+"."+suffix), []byte(strconv.FormatInt(time.Now().UnixNano(), 10)), 0o644)
	}
	exists := func(f string) bool { _, err := os.Stat(filepath.Join(out, f)); return err == nil }
	cwd, _ := os.Getwd()
	_ = os.WriteFile(filepath.Join(out, name+".cwd"), []byte(cwd), 0o644)
	stamp("start")
	if name == "c" {
		if !exists("a.end") || !exists("b.end") {
			os.Exit(4)
		}
	} else {
		deadline := time.Now().Add(15 * time.Second)
		for !(exists("a.start") && exists("b.start")) {
			if time.Now().After(deadline) {
				os.Exit(3)
			}
			time.Sleep(20 * time.Millisecond)
		}
	}
	if name == os.Getenv("GO_SWARM_RUN_FAIL") {
		os.Exit(7)
	}
	if name != os.Getenv("GO_SWARM_RUN_NOFILE") {
		_ = os.MkdirAll("out", 0o755)
		_ = os.WriteFile(filepath.Join("out", name+".md"), []byte("ok\n"), 0o644)
	}
	stamp("end")
	os.Exit(0)
}

func writeRunDAG(t *testing.T, dir string) string {
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
		"$schema": "https://legoloop.dev/schemas/dag.v1.json", "id": "demo-run", "name": "demo-run", "version": "1.0.0",
		"trigger": map[string]any{"type": "manual"},
		"nodes": []any{
			node("a", "engenheiro-backend"), node("b", "engenheiro-tdd"), node("c", "auditor-qualidade", "a", "b"),
		},
	}
	b, _ := json.MarshalIndent(d, "", "  ")
	p := filepath.Join(dir, "pipeline.dag.json")
	if err := os.WriteFile(p, b, 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func runSwarmRun(t *testing.T, repo, dagPath, fail string, extra ...string) (int, string, string) {
	t.Helper()
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	t.Setenv("GO_SWARM_RUN_HELPER", "1")
	t.Setenv("GO_SWARM_RUN_OUT", out)
	t.Setenv("GO_SWARM_RUN_FAIL", fail)
	var so, se bytes.Buffer
	args := []string{"swarm", "run", dagPath, "-dir", repo, "-repo", repo, "-timeout", "40s"}
	args = append(append(args, extra...), "--", exe, "-test.run=^TestHelperSwarmRunWorker$")
	code := Run(context.Background(), args, Deps{Stdout: &so, Stderr: &se})
	return code, so.String() + "\x00" + se.String() + "\x00" + out, out
}

func stampOf(t *testing.T, out, file string) int64 {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(out, file))
	if err != nil {
		t.Fatalf("faltou %s: %v", file, err)
	}
	n, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	return n
}

func assertBoardClean(t *testing.T, repo string, names ...string) {
	t.Helper()
	sb, err := board.OpenSwarmBoard(repo)
	if err != nil {
		t.Fatal(err)
	}
	st, err := sb.Status()
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range st.Workers {
		if w.Status != board.WorkerKilled || w.PID != 0 {
			t.Errorf("worker %s deveria estar killed com PID 0: %+v", w.Name, w)
		}
		if _, err := os.Stat(w.WorktreePath); !os.IsNotExist(err) {
			t.Errorf("worktree de %s sobrou: %v", w.Name, err)
		}
	}
	c := exec.Command("git", "worktree", "list", "--porcelain")
	c.Dir = repo
	if o, err := c.Output(); err != nil || bytes.Count(o, []byte("worktree ")) != 1 {
		t.Errorf("git worktree list deveria ter so o principal (err=%v): %s", err, o)
	}
}

// TestSwarmRunPorNivelParaleloEKillLimpo: a e b rodam juntos (barreira), c so depois, e tudo termina limpo.
func TestSwarmRunPorNivelParaleloEKillLimpo(t *testing.T) {
	if testing.Short() {
		t.Skip("E2E real sobe processos")
	}
	repo := e2eGitRepo(t)
	dagPath := writeRunDAG(t, t.TempDir())
	code, joined, out := runSwarmRun(t, repo, dagPath, "")
	parts := strings.Split(joined, "\x00")
	if code != 0 {
		t.Fatalf("swarm run saiu %d\nstdout=%s\nstderr=%s", code, parts[0], parts[1])
	}
	var res struct {
		Levels []board.RunLevelResult `json:"levels"`
	}
	if err := json.Unmarshal([]byte(parts[0]), &res); err != nil {
		t.Fatalf("json: %v\n%s", err, parts[0])
	}
	if len(res.Levels) != 2 || len(res.Levels[0].Workers) != 2 || len(res.Levels[1].Workers) != 1 {
		t.Fatalf("niveis inesperados: %+v", res.Levels)
	}
	for _, l := range res.Levels {
		for _, w := range l.Workers {
			if w.ExitCode != 0 {
				t.Errorf("worker %s saiu %d", w.Name, w.ExitCode)
			}
		}
	}
	if res.Levels[1].StartedAt.Before(res.Levels[0].EndedAt) {
		t.Errorf("nivel 1 comecou antes do fim do nivel 0: %+v", res.Levels)
	}
	if stampOf(t, out, "c.start") < max(stampOf(t, out, "a.end"), stampOf(t, out, "b.end")) {
		t.Error("c iniciou antes de a e b terminarem")
	}
	cwdA, _ := os.ReadFile(filepath.Join(out, "a.cwd"))
	cwdB, _ := os.ReadFile(filepath.Join(out, "b.cwd"))
	if string(cwdA) == string(cwdB) || !strings.Contains(string(cwdA), "worktrees") || !strings.Contains(string(cwdB), "worktrees") {
		t.Errorf("a e b deveriam rodar em worktrees distintos: %q %q", cwdA, cwdB)
	}
	assertBoardClean(t, repo)
}

// TestSwarmRunFalhaEncerraELimpa: worker com codigo 7 derruba o nivel; o nivel seguinte nao comeca e nada sobra.
func TestSwarmRunFalhaEncerraELimpa(t *testing.T) {
	if testing.Short() {
		t.Skip("E2E real sobe processos")
	}
	repo := e2eGitRepo(t)
	dagPath := writeRunDAG(t, t.TempDir())
	code, joined, out := runSwarmRun(t, repo, dagPath, "b")
	parts := strings.Split(joined, "\x00")
	if code != 1 || !strings.Contains(parts[1], "saiu com codigo 7") {
		t.Fatalf("esperava codigo 1 e erro do worker b: code=%d stderr=%s", code, parts[1])
	}
	if _, err := os.Stat(filepath.Join(out, "c.start")); err == nil {
		t.Error("o nivel 1 nao devia ter comecado apos a falha do nivel 0")
	}
	assertBoardClean(t, repo)
}

// TestSwarmRunGatesAprovaQuandoTodosProduzemArtefato: com -gates, o gate de cada no e avaliado no worktree do worker.
func TestSwarmRunGatesAprovaQuandoTodosProduzemArtefato(t *testing.T) {
	if testing.Short() {
		t.Skip("E2E real sobe processos")
	}
	repo := e2eGitRepo(t)
	code, joined, _ := runSwarmRun(t, repo, writeRunDAG(t, t.TempDir()), "", "-gates")
	parts := strings.Split(joined, "\x00")
	if code != 0 {
		t.Fatalf("esperava 0: stdout=%s stderr=%s", parts[0], parts[1])
	}
	var res struct {
		Levels []board.RunLevelResult `json:"levels"`
	}
	if err := json.Unmarshal([]byte(parts[0]), &res); err != nil {
		t.Fatal(err)
	}
	for _, l := range res.Levels {
		for _, w := range l.Workers {
			if w.GateOK == nil || !*w.GateOK {
				t.Errorf("worker %s deveria ter gate_ok=true: %+v", w.Name, w)
			}
		}
	}
	assertBoardClean(t, repo)
}

// TestSwarmRunGatesReprovaNivel: worker com codigo 0 mas sem o artefato do gate derruba o nivel e limpa tudo.
func TestSwarmRunGatesReprovaNivel(t *testing.T) {
	if testing.Short() {
		t.Skip("E2E real sobe processos")
	}
	repo := e2eGitRepo(t)
	t.Setenv("GO_SWARM_RUN_NOFILE", "b")
	code, joined, out := runSwarmRun(t, repo, writeRunDAG(t, t.TempDir()), "", "-gates")
	parts := strings.Split(joined, "\x00")
	if code != 1 || !strings.Contains(parts[1], "gate reprovou") || !strings.Contains(parts[1], "b: file_exists('out/b.md')") {
		t.Fatalf("esperava codigo 1 e gate reprovado no no b: code=%d stderr=%s", code, parts[1])
	}
	if strings.Contains(parts[1], "a: file_exists") {
		t.Errorf("o no a produziu o artefato e nao deveria reprovar: %s", parts[1])
	}
	if _, err := os.Stat(filepath.Join(out, "c.start")); err == nil {
		t.Error("o nivel 1 nao devia ter comecado apos o gate reprovar")
	}
	assertBoardClean(t, repo)
}

// Sem -gates o comportamento antigo vale: gate do no nao e avaliado.
func TestSwarmRunSemGatesIgnoraGateDoNo(t *testing.T) {
	if testing.Short() {
		t.Skip("E2E real sobe processos")
	}
	repo := e2eGitRepo(t)
	t.Setenv("GO_SWARM_RUN_NOFILE", "b")
	if code, joined, _ := runSwarmRun(t, repo, writeRunDAG(t, t.TempDir()), ""); code != 0 {
		t.Fatalf("sem -gates o run deveria passar: %s", joined)
	}
}

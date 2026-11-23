package cli

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/board"
)

func compilarFakeWorker(t *testing.T) string {
	t.Helper()
	src := filepath.Join(t.TempDir(), "worker_proc.go")
	if err := os.WriteFile(src, []byte("package main\nfunc main() {}\n"), 0o600); err != nil {
		t.Fatalf("falha ao gravar fonte de fake worker: %v", err)
	}
	bin := filepath.Join(t.TempDir(), "fake_worker")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	if out, err := exec.Command("go", "build", "-o", bin, src).CombinedOutput(); err != nil {
		t.Fatalf("falha ao compilar fake worker: %v\n%s", err, out)
	}
	return bin
}

func gravarDAGNativeSwarm(t *testing.T, dir, helperBin string) string {
	t.Helper()
	dagData := map[string]any{
		"name":    "dag-native-swarm",
		"version": "1.0.0",
		"nodes": []any{
			map[string]any{
				"id":        "no-worker-1",
				"type":      "script",
				"role":      "engenheiro-tdd",
				"command":   []string{helperBin},
				"isolation": map[string]any{"executor": "native_swarm"},
				"routing":   map[string]any{"mode": "manual", "model": "m1"},
				"output":    map[string]any{"artifact": "saida.md", "gate": "exit_code == 0"},
			},
		},
	}
	raw, err := json.MarshalIndent(dagData, "", "  ")
	if err != nil {
		t.Fatalf("falha ao serializar DAG: %v", err)
	}
	p := filepath.Join(dir, "native_swarm.dag.json")
	if err := os.WriteFile(p, raw, 0o600); err != nil {
		t.Fatalf("falha ao gravar DAG: %v", err)
	}
	return p
}

func gravarDAGSemSwarm(t *testing.T, dir string) string {
	t.Helper()
	dagData := map[string]any{
		"name":    "dag-sem-swarm",
		"version": "1.0.0",
		"nodes": []any{
			map[string]any{
				"id":        "no-stop-1",
				"type":      "stop_gate",
				"isolation": map[string]any{"executor": "local_process"},
				"routing":   map[string]any{"mode": "manual", "model": "m1"},
				"output":    map[string]any{"artifact": "stop.txt", "gate": "exit_code == 0"},
			},
		},
	}
	raw, err := json.MarshalIndent(dagData, "", "  ")
	if err != nil {
		t.Fatalf("falha ao serializar DAG: %v", err)
	}
	p := filepath.Join(dir, "sem_swarm.dag.json")
	if err := os.WriteFile(p, raw, 0o600); err != nil {
		t.Fatalf("falha ao gravar DAG: %v", err)
	}
	return p
}

func verificarWorkerRegistrado(t *testing.T, repo string) {
	t.Helper()
	sb, err := board.OpenSwarmBoard(repo)
	if err != nil {
		t.Fatalf("falha ao abrir SwarmBoard: %v", err)
	}
	status, err := sb.Status()
	if err != nil {
		t.Fatalf("falha ao obter status do Swarm: %v", err)
	}
	if len(status.Workers) == 0 {
		t.Fatalf("esperava ao menos um worker registrado no SwarmBoard")
	}
	for _, w := range status.Workers {
		if strings.Contains(w.Name, "no-worker-1") || w.Role == "engenheiro-tdd" {
			return
		}
	}
	t.Fatalf("worker no-worker-1 não encontrado no SwarmBoard: %+v", status.Workers)
}

func TestRunNativeSwarmBinario(t *testing.T) {
	bin := buildBinarioTeste(t)
	repo := e2eGitRepo(t)
	workerBin := compilarFakeWorker(t)
	dagPath := gravarDAGNativeSwarm(t, repo, workerBin)
	relatorioPath := filepath.Join(repo, "relatorio.md")

	code, saida := executarBinario(bin, "run", dagPath, "--workdir", repo, "--allow-scripts", "--report", relatorioPath)
	if code != 0 {
		t.Fatalf("esperava exit code 0, obteve %d\nsaida:\n%s", code, saida)
	}
	relatorio, err := os.ReadFile(relatorioPath)
	if err != nil {
		t.Fatalf("falha ao ler arquivo de relatório: %v", err)
	}
	if !strings.Contains(string(relatorio), "no-worker-1") || !strings.Contains(strings.ToLower(string(relatorio)), "passed") {
		t.Fatalf("relatório deve conter nó passed:\n%s", string(relatorio))
	}
	verificarWorkerRegistrado(t, repo)
}

func TestRunSemNativeSwarmNaoCriaDiretorio(t *testing.T) {
	bin := buildBinarioTeste(t)
	dir := t.TempDir()
	dagPath := gravarDAGSemSwarm(t, dir)

	code, saida := executarBinario(bin, "run", dagPath, "--workdir", dir, "--yes")
	if code != 0 {
		t.Fatalf("esperava exit code 0, obteve %d\nsaida:\n%s", code, saida)
	}
	swarmDir := filepath.Join(dir, ".legoloop-swarm")
	if _, err := os.Stat(swarmDir); !os.IsNotExist(err) {
		t.Fatalf("diretório do Swarm (.legoloop-swarm) não deve ser criado quando não há native_swarm: %v", err)
	}
}

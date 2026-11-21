package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"legoloop/internal/caminhos"
	"legoloop/internal/runner"
)

const dagTresScriptsDeterminismo = `{
  "name": "determinismo-tres-scripts",
  "nodes": [
    {
      "id": "n1",
      "type": "script",
      "command": ["go", "version"],
      "timeout_sec": 60,
      "routing": {"mode": "manual", "model": "none"},
      "output": {"artifact": "n1.txt", "gate": "exit_code == 0 && file_exists('n1.txt')"}
    },
    {
      "id": "n2",
      "type": "script",
      "command": ["go", "env", "GOOS"],
      "timeout_sec": 60,
      "routing": {"mode": "manual", "model": "none"},
      "output": {"artifact": "n2.txt", "gate": "exit_code == 0 && file_exists('n2.txt')"}
    },
    {
      "id": "n3",
      "type": "script",
      "depends_on": ["n1", "n2"],
      "command": ["go", "env", "GOARCH"],
      "timeout_sec": 60,
      "routing": {"mode": "manual", "model": "none"},
      "output": {"artifact": "n3.txt", "gate": "exit_code == 0 && file_exists('n3.txt')"}
    }
  ]
}`

// normalizarRelatorioRun devolve uma cópia do relatório zerando apenas a duração de cada nó.
func normalizarRelatorioRun(rep runner.Report) runner.Report {
	out := rep
	out.Nodes = make([]runner.NodeReport, len(rep.Nodes))
	copy(out.Nodes, rep.Nodes)
	for i := range out.Nodes {
		out.Nodes[i].Duration = 0
	}
	return out
}

func executarRunDeterminismo(t *testing.T, bin, dagPath, root string) (runner.Report, string) {
	t.Helper()
	router := binRouter(t, 100)
	work := filepath.Join(root, "work")
	vault := filepath.Join(root, "vault")
	repPath := filepath.Join(root, "report.md")
	res := runBinaryEm(t, bin, router, root, work,
		"run", dagPath, "--workdir", work, "--skills", vault,
		"--allow-scripts", "--json", "--report", repPath)
	if res.code != 0 {
		t.Fatalf("legoloop run saiu %d\nstdout:\n%s\nstderr:\n%s", res.code, res.out, res.errb)
	}
	var rep runner.Report
	if err := json.Unmarshal([]byte(res.out), &rep); err != nil {
		t.Fatalf("json inválido do relatório: %v\nsaída:\n%s", err, res.out)
	}
	factual, err := os.ReadFile(repPath)
	if err != nil {
		t.Fatalf("ler relatório factual: %v", err)
	}
	return rep, string(factual)
}

// TestRunScriptsDuasVezesRelatorioIdentico (AC-E01.2): executar o binário real
// duas vezes sobre um DAG de três nós script gera relatórios idênticos,
// descontadas apenas as durações.
func TestRunScriptsDuasVezesRelatorioIdentico(t *testing.T) {
	bin := buildBinary(t)
	dagDir := t.TempDir()
	dagPath := caminhos.CaminhoDAG(dagDir)
	if err := os.WriteFile(dagPath, []byte(dagTresScriptsDeterminismo), 0o644); err != nil {
		t.Fatal(err)
	}
	rep1, factual1 := executarRunDeterminismo(t, bin, dagPath, t.TempDir())
	rep2, factual2 := executarRunDeterminismo(t, bin, dagPath, t.TempDir())
	if !rep1.OK || !rep2.OK || len(rep1.Nodes) != 3 || len(rep2.Nodes) != 3 {
		t.Fatalf("esperava 3 nós aprovados nas duas execuções: rep1=%+v rep2=%+v", rep1, rep2)
	}
	n1, n2 := normalizarRelatorioRun(rep1), normalizarRelatorioRun(rep2)
	if !reflect.DeepEqual(n1, n2) {
		t.Fatalf("relatórios divergiram (descontadas as durações):\nexec1: %+v\nexec2: %+v", n1, n2)
	}
	if factual1 != factual2 {
		t.Fatalf("relatórios factuais divergiram:\nexec1:\n%s\nexec2:\n%s", factual1, factual2)
	}
}

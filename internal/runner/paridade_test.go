package runner

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/dag"
)

func TestParidadeDoRelatorio(t *testing.T) {
	d := montarDAGDeParidade(t)
	workdir := t.TempDir()
	fExec := &fakeParidadeExec{}
	opt := Options{
		Workdir:     workdir,
		Executor:    fExec,
		Parallelism: 1,
		AutoApprove: true,
	}

	rep, err := Run(context.Background(), d, opt)
	if err != nil {
		t.Fatalf("Run falhou: %v", err)
	}
	normalizarRelatorio(&rep)
	verificarGoldenRelatorio(t, rep)
}

func montarDAGDeParidade(t *testing.T) *dag.DAG {
	t.Helper()
	src := `{
		"name": "paridade-executores-atuais",
		"nodes": [
			{
				"id": "gate-aprovacao",
				"type": "stop_gate",
				"routing": {"mode": "manual", "model": "human/operator"},
				"output": {
					"artifact": "scratch/gate.json",
					"gate": "prompt('Aprova?')"
				}
			},
			{
				"id": "tarefa-hermetic",
				"type": "script",
				"command": ["echo", "hermetic"],
				"routing": {"mode": "manual", "model": "none"},
				"depends_on": ["gate-aprovacao"],
				"output": {
					"artifact": "out2.txt",
					"gate": "file_exists('out2.txt') && exit_code == 0"
				}
			}
		]
	}`
	d, err := dag.Parse([]byte(src))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	return d
}

type fakeParidadeExec struct{}

func (f *fakeParidadeExec) Execute(ctx context.Context, n dag.Node, workdir string) (Result, error) {
	if n.Output.Artifact != "" {
		p := filepath.Join(workdir, n.Output.Artifact)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return Result{}, err
		}
		if err := os.WriteFile(p, []byte("ok "+n.ID), 0o644); err != nil {
			return Result{}, err
		}
	}
	return Result{
		ExitCode: 0,
		Tokens:   TokenUsage{Input: 100, Output: 50},
		Model:    "modelo-fixo",
		Notes:    []string{"executado " + n.ID},
	}, nil
}

func normalizarRelatorio(r *Report) {
	for i := range r.Nodes {
		r.Nodes[i].Duration = 0
	}
}

func verificarGoldenRelatorio(t *testing.T, rep Report) {
	t.Helper()
	dados, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		t.Fatalf("Marshal falhou: %v", err)
	}
	goldenPath := filepath.Join("testdata", "paridade_report.golden.json")
	esperado, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("falha ao ler golden %s: %v", goldenPath, err)
	}
	if string(dados) != string(esperado) {
		t.Errorf("relatorio diverge do golden:\nobtido:\n%s\nesperado:\n%s", string(dados), string(esperado))
	}
}

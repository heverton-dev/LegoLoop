package runner_test

import (
	"context"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/dag"
	"github.com/heverton-dev/LegoLoop/internal/orchestrix"
	"github.com/heverton-dev/LegoLoop/internal/runner"
)

func TestSBARDeRecusaPreExecucao(t *testing.T) {
	node := dag.Node{
		ID:      "step_script_recusado",
		Type:    dag.TypeScript,
		Command: []string{"echo", "nao-deve-rodar"},
		Output:  dag.HandoffContract{Artifact: "out.txt", Gate: "exit_code == 0"},
	}
	ox := &orchestrix.Orchestrix{AllowScripts: false}

	_, sbar, err := runner.ExecuteWithAndon(context.Background(), ox, node, t.TempDir())
	if err == nil {
		t.Fatal("esperava recusa por falta de --allow-scripts, mas obteve nil")
	}
	if !strings.Contains(err.Error(), "--allow-scripts") {
		t.Fatalf("erro deve explicar falta de --allow-scripts, obteve: %v", err)
	}

	bgLower := strings.ToLower(sbar.Background)
	if !strings.Contains(bgLower, "não foi executado") {
		t.Errorf("BACKGROUND deve indicar que o nó não foi executado, obteve: %q", sbar.Background)
	}
	if strings.Contains(bgLower, "exit code 0") {
		t.Errorf("BACKGROUND não pode afirmar 'exit code 0' em recusa pré-execução, obteve: %q", sbar.Background)
	}

	formatted := sbar.Format()
	if strings.Contains(strings.ToLower(formatted), "exit code 0") {
		t.Errorf("SBAR formatado não pode conter 'exit code 0', obteve:\n%s", formatted)
	}
}

func TestFormatNodeErrorComExitCodePositivo(t *testing.T) {
	sbar := runner.FormatNodeError("step_cmd", 2, "comando falhou")
	if !strings.Contains(sbar.Background, "exit code 2") {
		t.Errorf("BACKGROUND deve informar 'exit code 2' quando o nó executou e falhou, obteve: %q", sbar.Background)
	}
}

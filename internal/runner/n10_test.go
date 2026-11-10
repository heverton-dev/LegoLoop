package runner_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"legoloop/internal/runner"
)

func TestRedTeamBlueTeamIsolation(t *testing.T) {
	t.Run("BlueTeamProibidoDeAlterarTestes", func(t *testing.T) {
		contextIn := []string{"internal/calc/calc_test.go"}
		contextOut := []string{"internal/calc/calc.go", "internal/calc/calc_test.go"}

		err := runner.ValidateRedBlueBoundary("blue_team", contextIn, contextOut)
		if err == nil {
			t.Fatal("esperava erro ao tentar modificar arquivo de teste no blue_team, mas passou")
		}
		if !strings.Contains(err.Error(), "reward_hacking") {
			t.Errorf("erro deve mencionar reward_hacking, obteve: %v", err)
		}
	})

	t.Run("BlueTeamValidoSemModificarTestes", func(t *testing.T) {
		contextIn := []string{"internal/calc/calc_test.go"}
		contextOut := []string{"internal/calc/calc.go"}

		err := runner.ValidateRedBlueBoundary("blue_team", contextIn, contextOut)
		if err != nil {
			t.Fatalf("esperava sucesso para blue_team sem modificar testes, obteve: %v", err)
		}
	})
}

func TestWorktreeAstMerge(t *testing.T) {
	baseGo := []byte(`package sample

func First() int {
	return 1
}
`)
	incomingGo := []byte(`package sample

func Second() int {
	return 2
}
`)

	merged, err := runner.MergeAST(baseGo, incomingGo)
	if err != nil {
		t.Fatalf("falha ao mesclar AST: %v", err)
	}

	mergedStr := string(merged)
	if !strings.Contains(mergedStr, "func First()") || !strings.Contains(mergedStr, "func Second()") {
		t.Fatalf("AST mesclada deve conter First e Second, obteve:\n%s", mergedStr)
	}
}

func TestFlightRecordOTel(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "flight_record.json")

	rec := runner.FlightRecord{
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
		PipelineName: "test-processo",
		TotalTokens: runner.TokenUsage{
			Input:  1500,
			Output: 400,
		},
		NodesPassed:    3,
		NodesFailed:    0,
		FirstPassYield: 1.0,
	}

	if err := rec.SaveFlightRecord(outPath); err != nil {
		t.Fatalf("falha ao salvar flight_record: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("falha ao ler flight_record gravado: %v", err)
	}

	if !strings.Contains(string(data), `"first_pass_yield": 1`) {
		t.Errorf("conteúdo gravado deve conter first_pass_yield, obteve:\n%s", string(data))
	}
}

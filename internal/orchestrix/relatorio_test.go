package orchestrix

import (
	"context"
	"strings"
	"testing"

	"legoloop/internal/roteamento"
	"legoloop/internal/runner"
)

func TestRelatorioFactual(t *testing.T) {
	nos := []RelatorioNo{
		{ID: "TICK-1", Commit: "abc1234", CustoUSD: 0.25, Gate: "PASS", Harness: "claude", Modelo: "haiku"},
		{ID: "TICK-2", Commit: "def5678", CustoUSD: 0.5, Gate: "PASS", Harness: "claude", Modelo: "sonnet"},
		{ID: "TICK-3", Gate: "FAIL", Harness: "claude", Modelo: "opus"},
	}
	divs := []RelatorioDivergencia{{No: "TICK-3", Descricao: "modelo diferente do PLAN"}}
	out := GerarRelatorioFactual(nos, divs)

	linhas := map[string]string{
		"TICK-1": "| TICK-1 | abc1234 | 0.2500 | PASS | claude | - | haiku | - |",
		"TICK-2": "| TICK-2 | def5678 | 0.5000 | PASS | claude | - | sonnet | - |",
		"TICK-3": "| TICK-3 | - | 0.0000 | FAIL | claude | - | opus | - |",
	}
	for id, linha := range linhas {
		if !strings.Contains(out, linha) {
			t.Errorf("linha do nó %s ausente:\n%s", id, out)
		}
	}
	if !strings.Contains(out, "Custo total (USD): 0.7500") {
		t.Errorf("custo total ausente:\n%s", out)
	}
	if !strings.Contains(out, "- TICK-3: modelo diferente do PLAN") {
		t.Errorf("divergência ausente:\n%s", out)
	}
	for _, secao := range []string{"Riscos", "Próximo passo", "Divergências aceitas"} {
		if strings.Contains(out, secao) {
			t.Errorf("seção interpretativa %q não deve ser escrita", secao)
		}
	}
}

func TestRelatorioFactualSemDivergencia(t *testing.T) {
	out := GerarRelatorioFactual(nil, nil)
	if !strings.Contains(out, "Nenhuma.") {
		t.Errorf("esperava marcador de ausência de divergências:\n%s", out)
	}
}

func TestRelatorioFactualEscapaCelula(t *testing.T) {
	out := GerarRelatorioFactual([]RelatorioNo{{ID: "a|b", Gate: "PASS"}}, nil)
	if !strings.Contains(out, "a\\|b") {
		t.Errorf("barra vertical não escapada:\n%s", out)
	}
}

// TestRelatorioComHarnessReal roda um nó degradado pelo caminho real e confere a tabela factual.
func TestRelatorioComHarnessReal(t *testing.T) {
	var avisos []string
	res, err := orchestrixDegrada(&avisos).Execute(context.Background(), noDaClasse("TICK-D", roteamento.ClasseMicro), newRepoT(t))
	if err != nil {
		t.Fatal(err)
	}
	rep := runner.Report{OK: true, Nodes: []runner.NodeReport{
		{ID: "TICK-D", Status: runner.StatusPassed, Tokens: runner.TokenUsage{CostUSD: 0.1234}, Model: res.Model, Harness: res.Harness, Via: res.Via, Notes: res.Notes},
		{ID: "TICK-S", Status: runner.StatusPassed, Model: "sonnet", Harness: "claude", Via: "oficial"},
	}}
	out := GerarRelatorioFactual(RelatorioDoRunner(rep, map[string]string{"TICK-S": "abc1234"}), nil)
	if !strings.Contains(out, "| TICK-D | - | 0.1234 | passed | omp | 9router | omp:y | TICK-D usou omp/9router/omp:y no lugar de omp/oficial/omp:x: ") {
		t.Fatalf("linha do nó degradado ausente:\n%s", out)
	}
	if !strings.Contains(out, "| TICK-S | abc1234 | 0.0000 | passed | claude | oficial | sonnet | - |") {
		t.Fatalf("linha do nó sem troca ausente:\n%s", out)
	}
}

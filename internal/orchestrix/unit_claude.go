package orchestrix

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"legoloop/internal/dag"
	"legoloop/internal/runner"
)

// unitTools é a lista fechada de ferramentas de um unit_task. Não contém Bash (o `command`
// do nó roda fora do claude) nem a ferramenta de agentes (sem sub-subagentes).
var unitTools = []string{"Read", "Edit", "Write", "Glob", "Grep"}

// unitDeniedTools bloqueia explicitamente a ferramenta de agentes, mesmo que outra flag a libere.
var unitDeniedTools = []string{"Agent", "Task"}

// checkUnitTools recusa lista vazia e qualquer ferramenta fora de unitTools.
func checkUnitTools(tools []string) error {
	if len(tools) == 0 {
		return errors.New("lista de ferramentas vazia")
	}
	ok := map[string]bool{}
	for _, t := range unitTools {
		ok[t] = true
	}
	for _, t := range tools {
		if !ok[t] {
			return fmt.Errorf("ferramenta %q fora da lista permitida (%s)", t, strings.Join(unitTools, ","))
		}
	}
	return nil
}

// unitClaudeArgs monta os argumentos exatos do `claude -p` de um unit_task.
func unitClaudeArgs(model string, tools []string, lim unitLimits) []string {
	list := strings.Join(tools, ",")
	return []string{"-p", "--model", model, "--output-format", "json", "--no-session-persistence",
		"--strict-mcp-config", "--safe-mode", "--setting-sources", "user",
		"--tools", list, "--allowed-tools", list, "--disallowed-tools", strings.Join(unitDeniedTools, ","),
		"--disable-slash-commands", "--effort", lim.effort, "--max-budget-usd", lim.budget}
}

// unitLimits limita o raciocínio (effort) e o gasto em dólares de cada chamada do claude.
type unitLimits struct {
	effort string
	budget string
}

// unitLimitsFor escolhe o limite pela natureza do nó: micro e nano raciocinam pouco.
func unitLimitsFor(n dag.Node) unitLimits {
	switch n.Nature {
	case "micro", "nano":
		return unitLimits{effort: "low", budget: "0.80"}
	}
	return unitLimits{effort: "medium", budget: "2.50"}
}

// callClaudeTools roda `claude -p` DENTRO da worktree com ferramentas de arquivo fechadas.
// Um `tools` nulo usa unitTools. Prompt via stdin. A worktree é o único cwd do processo.
func (o *Orchestrix) callClaudeTools(ctx context.Context, model, prompt, worktree string, tools []string) (string, runner.TokenUsage, error) {
	return o.callClaudeToolsLim(ctx, model, prompt, worktree, tools, unitLimitsFor(dag.Node{}))
}

// callClaudeToolsLim é o callClaudeTools com limite de esforço e de gasto explícito.
func (o *Orchestrix) callClaudeToolsLim(ctx context.Context, model, prompt, worktree string, tools []string, lim unitLimits) (string, runner.TokenUsage, error) {
	if worktree == "" {
		return "", runner.TokenUsage{}, errors.New("callClaudeTools: worktree vazia")
	}
	if tools == nil {
		tools = unitTools
	}
	if err := checkUnitTools(tools); err != nil {
		return "", runner.TokenUsage{}, fmt.Errorf("callClaudeTools: %w", err)
	}
	out, code, err := o.spawn(ctx, o.Inventory.CLIPath("claude"), unitClaudeArgs(model, tools, lim), worktree, strings.NewReader(prompt))
	if err != nil {
		return "", runner.TokenUsage{}, err
	}
	r := parseClaudeOutput(out)
	if code != 0 || r.isError {
		return "", r.usage, fmt.Errorf("claude saiu com código %d: %s", code, r.reason())
	}
	return r.text, r.usage, nil
}

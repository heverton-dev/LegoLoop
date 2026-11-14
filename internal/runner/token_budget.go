package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"legoloop/internal/dag"
	"legoloop/internal/triagem"
)

func estimatePromptTokens(n dag.Node, workdir string) int {
	totalBytes := 0
	for _, v := range n.Inputs {
		totalBytes += len(v)
	}
	for _, rel := range n.ContextIn {
		if data, err := os.ReadFile(filepath.Join(workdir, rel)); err == nil {
			totalBytes += len(data)
		}
	}
	if totalBytes == 0 {
		return 0
	}
	tokens := totalBytes / 4
	if tokens == 0 {
		return 1
	}
	return tokens
}

func parseMaxPromptTokens(budget map[string]any) int {
	if budget == nil {
		return 0
	}
	raw, ok := budget["max_prompt_tokens"]
	if !ok {
		return 0
	}
	switch v := raw.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}

func budgetPolicy(budget map[string]any) string {
	if budget == nil {
		return "abort"
	}
	if p, ok := budget["policy"].(string); ok && p != "" {
		return strings.ToLower(p)
	}
	return "abort"
}

func truncateNodeInputs(n *dag.Node, maxTokens int) {
	if len(n.Inputs) == 0 || maxTokens <= 0 {
		return
	}
	limit := maxTokens * 4
	cloned := make(map[string]string, len(n.Inputs))
	for k, v := range n.Inputs {
		if len(v) > limit {
			cloned[k] = v[:limit]
		} else {
			cloned[k] = v
		}
	}
	n.Inputs = cloned
}

func recordBudgetTelemetry(memDir string, tel triagem.TokenBudgetTelemetry) {
	if memDir == "" {
		return
	}
	mem, err := triagem.LerMemoria(memDir)
	if err != nil {
		mem = triagem.MemoriaExecucao{
			EtapaAtual:   tel.NodeID,
			ProximoPasso: tel.Action,
		}
	}
	mem.TokenBudget = &tel
	_ = triagem.SalvarMemoria(memDir, mem)
}

func enforceTokenBudget(n *dag.Node, workdir, memDir string) error {
	maxTokens := parseMaxPromptTokens(n.TokenBudget)
	if maxTokens <= 0 {
		return nil
	}
	est := estimatePromptTokens(*n, workdir)
	tel := triagem.TokenBudgetTelemetry{
		NodeID:          n.ID,
		EstimatedTokens: est,
		MaxPromptTokens: maxTokens,
		Action:          "pass",
	}
	if est <= maxTokens {
		recordBudgetTelemetry(memDir, tel)
		return nil
	}
	policy := budgetPolicy(n.TokenBudget)
	if policy == "truncate" || policy == "compact" {
		truncateNodeInputs(n, maxTokens)
		tel.Action = "truncate"
		recordBudgetTelemetry(memDir, tel)
		return nil
	}
	tel.Action = "abort"
	recordBudgetTelemetry(memDir, tel)
	return fmt.Errorf("andon cord: token_budget excedido no nó %s (estimado=%d, max_prompt_tokens=%d)", n.ID, est, maxTokens)
}

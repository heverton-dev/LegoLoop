package runner

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/dag"
	"github.com/heverton-dev/LegoLoop/internal/triagem"
)

// novaSessaoID gera o identificador único da sessão (32 hex de crypto/rand).
func novaSessaoID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fmt.Sprintf("sess-%d", os.Getpid())
	}
	return hex.EncodeToString(b[:])
}

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

// recordBudgetTelemetry grava a telemetria na memória da sessão da execução.
func recordBudgetTelemetry(memDir, sessaoID string, tel triagem.TokenBudgetTelemetry) {
	if memDir == "" {
		return
	}
	mem, err := triagem.LerMemoria(memDir, sessaoID)
	if err != nil {
		mem = triagem.MemoriaExecucao{
			EtapaAtual:   tel.NodeID,
			ProximoPasso: tel.Action,
		}
	}
	mem.SessaoID = sessaoID
	mem.PID = os.Getpid()
	mem.TokenBudget = &tel
	if err := triagem.SalvarMemoria(memDir, mem); err != nil {
		log.Printf("token_budget: falha ao salvar memoria em %s: %v", memDir, err)
	}
}

func enforceTokenBudget(n *dag.Node, workdir, memDir, sessaoID string) error {
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
		recordBudgetTelemetry(memDir, sessaoID, tel)
		return nil
	}
	policy := budgetPolicy(n.TokenBudget)
	if policy == "truncate" || policy == "compact" {
		truncateNodeInputs(n, maxTokens)
		tel.Action = "truncate"
		recordBudgetTelemetry(memDir, sessaoID, tel)
		return nil
	}
	tel.Action = "abort"
	recordBudgetTelemetry(memDir, sessaoID, tel)
	return fmt.Errorf("andon cord: token_budget excedido no nó %s (estimado=%d, max_prompt_tokens=%d)", n.ID, est, maxTokens)
}

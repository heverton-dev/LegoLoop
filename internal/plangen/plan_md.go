package plangen

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"legoloop/docs/templates"
)

// renderPlan materializa o TEMPLATE_PLAN (mesmas seções, mesma ordem). Cada nó vem da mesma fonte do ROUTE,
// então os dois documentos são espelhos por construção.
func renderPlan(m planMeta, nodes []routeNode) string {
	var b strings.Builder
	writePlanHeader(&b, m)
	writePlanStrategy(&b, m)
	writeEpicTable(&b, m, nodes)
	b.WriteString("## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)\n\n")
	writePlanNodes(&b, m, nodes)
	writePlanFinal(&b, nodes)
	return b.String()
}

func writePlanHeader(b *strings.Builder, m planMeta) {
	fmt.Fprintf(b, "# 🗺️ PLAN: %s\n\n", m.Theme)
	b.WriteString("> **Status:** Em Planejamento  \n")
	fmt.Fprintf(b, "> **Data:** %s  \n> **Versão:** 1.0.0  \n", m.Date)
	b.WriteString("> **Autor:** Antigravity AI + Usuário  \n")
	if m.IntentRef != "" {
		fmt.Fprintf(b, "> **Referência da Intenção:** [%s](%s)  \n", m.IntentRef, m.IntentLink)
	} else {
		b.WriteString("> **Referência da Intenção:** Não informada na SPEC  \n")
	}
	fmt.Fprintf(b, "> **Referência da Especificação:** [%s](%s)  \n", m.SpecRel, m.SpecLink)
	fmt.Fprintf(b, "> **Espelho Executável JSON:** [%s](%s)  \n", m.RouteRel, m.RouteLink)
	fmt.Fprintf(b, "> **Template:** TEMPLATE_PLAN.md v%s (espelha TEMPLATE_ROUTE.json v%s)  \n\n---\n\n", templates.Version, templates.Version)
}

func writePlanStrategy(b *strings.Builder, m planMeta) {
	b.WriteString("## 1. Estratégia de Execução & Governança de Loteamento\n\n")
	b.WriteString("* **Engine:** `orchestrix_hermetic` (execução isolada por bloco, sem context bloat)\n")
	fmt.Fprintf(b, "* **Nível de Supervisão:** `%s`\n", m.Level)
	b.WriteString("* **Permissão de Subagentes:** `false` (AGENTS.md, regra 4: agentes de harness proibidos; subprocessos estéreis)\n")
	b.WriteString("* **TDD Estrito:** `true` (Red ➔ Green ➔ Refactor obrigatório antes de qualquer código de produção)\n")
	fmt.Fprintf(b, "* **Limite de Workers Paralelos:** `2`\n* **Timeout Padrão por Bloco:** `%ds`\n", defaultTimeout)
	b.WriteString("* **Política de Rollback:** global `abort_pipeline`; por nó `stash_changes`\n\n---\n\n")
}

func writeEpicTable(b *strings.Builder, m planMeta, nodes []routeNode) {
	b.WriteString("## 2. Mapa Geral de Épicos & Rastreabilidade de Critérios\n\n")
	b.WriteString("| Épico ID | Nome do Épico | Total Tickets | Critérios de Aceite Cobertos | Natureza Dominante |\n|---|---|---|---|---|\n")
	for _, epic := range epicsOf(nodes) {
		total, acs, nature := epicSummary(nodes, epic)
		fmt.Fprintf(b, "| **%s** | %s | %d | %s | `%s` |\n", epic, m.epicName(epic), total, acs, nature)
	}
	b.WriteString("\n> Regra de cobertura: a união da coluna \"Critérios de Aceite Cobertos\" deve ser igual ao conjunto de `AC-*` da SPEC.\n\n---\n\n")
}

func epicsOf(nodes []routeNode) []string {
	var order []string
	for _, n := range nodes {
		if n.Type == "unit_task" && !slices.Contains(order, n.epic) {
			order = append(order, n.epic)
		}
	}
	return order
}

func epicSummary(nodes []routeNode, epic string) (int, string, string) {
	total, counts := 0, map[string]int{}
	var acs []string
	for _, n := range nodes {
		if n.Type != "unit_task" || n.epic != epic {
			continue
		}
		total++
		counts[n.Nature]++
		if n.SpecRef != nil {
			for c := range strings.SplitSeq(*n.SpecRef, ", ") {
				acs = append(acs, "`"+c+"`")
			}
		}
	}
	nature := "micro"
	if counts["nano"] > counts["micro"] {
		nature = "nano"
	}
	return total, strings.Join(acs, ", "), nature
}

func writePlanNodes(b *strings.Builder, m planMeta, nodes []routeNode) {
	current := ""
	for _, n := range nodes {
		if n.ID == finalGateID {
			continue
		}
		if n.epic != current {
			current = n.epic
			fmt.Fprintf(b, "### Épico %02d — %s\n\n", epicNumber(n.epic), m.epicName(n.epic))
		}
		if n.Type == "stop_gate" {
			writePlanGate(b, n)
			continue
		}
		writePlanTicket(b, n)
	}
}

func jsonOf(v any) string {
	out, _ := json.Marshal(v)
	return string(out)
}

func writePlanTicket(b *strings.Builder, n routeNode) {
	fmt.Fprintf(b, "#### Ticket `%s` — %s\n", n.ID, n.Label)
	fmt.Fprintf(b, "* **Nature:** `%s` | **Tipo:** `%s` | **Modo de Execução:** `%s` | **Status:** `%s`\n", n.Nature, n.Type, n.ExecutionMode, n.Status)
	fmt.Fprintf(b, "* **Dependências:** `%s`\n", jsonOf(n.DependsOn))
	fmt.Fprintf(b, "* **Isolamento:** Workspace `%s` · Executor `%s` · Subagentes `%t`\n", n.Isolation.Workspace, n.Isolation.Executor, n.Isolation.SubagentsAllowed)
	fmt.Fprintf(b, "* **Roteamento:** Harness `%s` · Modelo `%s` · Effort `%s` · Fallback Harness `%s` · Modelo `%s`\n",
		n.Routing.Harness, n.Routing.Model, *n.Routing.Effort, n.Routing.Fallback.Harness, n.Routing.Fallback.Model)
	fmt.Fprintf(b, "* **Tooling:** Runtime `%s` · MCPs Requeridos `%t` · Allowlist `%s`\n", n.Tooling.Runtime, n.Tooling.McpsRequired, jsonOf(n.Tooling.McpsAllowed))
	fmt.Fprintf(b, "* **Limites:** Timeout `%ds` · Teto de saída `%d` tokens · Rollback `%s`\n", n.TimeoutSec, n.TokenBudget.Max, n.RollbackStrategy)
	fmt.Fprintf(b, "* **Inputs:** Source Files `%s` · Spec Anchor `%s` · Context Artifacts `%s`\n", strings.Join(n.Inputs.SourceFiles, ", "), deref(n.Inputs.SpecAnchor), jsonOf(n.Inputs.ContextArtifacts))
	fmt.Fprintf(b, "* **Alvo de Teste TDD:** `%s`\n", deref(n.TestTarget))
	fmt.Fprintf(b, "* **Outputs:** Artifacts `%s`\n", strings.Join(n.Output.Artifacts, ", "))
	fmt.Fprintf(b, "* **Verificação (argv, sem shell):** %s\n", verifyText(n.Output.Verify))
	fmt.Fprintf(b, "* **Gate (DSL):** `%s`\n* **Spec Ref:** `%s`\n\n", n.Output.Gate, deref(n.SpecRef))
}

func deref(s *string) string {
	if s == nil {
		return "nenhum"
	}
	return *s
}

func verifyText(verify [][]string) string {
	if len(verify) == 0 {
		return "nenhuma (confirmação manual pelo gate)"
	}
	parts := make([]string, len(verify))
	for i, v := range verify {
		parts[i] = "`" + strings.Join(v, " ") + "`"
	}
	return strings.Join(parts, " · ")
}

func writePlanGate(b *strings.Builder, n routeNode) {
	fmt.Fprintf(b, "### 🛑 STOP-GATE HUMANO (`%s`)\n* **Gate (DSL):** `%s`\n", n.ID, n.Output.Gate)
	b.WriteString("- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).\n")
	b.WriteString("- [ ] Evidências de testes unitários verdes apresentadas ao usuário.\n")
	b.WriteString("- [ ] Autorização explícita do usuário concedida para prosseguir.\n\n---\n\n")
}

func writePlanFinal(b *strings.Builder, nodes []routeNode) {
	b.WriteString("## 4. Homologação Final & Fechamento da Demanda (`SG-FINAL`)\n\n")
	b.WriteString("Ao concluir todos os épicos, executar nesta ordem (cada comando é um argv, sem shell):\n")
	b.WriteString("1. Suíte de regressão global: `go test ./...`\n")
	b.WriteString("2. Estática e formatação: `go vet ./...` e `gofmt -l .` (saída vazia). `golangci-lint run ./...` apenas se instalado e configurado.\n")
	fmt.Fprintf(b, "3. Recompilação obrigatória do binário: `go build -o %s ./cmd/legoloop`\n", finalBinary)
	b.WriteString("4. Emissão do relatório formal em `docs/reports/`\n")
	b.WriteString("5. Registro de histórico da sessão via `faz-session --id <ID> --format all`\n\n")
	fmt.Fprintf(b, "* **Gate (DSL):** `%s`\n\n---\n\n", nodes[len(nodes)-1].Output.Gate)
	b.WriteString("🛑 **STOP-GATE 03 DE AUDITORIA HUMANA:**  \n")
	fmt.Fprintf(b, "*Plano de execução fatiado com %d tickets. Deseja aprovar o plano e liberar a execução autônoma do BUILD?*  \n", countTickets(nodes))
	b.WriteString("*Aprovado este documento e gerado o ROUTE correspondente, o BUILD pode ser iniciado.*\n")
}

func countTickets(nodes []routeNode) int {
	n := 0
	for _, x := range nodes {
		if x.Type == "unit_task" {
			n++
		}
	}
	return n
}

---
id: PLAN-2026-11-21_ROUTE_INICIATIVA_3_I03_DO_MACROPROJE
status: DRAFT
date: 2026-11-21
schema_version: 2
upstream: [2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE]
covers: [AC-E01.1, AC-E01.2, AC-E01.3, AC-E01.4, AC-E01.5, AC-E02.1, AC-E02.2, AC-E02.3, AC-E03.1, AC-E03.2, AC-E03.3, AC-E03.4, AC-E03.5, AC-E03.6, AC-E04.1, AC-E04.2, AC-E04.3, AC-E05.1, AC-E05.2, AC-E05.3, AC-E05.4, AC-NF.1, AC-NF.2]
---

# 🗺️ PLAN: Despachante de executores no runner (I03)

> **Status:** Em Planejamento  
> **Data:** 2026-11-21  
> **Versão:** 1.0.0  
> **Autor:** Antigravity AI + Usuário  
> **Referência da Intenção:** [docs/intents/2026-11-21_INTENT_INICIATIVA_3_I03_DO_MACROPROJE.md](../intents/2026-11-21_INTENT_INICIATIVA_3_I03_DO_MACROPROJE.md)  
> **Referência da Especificação:** [docs/specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md](../specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-21_ROUTE_INICIATIVA_3_I03_DO_MACROPROJE.json](../routes/2026-11-21_ROUTE_INICIATIVA_3_I03_DO_MACROPROJE.json)  
> **Template:** TEMPLATE_PLAN.md v2.2.0 (espelha TEMPLATE_ROUTE.json v2.2.0)  

---

## 1. Estratégia de Execução & Governança de Loteamento

* **Engine:** `orchestrix_hermetic` (execução isolada por Peça, sem context bloat)
* **Nível de Supervisão:** `epic`
* **Permissão de Subagentes:** `false` ou `ephemeral` (AGENTS.md, regra 4: agente de harness proibido; subagente efêmero só com as seis condições)
* **TDD Estrito:** `true` (Red ➔ Green ➔ Refactor obrigatório antes de qualquer código de produção)
* **Limite de Workers Paralelos:** `2`
* **Timeout Padrão por Peça:** `180s`
* **Política de Rollback:** global `abort_processo`; por nó `stash_changes`

---

## 2. Mapa Geral de Épicos & Rastreabilidade de Critérios

| Épico ID | Nome do Épico | Total Tickets | Critérios de Aceite Cobertos | Natureza Dominante |
|---|---|---|---|---|
| **E01** | Vocabulário e validação (`internal/dag`) | 4 | `AC-E01.1`, `AC-E01.2`, `AC-E01.3`, `AC-E01.4`, `AC-E01.5` | `micro` |
| **E02** | Contrato no `plangen` e no `routedag` | 3 | `AC-E02.1`, `AC-E02.2`, `AC-E02.3` | `micro` |
| **E03** | Despachante no `runner` | 6 | `AC-E03.1`, `AC-E03.2`, `AC-NF.1`, `AC-E03.3`, `AC-E03.4`, `AC-E03.5`, `AC-E03.6` | `micro` |
| **E04** | Ligação real com o `board` | 3 | `AC-E04.1`, `AC-E04.2`, `AC-E04.3` | `micro` |
| **E05** | Paridade e governança viva | 4 | `AC-E05.1`, `AC-E05.2`, `AC-E05.3`, `AC-E05.4`, `AC-NF.2` | `nano` |

> Regra de cobertura: a união da coluna "Critérios de Aceite Cobertos" deve ser igual ao conjunto de `AC-*` da SPEC.

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Vocabulário e validação (`internal/dag`)

#### Ticket `TICK-01.1` — Criar o vocabulário de executores, ErrExecutorDesconhecido e dag.ValidarExecutores que lista todos os nós inválidos
* **Iniciativa:** `I03` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/dag/executor.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md#AC-E01.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/dag/executor_test.go`
* **Outputs:** Artifacts `internal/dag/executor.go, internal/dag/executor_test.go`
* **Verificação (argv, sem shell):** nenhuma (confirmação manual pelo gate)
* **Gate (DSL):** `prompt('Execute manualmente e confirme exit code 0: go test ./internal/dag/ -run "TestExecutorAceito|TestValidarExecutoresListaTodos"')`
* **Spec Ref:** `AC-E01.1, AC-E01.2`

#### Ticket `TICK-01.2` — Recusar subagents_allowed ephemeral fora do native_swarm
* **Iniciativa:** `I03` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/dag/executor.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md#AC-E01.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/dag/executor_test.go`
* **Outputs:** Artifacts `internal/dag/executor.go, internal/dag/executor_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/dag/ -run TestSubagentesEfemerosSoNoSwarm`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/dag/executor.go') && file_exists('internal/dag/executor_test.go')`
* **Spec Ref:** `AC-E01.3`

#### Ticket `TICK-01.3` — Exigir papel e command e recusar stop_gate no native_swarm
* **Iniciativa:** `I03` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/dag/executor.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md#AC-E01.4` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/dag/executor_test.go`
* **Outputs:** Artifacts `internal/dag/executor.go, internal/dag/executor_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/dag/ -run TestNativeSwarmExigePapelECommand`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/dag/executor.go') && file_exists('internal/dag/executor_test.go')`
* **Spec Ref:** `AC-E01.4`

#### Ticket `TICK-01.4` — Fazer dag.Parse e legoloop validate recusarem executor desconhecido
* **Iniciativa:** `I03` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.3"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/dag/parse.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md#AC-E01.5` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/dag/executor_test.go`
* **Outputs:** Artifacts `internal/dag/parse.go, internal/dag/executor_test.go, internal/cli/validate_executor_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/dag/... ./internal/cli/ -run TestValidateRecusaExecutor`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/dag/parse.go') && file_exists('internal/dag/executor_test.go') && file_exists('internal/cli/validate_executor_test.go')`
* **Spec Ref:** `AC-E01.5`

### 🛑 STOP-GATE HUMANO (`SG-E01`)
* **Gate (DSL):** `prompt('Evidências do Épico 01 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 02 — Contrato no `plangen` e no `routedag`

#### Ticket `TICK-02.1` — Criar o campo executor do Ticket e gravá-lo na ROUTE
* **Iniciativa:** `I03` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.4","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/plangen/plangen.go, internal/plangen/route.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md#AC-E02.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/plangen/ticket_executor_test.go`
* **Outputs:** Artifacts `internal/plangen/plangen.go, internal/plangen/route.go, internal/plangen/ticket_executor_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/plangen/ -run TestTicketExecutor`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/plangen/plangen.go') && file_exists('internal/plangen/route.go') && file_exists('internal/plangen/ticket_executor_test.go')`
* **Spec Ref:** `AC-E02.1`

#### Ticket `TICK-02.2` — Recusar Ticket com executor desconhecido sem gravar PLAN nem ROUTE
* **Iniciativa:** `I03` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.1","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/plangen/plangen.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md#AC-E02.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/plangen/ticket_executor_test.go`
* **Outputs:** Artifacts `internal/plangen/plangen.go, internal/plangen/ticket_executor_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/plangen/ -run TestTicketExecutorDesconhecido`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/plangen/plangen.go') && file_exists('internal/plangen/ticket_executor_test.go')`
* **Spec Ref:** `AC-E02.2`

#### Ticket `TICK-02.3` — Provar que a projeção do DAG preserva native_swarm e recusa valor desconhecido
* **Iniciativa:** `I03` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.2","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/routedag/routedag.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md#AC-E02.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/routedag/executor_test.go`
* **Outputs:** Artifacts `internal/routedag/routedag.go, internal/routedag/executor_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/routedag/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/routedag/routedag.go') && file_exists('internal/routedag/executor_test.go')`
* **Spec Ref:** `AC-E02.3`

### 🛑 STOP-GATE HUMANO (`SG-E02`)
* **Gate (DSL):** `prompt('Evidências do Épico 02 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 03 — Despachante no `runner`

#### Ticket `TICK-03.1` — Criar Options.Swarm e a validação de executores antes do nível 0 em runner.Run
* **Iniciativa:** `I03` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.3","SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/runner/runner.go, internal/runner/swarm.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md#AC-E03.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/runner/swarm_test.go`
* **Outputs:** Artifacts `internal/runner/runner.go, internal/runner/swarm.go, internal/runner/swarm_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/runner/ -run TestRunRecusaAntesDoNivelZero`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/runner/runner.go') && file_exists('internal/runner/swarm.go') && file_exists('internal/runner/swarm_test.go')`
* **Spec Ref:** `AC-E03.1`

#### Ticket `TICK-03.2` — Despachar por nível: native_swarm ao SwarmRunner e os demais ao Executor
* **Iniciativa:** `I03` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.1","SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/runner/runner.go, internal/runner/swarm.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md#AC-E03.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/runner/swarm_test.go`
* **Outputs:** Artifacts `internal/runner/runner.go, internal/runner/swarm.go, internal/runner/swarm_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/runner/ -run TestDespachoPorNivel`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/runner/runner.go') && file_exists('internal/runner/swarm.go') && file_exists('internal/runner/swarm_test.go')`
* **Spec Ref:** `AC-E03.2, AC-NF.1`

#### Ticket `TICK-03.3` — Rodar o nível misto em paralelo e terminar só quando os dois lados terminam
* **Iniciativa:** `I03` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.2","SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/runner/runner.go, internal/runner/swarm.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md#AC-E03.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/runner/swarm_test.go`
* **Outputs:** Artifacts `internal/runner/runner.go, internal/runner/swarm.go, internal/runner/swarm_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/runner/ -run TestNivelMistoEsperaOsDoisLados`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/runner/runner.go') && file_exists('internal/runner/swarm.go') && file_exists('internal/runner/swarm_test.go')`
* **Spec Ref:** `AC-E03.3`

#### Ticket `TICK-03.4` — Converter o resultado do Swarm em NodeReport e marcar os níveis seguintes como skipped
* **Iniciativa:** `I03` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.3","SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/runner/runner.go, internal/runner/swarm.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md#AC-E03.4` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/runner/swarm_test.go`
* **Outputs:** Artifacts `internal/runner/runner.go, internal/runner/swarm.go, internal/runner/swarm_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/runner/ -run TestFalhaNoSwarmMarcaSkipped`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/runner/runner.go') && file_exists('internal/runner/swarm.go') && file_exists('internal/runner/swarm_test.go')`
* **Spec Ref:** `AC-E03.4`

#### Ticket `TICK-03.5` — Garantir que DAG sem native_swarm não toca no Swarm
* **Iniciativa:** `I03` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.4","SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/runner/runner.go, internal/runner/swarm.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md#AC-E03.5` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/runner/swarm_test.go`
* **Outputs:** Artifacts `internal/runner/runner.go, internal/runner/swarm.go, internal/runner/swarm_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/runner/ -run TestSemNativeSwarmNaoUsaSwarm`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/runner/runner.go') && file_exists('internal/runner/swarm.go') && file_exists('internal/runner/swarm_test.go')`
* **Spec Ref:** `AC-E03.5`

#### Ticket `TICK-03.6` — Provar que legoloop run recusa executor inválido com exit 1 e sem criar worktree
* **Iniciativa:** `I03` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.5","SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/cli/commands.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md#AC-E03.6` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/cli/run_executor_test.go`
* **Outputs:** Artifacts `internal/cli/commands.go, internal/cli/run_executor_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/cli/ -run TestRunRecusaExecutorSemWorktree`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/commands.go') && file_exists('internal/cli/run_executor_test.go')`
* **Spec Ref:** `AC-E03.6`

### 🛑 STOP-GATE HUMANO (`SG-E03`)
* **Gate (DSL):** `prompt('Evidências do Épico 03 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 04 — Ligação real com o `board`

#### Ticket `TICK-04.1` — Extrair board.RunNodeOf de LevelsFromDAG sem mudar o comportamento
* **Iniciativa:** `I03` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.6","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/board/swarm_dag.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md#AC-E04.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/board/swarm_dag_test.go`
* **Outputs:** Artifacts `internal/board/swarm_dag.go, internal/board/swarm_dag_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/board/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/board/swarm_dag.go') && file_exists('internal/board/swarm_dag_test.go')`
* **Spec Ref:** `AC-E04.1`

#### Ticket `TICK-04.2` — Ligar legoloop run ao SwarmBoard real quando houver nó native_swarm
* **Iniciativa:** `I03` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.1","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/cli/commands.go, internal/runner/swarm.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md#AC-E04.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/cli/run_swarm_test.go`
* **Outputs:** Artifacts `internal/cli/commands.go, internal/runner/swarm.go, internal/cli/run_swarm_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/cli/ -run TestRunNativeSwarmBinario`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/commands.go') && file_exists('internal/runner/swarm.go') && file_exists('internal/cli/run_swarm_test.go')`
* **Spec Ref:** `AC-E04.2`

#### Ticket `TICK-04.3` — Provar que sem native_swarm o diretório do Swarm não é criado
* **Iniciativa:** `I03` | **Épico:** `E04` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.2","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/cli/run_swarm_test.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md#AC-E04.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/cli/run_swarm_test.go`
* **Outputs:** Artifacts `internal/cli/run_swarm_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/cli/ -run TestRunSemNativeSwarmNaoCriaDiretorio`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/run_swarm_test.go')`
* **Spec Ref:** `AC-E04.3`

### 🛑 STOP-GATE HUMANO (`SG-E04`)
* **Gate (DSL):** `prompt('Evidências do Épico 04 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 05 — Paridade e governança viva

#### Ticket `TICK-05.1` — Provar que todos os DAGs do repositório continuam válidos
* **Iniciativa:** `I03` | **Épico:** `E05` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.3","SG-E04"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/dag/executor_repo_test.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md#AC-E05.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/dag/executor_repo_test.go`
* **Outputs:** Artifacts `internal/dag/executor_repo_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/dag/ -run TestDAGsDoRepositorioValidam`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/dag/executor_repo_test.go')`
* **Spec Ref:** `AC-E05.1`

#### Ticket `TICK-05.2` — Provar com golden que o relatório de DAG sem native_swarm não mudou
* **Iniciativa:** `I03` | **Épico:** `E05` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.1","SG-E04"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/runner/paridade_test.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md#AC-E05.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/runner/paridade_test.go`
* **Outputs:** Artifacts `internal/runner/paridade_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/runner/ -run TestParidadeDoRelatorio`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/runner/paridade_test.go')`
* **Spec Ref:** `AC-E05.2`

#### Ticket `TICK-05.3` — Atualizar AGENTS, README e MEMORY com o despachante e o native_swarm
* **Iniciativa:** `I03` | **Épico:** `E05` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.2","SG-E04"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `AGENTS.md, docs/governanca/AGENTS_COMPLETO.md, README.md, MEMORY.md` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md#AC-E05.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `AGENTS.md, docs/governanca/AGENTS_COMPLETO.md, README.md, MEMORY.md`
* **Verificação (argv, sem shell):** nenhuma (confirmação manual pelo gate)
* **Gate (DSL):** `prompt('Execute manualmente e confirme exit code 0: grep -c "native_swarm" AGENTS.md README.md MEMORY.md')`
* **Spec Ref:** `AC-E05.3`

#### Ticket `TICK-05.4` — Registrar os anticorpos, fechar a Emenda 1, gerar o prompt da I03 e rodar o portão completo
* **Iniciativa:** `I03` | **Épico:** `E05` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.3","SG-E04"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `docs/intents/2026-11-21_EMENDA_1_MACROPLANO.md, docs/prompts` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md#AC-E05.4` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `docs/intents/2026-11-21_EMENDA_1_MACROPLANO.md, docs/prompts`
* **Verificação (argv, sem shell):** `go test -count=1 ./...`
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/intents/2026-11-21_EMENDA_1_MACROPLANO.md') && file_exists('docs/prompts')`
* **Spec Ref:** `AC-E05.4, AC-NF.2`

## 4. Homologação Final & Fechamento da Demanda (`SG-FINAL`)

Ao concluir todos os épicos, executar nesta ordem (cada comando é um argv, sem shell):
1. Suíte de regressão global: `go test ./...`
2. Estática e formatação: `go vet ./...` e `gofmt -l .` (saída vazia). `golangci-lint run ./...` apenas se instalado e configurado.
3. Recompilação obrigatória do binário: `go build -o bin/legoloop.exe ./cmd/legoloop`
4. Emissão do relatório formal em `docs/reports/`
5. Registro de histórico da sessão via `faz-session --id <ID> --format all`

* **Gate (DSL):** `exit_code == 0 && file_exists('bin/legoloop.exe')`

---

🛑 **STOP-GATE 03 DE AUDITORIA HUMANA:**  
*Plano de execução fatiado com 20 tickets. Deseja aprovar o plano e liberar a execução autônoma do BUILD?*  
*Aprovado este documento e gerado o ROUTE correspondente, o BUILD pode ser iniciado.*

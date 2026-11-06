# 🗺️ PLAN: Alternador Universal e Granular de Governança, Economia de Tokens e Output STE-PT para Harnesses

> **Status:** APPROVED / EM EXECUÇÃO  
> **Data:** 2026-11-05 (revisado em 2026-11-06)  
> **Versão:** 2.0.0  
> **Autor:** Antigravity (Pair Programming com Usuário)  
> **Referência da Intenção:** [docs/intents/2026-11-05_INTENT_STEPT_SWITCH.md](../intents/2026-11-05_INTENT_STEPT_SWITCH.md)  
> **Referência da Especificação:** [docs/specs/2026-11-05_SPEC_STEPT_SWITCH.md](../specs/2026-11-05_SPEC_STEPT_SWITCH.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-05_ROUTE_STEPT_SWITCH.json](../routes/2026-11-05_ROUTE_STEPT_SWITCH.json)  
> **Template:** TEMPLATE_PLAN.md v2.0.0 (espelha TEMPLATE_ROUTE.json v2.0.0)  

---

## 1. Estratégia de Execução & Governança de Loteamento

* **Engine:** `orchestrix_hermetic` (execução isolada por bloco, sem context bloat)
* **Nível de Supervisão:** `epic` (um stop-gate humano por épico, mais o `SG-FINAL`)
* **Permissão de Subagentes:** `false` (AGENTS.md, regra 4: agentes de harness proibidos; subprocessos estéreis)
* **TDD Estrito:** `true` (Red ➔ Green ➔ Refactor; o teste vermelho é registrado antes da implementação e **cada teste crítico é provado por mutação**)
* **Limite de Workers Paralelos:** `1`
* **Timeout Padrão por Bloco:** `300s` (cada ticket declara o seu)
* **Política de Rollback:** global `abort_pipeline`; por nó `stash_changes`
* **Convenção de testes:** o nome em `test_target` é prefixo. Todo teste do ticket começa com ele e o `verify` roda `-run ^<Prefixo>`.
* **Verificação por ticket Go:** teste do ticket, `go vet` do pacote, `funclen` do pacote e `go build -o bin/legoloop.exe ./cmd/legoloop` (AGENTS.md, regra 5).
* **Prova pelo binário e mutações (AGENTS.md, regra 8):** suíte hermética de mutações deliberadas no Épico 04 cobrindo os 4 pontos críticos do sistema.
* **Matriz de Roteamento de 06/11:**
  * `macro`: Harness `claude` (modelo `claude-sonnet-5-5`), fallback `9router` (`code-pro`), effort `high`
  * `micro`: Harness `claude` (modelo `claude-sonnet-5-5`), fallback `9router` (`task-micro`), effort `medium`
  * `nano`: Harness `claude` (modelo `claude-haiku-5-5`), fallback `9router` (`code-tests`), effort `low`

---

## 2. Mapa Geral de Épicos & Rastreabilidade de Critérios

| Épico ID | Nome do Épico | Total Tickets | Critérios de Aceite Cobertos | Natureza Dominante |
|---|---|---|---|---|
| **E01** | Motor de Configuração, Persistência e Regras STE-PT | 3 | `AC-E01.1`, `AC-E01.2`, `AC-E01.3` | `micro` |
| **E02** | Scanner Dinâmico, Sincronizador de Harnesses & CLI Integrada | 3 | `AC-E02.1`, `AC-E02.2`, `AC-E02.3` | `micro` |
| **E03** | Projeto Standalone em Go no Desktop (stept-cli) | 3 | `AC-E03.1`, `AC-E03.2`, `AC-E03.3`, `AC-NF.1` | `micro` |
| **E04** | Testes de Mutação (Regra 8) & Homologação Final | 2 | `AC-E04.1`, `AC-E04.2`, `AC-E04.3`, `AC-E04.4`, `AC-NF.1` | `micro` |

> Regra de cobertura: a união da coluna "Critérios de Aceite Cobertos" é estritamente igual ao conjunto de `AC-*` da SPEC.

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Motor de Configuração, Persistência e Regras STE-PT

#### Ticket `TICK-01.1` — Modelos de Dados, Enums e Validação
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Modo `suggested` · Harness `claude` · Modelo `claude-sonnet-5-5` · Effort `medium` · Fallback Harness `9router` · Modelo `task-micro`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `300s` · Teto de saída `800` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/ste/switch.go` · Spec Anchor `docs/specs/2026-11-05_SPEC_STEPT_SWITCH.md#AC-E01.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/ste/switch_test.go:TestConfigPerfil_Validate`
* **Outputs:** Artifacts `internal/ste/switch.go`, `internal/ste/switch_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 -v -run ^TestConfigPerfil_Validate ./internal/ste/...`; `go vet ./internal/ste/...`; `go run ./tools/funclen internal/ste`; `go build -o bin/legoloop.exe ./cmd/legoloop`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/ste/switch.go') && file_exists('internal/ste/switch_test.go')`
* **Spec Ref:** `AC-E01.2`

#### Ticket `TICK-01.2` — Persistência JSON e Resolução de Herança
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Modo `suggested` · Harness `claude` · Modelo `claude-sonnet-5-5` · Effort `medium` · Fallback Harness `9router` · Modelo `task-micro`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `300s` · Teto de saída `800` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/ste/switch.go` · Spec Anchor `docs/specs/2026-11-05_SPEC_STEPT_SWITCH.md#AC-E01.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/ste/switch_test.go:TestConfig_ResolveEHeranca`
* **Outputs:** Artifacts `internal/ste/switch.go`, `internal/ste/switch_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 -v -run ^TestConfig_ResolveEHeranca ./internal/ste/...`; `go vet ./internal/ste/...`; `go run ./tools/funclen internal/ste`; `go build -o bin/legoloop.exe ./cmd/legoloop`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/ste/switch.go') && file_exists('internal/ste/switch_test.go')`
* **Spec Ref:** `AC-E01.1`

#### Ticket `TICK-01.3` — Geração do Prompt Canônico (Tríade Input + Caveman Thinking + Cartão STE-PT)
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Modo `suggested` · Harness `claude` · Modelo `claude-sonnet-5-5` · Effort `medium` · Fallback Harness `9router` · Modelo `task-micro`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `300s` · Teto de saída `800` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/ste/prompt.go` · Spec Anchor `docs/specs/2026-11-05_SPEC_STEPT_SWITCH.md#AC-E01.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/ste/prompt_test.go:TestPromptRegra`
* **Outputs:** Artifacts `internal/ste/prompt.go`, `internal/ste/prompt_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 -v -run ^TestPromptRegra ./internal/ste/...`; `go vet ./internal/ste/...`; `go run ./tools/funclen internal/ste`; `go build -o bin/legoloop.exe ./cmd/legoloop`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/ste/prompt.go') && file_exists('internal/ste/prompt_test.go')`
* **Spec Ref:** `AC-E01.3`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 01 (`SG-E01`)
* **Gate (DSL):** `prompt('Evidências do Épico 01 aprovadas? Autoriza iniciar o Épico 02?')`
- [ ] 100% dos tickets do Épico 01 executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários e de integração verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para iniciar o próximo Épico.

---

### Épico 02 — Scanner Dinâmico, Sincronizador de Harnesses & CLI Integrada

#### Ticket `TICK-02.1` — Scanner Dinâmico de Harnesses no SO
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Modo `suggested` · Harness `claude` · Modelo `claude-sonnet-5-5` · Effort `medium` · Fallback Harness `9router` · Modelo `task-micro`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `300s` · Teto de saída `800` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/ste/scan.go` · Spec Anchor `docs/specs/2026-11-05_SPEC_STEPT_SWITCH.md#AC-E02.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/ste/scan_test.go:TestScanHarnesses`
* **Outputs:** Artifacts `internal/ste/scan.go`, `internal/ste/scan_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 -v -run ^TestScanHarnesses ./internal/ste/...`; `go vet ./internal/ste/...`; `go run ./tools/funclen internal/ste`; `go build -o bin/legoloop.exe ./cmd/legoloop`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/ste/scan.go') && file_exists('internal/ste/scan_test.go')`
* **Spec Ref:** `AC-E02.1`

#### Ticket `TICK-02.2` — Sincronizador Atômico e Remoção Limpa por Harness
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Modo `suggested` · Harness `claude` · Modelo `claude-sonnet-5-5` · Effort `medium` · Fallback Harness `9router` · Modelo `task-micro`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `300s` · Teto de saída `800` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/ste/sync.go` · Spec Anchor `docs/specs/2026-11-05_SPEC_STEPT_SWITCH.md#AC-E02.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/ste/sync_test.go:TestSyncHarnesses`
* **Outputs:** Artifacts `internal/ste/sync.go`, `internal/ste/sync_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 -v -run ^TestSyncHarnesses ./internal/ste/...`; `go vet ./internal/ste/...`; `go run ./tools/funclen internal/ste`; `go build -o bin/legoloop.exe ./cmd/legoloop`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/ste/sync.go') && file_exists('internal/ste/sync_test.go')`
* **Spec Ref:** `AC-E02.2`

#### Ticket `TICK-02.3` — Subcomando CLI legoloop stept e Tabela de Status
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Modo `suggested` · Harness `claude` · Modelo `claude-sonnet-5-5` · Effort `medium` · Fallback Harness `9router` · Modelo `task-micro`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `300s` · Teto de saída `800` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/cli/stept.go` · Spec Anchor `docs/specs/2026-11-05_SPEC_STEPT_SWITCH.md#AC-E02.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/cli/stept_test.go:TestCLIStept`
* **Outputs:** Artifacts `internal/cli/stept.go`, `internal/cli/stept_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 -v -run ^TestCLIStept ./internal/cli/...`; `go vet ./internal/cli/...`; `go run ./tools/funclen internal/cli`; `go build -o bin/legoloop.exe ./cmd/legoloop`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/stept.go') && file_exists('internal/cli/stept_test.go')`
* **Spec Ref:** `AC-E02.3`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 02 (`SG-E02`)
* **Gate (DSL):** `prompt('Evidências do Épico 02 aprovadas? Autoriza iniciar o Épico 03?')`
- [ ] 100% dos tickets do Épico 02 executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários e de integração verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para iniciar o próximo Épico.

---

### Épico 03 — Projeto Standalone em Go no Desktop (stept-cli)

#### Ticket `TICK-03.1` — Estrutura e Código do Projeto Standalone em Go
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Modo `suggested` · Harness `claude` · Modelo `claude-sonnet-5-5` · Effort `medium` · Fallback Harness `9router` · Modelo `task-micro`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `300s` · Teto de saída `800` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `C:/Users/trcnologia/Desktop/stept-cli/go.mod`, `C:/Users/trcnologia/Desktop/stept-cli/cmd/stept/main.go` · Spec Anchor `docs/specs/2026-11-05_SPEC_STEPT_SWITCH.md#AC-E03.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `cmd/stept/main_test.go:TestStandaloneMain`
* **Outputs:** Artifacts `C:/Users/trcnologia/Desktop/stept-cli/go.mod`, `C:/Users/trcnologia/Desktop/stept-cli/cmd/stept/main.go`, `C:/Users/trcnologia/Desktop/stept-cli/cmd/stept/main_test.go`
* **Verificação (argv, sem shell):** `cd C:/Users/trcnologia/Desktop/stept-cli`; `go test -count=1 ./...`; `go build -o bin/stept.exe ./cmd/stept`
* **Gate (DSL):** `exit_code == 0 && file_exists('C:/Users/trcnologia/Desktop/stept-cli/go.mod') && file_exists('C:/Users/trcnologia/Desktop/stept-cli/cmd/stept/main_test.go')`
* **Spec Ref:** `AC-E03.1`

#### Ticket `TICK-03.2` — Instalador Universal de PATH (stept install)
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Modo `suggested` · Harness `claude` · Modelo `claude-sonnet-5-5` · Effort `medium` · Fallback Harness `9router` · Modelo `task-micro`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `300s` · Teto de saída `800` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `C:/Users/trcnologia/Desktop/stept-cli/cmd/stept/install.go` · Spec Anchor `docs/specs/2026-11-05_SPEC_STEPT_SWITCH.md#AC-E03.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `cmd/stept/install_test.go:TestStandaloneInstall`
* **Outputs:** Artifacts `C:/Users/trcnologia/Desktop/stept-cli/cmd/stept/install.go`, `C:/Users/trcnologia/Desktop/stept-cli/cmd/stept/install_test.go`
* **Verificação (argv, sem shell):** `cd C:/Users/trcnologia/Desktop/stept-cli`; `go test -count=1 ./...`; `go build -o bin/stept.exe ./cmd/stept`
* **Gate (DSL):** `exit_code == 0 && file_exists('C:/Users/trcnologia/Desktop/stept-cli/cmd/stept/install.go') && file_exists('C:/Users/trcnologia/Desktop/stept-cli/cmd/stept/install_test.go')`
* **Spec Ref:** `AC-E03.2`

#### Ticket `TICK-03.3` — Documentação Comercial, Manual e Prova de Latência Sub-10ms
* **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Modo `suggested` · Harness `claude` · Modelo `claude-haiku-5-5` · Effort `low` · Fallback Harness `9router` · Modelo `code-tests`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `C:/Users/trcnologia/Desktop/stept-cli/README.md`, `C:/Users/trcnologia/Desktop/stept-cli/MANUAL.md`, `C:/Users/trcnologia/Desktop/stept-cli/bin/stept.exe` · Spec Anchor `docs/specs/2026-11-05_SPEC_STEPT_SWITCH.md#AC-E03.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `cmd/stept/perf_test.go:TestStandalonePerformance`
* **Outputs:** Artifacts `C:/Users/trcnologia/Desktop/stept-cli/README.md`, `C:/Users/trcnologia/Desktop/stept-cli/MANUAL.md`, `C:/Users/trcnologia/Desktop/stept-cli/bin/stept.exe`
* **Verificação (argv, sem shell):** `cd C:/Users/trcnologia/Desktop/stept-cli`; `go test -count=1 ./...`; `go build -o bin/stept.exe ./cmd/stept`
* **Gate (DSL):** `exit_code == 0 && file_exists('C:/Users/trcnologia/Desktop/stept-cli/README.md')`
* **Spec Ref:** `AC-E03.3` (Também fecha: `AC-NF.1`)

---

### 🛑 STOP-GATE HUMANO — ÉPICO 03 (`SG-E03`)
* **Gate (DSL):** `prompt('Evidências do Épico 03 aprovadas? Autoriza iniciar o Épico 04?')`
- [ ] 100% dos tickets do Épico 03 executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários e de integração verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para iniciar o próximo Épico.

---

### Épico 04 — Testes de Mutação (Regra 8) & Homologação Final

#### Ticket `TICK-04.1` — Prova de Mutações Deliberadas nos 4 Pontos Críticos (M1 a M4)
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Modo `suggested` · Harness `claude` · Modelo `claude-sonnet-5-5` · Effort `medium` · Fallback Harness `9router` · Modelo `task-micro`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `300s` · Teto de saída `800` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `[]` · Spec Anchor `docs/specs/2026-11-05_SPEC_STEPT_SWITCH.md#AC-E04.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/ste/mutations_test.go:TestMutationsCoverage`
* **Outputs:** Artifacts `internal/ste/mutations_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 -v -run ^TestMutationsCoverage ./internal/ste/...`; `go vet ./internal/ste/...`; `go run ./tools/funclen internal/ste`; `go build -o bin/legoloop.exe ./cmd/legoloop`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/ste/mutations_test.go')`
* **Spec Ref:** `AC-E04.1` (Também fecha: `AC-E04.2`, `AC-E04.3`, `AC-E04.4`)

#### Ticket `TICK-04.2` — Homologação Final, Varredura Limpa e Relatório Consolidado
* **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Modo `suggested` · Harness `claude` · Modelo `claude-haiku-5-5` · Effort `low` · Fallback Harness `9router` · Modelo `code-tests`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `[]` · Spec Anchor `docs/specs/2026-11-05_SPEC_STEPT_SWITCH.md#AC-NF.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `docs/reports/2026-11-06_REPORT_STEPT_SWITCH.md`
* **Outputs:** Artifacts `docs/reports/2026-11-06_REPORT_STEPT_SWITCH.md`
* **Verificação (argv, sem shell):** `go test -count=1 ./...`; `go vet ./...`; `go run ./tools/funclen internal`; `go build -o bin/legoloop.exe ./cmd/legoloop`
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/reports/2026-11-06_REPORT_STEPT_SWITCH.md')`
* **Spec Ref:** `AC-NF.1`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 04 (`SG-E04`)
* **Gate (DSL):** `prompt('Evidências do Épico 04 e testes de mutação aprovados? Autoriza a homologação final?')`
- [ ] 100% dos tickets do Épico 04 executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários e de integração verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para iniciar o próximo Épico.

---

### 🛑 STOP-GATE HUMANO FINAL (`SG-FINAL`)
* **Nature:** `nano` | **Tipo:** `stop_gate` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["SG-E04"]`
* **Inputs:** `bin/legoloop.exe`, `internal/ste/switch.go`, `internal/cli/stept.go`, `C:/Users/trcnologia/Desktop/stept-cli/bin/stept.exe`
* **Outputs:** `docs/reports/2026-11-06_REPORT_STEPT_SWITCH.md`
* **Verificação (argv, sem shell):** `go test -count=1 ./...` e `go run ./tools/funclen internal`
* **Gate (DSL):** `prompt('Todos os 4 Épicos executados, binários compilados e provados por mutações. Aprova a entrega final?')`
- [ ] Suíte completa do LegoLoop verde (`go test -count=1 ./...`).
- [ ] 100% das 4 mutações deliberadas detectadas e comprovadas.
- [ ] Binário `bin/legoloop.exe` recompilado e atualizado.
- [ ] Projeto standalone `stept-cli` compilado, testado e com instalador no Desktop.
- [ ] Relatório final de homologação `docs/reports/2026-11-06_REPORT_STEPT_SWITCH.md` gerado.
- [ ] Aprovação explícita do usuário.

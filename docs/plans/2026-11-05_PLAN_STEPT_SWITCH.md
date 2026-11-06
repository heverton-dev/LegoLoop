# 🗺️ PLAN: Alternador Global e Granular de Output STE-PT para Harnesses

> **Status:** Proposta para Revisão  
> **Data:** 2026-11-05  
> **Versão:** 1.0.0  
> **Autor:** Antigravity (Pair Programming com Usuário)  
> **Referência da Intenção:** [docs/intents/2026-11-05_INTENT_STEPT_SWITCH.md](../intents/2026-11-05_INTENT_STEPT_SWITCH.md)  
> **Referência da Especificação:** [docs/specs/2026-11-05_SPEC_STEPT_SWITCH.md](../specs/2026-11-05_SPEC_STEPT_SWITCH.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-05_ROUTE_STEPT_SWITCH.json](../routes/2026-11-05_ROUTE_STEPT_SWITCH.json)  
> **Template:** TEMPLATE_PLAN.md v2.0.0 (espelha TEMPLATE_ROUTE.json v2.0.0)  

---

## 1. Estratégia de Execução & Governança de Loteamento

* **Engine:** `orchestrix_hermetic` (execução isolada por bloco, sem context bloat)
* **Nível de Supervisão:** `micro`
* **Permissão de Subagentes:** `false` (AGENTS.md, regra 4: agentes de harness proibidos; subprocessos estéreis)
* **TDD Estrito:** `true` (Red ➔ Green ➔ Refactor obrigatório antes de qualquer código de produção)
* **Limite de Workers Paralelos:** `2`
* **Timeout Padrão por Bloco:** `180s`
* **Política de Rollback:** global `abort_pipeline`; por nó `stash_changes`

---

## 2. Mapa Geral de Épicos & Rastreabilidade de Critérios

| Épico ID | Nome do Épico | Total Tickets | Critérios de Aceite Cobertos | Natureza Dominante |
|---|---|---|---|---|
| **E01** | Motor de Configuração, Persistência e Regras STE-PT | 3 | `AC-E01.1`, `AC-E01.2`, `AC-E01.3` | `micro` |
| **E02** | Sincronizador de Harnesses & CLI Integrada | 2 | `AC-E02.1`, `AC-E02.2`, `AC-E02.3`, `AC-NF.1` | `micro` |
| **E03** | Projeto Standalone no Desktop (`stept-cli`) | 3 | `AC-E03.1`, `AC-E03.2`, `AC-NF.1` | `micro` |

> Regra de cobertura: a união da coluna "Critérios de Aceite Cobertos" é estritamente igual ao conjunto de `AC-*` da SPEC.

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Motor de Configuração, Persistência e Regras STE-PT

#### Ticket `TICK-01.1` — Modelos de Dados, Enums e Validação
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:**
  * **Workspace:** `local`
  * **Executor:** `orchestrix_hermetic`
  * **Subagentes Permitidos:** `false`
* **Roteamento Cognitivo:**
  * **Modo:** `suggested`
  * **Primário:** Harness `gemini` · Modelo `flash` · Effort `low`
  * **Fallback:** Harness `antigravity` · Modelo `flash`
* **Tooling:**
  * **Runtime:** `powershell`
  * **MCPs Requeridos:** `false` · **Allowlist:** `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:**
  * **Source Files:** `internal/ste/switch.go`
  * **Spec Anchor:** `docs/specs/2026-11-05_SPEC_STEPT_SWITCH.md#AC-E01.2`
  * **Context Artifacts:** `[]`
* **Alvo de Teste TDD:** `internal/ste/switch_test.go:TestConfigPerfil_Validate`
* **Outputs:**
  * **Artifacts Produzidos:** `internal/ste/switch.go`, `internal/ste/switch_test.go`
  * **Verificação (argv, sem shell):** `go test -v ./internal/ste/... -run TestConfigPerfil_Validate`
  * **Gate (DSL):** `exit_code == 0 && file_exists('internal/ste/switch.go') && file_exists('internal/ste/switch_test.go')`
* **Spec Ref:** `AC-E01.2`

#### Ticket `TICK-01.2` — Persistência JSON e Resolução de Herança
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Harness `gemini` · Modelo `flash` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/ste/switch.go` · Spec Anchor `docs/specs/2026-11-05_SPEC_STEPT_SWITCH.md#AC-E01.1`
* **Alvo de Teste TDD:** `internal/ste/switch_test.go:TestConfig_ResolveEHeranca`
* **Outputs:** `internal/ste/switch.go`, `internal/ste/switch_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/ste/... -run TestConfig_ResolveEHeranca`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E01.1`

#### Ticket `TICK-01.3` — Geração do Prompt Canônico ASD / STE-PT
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Harness `gemini` · Modelo `flash` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/ste/prompt.go` · Spec Anchor `docs/specs/2026-11-05_SPEC_STEPT_SWITCH.md#AC-E01.3`
* **Alvo de Teste TDD:** `internal/ste/prompt_test.go:TestPromptRegra`
* **Outputs:** `internal/ste/prompt.go`, `internal/ste/prompt_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/ste/... -run TestPromptRegra`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/ste/prompt.go')`
* **Spec Ref:** `AC-E01.3`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 01 (`SG-E01`)
* **Gate (DSL):** `prompt('Evidências do Épico 01 aprovadas? Autoriza iniciar o Épico 02?')`
- [ ] 100% dos tickets de E01 executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para iniciar o próximo Épico.

---

### Épico 02 — Sincronizador de Harnesses & CLI Integrada

#### Ticket `TICK-02.1` — Sincronizador de Harnesses (Injeção e Remoção Atômica)
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Harness `gemini` · Modelo `flash` · Effort `medium` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/ste/sync.go` · Spec Anchor `docs/specs/2026-11-05_SPEC_STEPT_SWITCH.md#AC-E02.1`
* **Alvo de Teste TDD:** `internal/ste/sync_test.go:TestSyncHarnesses`
* **Outputs:** `internal/ste/sync.go`, `internal/ste/sync_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/ste/... -run TestSyncHarnesses`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/ste/sync.go')`
* **Spec Ref:** `AC-E02.1`

#### Ticket `TICK-02.2` — Subcomando CLI `legoloop stept` e Tabela de Status
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Harness `gemini` · Modelo `flash` · Effort `medium` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/cli/stept.go`, `cmd/legoloop/main.go` · Spec Anchor `docs/specs/2026-11-05_SPEC_STEPT_SWITCH.md#AC-E02.3`
* **Alvo de Teste TDD:** `internal/cli/stept_test.go:TestCLIStept_StatusESet`
* **Outputs:** `internal/cli/stept.go`, `internal/cli/stept_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/cli/... -run TestCLIStept`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/stept.go')`
* **Spec Ref:** `AC-E02.3`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 02 (`SG-E02`)
* **Gate (DSL):** `prompt('Evidências do Épico 02 aprovadas? Autoriza iniciar o Épico 03?')`
- [ ] 100% dos tickets de E02 executados com sucesso (PASS nos gates mecânicos).
- [ ] Subcomando `legoloop stept status` provado em execução hermética.
- [ ] Autorização explícita do usuário concedida para iniciar o próximo Épico.

---

### Épico 03 — Projeto Standalone no Desktop (`stept-cli`)

#### Ticket `TICK-03.1` — Estrutura e Código do Projeto Standalone
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Harness `gemini` · Modelo `flash` · Effort `medium` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `C:\Users\trcnologia\Desktop\stept-cli\` · Spec Anchor `docs/specs/2026-11-05_SPEC_STEPT_SWITCH.md#AC-E03.1`
* **Alvo de Teste TDD:** `cmd/stept/main_test.go:TestStandaloneMain`
* **Outputs:** `C:\Users\trcnologia\Desktop\stept-cli\go.mod`, `C:\Users\trcnologia\Desktop\stept-cli\cmd\stept\main.go`
* **Verificação (argv, sem shell):** `go test -v ./...`
* **Gate (DSL):** `exit_code == 0 && dir_exists('C:/Users/trcnologia/Desktop/stept-cli')`
* **Spec Ref:** `AC-E03.1`

#### Ticket `TICK-03.2` — Documentação Comercial e Manual de Referência
* **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Harness `gemini` · Modelo `flash` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Spec Anchor `docs/specs/2026-11-05_SPEC_STEPT_SWITCH.md#AC-E03.2`
* **Outputs:** `C:\Users\trcnologia\Desktop\stept-cli\README.md`, `C:\Users\trcnologia\Desktop\stept-cli\MANUAL.md`
* **Verificação (argv, sem shell):** `go test -v ./...`
* **Gate (DSL):** `file_exists('C:/Users/trcnologia/Desktop/stept-cli/README.md') && file_exists('C:/Users/trcnologia/Desktop/stept-cli/MANUAL.md')`
* **Spec Ref:** `AC-E03.2`

#### Ticket `TICK-03.3` — Compilação e Prova de Latência Sub-10ms
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Harness `gemini` · Modelo `flash` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `C:\Users\trcnologia\Desktop\stept-cli` · Spec Anchor `docs/specs/2026-11-05_SPEC_STEPT_SWITCH.md#AC-NF.1`
* **Outputs:** `C:\Users\trcnologia\Desktop\stept-cli\bin\stept.exe`
* **Verificação (argv, sem shell):** `go build -o bin/stept.exe ./cmd/stept`
* **Gate (DSL):** `exit_code == 0 && file_exists('C:/Users/trcnologia/Desktop/stept-cli/bin/stept.exe')`
* **Spec Ref:** `AC-NF.1`

---

### 🛑 STOP-GATE HUMANO FINAL (`SG-FINAL`)
* **Nature:** `macro` | **Tipo:** `stop_gate` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.3"]`
* **Inputs:** `bin/legoloop.exe`, `internal/ste/switch.go`, `internal/cli/stept.go`, `C:/Users/trcnologia/Desktop/stept-cli/bin/stept.exe`, `C:/Users/trcnologia/Desktop/stept-cli/README.md`, `C:/Users/trcnologia/Desktop/stept-cli/MANUAL.md`
* **Outputs:** `docs/reports/2026-11-05_REPORT_STEPT_SWITCH.md`
* **Verificação (argv, sem shell):** `go test -count=1 ./...` e `go run ./tools/funclen internal`
* **Gate (DSL):** `prompt('Todos os 3 Épicos executados, binários compilados e provados. Aprova a entrega final?')`
- [ ] Suíte completa do LegoLoop verde (`go test -count=1 ./...`).
- [ ] Binário `bin/legoloop.exe` recompilado e atualizado.
- [ ] Projeto standalone `stept-cli` compilado, testado e com manuais no Desktop.
- [ ] Relatório final de homologação `docs/reports/2026-11-05_REPORT_STEPT_SWITCH.md` gerado.
- [ ] Aprovação explícita do usuário.

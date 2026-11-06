# 🗺️ PLAN: Governança Lazy Universal de Ferramentas e Conhecimentos no LegoLoop

> **Status:** Aprovado pelo Usuário  
> **Data:** 2026-11-06  
> **Versão:** 1.0.0  
> **Autor:** Antigravity & Usuário  
> **Referência da Intenção:** [docs/intents/2026-11-06_INTENT_GOVERNANCA_LAZY_UNIVERSAL.md](../intents/2026-11-06_INTENT_GOVERNANCA_LAZY_UNIVERSAL.md)  
> **Referência da Especificação:** [docs/specs/2026-11-06_SPEC_GOVERNANCA_LAZY_UNIVERSAL.md](../specs/2026-11-06_SPEC_GOVERNANCA_LAZY_UNIVERSAL.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-06_ROUTE_GOVERNANCA_LAZY_UNIVERSAL.json](../routes/2026-11-06_ROUTE_GOVERNANCA_LAZY_UNIVERSAL.json)  
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
| **E01** | Meta-MCP Router STDIO (`internal/mcp`) | 3 | `AC-E01.1`, `AC-E01.2`, `AC-E01.3` | `micro` |
| **E02** | Lazy Skill Engine com Metadados Enriquecidos | 2 | `AC-E02.1`, `AC-E02.2` | `micro` |
| **E03** | Sincronizador de Harnesses & Recompilação | 2 | `AC-E03.1` | `micro` |

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Meta-MCP Router STDIO (`internal/mcp`)

#### Ticket `TICK-01.1` — Handshake JSON-RPC e Ferramentas Mínimas (2 tools)
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PASSED`
* **Dependências:** `[]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:**
  * **Source Files:** `internal/mcp/router.go`
  * **Spec Anchor:** `docs/specs/2026-11-06_SPEC_GOVERNANCA_LAZY_UNIVERSAL.md#AC-E01.1`
  * **Context Artifacts:** `[]`
* **Alvo de Teste TDD:** `internal/mcp/router_test.go:TestRouterHandshakeAndListTools`
* **Outputs:**
  * **Artifacts Produzidos:** `internal/mcp/router.go`, `internal/mcp/router_test.go`
  * **Verificação (argv, sem shell):** `go test -v ./internal/mcp/...`
  * **Gate (DSL):** `exit_code == 0 && file_exists('internal/mcp/router.go') && file_exists('internal/mcp/router_test.go')`
* **Spec Ref:** `AC-E01.1`
* **Escopo:** Prova que o router responde ao handshake inicial JSON-RPC 2.0 e que `tools/list` retorna estritamente `mcp_search` e `call_mcp_tool`.

#### Ticket `TICK-01.2` — Despacho Dinâmico de Ferramentas (`mcp_search` e `call_mcp_tool`)
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PASSED`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/mcp/router.go` · Spec Anchor `docs/specs/2026-11-06_SPEC_GOVERNANCA_LAZY_UNIVERSAL.md#AC-E01.2`
* **Outputs:** `internal/mcp/router.go`, `internal/mcp/router_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/mcp/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E01.2`
* **Escopo:** Prova que `mcp_search` busca e filtra ferramentas do cofre sem vazar ferramentas desconexas, e que chamadas inválidas retornam `-32601`.

#### Ticket `TICK-01.3` — Tratamento de Erros e Verificação Estática de E01
* **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PASSED`
* **Dependências:** `["TICK-01.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Harness `9router` · Modelo `code-nano` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** `internal/mcp/router.go` · Spec Anchor `docs/specs/2026-11-06_SPEC_GOVERNANCA_LAZY_UNIVERSAL.md#AC-E01.3`
* **Outputs:** `internal/mcp/router.go`
* **Verificação (argv, sem shell):** `go vet ./internal/mcp/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E01.3`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 01 (`SG-E01`)
* **Gate (DSL):** `prompt('Evidências do Épico 01 aprovadas? Autoriza iniciar o Épico 02?')`
- [x] 100% dos tickets de E01 executados com sucesso (PASS nos gates mecânicos).
- [x] Evidências de testes unitários verdes apresentadas ao usuário.
- [x] Autorização explícita do usuário concedida para iniciar o próximo Épico.

---

### Épico 02 — Lazy Skill Engine com Metadados Enriquecidos

#### Ticket `TICK-02.1` — Motor de Busca com Taxonomia (Tier, Natureza, Categoria)
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PASSED`
* **Dependências:** `["SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/mcp/skills.go` · Spec Anchor `docs/specs/2026-11-06_SPEC_GOVERNANCA_LAZY_UNIVERSAL.md#AC-E02.1`
* **Outputs:** `internal/mcp/skills.go`, `internal/mcp/skills_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/mcp/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/mcp/skills.go')`
* **Spec Ref:** `AC-E02.1`

#### Ticket `TICK-02.2` — Poda de Pastas Upstream com Prefixo de Ponto
* **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PASSED`
* **Dependências:** `["TICK-02.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Harness `9router` · Modelo `code-nano` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** `internal/mcp/skills.go` · Spec Anchor `docs/specs/2026-11-06_SPEC_GOVERNANCA_LAZY_UNIVERSAL.md#AC-E02.2`
* **Outputs:** `internal/mcp/skills.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/mcp/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E02.2`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 02 (`SG-E02`)
* **Gate (DSL):** `prompt('Evidências do Épico 02 aprovadas? Autoriza iniciar o Épico 03?')`
- [x] 100% dos tickets de E02 executados com sucesso (PASS nos gates mecânicos).
- [x] Evidências de testes unitários verdes apresentadas ao usuário.
- [x] Autorização explícita do usuário concedida para iniciar o próximo Épico.

---

### Épico 03 — Sincronizador de Harnesses & Recompilação

#### Ticket `TICK-03.1` — CLI `legoloop sync-harnesses` e Distribuição de Configs
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PASSED`
* **Dependências:** `["SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/cli/sync_harnesses.go` · Spec Anchor `docs/specs/2026-11-06_SPEC_GOVERNANCA_LAZY_UNIVERSAL.md#AC-E03.1`
* **Outputs:** `internal/cli/sync_harnesses.go`, `internal/cli/sync_harnesses_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/cli/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/sync_harnesses.go')`
* **Spec Ref:** `AC-E03.1`

#### Ticket `TICK-03.2` — Recompilação Fresca e Integridade da Suíte Global
* **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PASSED`
* **Dependências:** `["TICK-03.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Harness `9router` · Modelo `code-nano` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** `cmd/legoloop/main.go`
* **Outputs:** `bin/legoloop.exe`
* **Verificação (argv, sem shell):** `go build -o bin/legoloop.exe ./cmd/legoloop`
* **Gate (DSL):** `exit_code == 0 && file_exists('bin/legoloop.exe')`
* **Spec Ref:** `AC-E03.1`

---

### 🛑 STOP-GATE FINAL DE AUDITORIA HUMANA (`SG-FINAL`)
* **Gate (DSL):** `prompt('Todos os 3 Épicos concluídos e binário atualizado. Autoriza entrega final?')`
- [x] 100% dos tickets executados com sucesso (PASS nos gates mecânicos).
- [x] Suíte de testes `go test -count=1 ./...` 100% verde.
- [x] Verificação de funções enxutas (`funclen`) limpa.
- [x] Binário `bin/legoloop.exe` recompilado e testado.
- [x] Autorização explícita do usuário concedida e `legoloop sync-harnesses --apply` executado com backups criados.

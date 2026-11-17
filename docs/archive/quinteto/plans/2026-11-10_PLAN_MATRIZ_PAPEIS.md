---
id: PLAN-MATRIZ_PAPEIS
status: APPROVED
date: 2026-11-10
schema_version: 2
upstream: [SPEC-MATRIZ_PAPEIS]
covers: [AC-I5-01.1, AC-I5-01.2, AC-I5-02.1, AC-I5-02.2, AC-I5-02.3, AC-I5-03.1, AC-I5-03.2]
---

# 🗺️ PLAN: Matriz Canônica de Papéis e Casamento Determinístico de Peças

> **Status:** Aprovado pelo Usuário  
> **Data:** 2026-11-10  
> **Versão:** 2.2.0  
> **Autor:** LegoLoop  
> **Iniciativa:** Fábrica Determinística LegoLoop (I5)  
> **Referência da Intenção:** [docs/intents/2026-11-10_INTENT_MATRIZ_PAPEIS.md](../intents/2026-11-10_INTENT_MATRIZ_PAPEIS.md)  
> **Referência da Especificação:** [docs/specs/2026-11-10_SPEC_MATRIZ_PAPEIS.md](../specs/2026-11-10_SPEC_MATRIZ_PAPEIS.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-10_ROUTE_MATRIZ_PAPEIS.json](../routes/2026-11-10_ROUTE_MATRIZ_PAPEIS.json)  
> **Template:** TEMPLATE_PLAN.md v2.2.0 (espelha TEMPLATE_ROUTE.json v2.2.0)  

---

## 1. Estratégia de Execução & Governança de Loteamento

* **Engine:** `orchestrix_hermetic` (execução isolada por Ticket, sem context bloat)
* **Nível de Supervisão:** `micro`
* **Permissão de Subagentes:** `false`
* **TDD Estrito:** `true` (Red ➔ Green ➔ Refactor obrigatório antes de qualquer código de produção)
* **Limite de Workers Paralelos:** `2`
* **Timeout Padrão por Ticket:** `180s`
* **Política de Rollback:** global `abort_pipeline`; por nó `stash_changes`

---

## 2. Mapa Geral de Épicos & Rastreabilidade de Critérios

| Épico ID | Nome do Épico | Total Tickets | Critérios de Aceite Cobertos | Natureza Dominante |
|---|---|---|---|---|
| **E01** | Catálogo de Papéis e Validação Estrita (Taxonomia) | 2 | `AC-I5-01.1`, `AC-I5-01.2` | `micro` |
| **E02** | Contrato da Peça na Matriz e Filtros na CLI | 2 | `AC-I5-02.1`, `AC-I5-02.2`, `AC-I5-02.3` | `micro` |
| **E03** | Casamento Determinístico no Gerador (`internal/fabricar`) | 2 | `AC-I5-03.1`, `AC-I5-03.2` | `micro` |

| R | Requisito | AC | Ticket |
| :--- | :--- | :--- | :--- |
| R1 | Catálogo oficial dos 18 papéis canônicos em 7 guildas em `internal/taxonomy` | `AC-I5-01.1`, `AC-I5-01.2` | `TICK-01.1`, `TICK-01.2` |
| R2 | Contrato de peça na Matriz exige `papel_primario` e valida `papeis_secundarios` | `AC-I5-02.1` | `TICK-02.1` |
| R3 | Disjuntor estrito retorna erro sentinela para papel desconhecido | `AC-I5-01.2`, `AC-I5-02.1` | `TICK-01.1`, `TICK-02.1` |
| R4 | Filtros `--papel` e `--guilda` em `legoloop matriz ver` (`internal/cli/matriz_ver.go`) | `AC-I5-02.2`, `AC-I5-02.3` | `TICK-02.2` |
| R5 | Casamento determinístico em 2 níveis em `internal/fabricar/merger_fit.go` | `AC-I5-03.1`, `AC-I5-03.2` | `TICK-03.1`, `TICK-03.2` |

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Catálogo de Papéis e Validação Estrita (Taxonomia)

#### Ticket `TICK-01.1` — Implementar Catálogo dos 18 Papéis e Validador em Go
* **Iniciativa:** `Fábrica Determinística LegoLoop (I5)` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:**
  * **Workspace:** `local`
  * **Executor:** `orchestrix_hermetic`
  * **Subagentes Permitidos:** `false`
  * **Write Boundary:** `strict`
  * **Extreme:** `false`
* **Roteamento Cognitivo:**
  * **Modo:** `suggested`
  * **Primário:** Harness `claude` · Via `oficial` · Modelo `sonnet` · Effort `high`
  * **Fallback:** Harness `claude` · Via `oficial` · Modelo `haiku`
* **Tooling:**
  * **Runtime:** `powershell`
  * **MCPs Requeridos:** `false` · **Allowlist:** `[]`
* **Limites & Resiliência:**
  * **Timeout:** `180s` · **Retries:** `1`
  * **Orçamento de Tokens:** Entrada `16000` tokens · Teto de saída `600` tokens
  * **Rollback:** `stash_changes`
* **Inputs:**
  * **Source Files:** `internal/taxonomy/taxonomy.go`, `internal/taxonomy/roles.go`
  * **Spec Anchor:** `docs/specs/2026-11-10_SPEC_MATRIZ_PAPEIS.md#AC-I5-01.1`
  * **Context Artifacts:** `[]`
* **Alvo de Teste TDD:** `internal/taxonomy/roles_test.go:TestValidarPapel`
* **Outputs:**
  * **Artifacts Produzidos:** `internal/taxonomy/roles.go`, `internal/taxonomy/roles_test.go`
  * **Verificação (argv, sem shell):** `go test -count=1 ./internal/taxonomy/...`
  * **Gate (DSL):** `exit_code == 0 && file_exists('internal/taxonomy/roles.go') && file_exists('internal/taxonomy/roles_test.go')`
* **Spec Ref:** `AC-I5-01.1`, `AC-I5-01.2`

#### Ticket `TICK-01.2` — Sincronizar Dicionário JSON e Schemas Oficiais
* **Iniciativa:** `Fábrica Determinística LegoLoop (I5)` | **Épico:** `E01` | **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `mimo` · Via `oficial` · Modelo `xiaomi/mimo-v2` · Effort `low` · Fallback Harness `claude` · Via `oficial` · Modelo `haiku`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** `internal/taxonomy/dictionary.json`, `docs/vocab/2026-11-10_ASD_STE100_CANONICAL_DICTIONARY.json`, `schemas/dag.v2.json` · Spec Anchor `docs/specs/2026-11-10_SPEC_MATRIZ_PAPEIS.md#AC-I5-01.1`
* **Outputs:** `internal/taxonomy/dictionary.json`, `docs/vocab/2026-11-10_ASD_STE100_CANONICAL_DICTIONARY.json`, `schemas/dag.v2.json`
* **Verificação (argv, sem shell):** `go run ./tools/taxlinter internal cmd loops skills docs/templates`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-I5-01.1`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 01 (`SG-E01`)
* **Dependências:** `["TICK-01.2"]`
* **Gate (DSL):** `prompt('Evidências do Épico 01 aprovadas? Autoriza iniciar o Épico 02?')`
* **Roteamento:** Harness `human` · Via `oficial` · Modelo `human/operator` · Effort `low`
- [ ] 100% dos tickets de E01 executados com sucesso (PASS nos gates mecânicos).

---

### Épico 02 — Contrato da Peça na Matriz e Filtros na CLI

#### Ticket `TICK-02.1` — Atualizar Contrato da Peça e Validação na Matriz
* **Iniciativa:** `Fábrica Determinística LegoLoop (I5)` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `claude` · Via `oficial` · Modelo `sonnet` · Effort `high` · Fallback Harness `mimo` · Via `oficial` · Modelo `xiaomi/mimo-v2`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/matriz/perfil.go`, `internal/matriz/catalogador.go` · Spec Anchor `docs/specs/2026-11-10_SPEC_MATRIZ_PAPEIS.md#AC-I5-02.1`
* **Outputs:** `internal/matriz/perfil.go`, `internal/matriz/catalogador.go`, `internal/matriz/catalogador_test.go`
* **Alvo de Teste TDD:** `internal/matriz/catalogador_test.go:TestPecaValidateRoles`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/matriz/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-I5-02.1`

#### Ticket `TICK-02.2` — Adicionar Filtros `--papel` e `--guilda` na CLI da Matriz
* **Iniciativa:** `Fábrica Determinística LegoLoop (I5)` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `opencode` · Via `oficial` · Modelo `default` · Effort `high` · Fallback Harness `claude` · Via `oficial` · Modelo `sonnet`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/cli/matriz_ver.go`, `cmd/legoloop/main.go` · Spec Anchor `docs/specs/2026-11-10_SPEC_MATRIZ_PAPEIS.md#AC-I5-02.2`
* **Outputs:** `internal/cli/matriz_ver.go`, `internal/cli/matriz_ver_test.go`
* **Alvo de Teste TDD:** `internal/cli/matriz_ver_test.go:TestMatrizVerFiltrosPapelGuilda`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/cli/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-I5-02.2`, `AC-I5-02.3`

---

### Épico 03 — Casamento Determinístico no Gerador (`internal/fabricar`)

#### Ticket `TICK-03.1` — Implementar Algoritmo de Casamento em 2 Níveis no Fabricar
* **Iniciativa:** `Fábrica Determinística LegoLoop (I5)` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `claude` · Via `oficial` · Modelo `sonnet` · Effort `high` · Fallback Harness `claude` · Via `oficial` · Modelo `haiku`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/fabricar/merger_fit.go`, `internal/fabricar/engine.go` · Spec Anchor `docs/specs/2026-11-10_SPEC_MATRIZ_PAPEIS.md#AC-I5-03.1`
* **Outputs:** `internal/fabricar/merger_fit.go`, `internal/fabricar/merger_rank_test.go`
* **Alvo de Teste TDD:** `internal/fabricar/merger_rank_test.go:TestCasamentoPapeis`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/fabricar/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-I5-03.1`, `AC-I5-03.2`

---

### 🛑 STOP-GATE HUMANO — ÉPICOS 02 E 03 (`SG-E02`)
* **Dependências:** `["TICK-02.2", "TICK-03.1"]`
* **Gate (DSL):** `prompt('Evidências paralelas de CLI (TICK-02.2) e Fabricar (TICK-03.1) aprovadas? Autoriza recompilação final?')`
* **Roteamento:** Harness `human` · Via `oficial` · Modelo `human/operator` · Effort `low`
- [ ] 100% dos tickets de E02 e E03 (`TICK-02.1`, `TICK-02.2`, `TICK-03.1`) executados com sucesso (PASS nos gates mecânicos).

---

#### Ticket `TICK-03.2` — Verificação Integrada e Recompilação Fresca
* **Iniciativa:** `Fábrica Determinística LegoLoop (I5)` | **Épico:** `E03` | **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `claude` · Via `oficial` · Modelo `haiku` · Effort `low` · Fallback Harness `opencode` · Via `oficial` · Modelo `default`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** `cmd/legoloop/main.go` · Spec Anchor `docs/specs/2026-11-10_SPEC_MATRIZ_PAPEIS.md#AC-I5-03.2`
* **Outputs:** `bin/legoloop.exe`
* **Verificação (argv, sem shell):** `go build -o bin/legoloop.exe ./cmd/legoloop`
* **Gate (DSL):** `exit_code == 0 && file_exists('bin/legoloop.exe')`
* **Spec Ref:** `AC-I5-03.2`

---

### 🛑 STOP-GATE HUMANO FINAL (`SG-FINAL`)
* **Dependências:** `["TICK-03.2"]`
* **Gate (DSL):** `prompt('Todos os Épicos foram concluídos e validados? Autoriza o fechamento da demanda?')`
* **Roteamento:** Harness `human` · Via `oficial` · Modelo `human/operator` · Effort `low`
- [ ] 100% dos tickets executados com sucesso.
- [ ] Binário `bin/legoloop.exe` atualizado e testado.

---

## 4. Homologação Final
- [ ] `go test -count=1 ./...` verde sem avisos.
- [ ] `go run ./cmd/legoloop validate docs/dags/2026-11-10_DAG_MATRIZ_PAPEIS.json` verde.
- [ ] `go run ./tools/ticketgate -trace docs/intents/2026-11-10_INTENT_MATRIZ_PAPEIS.md,docs/specs/2026-11-10_SPEC_MATRIZ_PAPEIS.md,docs/plans/2026-11-10_PLAN_MATRIZ_PAPEIS.md,docs/routes/2026-11-10_ROUTE_MATRIZ_PAPEIS.json,docs/dags/2026-11-10_DAG_MATRIZ_PAPEIS.json` verde.
- [ ] `legoloop fechar` executado após aprovação humana em `SG-FINAL`.

---
id: PLAN-AUDITORIA-SWARM-01
status: DRAFT
date: 2026-11-14
schema_version: 2
upstream: [SPEC-AUDITORIA-SWARM-01]
covers: [AC-E01.1, AC-E02.1, AC-E03.1, AC-E04.1]
---

# 🗺️ PLAN: Auditoria Profunda de Código e UI/UX em Modo Swarm (4 Frentes)

> **Status:** Em Planejamento  
> **Data:** 2026-11-14  
> **Versão:** 1.0.0  
> **Autor:** Antigravity (Coordinator)  
> **Iniciativa:** Auditoria Profunda Swarm  
> **Referência da Intenção:** [docs/intents/2026-11-14_INTENT_AUDITORIA_PROFUNDA_SWARM.md](../intents/2026-11-14_INTENT_AUDITORIA_PROFUNDA_SWARM.md)  
> **Referência da Especificação:** [docs/specs/2026-11-14_SPEC_AUDITORIA_PROFUNDA_SWARM.md](../specs/2026-11-14_SPEC_AUDITORIA_PROFUNDA_SWARM.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-14_ROUTE_AUDITORIA_PROFUNDA_SWARM.json](../routes/2026-11-14_ROUTE_AUDITORIA_PROFUNDA_SWARM.json)  
> **Template:** TEMPLATE_PLAN.md v2.2.0  

---

## 1. Estratégia de Execução & Governança de Loteamento

* **Engine:** `orchestrix_hermetic` (execução isolada por Ticket em modo Swarm)
* **Nível de Supervisão:** `micro`
* **Permissão de Subagentes:** `ephemeral` (1 worker estéril por frente)
* **TDD Estrito:** `false` (missão de auditoria somente leitura)
* **Limite de Workers Paralelos:** `3`
* **Timeout Padrão por Ticket:** `180s`
* **Política de Rollback:** global `abort_pipeline`; por nó `stash_changes`

---

## 2. Mapa Geral de Épicos & Rastreabilidade de Critérios

| Épico ID | Nome do Épico | Total Tickets | Critérios de Aceite Cobertos | Natureza Dominante |
|---|---|---|---|---|
| **E01** | Frente 1 — Auditoria de Código Go | 1 | `AC-E01.1` | `micro` |
| **E02** | Frente 2 — Auditoria de UI/UX Studio | 1 | `AC-E02.1` | `micro` |
| **E03** | Frente 3 — Auditoria de Contratos e DAGs | 1 | `AC-E03.1` | `micro` |
| **E04** | Frente 4 — Consolidação Swarm | 1 | `AC-E04.1` | `micro` |

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Frente 1: Auditoria de Código Go (`TICK-01.1`)
* **Dependências:** `[]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low`
* **Verificação:** `go test -count=1 ./...` e `go run ./tools/funclen internal`
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/reports/audit_f1_go.json')`
* **Spec Ref:** `AC-E01.1`

### Épico 02 — Frente 2: Auditoria de UI/UX Studio (`TICK-02.1`)
* **Dependências:** `[]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `medium`
* **Verificação:** Inspeção de `internal/studio/` e registro de achados
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/reports/audit_f2_uiux.json')`
* **Spec Ref:** `AC-E02.1`

### Épico 03 — Frente 3: Auditoria de Contratos e DAGs (`TICK-03.1`)
* **Dependências:** `[]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low`
* **Verificação:** `legoloop validate` nos DAGs em `loops/`
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/reports/audit_f3_dags.json')`
* **Spec Ref:** `AC-E03.1`

### Épico 04 — Frente 4: Consolidação Swarm (`TICK-04.1`)
* **Dependências:** `["TICK-01.1", "TICK-02.1", "TICK-03.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high`
* **Verificação:** `Test-Path docs/reports/2026-11-14_REPORT_AUDITORIA_PROFUNDA_SWARM.md`
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/reports/2026-11-14_REPORT_AUDITORIA_PROFUNDA_SWARM.md')`
* **Spec Ref:** `AC-E04.1`

---

🛑 **STOP-GATE HUMANO (`SG-PLAN`):**  
*Autoriza o disparo dos 3 workers paralelos do Swarm conforme este plano?*

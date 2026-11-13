---
id: PLAN-I6
status: DRAFT
date: 2026-11-13
schema_version: 2
upstream: [SPEC-I6]
covers: [AC-I6-01.1, AC-I6-01.2, AC-I6-01.3, AC-I6-02.1, AC-I6-02.2, AC-I6-02.3, AC-I6-03.1]
---

# 🗺️ PLAN: Migração dos DAGs Legados (Iniciativa I6)

> **Status:** Em Planejamento  
> **Data:** 2026-11-13  
> **Versão:** 1.0.0  
> **Autor:** i6-preparer (Arquiteto I6)  
> **Iniciativa:** I6 — Migração dos Legados  
> **Referência da Intenção:** [docs/intents/2026-11-13_INTENT_MIGRACAO_LEGADOS.md](../intents/2026-11-13_INTENT_MIGRACAO_LEGADOS.md)  
> **Referência da Especificação:** [docs/specs/2026-11-13_SPEC_MIGRACAO_LEGADOS.md](../specs/2026-11-13_SPEC_MIGRACAO_LEGADOS.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-13_ROUTE_MIGRACAO_LEGADOS.json](../routes/2026-11-13_ROUTE_MIGRACAO_LEGADOS.json)  
> **Template:** TEMPLATE_PLAN.md v2.2.0  

---

## 1. Estratégia de Execução & Governança de Loteamento

* **Engine:** `orchestrix_hermetic` (execução isolada por Ticket, sem context bloat)
* **Nível de Supervisão:** `micro`
* **Permissão de Subagentes:** `false` (camadas paralelas coordenadas pelo motor DAG com `write_boundary: strict`)
* **TDD Estrito:** `false` (alteração exclusiva de manifestos JSON validados por `legoloop validate`)
* **Limite de Workers Paralelos:** `3`
* **Timeout Padrão por Ticket:** `180s`
* **Política de Rollback:** global `abort_pipeline`; por nó `stash_changes`

---

## 2. Mapa Geral de Épicos & Rastreabilidade de Critérios

| Épico ID | Nome do Épico | Total Tickets | Critérios de Aceite Cobertos | Natureza Dominante |
|---|---|---|---|---|
| **E01** | Inventário e Migração de `docs/dags/` | 3 | `AC-I6-01.1`, `AC-I6-01.2`, `AC-I6-01.3` | `micro` |
| **E02** | Migração Paralela de `loops/*/pipeline.dag.json` | 3 | `AC-I6-02.1`, `AC-I6-02.2`, `AC-I6-02.3` | `micro` |
| **E03** | Validação Global e Homologação | 1 | `AC-I6-03.1` | `nano` |

| R | Requisito | AC | Ticket |
| :--- | :--- | :--- | :--- |
| R1 | Inventário e migração dos 9 DAGs de `docs/dags/` | `AC-I6-01.1`, `AC-I6-01.2`, `AC-I6-01.3` | `TICK-I6-01.1`, `TICK-I6-01.2`, `TICK-I6-01.3` |
| R2 | Migração paralela dos 10 DAGs de `loops/*/` | `AC-I6-02.1`, `AC-I6-02.2`, `AC-I6-02.3` | `TICK-I6-02.1`, `TICK-I6-02.2`, `TICK-I6-02.3` |
| R3 | Validação mecânica integral dos 19 DAGs | `AC-I6-03.1` | `TICK-I6-03.1` |

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Inventário e Migração de `docs/dags/`

#### Ticket `TICK-I6-01.1` — Inventário determinístico dos 19 DAGs legados
* **Iniciativa:** `I6` | **Épico:** `E01` | **Nature:** `nano` | **Tipo:** `script` | **Modo de Execução:** `isolated` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `claude` · Via `oficial` · Modelo `haiku` · Effort `low`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `120s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** `docs/dags/2026-11-07_DAG_LEGOLOOP_FECHAR.json`, `loops/legoloop-fechar/pipeline.dag.json` · Spec Anchor `docs/specs/2026-11-13_SPEC_MIGRACAO_LEGADOS.md#AC-I6-01.1`
* **Outputs:** `scratch/i6_inventario_legados.json`
* **Verificação:** `go run ./cmd/legoloop validate docs/dags/2026-11-13_DAG_MIGRACAO_LEGADOS.json`
* **Gate (DSL):** `exit_code == 0 && file_exists('scratch/i6_inventario_legados.json')`
* **Spec Ref:** `AC-I6-01.1`

#### Ticket `TICK-I6-01.2` — Migrar Lote A de `docs/dags/` (4 DAGs operacionais)
* **Iniciativa:** `I6` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `parallel` | **Status:** `PENDING`
* **Dependências:** `["TICK-I6-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `claude` · Via `oficial` · Modelo `sonnet` · Effort `medium`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `scratch/i6_inventario_legados.json` · Spec Anchor `docs/specs/2026-11-13_SPEC_MIGRACAO_LEGADOS.md#AC-I6-01.2`
* **Outputs:** `docs/dags/2026-11-07_DAG_LEGOLOOP_FECHAR.json`, `docs/dags/2026-11-07_DAG_LEGOLOOP_VIDEO_MOTION.json`, `docs/dags/2026-11-08_DAG_LEGOLOOP_VIDEO_EDITION.json`, `docs/dags/2026-11-09_DAG_LEGOLOOP_VIDEO_REA.json`
* **Verificação:** `go run ./cmd/legoloop validate docs/dags/2026-11-07_DAG_LEGOLOOP_FECHAR.json`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-I6-01.2`

#### Ticket `TICK-I6-01.3` — Migrar Lote B de `docs/dags/` (5 DAGs de governança e I1-I3)
* **Iniciativa:** `I6` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `parallel` | **Status:** `PENDING`
* **Dependências:** `["TICK-I6-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `claude` · Via `oficial` · Modelo `sonnet` · Effort `medium`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `24000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `scratch/i6_inventario_legados.json` · Spec Anchor `docs/specs/2026-11-13_SPEC_MIGRACAO_LEGADOS.md#AC-I6-01.3`
* **Outputs:** `docs/dags/2026-11-09_DAG_QUINTETO_GOVERNANCA_ISO.json`, `docs/dags/2026-11-10_DAG_DISSECAR.json`, `docs/dags/2026-11-10_DAG_NORMALIZACAO_TAXONOMICA_LEI15.json`, `docs/dags/2026-11-12_DAG_ROTEAMENTO_HONESTO.json`, `docs/dags/2026-11-12_DAG_SEXTETO_E_TEMPLATES.json`
* **Verificação:** `go run ./cmd/legoloop validate docs/dags/2026-11-12_DAG_ROTEAMENTO_HONESTO.json`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-I6-01.3`

### 🛑 STOP-GATE HUMANO — ÉPICO 01 (`SG-I6-E01`)
* **Gate (DSL):** `prompt('Os 9 DAGs de docs/dags/ validam em v2? Autoriza iniciar E02 (loops/)?')`

---

### Épico 02 — Migração Paralela de `loops/*/pipeline.dag.json`

#### Ticket `TICK-I6-02.1` — Migrar Lote Vídeo I em `loops/` (4 arquivos)
* **Iniciativa:** `I6` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `parallel` | **Status:** `PENDING`
* **Dependências:** `["SG-I6-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `claude` · Via `oficial` · Modelo `sonnet` · Effort `medium`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `scratch/i6_inventario_legados.json` · Spec Anchor `docs/specs/2026-11-13_SPEC_MIGRACAO_LEGADOS.md#AC-I6-02.1`
* **Outputs:** `loops/ai-motion-graphics/pipeline.dag.json`, `loops/legoloop-video-edition/pipeline.dag.json`, `loops/legoloop-video-master/pipeline.dag.json`, `loops/legoloop-video-motion/pipeline.dag.json`
* **Verificação:** `go run ./cmd/legoloop validate loops/ai-motion-graphics/pipeline.dag.json`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-I6-02.1`

#### Ticket `TICK-I6-02.2` — Migrar Lote Vídeo II em `loops/` (3 arquivos)
* **Iniciativa:** `I6` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `parallel` | **Status:** `PENDING`
* **Dependências:** `["SG-I6-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `claude` · Via `oficial` · Modelo `sonnet` · Effort `medium`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `scratch/i6_inventario_legados.json` · Spec Anchor `docs/specs/2026-11-13_SPEC_MIGRACAO_LEGADOS.md#AC-I6-02.2`
* **Outputs:** `loops/legoloop-canal-redc/pipeline.dag.json`, `loops/legoloop-cartoon-shorts/pipeline.dag.json`, `loops/legoloop-video-rea/pipeline.dag.json`
* **Verificação:** `go run ./cmd/legoloop validate loops/legoloop-canal-redc/pipeline.dag.json`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-I6-02.2`

#### Ticket `TICK-I6-02.3` — Migrar Lote Governança em `loops/` (3 arquivos)
* **Iniciativa:** `I6` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `parallel` | **Status:** `PENDING`
* **Dependências:** `["SG-I6-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `claude` · Via `oficial` · Modelo `sonnet` · Effort `medium`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `scratch/i6_inventario_legados.json` · Spec Anchor `docs/specs/2026-11-13_SPEC_MIGRACAO_LEGADOS.md#AC-I6-02.3`
* **Outputs:** `loops/legoloop-dissecar/pipeline.dag.json`, `loops/legoloop-fechar/pipeline.dag.json`, `loops/legoloop-n10-foundry/pipeline.dag.json`
* **Verificação:** `go run ./cmd/legoloop validate loops/legoloop-fechar/pipeline.dag.json`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-I6-02.3`

### 🛑 STOP-GATE HUMANO — ÉPICO 02 (`SG-I6-E02`)
* **Gate (DSL):** `prompt('Os 10 pipeline.dag.json de loops/ validam em v2? Autoriza validação final E03?')`

---

### Épico 03 — Validação Global e Homologação

#### Ticket `TICK-I6-03.1` — Auditoria mecânica integral dos 19 DAGs migrados
* **Iniciativa:** `I6` | **Épico:** `E03` | **Nature:** `nano` | **Tipo:** `script` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["SG-I6-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `claude` · Via `oficial` · Modelo `haiku` · Effort `low`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** `scratch/i6_inventario_legados.json` · Spec Anchor `docs/specs/2026-11-13_SPEC_MIGRACAO_LEGADOS.md#AC-I6-03.1`
* **Outputs:** `scratch/i6_validacao_final.json`
* **Verificação:** `go test -count=1 ./internal/dag/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('scratch/i6_validacao_final.json')`
* **Spec Ref:** `AC-I6-03.1`

---

## 4. Homologação Final & Fechamento da Demanda (`SG-I6-FINAL`)

1. Suíte de regressão: `go test -count=1 ./internal/dag/...`
2. Linter taxonômico: `go run ./tools/taxlinter docs/intents/2026-11-13_INTENT_MIGRACAO_LEGADOS.md docs/specs/2026-11-13_SPEC_MIGRACAO_LEGADOS.md docs/plans/2026-11-13_PLAN_MIGRACAO_LEGADOS.md`
3. Rastreabilidade do Quinteto: `go run ./tools/ticketgate -trace docs/intents/2026-11-13_INTENT_MIGRACAO_LEGADOS.md,docs/specs/2026-11-13_SPEC_MIGRACAO_LEGADOS.md,docs/plans/2026-11-13_PLAN_MIGRACAO_LEGADOS.md,docs/routes/2026-11-13_ROUTE_MIGRACAO_LEGADOS.json,docs/dags/2026-11-13_DAG_MIGRACAO_LEGADOS.json`

* **Gate (DSL):** `exit_code == 0`

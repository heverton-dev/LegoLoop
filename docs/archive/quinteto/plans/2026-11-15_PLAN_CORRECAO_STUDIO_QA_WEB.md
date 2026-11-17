---
id: PLAN-CORRECAO_STUDIO_QA_WEB
status: DRAFT
date: 2026-11-15
schema_version: 2
upstream: [SPEC-CORRECAO_STUDIO_QA_WEB]
covers: [AC-E01.1, AC-E01.2, AC-E02.1, AC-E02.2, AC-E03.1, AC-E03.2]
---

# 🗺️ PLAN: Correção dos Achados QA Web do LegoLoop Studio v2

> **Status:** Pronto para Aprovação de Execução (SWARM -> BUILD)  
> **Data:** 2026-11-15  
> **Versão:** 1.0.0  
> **Autor:** Antigravity (LegoLoop Governança)  
> **Iniciativa:** CORRECAO_STUDIO_QA_WEB  
> **Referência da Intenção:** [2026-11-15_INTENT_CORRIGIR_OS_ACHADOS_DO_QA_WEB_.md](../intents/2026-11-15_INTENT_CORRIGIR_OS_ACHADOS_DO_QA_WEB_.md)  
> **Referência da Especificação:** [2026-11-15_SPEC_CORRECAO_STUDIO_QA_WEB.md](../specs/2026-11-15_SPEC_CORRECAO_STUDIO_QA_WEB.md)  
> **Espelho Executável JSON:** [2026-11-15_ROUTE_CORRECAO_STUDIO_QA_WEB.json](../routes/2026-11-15_ROUTE_CORRECAO_STUDIO_QA_WEB.json)  
> **Template:** TEMPLATE_PLAN.md v2.2.0

---

## 1. Estratégia de Execução & Governança de Loteamento

* **Engine:** `native_swarm` + `orchestrix_hermetic`
* **Nível de Supervisão:** `micro`
* **Permissão de Subagentes:** `ephemeral` (workers nativos Swarm em escopos de arquivo disjuntos)
* **TDD Estrito:** `true` (Red ➔ Green ➔ Refactor obrigatório)
* **Limite de Workers Paralelos:** `3`
* **Timeout Padrão por Ticket:** `180s`
* **Política de Rollback:** global `abort_processo`; por nó `stash_changes`

---

## 2. Mapa Geral de Épicos & Rastreabilidade de Critérios

| Épico ID | Nome do Épico | Total Tickets | Critérios de Aceite Cobertos | Natureza Dominante |
|---|---|---|---|---|
| **E01** | Concorrência Thread-Safe no Studio | 1 | `AC-E01.1`, `AC-E01.2` | `micro` |
| **E02** | Acessibilidade ARIA e Layout RenderLint | 1 | `AC-E02.1`, `AC-E02.2` | `micro` |
| **E03** | Extração SSOT de Design Tokens | 1 | `AC-E03.1`, `AC-E03.2` | `micro` |
| **E04** | Build e Reauditoria QA Web Real | 1 | `AC-E01.1`..`AC-E03.2` | `nano` |

| R | Requisito | AC | Ticket |
| :--- | :--- | :--- | :--- |
| `R1` | Eliminar `concurrent map writes` em `assets_api.go` | `AC-E01.1`, `AC-E01.2` | `TICK-01.1` |
| `R2` | Corrigir `aria-required-children` e violações `axe-core` | `AC-E02.1` | `TICK-02.1` |
| `R3` | Corrigir falhas de layout e controles (`renderlint`) | `AC-E02.2` | `TICK-02.1` |
| `R4` | Migrar cores hardcoded para `tokens.css` (`0` no linter) | `AC-E03.1`, `AC-E03.2` | `TICK-03.1`, `TICK-04.1` |

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Concorrência Thread-Safe no Studio

#### Ticket `TICK-01.1` — Clonar `it.Meta` em `present` e validar acesso concorrente em `internal/studio/`
* **Iniciativa:** `CORRECAO_STUDIO_QA_WEB` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:** `workspace: inherit`, `executor: native_swarm`, `write_boundary: strict` (`internal/studio/assets_api.go`, `internal/studio/assets_api_test.go`)
* **Spec Ref:** `AC-E01.1`, `AC-E01.2`
* **Gate:** `go test -count=1 ./internal/studio/...` (`exit_code == 0`)

### Épico 02 — Acessibilidade ARIA e Layout RenderLint

#### Ticket `TICK-02.1` — Corrigir `aria-required-children` e violações de layout em HTML/JS do Studio
* **Iniciativa:** `CORRECAO_STUDIO_QA_WEB` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:** `workspace: inherit`, `executor: native_swarm`, `write_boundary: strict` (`internal/studio/static/index.html`, `internal/studio/static/inventario.html`, `internal/studio/static/swarm.html`, `internal/studio/static/palette.js`, `internal/studio/static/app.js`, `internal/studio/static_test.go`)
* **Spec Ref:** `AC-E02.1`, `AC-E02.2`
* **Gate:** `go test -count=1 ./internal/studio/...` (`exit_code == 0`)

### Épico 03 — Extração SSOT de Design Tokens

#### Ticket `TICK-03.1` — Extrair `:root`/dark para `tokens.css` e eliminar 328 cores hardcoded em `style.css` e `swarm.css`
* **Iniciativa:** `CORRECAO_STUDIO_QA_WEB` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.1"]`
* **Isolamento:** `workspace: inherit`, `executor: native_swarm`, `write_boundary: strict` (`internal/studio/static/tokens.css`, `internal/studio/static/style.css`, `internal/studio/static/swarm.css`, HTMLs `<link rel="stylesheet" href="/tokens.css">`)
* **Spec Ref:** `AC-E03.1`, `AC-E03.2`
* **Gate:** `runDesignLinterAudit` retorna `violations == 0`

---

## 4. Homologação Final

#### Ticket `TICK-04.1` — Recompilar `bin/legoloop.exe`, rodar portões e executar `./bin/legoloop.exe web-qa run --project . --port 7895 --alvo legoloop-studio-v3`
* **Iniciativa:** `CORRECAO_STUDIO_QA_WEB` | **Épico:** `E04` | **Nature:** `nano` | **Tipo:** `script` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1", "TICK-02.1", "TICK-03.1"]`
* **Spec Ref:** `AC-E01.1`, `AC-E02.1`, `AC-E02.2`, `AC-E03.1`, `AC-E03.2`
* **Gate:** `go vet ./...`, `funclen`, `gofmt`, `go test -count=1 ./...` e `RELATORIO_QA_WEB.json` com `status == "PASSED"` (`exit_code == 0`)

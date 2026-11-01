# 🚀 LegoLoop — Plano de Execução & Micro/Nano Tickets: Super-Board Integration

> **Documento:** `docs/plans/2026-11-01_PLAN_SUPER_BOARD_INTEGRATION.md`  
> **Status:** Pronto para Execução  
> **Metodologia:** Tracer Bullet Development + TDD Estrito (Red-Green-Refactor) + Zero Stubs  
> **Spec de Referência:** `docs/specs/2026-11-01_SPEC_SUPER_BOARD_INTEGRATION.md`  
> **Intent de Referência:** `docs/intents/2026-11-01_INTENT_SUPER_BOARD_INTEGRATION.md`  

---

## 🗺️ Mapa de Épicos

| Épico | Título | Objetivo Central |
|---|---|---|
| **SB-01** | **GitHub Board Ingestor & Rate Guard** | Tipos de cards, cliente GitHub GraphQL com *token bucket* defensivo (Exit 79). |
| **SB-02** | **Wave Planner & File Overlap Engine** | Detector de arquivos tocados e particionamento de cards em ondas mutuamente disjuntas. |
| **SB-03** | **Issue Lifecycle Worktree & PR Runner** | Extensão de worktrees em Go para branch dedicada, draft PR e merge gate. |
| **SB-04** | **Guard Hooks & Canonização no Studio** | Transplante dos 4 guard hooks e inclusão do bloco de trigger no canvas N8N. |

---

## 📋 Decomposição em Micro e Nano Tickets

### SB-01: GitHub Board Ingestor & Rate Guard
- [x] **TICK-SB-101 (Nano):** Criar `internal/board/types.go` com structs `Card`, `Wave`, `BoardState` e interfaces do client.
- [x] **TICK-SB-102 (Nano):** Implementar `internal/board/guard.go` com verificação de quota GraphQL e teste unitário com mocks (`internal/board/guard_test.go`).
- [x] **TICK-SB-103 (Micro):** Implementar `internal/board/client.go` para consulta de colunas do GitHub Projects e transição de status.

### SB-02: Wave Planner & File Overlap Engine
- [x] **TICK-SB-201 (Nano):** Criar `internal/wave/overlap.go` para extração de caminhos de arquivos do markdown de issues (`internal/wave/overlap_test.go`).
- [x] **TICK-SB-202 (Micro):** Implementar `internal/wave/planner.go` para coloração/partição em ondas e testes de grafos com conflito vs disjuntos.

### SB-03: Issue Lifecycle Worktree & PR Runner
- [x] **TICK-SB-301 (Micro):** Implementar `internal/runner/lifecycle.go` para provisionamento de branch dedicada e testes de criação/cleanup (`internal/runner/lifecycle_test.go`).
- [x] **TICK-SB-302 (Nano):** Integrar validação pré-merge com o motor de gates do LegoLoop (`internal/gate`).

### SB-04: Guard Hooks & Canonização no Studio
- [x] **TICK-SB-401 (Nano):** Portar os 4 guard hooks para `.aidd/canonical/hooks/` com testes de validação.
- [x] **TICK-SB-402 (Nano):** Atualizar `internal/assets/hooks.go` para indexação profunda no scanner de ativos.
- [x] **TICK-SB-403 (Micro):** Registrar a peça `Trigger: GitHub Project Kanban` em `internal/studio/pieces.go` e expor na paleta do canvas N8N.

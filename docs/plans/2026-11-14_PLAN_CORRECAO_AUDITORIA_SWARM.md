---
id: PLAN-CORRECAO-AUDITORIA-01
status: DRAFT
date: 2026-11-14
schema_version: 2
upstream: [SPEC-CORRECAO-AUDITORIA-01]
covers: [AC-E01.1, AC-E02.1, AC-E03.1]
---

# 🗺️ PLAN: Ciclo de Correção Completa da Auditoria Profunda

> **Status:** Em Planejamento  
> **Data:** 2026-11-14  
> **Versão:** 1.0.0  
> **Template:** TEMPLATE_PLAN.md v2.2.0  

---

## 1. Decomposição Estruturada de Tickets

### Ticket `TICK-01.1` — Remover emojis de `lego-item-modal.js` (`AUD-03`)
* **Dependências:** `[]`
* **Spec Ref:** `AC-E01.1`
* **Gate:** `exit_code == 0`

### Ticket `TICK-02.1` — Adicionar breakpoint `<= 768px` em `style.css` (`AUD-04`)
* **Dependências:** `[]`
* **Spec Ref:** `AC-E02.1`
* **Gate:** `exit_code == 0`

### Ticket `TICK-03.1` — Padronizar `schema_version` nos `.loop.json` (`AUD-05`)
* **Dependências:** `[]`
* **Spec Ref:** `AC-E03.1`
* **Gate:** `exit_code == 0` (`go test -count=1 ./internal/dag/...`)

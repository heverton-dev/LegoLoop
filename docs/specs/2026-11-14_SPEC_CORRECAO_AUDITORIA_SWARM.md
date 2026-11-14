---
id: SPEC-CORRECAO-AUDITORIA-01
status: DRAFT
date: 2026-11-14
schema_version: 2
upstream: [INTENT-CORRECAO-AUDITORIA-01]
covers: [R1, R2, R3]
---

# 📋 SPEC: Ciclo de Correção Completa da Auditoria Profunda

> **Status:** Proposta para Revisão  
> **Data:** 2026-11-14  
> **Versão:** 1.0.0  
> **Autor:** Antigravity (Coordinator)  
> **Alvo Primário:** `internal/studio/static/`, `loops/`  
> **Template:** TEMPLATE_SPEC.md v2.0.0  

---

## 1. Critérios de Aceitação Binários (EARS)

| ID | Requisito | Critério Binário (EARS) |
| :--- | :--- | :--- |
| `AC-E01.1` | `R1` (`AUD-03`) | QUANDO `lego-item-modal.js` for inspecionado, ENTÃO zero emojis Unicode devem existir na interface, substituídos por ícones SVG de `icons.js`. |
| `AC-E02.1` | `R2` (`AUD-04`) | QUANDO `style.css` for inspecionado, ENTÃO deve existir `@media (max-width: 768px)` adaptando a sidebar e o grid principal. |
| `AC-E03.1` | `R3` (`AUD-05`) | QUANDO `go test -count=1 ./internal/dag/...` e `go test -count=1 ./...` forem executados após padronizar `schema_version` nos `.loop.json`, ENTÃO `exit_code == 0`. |

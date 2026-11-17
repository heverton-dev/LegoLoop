---
id: INTENT-CORRECAO-AUDITORIA-01
status: DRAFT
date: 2026-11-14
schema_version: 2
upstream: [REPORT-AUDITORIA-SWARM-01]
covers: [Q1]
---

# 🎯 INTENT: Ciclo de Correção Completa da Auditoria Profunda (AUD-03, AUD-04, AUD-05)

> **Status:** DRAFT  
> **Data:** 2026-11-14  
> **Autor:** Antigravity (Coordinator)  
> **Iniciativa:** Correção Auditoria Swarm  
> **Alvo:** `internal/studio/static/components/lego-item-modal.js`, `internal/studio/static/style.css`, `loops/*/*.loop.json`  
> **Sessão:** `2b93726e-e3b1-4f72-aa2a-d2ddcd8ca025`  
> **Confiança da Entrevista:** 98%  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** A auditoria profunda em 4 frentes identificou 3 achados remanescentes: emojis Unicode em `lego-item-modal.js` (`AUD-03`), ausência de breakpoint `<= 768px` em `style.css` (`AUD-04`) e 15 `.loop.json` sem `schema_version` (`AUD-05`).
* **Quem Sofre:** Operadores do LegoLoop Studio em telas menores e validadores de contrato legado.
* **O Porquê:** Zerar todos os achados da auditoria garante conformidade visual (sem emojis) e estrutural (responsividade e schema).

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo:
* [ ] Objetivo 1 (`AUD-03`): Substituir emojis Unicode em `internal/studio/static/components/lego-item-modal.js` por ícones SVG de `icons.js` ou texto limpo.
* [ ] Objetivo 2 (`AUD-04`): Adicionar media query `<= 768px` em `internal/studio/static/style.css` para recolher/adaptar a sidebar.
* [ ] Objetivo 3 (`AUD-05`): Padronizar `schema_version` nos 15 arquivos `.loop.json` em `loops/` mantendo `go test -count=1 ./...` 100% verde.

### O que NÃO ESTÁ no Escopo:
* ⛔ Não alterar contratos de API HTTP do Studio nem adicionar dependências JS/Go externas.

---

## 3. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-14_SPEC_CORRECAO_AUDITORIA_SWARM.md`](../specs/2026-11-14_SPEC_CORRECAO_AUDITORIA_SWARM.md)

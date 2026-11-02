# 🎯 INTENT: Auditoria Exaustiva Fase B — LegoLoop Full-Stack

> **Status:** Aprovado pelo Usuário
> **Data:** 2026-11-02
> **Autor:** Antigravity AI + Usuário
> **Alvo:** `internal/studio/static/` · `internal/studio/` · `internal/canon/` · `bin/legoloop.exe`

---

## 1. Contexto e Motivação

A Fase A catalogou 10 bugs superficiais. O usuário confirmou que a auditoria deve ser exaustiva e camada por camada — cobrindo cada subsistema com granularidade máxima: **arquivo + linha + causa-raiz + passos de reprodução + proposta cirúrgica de correção**.

A metodologia usa combo triplo por épico:
- **REA MCP** (`codebase-memory-mcp`) — análise estática, localiza causa no código sem servidor
- **Obscura** — smoke test headless (~85ms, ~30MB), verifica estado inicial do DOM
- **Playwright** — interação profunda em runtime, reproduz e documenta bugs

Servidor (`legoloop.exe studio`) **nasce e morre por épico** — zero resquícios em background.

---

## 2. Objetivos Principais

1. Cobrir 100% dos subsistemas do Studio UI com granularidade máxima
2. Produzir relatório master com matriz P0-P3 linkada a arquivo+linha
3. Servir como base para Fase C (esteira de correção TDD)

---

## 3. Escopo

* **O que entra:** Canvas · Conexões · Undo/Redo · Teclas de Atalho · Paleta & DnD · Inspetor · Inventário · Templates · Modal Run · Vault · Consolidação
* **O que NÃO entra:** Correção de código (Fase C) · Backend Go fora do Studio · CLI engines

---

## 4. Épicos — IDs Canônicos (Pivot de Rastreabilidade)

| ID | Nome | Método | Servidor? |
|----|------|--------|-----------|
| E01 | Canvas — Zoom, Pan, Limites, Ctrl+0 | REA MCP | ❌ |
| E02 | Conexões & Arestas — Arrastar, Vazio, Religar | REA + Obscura + Playwright | ✅ |
| E03 | Undo/Redo — Mutações, Snapshots, Tamanho | REA MCP | ❌ |
| E04 | Teclas de Atalho — Delete/Backspace, Escape, Ctrl+Z/Y/S | REA + Playwright | ✅ |
| E05 | Paleta & Drag-and-Drop — Offset, Virtualização, Busca | REA + Obscura + Playwright | ✅ |
| E06 | Inspetor — Reatividade, Debounce, Sugestões | REA + Playwright | ✅ |
| E07 | Inventário — Abas, Busca, Modal, Canonização | REA + Obscura + Playwright | ✅ |
| E08 | Templates — Carregamento, Confirmação, Fallback | REA + Playwright | ✅ |
| E09 | Modal Run — Polling/SSE, Badge, Logs, Fechar | REA + Playwright | ✅ |
| E10 | Vault — Dry-run, Apply, Feedback, Erro 503 | REA + Playwright | ✅ |
| E11 | Consolidação Final — Matriz de Bugs + Relatório Master | Síntese | ❌ |

---

## 5. Restrições Operacionais

- ❌ Zero `invoke_subagent`
- ❌ Zero polling de background tasks
- ✅ Servidor sobe e morre por épico (no Nano consolidador)
- ✅ Scripts temporários em `scratch/` removidos no Nano consolidador do mesmo épico
- ✅ Prefixo ISO em todos os artefatos persistentes
- ✅ Idioma PT-BR 100%
- ✅ Granularidade mínima: arquivo + linha + causa-raiz + reprodução + proposta cirúrgica

---

🛑 **Stop-Gate 01:** ✅ Aprovado pelo usuário → SPEC gerada.

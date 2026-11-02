# 🎯 Documento de Intenção (Intent) — Esteira de Correção Full-Stack LegoLoop

> **Prefixo ISO Canônico:** `2026-11-02_INTENT_ESTEIRA_CORRECOES_LEGOLOOP.md`  
> **Sessão:** `d4cfbcb7-ebf7-4f59-8d77-7c13903afc72`  
> **Status:** Proposto para Homologação (Stop-Gate SG-C01)  
> **Origem:** [`docs/reports/2026-11-02_REPORT_AUDITORIA_EXAUSTIVA.md`](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/reports/2026-11-02_REPORT_AUDITORIA_EXAUSTIVA.md)

---

## 1. Declaração de Propósito & Missão

Eliminar cirurgicamente todos os **20 defeitos factuais** mapeados durante a Auditoria Exaustiva da Fase B (1 P0, 4 P1, 10 P2 e 5 P3) nos subsistemas do LegoLoop (Studio Web UI, Web Components, rotas Go, CSP e integração de SO), operando em esteira de execução modular estéril por bloco, governada por um grafo acíclico direcionado (DAG) com roteamento de modelos cognitivos e verificação rigorosa por TDD Estrito (Red-Green-Refactor).

---

## 2. Invariantes Invioláveis da Fase C

1. **Zero Stubs:** Proibido introduzir `TODO`, `pass`, mocks parciais ou lógica simulada. Cada correção deve ser definitiva, funcional e testada.
2. **TDD Estrito:** Todo ticket deve possuir teste automatizado (Go ou Playwright E2E) comprovando a falha em vermelho antes da edição e a aprovação em verde após a edição.
3. **Isolamento Estéril por Bloco:** Cada nó do grafo DAG deve ler estritamente seu contrato de entrada e produzir seu artefato de saída namespaced em `scratch/C{N}_*.json`.
4. **Recompilação Fresca Obrigatória:** Qualquer ticket que modifique arquivos no backend Go ou arquivos estáticos embarcados (`staticFS`) deve acionar `go build -o bin/legoloop.exe ./cmd/legoloop`.
5. **Roteamento Cognitivo por Complexidade:**
   - Tarefas pontuais e sintáticas (ex: parâmetros, fallback de texto, atalhos simples): modelo ágil (`flash`).
   - Tarefas estruturais e concorrentes (ex: debounce de rede, AbortController, CSP, refatoração de DOM): modelo analítico profundo (`pro`).

---

## 3. Escopo dos 3 Lotes de Correção

* **Lote 01 — Bloqueadores & Segurança (P0 & P1):**
  - `TICK-C01`: `BUG-RUN-01` (`validate()` booleano) + `BUG-RUN-02` (`AbortController` no SSE).
  - `TICK-C02`: `BUG-SRV-01` (`icons.js` no mux Go) + Recompilação.
  - `TICK-C03`: `BUG-CSP-01` (Eliminação de inline style em `lego-item-modal.js`).
  - `TICK-C04`: `BUG-UNDO-01` (`renderAllNodes()`) + `BUG-UNDO-02` (`pushState()` em mutações).
* **Lote 02 — Consistência de UI, Concorrência e Performance (P2):**
  - `TICK-C05`: `BUG-CONFIRM-01` (`message:` no diálogo) + `BUG-TMPL-01` (Fallback de templates).
  - `TICK-C06`: `BUG-VAULT-01` (Renderização de grupos) + `BUG-VAULT-03` (Reset de estado disabled).
  - `TICK-C07`: `BUG-INSP-01` (Debounce e cancelamento em `/api/suggest`).
  - `TICK-C08`: `BUG-CANVAS-02` (Clamp no Pan) + `BUG-EDGE-01` (Escape no arraste de arestas).
  - `TICK-C09`: `BUG-PALETTE-02` (Virtualização do treeview) + `BUG-ITEM-MODAL-01` (Refatoração de DOM).
* **Lote 03 — Ergonomia e Atalhos (P3):**
  - `TICK-C10`: `BUG-KEY-01` (Escape desmarca seleção) + `BUG-KEY-02` (Atalho Ctrl+S) + `BUG-PALETTE-01` (Offset centralizado no Drop) + `BUG-CANVAS-03` (Pivot de Zoom) + `BUG-VAULT-02` (Toast em 503).

# 📋 SPEC: Auditoria Exaustiva Fase B — LegoLoop Full-Stack

> **Status:** Proposta para Revisão
> **Data:** 2026-11-02
> **Referência:** [2026-11-02_INTENT_AUDITORIA_EXAUSTIVA_LEGOLOOP.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/2026-11-02_INTENT_AUDITORIA_EXAUSTIVA_LEGOLOOP.md)

---

## 1. Visão Geral e Arquitetura

Auditoria cirúrgica profunda por decomposição em 11 Épicos e 34 Nodes, avaliando todo o ciclo de vida do Studio Web, bindings DOM, reatividade, segurança CSP/XSS e endpoints Go. Cada Node possui isolamento contratual estéril e critérios binários indexados pelo seu `node_id`.

---

## 2. Contratos de Dados & I/O

- **Entradas:** Grafo estático AST via `codebase-memory-mcp`, DOM live via `obscura`, interação Playwright via CDP (`:9222` ou standalone `:8080`).
- **Saídas:** Arquivos isolados em `scratch/E{N}_{M}_{desc}.json` e consolidados em `scratch/audit_E{N}_{tema}.json`.
- **Relatório Final:** [docs/reports/2026-11-02_REPORT_AUDITORIA_EXAUSTIVA.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/reports/2026-11-02_REPORT_AUDITORIA_EXAUSTIVA.md).

---

## 3. Critérios Binários de Aceite por Node

### E01 — Canvas (Zoom, Pan, Limites, Ctrl+0)
* [ ] `AC-E01.1` Identificadas funções de zoom (wheel, botões +/-) com arquivo e linha exatos em `app.js`.
* [ ] `AC-E01.2` Identificada função de pan e verificado se existe clamp de limites no canvas.
* [ ] `AC-E01.3` Auditado `Ctrl+0` / `btn-zoom-fit` avaliando se calcula bounding-box real de nós ou mero reset.
* [ ] `AC-E01.4` Artefato `scratch/audit_E01_canvas.json` gerado com exit code 0 e consolidação completa.

### E02 — Conexões & Arestas
* [ ] `AC-E02.1` Rastreamento completo de criação de arestas (mousedown port, drag, drop, guardas de aresta fantasma).
* [ ] `AC-E02.2` Smoke test via Obscura confirmando subida do Studio na porta 8080 e captura de tela inicial.
* [ ] `AC-E02.3` Execução Playwright validando: cabo solto no vazio, cancelamento com Escape e religação de aresta.
* [ ] `AC-E02.4` Encerramento forçado do processo do servidor (`legoloop.exe`) e persistência de `scratch/audit_E02_conexoes.json`.

### E03 — Undo/Redo
* [ ] `AC-E03.1` Mapeamento da pilha de histórico, mutações registradas e callers em `app.js`.
* [ ] `AC-E03.2` Confirmação se snapshots ocorrem em `mouseup` ou `mousemove` (risco de saturação).
* [ ] `AC-E03.3` Persistência de `scratch/audit_E03_undoredo.json`.

### E04 — Teclas de Atalho
* [ ] `AC-E04.1` Mapeamento do listener `keydown` global e verificação de guard contra target em input/textarea.
* [ ] `AC-E04.2` Validação Playwright de Delete com input focado, Escape em modais e Ctrl+Z em execução ativa.
* [ ] `AC-E04.3` Servidor finalizado e persistência de `scratch/audit_E04_atalhos.json`.

### E05 — Paleta & Drag-and-Drop
* [ ] `AC-E05.1` Análise de cálculo de offset no drop (compensação de zoom + pan).
* [ ] `AC-E05.2` Análise de virtualização (`virtual.js`) com filtragem concorrente.
* [ ] `AC-E05.3` Teste Playwright com drag-and-drop e medição de precisão de coordenadas.
* [ ] `AC-E05.4` Limpeza de processo e persistência de `scratch/audit_E05_paleta.json`.

### E06 — Inspetor
* [ ] `AC-E06.1` Análise estática de reatividade do inspetor e falta de debounce em `/api/suggest`.
* [ ] `AC-E06.2` Teste Playwright medindo rajadas de requests e race condition ao digitar no inspetor.
* [ ] `AC-E06.3` Limpeza de processo e persistência de `scratch/audit_E06_inspetor.json`.

### E07 — Inventário
* [ ] `AC-E07.1` Auditoria de segurança de injeção HTML (`innerHTML`) em `lego-item-modal.js`.
* [ ] `AC-E07.2` Teste dinâmico das abas de inventário, pesquisa acentuada e botões de ação do Web Component.
* [ ] `AC-E07.3` Limpeza de processo e persistência de `scratch/audit_E07_inventario.json`.

### E08 — Templates
* [ ] `AC-E08.1` Causa-raiz da duplicação de nomes e fallback de ícones mapeados em `templates.go` / `app.js`.
* [ ] `AC-E08.2` Teste de confirmação de sobreposição ao carregar templates em canvas não-vazio.
* [ ] `AC-E08.3` Limpeza de processo e persistência de `scratch/audit_E08_templates.json`.

### E09 — Modal de Execução (Run)
* [ ] `AC-E09.1` Mapeamento do fluxo de disparo de `#btn-run`, impacto do bug de `validate()` e polling vs SSE.
* [ ] `AC-E09.2` Análise do backend `run_api.go` sobre o protocolo de streaming de logs.
* [ ] `AC-E09.3` Teste Playwright validando transições de badge e fechamento do modal durante run ativo.
* [ ] `AC-E09.4` Limpeza de processo e persistência de `scratch/audit_E09_run.json`.

### E10 — Cofre Canônico (Vault)
* [ ] `AC-E10.1` Mapeamento dos handlers `/api/canonize/plan` e `/apply` e tratamento de `s.Assets == nil`.
* [ ] `AC-E10.2` Teste Playwright de simulação de erro 503 e feedback de dry-run/apply.
* [ ] `AC-E10.3` Limpeza de processo e persistência de `scratch/audit_E10_vault.json`.

### E11 — Consolidação Final
* [ ] `AC-E11.1` Ingestão e unificação de todos os 10 JSONs de auditoria com deduplicação.
* [ ] `AC-E11.2` Produção do relatório exaustivo master com matriz P0-P3, causas e propostas cirúrgicas.
* [ ] `AC-E11.3` Limpeza total de arquivos temporários de `scratch/` e confirmação de zero processos órfãos.

---

## 4. Critério de Done Global

> A demanda é **CONCLUÍDA** quando 100% dos critérios `AC-E01.1` até `AC-E11.3` estiverem marcados como `[x]` e validados no Stop-Gate 04.

---
🛑 **Stop-Gate 02:** Aguardando aprovação do usuário para o PLAN e ROUTE.

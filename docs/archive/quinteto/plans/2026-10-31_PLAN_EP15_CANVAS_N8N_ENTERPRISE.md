# PLAN_EP15: Plano de Execução — Canvas N8N Enterprise, Controle de Fluxo & Ergonomia

> **Referência:** [SPEC_EP15_CANVAS_N8N_ENTERPRISE.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/SPEC_EP15_CANVAS_N8N_ENTERPRISE.md)  
> **Status:** Concluído com Sucesso  
> **Estratégia:** Micro-Nano Tickets, TDD Estrito, Zero Stubs, `funclen` ≤ 35 linhas.

---

### 📦 Fases e Micro-Nano Tickets

#### Fase 1: Backend Go & Ontologia (TICK-1501 a TICK-1504)
- [x] **TICK-1501:** Flexibilizar `ParseArtifactType` em `internal/graph/ontology.go` para aceitar nomes sem prefixo (`URL`, `Spec`, etc.). Teste unitário em `graph_test.go`.
- [x] **TICK-1502:** Cadastrar nós estruturais de controle (`flow_if`, `flow_loop`, `flow_merge`) em `internal/graph/ontology.go` e `index.go`.
- [x] **TICK-1503:** Implementar endpoint `POST /api/assets/reveal` em `internal/studio/reveal.go` para abrir pastas no Windows Explorer com segurança contra traversal. Teste unitário em `reveal_test.go`.
- [x] **TICK-1504:** Implementar endpoint `POST /api/canonize/piece` em `internal/studio/canon_piece.go` para canonizar apenas o hash da peça selecionada. Teste unitário correspondente.

#### Fase 2: Layout Reativo e Ergonomia (TICK-1505 a TICK-1507)
- [x] **TICK-1505:** Expandir Sidebar para 320px em `style.css` e implementar o CSS Grid do `.workspace` com suporte a variáveis `--sidebar-w` e `--inspector-w` com transição suave.
- [x] **TICK-1506:** Adicionar botão de colapso no Inspector direito (`btn-toggle-inspector`) com alternância de classe `.is-inspector-collapsed` em `index.html`, `style.css` e `app.js`.
- [x] **TICK-1507:** Agrupar cópias por SHA-256 no catálogo de peças do Canvas (`palette.js`), exibindo badge de contagem idêntico ao do Inventário e nova categoria "Controle".

#### Fase 3: Canvas N8N Core & Desconexão (TICK-1508 a TICK-1510)
- [x] **TICK-1508:** Alterar cursor padrão do Canvas para `default`; implementar listener de `Spacebar` para modo Pan (`grab`/`grabbing`) com pan contínuo sem limites de $X/Y$.
- [x] **TICK-1509:** Aprimorar o botão circular `✕` nas arestas com área de clique confortável (raio 12px), highlight no hover da aresta e exclusão instantânea no clique.
- [x] **TICK-1510:** Implementar animação vetorial luminosa de fluxo nas arestas em `style.css` via `stroke-dasharray` e `stroke-dashoffset` (`.edge--animated`).

#### Fase 4: Nós de Controle de Fluxo & Micro-Governança (TICK-1511 a TICK-1513)
- [x] **TICK-1511:** Cadastrar blocos de controle de fluxo algorítmico (`IF`, `Loop`, `Merge`) disponíveis na paleta do Studio.
- [x] **TICK-1512:** Exibir badges discretos no rodapé dos nós (`rules`, `pre`, `gates`, `post`, `scripts`).
- [x] **TICK-1513:** No Inspector, separar claramente slots `[Nativo]` de slots `[Custom]` e permitir inclusão de Regras (`Rules`).

#### Fase 5: Modal Rico Enterprise com Dimensão Fixa (TICK-1514 a TICK-1516)
- [x] **TICK-1514:** Fixar dimensões do modal (880px × 640px) no CSS; adicionar botões no cabeçalho *"Canonizar Peça"* e *"Abrir no Explorer"*.
- [x] **TICK-1515:** Implementar diagrama visual **IPO** (Input ➔ Processamento ➔ Output) e badges de artefatos na sub-aba "Visão Geral".
- [x] **TICK-1516:** Implementar filtros por categoria e paginação na sub-aba "Peças Compatíveis"; adicionar botão "Abrir Pasta" por item na sub-aba "Ocorrências".

#### Fase 6: Produtividade Visual & Histórico (TICK-1517 a TICK-1519)
- [x] **TICK-1517:** Implementar histórico de Desfazer/Refazer (`Ctrl+Z` / `Ctrl+Y`) para nós e arestas no Canvas.
- [x] **TICK-1518:** Implementar botão de Auto-Layout (organização hierárquica automática de nós Sugiyama).
- [x] **TICK-1519:** Adicionar botões flutuantes de produtividade no Canvas (Undo, Redo, Auto-Layout, Exportar Imagem).

#### Fase 7: Validação & Auditoria E2E (TICK-1520)
- [x] **TICK-1520:** Executar `go test ./...`, validar linter de `innerHTML` e conduzir auditoria ponta a ponta via **Obscura + Playwright MCP** com geração de capturas de tela.

# 📋 PLAN_EP17: Triggers de Primeira Classe no Catálogo Lateral e Canvas N8N

> **Documento:** `docs/plans/2026-11-01_PLAN_EP17_TRIGGERS_CATALOG_CANVAS.md`  
> **Status:** Pronto para Execução  
> **Data:** 01/11/2026  
> **Mega-Pipeline Loop 3:** Engenharia & Decomposição em Micro/Nano-tickets

---

## Matriz de Roteamento de Tickets

| Ticket | Escopo & Descrição | Categoria | Modelo | Arquivos Alvo | Gate de Sucesso |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **TICK-1701** | **Ontologia Go: `TriggerBlocks`**<br>Definir os 5 blocos canônicos de Trigger em `internal/graph/ontology.go` e incluir em `FullCatalog()`. | Nano | `flash` | [`internal/graph/ontology.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/graph/ontology.go)<br>[`internal/studio/pieces_test.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/pieces_test.go) | `go test ./internal/graph ./internal/studio` verde. |
| **TICK-1702** | **Paleta & Treeview no Menu Lateral**<br>Adicionar categoria `Triggers` no topo de `PALETTE_ORDER`, ícone SVG de raio e classificação automática de blocos `trigger_*`. | Micro | `flash` | [`internal/studio/static/palette.js`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/palette.js)<br>[`internal/studio/static/style.css`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/style.css) | Menu lateral exibe "Triggers" com badge de contagem e os 5 blocos. |
| **TICK-1703** | **Canvas N8N: Estilização do Nó Raiz**<br>Renderizar nós trigger com estilo `.node--trigger`, omitir porta de entrada (`node__port--in`) e impor restrição de gatilho único no canvas. | Micro | `flash` | [`internal/studio/static/app.js`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/app.js)<br>[`internal/studio/static/style.css`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/style.css) | Nó trigger aparece destacado em âmbar/verde sem porta de entrada. |
| **TICK-1704** | **Inspetor de Gatilho & Serialização DAG v2**<br>Painel no Inspetor para parâmetros do trigger (`cron`, `path`, `loop`) e injeção do objeto `trigger` no payload de `/api/compose`. | Micro | `flash` | [`internal/studio/static/app.js`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/app.js)<br>[`internal/studio/compose.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/compose.go) | Compor esteira envia `trigger` correto ao backend. |
| **TICK-1705** | **Homologação, Testes Gerais & Recompilação Go**<br>Executar suíte de testes de todo o repositório, recompilar o binário `bin/legoloop.exe` e emitir relatório ISO. | Nano | `flash` | [`bin/legoloop.exe`](file:///C:/Users/trcnologia/Desktop/LegoLoop/bin/legoloop.exe)<br>[`docs/reports/2026-11-01_REPORT_EP17_TRIGGERS_CATALOG_CANVAS.md`](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/reports/2026-11-01_REPORT_EP17_TRIGGERS_CATALOG_CANVAS.md) | `go test ./...` 100% verde; binário atualizado. |

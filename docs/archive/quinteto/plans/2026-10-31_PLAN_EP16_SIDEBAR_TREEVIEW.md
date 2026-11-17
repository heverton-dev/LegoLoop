# 📋 PLAN_EP16_SIDEBAR_TREEVIEW — Plano de Execução do Mega-Pipeline

> **Caso de Uso:** Refatoração Estrutural e Arquitetura Sanfona (Treeview) do Menu Lateral Esquerdo do LegoLoop Studio.  
> **Status:** Em Execução  
> **Data:** 2026-10-31  
> **Regras:** Zero Stubs | Caveman Thinking | Roteamento de Modelos | Execução Híbrida (In-Line / Orchestrix) | Zero Lixo de Screenshots

---

## 🧭 Tabela de Micro e Nano Tickets

| Ticket | Escopo & Descrição | Modelo Sugerido | Modo de Execução | Arquivos Impactados | Gate / Critério de Aceite |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **TICK-1601** | **Confinamento CSS & Viewport**<br>Eliminar aninhamento que causava altura de 208.620px. Confinar `.sidebar` a `calc(100vh - 56px)`, container flex vertical estrito com `overflow-hidden` e sub-container com `flex: 1; min-height: 0; overflow-y: auto`. | `flash` | `In-Line` | [`internal/studio/static/style.css`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/style.css), [`internal/studio/static/index.html`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/index.html) | Altura do `#palette-list` <= 600px; zero scrollbar dupla; zero vazamento de viewport. |
| **TICK-1602** | **Treeview de Peças Disponíveis**<br>Agrupar as peças em subpastas sanfona por tipo (`Controle`, `Skills`, `MCPs`, `Loops`, `Scripts`, `Hooks`, `Gates`) com setas `▾` / `▸`, ícones distintivos e badge com contagem em cada grupo. | `flash` | `In-Line` | [`internal/studio/static/palette.js`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/palette.js) | Renderizar 7 grupos retráteis; badge do botão principal `Peças Disponíveis (<total>)`. |
| **TICK-1603** | **Treeview de Modelos Vivos por Harness**<br>Agrupar os modelos vivos por Provedor / Harness (`9Router VPS`, `Claude CLI`, `Ollama Local`, `Gemini / AGY`) em sanfonas individuais com status pulsante verde e contadores parciais. | `flash` | `In-Line` | [`internal/studio/static/palette.js`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/palette.js) | Agrupamento de 225+ modelos por harness com contagens parciais e badge global no botão. |
| **TICK-1604** | **Lógica de Sanfona Exclusiva & Ações de Clique**<br>Ao expandir Peças, colapsar Modelos (e vice-versa); clique simples em uma peça adiciona o bloco no centro do canvas; duplo clique abre o Modal Rico de Raio-X; arrastar/soltar mantido. | `flash_lite` | `In-Line` | [`internal/studio/static/palette.js`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/palette.js), [`internal/studio/static/app.js`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/app.js) | Alternância sanfona suave; clique simples insere nó no canvas; duplo clique abre modal raio-x. |
| **TICK-1605** | **Impeccable Design Pass (Polimento Visual)**<br>Tipografia de alta legibilidade, espaçamentos compactos (4px / 8px), contrastes refinados para tema escuro e claro, hover states táteis e micro-animações nas setas sanfona. | `pro` / `inherit` | `In-Line` | [`internal/studio/static/style.css`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/style.css) | Hierarquia visual impecável; ausência de quebras de layout ou textos truncados indevidamente. |
| **TICK-1606** | **Gate de Segurança & Linting**<br>Verificação de injeção de HTML não sanitizado (`TestNoInnerHTMLInScripts`), conformidade de tamanho de função (`funclen`) e build dos assets do Studio. | `inherit` | `In-Line` | Go Tests & Linter | `go test ./...` verde; zero innerHTML vulnerável; zero stubs. |
| **TICK-1607** | **Homologação E2E (Obscura + Playwright MCP)**<br>Validação sem prints desnecessários (Zero Lixo de Drive). Auditoria via asserção de DOM em memória: rolagem, expansão/colapso, contagens e clique de inserção no canvas. | `inherit` | `In-Line` | E2E Browser | Todas as asserções de DOM e usabilidade aprovadas no Studio ao vivo (`:7878`). |

---

## 🔄 Fluxo de Retorno e Retroalimentação
Se qualquer ticket falhar em seus critérios de aceitação:
1. Reabrir ticket específico com análise da causa raiz;
2. Executar refatoração corretiva;
3. Re-executar baterias do Gate e E2E até 100% verde.

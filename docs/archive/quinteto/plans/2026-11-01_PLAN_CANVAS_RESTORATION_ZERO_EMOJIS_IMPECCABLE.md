# 🗺️ PLAN: Restauração Funcional do Canvas Studio (REA) & Design Pass Impeccable (Zero Emojis)

> **Status:** Proposta para Revisão (STOP-GATE HUMANO 02)  
> **Data:** 2026-11-01  
> **Especificação Base:** `docs/specs/2026-11-01_SPEC_CANVAS_RESTORATION_ZERO_EMOJIS_IMPECCABLE.md`  
> **Metodologia:** Tracer Bullet Development + TDD Estrito (Red ➔ Green ➔ Refactor) + Silent Execution

---

## 🏛️ ÉPICO 18: Engenharia Reversa (REA) & Restauração Funcional do Canvas

> **Objetivo:** Auditar, corrigir e estabilizar todos os event listeners do Canvas 2D (pan, zoom, drag, clique e atalhos).  
> **Modelo Recomendado:** Claude 3.7 Sonnet / Gemini 2.0 Pro (`--effort high`)

### 🔹 Micro-ticket TICK-1801: Auditoria Forense e Estabilização de Pan e Zoom 2D
* **Tipo:** Micro-ticket
* **Arquivo Alvo:** `internal/studio/static/app.js`, `internal/studio/static/style.css`
* **Rota:** Gemini 2.0 Flash (`--effort medium`)
* **Entregáveis & Testes TDD:**
  * Refatorar o handler de `wheel` para suportar tanto trackpad (pan natural em X/Y) quanto mouse wheel com e sem Ctrl.
  * Estabilizar o listener de `pointerdown` no `#canvas` e `#canvas-surface`: suporte a botão do meio (`button === 1`), barra de espaço (`state.spacePan`) e clique no fundo vazio sem dependência de classes frágeis.
  * Suíte de testes de regressão de eventos no JS (`tests/canvas_events_test.js` ou teste Go de assets).

### 🔹 Micro-ticket TICK-1802: Unificação e Estabilização do Drag-and-Drop de Cards e Seleção
* **Tipo:** Micro-ticket
* **Arquivo Alvo:** `internal/studio/static/app.js`
* **Rota:** Gemini 2.0 Flash (`--effort medium`)
* **Entregáveis & Testes TDD:**
  * Permitir que o clique em qualquer parte não interativa do card selecione o nó e abra o Inspetor.
  * Corrigir `startNodeDrag` para ligar `pointermove` e `pointerup` no `window` com `setPointerCapture` seguro, evitando nós "presos" ou perdendo o foco quando o mouse acelera.
  * Garantir que o clique nas portas de I/O continue disparando a criação de conexões sem acionar o arrasto do card.

### 🔸 Nano-ticket TICK-1803: Restauração de Atalhos de Teclado Globais
* **Tipo:** Nano-ticket
* **Arquivo Alvo:** `internal/studio/static/app.js`
* **Rota:** Gemini 2.0 Flash-Lite (`--effort low`)
* **Entregáveis & Testes TDD:**
  * Garantir que `Delete`, `Backspace`, `Ctrl+Z`, `Ctrl+Y`, `Ctrl+0` funcionem a partir de qualquer ponto da tela quando nenhum formulário (`input, textarea, select`) estiver ativo.

### 🛑 STOP-GATE DE AUDITORIA HUMANA — ÉPICO 18
* [ ] Pan, zoom, atalhos e seleção de nós 100% operantes e fluidos no navegador.
* [ ] Auditoria do usuário aprovada para avançar ao Épico 19.

---

## 🏛️ ÉPICO 19: Design Pass Impeccable (Regra Zero Emojis & Cores Ontológicas)

> **Objetivo:** Erradicar 100% dos emojis, adotar biblioteca SVG pura e implantar código de cores estrutural.  
> **Modelo Recomendado:** Claude 3.7 Sonnet / Gemini 2.0 Pro (`--effort high`)

### 🔹 Micro-ticket TICK-1901: Barreira Automatizada Anti-Emojis no Go (TDD Estrito)
* **Tipo:** Micro-ticket
* **Arquivo Alvo:** `internal/studio/emojis_test.go`
* **Rota:** Gemini 2.0 Flash (`--effort medium`)
* **Entregáveis & Testes TDD:**
  * Escrever teste Go (`TestZeroEmojisInFrontendAndTemplates`) que lê iterativamente todos os arquivos de `internal/studio/static/` e `templates.go`.
  * O teste DEVE falhar inicialmente (Red) comprovando a existência de emojis nos arquivos atuais.

### 🔹 Micro-ticket TICK-1902: Erradicação de Emojis e Implantação de SVGs Canônicos
* **Tipo:** Micro-ticket
* **Arquivo Alvo:** `internal/studio/static/index.html`, `internal/studio/static/app.js`, `internal/studio/templates.go`, `internal/studio/static/icons.js`
* **Rota:** Gemini 2.0 Flash (`--effort medium`)
* **Entregáveis & Testes TDD:**
  * Criar catálogo canônico de ícones SVG inline (`icons.js` / helpers) para `templates`, `vault`, `camera`, `zoom`, `delete`, `info`, etc.
  * Substituir todos os emojis em botões e templates Go por identificadores de ícones SVG.
  * Rodar o teste `TestZeroEmojisInFrontendAndTemplates` até ficar 100% verde (Green).

### 🔹 Micro-ticket TICK-1903: Taxonomia Estrutural de Cores por Tipo de Peça no Canvas
* **Tipo:** Micro-ticket
* **Arquivo Alvo:** `internal/studio/static/style.css`, `internal/studio/static/app.js`
* **Rota:** Gemini 2.0 Flash (`--effort medium`)
* **Entregáveis & Testes TDD:**
  * Aplicar cores CSS semânticas específicas por categoria estrutural (Triggers, Controles, Skills, MCPs, Loops, Scripts, Hooks e Gates).
  * Estilizar acentos de cabeçalho, borda esquerda e portas com a cor respectiva.

### 🛑 STOP-GATE DE AUDITORIA HUMANA — ÉPICO 19
* [ ] Suíte de testes `go test ./...` 100% verde (incluindo teste anti-emojis).
* [ ] Zero emojis visíveis no Studio e paleta de cores corporativa ativa no Canvas.
* [ ] Recompilação obrigatória: `go build -o bin/legoloop.exe ./cmd/legoloop`.

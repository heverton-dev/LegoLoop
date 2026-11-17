# 📋 SPEC: Restauração Funcional do Canvas Studio (REA) & Design Pass Impeccable (Zero Emojis)

> **Status:** Proposta para Revisão (STOP-GATE HUMANO 02)  
> **Data:** 2026-11-01  
> **Épicos Alvo:** ÉPICO 18 (Engenharia Reversa e Eventos do Canvas) e ÉPICO 19 (Design Pass Impeccable / Zero Emojis)  
> **Referência:** `docs/intents/2026-11-01_INTENT_ONTOLOGIA_ASSETS_SO_E_MEGA_PIPELINE_TEMPLATES.md`

---

## 1. Visão Geral e Arquitetura

Esta especificação formaliza a resolução de dois pilares críticos do **LegoLoop Studio**:
1. **Restauração Completa da Interação do Canvas (ÉPICO 18):** Diagnóstico e conserto dos *event listeners* do Canvas (`internal/studio/static/app.js`), garantindo fluidez absoluta em:
   * Pan/Scroll 2D livre (Espaço + Drag, Botão do meio do mouse, arrasto no fundo do canvas e suporte nativo a trackpad/touchpad).
   * Zoom contínuo centrado no cursor via `Ctrl + Wheel` e botões de zoom in/out/fit.
   * Seleção e clique em cards de nós em qualquer área não interativa (cabeçalho, corpo ou badges).
   * Drag-and-drop de nós sem perda de ponteiro (`window.setPointerCapture` e handlers na janela).
   * Atalhos de teclado universais operantes (`Delete`/`Backspace` para nós/arestas, `Ctrl+Z`/`Ctrl+Y`, `Ctrl+0`, `Ctrl+B`).
2. **Design Pass Impeccable & Lei Zero Emojis (ÉPICO 19):**
   * **Erradicação absoluta de emojis** em todo o frontend (`index.html`, `app.js`, `inspector.js`, `inventory.js`, `palette.js`) e no backend Go (`templates.go`).
   * Implantação de um **Sistema Canônico de Ícones SVG Inline** (leves, nítidos, 100% vetoriais, sem dependências externas).
   * Implantação da **Taxonomia Estrutural de Cores** no Canvas por tipo de peça (Triggers, Controles, Skills, MCPs, Loops, Scripts, Hooks e Gates).

---

## 2. Critérios de Aceitação Binários (Gates Mecânicos)

### 2.1 ÉPICO 18: Engenharia Reversa & Restauração Funcional do Canvas
* [ ] **CA-18.1 (Pan 2D Confiável):** Pressionar e segurar `Espaço` + arrastar com botão esquerdo, ou arrastar diretamente com o botão do meio (scroll click), ou arrastar o fundo vazio do canvas realiza pan infinito contínuo e suave sem travar.
* [ ] **CA-18.2 (Trackpad & Wheel Pan):** Rolar o mouse/trackpad sem teclas modificadoras desloca o canvas em X e Y (`deltaX` e `deltaY`) com propagação interrompida para o container pai.
* [ ] **CA-18.3 (Zoom Centrado no Cursor):** `Ctrl + Wheel` (ou `Cmd + Wheel`) executa zoom contínuo mantendo o ponto exato sob o cursor estável na tela, com clamping rígido entre $0.2\times$ e $3.0\times$.
* [ ] **CA-18.4 (Seleção de Nós Holística):** Clicar em qualquer parte do card do nó (cabeçalho, corpo, borda, exceto conectores de portas) seleciona o nó imediatamente e abre seus dados no Inspetor direito.
* [ ] **CA-18.5 (Arrasto de Card Estável):** Iniciar o arrasto em qualquer área do cabeçalho ou corpo do nó (exceto slots/botões interativos) move o nó com fidelidade $1:1$ com o cursor; o arrasto não se perde se o ponteiro sair rapidamente do card (`window.pointermove`/`pointerup`).
* [ ] **CA-18.6 (Atalhos de Teclado Globais):** As teclas `Delete`/`Backspace` removem o nó ou aresta selecionada (a menos que um `input`, `textarea` ou `select` esteja focado). `Ctrl+Z` desfaz, `Ctrl+Y` refaz, `Ctrl+0` enquadra à tela.

### 2.2 ÉPICO 19: Design Pass Impeccable, Regra Zero Emojis & Cores Ontológicas
* [ ] **CA-19.1 (Barreira Automatizada Anti-Emojis):** Teste unitário em Go (`internal/studio/emojis_test.go`) varre compulsoriamente todos os arquivos de `internal/studio/static/` (`.html`, `.js`, `.css`) e `templates.go`, falhando se encontrar qualquer caractere no range de emojis unicode (`\x{1F300}-\x{1F9FF}`, `\x{2600}-\x{26FF}`, `\x{2700}-\x{27BF}`, etc.).
* [ ] **CA-19.2 (Erradicação Total no Código):** Zero ocorrências de emojis literais em botões de cabeçalho (`btn-templates`, `btn-vault`, `btn-export-png`), paleta de templates ou cards.
* [ ] **CA-19.3 (Biblioteca Canônica de Ícones SVG):** Criação do helper/componente SVG reutilizável (`icons.js` ou template SVG) fornecendo ícones monocromáticos otimizados (`icon-template`, `icon-vault`, `icon-camera`, `icon-trigger`, `icon-skill`, `icon-mcp`, `icon-loop`, `icon-script`, `icon-hook`, `icon-gate`).
* [ ] **CA-19.4 (Taxonomia Estrutural de Cores no Canvas):** Nós recebem classes CSS determinísticas e acentos de cor na borda e conectores baseados no tipo estrutural da peça:
  * Triggers: `#F97316` (Laranja)
  * Controles: `#F59E0B` (Âmbar)
  * Skills: `#8B5CF6` (Roxo)
  * MCPs: `#06B6D4` (Ciano)
  * Loops: `#2563EB` (Azul)
  * Scripts: `#10B981` (Verde)
  * Hooks: `#EC4899` (Rosa)
  * Gates: `#EF4444` (Vermelho)
* [ ] **CA-19.5 (Conformidade de Estilo Corporativo):** O layout visual atende às diretrizes do *Impeccable*: tipografia limpa, micro-interações refinadas, espaçamento coerente e visual enterprise.

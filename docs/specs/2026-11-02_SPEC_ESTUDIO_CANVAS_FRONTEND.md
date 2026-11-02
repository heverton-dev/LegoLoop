# SPEC — Correções Estruturais de Front-End, Canvas e Modais no LegoLoop Studio

**Data:** 2026-11-02  
**Referência Intent:** [docs/intents/2026-11-02_INTENT_ESTUDIO_CANVAS_FRONTEND.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/2026-11-02_INTENT_ESTUDIO_CANVAS_FRONTEND.md)  
**Status:** Proposto  

---

## 1. Requisitos Funcionais

### RF-01: Desbloqueio de CSP para Estilos Dinâmicos do Canvas
- O servidor Go em `internal/studio/security.go` deve fornecer diretivas CSP (`style-src 'self' 'unsafe-inline'`) permitindo que o Canvas manipule `style.transform`, `style.left`, `style.top` e `style.display` sem lançar erros no console do navegador.

### RF-02: Normalização do Scroll e Pan do Canvas
- O elemento `#canvas` deve utilizar rolagem virtual infinita (`overflow: hidden` sem scrollbars nativas do SO interferindo no bounding box).
- O cálculo `canvasPoint(ev)` deve ser 100% determinístico com zoom e pan, garantindo precisão milimétrica de cliques em nós, conexões de portas e drag-and-drop.
- O scroll do mouse wheel deve suportar:
  - Scroll vertical/horizontal para Pan do canvas.
  - `Ctrl + Wheel` para Zoom in / Zoom out focado na posição do cursor.

### RF-03: Toggle de Ligar/Desligar Card (Ativar / Desativar Peça)
- Cada card renderizado no canvas deve possuir no cabeçalho (`.node__head`) um controle visual de interruptor/toggle (Ligar / Desligar).
- Quando desativado:
  - O card recebe a classe CSS `.is-disabled` (opacidade reduzida, badge de "Inativo / Bypass").
  - O estado de ativação é persistido no nó (`node.enabled = false`) e no payload da DAG/composição.
  - Conexões conectadas ao nó desativado permanecem visíveis com estilo tracejado/desbotado.
  - Ao executar, o nó desativado é tratado como bypass ou pulado sem quebrar o grafo.

### RF-04: Robustez e Visibilidade de Modais
- O modal de Execução (`#modal-run`) deve abrir corretamente e exibir status, mesmo se o canvas estiver vazio (exibindo empty state descritivo em vez de apenas toast genérico) ou com nós válidos.
- O modal de Templates (`#modal-templates`) deve carregar e exibir os templates canônicos imediatamente, garantindo que os cards de template estejam visíveis e clicáveis.
- O modal do Cofre (`#modal-vault`) e Importação (`#modal-import`) devem fechar suavemente via tecla `ESC` ou botão `✕`.

---

## 2. Critérios Binários de Aceite (Stop Gates)

- [ ] **G1 (Zero Erros de CSP):** Ao carregar `http://127.0.0.1:7878`, o console do navegador não deve registrar nenhuma violação de Content Security Policy.
- [ ] **G2 (Scroll/Pan Estável):** Arrastar o canvas ou usar wheel move o `canvas-surface` suavemente sem quebrar o alinhamento de cliques em portas e nós.
- [ ] **G3 (Toggle de Peça Operacional):** Clicar no interruptor de ligar/desligar de qualquer card alterna seu estado visual e lógico entre ativo e inativo.
- [ ] **G4 (Modais Funcionais):** Todos os 4 modais principais abrem, renderizam conteúdo visível e fecham via botão de fechar e ESC.
- [ ] **G5 (Suíte Go & E2E Verde):** `go test ./internal/studio/...` e teste visual automatizado executam com 100% de sucesso.

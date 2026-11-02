# PLAN — Plano de Implementação: Correções de Front-End e UX do Studio

**Data:** 2026-11-02  
**Referência Spec:** [docs/specs/2026-11-02_SPEC_ESTUDIO_CANVAS_FRONTEND.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-11-02_SPEC_ESTUDIO_CANVAS_FRONTEND.md)  
**Status:** Planejado  

---

## Decomposição de Tickets (Tracer-Bullets)

### Ticket 1: Ajuste de CSP e Headers de Segurança
* **Alvo:** [`internal/studio/security.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/security.go) e [`internal/studio/security_test.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/security_test.go)
* **Ação:** Atualizar a diretiva `style-src` no cabeçalho CSP para permitir inline styles (`'self' 'unsafe-inline'`), necessários para a computação dinâmica do canvas no DOM.
* **Critério de Aceite:** Console do browser sem erros de CSP ao inicializar e manipular o canvas. Testes unitários de `security_test.go` atualizados e verdes.

### Ticket 2: Normalização de Scroll, Viewport e Pan 2D do Canvas
* **Alvo:** [`internal/studio/static/style.css`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/style.css) e [`internal/studio/static/app.js`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/app.js)
* **Ação:**
  1. Remover `overflow: auto` conflitante do CSS de `.canvas`, definindo rolagem puramente virtual com `overflow: hidden`.
  2. Ajustar os manipuladores de `wheel` e `pointermove` para pan contínuo e zoom sem offsets espúrios.
  3. Garantir consistência milimétrica no cálculo de coordenadas de drag de nós e portas.
* **Critério de Aceite:** Scroll suave via mouse wheel ou arraste de tela sem surgimento de barras de rolagem nativas e sem desalinhamento das conexões.

### Ticket 3: Implementação do Switch de Ligar/Desligar Card (Enable/Disable Toggle)
* **Alvo:** [`internal/studio/static/app.js`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/app.js) e [`internal/studio/static/style.css`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/style.css)
* **Ação:**
  1. Adicionar componente de interruptor/toggle no `.node__head` com tooltip explicativo ("Ligar / Desligar peça").
  2. Implementar evento de clique que alterna `node.enabled` (padrão: `true`).
  3. Aplicar estilo visual de card desligado (`.is-disabled`, visual translúcido/faded e indicador desligado).
  4. Manter o estado serializável na composição do loop.
* **Critério de Aceite:** Clicar no interruptor alterna visualmente e no modelo o estado ativo/inativo do nó sem disparar arraste acidental do card.

### Ticket 4: Correção de Modais (Execução, Templates e Visibilidade)
* **Alvo:** [`internal/studio/static/app.js`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/app.js), [`internal/studio/static/style.css`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/style.css) e [`internal/studio/static/index.html`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/index.html)
* **Ação:**
  1. Garantir que o modal de templates exiba a lista de templates de forma visível e clicável, sem timeouts de renderização.
  2. Ajustar o fluxo do botão Executar para abrir o modal com feedback explicativo se o canvas estiver vazio.
  3. Padronizar o fechamento e acessibilidade de todos os `<dialog>` do Studio.
* **Critério de Aceite:** Todos os modais abrem, exibem seu conteúdo de forma visível e fecham perfeitamente via ESC e botões dedicados.

### Ticket 5: Homologação Visual End-to-End no Navegador Real
* **Alvo:** Script de homologação em [`scratch/audit_studio_interactions.py`](file:///C:/Users/trcnologia/Desktop/LegoLoop/scratch/audit_studio_interactions.py)
* **Ação:** Executar a bateria de testes interativos no Chrome aberto na tela e validar que todas as interações passam sem erros no console.
* **Critério de Aceite:** Relatório factual 100% verde gerado em `docs/reports/`.

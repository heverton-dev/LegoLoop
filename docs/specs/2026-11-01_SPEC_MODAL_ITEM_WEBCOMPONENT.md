# 📐 SPEC-2026-11-01: Refatoração Canônica do Modal de Peça (<lego-item-modal>)

## 1. Contexto & Diagnóstico
Atualmente, o modal de detalhes da peça (`#modal-item`) está duplicado em dois arquivos HTML:
- `internal/studio/static/index.html` (Canvas Studio)
- `internal/studio/static/inventario.html` (Inventário Completo)

Além disso, a lógica de controle (`openItemModal`) encontra-se fragmentada em:
- `internal/studio/static/app.js` (Canvas): lógica simplificada sem filtros de compatibilidade nem tradução completa assíncrona.
- `internal/studio/static/inventory.js` (Inventário): lógica com paginação e tradução, mas sem suporte a contexto do canvas.

### Problemas Detectados:
1. **Duplicação de Código (Anti-DRY):** Manter o HTML e os estilos duplicados ou sincronizados manualmente é custoso e propenso a erros.
2. **Acessibilidade Incompleta:** Falta controle WAI-ARIA estrito para `role="tablist"`, `role="tab"`, navegação via teclado (setas, Home, End), foco inicial e fechamento seguro.
3. **Design e Hierarquia Visual:** O design precisa ser elevado ao padrão visual Linear/n8n dark mode (tokens consolidados, bordas sutis, tipografia mono refinada, diagrama IPO visual interativo, visualizador de código com numeração de linha e botão de cópia de JSON/Markdown/ID).
4. **Recursos Ausentes:** Não há inspeção explícita de schemas de I/O, telemetria/parâmetros, nem sandbox de simulação rápida de execução.

---

## 2. Objetivos & Critérios de Aceite Binários

- [ ] **CA-01 (Web Component Nativo):** Criar `internal/studio/static/components/lego-item-modal.js` exportando `<lego-item-modal>` que encapsula marcação, estado, ciclo de vida e renderização.
- [ ] **CA-02 (Eliminação de Duplicação HTML):** Substituir o bloco estático `<dialog id="modal-item">` em `index.html` e `inventario.html` pela tag canônica `<lego-item-modal id="modal-item"></lego-item-modal>`.
- [ ] **CA-03 (Unificação da API de Controle):** Expor método público `.open(item, options)` no componente, atendendo tanto ao Canvas Studio (com callback/ação de "Adicionar ao Canvas") quanto ao Inventário.
- [ ] **CA-04 (Abas WAI-ARIA & Teclado):** Navegação entre abas via setas do teclado (Esquerda/Direita), `Home`/`End`, `aria-selected`, `aria-controls` e `tabindex`.
- [ ] **CA-05 (Redesign Visual Linear/n8n):**
  - Diagrama IPO modernizado com nós visuais, conectores e status de artefatos.
  - Sub-aba ou seção de **Schema & I/O** detalhando tipos e validações de portas.
  - Sub-aba de **Simulação & Parâmetros** permitindo testar payload de entrada mock e inspecionar resposta.
  - Suporte a copiar ID, JSON da peça e Markdown com 1 clique (feedback visual instantâneo).
- [ ] **CA-06 (Compatibilidade & Resiliência):** Continua funcionando com a API `/api/assets/compatible`, `/api/assets/translate`, `/api/canonize/piece` e `/api/assets/reveal`.
- [ ] **CA-07 (Recompilação & Testes Verdes):** `go test ./...` 100% aprovado e binário compilado sem erros.

---

## 3. Arquitetura do Componente `<lego-item-modal>`

### 3.1 Ciclo de Vida e Métodos Públicos
```javascript
class LegoItemModal extends HTMLElement {
  connectedCallback()
  open(item, options = { showCanvasAction: false, onAddToCanvas: null })
  close()
}
customElements.define('lego-item-modal', LegoItemModal);
```

### 3.2 Estrutura de Abas Refatorada
1. **Visão Geral & Fluxo IPO:** Diagrama visual interativo Entrada ➔ Processamento ➔ Saída + Propriedades (O que é, O que faz, O que entrega, Artefatos).
2. **Schema & Portas de I/O:** Visualizador estruturado de campos esperados, tipagem e artefatos de entrada/saída.
3. **Peças Compatíveis:** Grafo ontológico com filtros por categoria (Skills, MCPs, Scripts, Hooks, Gates) e paginação.
4. **Ocorrências no SO:** Ocorrências em disco, status de canonização, atalhos de "Abrir Pasta" e "Canonizar".
5. **Simulação & Sandbox:** Interface rápida para rodar/simular execução da peça com mock de payload.
6. **Código & Instruções:** Visualizador com syntax styling, tradução automática PT-BR e atalho de cópia.

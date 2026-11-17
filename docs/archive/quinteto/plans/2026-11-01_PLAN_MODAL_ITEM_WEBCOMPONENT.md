# 📋 PLAN-2026-11-01: Plano de Execução do Modal Canônico <lego-item-modal>

## Fatiamento em Micro-Tickets (Mega Pipeline - Loops 3 e 4)

### 🎫 TICKET-01: Criação do Web Component `<lego-item-modal>`
* **Alvo:** `internal/studio/static/components/lego-item-modal.js`
* **Escopo:**
  * Implementar classe `LegoItemModal extends HTMLElement`.
  * Encapsular template HTML limpo e semântico com padrão de abas WAI-ARIA (`role="tablist"`, `role="tab"`, `aria-selected`, `aria-controls`).
  * Suporte a navegação por teclado: ArrowLeft / ArrowRight / Home / End nas abas.
  * Suporte a ESC para fechar e clique no backdrop.
  * Expor métodos `.open(item, options)` e `.close()`.
  * Integrar ações contextuais: "Adicionar ao Canvas", "Canonizar Peça", "Abrir no Explorer", "Copiar ID", "Copiar JSON".
* **Critério de Aceite:** Componente se registra como Custom Element e responde ao método `.open()`.

### 🎫 TICKET-02: Seções e Abas Expandidas (IPO, Schema, Compatíveis, Ocorrências, Sandbox, Código)
* **Alvo:** `internal/studio/static/components/lego-item-modal.js`
* **Escopo:**
  * **Aba 1 (Visão Geral & IPO):** Diagrama estilizado com badges, setas e cards de propriedades.
  * **Aba 2 (Schema & Portas):** Tabela de I/O inputs/outputs com tipos e tags.
  * **Aba 3 (Peças Compatíveis):** Consumo assíncrono de `/api/assets/compatible`, filtros por categoria e paginação integrada.
  * **Aba 4 (Ocorrências no SO):** Lista de caminhos com botões de revelar pasta e status do cofre.
  * **Aba 5 (Simulação & Sandbox):** Formulário mock para testar execução e inspecionar resposta gerada.
  * **Aba 6 (Código & Fonte):** Exibição do fonte traduzido ou bruto, botão de tradução PT-BR e cópia para clipboard.
* **Critério de Aceite:** Todas as 6 abas funcionam com troca dinâmica e sem quebra de estado.

### 🎫 TICKET-03: Design System & Estilização Linear/n8n Dark Mode
* **Alvo:** `internal/studio/static/style.css`
* **Escopo:**
  * Atualizar seletores de `.modal-item` e `lego-item-modal`.
  * Estilizar abas com feedback visual nítido (focus-visible, hover, active).
  * Estilizar novo visualizador de IPO com gradientes sutis e bordas modernas.
  * Estilizar visualizador de schema e sandbox de simulação.
  * Garantir responsividade e barra de rolagem estilizada (`overflow-y: auto`).
* **Critério de Aceite:** Modal com layout impecável em resoluções normais e mobile/compactas.

### 🎫 TICKET-04: Unificação e Integração em `index.html` e `inventario.html`
* **Alvo:** `internal/studio/static/index.html`, `internal/studio/static/inventario.html`, `internal/studio/static/app.js`, `internal/studio/static/inventory.js`
* **Escopo:**
  * Incluir script `<script src="components/lego-item-modal.js"></script>` em ambos os HTMLs.
  * Substituir o `<dialog id="modal-item">` estático por `<lego-item-modal id="modal-item"></lego-item-modal>`.
  * Refatorar `openItemModal` em `app.js` e `inventory.js` para delegar para `modal.open(item, options)`.
* **Critério de Aceite:** Ambos os pontos de entrada usam a mesma instância sem duplicar HTML.

### 🎫 TICKET-05: Homologação, Recompilação Go e Testes Finais
* **Alvo:** Testes Go e compilação do binário
* **Escopo:**
  * Executar `go test ./...`.
  * Recompilar `go build -o bin/legoloop.exe ./cmd/legoloop`.
  * Validar integridade sem regressões.
* **Critério de Aceite:** Testes 100% verdes e binário atualizado.

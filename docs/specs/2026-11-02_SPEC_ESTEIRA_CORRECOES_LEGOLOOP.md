# 📋 Especificação Executável (Spec) — Esteira de Correção Full-Stack LegoLoop

> **Prefixo ISO Canônico:** `2026-11-02_SPEC_ESTEIRA_CORRECOES_LEGOLOOP.md`  
> **Sessão:** `d4cfbcb7-ebf7-4f59-8d77-7c13903afc72`  
> **Status:** Proposto para Homologação (Stop-Gate SG-C01)  
> **Origem:** [`docs/intents/2026-11-02_INTENT_ESTEIRA_CORRECOES_LEGOLOOP.md`](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/2026-11-02_INTENT_ESTEIRA_CORRECOES_LEGOLOOP.md)

---

## 1. Critérios de Aceite Binários (AC-C01 a AC-C10)

### Lote 01 — Bloqueadores & Segurança (P0 & P1)
- [x] **AC-C01.1 (`BUG-RUN-01`):** `validate()` em `app.js` retorna booleano estrito (`true` ou `false`). Se o canvas estiver válido, `await validate()` avalia como truthy e abre `#modal-run`.
- [x] **AC-C01.2 (`BUG-RUN-02`):** `initRunModal()` instancia `AbortController`. Ao clicar em `#btn-modal-run-close` ou disparar o evento `cancel` do dialog, `runAbortController.abort()` é chamado, encerrando a conexão SSE.
- [x] **AC-C02.1 (`BUG-SRV-01`):** Requisição HTTP GET para `http://localhost:7878/icons.js` responde código 200 com cabeçalho `Content-Type: text/javascript`.
- [x] **AC-C03.1 (`BUG-CSP-01`):** Nenhum elemento do componente `<lego-item-modal>` possui o atributo `style="..."` inline. O console do navegador apresenta zero violações de CSP ao abrir `/inventario`.
- [x] **AC-C04.1 (`BUG-UNDO-01`):** Pressionar `Ctrl+Z` após modificar o canvas restaura o estado sem gerar `ReferenceError: renderAll is not defined` no console.
- [x] **AC-C04.2 (`BUG-UNDO-02`):** As funções `register()`, `removeNode()`, `addEdge()`, `disconnect()` e o evento `pointerup` de `startNodeDrag()` registram snapshot no histórico via `pushState()`.

### Lote 02 — Consistência de UI, Concorrência e Performance (P2)
- [ ] **AC-C05.1 (`BUG-CONFIRM-01`):** Ao clicar em "Aplicar Canonização" em `/inventario`, o corpo do diálogo exibe a mensagem de impacto explicativa (`message`) em vez de ficar vazio.
- [ ] **AC-C05.2 (`BUG-TMPL-01`):** Ao abrir o modal `#modal-templates`, os cards de template exibem o ícone SVG correspondente ou fallback neutro (`📦`), sem duplicar o nome do template lado a lado.
- [ ] **AC-C06.1 (`BUG-VAULT-01`):** Ao clicar em "Auditar Desduplicação (Dry-Run)" no modal do cofre, o container `#vault-groups-list` renderiza a lista detalhada dos grupos de arquivos redundantes encontrados.
- [ ] **AC-C06.2 (`BUG-VAULT-03`):** Se a requisição `/api/canonize/apply` falhar, o botão `#btn-vault-apply` tem seu estado redefinido para `disabled = false` no bloco `finally`.
- [ ] **AC-C07.1 (`BUG-INSP-01`):** Clicar em nós alternados no canvas dentro de um intervalo de 150ms dispara no máximo 1 requisição para `/api/suggest`, cancelando requisições anteriores via debounce e `AbortController`.
- [ ] **AC-C08.1 (`BUG-CANVAS-02`):** O arraste do canvas (pan) é limitado a uma distância máxima de $\pm 10.000\text{px}$, impedindo que nós se percam infinitamente fora da visão.
- [ ] **AC-C08.2 (`BUG-EDGE-01`):** Pressionar a tecla `Escape` durante a criação/arraste de uma aresta temporária remove a linha guia e cancela a ligação imediatamente.
- [ ] **AC-C09.1 (`BUG-PALETTE-02`):** Abrir uma categoria na paleta com mais de 50 itens (ex: Skills 773) renderiza inicialmente uma fatia controlada (janela virtual ou primeiras 40 peças + carregar mais), sem travar a thread de animação da UI.
- [ ] **AC-C09.2 (`BUG-ITEM-MODAL-01`):** O método `render()` de `lego-item-modal.js` deixa de utilizar `innerHTML` massivo e emojis hardcoded, operando com manipulação estruturada do DOM.

### Lote 03 — Ergonomia e Atalhos (P3)
- [ ] **AC-C10.1 (`BUG-KEY-01`):** Pressionar a tecla `Escape` no canvas com um nó ou aresta selecionado desseleciona o elemento ativo (`select(null)`).
- [ ] **AC-C10.2 (`BUG-KEY-02`):** Pressionar o atalho universal `Ctrl + S` ou `Cmd + S` com foco no Studio intercepta o comportamento nativo do navegador e aciona `#btn-materialize`.
- [ ] **AC-C10.3 (`BUG-PALETTE-01`):** Ao soltar uma peça arrastada no canvas, a coordenada vertical aplica compensação centralizada baseada na altura média do card (`p.y - 50`).
- [ ] **AC-C10.4 (`BUG-CANVAS-03`):** Clicar nos botões `+` e `-` de zoom utiliza a coordenada do nó selecionado ou o centro geométrico dos nós visíveis como ponto de ancoragem (pivot).
- [ ] **AC-C10.5 (`BUG-VAULT-02`):** Se `/api/canonize/status` retornar código HTTP diferente de 200 (ex: 503), a UI emite toast de aviso explicativo e não abre o modal em estado inconsistente.

---

## 2. Matriz de Roteamento Cognitivo por Ticket

| Ticket | Escopo / Defeitos | Complexidade | Modo DAG | Modelo Sugerido | Esforço |
| :--- | :--- | :---: | :---: | :---: | :---: |
| `TICK-C01` | `BUG-RUN-01`, `BUG-RUN-02` | Média | Sequencial | `pro` | High |
| `TICK-C02` | `BUG-SRV-01` | Baixa | Sequencial | `flash` | Low |
| `TICK-C03` | `BUG-CSP-01` | Média | Sequencial | `pro` | Medium |
| `TICK-C04` | `BUG-UNDO-01`, `BUG-UNDO-02` | Média | Sequencial | `pro` | Medium |
| `TICK-C05` | `BUG-CONFIRM-01`, `BUG-TMPL-01` | Baixa | Paralelo | `flash` | Low |
| `TICK-C06` | `BUG-VAULT-01`, `BUG-VAULT-03` | Média | Sequencial | `pro` | Medium |
| `TICK-C07` | `BUG-INSP-01` | Alta | Sequencial | `pro` | High |
| `TICK-C08` | `BUG-CANVAS-02`, `BUG-EDGE-01` | Média | Paralelo | `pro` | Medium |
| `TICK-C09` | `BUG-PALETTE-02`, `BUG-ITEM-MODAL-01` | Alta | Sequencial | `pro` | High |
| `TICK-C10` | `BUG-KEY-01..02`, `BUG-PALETTE-01`, `BUG-CANVAS-03`, `BUG-VAULT-02` | Média | Sequencial | `pro` | Medium |

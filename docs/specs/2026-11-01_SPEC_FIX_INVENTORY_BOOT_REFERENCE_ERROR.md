# 📋 SPEC: Correção do ReferenceError no Boot do Inventário do Studio

> **Status:** Proposta Técnica para Revisão e Validação Humana (Stop-Gate 02)  
> **Data:** 2026-11-01  
> **Referência de Intenção:** [docs/intents/2026-11-01_INTENT_FIX_INVENTORY_BOOT_REFERENCE_ERROR.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/2026-11-01_INTENT_FIX_INVENTORY_BOOT_REFERENCE_ERROR.md)  
> **Módulos Afetados:** [`internal/studio/static/inventory.js`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/inventory.js), [`internal/studio/static_test.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static_test.go)

---

## 1. Visão Geral e Arquitetura

O modal de inspeção de peças do Studio foi padronizado no Web Component autônomo `<lego-item-modal>`. O componente encapsula internamente:
* Gerenciamento de abas (`.modal-tab` / `this.switchTab`).
* Ação de fechamento (`#btn-modal-close` / `this.close()`).
* Ação de tradução (`#btn-translate` / `this.translateActiveItem()`).

O arquivo [`inventory.js`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/inventory.js) é responsável exclusivamente pela listagem, filtros, paginação, agrupamento por hash e invocação do modal via `modal.open(it, options)`. Portanto, os ouvintes residuais do modal legado presentes na função `boot()` devem ser eliminados.

---

## 2. Contratos & Modificações Fatuais

### 2.1. Saneamento de [`internal/studio/static/inventory.js`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/inventory.js)

Remover o trecho órfão em `boot()`:
```javascript
// REMOVER:
// Modal subtabs
document.querySelectorAll(".modal-tab").forEach((btn) => {
  btn.addEventListener("click", () => switchModalTab(btn.dataset.tab));
});
$("btn-modal-close").addEventListener("click", () => $("modal-item").close());
$("btn-translate").addEventListener("click", translateActiveItem);
```

### 2.2. Guardrail Automatizado em [`internal/studio/static_test.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static_test.go)

Adicionar teste de integridade estrutural `TestInventoryScriptNoOrphanModalListeners` que garanta:
* `inventory.js` não contém referências a `translateActiveItem` ou `switchModalTab`.
* O arquivo contém a delegação correta para `$("modal-item").open(it`.

---

## 3. Critérios Binários de Aceite (Gates Mecânicos)

* [ ] **G1 (Zero Referências Órfãs):** `inventory.js` não contém referências aos identificadores inexistentes `translateActiveItem` e `switchModalTab`.
* [ ] **G2 (Delegação ao Web Component Preservada):** `openItemModal(it)` continua delegando a abertura e canonização para o componente `<lego-item-modal>`.
* [ ] **G3 (Guardrail de Teste em Go):** Teste unitário em `internal/studio/static_test.go` valida a ausência de listeners órfãos no script de inventário.
* [ ] **G4 (Suíte 100% Verde):** `go test ./internal/studio/...` executa com sucesso sem falhas.

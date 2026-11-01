# 🗺️ PLAN: Correção do ReferenceError no Boot do Inventário

> **Status:** Proposta de Plano Único para Validação Humana (Stop-Gate 02)  
> **Data:** 2026-11-01  
> **Especificação Base:** [docs/specs/2026-11-01_SPEC_FIX_INVENTORY_BOOT_REFERENCE_ERROR.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-11-01_SPEC_FIX_INVENTORY_BOOT_REFERENCE_ERROR.md)  
> **Referência de Intenção:** [docs/intents/2026-11-01_INTENT_FIX_INVENTORY_BOOT_REFERENCE_ERROR.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/2026-11-01_INTENT_FIX_INVENTORY_BOOT_REFERENCE_ERROR.md)

---

## 🏛️ ÉPICO 01: Saneamento do Boot do Inventário e Guardrail de Integridade
> **Objetivo:** Eliminar o `ReferenceError` na inicialização do inventário e blindar o script contra chamadas órfãs ao modal legado.  
> **Modelo Recomendado:** `Claude 3.5 Sonnet` ou `Gemini 2.0 Flash` (`--effort medium`)

### 🔹 Micro-ticket TICK-FIX-01.1: Guardrail de Integridade Estática no Studio [CONCLUÍDO]
* **Tipo:** Micro-ticket (Teste de Regressão TDD)
* **Status:** ✅ Concluído
* **Arquivos Alvo:** [`internal/studio/static_test.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static_test.go)
* **Rota:** `Gemini 2.0 Flash` (`--effort medium`)
* **Entregáveis & Testes TDD:**
  1. Adicionar teste `TestInventoryScriptNoOrphanModalListeners`.
  2. Comprovar falha (Red) enquanto `inventory.js` contiver referências a `translateActiveItem`.

### 🔸 Nano-ticket TICK-FIX-01.2: Remoção dos Listeners Residuais em `inventory.js` [CONCLUÍDO]
* **Tipo:** Nano-ticket (Saneamento Cirúrgico)
* **Status:** ✅ Concluído
* **Arquivos Alvo:** [`internal/studio/static/inventory.js`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/inventory.js)
* **Rota:** `Gemini 2.0 Flash-Lite` (`--effort low`)
* **Entregáveis & Testes TDD:**
  1. Excluir o bloco residual de listeners do modal dentro de `boot()`.
  2. Comprovar sucesso do teste (Green).

### 🔸 Nano-ticket TICK-FIX-01.3: Homologação da Suíte do Studio e Fechamento [CONCLUÍDO]
* **Tipo:** Nano-ticket (Homologação & Relatório)
* **Status:** ✅ Concluído
* **Arquivos Alvo:** [`docs/reports/2026-11-01_REPORT_FIX_INVENTORY_BOOT_REFERENCE_ERROR.md`](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/reports/2026-11-01_REPORT_FIX_INVENTORY_BOOT_REFERENCE_ERROR.md)
* **Rota:** `Gemini 2.0 Flash-Lite` (`--effort low`)
* **Entregáveis:**
  1. Validação de `go test ./internal/studio/...`.
  2. Emissão do relatório de homologação final.

### 🛑 STOP-GATE DE AUDITORIA HUMANA FINAL
* [x] Suíte 100% verde no Studio.
* [x] Zero erros de inicialização no inventário.

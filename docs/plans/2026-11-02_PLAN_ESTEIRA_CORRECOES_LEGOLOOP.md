# 🗺️ Plano Master de Execução (Plan) — Esteira de Correção Full-Stack LegoLoop

> **Prefixo ISO Canônico:** `2026-11-02_PLAN_ESTEIRA_CORRECOES_LEGOLOOP.md`  
> **Sessão:** `d4cfbcb7-ebf7-4f59-8d77-7c13903afc72`  
> **Status:** Proposto para Homologação (Stop-Gate SG-C01)  
> **Tríade:** [`docs/intents/2026-11-02_INTENT_ESTEIRA_CORRECOES_LEGOLOOP.md`](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/2026-11-02_INTENT_ESTEIRA_CORRECOES_LEGOLOOP.md) | [`docs/specs/2026-11-02_SPEC_ESTEIRA_CORRECOES_LEGOLOOP.md`](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-11-02_SPEC_ESTEIRA_CORRECOES_LEGOLOOP.md)

---

## 1. Topologia do Grafo DAG de Execução

```mermaid
graph TD
    classDef done fill:#1b4d3e,stroke:#2ecc71,stroke-width:2px,color:#fff;
    classDef pending fill:#2c3e50,stroke:#3498db,stroke-width:2px,color:#fff;
    classDef gate fill:#e67e22,stroke:#d35400,stroke-width:2px,color:#fff;

    SG_C01["🛑 SG-C01: Aprovação da Tríade Fase C"]:::gate

    subgraph Lote01 ["Lote 01 — Bloqueadores & Segurança"]
        C01["TICK-C01: Execução & SSE (P0/P2)"]:::done
        C02["TICK-C02: Servidor Go /icons.js (P1)"]:::done
        C03["TICK-C03: Segurança CSP Modal (P1)"]:::done
        C04["TICK-C04: Undo/Redo & Mutações (P1)"]:::done
    end

    SG_L01["🛑 SG-L01: Gate Lote 01 Concluído"]:::gate

    subgraph Lote02 ["Lote 02 — Consistência de UI & Performance"]
        C05["TICK-C05: Diálogos & Templates (P2)"]:::pending
        C06["TICK-C06: Cofre Canônico & Vault (P2)"]:::pending
        C07["TICK-C07: Debounce Inspetor /api/suggest (P2)"]:::pending
        C08["TICK-C08: Limites Canvas Pan & Escape Cabo (P2)"]:::pending
        C09["TICK-C09: Virtualização Treeview & Refatoração DOM (P2)"]:::pending
    end

    SG_L02["🛑 SG-L02: Gate Lote 02 Concluído"]:::gate

    subgraph Lote03 ["Lote 03 — Ergonomia & Atalhos"]
        C10["TICK-C10: Atalhos Ctrl+S/Esc & Polimento (P3)"]:::pending
    end

    SG_FINAL["🛑 SG-C04: Homologação Final Full-Stack"]:::gate

    SG_C01 --> C01 --> C02 --> C03 --> C04 --> SG_L01
    SG_L01 --> C05
    SG_L01 --> C06
    C05 --> C07
    C06 --> C07
    C07 --> C08 --> C09 --> SG_L02
    SG_L02 --> C10 --> SG_FINAL
```

---

## 2. Dicionário de Nós & Contratos Estéreis

### Nó C05 — Diálogos & Templates (`BUG-CONFIRM-01`, `BUG-TMPL-01`)
* **Modelo / Modo:** `flash` (esforço: `low`) | Execução paralela com C06.
* **Entrada:** `internal/studio/static/inventory.js`, `internal/studio/static/app.js`.
* **Critérios de Teste:** `TestConfirmDialogMessageParam` e `TestTemplatesFallbackWithoutDoubleName`.
* **Saída:** `scratch/C05_confirm_templates.json`.

### Nó C06 — Cofre Canônico (`BUG-VAULT-01`, `BUG-VAULT-03`)
* **Modelo / Modo:** `pro` (esforço: `medium`) | Execução paralela com C05.
* **Entrada:** `internal/studio/static/app.js`.
* **Critérios de Teste:** `TestVaultRendersGroupsList` e `TestVaultApplyResetsDisabledInFinally`.
* **Saída:** `scratch/C06_vault_ui.json`.

### Nó C07 — Concorrência & Debounce (`BUG-INSP-01`)
* **Modelo / Modo:** `pro` (esforço: `high`) | Depende de C05 e C06.
* **Entrada:** `internal/studio/static/inspector.js`.
* **Critérios de Teste:** `TestInspectorSuggestDebounceAndAbort`.
* **Saída:** `scratch/C07_inspector_debounce.json`.

### Nó C08 — Limites de Canvas & Conexões (`BUG-CANVAS-02`, `BUG-EDGE-01`)
* **Modelo / Modo:** `pro` (esforço: `medium`) | Depende de C07.
* **Entrada:** `internal/studio/static/app.js`.
* **Critérios de Teste:** `TestCanvasPanClamping` e `TestEscapeCancelsLinking`.
* **Saída:** `scratch/C08_canvas_edge.json`.

### Nó C09 — Virtualização de Treeview & Refatoração DOM (`BUG-PALETTE-02`, `BUG-ITEM-MODAL-01`)
* **Modelo / Modo:** `pro` (esforço: `high`) | Depende de C08.
* **Entrada:** `internal/studio/static/palette.js`, `internal/studio/static/components/lego-item-modal.js`.
* **Critérios de Teste:** `TestPaletteTreeviewChunking` e `TestItemModalNoMassiveInnerHTML`.
* **Saída:** `scratch/C09_treeview_dom.json`.

### Nó C10 — Atalhos & Ergonomia (`BUG-KEY-01`, `BUG-KEY-02`, `BUG-PALETTE-01`, `BUG-CANVAS-03`, `BUG-VAULT-02`)
* **Modelo / Modo:** `pro` (esforço: `medium`) | Depende do Gate SG-L02.
* **Entrada:** `internal/studio/static/app.js`.
* **Critérios de Teste:** `TestKeyboardShortcutsCtrlSEsc` e `TestPaletteDropCenterOffset`.
* **Saída:** `scratch/C10_shortcuts_ergonomics.json`.

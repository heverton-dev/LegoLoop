# 🗺️ PLAN: Execução da Canonização Global do Sistema Operacional

> **Referência:** `docs/specs/SPEC_CANONIZACAO_SO.md`  
> **Status:** Pronto para Início da Fase 1 (Dry-Run)  
> **Metodologia:** Looping Engineering com Gate de Consentimento  

---

## 1. Decomposição de Tickets

| Ticket | Nível | Escopo & Ação | Critérios de Aceitação | Status |
|---|:---:|---|---|:---:|
| **TICK-CANON-01** | Micro | **Fase 1 (Dry-Run):** Rodar `legoloop canonize` em modo passivo sobre `.gemini`, `.claude` e `.agents`. Gravar telemetria em `docs/reports/RELATORIO_DRY_RUN_CANONIZACAO.json`. | CA-CANON-01, CA-CANON-02, CA-CANON-03 | ⏳ Pendente |
| **TICK-CANON-02** | Gate | **Auditoria & Aprovação:** Apresentar relatório quantitativo (arquivos escaneados, duplicatas encontradas, bytes a liberar) e aguardar o aval para aplicação real. | Gate Humano | ⏳ Pendente |
| **TICK-CANON-03** | Micro | **Fase 2 (Apply):** Executar `legoloop canonize --apply` com direcionamento ao cofre `~/.aidd/canonical`. | CA-CANON-04, CA-CANON-05 | ⏳ Pendente |
| **TICK-CANON-04** | Nano | **Validação & Commit:** Executar `go test ./...` e registrar commit semântico do relatório. | CA-CANON-06 | ⏳ Pendente |

---

## 2. Fluxo da Esteira

```mermaid
flowchart LR
    T1["TICK-CANON-01<br/>(Dry-Run & Relatório)"] --> T2{"TICK-CANON-02<br/>(Gate de Aprovação)"}
    T2 -->|Aprovado| T3["TICK-CANON-03<br/>(Apply no SO)"]
    T3 --> T4["TICK-CANON-04<br/>(Testes & Commit)"]
```

---

## 3. Próximo Passo Imediato

Iniciar a execução do **TICK-CANON-01** (Dry-Run passivo).

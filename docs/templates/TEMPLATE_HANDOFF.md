---
id: HANDOFF-{{ID}}
status: {{DRAFT | UNDER_REVIEW | APPROVED}}
date: {{AAAA-MM-DD}}
schema_version: 2
upstream: [REPORT-{{ID}}]
covers: [{{IDS_R_E_AC}}]
interpretative: [Estado, Bloqueios, Stop-Gates pendentes, Prompt de retomada]
---

# 🤝 HANDOFF: {{TITULO_DA_DEMANDA_OU_PROCESSO}}

> **Data:** {{DATA}}
> **Autor:** {{AUTOR}}
> **Template:** TEMPLATE_HANDOFF.md v2.2.0
> **Referências:** [REPORT]({{LINK_REPORT}}) · [PLAN]({{LINK_PLAN}}) · [ROUTE]({{LINK_ROUTE}})

Este arquivo recorta do REPORT o que a sessão seguinte precisa. A sessão seguinte o lê e retoma sem depender de memória de conversa.

---

## Estado

| Nó | Estado | Commit | Prova |
| :--- | :--- | :--- | :--- |
| `{{ID_DO_NO}}` | {{CONCLUÍDO | PARCIAL | NÃO INICIADO}} | `{{COMMIT_SHA}}` | {{GATE_E_SAIDA}} |

## Bloqueios

* {{Bloqueio, causa e quem o resolve.}}

## Stop-Gates pendentes

| Gate | Estado | Quem decide |
| :--- | :--- | :--- |
| `{{ID_DO_STOP_GATE}}` | {{PENDENTE | APROVADO}} | Usuário |

Nenhum agente aprova Stop-Gate.

## Prompt de retomada

```text
{{Prompt que a sessão seguinte executa para retomar: artefatos a ler, nó a executar, comando do gate.}}
```

---
schema: progresso.v1
iniciativa: "{{INICIATIVA}}"
plano: "{{PLANO}}"
atualizado_em: "{{ATUALIZADO_EM}}"
---

# Progresso Vivo: {{INICIATIVA}}

> Gerado automaticamente a partir do estado validado por Gate. Não edite à mão.

## Resumo

| Total | Aprovados (PASSED) | Em andamento (IN_PROGRESS) | Reprovados (FAILED) | Bloqueados (BLOCKED) | Pendentes (PENDING) |
|---:|---:|---:|---:|---:|---:|
| {{TOTAL}} | {{APROVADOS}} | {{EM_ANDAMENTO}} | {{REPROVADOS}} | {{BLOQUEADOS}} | {{PENDENTES}} |

## Ponto de Retomada (Em Andamento)

| Ticket | Épico | Responsável | Último registro |
|---|---|---|---|
| {{TICKET_ATIVO}} | {{EPICO_ATIVO}} | {{RESPONSAVEL}} | {{PROVA_ATUAL}} |

## Reprovados ou Bloqueados

| Ticket | Épico | Estado | Exit | Prova |
|---|---|---|---|---|
| {{TICKET_FALHA}} | {{EPICO_FALHA}} | {{ESTADO_FALHA}} | {{EXIT_FALHA}} | {{PROVA_FALHA}} |

## Pendentes

| Ticket | Épico | Título | Critérios |
|---|---|---|---|
| {{TICKET_PENDENTE}} | {{EPICO_PENDENTE}} | {{TITULO_PENDENTE}} | {{CRITERIOS_PENDENTE}} |

## Aprovados com Prova

| Ticket | Épico | Exit | Prova | Atualizado em |
|---|---|---:|---|---|
| {{TICKET_APROVADO}} | {{EPICO_APROVADO}} | 0 | {{PROVA_APROVADO}} | {{DATA_APROVADO}} |

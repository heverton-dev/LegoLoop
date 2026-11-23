---
name: legoloop-ciclo
description: Processo composto que encadeia concepção, especificação e planejamento.
consumes: Artifact:Home
produces: Artifact:Plan
artifact: docs/plans/PLAN.md
gate: file_exists('docs/plans/PLAN.md')
suggested_model: sonnet
---

# 🔄 legoloop-ciclo — Processo Composto de Concepção, Especificação e Planejamento

> **Versão:** 1.0.0
> **Motor oficial:** Composição determinística de subprocessos (`legoloop-concepcao`, `legoloop-especificacao`, `legoloop-planejamento`).
> **Regra de ouro:** execute os três subprocessos em cadeia com validação mecânica de handoff entre contratos e gates de aprovação humana.

---

## 1. Contexto & Contratos de I/O

**Quando usar:** quando for necessário executar o ciclo completo de concepção, especificação e planejamento técnico de uma demanda do início ao plano pronto para execução.
**Quando NÃO usar:** para execução de planos existentes ou para Tickets atômicos que pertençam exclusivamente a uma etapa isolada.

* **Consumes (Entradas):**
  - `Artifact:Home`: ideia inicial do usuário em linguagem natural.
* **Produces (Saídas Mecânicas):**
  - `Artifact:Plan`: plano final aprovado contendo PLAN, ROUTE e DAG formalizados.

---

## 2. Procedimento (passos determinísticos)

1. **Concepção (`legoloop-concepcao`):**
   - Entrevista socrática e alinhamento de intenção até atingir confiança estrita.
   - Gera e aprova a INTENT canônica (`docs/intents/INTENT.md`) no Stop-Gate 01.
2. **Especificação (`legoloop-especificacao`):**
   - Formalização de requisitos binários EARS/ASD-STE100 a partir da INTENT.
   - Gera e aprova a SPEC formal (`docs/specs/SPEC.md`) no Stop-Gate 02.
3. **Planejamento (`legoloop-planejamento`):**
   - Decomposição técnica em tickets e validação da cadeia de upstream.
   - Gera PLAN, ROUTE e DAG e aprova no Stop-Gate 03.

---

## 3. Invariantes & Regras Negativas (Negative Guardrails)

1. **NUNCA** execute a especificação sem a INTENT ter sido aprovada no Stop-Gate 01.
2. **NUNCA** execute o planejamento sem a SPEC ter sido aprovada no Stop-Gate 02.
3. **NUNCA** pule nenhum dos três subprocessos do ciclo canônico.

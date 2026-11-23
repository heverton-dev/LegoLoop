---
name: legoloop-planejamento
description: Use quando precisar executar o processo legoloop-planejamento: Geração formal de PLAN, ROUTE e DAG a partir de uma SPEC aprovada e tickets técnicos com trava de integridade upstream e aprovação no Stop-Gate 03 humano. NÃO use para outros Processos.
consumes: Artifact:Spec
produces: Artifact:Plan
artifact: docs/plans/PLAN.md
gate: file_exists('docs/plans/PLAN.md')
suggested_model: sonnet
---

# 📋 legoloop-planejamento — Geração Formal de Planejamento (PLAN, ROUTE e DAG)

> **Versão:** 1.0.0
> **Motor oficial:** `legoloop fabricar plano` (implementação em `internal/plangen`, `internal/routedag` e `internal/cli/fabricar_plano.go`).
> **Regra de ouro:** execute a geração formal de plano via motor oficial. A SPEC informada em `--spec` DEVE ter cadeia de aprovação válida upstream (trava `upstream_ok`).

---

## 1. Contexto & Contratos de I/O

**Quando usar:** quando houver uma SPEC aprovada no Stop-Gate 02 e tickets técnicos estruturados, e for necessário gerar atomicamente o PLAN.md, ROUTE.json e DAG.json.
**Quando NÃO usar:** para concepção/ideação (use `legoloop-concepcao`) ou especificação técnica (use `legoloop-especificacao`).

* **Consumes (Entradas):**
  - `Artifact:Spec`: documento de especificação técnica formal aprovada em `docs/specs/` (status `APPROVED` ou Stop-Gate 02 aprovado).
  - `tickets.json`: arquivo JSON contendo a lista formal de tickets cobrindo os critérios da SPEC.
* **Produces (Saídas Mecânicas):**
  - `Artifact:Plan`: documento formal `docs/plans/YYYY-MM-DD_PLAN_<TEMA>.md` e alias canônico `docs/plans/PLAN.md`.
  - `ROUTE`: arquivo `docs/routes/YYYY-MM-DD_ROUTE_<TEMA>.json`.
  - `DAG`: arquivo `docs/dags/YYYY-MM-DD_DAG_<TEMA>.json`.
  - `docs/plans/PLAN_APROVADO.md`: registro de aprovação humana no Stop-Gate 03.

---

## 2. Procedimento (passos determinísticos)

1. **Verificar Pré-requisito:** certificar que a SPEC de entrada possui cadeia upstream válida (`status: APPROVED` na INTENT e na SPEC).
2. **Executar Motor Oficial:**
   ```bash
   legoloop fabricar plano --spec docs/specs/SPEC.md --tickets docs/plans/tickets.json
   ```
   O motor valida a trava upstream, sintetiza os tickets, converte via `internal/routedag` e gera atomicamente os três artefatos.
3. **Validar Critérios Mecânicos:** conferir integridade do DAG gerado com `legoloop validate`.
4. **Stop-Gate 03 Humano:** apresentar o plano ao usuário e colher autorização explícita antes de qualquer execução de tickets.

---

## 3. Invariantes & Regras Negativas (Negative Guardrails)

1. **NUNCA** execute `legoloop fabricar plano` sem SPEC válida com upstream aprovado (a trava mecânica rejeita com exit code 2).
2. **NUNCA** adicione stubs, mocks simulados ou tickets sem critério de aceite binário (Lei #1 e Lei #5).
3. **NUNCA** inicie a execução de tickets sem o Stop-Gate 03 explicitamente aprovado pelo usuário.
4. **NUNCA** duplique lógica de conversão fora de `internal/routedag`.

---

## 4. Comandos de Verificação & Testes de Terminal

```bash
legoloop validate processos/legoloop-planejamento/*.dag.json
legoloop espinha verificar
```

---

## 5. Critério de Parada (Stopping Checklist)

- [ ] SPEC upstream validada com cadeia `status: APPROVED`.
- [ ] PLAN, ROUTE e DAG gerados com sucesso via `legoloop fabricar plano`.
- [ ] O DAG do processo passa com exit code 0 em `legoloop validate`.
- [ ] Registro da espinha verificado com exit code 0 em `legoloop espinha verificar`.
- [ ] Stop-Gate 03 de aprovação humana registrado formalmente.

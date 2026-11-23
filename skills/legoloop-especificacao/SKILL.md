---
name: legoloop-especificacao
description: Use quando precisar executar o processo legoloop-especificacao: Geração formal da SPEC a partir de uma INTENT aprovada com trava de integridade upstream, validação de critérios binários de aceite (EARS/ASD-STE100) e aprovação no Stop-Gate 02 humano. NÃO use para outros Processos.
consumes: Artifact:Intent
produces: Artifact:Spec
artifact: docs/specs/SPEC.md
gate: file_exists('docs/specs/SPEC.md')
suggested_model: sonnet
---

# 📋 legoloop-especificacao — Geração Formal de Especificação Técnica (SPEC)

> **Versão:** 1.0.0
> **Motor oficial:** `legoloop fabricar spec` (implementação em `internal/specgen` e `internal/cli/fabricar_spec.go`).
> **Regra de ouro:** execute a geração formal de SPEC via motor oficial. A INTENT informada em `--intent` DEVE estar com status `APPROVED` (trava `upstream_ok`).

---

## 1. Contexto & Contratos de I/O

**Quando usar:** quando houver uma INTENT aprovada no Stop-Gate 01 e for necessário formalizar a arquitetura técnica, requisitos binários (EARS/ASD-STE100), análise FMEA e critérios de aceite na SPEC.
**Quando NÃO usar:** para concepção/ideação (use `legoloop-concepcao`) ou planejamento em tickets/DAG (use `legoloop-planejamento`).

* **Consumes (Entradas):**
  - `Artifact:Intent`: arquivo de INTENT aprovado em `docs/intents/` (status `APPROVED`).
* **Produces (Saídas Mecânicas):**
  - `Artifact:Spec`: documento formal `docs/specs/YYYY-MM-DD_SPEC_<TEMA>.md` e alias canônico `docs/specs/SPEC.md`.
  - `docs/specs/SPEC_APROVADA.md`: registro de aprovação humana no Stop-Gate 02.

---

## 2. Procedimento (passos determinísticos)

1. **Verificar Pré-requisito:** certificar que a INTENT de entrada em `docs/intents/` possui `status: APPROVED` no cabeçalho formal.
2. **Executar Motor Oficial:**
   ```bash
   legoloop fabricar spec --intent docs/intents/INTENT.md
   ```
   O motor valida a trava upstream, sintetiza os critérios binários e grava o documento formal em `docs/specs/`.
3. **Validar Critérios Mecânicos:** conferir que o arquivo gerado atende aos padrões ASD-STE100 e possui critérios de aceite binários verificáveis.
4. **Stop-Gate 02 Humano:** apresentar a proposta técnica ao usuário e colher autorização explícita antes de qualquer etapa de planejamento (`legoloop-planejamento`).

---

## 3. Invariantes & Regras Negativas (Negative Guardrails)

1. **NUNCA** execute `legoloop fabricar spec` com INTENT em `UNDER_REVIEW`, `DRAFT` ou `REJECTED` (a trava mecânica rejeita com exit code 2).
2. **NUNCA** adicione stubs, mocks simulados ou pendências `TODO` na especificação gerada (Lei #1 e Lei #16).
3. **NUNCA** avance para o planejamento ou execução sem o Stop-Gate 02 explicitamente aprovado pelo usuário.
4. **NUNCA** altere o documento da SPEC gerado para contornar a validação upstream.

---

## 4. Comandos de Verificação & Testes de Terminal

```bash
legoloop validate processos/legoloop-especificacao/*.dag.json
legoloop espinha verificar
```

---

## 5. Critério de Parada (Stopping Checklist)

- [ ] INTENT upstream validada com `status: APPROVED`.
- [ ] SPEC gerada com sucesso via `legoloop fabricar spec --intent ...`.
- [ ] O DAG do processo passa com exit code 0 em `legoloop validate`.
- [ ] Registro da espinha verificado com exit code 0 em `legoloop espinha verificar`.
- [ ] Stop-Gate 02 de aprovação humana registrado formalmente.
